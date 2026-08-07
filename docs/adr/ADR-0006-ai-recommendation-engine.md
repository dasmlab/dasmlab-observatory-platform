# ADR-0006 — AI & Recommendation Engine

**Status:** Accepted  
**Date:** 2026-08-07  
**Depends on:** [ADR-0004](./ADR-0004-analytics-correlation-engine.md), [ADR-9999](./ADR-9999-innovation-principles.md)

## Context

Passive dashboards stop at L1–L2 maturity. DASMLAB targets L3–L5 (explain, predict, recommend) with a path to L6–L7. Advice must never be opaque.

## Decision

### Recommendation object (normative)

```json
{
  "id": "rec_…",
  "tenant": "dasmlab.org",
  "title": "Republish sitemap and fix two 404s",
  "summary": "Crawl frequency and impressions dropped after sitemap staleness and broken URLs.",
  "evidence": [
    {"type": "metric", "ref": "googlebot_fetches", "delta": "-22%", "window": "7d"},
    {"type": "event", "ref": "http_404", "count": 2}
  ],
  "confidence": 0.82,
  "expected_impact": {
    "scores": {"seo": "+3..+8", "overall": "+1..+3"},
    "horizon": "14d"
  },
  "estimated_effort": "S",
  "actions": [
    {"kind": "manual", "description": "Fix /docs/old-path 404"},
    {"kind": "manual", "description": "Regenerate and submit sitemap"}
  ],
  "supporting_metrics": ["impressions", "coverage_errors"],
  "created_at": "…"
}
```

Confidence ∈ [0,1]; effort ∈ {XS,S,M,L,XL}. Missing evidence → do not emit.

### Observatory AI Analyst

- Input: natural language + tenant scope.
- Tools: metrics query, insights, graph neighborhood, doc/search hits.
- Output: narrative + zero-or-more Recommendation objects.
- Models: BYO keys; no scraping of consumer chat UIs.
- All prompts/responses for GEO collectors are **data plane**; Analyst is **control plane**.

### Prediction (Phase 3)

What-if scenarios (“publish 20 articles”) produce forecast bands with methodology disclosure—not false precision.

### Automation (L6+)

Actions marked `kind: automated` require explicit tenant policy + approval workflow ([ADR-0002](./ADR-0002-common-reference-architecture.md) Workflow service). Default is recommend-only.

## Consequences

- Products may ship rule-based recommendations before LLM Analyst.
- Every AI feature must pass [ADR-9999](./ADR-9999-innovation-principles.md) explainability clauses.

## Related

- [ADR-0400](./ADR-0400-digital-presence-observatory.md) AI Observatory phase
