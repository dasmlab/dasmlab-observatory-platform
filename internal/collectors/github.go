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

type githubCollector struct {
	tenant string
	buf    []collector.Event
}

func (c *githubCollector) Name() string { return "github" }

func (c *githubCollector) Discover(ctx context.Context) error { return ctx.Err() }

func (c *githubCollector) Health(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if token := strings.TrimSpace(os.Getenv("GITHUB_TOKEN")); token == "" && os.Getenv("DPO_REQUIRE_LIVE_CREDS") == "1" {
		return errMissing("GITHUB_TOKEN")
	}
	return nil
}

func (c *githubCollector) Collect(ctx context.Context) error {
	token := strings.TrimSpace(os.Getenv("GITHUB_TOKEN"))
	repos := parseRepoAllowlist()
	if token == "" || len(repos) == 0 {
		// Demo fallback
		c.buf = []collector.Event{evt(c.tenant, "github", "repo", "github_stars", 42, map[string]string{"mode": "demo"})}
		return ctx.Err()
	}

	client := &http.Client{Timeout: 20 * time.Second}
	var events []collector.Event
	var totalStars, totalForks float64
	for _, full := range repos {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://api.github.com/repos/"+full, nil)
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
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
		_ = resp.Body.Close()
		if resp.StatusCode == http.StatusNotFound {
			continue
		}
		if resp.StatusCode >= 300 {
			return fmt.Errorf("github %s: %s", full, resp.Status)
		}
		var data struct {
			FullName        string `json:"full_name"`
			StargazersCount int    `json:"stargazers_count"`
			ForksCount      int    `json:"forks_count"`
			OpenIssuesCount int    `json:"open_issues_count"`
		}
		if err := json.Unmarshal(body, &data); err != nil {
			return err
		}
		totalStars += float64(data.StargazersCount)
		totalForks += float64(data.ForksCount)
		dims := map[string]string{"mode": "live", "repository": data.FullName}
		events = append(events,
			evt(c.tenant, "github", "repo", "github_stars", float64(data.StargazersCount), dims),
			evt(c.tenant, "github", "repo", "github_forks", float64(data.ForksCount), dims),
			evt(c.tenant, "github", "repo", "github_open_issues", float64(data.OpenIssuesCount), dims),
		)
	}
	events = append(events, evt(c.tenant, "github", "rollup", "github_stars", totalStars, map[string]string{"mode": "live", "scope": "allowlist"}))
	_ = totalForks
	c.buf = events
	return ctx.Err()
}

func (c *githubCollector) Normalize(ctx context.Context) ([]collector.Event, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	out := c.buf
	c.buf = nil
	return out, nil
}

func parseRepoAllowlist() []string {
	raw := strings.TrimSpace(os.Getenv("GITHUB_REPOS"))
	if raw == "" {
		raw = "dasmlab/dasmlab-observatory-platform,lmcdasm/cheapcloud,lmcdasm/dasmlab_home"
	}
	var out []string
	for _, p := range strings.Split(raw, ",") {
		p = strings.TrimSpace(p)
		if p != "" && strings.Contains(p, "/") {
			out = append(out, p)
		}
	}
	return out
}

func evt(tenant, collectorName, typ, metric string, value float64, dims map[string]string) collector.Event {
	return collector.Event{
		SchemaVersion: "1",
		Collector:     collectorName,
		Type:          typ,
		Timestamp:     time.Now().UTC(),
		Tenant:        tenant,
		Entity:        tenant,
		Metric:        metric,
		Value:         value,
		Unit:          "count",
		Dims:          dims,
	}
}
