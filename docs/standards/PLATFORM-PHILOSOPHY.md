# Platform philosophy

Engineering platforms that make complex systems **observable, understandable, automatable, and ultimately self-improving**.

## One platform, many products

We do not build Observatory products independently. We build **one platform** that manifests as multiple products — the same pattern as vSphere, HashiCorp, Grafana Labs, OpenShift, Dynatrace, Splunk.

Every DXX (DCO, DNO, DSO, DPO, DAO, DAOps, DIO, DUO) is a **specialization** of DOP: shared collectors, events, analytics, recommendations, identity, and UI shell. Thin products. Heavy platform.

## Mission

Build engineering observability systems that measure complex domains in ways that are currently **unavailable, fragmented, or prohibitively expensive**.

| Not | Instead |
|-----|---------|
| Build dashboards | Invent **new forms of observability** |
| Build monitoring | Answer **what should happen next** |
| Another SEO / APM tool | Engineering observability for that domain |

## Vs Grafana-class tools

Grafana (and peers) answer: **“What is happening.”**

We answer: **“What should happen next.”**

That is why every Observatory is bound to the [five questions](./ENGINEERING-PRINCIPLES.md) and [ADR-9999](../adr/ADR-9999-innovation-principles.md). Commodity signals (CPU, CTR, CVE lists, bandwidth) may be **inputs**. They must not be the product story.

## Digital presence example

Do not ask only “How do I see my SEO?”

Ask: **“How do I build an observability platform for my digital presence?”** — treat presence like you would observe a Kubernetes cluster (collectors, correlation, explain, recommend). GEO/SEO are telemetry sources, not the product.

## Normative refs

- [ADR-0001 Vision](../adr/ADR-0001-observatory-platform-vision.md)
- [ADR-0002 Architecture](../adr/ADR-0002-common-reference-architecture.md)
- [ADR-9999 Innovation Principles](../adr/ADR-9999-innovation-principles.md)
- [Observatory family](../OBSERVATORY-FAMILY.md)
- [Documentation as product](./DOCUMENTATION-AS-PRODUCT.md)
- [Knowledge model](./KNOWLEDGE-MODEL.md)
- [Flywheel](./FLYWHEEL.md)
