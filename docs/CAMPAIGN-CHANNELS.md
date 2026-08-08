# Campaign channels — account & API matrix

Wave-1 channels for DPO Campaign Orchestration (ADR-0402).  
**Dry-run works without accounts.** Live send needs the checklist below.

| id | kind | Official API / SDK | Env / secrets | Account checklist | Live status |
|----|------|-------------------|---------------|-------------------|-------------|
| `web_slash` | web | Home GitOps publish | — | Existing home CI | ready (via `/launch` ship) |
| `linkedin` | social | LinkedIn Community Management / UGC Posts | `LINKEDIN_ACCESS_TOKEN`, `LINKEDIN_ORG_URN` | Company Page → [LinkedIn Developers](https://www.linkedin.com/developers/) app → OAuth (w_member_social / w_organization_social) | blocked until secrets |
| `x_twitter` | social | X API v2 manage tweets | `X_BEARER_TOKEN` or `X_API_KEY`+`X_API_SECRET`+`X_ACCESS_TOKEN`+`X_ACCESS_SECRET` | [developer.x.com](https://developer.x.com/) project + write app | blocked until secrets |
| `bluesky` | social | AT Protocol `com.atproto.repo.createRecord` | `BLUESKY_HANDLE`, `BLUESKY_APP_PASSWORD` | bsky.app → App Passwords | blocked until secrets |
| `mastodon` | social | Mastodon REST `/api/v1/statuses` | `MASTODON_BASE_URL`, `MASTODON_ACCESS_TOKEN` | Instance account → Preferences → Development → New application | blocked until secrets |
| `meta_threads` | social | Threads API | `THREADS_ACCESS_TOKEN`, `THREADS_USER_ID` | Meta Developer app + Threads product + test user | blocked until secrets |
| `facebook_page` | social | Meta Graph `/{page-id}/feed` | `FACEBOOK_PAGE_ID`, `FACEBOOK_PAGE_TOKEN` | Meta Business + Page + app (often same as Threads) | blocked until secrets |
| `reddit` | social | Reddit API submit | `REDDIT_CLIENT_ID`, `REDDIT_CLIENT_SECRET`, `REDDIT_USERNAME`, `REDDIT_PASSWORD`, `REDDIT_USER_AGENT` | reddit.com/prefs/apps → script app; respect subreddit rules | blocked until secrets |
| `email` | email | Resend API | `RESEND_API_KEY`, `RESEND_FROM` | [resend.com](https://resend.com) + DNS for dasmlab.org | blocked until secrets |
| `sms` | sms | Twilio Messages | `TWILIO_ACCOUNT_SID`, `TWILIO_AUTH_TOKEN`, `TWILIO_FROM_NUMBER` | Twilio account + **opt-in recipient list only** (no cold SMS) | blocked until secrets |

## Wave-2 (catalog only)

YouTube Community, Instagram Graph, WhatsApp Business Cloud, Discord webhook, Slack webhook, Google Business Profile posts.

## Character limits (dry-run)

| Channel | Soft limit |
|---------|------------|
| linkedin | 3000 |
| x_twitter | 280 |
| bluesky | 300 |
| mastodon | 500 |
| meta_threads | 500 |
| facebook_page | 63206 (preview truncates ~500) |
| reddit | title 300 / body 40000 |
| email | subject 78 / body unlimited (preview) |
| sms | 160 (GSM) / segments noted |
| web_slash | SEO title ~60 / description ~155 |

## Setup order (recommended)

1. Bluesky + Mastodon (fastest tokens)  
2. Resend DNS for dasmlab.org  
3. LinkedIn Company Page app  
4. X developer app  
5. Meta (Threads + Facebook Page)  
6. Reddit (if community post planned)  
7. Twilio + opt-in list  

Wire secrets into `dpo-secrets` the same way as `GSC_CREDENTIALS_JSON`.
