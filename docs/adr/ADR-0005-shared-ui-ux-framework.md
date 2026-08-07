# ADR-0005 — Shared UI/UX Framework

**Status:** Accepted  
**Date:** 2026-08-07  
**Depends on:** [ADR-0001](./ADR-0001-observatory-platform-vision.md), [ADR-0002](./ADR-0002-common-reference-architecture.md)

## Context

Observatories must feel like one product family. DASMLAB already uses Vue/Quasar on home and sibling apps.

## Decision

### Stack

- **Vue 3 + Quasar** for all Observatory front ends.
- Shared package: `platform/ui-components` (design tokens, shell, score widgets, insight cards, source health).
- Visual language: align with dasmlab.org family (expressive typography, atmospheric backgrounds)—avoid generic “AI purple SaaS” clichés; prefer DASMLAB lab aesthetic.

### Shell contract

Every Observatory app includes:

1. **Product switcher** (links to other observatories / DUO when available).
2. **Tenant context** (pilot single-tenant OK).
3. **Five-question navigation affordance** (Exists / Changed / Why / Forecast / Act).
4. **Score strip** (Overall + primary sub-scores).
5. **Source health** (collector plugins Health()).
6. **Insight / recommendation panel** (explainable cards).

### Dashboard types (shared patterns)

| Pattern | Purpose |
|---------|---------|
| Executive | Multi-score health, sparklines, top insights |
| Engineering | Technical signals, errors, bots, SLOs |
| Domain | Product-specific (GEO, Content, Cluster, …) |
| Entity detail | One repo / article / cluster deep dive |

### Auth UX

- Keycloak OIDC (realm `dasmlab`) per [ADR-0007](./ADR-0007-security-multi-tenancy.md).
- Anonymous/public executive views are product decisions (DPO pilot: authenticated first).

### Accessibility & motion

- Prefer intentional motion (2–3) for hierarchy; avoid noise.
- Score changes announce direction (↑/↓) with text, not color alone.

## Consequences

- Product UIs import shared shell; they do not fork Quasar bootstraps long-term.
- Mobile / thin executive layouts are Phase 3.

## Related

- [ADR-0400](./ADR-0400-digital-presence-observatory.md) dashboard list
