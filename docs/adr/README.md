# ADR Index — DASMLAB Observatory Platform

Numbering is **platform-global**. Do not reuse these numbers in `dasmlab_home` (home keeps local ADRs such as CDN-mgr GEO).

## Platform core

| ADR | Title | Status |
|-----|-------|--------|
| [0001](./ADR-0001-observatory-platform-vision.md) | Observatory Platform Vision | Accepted |
| [0002](./ADR-0002-common-reference-architecture.md) | Common Reference Architecture | Accepted |
| [0003](./ADR-0003-collector-sdk.md) | Collector SDK | Accepted |
| [0004](./ADR-0004-analytics-correlation-engine.md) | Analytics & Correlation Engine | Accepted |
| [0005](./ADR-0005-shared-ui-ux-framework.md) | Shared UI/UX Framework | Accepted |
| [0006](./ADR-0006-ai-recommendation-engine.md) | AI & Recommendation Engine | Accepted |
| [0007](./ADR-0007-security-multi-tenancy.md) | Security & Multi-Tenancy | Accepted |
| [0008](./ADR-0008-plugin-marketplace.md) | Plugin Marketplace | Accepted |
| [9999](./ADR-9999-innovation-principles.md) | Innovation Principles (North Star) | Accepted |

## Product specializations

| ADR | Product | Status |
|-----|---------|--------|
| [0100](./ADR-0100-cloud-observatory.md) | DCO — Cloud | Stub |
| [0200](./ADR-0200-network-observatory.md) | DNO — Network | Stub |
| [0300](./ADR-0300-security-observatory.md) | DSO — Security | Stub |
| [0400](./ADR-0400-digital-presence-observatory.md) | DPO — Digital Presence | Accepted (MVP scope) |
| [0500](./ADR-0500-ai-observatory.md) | DAO — AI | Stub |
| [0600](./ADR-0600-infrastructure-observatory.md) | DIO — Infrastructure | Stub |
| [0650](./ADR-0650-devops-observatory.md) | DAOps — DevOps | Stub |
| [0700](./ADR-0700-unified-observatory.md) | DUO — Unified | Stub |

## How agents should consume ADRs

1. Read **0001 + 9999** before proposing features.
2. Implement collectors against **0003**; storage/APIs against **0002**.
3. Product work: read the product ADR (e.g. **0400**) plus **0002–0007**.
4. Reject work that fails the **9999** gate unless explicitly marked experimental Labs.
