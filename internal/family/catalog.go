package family

// Catalog is the normative Observatory Platform family view (ADR-0001 / 9999).
// DPO is the first live specialization; siblings are scaffolded until their ADR expands.

type Product struct {
	Code           string   `json:"code"`
	Name           string   `json:"name"`
	ADR            string   `json:"adr"`
	Status         string   `json:"status"` // live | scaffold | planned
	CommodityAvoid []string `json:"commodity_avoid"`
	NovelScores    []string `json:"novel_scores"`
}

type Layer struct {
	Name     string   `json:"name"`
	Elements []string `json:"elements"`
}

type Catalog struct {
	Platform           string    `json:"platform"`
	Tagline            string    `json:"tagline"`
	Differentiator     string    `json:"differentiator"`
	InnovationGate     []string  `json:"innovation_gate"`
	ArchitecturePipe   []string  `json:"architecture_pipeline"`
	MaturityLevels     []string  `json:"maturity_levels"`
	Products           []Product `json:"products"`
	PlatformLayers     []Layer   `json:"platform_layers"`
	SharedSDKs         []string  `json:"shared_sdks"`
	ResearchBacklog    []string  `json:"research_backlog"`
	ActiveProduct      string    `json:"active_product"`
}

func Default() Catalog {
	return Catalog{
		Platform:       "dop",
		Tagline:        "Engineering platforms that make complex systems observable, understandable, automatable, and ultimately self-improving.",
		Differentiator: "Every product discovers things nobody measures today — not another commodity dashboard.",
		InnovationGate: []string{
			"Measure something no one measures",
			"Correlate things no one correlates",
			"Predict something no one predicts",
			"Visualize something no one visualizes",
			"Automate something no one automates",
		},
		ArchitecturePipe: []string{
			"Collectors", "Normalization", "Correlation", "Analytics", "AI",
			"Storage", "API", "Dashboard", "Recommendations",
		},
		MaturityLevels: []string{
			"L1 Observe", "L2 Understand", "L3 Explain", "L4 Predict",
			"L5 Recommend", "L6 Automate", "L7 Learn",
		},
		ActiveProduct: "dpo",
		Products: []Product{
			{
				Code: "dpo", Name: "Digital Presence Observatory", ADR: "0400", Status: "live",
				CommodityAvoid: []string{"Clicks", "CTR", "Rank", "Traffic"},
				NovelScores: []string{
					"AI Citation Velocity", "Authority Growth", "Topic Coverage",
					"Knowledge Graph Density", "Engineering Trust", "Problem Ownership",
					"Content Originality", "Research Influence", "Innovation Score",
				},
			},
			{
				Code: "dco", Name: "Cloud Observatory", ADR: "0100", Status: "scaffold",
				CommodityAvoid: []string{"CPU", "Memory", "Pods"},
				NovelScores: []string{
					"Operational Complexity", "Cluster Maintainability", "Technical Debt",
					"Deployment Confidence", "Recovery Confidence", "Automation Maturity",
					"Engineering Efficiency",
				},
			},
			{
				Code: "dso", Name: "Security Observatory", ADR: "0300", Status: "scaffold",
				CommodityAvoid: []string{"CVE lists alone"},
				NovelScores: []string{
					"Attack Surface Evolution", "Risk Momentum", "Patch Confidence",
					"Exploit Probability", "Blast Radius", "Privilege Complexity", "Secrets Hygiene",
				},
			},
			{
				Code: "dao", Name: "AI Observatory", ADR: "0500", Status: "scaffold",
				CommodityAvoid: []string{"Prompt logs alone"},
				NovelScores: []string{
					"Model Drift", "Prompt Effectiveness", "Knowledge Freshness",
					"Citation Accuracy", "Reasoning Stability", "Inference Cost", "AI Trust Score",
				},
			},
			{
				Code: "dno", Name: "Network Observatory", ADR: "0200", Status: "scaffold",
				CommodityAvoid: []string{"Bandwidth alone"},
				NovelScores: []string{
					"Path Stability", "Routing Confidence", "Failure Prediction",
					"Policy Complexity", "Service Reachability", "Intent Compliance",
				},
			},
			{
				Code: "dio", Name: "Infrastructure Observatory", ADR: "0600", Status: "scaffold",
				CommodityAvoid: []string{"Disk / host metrics alone"},
				NovelScores: []string{
					"Capacity Confidence", "Hardware Debt", "Failover Readiness",
					"Facility Risk", "Lifecycle Coverage",
				},
			},
			{
				Code: "daops", Name: "DevOps Observatory", ADR: "0650", Status: "scaffold",
				CommodityAvoid: []string{"Pipeline green/red alone"},
				NovelScores: []string{
					"Delivery Confidence", "Change Failure Momentum", "Toil Ratio",
					"Platform Friction", "Developer Experience Index",
				},
			},
			{
				Code: "duo", Name: "Unified Observatory", ADR: "0700", Status: "scaffold",
				CommodityAvoid: []string{"Aggregated dashboards"},
				NovelScores: []string{
					"Business Impact", "Engineering Impact", "Operational Impact", "Recommended Actions",
				},
			},
		},
		PlatformLayers: []Layer{
			{Name: "Applications", Elements: []string{"DUO", "DPO", "DCO", "DAO", "DSO", "DNO", "DIO", "DAOps"}},
			{Name: "Platform services", Elements: []string{"AI", "Analytics", "Graph", "Workflow", "Search", "Collectors", "Storage", "Notifications"}},
			{Name: "Shared SDK", Elements: []string{"UI", "REST", "Plugins", "Security", "Licensing", "CLI", "Configuration"}},
			{Name: "Infrastructure", Elements: []string{"PostgreSQL", "VictoriaMetrics", "OpenSearch", "Neo4j", "Redis", "NATS", "S3", "Keycloak", "Kubernetes"}},
		},
		SharedSDKs: []string{
			"observatory-sdk", "collector-sdk", "analytics-sdk", "plugin-sdk", "dashboard-sdk",
			"storage-sdk", "graph-sdk", "search-sdk", "notification-sdk", "workflow-sdk", "ai-sdk", "licensing-sdk",
		},
		ResearchBacklog: []string{
			"AI Citation Index 2026",
			"Open Source Visibility Report",
			"Kubernetes Adoption Report",
			"Cloud Complexity Report",
			"Engineering Observability Whitepaper",
			"AI Trust Benchmark",
			"Developer Experience Index",
		},
	}
}
