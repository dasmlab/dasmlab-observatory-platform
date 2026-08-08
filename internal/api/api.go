package api

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"github.com/dasmlab/dasmlab-observatory-platform/internal/auth"
	"github.com/dasmlab/dasmlab-observatory-platform/internal/content"
	"github.com/dasmlab/dasmlab-observatory-platform/internal/duo"
	"github.com/dasmlab/dasmlab-observatory-platform/internal/family"
	"github.com/dasmlab/dasmlab-observatory-platform/internal/scheduler"
	"github.com/dasmlab/dasmlab-observatory-platform/internal/score"
	"github.com/dasmlab/dasmlab-observatory-platform/internal/store"
	collector "github.com/dasmlab/dasmlab-observatory-platform/platform/collector-sdk"
)

type Deps struct {
	Store     *store.Store
	Registry  *collector.Registry
	Engine    *score.Engine
	Scheduler *scheduler.Scheduler
	Spine     *content.Spine
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
	r.Use(auth.Middleware)
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
		r.Get("/family", s.family)
		r.Get("/duo/impact", s.duoImpact)
		r.Get("/duo/recommend", s.duoRecommend)
		r.Get("/engineering", s.engineering)
		r.Get("/content", s.content)
		r.Get("/baselines", s.listBaselines)
		r.Post("/baseline", s.createBaseline)
		r.Get("/baseline/diff", s.baselineDiff)
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


func (s *Server) family(w http.ResponseWriter, _ *http.Request) {
	c := family.Default()
	c.ActiveProduct = "dpo"
	writeJSON(w, http.StatusOK, c)
}

func (s *Server) duoImpact(w http.ResponseWriter, _ *http.Request) {
	overall := 55.0
	live := false
	if snap, err := s.d.Store.LatestScore(s.d.Tenant, "overall"); err == nil && snap != nil {
		overall = snap.Value
		live = true
	}
	writeJSON(w, http.StatusOK, duo.Compose(s.d.Tenant, overall, live))
}

func (s *Server) duoRecommend(w http.ResponseWriter, _ *http.Request) {
	overall := 55.0
	if snap, err := s.d.Store.LatestScore(s.d.Tenant, "overall"); err == nil && snap != nil {
		overall = snap.Value
	}
	writeJSON(w, http.StatusOK, duo.Recommend(s.d.Tenant, overall))
}

func (s *Server) meta(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{
		"tenant":  s.d.Tenant,
		"version": s.d.Version,
		"product": "dpo",
		"platform": "dop",
		"family":   "/api/v1/family",
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
			"googlebot_fetches":    metrics["googlebot_fetches"],
			"gptbot_fetches":       metrics["gptbot_fetches"],
			"claudebot_fetches":    metrics["claudebot_fetches"],
			"perplexitybot_fetches": metrics["perplexitybot_fetches"],
			"bingbot_fetches":      metrics["bingbot_fetches"],
			"bot_hits":             metrics["bot_hits"],
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
		"search": map[string]any{
			"gsc_impressions": metrics["gsc_impressions"],
			"gsc_clicks":      metrics["gsc_clicks"],
			"gsc_ctr":         metrics["gsc_ctr"],
			"gsc_position":    metrics["gsc_position"],
		},
	})
}

func (s *Server) content(w http.ResponseWriter, _ *http.Request) {
	if s.d.Spine != nil {
		_ = s.d.Spine.RefreshFromSitemap(context.Background())
	}
	paths, err := s.d.Store.TopPaths(s.d.Tenant, 30)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	ents, _ := s.d.Store.ListEntities(s.d.Tenant, "path")
	writeJSON(w, http.StatusOK, map[string]any{
		"tenant":       s.d.Tenant,
		"paths":        paths,
		"sitemap_urls": ents,
	})
}

func (s *Server) createBaseline(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Label string `json:"label"`
	}
	_ = json.NewDecoder(r.Body).Decode(&body)
	if body.Label == "" {
		body.Label = "snapshot-" + time.Now().UTC().Format("20060102-1504")
	}
	score, _ := s.d.Store.LatestScore(s.d.Tenant, "overall")
	paths, _ := s.d.Store.TopPaths(s.d.Tenant, 20)
	metrics, _ := s.d.Store.LatestMetrics(s.d.Tenant, time.Now().Add(-48*time.Hour))
	sts, _ := s.d.Store.ListCollectorStatus()
	payload := map[string]any{
		"score":      score,
		"paths":      paths,
		"metrics":    metrics,
		"collectors": sts,
	}
	if err := s.d.Store.SaveBaseline(s.d.Tenant, body.Label, payload); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	writeJSON(w, http.StatusCreated, map[string]any{"label": body.Label, "payload": payload})
}

func (s *Server) listBaselines(w http.ResponseWriter, _ *http.Request) {
	list, err := s.d.Store.ListBaselines(s.d.Tenant)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	writeJSON(w, http.StatusOK, list)
}

func (s *Server) baselineDiff(w http.ResponseWriter, r *http.Request) {
	from := r.URL.Query().Get("from")
	to := r.URL.Query().Get("to")
	if from == "" || to == "" {
		http.Error(w, "from and to labels required", 400)
		return
	}
	a, err := s.d.Store.GetBaseline(s.d.Tenant, from)
	if err != nil {
		http.Error(w, "from baseline: "+err.Error(), 404)
		return
	}
	b, err := s.d.Store.GetBaseline(s.d.Tenant, to)
	if err != nil {
		http.Error(w, "to baseline: "+err.Error(), 404)
		return
	}
	var pa, pb map[string]any
	_ = json.Unmarshal(a.Payload, &pa)
	_ = json.Unmarshal(b.Payload, &pb)
	writeJSON(w, http.StatusOK, map[string]any{
		"from": a,
		"to":   b,
		"score_delta": map[string]any{
			"from": pa["score"],
			"to":   pb["score"],
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
