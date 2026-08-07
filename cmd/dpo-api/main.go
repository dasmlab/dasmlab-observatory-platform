package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/lmcdasm/dasmlab-observatory-platform/internal/api"
	"github.com/lmcdasm/dasmlab-observatory-platform/internal/collectors"
	"github.com/lmcdasm/dasmlab-observatory-platform/internal/scheduler"
	"github.com/lmcdasm/dasmlab-observatory-platform/internal/score"
	"github.com/lmcdasm/dasmlab-observatory-platform/internal/store"
	collector "github.com/lmcdasm/dasmlab-observatory-platform/platform/collector-sdk"
)

func main() {
	addr := env("DPO_LISTEN", ":8080")
	dataDir := env("DPO_DATA_DIR", "/data")
	tenant := env("DPO_TENANT", "dasmlab.org")
	staticDir := env("DPO_STATIC_DIR", "/app/web")
	version := env("DPO_BUILD_VERSION", "dev")

	if err := os.MkdirAll(dataDir, 0o755); err != nil {
		log.Fatalf("data dir: %v", err)
	}

	st, err := store.Open(filepath.Join(dataDir, "dpo.sqlite"))
	if err != nil {
		log.Fatalf("store: %v", err)
	}
	defer st.Close()

	reg := collector.NewRegistry()
	for _, c := range collectors.Default(tenant, dataDir) {
		reg.Register(c)
	}

	engine := score.NewEngine(st, tenant)
	sched := scheduler.New(reg, st, engine)
	sched.Start()

	// Seed a score on boot so the dashboard is never empty.
	if err := engine.Recompute(); err != nil {
		log.Printf("initial score: %v", err)
	}

	srv := api.New(api.Deps{
		Store:     st,
		Registry:  reg,
		Engine:    engine,
		Scheduler: sched,
		Tenant:    tenant,
		Version:   version,
		StaticDir: staticDir,
	})

	log.Printf("dpo-api listening on %s tenant=%s version=%s", addr, tenant, version)
	if err := http.ListenAndServe(addr, srv.Handler()); err != nil {
		log.Fatal(err)
	}
}

func env(k, def string) string {
	if v := strings.TrimSpace(os.Getenv(k)); v != "" {
		return v
	}
	return def
}
