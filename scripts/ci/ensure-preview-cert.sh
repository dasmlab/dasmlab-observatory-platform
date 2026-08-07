#!/usr/bin/env bash
# Ensure HAProxy CERT for a DPO preview / dev FQDN (same host-first pattern).
# Wrapper kept separate so preview pipelines can call it without implying prod vanity hosts.
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
FQDN="${1:?usage: ensure-preview-cert.sh <fqdn>}"

if [ "${SKIP_PREVIEW_CERT:-}" = "true" ]; then
  echo "SKIP_PREVIEW_CERT=true — not updating HAProxy CERT for ${FQDN}"
  exit 0
fi

exec bash "${ROOT}/scripts/ci/ensure-prod-cert.sh" "${FQDN}"
