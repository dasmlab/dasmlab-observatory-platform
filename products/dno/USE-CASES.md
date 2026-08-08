# DNO use cases

## F1 — Service Reachability

**Given** pilot tenant `dasmlab.org` and DOP collectors/stubs registered  
**When** the F1 path runs (live or scaffold stub)  
**Then** novel signal **Service Reachability** is visible via Family API `features` and product scores  
**And** five-question focus covers: What exists?, What changed?, What will happen?, What should I do?  
**And** chain role holds: Outage narratives for DUO; complements DPO edge bots

Proof: Probe HAProxy→Route→Service for dasmlab.org / DPO FQDN

## F2 — Intent Compliance

**Given** same tenant  
**When** the F2 path runs  
**Then** novel signal **Intent Compliance** is visible  
**And** five-question focus covers: What exists?, Why did it change?, What should I do?  
**And** chain role holds: Policy complexity → DUO operational impact

Proof: CERT*/ACL intent vs ensure-prod-cert outcomes

## Feeds E2E story

See [docs/plans/E2E-PLATFORM-STORY.md](../../docs/plans/E2E-PLATFORM-STORY.md) — Edge path health underpins DPO public truth + DUO outages
