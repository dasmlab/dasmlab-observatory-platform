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

	collector "github.com/dasmlab/dasmlab-observatory-platform/platform/collector-sdk"
)

type daopsCollector struct {
	tenant string
	buf    []collector.Event
}

func (c *daopsCollector) Name() string                       { return "daops" }
func (c *daopsCollector) Discover(ctx context.Context) error { return ctx.Err() }
func (c *daopsCollector) Health(ctx context.Context) error   { return ctx.Err() }
func (c *daopsCollector) Normalize(ctx context.Context) ([]collector.Event, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	out := c.buf
	c.buf = nil
	return out, nil
}

func (c *daopsCollector) Collect(ctx context.Context) error {
	token := strings.TrimSpace(os.Getenv("GITHUB_TOKEN"))
	repo := envOr("DAOPS_GITHUB_REPO", "dasmlab/dasmlab-observatory-platform")
	if token == "" {
		c.buf = []collector.Event{
			evt(c.tenant, "daops", "score", "delivery_confidence", 65, map[string]string{"mode": "demo"}),
			evt(c.tenant, "daops", "score", "toil_ratio", 35, map[string]string{"mode": "demo"}),
		}
		return ctx.Err()
	}

	client := &http.Client{Timeout: 25 * time.Second}
	url := fmt.Sprintf("https://api.github.com/repos/%s/actions/runs?per_page=30", repo)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	body, _ := io.ReadAll(io.LimitReader(resp.Body, 2<<20))
	_ = resp.Body.Close()
	if resp.StatusCode >= 300 {
		c.buf = []collector.Event{
			evt(c.tenant, "daops", "score", "delivery_confidence", 50, map[string]string{"mode": "demo", "http": resp.Status}),
			evt(c.tenant, "daops", "score", "toil_ratio", 40, map[string]string{"mode": "demo"}),
		}
		return ctx.Err()
	}
	var payload struct {
		WorkflowRuns []struct {
			Conclusion string `json:"conclusion"`
			Status     string `json:"status"`
			Event      string `json:"event"`
			Name       string `json:"name"`
		} `json:"workflow_runs"`
	}
	if err := json.Unmarshal(body, &payload); err != nil {
		return err
	}
	total, success, dispatch := 0, 0, 0
	for _, r := range payload.WorkflowRuns {
		if r.Status != "completed" {
			continue
		}
		total++
		if r.Conclusion == "success" {
			success++
		}
		if r.Event == "workflow_dispatch" {
			dispatch++
		}
	}
	delivery := 65.0
	toil := 35.0
	mode := "live"
	if total == 0 {
		mode = "demo"
	} else {
		delivery = 100 * float64(success) / float64(total)
		toil = 100 * float64(dispatch) / float64(total)
	}
	dims := map[string]string{"mode": mode, "repo": repo}
	c.buf = []collector.Event{
		evt(c.tenant, "daops", "score", "delivery_confidence", round1(delivery), dims),
		evt(c.tenant, "daops", "score", "toil_ratio", round1(toil), dims),
	}
	return ctx.Err()
}
