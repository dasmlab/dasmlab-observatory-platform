package collectors

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	collector "github.com/lmcdasm/dasmlab-observatory-platform/platform/collector-sdk"
)

// activityCollector pulls Surfing Activity via machine token, or falls back to demo.
type activityCollector struct {
	tenant string
	buf    []collector.Event
}

func (c *activityCollector) Name() string { return "activity" }

func (c *activityCollector) Discover(ctx context.Context) error { return ctx.Err() }

func (c *activityCollector) Health(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	url := strings.TrimSpace(os.Getenv("DPO_ACTIVITY_URL"))
	if url != "" && strings.TrimSpace(os.Getenv("DPO_ACTIVITY_TOKEN")) == "" && os.Getenv("DPO_REQUIRE_LIVE_CREDS") == "1" {
		return errMissing("DPO_ACTIVITY_TOKEN")
	}
	return nil
}

func (c *activityCollector) Collect(ctx context.Context) error {
	url := strings.TrimSpace(os.Getenv("DPO_ACTIVITY_URL"))
	token := strings.TrimSpace(os.Getenv("DPO_ACTIVITY_TOKEN"))
	if url == "" {
		c.buf = []collector.Event{evt(c.tenant, "activity", "journey", "engaged_sessions", 18, map[string]string{"mode": "demo"})}
		return ctx.Err()
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(io.LimitReader(resp.Body, 4<<20))
	if resp.StatusCode >= 300 {
		return fmt.Errorf("activity bridge %s", resp.Status)
	}
	var payload struct {
		Events []struct {
			Type      string `json:"type"`
			Path      string `json:"path"`
			EngagedMs int64  `json:"engagedMs"`
			Bot       bool   `json:"bot"`
		} `json:"events"`
	}
	if err := json.Unmarshal(body, &payload); err != nil {
		return err
	}
	var pages, engaged, bots float64
	pathViews := map[string]float64{}
	for _, e := range payload.Events {
		if e.Bot {
			bots++
			continue
		}
		switch e.Type {
		case "page", "navigate":
			pages++
			if e.Path != "" {
				pathViews[e.Path]++
			}
		case "engaged":
			engaged++
		}
	}
	dims := map[string]string{"mode": "live", "source": "surfing"}
	c.buf = []collector.Event{
		evt(c.tenant, "activity", "journey", "engaged_sessions", engaged, dims),
		evt(c.tenant, "activity", "journey", "page_views", pages, dims),
		evt(c.tenant, "activity", "journey", "bot_hits", bots, dims),
		evt(c.tenant, "activity", "journey", "activity_events", float64(len(payload.Events)), dims),
	}
	for p, v := range pathViews {
		c.buf = append(c.buf, evt(c.tenant, "activity", "journey", "path_page_views", v, map[string]string{
			"mode": "live", "source": "surfing", "path": p,
		}))
	}
	return ctx.Err()
}

func (c *activityCollector) Normalize(ctx context.Context) ([]collector.Event, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	out := c.buf
	c.buf = nil
	return out, nil
}
