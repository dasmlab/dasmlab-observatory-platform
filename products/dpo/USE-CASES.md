# DPO use cases

## F1 — Content spine + baseline diff

**Given** pilot tenant `dasmlab.org` and DOP collectors/stubs registered  
**When** the F1 path runs (live or scaffold stub)  
**Then** novel signal **Topic Coverage** is visible via Family API `features` and product scores  
**And** five-question focus covers: What exists?, What changed?, Why did it change?, What should I do?  
**And** chain role holds: Emits presence health + baseline labels for DUO and Research

Proof: Sitemap path entities joined to Activity/GSC/edge; POST /baseline + GET /baseline/diff

## F2 — Citation / crawl story

**Given** same tenant  
**When** the F2 path runs  
**Then** novel signal **AI Citation Velocity inputs** is visible  
**And** five-question focus covers: What exists?, What will happen?, What should I do?  
**And** chain role holds: Feeds Research Citation Index and DUO presence line

Proof: Edge UA classes (Googlebot, GPTBot, ClaudeBot, …) + GEO prompt pack scaffold

## Feeds E2E story

See [docs/plans/E2E-PLATFORM-STORY.md](../../docs/plans/E2E-PLATFORM-STORY.md) — After home ship → DPO scores/baselines → feeds DUO + Research
