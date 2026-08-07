# ADR-0007 — Security & Multi-Tenancy

**Status:** Accepted  
**Date:** 2026-08-07  
**Depends on:** [ADR-0002](./ADR-0002-common-reference-architecture.md)

## Context

Observatories handle credentials (GSC, GitHub, AI APIs), first-party logs, and eventually multi-org SaaS. Identity must be shared, not reinvented per product.

## Decision

### Identity

- **Keycloak** realm `dasmlab` (shared with home/mock-me pattern).
- OIDC for UIs; API keys / service accounts for collectors and agents.
- Concepts: Organization → Teams → Projects → Tenants (digital properties / clusters).

### RBAC (v0 roles)

| Role | Capabilities |
|------|----------------|
| `viewer` | Read scores, insights, dashboards |
| `analyst` | Run scans, manage prompt packs |
| `operator` | Configure collectors, credentials refs |
| `admin` | Org/tenant, users, licensing |

### Tenancy

- Phase 1: single pilot tenant `dasmlab.org` (hard-coded allowed).
- Phase 2+: org-scoped tenants; row-level `tenant` on all events/scores.
- No cross-tenant queries without explicit federation (DUO later).

### Secrets

- Collector credentials in K8s secrets / sealed secrets / vault; referenced by name in config.
- Never commit tokens; never log secret values.
- AI API keys are tenant-owned (BYO).

### Licensing (CE path)

- Dual license: Apache for pure libs where possible; MPL + commercial for CE apps (align with dasmlab_home open-core policy).
- License SDK checks feature flags for EE (marketplace, SSO extras, multi-tenant).

### Compliance posture

- Prefer official APIs; document ToS constraints in product ADRs.
- Activity/PII: follow Surfing Activity Phase A contract; minimize retention; document TTL.

## Consequences

- All Observatory APIs validate Bearer/OIDC or API key.
- CDN-mgr and home keep their own deploy auth where already established; DOP does not replace cluster IdP.

## Related

- [ADR-0008](./ADR-0008-plugin-marketplace.md)
- Home: `docs/KEYCLOAK_SETUP.md`, `docs/OPEN-CORE.md`
