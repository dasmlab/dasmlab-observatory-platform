# ADR-0004 — Analytics & Correlation Engine

**Status:** Accepted  
**Date:** 2026-08-07  
**Depends on:** [ADR-0002](./ADR-0002-common-reference-architecture.md), [ADR-0003](./ADR-0003-collector-sdk.md)

## Context

Plots alone are not differentiation. DASMLAB must correlate isolated datasets and answer **why** metrics moved.

## Decision

### Pipeline stages

```text
Events → Enrichment → Entity Resolution → Correlation → Aggregation → Scores → Insights
```

| Stage | Responsibility |
|-------|----------------|
| Enrichment | Geo, bot class, content type, tech tags |
| Entity Resolution | Map strings → graph/Postgres entity IDs |
| Correlation | Time-aligned rules + optional ML later |
| Aggregation | Rollups (hour/day/week) in PG + VM |
| Scores | Multi-score engine (product-defined weights) |
| Insights | Human-readable WHY chains |

### Correlation (P1: rules; P2+: ML assist)

Rule example (DPO):

```yaml
id: traffic-up-crawl-up
when:
  - metric: page_views
    change: ">= +15%"
    window: 7d
then_look_for:
  - metric: googlebot_fetches
    change: ">= +10%"
  - metric: github_stars_delta
    change: "> 0"
  - metric: ai_citations_delta
    change: "> 0"
  - event: content_published
emit:
  insight_template: "Traffic rose with crawl, authority, and/or AI citation signals"
```

Insights always cite contributing events (explainability).

### Score engine

- Products register named scores (Overall Health, SEO, GEO, …).
- Each score = weighted blend of normalized metric components (0–100).
- Weights live in versioned config; changes emit audit events.
- Nightly (or on-demand) snapshot to Postgres; series to VictoriaMetrics.

### Storage mapping

- Raw/normalized events: Postgres (hot) + optional S3 archive.
- Metrics: VictoriaMetrics.
- Textable text (AI responses, articles): OpenSearch P2.
- Relationships: Neo4j P2.

### API

- `GET /api/v1/scores/{name}`
- `GET /api/v1/insights?from=&to=`
- `GET /api/v1/metrics/query` (PromQL-like or simplified)
- `GET /api/v1/entities/{id}/timeline`

## Consequences

- Product ADRs define score catalogs and correlation packs.
- AI narrative layer is [ADR-0006](./ADR-0006-ai-recommendation-engine.md); this ADR owns deterministic analytics.

## Related

- [ADR-0400](./ADR-0400-digital-presence-observatory.md) score table
