package score

import (
	"math"
	"time"

	"github.com/lmcdasm/dasmlab-observatory-platform/internal/store"
)

// Weights locked in ADR-0400 (GEO neutral 50 until P2).
var weights = map[string]float64{
	"seo":        0.25,
	"geo":        0.25,
	"authority":  0.20,
	"engagement": 0.15,
	"freshness":  0.10,
	"trust":      0.05,
}

type Engine struct {
	st     *store.Store
	tenant string
}

func NewEngine(st *store.Store, tenant string) *Engine {
	return &Engine{st: st, tenant: tenant}
}

func (e *Engine) Recompute() error {
	comps := map[string]float64{
		"seo":        e.mapMetric(70, "gsc_impressions", func(v float64) float64 { return 20 + v*0.01 }),
		"geo":        50, // placeholder until AI collectors
		"authority":  e.mapMetric(60, "github_stars", func(v float64) float64 { return 45 + v*8 }),
		"engagement": e.mapMetric(65, "engaged_sessions", func(v float64) float64 { return 40 + v*3 }),
		"freshness":  clamp(e.metricOr(80, "sitemap_freshness_pct", 1)),
		"trust":      clamp(e.metricOr(90, "tech_health", 1)),
	}
	var overall float64
	for k, w := range weights {
		overall += comps[k] * w
	}
	overall = math.Round(overall*10) / 10
	return e.st.SaveScore(store.ScoreSnapshot{
		Tenant:     e.tenant,
		Name:       "overall",
		Value:      overall,
		Components: comps,
		CreatedAt:  time.Now().UTC(),
	})
}

func (e *Engine) metricOr(fallback float64, metric string, scale float64) float64 {
	v, err := e.st.MetricAvg(e.tenant, metric, time.Now().Add(-7*24*time.Hour))
	if err != nil {
		return fallback
	}
	return clamp(v * scale)
}

func (e *Engine) mapMetric(fallback float64, metric string, fn func(float64) float64) float64 {
	v, err := e.st.MetricAvg(e.tenant, metric, time.Now().Add(-7*24*time.Hour))
	if err != nil {
		return fallback
	}
	return clamp(fn(v))
}

func clamp(v float64) float64 {
	if v < 0 {
		return 0
	}
	if v > 100 {
		return 100
	}
	return math.Round(v*10) / 10
}
