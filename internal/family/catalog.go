package family

// Catalog is the normative Observatory Platform family view (ADR-0001 / 9999).

type Feature struct {
	ID           string   `json:"id"`
	Name         string   `json:"name"`
	NovelSignal  string   `json:"novel_signal"`
	FiveQFocus   []string `json:"five_q_focus"`
	Proof        string   `json:"proof"`
	ChainRole    string   `json:"chain_role"`
}

type Product struct {
	Code           string    `json:"code"`
	Name           string    `json:"name"`
	ADR            string    `json:"adr"`
	Status         string    `json:"status"` // live | scaffold | planned
	CommodityAvoid []string  `json:"commodity_avoid"`
	NovelScores    []string  `json:"novel_scores"`
	Features       []Feature `json:"features"`
	Blueprint      string    `json:"blueprint"`
}

type Layer struct {
	Name     string   `json:"name"`
	Elements []string `json:"elements"`
}

type Catalog struct {
	Platform         string    `json:"platform"`
	Tagline          string    `json:"tagline"`
	Differentiator   string    `json:"differentiator"`
	InnovationGate   []string  `json:"innovation_gate"`
	FiveQuestions    []string  `json:"five_questions"`
	ArchitecturePipe []string  `json:"architecture_pipeline"`
	MaturityLevels   []string  `json:"maturity_levels"`
	Products         []Product `json:"products"`
	PlatformLayers   []Layer   `json:"platform_layers"`
	SharedSDKs       []string  `json:"shared_sdks"`
	ResearchBacklog  []string  `json:"research_backlog"`
	ActiveProduct    string    `json:"active_product"`
	E2EStory         string    `json:"e2e_story"`
}

func Default() Catalog {
	return Catalog{
		Platform:       "dop",
		Tagline:        "Engineering platforms that make complex systems observable, understandable, automatable, and ultimately self-improving.",
		Differentiator: "Every product discovers things nobody measures today — not another commodity dashboard.",
		E2EStory:       "/docs/plans/E2E-PLATFORM-STORY.md",
		FiveQuestions: []string{
			"What exists?", "What changed?", "Why did it change?", "What will happen?", "What should I do?",
		},
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
				Blueprint:      "products/dpo/BLUEPRINT.md",
				CommodityAvoid: []string{"Clicks", "CTR", "Rank", "Traffic"},
				NovelScores: []string{
					"AI Citation Velocity", "Authority Growth", "Topic Coverage",
					"Knowledge Graph Density", "Engineering Trust", "Problem Ownership",
					"Content Originality", "Research Influence", "Innovation Score",
				},
				Features: []Feature{
					{
						ID: "dpo-f1", Name: "Content spine + baseline diff", NovelSignal: "Topic Coverage",
						FiveQFocus: []string{"exists", "changed", "why", "should do"},
						Proof:      "Path entities from sitemap joined to Activity/GSC/edge; labeled baselines + diff API",
						ChainRole:  "Emits presence health + baseline labels for DUO and Research",
					},
					{
						ID: "dpo-f2", Name: "Citation / crawl story", NovelSignal: "AI Citation Velocity inputs",
						FiveQFocus: []string{"exists", "will happen", "should do"},
						Proof:      "HAProxy UA bot classes (Googlebot, GPTBot, ClaudeBot, …) + GEO prompt pack scaffold",
						ChainRole:  "Feeds Research Citation Index and DUO presence line",
					},
				},
			},
			{
				Code: "dco", Name: "Cloud Observatory", ADR: "0100", Status: "scaffold",
				Blueprint:      "products/dco/BLUEPRINT.md",
				CommodityAvoid: []string{"CPU", "Memory", "Pods"},
				NovelScores: []string{
					"Operational Complexity", "Cluster Maintainability", "Technical Debt",
					"Deployment Confidence", "Recovery Confidence", "Automation Maturity",
					"Engineering Efficiency",
				},
				Features: []Feature{
					{
						ID: "dco-f1", Name: "Deploy Confidence", NovelSignal: "Deployment Confidence",
						FiveQFocus: []string{"changed", "why", "will happen", "should do"},
						Proof:      "Correlate GitOps image rollout windows ↔ pod ready/restarts for dpo-system",
						ChainRole:  "Same change window as DPO ship → cluster cost of presence change",
					},
					{
						ID: "dco-f2", Name: "Operational Complexity", NovelSignal: "Operational Complexity",
						FiveQFocus: []string{"exists", "changed", "should do"},
						Proof:      "Count controllers, PVCs, Routes, secrets, CronJobs in watched namespaces",
						ChainRole:  "Complexity index feeds DUO engineering impact",
					},
				},
			},
			{
				Code: "dso", Name: "Security Observatory", ADR: "0300", Status: "scaffold",
				Blueprint:      "products/dso/BLUEPRINT.md",
				CommodityAvoid: []string{"CVE lists alone"},
				NovelScores: []string{
					"Attack Surface Evolution", "Risk Momentum", "Patch Confidence",
					"Exploit Probability", "Blast Radius", "Privilege Complexity", "Secrets Hygiene",
				},
				Features: []Feature{
					{
						ID: "dso-f1", Name: "Attack Surface Evolution", NovelSignal: "Attack Surface Evolution",
						FiveQFocus: []string{"exists", "changed", "why", "will happen"},
						Proof:      "Diff public Routes/Services/ports over time for observatory namespaces",
						ChainRole:  "FQDN/vanity changes → surface delta → DUO risk line",
					},
					{
						ID: "dso-f2", Name: "Secrets Hygiene", NovelSignal: "Secrets Hygiene",
						FiveQFocus: []string{"exists", "why", "should do"},
						Proof:      "Optional keys in dpo-secrets (e.g. missing GSC) counted as hygiene debt",
						ChainRole:  "Blocks silent demo GSC; DUO recommends credential fill",
					},
				},
			},
			{
				Code: "dao", Name: "AI Observatory", ADR: "0500", Status: "scaffold",
				Blueprint:      "products/dao/BLUEPRINT.md",
				CommodityAvoid: []string{"Prompt logs alone"},
				NovelScores: []string{
					"Model Drift", "Prompt Effectiveness", "Knowledge Freshness",
					"Citation Accuracy", "Reasoning Stability", "Inference Cost", "AI Trust Score",
				},
				Features: []Feature{
					{
						ID: "dao-f1", Name: "Prompt Effectiveness", NovelSignal: "Prompt Effectiveness",
						FiveQFocus: []string{"exists", "changed", "why", "should do"},
						Proof:      "Versioned prompt packs (docs/PROMPTS.md) with recorded Labs outcomes",
						ChainRole:  "Internal AI quality; distinct from DPO external GEO",
					},
					{
						ID: "dao-f2", Name: "AI Trust Score", NovelSignal: "AI Trust Score",
						FiveQFocus: []string{"exists", "will happen", "should do"},
						Proof:      "Citation accuracy / knowledge freshness stubs for internal assistants",
						ChainRole:  "DUO shows DAO (our AI) beside DPO (AI about us)",
					},
				},
			},
			{
				Code: "dno", Name: "Network Observatory", ADR: "0200", Status: "scaffold",
				Blueprint:      "products/dno/BLUEPRINT.md",
				CommodityAvoid: []string{"Bandwidth alone"},
				NovelScores: []string{
					"Path Stability", "Routing Confidence", "Failure Prediction",
					"Policy Complexity", "Service Reachability", "Intent Compliance",
				},
				Features: []Feature{
					{
						ID: "dno-f1", Name: "Service Reachability", NovelSignal: "Service Reachability",
						FiveQFocus: []string{"exists", "changed", "will happen", "should do"},
						Proof:      "Probe HAProxy→Route→Service for dasmlab.org / DPO FQDN; path stability series",
						ChainRole:  "Outage narratives for DUO; complements DPO edge bots",
					},
					{
						ID: "dno-f2", Name: "Intent Compliance", NovelSignal: "Intent Compliance",
						FiveQFocus: []string{"exists", "why", "should do"},
						Proof:      "CERT*/ACL expectations vs ensure-prod-cert outcomes",
						ChainRole:  "Policy complexity feeds DUO operational impact",
					},
				},
			},
			{
				Code: "dio", Name: "Infrastructure Observatory", ADR: "0600", Status: "scaffold",
				Blueprint:      "products/dio/BLUEPRINT.md",
				CommodityAvoid: []string{"Disk / host metrics alone"},
				NovelScores: []string{
					"Capacity Confidence", "Hardware Debt", "Failover Readiness",
					"Facility Risk", "Lifecycle Coverage",
				},
				Features: []Feature{
					{
						ID: "dio-f1", Name: "Capacity Confidence", NovelSignal: "Capacity Confidence",
						FiveQFocus: []string{"exists", "will happen", "should do"},
						Proof:      "PVC/LVMS usage for observatory namespaces (not vanity disk%)",
						ChainRole:  "Underpins DCO recovery confidence",
					},
					{
						ID: "dio-f2", Name: "Failover Readiness", NovelSignal: "Failover Readiness",
						FiveQFocus: []string{"exists", "changed", "should do"},
						Proof:      "Dual-cluster GitOps drift check 2026-prod-1 vs 2026-prod-2-1",
						ChainRole:  "Recovery input to DCO/DUO",
					},
				},
			},
			{
				Code: "daops", Name: "DevOps Observatory", ADR: "0650", Status: "scaffold",
				Blueprint:      "products/daops/BLUEPRINT.md",
				CommodityAvoid: []string{"Pipeline green/red alone"},
				NovelScores: []string{
					"Delivery Confidence", "Change Failure Momentum", "Toil Ratio",
					"Platform Friction", "Developer Experience Index",
				},
				Features: []Feature{
					{
						ID: "daops-f1", Name: "Delivery Confidence", NovelSignal: "Delivery Confidence",
						FiveQFocus: []string{"changed", "why", "will happen", "should do"},
						Proof:      "CX pipeline success / queue time on org runner for observatory + home",
						ChainRole:  "Explains risky DPO/DCO windows (queue, cancels)",
					},
					{
						ID: "daops-f2", Name: "Toil Ratio", NovelSignal: "Toil Ratio",
						FiveQFocus: []string{"exists", "why", "should do"},
						Proof:      "Manual workflow_dispatch / cert ensures vs auto GitOps publishes",
						ChainRole:  "Platform friction → DUO engineering impact",
					},
				},
			},
			{
				Code: "duo", Name: "Unified Observatory", ADR: "0700", Status: "scaffold",
				Blueprint:      "products/duo/BLUEPRINT.md",
				CommodityAvoid: []string{"Aggregated dashboards"},
				NovelScores: []string{
					"Business Impact", "Engineering Impact", "Operational Impact", "Recommended Actions",
				},
				Features: []Feature{
					{
						ID: "duo-f1", Name: "Impact chain view", NovelSignal: "Business→Engineering→Operational",
						FiveQFocus: []string{"exists", "changed", "why", "will happen", "should do"},
						Proof:      "Compose scores from DPO + sibling stubs into impact chain API",
						ChainRole:  "Executive ending of E2E demo",
					},
					{
						ID: "duo-f2", Name: "Recommended action card", NovelSignal: "Recommended Actions",
						FiveQFocus: []string{"why", "should do"},
						Proof:      "One explainable action with evidence from ≥2 products (ADR-0006 shape)",
						ChainRole:  "Closes the five-question loop for leadership",
					},
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
