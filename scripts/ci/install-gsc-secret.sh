#!/usr/bin/env bash
# Install Google Search Console SA JSON into dpo-secrets (cluster).
# Usage: bash scripts/ci/install-gsc-secret.sh /path/to/sa.json
set -euo pipefail

JSON_FILE="${1:?usage: install-gsc-secret.sh /path/to/sa.json}"
NS="${DPO_NAMESPACE:-dpo-system}"
SECRET_NAME="${DPO_SECRET_NAME:-dpo-secrets}"

if [ ! -f "${JSON_FILE}" ]; then
  echo "ERROR: file not found: ${JSON_FILE}"
  exit 1
fi
if ! command -v oc >/dev/null 2>&1; then
  echo "ERROR: oc not in PATH"
  exit 1
fi
if ! command -v jq >/dev/null 2>&1; then
  echo "ERROR: jq required"
  exit 1
fi

email="$(jq -r '.client_email // empty' "${JSON_FILE}")"
if [ -z "${email}" ] || [ "${email}" = "null" ]; then
  echo "ERROR: not a service-account JSON (missing client_email)"
  exit 1
fi

echo "Installing GSC_CREDENTIALS_JSON for SA: ${email}"
echo "Target: secret/${SECRET_NAME} in ${NS}"

# Preserve existing keys when present
tmp="$(mktemp)"
trap 'rm -f "${tmp}"' EXIT

existing_github=""
existing_activity=""
if oc get secret "${SECRET_NAME}" -n "${NS}" >/dev/null 2>&1; then
  existing_github="$(oc get secret "${SECRET_NAME}" -n "${NS}" -o jsonpath='{.data.GITHUB_TOKEN}' 2>/dev/null | base64 -d || true)"
  existing_activity="$(oc get secret "${SECRET_NAME}" -n "${NS}" -o jsonpath='{.data.ACTIVITY_MACHINE_TOKEN}' 2>/dev/null | base64 -d || true)"
fi

args=(--from-file="GSC_CREDENTIALS_JSON=${JSON_FILE}")
if [ -n "${existing_github}" ]; then
  args+=(--from-literal="GITHUB_TOKEN=${existing_github}")
elif [ -f /home/dasm/gh_token ]; then
  args+=(--from-literal="GITHUB_TOKEN=$(tr -d '\n\r' </home/dasm/gh_token)")
fi
if [ -n "${existing_activity}" ]; then
  args+=(--from-literal="ACTIVITY_MACHINE_TOKEN=${existing_activity}")
fi

oc create secret generic "${SECRET_NAME}" -n "${NS}" \
  "${args[@]}" \
  --dry-run=client -o yaml | oc apply -f -

oc rollout restart deploy/dpo -n "${NS}" || true
echo "Done. Add ${email} in Search Console → Users if not already. Then: curl -X POST .../api/v1/collect/run"
