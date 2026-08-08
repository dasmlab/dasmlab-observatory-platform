# ADR-0400 — DPO — Digital Presence Observatory

**Status:** Accepted (MVP scope)  
**Date:** 2026-08-07  
**Depends on:** ADR-0001 … 0008, ADR-9999  
**Pilot tenant:** `dasmlab.org` (+ preview constellation)

## Context

SEO/GEO markets are crowded (Ahrefs, Semrush, Rankflo, OpenGSC, …). DASMLAB will not ship “another SEO tool.” DPO is **engineering observability for digital presence**—first-party crawl telemetry, GSC/Bing, Activity CDP, GitHub authority, API-only AI visibility, multi-score health, correlation WHY, and (P2) digital twin.

CDN-mgr stays separate (home ADR-001). Home/Surfing remain Activity producers.

## Decision

### Mission

Observe, measure, and improve a company’s digital footprint across search engines, AI systems, code repositories, CDN/edge, and websites—with novel scores and explainable recommendations.

### Competitive niche (lock)

1. First-party edge bot/crawl telemetry.
2. Fuse Search Console + Activity CDP + GitHub.
3. AI visibility via **official APIs / BYO keys only**.
4. Multi-score Digital Presence health + WHY correlation.
5. Self-hosted, dual-license path later.

### Non-goals (through P1)

- Ahrefs-scale backlink index  
- Scraping consumer ChatGPT/Claude UIs  
- Folding into CDN-mgr  
- Multi-tenant billing SaaS  
- 99 SEO micro-tools  

### Novel scores (primary value)

Prefer: AI Citation Velocity, Authority Growth, Topic Coverage, Knowledge Graph Density, Engineering Trust, Content Originality, Research Influence, Innovation Score.

Commodity (CTR, clicks, rank) are **inputs**, not the product story.

### Multi-score catalog (v0)

| Score | P1 | Notes |
|-------|----|-------|
| Overall Health | Yes | Weighted composite hero |
| SEO | Yes | GSC-derived |
| Engagement | Yes | Activity CDP |
| Trust (technical) | Yes | Edge errors, HTTPS, sitemap freshness |
| Authority | Partial | GitHub until graph denser |
| GEO | Placeholder → P2 | AI mention/cite/rank |
| Popularity / Freshness / Community / Developer Adoption / Commercial Interest | Progressive | Fill as collectors land |

**Overall Health weights (P1 lock, tunable later):**  
Search/SEO 25% · GEO 25% (neutral 50 until P2) · Authority 20% · Engagement 15% · Freshness 10% · Technical Trust 5%.

Document runtime weights in `products/dpo/docs/SCORES.md` when code lands.

### Collectors (DPO)

| Collector | Phase |
|-----------|-------|
| `gsc` | P1 |
| `github` | P1 |
| `edge-logs` (HAProxy / CF / optional app logs) | P1 |
| `activity` (Surfing bridge) | P1 |
| AI engines (OpenAI, Anthropic, Gemini, Perplexity, …) | P2 |
| Bing, sitemap, robots, crawler, Reddit/YouTube/LinkedIn | P2+ |

### Phases

**Phase 1 — MVP (v0.1)**  
Scheduler + plugin framework; GSC + GitHub + edge + Activity; Postgres + VictoriaMetrics; REST; Keycloak; Quasar Executive + Engineering dashboards; Overall Health history.

**Phase 2 — AI Observatory + twin**  
Prompt library × engines; citation/ranking; GEO/Content/Repo dashboards; Neo4j digital twin; OpenSearch; competitors YAML; alerts; WHY overlays.

**Phase 3 — Predictive + suite**  
What-if planning; CDN-mgr publish hooks; MCP/GraphQL/SSE; CE extract; DUO-ready exports.

### Dashboards

| Dashboard | Phase |
|-----------|-------|
| Executive | P1 |
| Engineering (bots, index, errors, load, sitemap) | P1 |
| GEO / AI Observatory | P2 |
| Content | P2 |
| Repository | P2 |

### Digital twin entities (P2)

Website → Projects → Articles → Repositories → Videos → Authors → Technologies → Search Terms → Backlinks → AI Mentions → Visitors.

### Repo layout

```text
products/dpo/          # this product (thin)
# or sibling repo digital-presence-observatory — either OK if it imports platform SDKs
```

Interim: implement under `products/dpo` in this monorepo **or** `/home/dasm/digital-presence-observatory` depending on agent capacity—**must** import `collector-sdk` / event schema from DOP.

### Success criteria

- Nightly Overall Health for dasmlab.org with history  
- Executive + Engineering live with GSC + bots + Activity + GitHub on one event schema  
- Plugin registry lists collectors with Health()  
- P2: AI mention/cite trends; ≥2 twin queries answered  

## Consequences

- Implementation follows platform ADRs; DPO does not fork analytics or identity.
- Research flywheel: DPO metrics + `dasmlab-research` reports reinforce GEO.

## Related

- [ADR-0002](./ADR-0002-common-reference-architecture.md)
- [ADR-0003](./ADR-0003-collector-sdk.md)
- Home: `docs/ADR-001-CDN-MGR-GEO-ENGAGEMENT.md`, `docs/DEMO-VISITOR-CONTRACT.md`

## Blueprint

See [products/dpo/BLUEPRINT.md](../../products/dpo/BLUEPRINT.md).

