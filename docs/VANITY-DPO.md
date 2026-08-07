# Vanity host: dpo.dasmlab.org

Internet → lab edge → **HAProxy** (TLS via `CERT*=dpo.dasmlab.org`) → OCP Route backend.

CI already runs `scripts/ci/ensure-prod-cert.sh dpo.dasmlab.org` on the **self-hosted** runner (same host-first SSH as home/cheapcloud).

## Host-first checklist

1. HAProxy frontend ACL for `dpo.dasmlab.org`
2. Backend to cluster ingress / Route host (`dpo-dasmlab.apps.…`)
3. CERT ensure from runner; do not invent in-pod nginx for edge TLS
4. Confirm DNS A/CNAME at the edge, not only apps.* wildcard

See `docs/standards/EDGE-TOPOLOGY.md`.
