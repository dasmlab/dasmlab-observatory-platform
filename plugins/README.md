# Plugins

Collector / integration plugins. Same Observatory pipeline; different sources.

| Plugin | Feeds | Status |
|--------|-------|--------|
| google (GSC) | DPO | live (in-tree collector) |
| github | DPO / authority | live (in-tree) |
| haproxy-edge | DPO bots | live (in-tree + CI pull) |
| surfing-activity | DPO engagement | live (bridge) |
| kubernetes | DCO | planned |
| aws / azure | DCO | planned |
| linkedin | DPO campaign / GEO | dry-run live; send when `LINKEDIN_*` secrets |
| bluesky / mastodon | DPO campaign | dry-run + send when tokens present |
| x_twitter / meta / reddit | DPO campaign | dry-run; send stubs |
| resend (email) | DPO campaign | dry-run; send when `RESEND_*` |
| twilio (sms) | DPO campaign | dry-run; send when Twilio + opt-in `TWILIO_TO_NUMBER` |
| chatgpt / perplexity / gemini | DPO/DAO GEO | planned (official APIs only) |

Channel adapters: `internal/campaign/` (ADR-0402). Matrix: `docs/CAMPAIGN-CHANNELS.md`.

In-tree collectors under `internal/collectors/` until packages move to `plugins/<name>`.
