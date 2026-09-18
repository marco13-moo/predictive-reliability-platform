# ADR 0001: System Boundaries and Event-Driven Architecture

## Status
Accepted. This is the baseline architecture for the predictive reliability platform; changes require a replacement ADR and compatibility plan.

## Context
The platform combines telemetry, deployment/change records, dependency topology, SLO measurements, detector findings, correlation hypotheses, remediation recommendations, and audit evidence. These inputs arrive at different rates, with different delivery guarantees and schemas. A vendor-specific monolith would couple ingestion to analysis and make replay, testing, and source replacement expensive. The Go packages currently map naturally to contracts, ingest, detector, correlate, remediation, replay, SLO, and audit concerns. Safety-sensitive outputs must be explainable and must not silently acquire side effects.

## Decision
Use a contract-first, event-oriented pipeline with explicit boundaries: **ingest/normalize → analyze (SLO and detection) → correlate → recommend remediation → audit/replay**. Boundary-crossing data uses versioned canonical contracts and immutable values; adapters translate vendor payloads at the edge. Processing components are deterministic where practical, accept injected clocks/configuration, and return decisions or reports rather than performing external effects. A broker, queue, or HTTP transport may be added by an adapter, but core packages must remain usable synchronously and in replay.

Each boundary must document ownership, validation, error behavior, and idempotency expectations. Unknown schema versions are rejected or quarantined rather than guessed. Tenant and provenance context travels with every event. Audit records are emitted for decisions and policy outcomes, not only successful actions.

## Alternatives Considered
- **Vendor SDK monolith:** rejected because it entangles source schemas, credentials, transport retries, and domain logic.
- **Hard dependency on one broker/workflow engine:** rejected because it would make local tests, replay, and small deployments depend on infrastructure.
- **Shared mutable database as the integration contract:** rejected because it obscures ownership and weakens event provenance and version compatibility.
- **Synchronous request chain only:** rejected because bursty telemetry and replay need buffering and independent consumers, even if an adapter chooses synchronous execution.

## Consequences
Positive: packages can be tested with fixtures, sources can be replaced independently, and the same normalized input can drive online processing and historical replay. Explicit boundaries make failure and retry policy reviewable. Negative: contracts, schema evolution, deduplication, and observability require deliberate maintenance; eventual consistency means a finding may be revised when late evidence arrives. Operators must understand that an event accepted by ingest is not proof that downstream consumers have processed it.

## Explicit Non-goals
This ADR does not select a cloud, broker, workflow orchestrator, database, deployment topology, or vendor integration. It does not promise exactly-once delivery, automatic root-cause analysis, or autonomous production changes. It does not replace authentication, authorization, incident-management, or data-governance policy.

## Verification/Exit Criteria
- Every cross-boundary type has a documented version, owner, validation rule, tenant identifier, UTC event time, and provenance field.
- Unit tests can run the pipeline without a broker, network, credentials, or wall-clock dependence.
- A fixture can be normalized, analyzed, correlated, replayed, and audited with deterministic results (excluding explicitly documented nondeterminism).
- Duplicate and out-of-order events have tested behavior, and unsupported schema versions produce an observable rejection/quarantine outcome.
- A failure in one optional consumer does not prevent audit emission or corrupt the canonical event stream.
