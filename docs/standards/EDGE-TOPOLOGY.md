# Edge topology (keep in mind)

Internet traffic to DASMLAB apps generally flows:

```text
Client → router/firewall → HAProxy (edge TLS / FQDN / optional auth)
       → OpenShift Routes (cluster apps)
```

## Facts for agents

- **HAProxy is the public edge**, not the in-pod `nginx` used only by `dasmlab_home` to serve the Quasar SPA.
- **Host is source of truth** for HAProxy (see cheapcloud `docs/haproxy-host-first.md`): live config on the proxy (`dasmlab-internal/new_haproxy` / `10.20.1.10`), git is a snapshot.
- **Certs for public FQDNs** are often ensured via HAProxy `CERT*=fqdn` (e.g. `ensure-prod-cert.sh`), not only OpenShift Route edge certs.
- Some UIs sit behind **HAProxy basic auth**; service-to-service calls should prefer **cluster Routes / ClusterIP** (as surfing → cheapcloud already does).
- **DPO edge / bot telemetry** should prefer HAProxy (and/or CF) access logs over assuming every app has its own nginx access log. Home’s container nginx is an SPA static server, not the edge.

## Implications

- Fixes for public paths (`/sitemap.xml`, `/robots.txt`, vanity hosts, certs) may need **HAProxy and/or OpenShift**, not only app config.
- Do not introduce nginx into Go services (DPO, Surfing, cheapcloud) for edge concerns.
- **CI must ensure HAProxy CERT entries** for vanity / preview FQDNs (same pattern as cheapcloud / home):
  - `scripts/ci/ensure-prod-cert.sh <fqdn>`
  - `scripts/ci/ensure-preview-cert.sh <fqdn>` (dev/preview; respects `SKIP_PREVIEW_CERT=true`)
- Host-first: never `git pull` / reset the live HAProxy tree from CI — only edit `runme.sh` on the proxy and recreate.
