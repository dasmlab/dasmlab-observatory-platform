# ADR-0402: Campaign orchestration (DPO F3)

**Status:** Accepted  
**Date:** 2026-08-08  
**Product:** DPO (Digital Presence Observatory)  
**Extends:** ADR-0400, ADR-0401, ADR-0004, ADR-0006, ADR-0008  
**Innovation gate:** ADR-9999

## Context

DASMLAB 2.0 needs a first *real* presence campaign: announce the launch, drive traffic to a slash page, and prove that DPO observes the presence it helped create. Commodity social schedulers fail ADR-9999 (another marketing dashboard). Campaign orchestration must be **presence engineering**: plan → channel-native artifacts → dry-run → gated send → `content_published` events → collect → baseline → DUO recommend.

## Decision

1. Campaign orchestration is a **DPO feature (F3)**, not a ninth family product.
2. Outbound channels are **plugins** implementing `ChannelAdapter` (ADR-0008). Official APIs/SDKs only — no browser scraping.
3. Default mode is **`dry_run`**. Live send requires explicit **arm** + channel credentials. Missing creds → honest `blocked` + manual copy path (`mode=manual` receipts).
4. Channel kinds include **social**, **sms**, **email**, and **web** (slash page). SMS requires opt-in lists (no cold SMS).
5. Campaign `type` enum exists (`launch` | `presence` | `research` | …); this wave implements **`launch`** only.
6. First seeded campaign id: **`dasmlab-2.0-launch`**.

## Model

```text
Campaign { id, tenant, type, title, brief, utm, status, channels[] }
ChannelPlan { channel_id, status, credentials_ref?, artifacts[] }
Artifact { format, body, media_refs?, char_limit, preview_html }
Receipt { channel_id, mode: dry_run|manual|sent, external_id?, at }
```

Statuses: `draft` → `dry_run` → `armed` → `sent` | `partial`.

## Events (ADR-0004)

| Metric / type | When | Dims |
|---------------|------|------|
| `campaign_rendered` | after `/render` | campaign, channel, mode=dry_run |
| `content_published` | after send or manual receipt | campaign, channel, mode=sent\|manual |

## APIs

- `GET /api/v1/channels` — catalog + readiness from env
- `GET/POST /api/v1/campaigns`
- `GET /api/v1/campaigns/{id}`
- `POST /api/v1/campaigns/{id}/render`
- `POST /api/v1/campaigns/{id}/arm`
- `POST /api/v1/campaigns/{id}/send`

## Innovation-gate proof

| Gate | How campaigns satisfy it |
|------|---------------------------|
| Measure | Channel readiness + render/send receipts as presence telemetry |
| Correlate | Campaign id ↔ Activity UTM ↔ GSC/edge after `/launch` |
| Predict | (later) expected reach from prior campaign baselines |
| Visualize | Per-channel dry-run previews in DPO UI |
| Automate | Gated send via official APIs; never silent scrape |

## Non-goals

- Agency multi-tenant marketing suite
- Paid ad buy APIs this wave
- Replacing CDN-mgr (home ADR-001) as orchestrator

## Consequences

- Secrets for channels land in `dpo-secrets` (same pattern as GSC).
- Home `/launch` is the web_slash artifact (home ADR-002).
- Channel matrix: `docs/CAMPAIGN-CHANNELS.md`.
- Story: `docs/stories/DASMLAB-2.0-FIRST-CAMPAIGN.md`.
