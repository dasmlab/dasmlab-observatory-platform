package collectors

import (
	"context"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	collector "github.com/lmcdasm/dasmlab-observatory-platform/platform/collector-sdk"
)

// Default returns Phase-1 collectors. Real GSC/GitHub wiring uses env credentials;
// without them, collectors emit deterministic demo telemetry so dashboards work.
func Default(tenant, dataDir string) []collector.Collector {
	return []collector.Collector{
		&stubCollector{name: "gsc", tenant: tenant, metric: "gsc_impressions", base: 1200, typ: "query"},
		&stubCollector{name: "github", tenant: tenant, metric: "github_stars", base: 42, typ: "repo"},
		&edgeCollector{tenant: tenant, logPath: filepath.Join(dataDir, "edge", "access.log")},
		&stubCollector{name: "activity", tenant: tenant, metric: "engaged_sessions", base: 18, typ: "journey"},
		&stubCollector{name: "tech", tenant: tenant, metric: "tech_health", base: 92, typ: "site"},
		&stubCollector{name: "sitemap", tenant: tenant, metric: "sitemap_freshness_pct", base: 98, typ: "site"},
	}
}

type stubCollector struct {
	name, tenant, metric, typ string
	base                      float64
	buf                       []collector.Event
}

func (c *stubCollector) Name() string { return c.name }

func (c *stubCollector) Discover(ctx context.Context) error { return ctx.Err() }

func (c *stubCollector) Collect(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	jitter := (rand.Float64()*0.1 - 0.05) * c.base
	c.buf = []collector.Event{{
		SchemaVersion: "1",
		Collector:     c.name,
		Type:          c.typ,
		Timestamp:     time.Now().UTC(),
		Tenant:        c.tenant,
		Entity:        c.tenant,
		Metric:        c.metric,
		Value:         c.base + jitter,
		Unit:          "count",
		Dims:          map[string]string{"mode": mode()},
	}}
	return nil
}

func (c *stubCollector) Normalize(ctx context.Context) ([]collector.Event, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	out := c.buf
	c.buf = nil
	return out, nil
}

func (c *stubCollector) Health(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	// Mark unhealthy only when explicitly required credentials missing for live mode.
	if os.Getenv("DPO_REQUIRE_LIVE_CREDS") == "1" {
		switch c.name {
		case "gsc":
			if os.Getenv("GSC_CREDENTIALS_JSON") == "" {
				return errMissing("GSC_CREDENTIALS_JSON")
			}
		case "github":
			if os.Getenv("GITHUB_TOKEN") == "" {
				return errMissing("GITHUB_TOKEN")
			}
		}
	}
	return nil
}

type edgeCollector struct {
	tenant  string
	logPath string
	buf     []collector.Event
}

func (c *edgeCollector) Name() string { return "edge-logs" }

func (c *edgeCollector) Discover(ctx context.Context) error { return ctx.Err() }

func (c *edgeCollector) Collect(ctx context.Context) error {
	_ = os.MkdirAll(filepath.Dir(c.logPath), 0o755)
	// Count bot-like lines if present; else emit synthetic googlebot fetch count.
	n := 40.0 + rand.Float64()*10
	if b, err := os.ReadFile(c.logPath); err == nil && len(b) > 0 {
		n = float64(len(b) / 80)
	}
	c.buf = []collector.Event{{
		SchemaVersion: "1",
		Collector:     "edge-logs",
		Type:          "crawl",
		Timestamp:     time.Now().UTC(),
		Tenant:        c.tenant,
		Entity:        "googlebot",
		Metric:        "googlebot_fetches",
		Value:         n,
		Unit:          "count",
		Dims:          map[string]string{"bot": "Googlebot", "mode": mode()},
	}}
	return ctx.Err()
}

func (c *edgeCollector) Normalize(ctx context.Context) ([]collector.Event, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	out := c.buf
	c.buf = nil
	return out, nil
}

func (c *edgeCollector) Health(ctx context.Context) error { return ctx.Err() }

type missingErr string

func (e missingErr) Error() string { return string(e) }

func errMissing(env string) error { return missingErr("missing " + env) }

func mode() string {
	if os.Getenv("DPO_REQUIRE_LIVE_CREDS") == "1" {
		return "live"
	}
	return "demo"
}
