# Behind the Design: first campaign — DASMLAB 2.0 launch

**Campaign:** `dasmlab-2.0-launch`  
**ADRs:** Observatory [ADR-0402](../adr/ADR-0402-campaign-orchestration.md), Home [ADR-002](../../../dasmlab_home/docs/ADR-002-LAUNCH-SURFACE.md)  
**Channels:** [CAMPAIGN-CHANNELS.md](../CAMPAIGN-CHANNELS.md)

## What we wanted to prove

Not “another social scheduler.” That we can **observe digital presence** (DPO) and use that same system to **create measurable presence** — then watch GSC, edge bots, Activity UTMs, and DUO impact move.

## How we got here

1. **Family runtimes** — DCO/DSO/DNO/DAO/DAOps/DIO/DUO left scaffold and emitted live F1/F2 scores.
2. **GSC live** — Search Console domain property + SA; search metrics left demo mode.
3. **ADR-0402** — Campaign orchestration as DPO F3: dry-run → arm → gated send; channel adapters; `campaign_rendered` / `content_published` events.
4. **Channel matrix** — Wave-1 social + email (Resend) + SMS (Twilio) + `web_slash`; accounts opened per checklist; missing secrets stay `blocked` with copy-paste path.
5. **Slash page** — `https://dasmlab.org/launch` (alias `/2.0`) — brand-first composition, UTM `dasmlab-2.0-launch`.
6. **Dry-run** — `POST /api/v1/campaigns/dasmlab-2.0-launch/render` produces per-channel bodies (X 280, SMS 160, LinkedIn long-form, …).
7. **Baselines** — `pre-launch` → ship/share → collect → `post-launch` → `/api/v1/baseline/diff`.
8. **DUO** — Recommend cites the campaign when presence/reachability is the weakest live signal.

## Operator loop

```bash
BASE=https://dpo-dasmlab.apps.2026-prod-1.ocp.dasmlab.org
curl -sk $BASE/api/v1/channels | jq '.[]|{id,live_status,can_send}'
curl -sk -X POST $BASE/api/v1/campaigns/dasmlab-2.0-launch/render | jq '.channels[]|{channel_id,status,body:.artifacts[0].body}'
curl -sk -X POST $BASE/api/v1/baseline -H 'content-type: application/json' -d '{"label":"pre-launch"}'
# ship /launch on home, share dry-run copy or arm+send where creds exist
curl -sk -X POST $BASE/api/v1/collect/run
curl -sk -X POST $BASE/api/v1/baseline -H 'content-type: application/json' -d '{"label":"post-launch"}'
curl -sk "$BASE/api/v1/baseline/diff?from=pre-launch&to=post-launch" | jq .
curl -sk $BASE/api/v1/duo/recommend | jq .
```

## Meta

We used DPO to get the word out about DPO (and DASMLAB 2.0). That loop — invent observability, act, measure again — is the product.
