# DAOPS use cases

## F1 — Delivery Confidence

**Given** pilot tenant `dasmlab.org` and DOP collectors/stubs registered  
**When** the F1 path runs (live or scaffold stub)  
**Then** novel signal **Delivery Confidence** is visible via Family API `features` and product scores  
**And** five-question focus covers: What changed?, Why did it change?, What will happen?, What should I do?  
**And** chain role holds: Explains risky DPO/DCO windows

Proof: CX success rate / queue time on bld-249 for observatory + home

## F2 — Toil Ratio

**Given** same tenant  
**When** the F2 path runs  
**Then** novel signal **Toil Ratio** is visible  
**And** five-question focus covers: What exists?, Why did it change?, What should I do?  
**And** chain role holds: Friction → DUO engineering impact

Proof: Manual workflow_dispatch / cert ensures vs auto GitOps publishes

## Feeds E2E story

See [docs/plans/E2E-PLATFORM-STORY.md](../../docs/plans/E2E-PLATFORM-STORY.md) — CI pain explains deploy/presence risk windows
