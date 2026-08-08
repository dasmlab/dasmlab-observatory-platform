package content

import (
	"context"
	"encoding/xml"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/dasmlab/dasmlab-observatory-platform/internal/store"
)

type Spine struct {
	st     *store.Store
	tenant string
}

func New(st *store.Store, tenant string) *Spine {
	return &Spine{st: st, tenant: tenant}
}

func (s *Spine) RefreshFromSitemap(ctx context.Context) error {
	base := strings.TrimRight(os.Getenv("DPO_SITE_BASE"), "/")
	if base == "" {
		base = "https://dasmlab.org"
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, base+"/sitemap.xml", nil)
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", "DASMLAB-DPO/0.1")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(io.LimitReader(resp.Body, 2<<20))
	if resp.StatusCode >= 400 || !strings.Contains(string(body), "<urlset") {
		return nil // site collector already flags; spine waits for real sitemap
	}
	var doc struct {
		URLs []struct {
			Loc string `xml:"loc"`
		} `xml:"url"`
	}
	if err := xml.Unmarshal(body, &doc); err != nil {
		return err
	}
	day := time.Now().UTC().Format("2006-01-02")
	views, _ := s.st.SumMetricByDim(s.tenant, "path_page_views", "path", time.Now().Add(-7*24*time.Hour))
	// Activity may store path in dims from surfing — also try engaged
	engaged, _ := s.st.SumMetricByDim(s.tenant, "engaged_sessions", "path", time.Now().Add(-7*24*time.Hour))
	impr, _ := s.st.SumMetricByDim(s.tenant, "gsc_page_impressions", "path", time.Now().Add(-7*24*time.Hour))
	clicks, _ := s.st.SumMetricByDim(s.tenant, "gsc_page_clicks", "path", time.Now().Add(-7*24*time.Hour))
	bots, _ := s.st.SumMetricByDim(s.tenant, "edge_path_hits", "path", time.Now().Add(-7*24*time.Hour))

	for _, u := range doc.URLs {
		p := pathOnly(u.Loc)
		if p == "" {
			continue
		}
		_ = s.st.UpsertEntity(s.tenant, "path", p)
		_ = s.st.UpsertPathDaily(s.tenant, p, day, views[p], engaged[p], impr[p], clicks[p], bots[p])
	}
	return nil
}

func pathOnly(raw string) string {
	u, err := url.Parse(raw)
	if err != nil {
		return ""
	}
	if u.Path == "" {
		return "/"
	}
	return u.Path
}
