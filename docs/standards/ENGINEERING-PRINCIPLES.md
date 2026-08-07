# DASMLAB Engineering Principles

**Status:** Normative for all Observatory repositories  
**Related:** [ADR-0001](../adr/ADR-0001-observatory-platform-vision.md), [ADR-9999](../adr/ADR-9999-innovation-principles.md)

Copy this section into every repository README.

## The ten principles

1. **Everything is observable.** If it exists in production or research, it emits measurable signals.
2. **Everything exposes an API.** UIs are clients; agents and automation are first-class consumers.
3. **Everything emits events.** State changes become normalized events on the platform bus.
4. **Everything is measurable.** Prefer quantified scores and trends over anecdotes.
5. **Everything is pluggable.** Integrations implement SDK contracts; products do not hard-code vendors.
6. **Everything is automatable.** Observations can become recommendations, then approved actions.
7. **Everything is documented.** Docs are a product output (SEO/GEO/knowledge), not an afterthought.
8. **Everything is testable.** Collectors, scores, and recommendations have fixtures and golden tests.
9. **Everything is explainable.** Recommendations carry evidence, confidence, impact, and effort.
10. **Everything contributes knowledge.** Repos feed the knowledge graph, research loop, and Labs.

## Five questions (every Observatory UI/API)

1. What exists?
2. What changed?
3. Why did it change?
4. What will happen?
5. What should I do?

## Maturity model

| Level | Name | Capability |
|-------|------|------------|
| L1 | Observe | Collect metrics and events |
| L2 | Understand | Correlate across sources |
| L3 | Explain | Identify root causes (WHY) |
| L4 | Predict | Forecast outcomes |
| L5 | Recommend | Actionable, explainable advice |
| L6 | Automate | Execute approved actions |
| L7 | Learn | Improve from recommendation outcomes |

New features should advance (or solidify) a maturity level without regressing explainability.
