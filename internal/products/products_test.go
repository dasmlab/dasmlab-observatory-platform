package products_test

import (
	"testing"

	"github.com/dasmlab/dasmlab-observatory-platform/internal/duo"
	"github.com/dasmlab/dasmlab-observatory-platform/internal/products"
)

func TestAggregateMode(t *testing.T) {
	if products.AggregateMode([]products.Score{{Mode: "scaffold"}, {Mode: "demo"}}) != "demo" {
		t.Fatal("expected demo")
	}
	if products.AggregateMode([]products.Score{{Mode: "demo"}, {Mode: "live"}}) != "live" {
		t.Fatal("expected live")
	}
}

func TestComposeUsesLiveInputs(t *testing.T) {
	in := duo.Inputs{
		DPOOverall: 70,
		DPOLive:    true,
		ByProduct: map[string]map[string]products.Score{
			"dco": {
				"deploy_confidence":       {Name: "deploy_confidence", Value: 90, Mode: "live"},
				"operational_complexity":  {Name: "operational_complexity", Value: 20, Mode: "live"},
			},
			"daops": {
				"delivery_confidence": {Name: "delivery_confidence", Value: 80, Mode: "live"},
				"toil_ratio":          {Name: "toil_ratio", Value: 10, Mode: "live"},
			},
			"dno": {
				"service_reachability": {Name: "service_reachability", Value: 100, Mode: "live"},
				"intent_compliance":    {Name: "intent_compliance", Value: 90, Mode: "live"},
			},
			"dio": {
				"capacity_confidence": {Name: "capacity_confidence", Value: 85, Mode: "live"},
				"failover_readiness":  {Name: "failover_readiness", Value: 80, Mode: "live"},
			},
			"dao": {
				"prompt_effectiveness": {Name: "prompt_effectiveness", Value: 70, Mode: "live"},
				"ai_trust_score":       {Name: "ai_trust_score", Value: 65, Mode: "live"},
			},
			"dso": {
				"attack_surface_evolution": {Name: "attack_surface_evolution", Value: 40, Mode: "live"},
				"secrets_hygiene":          {Name: "secrets_hygiene", Value: 90, Mode: "live"},
			},
		},
	}
	chain := duo.Compose("dasmlab.org", in)
	for _, s := range chain.Sources {
		if s.Mode == "scaffold" {
			t.Fatalf("unexpected scaffold for %s.%s", s.Product, s.Name)
		}
	}
	rec := duo.Recommend("dasmlab.org", in, chain)
	if len(rec.Evidence) < 2 {
		t.Fatalf("expected ≥2 evidence, got %v", rec.Evidence)
	}
}
