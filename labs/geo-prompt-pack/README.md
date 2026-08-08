# Lab: GEO prompt pack (sunset)

**Hypothesis:** Versioned prompts about DASMLAB produce measurable citation/answer variance across official AI APIs (no UI scraping).  
**Status:** Declared experiment  
**Sunset:** 2026-11-07 (promote to DAO/DPO collectors or delete)

## Pack

See [docs/PROMPTS.md](../../docs/PROMPTS.md) pack v0.

## ADR-9999

Temporarily Labs-exempt from full gate if results are unpublished; must not scrape consumer chat UIs.

## Outcomes store

Record JSON outcomes under `labs/geo-prompt-pack/outcomes/` (gitignored samples OK). Schema: `{prompt_id, engine, date, cited_urls[], notes}`.
