package collectors

import (
	"context"
	"crypto/tls"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/dasmlab/dasmlab-observatory-platform/internal/kube"
	collector "github.com/dasmlab/dasmlab-observatory-platform/platform/collector-sdk"
)

type dnoCollector struct {
	tenant string
	buf    []collector.Event
}

func (c *dnoCollector) Name() string                       { return "dno" }
func (c *dnoCollector) Discover(ctx context.Context) error { return ctx.Err() }
func (c *dnoCollector) Health(ctx context.Context) error   { return ctx.Err() }
func (c *dnoCollector) Normalize(ctx context.Context) ([]collector.Event, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	out := c.buf
	c.buf = nil
	return out, nil
}

func (c *dnoCollector) Collect(ctx context.Context) error {
	targets := []string{
		envOr("DPO_SITE_BASE", "https://dasmlab.org"),
		envOr("DPO_PUBLIC_URL", "https://dpo-dasmlab.apps.2026-prod-1.ocp.dasmlab.org"),
		envOr("DPO_FAILOVER_URL", "https://dpo-dasmlab.apps.2026-prod-2-1.ocp.dasmlab.org"),
	}
	client := &http.Client{
		Timeout: 8 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via) >= 3 {
				return http.ErrUseLastResponse
			}
			return nil
		},
	}
	ok, tried := 0, 0
	for _, u := range targets {
		u = strings.TrimRight(strings.TrimSpace(u), "/")
		if u == "" {
			continue
		}
		tried++
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, u+"/healthz", nil)
		if err != nil {
			req, err = http.NewRequestWithContext(ctx, http.MethodHead, u+"/", nil)
		}
		if err != nil {
			continue
		}
		// Site base may not have /healthz — try both.
		resp, err := client.Do(req)
		if err != nil || resp.StatusCode >= 500 {
			req2, _ := http.NewRequestWithContext(ctx, http.MethodGet, u+"/", nil)
			if req2 != nil {
				resp2, err2 := client.Do(req2)
				if err2 == nil {
					if resp2.StatusCode < 500 {
						ok++
					}
					_ = resp2.Body.Close()
				}
			}
			if resp != nil {
				_ = resp.Body.Close()
			}
			continue
		}
		_ = resp.Body.Close()
		if resp.StatusCode < 500 {
			ok++
		}
	}
	reach := 50.0
	mode := "demo"
	if tried > 0 {
		mode = "live"
		reach = 100 * float64(ok) / float64(tried)
	}

	intent := 60.0
	intentMode := mode
	ns := envOr("DPO_WATCH_NS", "dpo-system")
	if kc, err := kube.InCluster(); err == nil {
		var route struct {
			Spec struct {
				Host string `json:"host"`
				TLS  *struct {
					Termination string `json:"termination"`
				} `json:"tls"`
			} `json:"spec"`
		}
		if err := kc.GetJSON(ctx, "/apis/route.openshift.io/v1/namespaces/"+ns+"/routes/dpo", &route); err == nil {
			intentMode = "live"
			intent = 40
			wantHost := os.Getenv("DPO_CLUSTER")
			if route.Spec.Host != "" {
				intent += 30
			}
			if route.Spec.TLS != nil && route.Spec.TLS.Termination != "" {
				intent += 30
			}
			_ = wantHost
		}
	}

	dims := map[string]string{"mode": mode}
	c.buf = []collector.Event{
		evt(c.tenant, "dno", "score", "service_reachability", round1(reach), dims),
		evt(c.tenant, "dno", "score", "intent_compliance", round1(intent), map[string]string{"mode": intentMode}),
	}
	return ctx.Err()
}
