# DIO use cases

## F1 — Capacity Confidence

**Given** pilot tenant `dasmlab.org` and DOP collectors/stubs registered  
**When** the F1 path runs (live or scaffold stub)  
**Then** novel signal **Capacity Confidence** is visible via Family API `features` and product scores  
**And** five-question focus covers: What exists?, What will happen?, What should I do?  
**And** chain role holds: Underpins DCO recovery confidence

Proof: PVC/LVMS usage for observatory namespaces

## F2 — Failover Readiness

**Given** same tenant  
**When** the F2 path runs  
**Then** novel signal **Failover Readiness** is visible  
**And** five-question focus covers: What exists?, What changed?, What should I do?  
**And** chain role holds: Recovery input to DCO/DUO

Proof: GitOps drift check 2026-prod-1 vs 2026-prod-2-1

## Feeds E2E story

See [docs/plans/E2E-PLATFORM-STORY.md](../../docs/plans/E2E-PLATFORM-STORY.md) — Infra readiness under DCO recovery + DUO ops impact
