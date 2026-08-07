# ADR-0003 — Collector SDK

**Status:** Accepted  
**Date:** 2026-08-07  
**Depends on:** [ADR-0002](./ADR-0002-common-reference-architecture.md)

## Context

Integrations must be Prometheus-exporter-like plugins, not hard-coded product modules. Coding agents and community contributors need a stable Go interface and lifecycle.

## Decision

### Go interface (normative)

```go
package collector

import "context"

type Event struct {
    SchemaVersion string            `json:"schema_version"`
    Collector     string            `json:"collector"`
    Type          string            `json:"type"`
    Timestamp     time.Time         `json:"timestamp"`
    Tenant        string            `json:"tenant"`
    Entity        string            `json:"entity"`
    Metric        string            `json:"metric"`
    Value         float64           `json:"value"`
    Unit          string            `json:"unit,omitempty"`
    Dims          map[string]string `json:"dims,omitempty"`
    TraceID       string            `json:"trace_id,omitempty"`
}

type Collector interface {
    Name() string
    Discover(ctx context.Context) error
    Collect(ctx context.Context) error
    Normalize(ctx context.Context) ([]Event, error)
    Health(ctx context.Context) error
}
```

### Lifecycle

1. **Register** with Plugin Manager (name, version, config schema, credential refs).
2. **Discover** — enumerate entities (sites, repos, prompts) for the tenant.
3. **Collect** — pull raw data into an internal buffer / object store.
4. **Normalize** — emit `[]Event` to the bus.
5. **Health** — return error if credentials, quotas, or upstream APIs fail.

Scheduler invokes Discover (as needed) → Collect → Normalize on cron or manual trigger. Failures are recorded by Collector Manager with backoff.

### Config & credentials

- Config: versioned YAML/JSON per tenant + collector (`config.schema.json` published by plugin).
- Credentials: never in git; Keycloak-linked secrets or K8s secrets referenced by name.
- Rate limits: plugins declare `RateLimit` hints; Scheduler enforces.

### Packaging

- In-tree plugins under `plugins/<name>` for first-party.
- External plugins: Go module implementing `Collector` + optional gRPC sidecar later (P3 marketplace).
- Language: **Go 1.22+** for v0 SDK. Other languages via gRPC adapter in marketplace phase.

### Minimum first-party collectors (platform + DPO)

| Plugin | Domain | Phase |
|--------|--------|-------|
| `gsc` | Google Search Console | DPO P1 |
| `github` | GitHub API | DPO P1 |
| `edge-logs` | nginx / Cloudflare bot & traffic logs | DPO P1 |
| `activity` | Surfing Activity CDP bridge | DPO P1 |
| `ai-openai` / `ai-anthropic` / … | Official AI APIs | DPO P2 |
| `bing` | Bing Webmaster | DPO P2 |
| `sitemap` / `robots` / `crawler` | Site hygiene | DPO P2 |

### Testing contract

- Golden fixtures: raw → events.
- `Health` unit tests with fake upstreams.
- No live credentials in CI; use recorded HTTP.

## Consequences

- Products depend on `platform/collector-sdk`; they do not redefine `Event`.
- See [ADR-0008](./ADR-0008-plugin-marketplace.md) for distribution.

## Related

- Package stub: `platform/collector-sdk/`
- [ADR-0400](./ADR-0400-digital-presence-observatory.md)
