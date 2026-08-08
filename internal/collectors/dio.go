package collectors

import (
	"context"
	"crypto/tls"
	"net/http"
	"strings"
	"time"

	"github.com/dasmlab/dasmlab-observatory-platform/internal/kube"
	collector "github.com/dasmlab/dasmlab-observatory-platform/platform/collector-sdk"
)

type dioCollector struct {
	tenant string
	buf    []collector.Event
}

func (c *dioCollector) Name() string                       { return "dio" }
func (c *dioCollector) Discover(ctx context.Context) error { return ctx.Err() }
func (c *dioCollector) Health(ctx context.Context) error   { return ctx.Err() }
func (c *dioCollector) Normalize(ctx context.Context) ([]collector.Event, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	out := c.buf
	c.buf = nil
	return out, nil
}

func (c *dioCollector) Collect(ctx context.Context) error {
	ns := envOr("DPO_WATCH_NS", "dpo-system")
	capacity := 75.0
	failover := 68.0
	mode := "demo"

	if client, err := kube.InCluster(); err == nil {
		var pvc struct {
			Status struct {
				Phase      string   `json:"phase"`
				Capacity   map[string]string `json:"capacity"`
				Conditions []struct {
					Type   string `json:"type"`
					Status string `json:"status"`
				} `json:"conditions"`
			} `json:"status"`
			Spec struct {
				Resources struct {
					Requests map[string]string `json:"requests"`
				} `json:"resources"`
			} `json:"spec"`
		}
		path := "/api/v1/namespaces/" + ns + "/persistentvolumeclaims/dpo-data"
		if err := client.GetJSON(ctx, path, &pvc); err == nil {
			mode = "live"
			capacity = 50
			if pvc.Status.Phase == "Bound" {
				capacity = 85
			}
			if _, ok := pvc.Status.Capacity["storage"]; ok {
				capacity = 90
			}
		}
	}

	// Failover readiness: both cluster Routes/health respond.
	urls := []string{
		envOr("DPO_PUBLIC_URL", "https://dpo-dasmlab.apps.2026-prod-1.ocp.dasmlab.org"),
		envOr("DPO_FAILOVER_URL", "https://dpo-dasmlab.apps.2026-prod-2-1.ocp.dasmlab.org"),
	}
	httpClient := &http.Client{
		Timeout: 8 * time.Second,
		Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}},
	}
	ok := 0
	for _, u := range urls {
		u = strings.TrimRight(strings.TrimSpace(u), "/")
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, u+"/healthz", nil)
		if err != nil {
			continue
		}
		resp, err := httpClient.Do(req)
		if err != nil {
			continue
		}
		_ = resp.Body.Close()
		if resp.StatusCode < 500 {
			ok++
		}
	}
	if ok > 0 {
		if mode != "live" {
			mode = "live"
		}
		failover = 50 * float64(ok) // 50 or 100
	}

	dims := map[string]string{"mode": mode, "namespace": ns}
	c.buf = []collector.Event{
		evt(c.tenant, "dio", "score", "capacity_confidence", round1(capacity), dims),
		evt(c.tenant, "dio", "score", "failover_readiness", round1(failover), dims),
	}
	return ctx.Err()
}
