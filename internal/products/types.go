package products

import "time"

// Score is one novel metric for a product (DUO source name).
type Score struct {
	Name  string  `json:"name"`
	Value float64 `json:"value"`
	Mode  string  `json:"mode"` // live | demo | scaffold
}

// Snapshot is the uniform product API payload.
type Snapshot struct {
	Code     string    `json:"code"`
	Status   string    `json:"status"` // live | demo | scaffold
	Mode     string    `json:"mode"`   // aggregate: live if any live, else demo, else scaffold
	Scores   []Score   `json:"scores"`
	Features []string  `json:"features,omitempty"`
	LastRun  string    `json:"last_run,omitempty"`
	Proof    []string  `json:"proof,omitempty"`
}

// MetricSpec binds a product to its F1/F2 metric names.
type MetricSpec struct {
	Code     string
	Metrics  []string
	Features []string
	Proof    []string
}

// Specs is the normative product → metric map (DUO contract).
func Specs() []MetricSpec {
	return []MetricSpec{
		{
			Code: "dco", Metrics: []string{"deploy_confidence", "operational_complexity"},
			Features: []string{"Deploy Confidence", "Operational Complexity"},
			Proof:    []string{"GitOps image ↔ pod ready/restarts", "Count Deploy/PVC/Route/Secret/CronJob"},
		},
		{
			Code: "dso", Metrics: []string{"attack_surface_evolution", "secrets_hygiene"},
			Features: []string{"Attack Surface Evolution", "Secrets Hygiene"},
			Proof:    []string{"Routes/Services/ports inventory", "dpo-secrets key presence"},
		},
		{
			Code: "dno", Metrics: []string{"service_reachability", "intent_compliance"},
			Features: []string{"Service Reachability", "Intent Compliance"},
			Proof:    []string{"HTTP probe public FQDNs + Route", "CERT/vanity host intent vs Route TLS"},
		},
		{
			Code: "dao", Metrics: []string{"prompt_effectiveness", "ai_trust_score"},
			Features: []string{"Prompt Effectiveness", "AI Trust Score"},
			Proof:    []string{"PROMPTS.md pack coverage", "Bot signals + prompt freshness blend"},
		},
		{
			Code: "daops", Metrics: []string{"delivery_confidence", "toil_ratio"},
			Features: []string{"Delivery Confidence", "Toil Ratio"},
			Proof:    []string{"GitHub Actions success rate", "workflow_dispatch share vs push"},
		},
		{
			Code: "dio", Metrics: []string{"capacity_confidence", "failover_readiness"},
			Features: []string{"Capacity Confidence", "Failover Readiness"},
			Proof:    []string{"PVC dpo-data capacity", "Dual-cluster health probe"},
		},
		{
			Code: "duo", Metrics: []string{"business_impact", "engineering_impact", "operational_impact"},
			Features: []string{"Impact chain", "Recommended action"},
			Proof:    []string{"Compose sibling scores", "ADR-0006 recommend with ≥2 evidence"},
		},
		{
			Code: "dpo", Metrics: []string{"overall"},
			Features: []string{"Content spine + baseline", "Citation / crawl"},
			Proof:    []string{"Score engine overall", "GSC/edge/activity collectors"},
		},
	}
}

func Spec(code string) (MetricSpec, bool) {
	for _, s := range Specs() {
		if s.Code == code {
			return s, true
		}
	}
	return MetricSpec{}, false
}

// AggregateMode picks live > demo > scaffold from score modes.
func AggregateMode(scores []Score) string {
	hasDemo := false
	for _, s := range scores {
		switch s.Mode {
		case "live":
			return "live"
		case "demo":
			hasDemo = true
		}
	}
	if hasDemo {
		return "demo"
	}
	return "scaffold"
}

// StatusFromMode maps mode to Family-facing status.
func StatusFromMode(mode string) string {
	switch mode {
	case "live", "demo":
		return "live"
	default:
		return "scaffold"
	}
}

// FreshWindow is how far back product events count as current.
func FreshWindow() time.Duration { return 48 * time.Hour }
