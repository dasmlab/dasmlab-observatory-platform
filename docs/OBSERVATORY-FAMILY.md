# Observatory family (DOP)

Engineering platforms that make complex systems **observable, understandable, automatable, and ultimately self-improving**.

One Digital Observatory / observability platform. Same architecture. Different plugins. Thin products.

## Differentiator

Every product discovers things **nobody measures today**. Commodity clicks/CPU/CVE boards fail ADR-9999.

Live machine-readable catalog: `GET /api/v1/family` (served from the DPO runtime until DUO ships).

## Products

| Code | Product | ADR | Status | Instead of commodity… |
|------|---------|-----|--------|------------------------|
| DPO | Digital Presence | 0400 | **Live** + blueprint | Citation Velocity, Topic Coverage, Engineering Trust, … |
| DCO | Cloud | 0100 | Blueprint + scaffold scores | Operational Complexity, Deployment Confidence, … |
| DSO | Security | 0300 | Blueprint + scaffold scores | Attack Surface Evolution, Risk Momentum, … |
| DAO | AI | 0500 | Blueprint + scaffold scores | Model Drift, AI Trust Score, … |
| DNO | Network | 0200 | Blueprint + scaffold scores | Path Stability, Intent Compliance, … |
| DIO | Infrastructure | 0600 | Blueprint + scaffold scores | Capacity Confidence, Failover Readiness, … |
| DAOps | DevOps | 0650 | Blueprint + scaffold scores | Delivery Confidence, Developer Experience Index, … |
| DUO | Unified | 0700 | Blueprint + compose API | Executive understanding + recommended actions |

## Architecture (identical per observatory)

Collectors → Normalization → Correlation → Analytics → AI → Storage → API → Dashboard → Recommendations

## Innovation gate (ADR-9999)

Measure / Correlate / Predict / Visualize / Automate — something no one does today for this domain.

## Maturity

L1 Observe → L2 Understand → L3 Explain → L4 Predict → L5 Recommend → L6 Automate → L7 Learn

## Repo layout (this monorepo)

```
platform/     shared SDKs (collector, graph, search, …)
products/     thin specializations (dpo live; siblings scaffold)
plugins/      collector plugins (google, github, k8s, …)
labs/         experiments
research/     quarterly indices / whitepapers backlog
docs/adr/     platform + product ADRs
```

Future: split SDKs / `dasmlab-research` into sibling repos under `github.com/dasmlab` when weight demands it.

## Research flywheel

Research → Experiment → Prototype → Product → Usage → Metrics → Insights → Research  
Reports become backlinks, SEO, and AI citations — DPO helps produce the research that improves DPO.
