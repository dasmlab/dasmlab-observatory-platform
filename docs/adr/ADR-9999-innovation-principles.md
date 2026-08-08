# ADR-9999 — Innovation Principles (North Star)

**Status:** Accepted — architectural North Star  
**Date:** 2026-08-07  
**Deciders:** Technologies DASMLAB Inc.  
**Depends on:** [ADR-0001](./ADR-0001-observatory-platform-vision.md)

## Context

Without an explicit gate, Observatory work drifts into commodity monitoring and marketing SEO tooling. The platform must continually produce **novel** capabilities that become research, articles, talks, and open-source knowledge—reinforcing DASMLAB authority and discoverability.

## Decision

### Feature acceptance gate

Every feature proposed for any DASMLAB Observatory **must satisfy at least one** of the following. Features that satisfy none are **rejected** or quarantined to **Labs** with an expiry date.

1. **Measure** something that existing platforms do not measure.
2. **Correlate** datasets that are normally isolated.
3. **Predict** an outcome before it becomes observable.
4. **Explain why** a change occurred, not just that it occurred.
5. **Transform** complex engineering telemetry into clear, actionable guidance.
6. **Produce** original research, benchmarks, or indices that contribute new knowledge to the industry.

### Innovation Framework (short form)

Also valid as a one-liner gate in PRs:

> Measure / Correlate / Predict / Visualize / Automate — something **no one** (credibly) does today for this domain.

### Anti-patterns (reject or demote)

- “Add another CTR chart like Semrush.”
- “Show CPU like Prometheus alone.”
- “List CVEs like a scanner UI.”
- “Aggregate three Grafana boards into DUO.”
- Opaque AI advice without evidence, confidence, impact, or effort.
- Scraping consumer AI UIs in violation of ToS (use official APIs / BYO keys only).

### Novel metric examples (illustrative, not exhaustive)

| Product | Prefer |
|---------|--------|
| DPO | AI Citation Velocity, Authority Growth, Topic Coverage, Knowledge Graph Density, Engineering Trust, Problem Ownership, Content Originality, Research Influence, Innovation Score |
| DCO | Operational Complexity, Cluster Maintainability, Technical Debt, Deployment Confidence, Recovery Confidence, Automation Maturity, Engineering Efficiency |
| DSO | Attack Surface Evolution, Risk Momentum, Patch Confidence, Exploit Probability, Blast Radius, Privilege Complexity, Secrets Hygiene |
| DAO | Model Drift, Prompt Effectiveness, Knowledge Freshness, Citation Accuracy, Reasoning Stability, Inference Cost, AI Trust Score |
| DNO | Path Stability, Routing Confidence, Failure Prediction, Policy Complexity, Service Reachability, Intent Compliance |
| DUO | Business / Engineering / Operational impact chains + recommended actions |

Commodity signals may feed these scores as inputs.

### Research obligation

Platform and product teams should backlog items that can become:

- Quarterly `dasmlab-research` publications
- Public Labs write-ups
- Conference / meetup talks
- Open benchmarks

DPO especially: research outputs are themselves digital-presence assets that DPO measures (citation flywheel).

### PR / ADR checklist (copy into contributions)

```text
[ ] Which gate criteria (1–6) does this satisfy?
[ ] What existing tool already does this? Why is ours different?
[ ] How is the output explainable (evidence / confidence)?
[ ] Does this feed research, Labs, or knowledge graph?
[ ] Which maturity level (L1–L7) does this advance?
```

## Consequences

- Product managers and coding agents treat this ADR as non-negotiable **North Star**.
- It ensures the platform never drifts into becoming just another monitoring suite and instead continually produces novel capabilities that can be turned into original articles, benchmark reports, conference talks, and open-source research—all of which reinforce the authority and discoverability of the DASMLAB ecosystem.
- Labs may violate the gate temporarily but must declare the experiment hypothesis and sunset.
- CE/EE packaging may include commodity connectors as *plumbing*, but marketing and roadmap prioritize gate-passing capabilities.

## Related

- [ADR-0001](./ADR-0001-observatory-platform-vision.md)
- [ENGINEERING-PRINCIPLES](../standards/ENGINEERING-PRINCIPLES.md)
- [ADR-0400](./ADR-0400-digital-presence-observatory.md)
