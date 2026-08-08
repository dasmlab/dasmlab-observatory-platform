package collectors

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	collector "github.com/dasmlab/dasmlab-observatory-platform/platform/collector-sdk"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

// gscCollector queries Search Console searchAnalytics when GSC_CREDENTIALS_JSON is set.
type gscCollector struct {
	tenant string
	buf    []collector.Event
}

func (c *gscCollector) Name() string { return "gsc" }

func (c *gscCollector) Discover(ctx context.Context) error { return ctx.Err() }

func (c *gscCollector) Health(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if strings.TrimSpace(os.Getenv("GSC_CREDENTIALS_JSON")) == "" {
		if os.Getenv("DPO_REQUIRE_LIVE_CREDS") == "1" {
			return errMissing("GSC_CREDENTIALS_JSON")
		}
		// Demo mode is healthy but not live.
		return nil
	}
	return nil
}

func (c *gscCollector) Collect(ctx context.Context) error {
	raw := strings.TrimSpace(os.Getenv("GSC_CREDENTIALS_JSON"))
	site := strings.TrimSpace(os.Getenv("GSC_SITE_URL"))
	if site == "" {
		// Domain property in Search Console (not URL-prefix https://…).
		site = "sc-domain:dasmlab.org"
	}
	if raw == "" {
		c.buf = []collector.Event{
			evt(c.tenant, "gsc", "query", "gsc_impressions", 1200, map[string]string{"mode": "demo"}),
			evt(c.tenant, "gsc", "query", "gsc_clicks", 40, map[string]string{"mode": "demo"}),
			evt(c.tenant, "gsc", "query", "gsc_ctr", 0.033, map[string]string{"mode": "demo"}),
			evt(c.tenant, "gsc", "query", "gsc_position", 18, map[string]string{"mode": "demo"}),
		}
		return ctx.Err()
	}

	ts, err := google.CredentialsFromJSON(ctx, []byte(raw), "https://www.googleapis.com/auth/webmasters.readonly")
	if err != nil {
		return fmt.Errorf("gsc credentials: %w", err)
	}
	client := oauth2.NewClient(ctx, ts.TokenSource)
	end := time.Now().UTC().AddDate(0, 0, -2) // GSC lag
	start := end.AddDate(0, 0, -7)

	// Totals
	total, err := gscQuery(ctx, client, site, start, end, nil, 1)
	if err != nil {
		return err
	}
	var impressions, clicks, ctr, position float64
	if len(total) > 0 {
		impressions = total[0].Impressions
		clicks = total[0].Clicks
		ctr = total[0].CTR
		position = total[0].Position
	}
	dims := map[string]string{"mode": "live", "site": site, "window": "7d"}
	events := []collector.Event{
		evt(c.tenant, "gsc", "query", "gsc_impressions", impressions, dims),
		evt(c.tenant, "gsc", "query", "gsc_clicks", clicks, dims),
		evt(c.tenant, "gsc", "query", "gsc_ctr", ctr, dims),
		evt(c.tenant, "gsc", "query", "gsc_position", position, dims),
	}

	// Top pages for content spine
	byPage, err := gscQuery(ctx, client, site, start, end, []string{"page"}, 25)
	if err != nil {
		return err
	}
	for _, row := range byPage {
		page := ""
		if len(row.Keys) > 0 {
			page = row.Keys[0]
		}
		pd := map[string]string{"mode": "live", "page": page, "path": pathFromURL(page)}
		events = append(events,
			evt(c.tenant, "gsc", "page", "gsc_page_impressions", row.Impressions, pd),
			evt(c.tenant, "gsc", "page", "gsc_page_clicks", row.Clicks, pd),
		)
	}

	// Top queries
	byQuery, err := gscQuery(ctx, client, site, start, end, []string{"query"}, 15)
	if err != nil {
		return err
	}
	for _, row := range byQuery {
		q := ""
		if len(row.Keys) > 0 {
			q = row.Keys[0]
		}
		qd := map[string]string{"mode": "live", "query": q}
		events = append(events, evt(c.tenant, "gsc", "query", "gsc_query_impressions", row.Impressions, qd))
	}

	c.buf = events
	return ctx.Err()
}

func (c *gscCollector) Normalize(ctx context.Context) ([]collector.Event, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	out := c.buf
	c.buf = nil
	return out, nil
}

type gscRow struct {
	Keys        []string
	Clicks      float64
	Impressions float64
	CTR         float64
	Position    float64
}

func gscQuery(ctx context.Context, client *http.Client, site string, start, end time.Time, dimensions []string, rowLimit int) ([]gscRow, error) {
	body := map[string]any{
		"startDate":  start.Format("2006-01-02"),
		"endDate":    end.Format("2006-01-02"),
		"rowLimit":   rowLimit,
		"dimensions": dimensions,
	}
	b, _ := json.Marshal(body)
	url := "https://www.googleapis.com/webmasters/v3/sites/" + urlPathEscape(site) + "/searchAnalytics/query"
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(b))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	raw, _ := io.ReadAll(io.LimitReader(resp.Body, 4<<20))
	if resp.StatusCode >= 300 {
		return nil, fmt.Errorf("gsc query %s: %s", resp.Status, truncate(string(raw), 200))
	}
	var parsed struct {
		Rows []struct {
			Keys        []string `json:"keys"`
			Clicks      float64  `json:"clicks"`
			Impressions float64  `json:"impressions"`
			CTR         float64  `json:"ctr"`
			Position    float64  `json:"position"`
		} `json:"rows"`
	}
	if err := json.Unmarshal(raw, &parsed); err != nil {
		return nil, err
	}
	out := make([]gscRow, 0, len(parsed.Rows))
	for _, r := range parsed.Rows {
		out = append(out, gscRow{Keys: r.Keys, Clicks: r.Clicks, Impressions: r.Impressions, CTR: r.CTR, Position: r.Position})
	}
	return out, nil
}

func urlPathEscape(site string) string {
	return strings.ReplaceAll(strings.ReplaceAll(site, ":", "%3A"), "/", "%2F")
}

func pathFromURL(u string) string {
	u = strings.TrimPrefix(u, "https://")
	u = strings.TrimPrefix(u, "http://")
	if i := strings.IndexByte(u, '/'); i >= 0 {
		return u[i:]
	}
	return "/"
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "…"
}
