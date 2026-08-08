# Product runtime contract

All DOP family products (except DPO historical collectors) ship inside **one** `dpo-api` binary.

## Layout

```
internal/products/{code}/   # collector + score helpers
internal/products/          # Snapshot types + registry + store projection
GET /api/v1/products
GET /api/v1/products/{code}
```

## Metric names (must match DUO sources)

| Product | Metrics |
|---------|---------|
| dco | `deploy_confidence`, `operational_complexity` |
| dso | `attack_surface_evolution`, `secrets_hygiene` |
| dno | `service_reachability`, `intent_compliance` |
| dao | `prompt_effectiveness`, `ai_trust_score` |
| daops | `delivery_confidence`, `toil_ratio` |
| dio | `capacity_confidence`, `failover_readiness` |
| dpo | `overall` (via score engine) |

## Event dims

Every product metric event **must** set `dims.mode` to `live` or `demo` (never leave blank). DUO and Family use that mode; hardcoded `scaffold` constants are forbidden once a collector is registered.

## Collector naming

Collector `Name()` equals product code (`dco`, `dso`, …) so `/api/v1/sources/status` and product APIs align.

## Family status

`live` when the product has emitted events with `mode=live|demo` in the last 48h. Catalog static status is updated to `live` for shipped runtimes.
