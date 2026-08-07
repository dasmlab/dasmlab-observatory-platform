# ADR-0401: GEO / AI visibility (scaffold)

**Status:** Proposed (Track B — not deployed until Track A acceptance)  
**Product:** DPO (Digital Presence Observatory)  
**Date:** 2026-08

## Decision

Measure generative-engine visibility with **official APIs and first-party signals only** (no scraping of chatbot UIs). Pair with prompt packs in `docs/PROMPTS.md` and future collectors under `internal/collectors/ai/`.

## Scope (later)

- Citation extract from permitted APIs / partner feeds
- Prompt pack versioning for regression of “how AIs describe DASMLAB”
- Scores: Citation Velocity, Topic Coverage (see `docs/SCORES.md`)

## Non-goals (now)

GEO scrapers, paid SERP farms, or silent third-party browser automation in prod.
