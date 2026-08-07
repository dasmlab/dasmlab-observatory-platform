package api

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"github.com/lmcdasm/dasmlab-observatory-platform/internal/scheduler"
	"github.com/lmcdasm/dasmlab-observatory-platform/internal/score"
	"github.com/lmcdasm/dasmlab-observatory-platform/internal/store"
	collector "github.com/lmcdasm/dasmlab-observatory-platform/platform/collector-sdk"
)

type Deps struct {
	Store     *store.Store
	Registry  *collector.Registry
	Engine    *score.Engine
	Scheduler *scheduler.Scheduler
	Tenant    string
	Version   string
	StaticDir string
}

type Server struct {
	d Deps
}

func New(d Deps) *Server { return &Server{d: d} }

func (s *Server) Handler() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.RequestID, middleware.RealIP, middleware.Logger, middleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders: []string{"Accept", "Authorization", "Content-Type"},
	}))

	r.Get("/healthz", s.healthz)
	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/score", s.score)
		r.Get("/score/history", s.scoreHistory)
		r.Get("/sources/status", s.sources)
		r.Get("/meta", s.meta)
		r.Get("/engineering", s.engineering)
		r.Post("/collect/run", s.runCollect)
	})

	if s.d.StaticDir != "" {
		if _, err := os.Stat(s.d.StaticDir); err == nil {
			fileServer(r, s.d.StaticDir)
		}
	}
	return r
}

func (s *Server) healthz(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{
		"status":  "ok",
		"service": "dpo-api",
		"version": s.d.Version,
		"time":    time.Now().UTC().Format(time.RFC3339),
	})
}

func (s *Server) meta(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{
		"tenant":  s.d.Tenant,
		"version": s.d.Version,
		"product": "dpo",
		"platform": "dop",
		"five_questions": []string{
			"What exists?", "What changed?", "Why did it change?", "What will happen?", "What should I do?",
		},
	})
}

func (s *Server) score(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		name = "overall"
	}
	snap, err := s.d.Store.LatestScore(s.d.Tenant, name)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	if snap == nil {
		writeJSON(w, http.StatusOK, map[string]any{"tenant": s.d.Tenant, "name": name, "value": nil})
		return
	}
	writeJSON(w, http.StatusOK, snap)
}

func (s *Server) scoreHistory(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		name = "overall"
	}
	hist, err := s.d.Store.ScoreHistory(s.d.Tenant, name, 30)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	writeJSON(w, http.StatusOK, hist)
}

func (s *Server) sources(w http.ResponseWriter, _ *http.Request) {
	sts, err := s.d.Store.ListCollectorStatus()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	// Ensure registered collectors appear even before first run.
	known := map[string]bool{}
	for _, st := range sts {
		known[st.Name] = true
	}
	for _, c := range s.d.Registry.List() {
		if !known[c.Name()] {
			sts = append(sts, store.CollectorStatus{Name: c.Name(), Healthy: true, Message: "registered"})
		}
	}
	writeJSON(w, http.StatusOK, sts)
}

func (s *Server) runCollect(w http.ResponseWriter, _ *http.Request) {
	go s.d.Scheduler.RunOnce()
	writeJSON(w, http.StatusAccepted, map[string]string{"status": "started"})
}

func (s *Server) engineering(w http.ResponseWriter, _ *http.Request) {
	metrics, err := s.d.Store.LatestMetrics(s.d.Tenant, time.Now().Add(-48*time.Hour))
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"tenant":  s.d.Tenant,
		"metrics": metrics,
		"bots": map[string]any{
			"googlebot_fetches": metrics["googlebot_fetches"],
			"bot_hits":          metrics["bot_hits"],
		},
		"index": map[string]any{
			"sitemap_urls":          metrics["sitemap_urls"],
			"sitemap_freshness_pct": metrics["sitemap_freshness_pct"],
			"tech_health":           metrics["tech_health"],
		},
		"authority": map[string]any{
			"github_stars":       metrics["github_stars"],
			"github_forks":       metrics["github_forks"],
			"github_open_issues": metrics["github_open_issues"],
		},
		"engagement": map[string]any{
			"engaged_sessions": metrics["engaged_sessions"],
			"page_views":       metrics["page_views"],
			"activity_events":  metrics["activity_events"],
		},
	})
}

func writeJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(v)
}

func fileServer(r chi.Router, dir string) {
	r.Get("/*", func(w http.ResponseWriter, req *http.Request) {
		path := strings.TrimPrefix(req.URL.Path, "/")
		if path == "" {
			path = "index.html"
		}
		full := filepath.Join(dir, path)
		if !strings.HasPrefix(full, filepath.Clean(dir)+string(os.PathSeparator)) && full != filepath.Clean(dir) {
			http.NotFound(w, req)
			return
		}
		if st, err := os.Stat(full); err == nil && !st.IsDir() {
			http.ServeFile(w, req, full)
			return
		}
		// SPA fallback
		http.ServeFile(w, req, filepath.Join(dir, "index.html"))
	})
}
