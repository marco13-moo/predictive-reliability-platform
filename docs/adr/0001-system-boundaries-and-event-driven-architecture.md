# ADR 0001: System Boundaries And Event-Driven Architecture

**Status**
Accepted

**Context**
Heterogeneous signals need replaceable integrations and auditable processing.

**Decision**
Use a contract-first event-driven pipeline with explicit normalize, analyze, correlate, remediate, and audit boundaries.

**Alternatives**
A vendor SDK monolith or hard broker dependency.

**Consequences**
Components stay testable and extensible.

**Explicit Non-goals**
Not a broker or workflow orchestrator.
