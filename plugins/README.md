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
| linkedin | DPO GEO | planned |
| chatgpt / perplexity / gemini | DPO/DAO GEO | planned (official APIs only) |

In-tree collectors under `internal/collectors/` until packages move to `plugins/<name>`.
