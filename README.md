# DASMLAB Observatory Platform (DOP)

**Technologies DASMLAB Inc.** — engineering platforms that make complex systems observable, understandable, automatable, and ultimately self-improving.

| | |
|--|--|
| Branch | `2026-dop-v0` (also builds from `main`) |
| First product | **DPO** — Digital Presence Observatory |
| Image | `ghcr.io/dasmlab/dpo` |
| Route (prod-1) | https://dpo-dasmlab.apps.2026-prod-1.ocp.dasmlab.org |
| GitOps | `lmcdasm/dasmlab-live-cicd` → `clusters/*/dpo/live` |
| Argo CD | Application `dpo` in `openshift-gitops` |

## Start here

1. [Engineering Principles](./docs/standards/ENGINEERING-PRINCIPLES.md)
2. [ADR Index](./docs/adr/README.md)
3. [ADR-0001 Platform Vision](./docs/adr/ADR-0001-observatory-platform-vision.md)
4. [ADR-0400 DPO](./docs/adr/ADR-0400-digital-presence-observatory.md)
5. [Edge topology (HAProxy)](./docs/standards/EDGE-TOPOLOGY.md)

## Run locally

```bash
go build -o bin/dpo-api ./cmd/dpo-api
DPO_STATIC_DIR=./web DPO_DATA_DIR=./.data ./bin/dpo-api
# http://127.0.0.1:8080/healthz  and  /
```

## Deploy flow

1. Push to `main` → GitHub Actions on **dasmlab org self-hosted runner** (buildah → `ghcr.io/dasmlab/dpo`).
2. `gitops-deploy` writes rendered manifest into `dasmlab-live-cicd` `clusters/<cluster>/dpo/live/dpo-v*.yaml`.
3. Argo CD Application `dpo` auto-syncs namespace `dpo-system`.

Bootstrap Argo (once per cluster):

```bash
oc apply -f deploy/argocd/dpo-2026-prod-1.yaml
# ensure GHCR pull secret exists in dpo-system (see OPENSHIFT_GITOPS_SETUP.md pattern)
```
