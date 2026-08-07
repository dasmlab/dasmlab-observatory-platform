# ADR-0010: OpenSearch

**Status:** Proposed (Track B scaffold)  
**Date:** 2026-08

## Decision

Use OpenSearch for full-text / log / event search across observatory telemetry. SQLite remains the track-home store until PG+VM+OpenSearch migration (ADR-0011).

## Package scaffold

`platform/search-sdk` — empty module until cutover.
