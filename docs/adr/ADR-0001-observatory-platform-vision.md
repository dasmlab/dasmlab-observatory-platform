# ADR-0001 — DASMLAB Observatory Platform Vision

**Status:** Accepted  
**Date:** 2026-08-07  
**Deciders:** Technologies DASMLAB Inc.  
**Branch:** `2026-dop-v0`

## Context

Over the past year DASMLAB has shipped many strong *projects* (home/Surfing, CDN-mgr, cheapcloud, mock-me, interview-me, etcd-synth, camera-scrape, observability stacks). They share patterns but not a single **platform philosophy**. Commodity monitoring and SEO vendors answer “what is happening.” DASMLAB needs a coherent ecosystem that invents **new forms of observability** and products that are specializations of one platform—not independent apps.

Comparable patterns: VMware vSphere, HashiCorp, Grafana Labs, Red Hat OpenShift, Dynatrace, Splunk.

## Decision

### 1. Root project is DOP

**DASMLAB Observatory Platform (DOP)** is the root. All observatories inherit shared architecture, SDKs, event model, identity, analytics, and UI shell.

### 2. Company vision

Technologies DASMLAB Inc. develops engineering platforms that make complex systems **observable, understandable, automatable, and ultimately self-improving**.

### 3. Platform mission

Build engineering observability systems that measure complex domains in ways that are currently **unavailable, fragmented, or prohibitively expensive**.

We do **not** optimize for “another dashboard” or “another monitoring suite.”

### 4. Five questions (normative UX/API contract)

Every Observatory must help the user answer:

1. **What exists?**
2. **What changed?**
3. **Why did it change?**
4. **What will happen?**
5. **What should I do?**

Differentiator vs Grafana-class tools: we prioritize **guidance and next action**, not only live state.

### 5. Product family

| Code | Name | Domain |
|------|------|--------|
| DCO | Cloud Observatory | Cloud / Kubernetes operations |
| DNO | Network Observatory | Networks and connectivity |
| DSO | Security Observatory | Security posture and risk |
| DPO | Digital Presence Observatory | SEO / GEO / authority / brand / AI visibility |
| DAO | AI Observatory | LLM usage, prompt quality, model telemetry |
| DAOps | DevOps Observatory | Delivery pipelines, CI/CD, release confidence |
| DIO | Infrastructure Observatory | Physical / virt / storage |
| DUO | Unified Observatory | Executive single pane across observatories |

**First specialization to implement:** DPO ([ADR-0400](./ADR-0400-digital-presence-observatory.md)).

### 6. Research → Platform → Product lifecycle

```text
Research area → DOP platform → Product → Community Edition → Enterprise → Managed SaaS
```

Research areas include: Cloud Native, AI, Digital Presence, Networking, Infrastructure, Cyber Security, Automation, Developer Experience, Observability, Knowledge Engineering, Distributed Systems.

### 7. Research & Labs flywheel

- **`dasmlab-research`**: quarterly original reports/indices (AI Citation Index, OSS Visibility, Cloud Complexity, AI Trust Benchmark, DX Index, …).
- **Labs**: explicit non-production experiments that become content, authority, and eventually product features.
- Product usage metrics feed research; research improves SEO/GEO; visibility brings users; users generate more data.

### 8. Novel measurement mandate

Every product strives to discover signals **nobody measures today**. Commodity metrics (CPU, CTR, CVE lists, bandwidth) may be collected as *inputs* but must not be the product’s primary value proposition. See tables in [ADR-9999](./ADR-9999-innovation-principles.md) and product ADRs.

### 9. Docs as product

Documentation is a primary engineering output: README, Getting Started, Architecture, ADR, API, SDK, Examples, Tutorials, Research, Benchmarks, Roadmap, FAQ, Contributing, Glossary, Reference, Release Notes. These feed developer docs, SEO/GEO pages, Behind the Design stories, and the knowledge graph.

## Consequences

- New work is gated by [ADR-9999](./ADR-9999-innovation-principles.md).
- Coding agents consume platform ADRs before product ADRs; products stay thin.
- `dasmlab_home` remains the public Engineering Knowledge Network and Activity producer; it is **not** the Observatory Platform repo.
- `dasmlab-cdn-mgr` remains object/realm/CDN platform; engagement may dual-write later (home ADR-001) but is not DOP.
- GitHub org layout targets `github.com/dasmlab/{platform,plugins,products,labs,research}`; interim hosting under `lmcdasm/` / local `/home/dasm/dasmlab-observatory-platform` is acceptable.

## Related

- [ADR-0002](./ADR-0002-common-reference-architecture.md)
- [ADR-9999](./ADR-9999-innovation-principles.md)
- [ENGINEERING-PRINCIPLES](../standards/ENGINEERING-PRINCIPLES.md)
