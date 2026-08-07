# Edge bot tokens (2026)

First-party crawl / AI visibility from HAProxy access logs (UA classification).

| Token | Metric | Notes |
|-------|--------|-------|
| Googlebot | googlebot_fetches | Classic search crawl |
| GoogleOther | googleother_fetches | Google non-search fetches |
| bingbot | bingbot_fetches | Bing |
| GPTBot | gptbot_fetches | OpenAI training/crawl |
| OAI-SearchBot | oai_searchbot_fetches | OpenAI search |
| ChatGPT-User | chatgpt_user_fetches | User-initiated |
| ClaudeBot | claudebot_fetches | Anthropic |
| Claude-SearchBot | claude_searchbot_fetches | Anthropic search |
| PerplexityBot | perplexitybot_fetches | Perplexity |
| Perplexity-User | perplexity_user_fetches | Perplexity user |

Pull sample from the edge host (self-hosted runner SSH):

```bash
bash scripts/ci/pull-haproxy-access-sample.sh ./.data/edge/access.log
# Mount/copy into DPO PVC at /data/edge/access.log (DPO_EDGE_LOG_PATH)
```

IP CIDR verification of bot authenticity is a Phase-1.5 enhancement.
