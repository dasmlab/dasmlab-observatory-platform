# E2E demo script (~10 minutes)

Target: https://dpo-dasmlab.apps.2026-prod-1.ocp.dasmlab.org (or local `dpo-api`).

## Prep

```bash
curl -sk https://dpo-dasmlab.apps.2026-prod-1.ocp.dasmlab.org/healthz | jq .
curl -sk https://dpo-dasmlab.apps.2026-prod-1.ocp.dasmlab.org/api/v1/family | jq '.products[]|{code,status,features:[.features[].name]}'
curl -sk https://dpo-dasmlab.apps.2026-prod-1.ocp.dasmlab.org/api/v1/products | jq '.[]|{code,status,mode,scores}'
```

## Walkthrough

| Min | Step | Show |
|-----|------|------|
| 0 | Open **Family** tab | All DXX cards with ≥2 features; live scores strip from `/api/v1/products` |
| 1 | **DPO** card | F1 content spine / F2 citation crawl — status live |
| 2 | DPO Engineering | Collectors include dco/dso/dno/dao/daops/dio + GSC; run collectors |
| 3 | Freeze baseline | Label `demo-pre`; note score |
| 4 | **DUO** tab | Business → Engineering → Operational chain from live/demo sources |
| 5 | DUO recommend | Evidence cites ≥2 products (dynamic weakest-signal action) |
| 6 | Sibling cards | Each product API returns F1/F2; modes `live` or `demo` (not hardcoded scaffold) |
| 7 | Story close | Research flywheel: Citation Index backlog → DPO authority |

## API checklist

```bash
BASE=https://dpo-dasmlab.apps.2026-prod-1.ocp.dasmlab.org
curl -sk -X POST $BASE/api/v1/collect/run
sleep 8
curl -sk $BASE/api/v1/products | jq '.[]|{code,mode,scores}'
curl -sk $BASE/api/v1/duo/impact | jq '{business:.business.score,engineering:.engineering.score,operational:.operational.score,modes:[.sources[].mode]|unique}'
curl -sk $BASE/api/v1/duo/recommend | jq '{title,confidence,evidence}'
```

## Pass criteria

- Family lists features for every product
- `/api/v1/products/{code}` returns F1/F2 for dco/dso/dno/dao/daops/dio
- DUO impact sources are `live` or `demo` after collect (scaffold only as pre-collect fallback)
- Recommendation includes evidence array length ≥2 and ADR-0006 fields
- No commodity-only hero (CPU/CTR/CVE) as the DUO title story

## Diagrams

See `diagrams/e2e-chain.svg` and [E2E-PLATFORM-STORY.md](./E2E-PLATFORM-STORY.md).
