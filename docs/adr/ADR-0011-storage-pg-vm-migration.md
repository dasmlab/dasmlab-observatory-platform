# ADR-0011: Storage migration (Postgres + VictoriaMetrics)

**Status:** Proposed (Track B scaffold)  
**Date:** 2026-08

## Decision

Migrate DPO (and siblings) from SQLite to:

- **Postgres** — entities, baselines, config, multi-tenant rows
- **VictoriaMetrics** — time-series metrics / score history

SQLite ships track-home MVP; compose + migrate notes land here before cutover.

## Non-goals

Blocking Track A on this migration.
