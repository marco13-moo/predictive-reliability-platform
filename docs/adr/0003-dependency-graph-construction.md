# ADR 0003: Dependency Graph Construction

**Status**
Accepted

**Context**
Hypotheses need dependency context while topology sources vary.

**Decision**
Provide a storage-neutral graph-builder extension fed by normalized relationship events and source timestamps.

**Alternatives**
Hard-coded topology or one catalog.

**Consequences**
Graphs are independently rebuildable and incomplete graphs are tolerated.

**Explicit Non-goals**
No network scanning or causal claims.
