#!/usr/bin/env bash
# One-time OpenShift bootstrap for DPO on the current oc context.
set -euo pipefail
ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
NS=dpo-system
GHCR_TOKEN="${GHCR_TOKEN:-}"
if [ -z "${GHCR_TOKEN}" ] && [ -f /home/dasm/gh_token ]; then
  GHCR_TOKEN="$(tr -d '\n\r' < /home/dasm/gh_token)"
fi

oc create namespace "${NS}" --dry-run=client -o yaml | oc apply -f -
oc apply -f "${ROOT}/deploy/argocd/namespace-rbac.yaml"

if [ -n "${GHCR_TOKEN}" ]; then
  oc create secret docker-registry dasmlab-ghcr-pull \
    --docker-server=ghcr.io \
    --docker-username="${GHCR_USERNAME:-lmcdasm}" \
    --docker-password="${GHCR_TOKEN}" \
    --docker-email=dasmlab-bot@dasmlab.org \
    -n "${NS}" \
    --dry-run=client -o yaml | oc apply -f -
fi

oc apply -f "${ROOT}/deploy/argocd/dpo-2026-prod-1.yaml"
# Fixed UID 65532 in the image requires anyuid (same pattern as surfing-service).
oc adm policy add-scc-to-user anyuid -z dpo -n "${NS}" || true

echo "Argo Application:"
oc get application.argoproj.io dpo -n openshift-gitops
echo "Pods/route:"
oc get pods,route -n "${NS}"
