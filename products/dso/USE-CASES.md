# DSO use cases

## F1 — Attack Surface Evolution

**Given** pilot tenant `dasmlab.org` and DOP collectors/stubs registered  
**When** the F1 path runs (live or scaffold stub)  
**Then** novel signal **Attack Surface Evolution** is visible via Family API `features` and product scores  
**And** five-question focus covers: What exists?, What changed?, Why did it change?, What will happen?  
**And** chain role holds: FQDN/vanity change → surface delta → DUO risk

Proof: Time-diff public Routes/Services/ports for observatory namespaces

## F2 — Secrets Hygiene

**Given** same tenant  
**When** the F2 path runs  
**Then** novel signal **Secrets Hygiene** is visible  
**And** five-question focus covers: What exists?, Why did it change?, What should I do?  
**And** chain role holds: DUO recommends credential fill; stops silent demo SEO

Proof: Optional/missing keys in dpo-secrets (e.g. GSC) = hygiene debt score

## Feeds E2E story

See [docs/plans/E2E-PLATFORM-STORY.md](../../docs/plans/E2E-PLATFORM-STORY.md) — DPO vanity/CERT → DSO surface + hygiene
