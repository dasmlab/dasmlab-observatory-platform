# Observatory family (DOP) roadmap

Technologies DASMLAB Inc. ships observatories as a **family**, not isolated apps.

| Code | Product | ADR | Status |
|------|---------|-----|--------|
| DCO | Cloud Observatory | 0100 | Stub |
| DNO | Network Observatory | 0200 | Stub |
| DSO | Security Observatory | 0300 | Stub |
| DPO | Digital Presence Observatory | 0400 (+0401 GEO) | Live skeleton + track-home |
| DAO | AI Observatory | 0500 | Stub |
| DIO | Infrastructure Observatory | 0600 | Stub |
| DAOps | DevOps Observatory | 0650 | Stub |
| DUO | Unified Observatory | 0700 | Stub |

Shared platform: collector SDK, scores, UI shell, edge topology, GitOps via `dasmlab-live-cicd`, **self-hosted CI buildah → GHCR**.

Implementation order: finish DPO track-home acceptance, then twin/storage (0009–0011), then GEO/DAO, then sibling products.
