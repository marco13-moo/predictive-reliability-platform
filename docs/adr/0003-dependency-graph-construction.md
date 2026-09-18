# ADR 0003: Dependency Graph Construction

## Status
Accepted. Graph construction is a replaceable enrichment capability and may produce partial graphs.

## Context
Detection and investigation need context about services, workloads, databases, queues, and ownership relationships. Topology may come from catalogs, deployment manifests, service discovery, or manually maintained inventories, each with different freshness and confidence. A single authoritative catalog is unrealistic, while hard-coded topology would make the platform brittle. Graph data also changes over time, so using current topology to explain a historical incident can be misleading.

## Decision
Provide a storage-neutral graph-builder interface that consumes normalized relationship events with source timestamps, effective intervals where known, tenant scope, provenance, and confidence. Nodes have stable typed identifiers; edges have relation type, direction, source, observed/effective times, and optional confidence. Builders must support deterministic construction for a requested tenant and time/as-of point, tolerate missing endpoints, and expose freshness and completeness metadata. Conflicting edges are retained with provenance and resolved only by an explicit, configurable precedence policy.

Graph enrichment is advisory input to correlation and display, never proof of causality. Rebuilds from the same event set must be equivalent. Consumers must be able to distinguish “no relationship,” “relationship unknown,” and “relationship observed but stale.”

## Alternatives Considered
- **Hard-coded graph in detector code:** rejected because topology changes faster than detector logic.
- **One central catalog as source of truth:** rejected because deployments often span systems and catalogs have different coverage and latency.
- **Live network scanning:** rejected for safety, permissions, repeatability, and production impact.
- **Overwrite conflicting edges:** rejected because it destroys evidence and hides source disagreement.

## Consequences
The platform can add topology sources without changing detectors and can replay historical graph state. Partial graphs preserve forward progress but may reduce confidence; freshness and completeness must be visible in findings. Storage-neutral interfaces require adapter work and explicit indexing in production. Graph size, cardinality, and tenant isolation become operational constraints.

## Explicit Non-goals
This ADR does not perform network discovery, infer undocumented dependencies, establish causal relationships, or select a graph database. It does not define ownership escalation, universal service identity, or a guarantee that every edge is current.

## Verification/Exit Criteria
- Builder tests cover out-of-order, duplicate, conflicting, stale, and missing-endpoint relationship events.
- Rebuilding the same tenant/as-of graph from the same fixture produces byte-for-byte stable serialized output or a documented canonical ordering.
- Every node and edge retains tenant, source, observed time, and confidence/freshness metadata.
- Consumers receive explicit completeness/freshness status and do not treat absent edges as proof of independence.
- A tenant-scoped test proves no node or edge from another tenant is observable through the builder interface.
