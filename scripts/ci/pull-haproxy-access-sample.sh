#!/usr/bin/env bash
# Pull a HAProxy access-log sample from the edge host (host-first).
# Same SSH pattern as ensure-prod-cert.sh — run from self-hosted CI / lab runner.
set -euo pipefail

PROXY_HOST="${PREVIEW_PROXY_HOST:-10.20.1.10}"
PROXY_USER="${PREVIEW_PROXY_USER:-dasm}"
# Override if your proxy stores logs elsewhere.
REMOTE_LOG="${HAPROXY_ACCESS_LOG:-/var/log/haproxy.log}"
OUT="${1:-./.data/edge/access.log}"
BYTES="${HAPROXY_LOG_BYTES:-5242880}" # last ~5MiB

mkdir -p "$(dirname "$OUT")"

ssh -o BatchMode=yes -o StrictHostKeyChecking=accept-new \
  "${PROXY_USER}@${PROXY_HOST}" bash -s -- "${REMOTE_LOG}" "${BYTES}" <<'EOS' > "${OUT}.tmp"
set -euo pipefail
LOG="$1"
N="$2"
if [ ! -f "$LOG" ]; then
  # Common alternates
  for c in /var/log/haproxy/haproxy.log /var/log/haproxy.log /home/dasm/dasmlab-internal/new_haproxy/logs/access.log; do
    if [ -f "$c" ]; then LOG="$c"; break; fi
  done
fi
if [ ! -f "$LOG" ]; then
  echo "ERROR: no HAProxy access log found on proxy" >&2
  exit 1
fi
echo "Using $LOG" >&2
tail -c "$N" "$LOG" || true
EOS

mv "${OUT}.tmp" "$OUT"
echo "Wrote $(wc -c < "$OUT") bytes to $OUT"
