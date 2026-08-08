# E2E demo script (~10 minutes)

Target: https://dpo-dasmlab.apps.2026-prod-1.ocp.dasmlab.org (or local `dpo-api`).

## Prep

```bash
curl -s https://dpo-dasmlab.apps.2026-prod-1.ocp.dasmlab.org/healthz | jq .
curl -s https://dpo-dasmlab.apps.2026-prod-1.ocp.dasmlab.org/api/v1/family | jq '.products[]|{code,status,features:[.features[].name]}'
```

## Walkthrough

| Min | Step | Show |
|-----|------|------|
| 0 | Open **Family** tab | All DXX cards with ≥2 features each; five questions; innovation gate |
| 1 | **DPO** card | F1 content spine / F2 citation crawl — status live |
| 2 | DPO Engineering | Collectors + bots + content spine; run collectors |
| 3 | Freeze baseline | Label `demo-pre`; note score |
| 4 | **DUO** tab | Business → Engineering → Operational chain |
| 5 | DUO recommend | Evidence cites ≥2 products (DPO + DSO/DNO) |
| 6 | Sibling cards | DCO/DSO/DNO/DAO/DAOps/DIO features = existence proofs (scaffold scores in DUO sources) |
| 7 | Story close | Research flywheel: Citation Index backlog → DPO authority |

## API checklist

```bash
BASE=https://dpo-dasmlab.apps.2026-prod-1.ocp.dasmlab.org
curl -s $BASE/api/v1/duo/impact | jq '{business:.business.score,engineering:.engineering.score,operational:.operational.score,n_sources:(.sources|length)}'
curl -s $BASE/api/v1/duo/recommend | jq '{title,confidence,evidence}'
curl -s $BASE/api/v1/content | jq .
curl -s -X POST $BASE/api/v1/collect/run
```

## Pass criteria

- Family lists features for every product (not empty READMEs only)
- DUO impact has live or scaffold sources from ≥2 products
- Recommendation includes evidence array length ≥2 and ADR-0006 fields
- No commodity-only hero (CPU/CTR/CVE) as the DUO title story

## Diagrams

See `diagrams/e2e-chain.svg` and [E2E-PLATFORM-STORY.md](./E2E-PLATFORM-STORY.md).
