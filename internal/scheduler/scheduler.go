package scheduler

import (
	"context"
	"log"
	"time"

	"github.com/lmcdasm/dasmlab-observatory-platform/internal/score"
	"github.com/lmcdasm/dasmlab-observatory-platform/internal/store"
	collector "github.com/lmcdasm/dasmlab-observatory-platform/platform/collector-sdk"
)

type Scheduler struct {
	reg    *collector.Registry
	st     *store.Store
	engine *score.Engine
	every  time.Duration
}

func New(reg *collector.Registry, st *store.Store, engine *score.Engine) *Scheduler {
	return &Scheduler{reg: reg, st: st, engine: engine, every: time.Hour}
}

func (s *Scheduler) Start() {
	go func() {
		t := time.NewTicker(s.every)
		defer t.Stop()
		s.RunOnce()
		for range t.C {
			s.RunOnce()
		}
	}()
}

func (s *Scheduler) RunOnce() {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	for _, c := range s.reg.List() {
		st := store.CollectorStatus{Name: c.Name(), LastRun: time.Now().UTC().Format(time.RFC3339)}
		if err := c.Health(ctx); err != nil {
			st.Healthy = false
			st.LastError = err.Error()
			st.Message = "unhealthy"
			_ = s.st.UpsertCollectorStatus(st)
			continue
		}
		if err := c.Discover(ctx); err != nil {
			st.Healthy = false
			st.LastError = err.Error()
			_ = s.st.UpsertCollectorStatus(st)
			continue
		}
		if err := c.Collect(ctx); err != nil {
			st.Healthy = false
			st.LastError = err.Error()
			_ = s.st.UpsertCollectorStatus(st)
			continue
		}
		events, err := c.Normalize(ctx)
		if err != nil {
			st.Healthy = false
			st.LastError = err.Error()
			_ = s.st.UpsertCollectorStatus(st)
			continue
		}
		if err := s.st.InsertEvents(events); err != nil {
			st.Healthy = false
			st.LastError = err.Error()
			_ = s.st.UpsertCollectorStatus(st)
			continue
		}
		st.Healthy = true
		st.Message = "ok"
		_ = s.st.UpsertCollectorStatus(st)
	}
	if err := s.engine.Recompute(); err != nil {
		log.Printf("score recompute: %v", err)
	}
}
