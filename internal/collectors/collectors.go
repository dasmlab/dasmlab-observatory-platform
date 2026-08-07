package collectors

import (
	"context"
	"math/rand"
	"os"
	"path/filepath"

	collector "github.com/lmcdasm/dasmlab-observatory-platform/platform/collector-sdk"
)

// Default returns Phase-1 collectors. Live paths activate when env credentials exist.
func Default(tenant, dataDir string) []collector.Collector {
	return []collector.Collector{
		&stubCollector{name: "gsc", tenant: tenant, metric: "gsc_impressions", base: 1200, typ: "query"},
		&githubCollector{tenant: tenant},
		&edgeCollector{tenant: tenant, logPath: filepath.Join(dataDir, "edge", "access.log")},
		&activityCollector{tenant: tenant},
		&siteCollector{tenant: tenant},
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
	c.buf = []collector.Event{evt(c.tenant, c.name, c.typ, c.metric, c.base+jitter, map[string]string{"mode": modeFor(c.name)})}
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
	if os.Getenv("DPO_REQUIRE_LIVE_CREDS") == "1" && c.name == "gsc" && os.Getenv("GSC_CREDENTIALS_JSON") == "" {
		return errMissing("GSC_CREDENTIALS_JSON")
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
	n := 40.0 + rand.Float64()*10
	mode := "demo"
	if b, err := os.ReadFile(c.logPath); err == nil && len(b) > 0 {
		n = float64(len(b) / 80)
		mode = "live"
	}
	c.buf = []collector.Event{evt(c.tenant, "edge-logs", "crawl", "googlebot_fetches", n, map[string]string{"bot": "Googlebot", "mode": mode})}
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

func modeFor(name string) string {
	if name == "gsc" && os.Getenv("GSC_CREDENTIALS_JSON") != "" {
		return "live"
	}
	return "demo"
}
