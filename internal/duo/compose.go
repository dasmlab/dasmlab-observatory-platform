package duo

import (
	"fmt"
	"time"

	"github.com/dasmlab/dasmlab-observatory-platform/internal/family"
	observatory "github.com/dasmlab/dasmlab-observatory-platform/platform/observatory-sdk"
)

// ImpactChain is DUO F1 — Business → Engineering → Operational (not dashboard mash).
type ImpactChain struct {
	Tenant     string         `json:"tenant"`
	Generated  string         `json:"generated_at"`
	Business   ImpactLayer    `json:"business"`
	Engineering ImpactLayer   `json:"engineering"`
	Operational ImpactLayer   `json:"operational"`
	Sources    []SourceScore  `json:"sources"`
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
	Mode    string  `json:"mode"` // live | scaffold
}

// Compose builds an impact chain from Family features + optional live DPO overall.
func Compose(tenant string, dpoOverall float64, dpoLive bool) ImpactChain {
	cat := family.Default()
	sources := []SourceScore{}
	if dpoLive {
		sources = append(sources, SourceScore{Product: "dpo", Name: "overall", Value: dpoOverall, Mode: "live"})
	} else {
		sources = append(sources, SourceScore{Product: "dpo", Name: "overall", Value: 55, Mode: "scaffold"})
	}
	// Scaffold novel scores for sibling existence proofs (never CPU/CTR heroes).
	scaffold := []SourceScore{
		{Product: "dco", Name: "deploy_confidence", Value: 72, Mode: "scaffold"},
		{Product: "dco", Name: "operational_complexity", Value: 41, Mode: "scaffold"},
		{Product: "dso", Name: "attack_surface_evolution", Value: 38, Mode: "scaffold"},
		{Product: "dso", Name: "secrets_hygiene", Value: 55, Mode: "scaffold"},
		{Product: "dno", Name: "service_reachability", Value: 88, Mode: "scaffold"},
		{Product: "dno", Name: "intent_compliance", Value: 70, Mode: "scaffold"},
		{Product: "dao", Name: "prompt_effectiveness", Value: 60, Mode: "scaffold"},
		{Product: "dao", Name: "ai_trust_score", Value: 58, Mode: "scaffold"},
		{Product: "daops", Name: "delivery_confidence", Value: 65, Mode: "scaffold"},
		{Product: "daops", Name: "toil_ratio", Value: 35, Mode: "scaffold"},
		{Product: "dio", Name: "capacity_confidence", Value: 75, Mode: "scaffold"},
		{Product: "dio", Name: "failover_readiness", Value: 68, Mode: "scaffold"},
	}
	sources = append(sources, scaffold...)

	biz := 0.4*dpoOverall + 0.3*60 + 0.3*58 // presence + prompt + trust blend
	if !dpoLive {
		biz = 55
	}
	eng := 0.35*72 + 0.35*65 + 0.3*41 // deploy + delivery − complexity pressure
	ops := 0.4*88 + 0.3*70 + 0.3*75   // reachability + intent + capacity

	_ = cat
	return ImpactChain{
		Tenant:    tenant,
		Generated: time.Now().UTC().Format(time.RFC3339),
		Business: ImpactLayer{
			Label:   "Business Impact",
			Score:   round1(biz),
			Summary: "Digital presence and AI-about-us signals drive brand/authority outcomes.",
			Signals: []string{"dpo.overall", "dao.ai_trust_score", "research.citation_index"},
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
func Recommend(tenant string, dpoOverall float64) observatory.Recommendation {
	return observatory.Recommendation{
		ID:      fmt.Sprintf("duo-rec-%s", time.Now().UTC().Format("20060102")),
		Product: observatory.ProductDUO,
		Title:   "Restore live presence truth and close GSC hygiene gap",
		Action:  "Ensure in-cluster fetch of https://dasmlab.org/sitemap.xml succeeds (or publish edge sample), add GSC_CREDENTIALS_JSON to dpo-secrets, freeze baseline post-collect, then re-check DUO impact.",
		Evidence: []string{
			fmt.Sprintf("DPO overall=%.1f (live or last known)", dpoOverall),
			"DSO secrets_hygiene scaffold flags optional GSC key debt",
			"DNO service_reachability depends on HAProxy→Route path for public FQDNs",
			"DPO content spine empty when cluster cannot hairpin to edge",
		},
		Confidence:      0.72,
		ExpectedImpact:  "Topic Coverage and SEO components leave demo/partial state; DUO business score gains explainable inputs",
		EstimatedEffort: "M (credential + network hairpin or edge-side collect; <1 day)",
		SupportingScores: []string{
			"dpo.topic_coverage", "dso.secrets_hygiene", "dno.service_reachability",
		},
		FiveQuestions: []string{
			"Why did it change?", "What should I do?",
		},
	}
}

func round1(v float64) float64 {
	return float64(int(v*10+0.5)) / 10
}
