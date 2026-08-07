# Track home improvements with DPO

Technologies DASMLAB Inc. builds observatories the same way it explores technology: hands-on, trend-aware, legacy versus future. **DPO** is the mirror for `dasmlab_home` — the Engineering Knowledge Network.

## Loop

1. Ship content / hubs / static truth on home.
2. Collectors run (GSC, HAProxy edge bots, Activity, website, GitHub).
3. Scores + content spine update.
4. Freeze a **baseline**, ship the next home change, **diff**.

## Operators

```bash
# After home sitemap is live:
curl -X POST https://dpo-dasmlab.apps.2026-prod-1.ocp.dasmlab.org/api/v1/collect/run

# Freeze
curl -X POST -H 'Content-Type: application/json' \
  -d '{"label":"pre-home-2.0"}' \
  https://dpo-dasmlab.apps.2026-prod-1.ocp.dasmlab.org/api/v1/baseline

# Diff later
curl 'https://dpo-dasmlab.apps.2026-prod-1.ocp.dasmlab.org/api/v1/baseline/diff?from=pre-home-2.0&to=post-home-2.0'
```

## Secrets

See [GSC-SETUP.md](./GSC-SETUP.md). Edge sample: `scripts/ci/pull-haproxy-access-sample.sh` from the self-hosted runner (host-first HAProxy).
