package collectors

import (
	"context"
	"strconv"

	"github.com/dasmlab/dasmlab-observatory-platform/internal/kube"
	collector "github.com/dasmlab/dasmlab-observatory-platform/platform/collector-sdk"
)

// dcoCollector — Deploy Confidence + Operational Complexity.
type dcoCollector struct {
	tenant string
	buf    []collector.Event
}

func (c *dcoCollector) Name() string                       { return "dco" }
func (c *dcoCollector) Discover(ctx context.Context) error { return ctx.Err() }
func (c *dcoCollector) Health(ctx context.Context) error   { return ctx.Err() }
func (c *dcoCollector) Normalize(ctx context.Context) ([]collector.Event, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	out := c.buf
	c.buf = nil
	return out, nil
}

func (c *dcoCollector) Collect(ctx context.Context) error {
	ns := envOr("DPO_WATCH_NS", "dpo-system")
	client, err := kube.InCluster()
	if err != nil {
		c.buf = []collector.Event{
			evt(c.tenant, "dco", "score", "deploy_confidence", 72, map[string]string{"mode": "demo", "reason": "not_in_cluster"}),
			evt(c.tenant, "dco", "score", "operational_complexity", 41, map[string]string{"mode": "demo", "reason": "not_in_cluster"}),
		}
		return ctx.Err()
	}

	mode := "live"
	deployConf := 50.0
	var deploy struct {
		Status struct {
			Replicas            int32 `json:"replicas"`
			ReadyReplicas       int32 `json:"readyReplicas"`
			UnavailableReplicas int32 `json:"unavailableReplicas"`
			UpdatedReplicas     int32 `json:"updatedReplicas"`
		} `json:"status"`
		Spec struct {
			Replicas *int32 `json:"replicas"`
		} `json:"spec"`
	}
	path := "/apis/apps/v1/namespaces/" + ns + "/deployments/dpo"
	if err := client.GetJSON(ctx, path, &deploy); err != nil {
		mode = "demo"
		deployConf = 55
	} else {
		desired := int32(1)
		if deploy.Spec.Replicas != nil {
			desired = *deploy.Spec.Replicas
		}
		ready := deploy.Status.ReadyReplicas
		updated := deploy.Status.UpdatedReplicas
		unavail := deploy.Status.UnavailableReplicas
		deployConf = 100
		if desired > 0 {
			deployConf = 40 + 50*float64(ready)/float64(desired) + 10*float64(updated)/float64(desired)
		}
		if unavail > 0 {
			deployConf -= 15 * float64(unavail)
		}
		deployConf = clamp(deployConf, 0, 100)
	}

	counts := 0
	liveCounts := 0
	resources := []string{
		"/apis/apps/v1/namespaces/" + ns + "/deployments",
		"/api/v1/namespaces/" + ns + "/persistentvolumeclaims",
		"/api/v1/namespaces/" + ns + "/secrets",
		"/apis/batch/v1/namespaces/" + ns + "/cronjobs",
		"/apis/route.openshift.io/v1/namespaces/" + ns + "/routes",
	}
	for _, p := range resources {
		n, err := client.CountResource(ctx, p)
		if err != nil {
			continue
		}
		liveCounts++
		counts += n
	}
	complexity := float64(counts)
	if liveCounts == 0 {
		mode = "demo"
		complexity = 41
	} else {
		complexity = float64(counts) * 8
		if complexity > 100 {
			complexity = 100
		}
	}

	dims := map[string]string{"mode": mode, "namespace": ns, "objects": strconv.Itoa(counts)}
	c.buf = []collector.Event{
		evt(c.tenant, "dco", "score", "deploy_confidence", round1(deployConf), dims),
		evt(c.tenant, "dco", "score", "operational_complexity", round1(complexity), dims),
	}
	return ctx.Err()
}
