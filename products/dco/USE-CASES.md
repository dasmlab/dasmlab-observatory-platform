# DCO use cases

## F1 — Deploy Confidence

**Given** pilot tenant `dasmlab.org` and DOP collectors/stubs registered  
**When** the F1 path runs (live or scaffold stub)  
**Then** novel signal **Deployment Confidence** is visible via Family API `features` and product scores  
**And** five-question focus covers: What changed?, Why did it change?, What will happen?, What should I do?  
**And** chain role holds: Same window as DPO ship → cluster cost of presence change

Proof: Correlate GitOps image rollout ↔ pod ready/restarts for dpo-system

## F2 — Operational Complexity

**Given** same tenant  
**When** the F2 path runs  
**Then** novel signal **Operational Complexity** is visible  
**And** five-question focus covers: What exists?, What changed?, What should I do?  
**And** chain role holds: Complexity index → DUO engineering impact

Proof: Count Deployments, PVCs, Routes, Secrets, CronJobs in watched ns

## Feeds E2E story

See [docs/plans/E2E-PLATFORM-STORY.md](../../docs/plans/E2E-PLATFORM-STORY.md) — DPO image change window → DCO deploy confidence
