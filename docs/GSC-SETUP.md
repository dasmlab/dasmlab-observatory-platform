# GSC credentials — plain English

## What it is

A **JSON key file** Google gives you for a *service account*. It looks like:

```json
{
  "type": "service_account",
  "project_id": "your-project",
  "private_key_id": "...",
  "private_key": "-----BEGIN PRIVATE KEY-----\n...\n-----END PRIVATE KEY-----\n",
  "client_email": "dpo-gsc@YOUR-PROJECT.iam.gserviceaccount.com",
  "client_id": "...",
  "token_uri": "https://oauth2.googleapis.com/token"
}
```

That whole file (as one string) is what we call **`GSC_CREDENTIALS_JSON`**.

I **cannot invent** a valid key — only Google can issue it. Once you download the file once, I can install it into the cluster for you.

## Where it goes

Kubernetes secret in the DPO namespace:

| | |
|--|--|
| Namespace | `dpo-system` |
| Secret name | `dpo-secrets` |
| Secret **key** | `GSC_CREDENTIALS_JSON` |
| Pod env | already wired in the Deployment |

Today that secret only has `GITHUB_TOKEN` and `ACTIVITY_MACHINE_TOKEN` — **no GSC key yet**, so GSC stays in demo mode.

## What you do (once)

1. Google Cloud Console → create/reuse a project → enable **Search Console API**.
2. Create a service account → **Keys → Add key → JSON** → save as e.g. `~/dpo-gsc-sa.json`.
3. Search Console → property **`dasmlab.org`** (domain = `sc-domain:dasmlab.org`) → **Users and permissions** → add the `client_email` from that JSON (Restricted is enough).
4. DPO must use **`GSC_SITE_URL=sc-domain:dasmlab.org`** (matches the domain property). Do **not** use `https://dasmlab.org/` unless that URL-prefix property also has the SA as a user.
5. Tell me the path to the file **or** run:

```bash
bash /home/dasm/dasmlab-observatory-platform/scripts/ci/install-gsc-secret.sh ~/dpo-gsc-sa.json
```

That script patches `dpo-secrets` and rolls the DPO pod. After that, collectors use live GSC (`mode=live`).

## Not this

- Not a Google password
- Not an API key string alone
- Not something to commit into git
