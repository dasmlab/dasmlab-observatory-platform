# ADR-0002 — Common Reference Architecture

**Status:** Accepted  
**Date:** 2026-08-07  
**Deciders:** Technologies DASMLAB Inc.  
**Depends on:** [ADR-0001](./ADR-0001-observatory-platform-vision.md)

## Context

Each Observatory must look identical at the architecture level so shared SDKs and services stay reusable. Without a common reference, products will reinvent collectors, storage, and dashboards.

## Decision

### Layered stack

```text
Applications   DUO · DPO · DCO · DAO · DSO · DNO · DAOps · DIO
─────────────
Platform       AI · Analytics · Graph · Workflow · Search · Collectors · Storage · Notifications
─────────────
Shared SDK     UI · REST/GraphQL · Plugins · Security · Licensing · CLI · Config
─────────────
Infrastructure PostgreSQL · VictoriaMetrics · OpenSearch · Neo4j · Redis · NATS · S3 · Keycloak · Kubernetes
```

Applications are **thin**. Domain logic beyond plugin packs and score definitions lives in platform services.

### Canonical pipeline (every Observatory)

```text
Collectors → Normalization → Correlation → Analytics → AI → Storage → API → Dashboard → Recommendations
```

### Collection control plane

| Component | Responsibility |
|-----------|----------------|
| **Scheduler** | Cron / interval / manual runs |
| **Collector Manager** | Registry, status, failures, versions, config, credentials refs |
| **Plugin Manager** | Load/enable/disable collector plugins |
| **Normalization Bus** | Convert plugin output to common `Event` schema |
| **Event Processor** | Enrichment, correlation, entity resolution |

### Normalized event schema (v0)

```json
{
  "schema_version": "1",
  "collector": "google",
  "type": "query",
  "timestamp": "2026-08-07T12:00:00Z",
  "tenant": "dasmlab.org",
  "entity": "kubernetes operator",
  "metric": "impression",
  "value": 1482,
  "unit": "count",
  "dims": {
    "page": "/projects/cheapcloud",
    "country": "CA"
  },
  "trace_id": "optional"
}
```

Rules:

- One logical observation → one event (or a batch of events with shared `trace_id`).
- `tenant` is required (pilot: `dasmlab.org`).
- Time series metrics also remote-write to VictoriaMetrics with labels mirroring `dims` + `collector` + `metric`.
- Entity identity resolution maps free-text `entity` to graph nodes (Phase 2+).

### Storage roles

| Store | Role | Phase |
|-------|------|-------|
| **PostgreSQL** | Entities, scan runs, config, score snapshots, recommendations | P1 |
| **VictoriaMetrics** | Time series gauges/counters/histograms | P1 |
| **OpenSearch** | Full-text search over responses, cites, docs, research | P2 |
| **Neo4j** | Knowledge / digital twin graph | P2 |
| **Redis** | Cache, rate limits, job locks | P2 as needed |
| **NATS** | Async event bus between services | P2 (P1 may use in-process + DB outbox) |
| **S3** | Raw payloads, log batches, research artifacts | P1 for edge log archives |

### Shared platform services

| Service | Owns |
|---------|------|
| **Identity** | Keycloak OIDC, orgs, teams, projects, RBAC, API keys |
| **Notifications** | Email, Slack, Discord, Teams, webhooks, SSE, push |
| **Analytics** | Aggregation, correlation rules, score computation |
| **Search** | Cross-entity search API |
| **Graph** | Neo4j projection + query API |
| **Recommendation** | Explainable recommendation objects |
| **Workflow** | Scheduled jobs, approval gates for automation (L6) |

No product owns a private analytics stack.

### API surface

| Protocol | Phase | Use |
|----------|-------|-----|
| REST `/api/v1` | P1 | Scores, sources, health, events query |
| GraphQL | P3 / late P2 | Flexible dashboard queries |
| WebSockets / SSE | P2–P3 | Live scan progress, alerts |
| MCP export | P3 | Agent-readable metrics |

### Knowledge graph (connective tissue)

Core node types: Person, Author, Project, Repository, Technology, Article, Video, Presentation, Issue, Research, Company, Website, SearchTerm, AIMention, VisitorCohort.

Edges encode authorship, depends-on, cites, ranks-for, mentions, drives-traffic-to, etc. Products project domain entities into this shared graph.

### Observatory AI Analyst

Each product may expose an AI Analyst that:

1. Accepts a natural-language question (Five Questions).
2. Pulls evidence from metrics, events, graph, and docs.
3. Returns explanation + [Recommendation objects](./ADR-0006-ai-recommendation-engine.md).

### Deploy topology (pilot)

- OpenShift on DASMLAB `2026-prod-1` (or successor), namespace per product + shared `dop-platform`.
- GitOps via `dasmlab-live-cicd` when ready.
- Local: docker-compose for PG + VM optional.

## Consequences

- Collector implementations follow [ADR-0003](./ADR-0003-collector-sdk.md).
- UI shells follow [ADR-0005](./ADR-0005-shared-ui-ux-framework.md).
- Security/tenancy follow [ADR-0007](./ADR-0007-security-multi-tenancy.md).
- P1 products (DPO) must ship with PG + VM only; Neo4j/OpenSearch are not MVP blockers.

## Related

- [ADR-0003](./ADR-0003-collector-sdk.md)
- [ADR-0004](./ADR-0004-analytics-correlation-engine.md)
- [ADR-0400](./ADR-0400-digital-presence-observatory.md)
