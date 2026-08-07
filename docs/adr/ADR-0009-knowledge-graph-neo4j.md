# ADR-0009: Knowledge graph (Neo4j)

**Status:** Proposed (Track B scaffold)  
**Date:** 2026-08

## Decision

Introduce Neo4j as the **entity/relationship twin** for observatory products (content paths, services, clouds, identities). Track-home MVP stays on SQLite `entities` / `path_daily`; this ADR locks the cutover shape.

## Package scaffold

`platform/graph-sdk` — empty module until cutover.
