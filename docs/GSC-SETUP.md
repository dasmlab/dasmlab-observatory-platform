# Google Search Console setup (DPO)

## Why

Live SEO scores require Search Console analytics — not demo impressions.

## Steps (2026)

1. In [Google Cloud Console](https://console.cloud.google.com), create/reuse a project.
2. Enable **Google Search Console API** (`searchconsole.googleapis.com` / webmasters).
3. Create a **service account** (no GCP roles required for GSC read).
4. Create a **JSON key**; store as cluster secret (never commit):

```bash
# From a machine with oc + the key file:
oc create secret generic dpo-secrets -n dpo-system \
  --from-file=GSC_CREDENTIALS_JSON=/path/to/sa.json \
  --from-literal=GITHUB_TOKEN="$(tr -d '\n\r' < /home/dasm/gh_token)" \
  --from-literal=ACTIVITY_MACHINE_TOKEN="$(oc get secret dpo-secrets -n dpo-system -o jsonpath='{.data.ACTIVITY_MACHINE_TOKEN}' | base64 -d)" \
  --dry-run=client -o yaml | oc apply -f -
```

Or patch an existing secret:

```bash
oc patch secret dpo-secrets -n dpo-system --type merge -p \
  "{\"stringData\":{\"GSC_CREDENTIALS_JSON\":$(jq -c -R -s . </path/to/sa.json)}}"
```

5. In [Search Console](https://search.google.com/search-console) → property `dasmlab.org` (URL-prefix `https://dasmlab.org/` or domain) → **Settings → Users and permissions → Add user**.
6. Paste the SA email (`client_email` in the JSON). Permission: **Restricted** (read-only) is enough.
7. Set env `GSC_SITE_URL` to the exact property string (default `https://dasmlab.org/`).
8. Rollout DPO / run collectors. `mode=live` appears on GSC events when the secret is present.

Without the secret, the collector stays in **demo** mode (healthy) unless `DPO_REQUIRE_LIVE_CREDS=1`.
