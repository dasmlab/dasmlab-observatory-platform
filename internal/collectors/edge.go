package collectors

import (
	"bufio"
	"context"
	"os"
	"path/filepath"
	"strings"

	collector "github.com/dasmlab/dasmlab-observatory-platform/platform/collector-sdk"
)

// 2026 crawl / AI bot UA tokens (first-party edge log analysis).
var botTokens = []struct {
	Token  string
	Metric string
	Label  string
}{
	{"Googlebot", "googlebot_fetches", "Googlebot"},
	{"GoogleOther", "googleother_fetches", "GoogleOther"},
	{"bingbot", "bingbot_fetches", "Bingbot"},
	{"GPTBot", "gptbot_fetches", "GPTBot"},
	{"OAI-SearchBot", "oai_searchbot_fetches", "OAI-SearchBot"},
	{"ChatGPT-User", "chatgpt_user_fetches", "ChatGPT-User"},
	{"ClaudeBot", "claudebot_fetches", "ClaudeBot"},
	{"Claude-SearchBot", "claude_searchbot_fetches", "Claude-SearchBot"},
	{"PerplexityBot", "perplexitybot_fetches", "PerplexityBot"},
	{"Perplexity-User", "perplexity_user_fetches", "Perplexity-User"},
}

type edgeCollector struct {
	tenant  string
	logPath string
	buf     []collector.Event
}

func (c *edgeCollector) Name() string { return "edge-logs" }

func (c *edgeCollector) Discover(ctx context.Context) error { return ctx.Err() }

func (c *edgeCollector) Health(ctx context.Context) error { return ctx.Err() }

func (c *edgeCollector) Collect(ctx context.Context) error {
	path := strings.TrimSpace(os.Getenv("DPO_EDGE_LOG_PATH"))
	if path == "" {
		path = c.logPath
	}
	_ = os.MkdirAll(filepath.Dir(path), 0o755)

	counts := map[string]float64{}
	pathHits := map[string]map[string]float64{} // path -> bot -> count
	mode := "demo"

	f, err := os.Open(path)
	if err != nil {
		// Demo synthetic until HAProxy sample is pulled.
		c.buf = []collector.Event{
			evt(c.tenant, "edge-logs", "crawl", "googlebot_fetches", 42, map[string]string{"bot": "Googlebot", "mode": "demo"}),
			evt(c.tenant, "edge-logs", "crawl", "gptbot_fetches", 3, map[string]string{"bot": "GPTBot", "mode": "demo"}),
		}
		return ctx.Err()
	}
	defer f.Close()
	mode = "live"
	sc := bufio.NewScanner(f)
	sc.Buffer(make([]byte, 0, 64*1024), 1024*1024)
	for sc.Scan() {
		line := sc.Text()
		ua := extractUA(line)
		reqPath := extractPath(line)
		for _, b := range botTokens {
			if strings.Contains(ua, b.Token) {
				counts[b.Metric]++
				if reqPath != "" {
					if pathHits[reqPath] == nil {
						pathHits[reqPath] = map[string]float64{}
					}
					pathHits[reqPath][b.Label]++
				}
				break
			}
		}
	}

	var events []collector.Event
	for _, b := range botTokens {
		if v := counts[b.Metric]; v > 0 {
			events = append(events, evt(c.tenant, "edge-logs", "crawl", b.Metric, v, map[string]string{"bot": b.Label, "mode": mode}))
		}
	}
	if len(events) == 0 {
		events = append(events, evt(c.tenant, "edge-logs", "crawl", "googlebot_fetches", 0, map[string]string{"bot": "Googlebot", "mode": mode}))
	}
	n := 0
	for p, bots := range pathHits {
		for bot, v := range bots {
			if n >= 50 {
				break
			}
			events = append(events, evt(c.tenant, "edge-logs", "crawl_path", "edge_path_hits", v, map[string]string{
				"mode": mode, "path": p, "bot": bot,
			}))
			n++
		}
	}
	c.buf = events
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

func extractUA(line string) string {
	// Prefer last quoted field (common httplog / combined).
	parts := strings.Split(line, `"`)
	if len(parts) >= 6 {
		return parts[len(parts)-2]
	}
	if len(parts) >= 2 {
		return parts[len(parts)-2]
	}
	return line
}

func extractPath(line string) string {
	// Look for "METHOD /path HTTP
	i := strings.Index(line, `"`)
	if i < 0 {
		return ""
	}
	rest := line[i+1:]
	j := strings.Index(rest, `"`)
	if j < 0 {
		return ""
	}
	req := rest[:j]
	fields := strings.Fields(req)
	if len(fields) >= 2 {
		p := fields[1]
		if q := strings.IndexByte(p, '?'); q >= 0 {
			p = p[:q]
		}
		return p
	}
	return ""
}
