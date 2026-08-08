package collectors

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/dasmlab/dasmlab-observatory-platform/internal/kube"
	collector "github.com/dasmlab/dasmlab-observatory-platform/platform/collector-sdk"
)

type dsoCollector struct {
	tenant string
	buf    []collector.Event
}

func (c *dsoCollector) Name() string                         { return "dso" }
func (c *dsoCollector) Discover(ctx context.Context) error   { return ctx.Err() }
func (c *dsoCollector) Health(ctx context.Context) error     { return ctx.Err() }
func (c *dsoCollector) Normalize(ctx context.Context) ([]collector.Event, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	out := c.buf
	c.buf = nil
	return out, nil
}

func (c *dsoCollector) Collect(ctx context.Context) error {
	ns := envOr("DPO_WATCH_NS", "dpo-system")
	client, err := kube.InCluster()
	if err != nil {
		c.buf = []collector.Event{
			evt(c.tenant, "dso", "score", "attack_surface_evolution", 38, map[string]string{"mode": "demo"}),
			evt(c.tenant, "dso", "score", "secrets_hygiene", 55, map[string]string{"mode": "demo"}),
		}
		return ctx.Err()
	}

	mode := "live"
	routes, errR := client.CountResource(ctx, "/apis/route.openshift.io/v1/namespaces/"+ns+"/routes")
	svcs, errS := client.CountResource(ctx, "/api/v1/namespaces/"+ns+"/services")
	if errR != nil && errS != nil {
		mode = "demo"
		routes, svcs = 1, 1
	}
	// Surface index: more public exposure → higher evolution pressure (novel, not CVE list).
	surface := float64(routes*25 + svcs*10)
	if surface > 100 {
		surface = 100
	}

	// Secrets hygiene: required vs optional keys on dpo-secrets (keys only).
	var secret struct {
		Data map[string]json.RawMessage `json:"data"`
	}
	hyg := 40.0
	secPath := "/api/v1/namespaces/" + ns + "/secrets/dpo-secrets"
	if err := client.GetJSON(ctx, secPath, &secret); err != nil {
		if mode == "live" {
			hyg = 35
		} else {
			hyg = 55
		}
	} else {
		required := []string{"GITHUB_TOKEN", "ACTIVITY_MACHINE_TOKEN"}
		optional := []string{"GSC_CREDENTIALS_JSON"}
		presentReq, presentOpt := 0, 0
		for _, k := range required {
			if _, ok := secret.Data[k]; ok {
				presentReq++
			}
		}
		for _, k := range optional {
			if _, ok := secret.Data[k]; ok {
				presentOpt++
			}
		}
		hyg = 50*float64(presentReq)/float64(len(required)) + 50*float64(presentOpt)/float64(len(optional))
	}

	dims := map[string]string{
		"mode": mode, "namespace": ns,
		"routes": strconv.Itoa(routes), "services": strconv.Itoa(svcs),
	}
	c.buf = []collector.Event{
		evt(c.tenant, "dso", "score", "attack_surface_evolution", round1(surface), dims),
		evt(c.tenant, "dso", "score", "secrets_hygiene", round1(hyg), dims),
	}
	return ctx.Err()
}
