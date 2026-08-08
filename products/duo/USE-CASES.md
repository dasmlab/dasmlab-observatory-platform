# DUO use cases

## F1 — Impact chain view

**Given** pilot tenant `dasmlab.org` and DOP collectors/stubs registered  
**When** the F1 path runs (live or scaffold stub)  
**Then** novel signal **Business→Engineering→Operational** is visible via Family API `features` and product scores  
**And** five-question focus covers: What exists?, What changed?, Why did it change?, What will happen?, What should I do?  
**And** chain role holds: Executive ending of E2E demo

Proof: GET /api/v1/duo/impact composes DPO + sibling stub scores

## F2 — Recommended action card

**Given** same tenant  
**When** the F2 path runs  
**Then** novel signal **Recommended Actions** is visible  
**And** five-question focus covers: Why did it change?, What should I do?  
**And** chain role holds: Closes five-question loop for leadership

Proof: GET /api/v1/duo/recommend — ADR-0006 shape with ≥2 product evidence

## Feeds E2E story

See [docs/plans/E2E-PLATFORM-STORY.md](../../docs/plans/E2E-PLATFORM-STORY.md) — Consumes all lanes; produces one action
