package collectors

import (
	"context"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	collector "github.com/lmcdasm/dasmlab-observatory-platform/platform/collector-sdk"
)

// siteCollector probes website sitemap/robots for technical trust + freshness.
type siteCollector struct {
	tenant string
	buf    []collector.Event
}

func (c *siteCollector) Name() string { return "website" }

func (c *siteCollector) Discover(ctx context.Context) error { return ctx.Err() }

func (c *siteCollector) Health(ctx context.Context) error { return ctx.Err() }

func (c *siteCollector) Collect(ctx context.Context) error {
	base := strings.TrimRight(strings.TrimSpace(os.Getenv("DPO_SITE_BASE")), "/")
	if base == "" {
		base = "https://dasmlab.org"
	}
	client := &http.Client{Timeout: 15 * time.Second, CheckRedirect: func(req *http.Request, via []*http.Request) error {
		if len(via) >= 5 {
			return http.ErrUseLastResponse
		}
		return nil
	}}

	tech := 100.0
	fresh := 50.0
	var urls float64
	dims := map[string]string{"mode": "live", "site": base}

	sitemapURL := base + "/sitemap.xml"
	scode, sbody, err := fetch(ctx, client, sitemapURL)
	if err != nil || scode >= 400 || !looksLikeSitemap(sbody) {
		tech -= 25
		dims["sitemap"] = "missing"
	} else {
		dims["sitemap"] = "ok"
		urls = float64(strings.Count(string(sbody), "<loc>"))
		if urls > 0 {
			fresh = 90
		}
		if urls > 20 {
			fresh = 98
		}
	}

	rcode, rbody, err := fetch(ctx, client, base+"/robots.txt")
	if err != nil || rcode >= 400 || looksLikeHTML(rbody) {
		tech -= 15
		dims["robots"] = "missing"
	} else {
		dims["robots"] = "ok"
		// Prefer Sitemap: directives from robots when /sitemap.xml is SPA-fallback.
		for _, line := range strings.Split(string(rbody), "\n") {
			line = strings.TrimSpace(line)
			if len(line) >= 8 && strings.EqualFold(line[:8], "sitemap:") {
				u := strings.TrimSpace(line[8:])
				if u == "" {
					continue
				}
				c2, b2, e2 := fetch(ctx, client, u)
				if e2 == nil && c2 < 400 && looksLikeSitemap(b2) {
					dims["sitemap"] = "ok"
					urls = float64(strings.Count(string(b2), "<loc>"))
					if urls > 0 {
						fresh = 90
					}
					if urls > 20 {
						fresh = 98
					}
					if tech < 100 {
						tech += 25
						if tech > 100 {
							tech = 100
						}
					}
				}
			}
		}
	}

	hcode, _, err := fetch(ctx, client, base+"/")
	if err != nil || hcode >= 500 {
		tech -= 30
		dims["home"] = "error"
	} else if hcode >= 400 {
		tech -= 10
		dims["home"] = "client_error"
	} else {
		dims["home"] = "ok"
	}
	if tech < 0 {
		tech = 0
	}

	c.buf = []collector.Event{
		evt(c.tenant, "website", "site", "tech_health", tech, dims),
		evt(c.tenant, "website", "site", "sitemap_freshness_pct", fresh, dims),
		evt(c.tenant, "website", "site", "sitemap_urls", urls, dims),
	}
	return ctx.Err()
}

func (c *siteCollector) Normalize(ctx context.Context) ([]collector.Event, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	out := c.buf
	c.buf = nil
	return out, nil
}

func fetch(ctx context.Context, client *http.Client, url string) (int, []byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return 0, nil, err
	}
	req.Header.Set("User-Agent", "DASMLAB-DPO/0.1 (+https://dasmlab.org)")
	resp, err := client.Do(req)
	if err != nil {
		return 0, nil, err
	}
	defer resp.Body.Close()
	b, _ := io.ReadAll(io.LimitReader(resp.Body, 2<<20))
	return resp.StatusCode, b, nil
}

func looksLikeHTML(b []byte) bool {
	s := strings.ToLower(string(b))
	return strings.Contains(s, "<html") || strings.Contains(s, "<!doctype html")
}

func looksLikeSitemap(b []byte) bool {
	if looksLikeHTML(b) {
		return false
	}
	s := string(b)
	return strings.Contains(s, "<urlset") || strings.Contains(s, "<sitemapindex") || strings.Contains(s, "<loc>")
}
