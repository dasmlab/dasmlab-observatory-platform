# ADR-0008 — Plugin Marketplace

**Status:** Accepted (phased)  
**Date:** 2026-08-07  
**Depends on:** [ADR-0003](./ADR-0003-collector-sdk.md), [ADR-0007](./ADR-0007-security-multi-tenancy.md)

## Context

Community and internal teams will add collectors (Search, Cloud, AI, Network). Distribution and lifecycle must be designed early even if the marketplace UI ships late.

## Decision

### Plugin kinds

- Collector plugins ([ADR-0003](./ADR-0003-collector-sdk.md))
- Exporter plugins (webhooks, MCP, SIEM)
- Analytics packs (correlation rule bundles)
- UI widgets (Quasar modules — EE/marketplace later)

### Lifecycle

1. **Discover** — registry metadata (name, version, checksum, capabilities, config schema).
2. **Install** — pin version to tenant/org.
3. **Configure** — validate against JSON schema; attach secret refs.
4. **Enable / Disable** — without uninstall.
5. **Upgrade / Rollback** — semver; breaking changes require major bump.
6. **Health** — surface in Engineering dashboard.

### Trust tiers

| Tier | Who | Phase |
|------|-----|-------|
| First-party | DASMLAB signed | P1 (in-tree) |
| Partner | Signed + review | P3 |
| Community | Unsigned allowed only in Labs / self-host override | P3 |

### Marketplace UI

- Deferred to Phase 3 / EE.
- Until then: in-repo `plugins/` + Helm/Kustomize enables.

### SDK packages (target names)

```text
dasmlab-observatory-sdk
collector-sdk
analytics-sdk
plugin-sdk
dashboard-sdk
storage-sdk
graph-sdk
notification-sdk
workflow-sdk
ai-sdk
licensing-sdk
```

## Consequences

- P1 ships first-party plugins only; marketplace is a contract + registry stub.
- Security review required before enabling untrusted collectors against production tenants.

## Related

- [ADR-0003](./ADR-0003-collector-sdk.md)
