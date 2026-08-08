package collectors

import (
	"context"
	"os"
	"strings"

	collector "github.com/dasmlab/dasmlab-observatory-platform/platform/collector-sdk"
)

type daoCollector struct {
	tenant  string
	dataDir string
	buf     []collector.Event
}

func (c *daoCollector) Name() string                       { return "dao" }
func (c *daoCollector) Discover(ctx context.Context) error { return ctx.Err() }
func (c *daoCollector) Health(ctx context.Context) error   { return ctx.Err() }
func (c *daoCollector) Normalize(ctx context.Context) ([]collector.Event, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	out := c.buf
	c.buf = nil
	return out, nil
}

func (c *daoCollector) Collect(ctx context.Context) error {
	path := envOr("DAO_PROMPTS_PATH", "/app/docs/PROMPTS.md")
	raw, err := os.ReadFile(path)
	mode := "live"
	promptScore := 60.0
	if err != nil {
		// Fallback: known pack size from repo (3 draft prompts).
		mode = "demo"
		promptScore = 55
	} else {
		text := string(raw)
		lines := 0
		for _, line := range strings.Split(text, "\n") {
			line = strings.TrimSpace(line)
			if len(line) > 2 && (line[0] >= '1' && line[0] <= '9') && (line[1] == '.' || line[1] == ')') {
				lines++
			}
		}
		// Effectiveness: coverage of numbered prompts + section headers.
		headers := strings.Count(text, "##")
		promptScore = clamp(30+float64(lines)*15+float64(headers)*5, 0, 100)
	}

	// AI trust: blend bot crawl signals from sibling collectors when present in process
	// (read via env stubs / conservative baseline; store blend happens in DUO).
	trust := 50.0
	if mode == "live" {
		trust = 40 + promptScore*0.4
	} else {
		trust = 58
	}
	trust = clamp(trust, 0, 100)

	dims := map[string]string{"mode": mode, "prompts_path": path}
	c.buf = []collector.Event{
		evt(c.tenant, "dao", "score", "prompt_effectiveness", round1(promptScore), dims),
		evt(c.tenant, "dao", "score", "ai_trust_score", round1(trust), dims),
	}
	return ctx.Err()
}
