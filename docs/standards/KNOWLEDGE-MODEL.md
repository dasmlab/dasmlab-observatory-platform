# Knowledge model — think like a university

University and research sites perform because they publish **interconnected knowledge**, not isolated pages. DOP and `dasmlab_home` follow the same pattern.

## Interconnected hierarchy (example)

```text
AI
  Tutorials · Projects · Videos · Source Code · Architecture · Benchmarks · Labs
Networking
  BGP · VXLAN · Calico · MetalLB
Cloud Native
  Kubernetes · OpenShift · Helm
```

Extend the same shape for Security, Infrastructure, Digital Presence, DevOps, Observability.

## Why this matters for Observatories

| Question | Knowledge-model answer |
|----------|------------------------|
| What exists? | Entities in the hierarchy (hubs, repos, ADRs, labs) |
| What changed? | New/updated nodes and edges |
| Why? | Correlation across docs ↔ code ↔ deploy ↔ crawl |
| What will happen? | Coverage gaps / citation risk |
| What should I do? | Fill a missing hub, republish sitemap, ship a benchmark |

## Graph destination

Phase 2+ Neo4j twin (ADR-0009): Person → Author → Project → Repository → Technology → Article → Issue → Research → Company.

## Related

- [DOCUMENTATION-AS-PRODUCT](./DOCUMENTATION-AS-PRODUCT.md)
- [FLYWHEEL](./FLYWHEEL.md)
- Home hubs inventory (dasmlab_home)
