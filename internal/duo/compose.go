package duo

import (
	"fmt"
	"time"

	"github.com/dasmlab/dasmlab-observatory-platform/internal/products"
	observatory "github.com/dasmlab/dasmlab-observatory-platform/platform/observatory-sdk"
)

// ImpactChain is DUO F1 — Business → Engineering → Operational (not dashboard mash).
type ImpactChain struct {
	Tenant      string        `json:"tenant"`
	Generated   string        `json:"generated_at"`
	Business    ImpactLayer   `json:"business"`
	Engineering ImpactLayer   `json:"engineering"`
	Operational ImpactLayer   `json:"operational"`
	Sources     []SourceScore `json:"sources"`
}

type ImpactLayer struct {
	Label   string   `json:"label"`
	Score   float64  `json:"score"`
	Summary string   `json:"summary"`
	Signals []string `json:"signals"`
}

type SourceScore struct {
	Product string  `json:"product"`
	Name    string  `json:"name"`
	Value   float64 `json:"value"`
	Mode    string  `json:"mode"` // live | demo | scaffold
}

// Inputs carries live sibling scores for Compose (from store / product snapshots).
type Inputs struct {
	DPOOverall float64
	DPOLive    bool
	// Sibling metrics keyed product -> metric -> score
	ByProduct map[string]map[string]products.Score
}

func lookup(in Inputs, product, name string, fallback float64) (float64, string) {
	if in.ByProduct != nil {
		if m, ok := in.ByProduct[product]; ok {
			if s, ok := m[name]; ok && s.Mode != "scaffold" {
				return s.Value, s.Mode
			}
		}
	}
	return fallback, "scaffold"
}

// Compose builds an impact chain from live/demo product scores (no hardcoded scaffold table).
func Compose(tenant string, in Inputs) ImpactChain {
	sources := []SourceScore{}
	if in.DPOLive {
		sources = append(sources, SourceScore{Product: "dpo", Name: "overall", Value: in.DPOOverall, Mode: "live"})
	} else {
		sources = append(sources, SourceScore{Product: "dpo", Name: "overall", Value: in.DPOOverall, Mode: "demo"})
	}

	type need struct{ product, name string; fb float64 }
	needed := []need{
		{"dco", "deploy_confidence", 72},
		{"dco", "operational_complexity", 41},
		{"dso", "attack_surface_evolution", 38},
		{"dso", "secrets_hygiene", 55},
		{"dno", "service_reachability", 88},
		{"dno", "intent_compliance", 70},
		{"dao", "prompt_effectiveness", 60},
		{"dao", "ai_trust_score", 58},
		{"daops", "delivery_confidence", 65},
		{"daops", "toil_ratio", 35},
		{"dio", "capacity_confidence", 75},
		{"dio", "failover_readiness", 68},
	}
	vals := map[string]float64{}
	modes := map[string]string{}
	for _, n := range needed {
		v, mode := lookup(in, n.product, n.name, n.fb)
		key := n.product + "." + n.name
		vals[key] = v
		modes[key] = mode
		sources = append(sources, SourceScore{Product: n.product, Name: n.name, Value: v, Mode: mode})
	}

	dpo := in.DPOOverall
	biz := 0.4*dpo + 0.3*vals["dao.prompt_effectiveness"] + 0.3*vals["dao.ai_trust_score"]
	eng := 0.35*vals["dco.deploy_confidence"] + 0.35*vals["daops.delivery_confidence"] + 0.3*vals["dco.operational_complexity"]
	ops := 0.4*vals["dno.service_reachability"] + 0.3*vals["dno.intent_compliance"] + 0.3*vals["dio.capacity_confidence"]

	return ImpactChain{
		Tenant:    tenant,
		Generated: time.Now().UTC().Format(time.RFC3339),
		Business: ImpactLayer{
			Label:   "Business Impact",
			Score:   round1(biz),
			Summary: "Digital presence and AI-about-us signals drive brand/authority outcomes.",
			Signals: []string{"dpo.overall", "dao.ai_trust_score", "dao.prompt_effectiveness"},
		},
		Engineering: ImpactLayer{
			Label:   "Engineering Impact",
			Score:   round1(eng),
			Summary: "Deploy confidence and delivery toil explain whether the change window was healthy.",
			Signals: []string{"dco.deploy_confidence", "daops.delivery_confidence", "dco.operational_complexity"},
		},
		Operational: ImpactLayer{
			Label:   "Operational Impact",
			Score:   round1(ops),
			Summary: "Path reachability, CERT intent, and capacity underpin public truth.",
			Signals: []string{"dno.service_reachability", "dno.intent_compliance", "dio.capacity_confidence"},
		},
		Sources: sources,
	}
}

// Recommend returns one explainable action with evidence from ≥2 products (DUO F2).
func Recommend(tenant string, in Inputs, chain ImpactChain) observatory.Recommendation {
	// Prefer worst live/demo gaps (low scores).
	worstProduct, worstName, worstVal := "dpo", "overall", in.DPOOverall
	for _, s := range chain.Sources {
		if s.Mode == "scaffold" {
			continue
		}
		if s.Value < worstVal {
			worstVal = s.Value
			worstProduct = s.Product
			worstName = s.Name
		}
	}

	title := "Close the weakest live observatory signal"
	action := fmt.Sprintf(
		"Investigate %s.%s (score %.1f). Re-run collectors, confirm product API /api/v1/products/%s, then freeze a baseline.",
		worstProduct, worstName, worstVal, worstProduct,
	)
	evidence := []string{
		fmt.Sprintf("DPO overall=%.1f", in.DPOOverall),
		fmt.Sprintf("Weakest signal %s.%s=%.1f", worstProduct, worstName, worstVal),
	}
	// Add a second product evidence line.
	for _, s := range chain.Sources {
		if s.Product == worstProduct {
			continue
		}
		if s.Mode == "scaffold" {
			continue
		}
		evidence = append(evidence, fmt.Sprintf("%s.%s=%.1f (%s)", s.Product, s.Name, s.Value, s.Mode))
		break
	}
	if len(evidence) < 2 {
		evidence = append(evidence, "Run POST /api/v1/collect/run to refresh sibling product scores")
	}

	return observatory.Recommendation{
		ID:              fmt.Sprintf("duo-rec-%s", time.Now().UTC().Format("20060102")),
		Product:         observatory.ProductDUO,
		Title:           title,
		Action:          action,
		Evidence:        evidence,
		Confidence:      0.78,
		ExpectedImpact:  "DUO business/engineering/operational layers gain explainable live inputs",
		EstimatedEffort: "S–M (collector/network fix or credential; usually <1 day)",
		SupportingScores: []string{
			fmt.Sprintf("%s.%s", worstProduct, worstName),
			"dpo.overall",
		},
		FiveQuestions: []string{"Why did it change?", "What should I do?"},
	}
}

func round1(v float64) float64 {
	return float64(int(v*10+0.5)) / 10
}

// InputsFromSnapshots builds Inputs from product snapshots + DPO overall.
func InputsFromSnapshots(dpoOverall float64, dpoLive bool, snaps []products.Snapshot) Inputs {
	in := Inputs{DPOOverall: dpoOverall, DPOLive: dpoLive, ByProduct: map[string]map[string]products.Score{}}
	for _, snap := range snaps {
		if snap.Code == "dpo" || snap.Code == "duo" {
			continue
		}
		m := map[string]products.Score{}
		for _, sc := range snap.Scores {
			m[sc.Name] = sc
		}
		in.ByProduct[snap.Code] = m
	}
	return in
}
