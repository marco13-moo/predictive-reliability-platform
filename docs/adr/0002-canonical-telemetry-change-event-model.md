# ADR 0002: Canonical Telemetry and Change-Event Model

## Status
Accepted. The canonical model is the compatibility boundary for all telemetry and change inputs.

## Context
Metrics, logs, traces, alerts, deployment records, feature-flag changes, and infrastructure events disagree about identity, timestamp precision, units, tenancy, labels, and provenance. If downstream packages consume vendor structs, every detector and replay fixture inherits those differences. Missing or ambiguous identity also creates cross-tenant leakage and unreliable correlation. The platform needs a small stable envelope while preserving source evidence for investigation.

## Decision
Normalize each accepted input into a versioned canonical event envelope containing: globally unique event ID (with source ID retained), event kind, schema version, tenant/service/resource identity, UTC occurrence time, ingestion time, typed measurements or change payload, normalized labels, source/provenance metadata, and optional correlation keys. Measurements declare units and semantic meaning; change events declare actor/source, affected resource, and observed before/after or intent where available. Raw vendor fields may be retained in a bounded, access-controlled extension, but consumers must use canonical fields.

Validation occurs at ingestion: required identity and tenant fields, UTC timestamps, supported schema version, finite numeric values, bounded label keys/values, and payload-kind consistency. Normalization is pure and returns structured validation errors; it does not infer tenant context or invent timestamps. Event IDs support deduplication, while occurrence time (not arrival order) drives temporal analysis. Late and corrected events are represented explicitly rather than mutating history.

## Alternatives Considered
- **Pass-through vendor payloads:** rejected because every consumer would implement incompatible semantics.
- **One untyped map for all data:** rejected because it defers validation and makes units and compatibility unverifiable.
- **Implicit tenant from connection/session:** rejected because replay, batching, and multi-tenant workers make implicit context unsafe.
- **Discard source payload after normalization:** rejected because investigations and parser regression tests need provenance, subject to minimization and retention controls.

## Consequences
Consumers gain portable, deterministic inputs and fixtures remain stable across integrations. Schema evolution is explicit and adapters can be upgraded independently. The envelope adds normalization cost and requires a versioning policy; some vendor-specific richness is unavailable to generic consumers. Clock skew, duplicate delivery, and corrections remain operational concerns and must be surfaced in metadata rather than hidden by normalization.

## Explicit Non-goals
The model does not guarantee semantic correctness of a source, solve identity resolution across every vendor, provide a universal ontology, or retain unrestricted raw data. It does not define transport serialization, storage technology, alert thresholds, or authorization policy.

## Verification/Exit Criteria
- Contract tests accept valid metric and change fixtures and reject missing tenant, identity, unit, timestamp, or unsupported-version inputs with field-level errors.
- Two adapters producing the same observation yield equivalent canonical events, including normalized UTC time and units.
- Duplicate IDs have documented idempotent behavior; late and correction events remain distinguishable in replay tests.
- Every emitted event carries provenance sufficient to identify source, parser/version, and ingestion time.
- Fuzz/property tests demonstrate bounded labels/payloads and no panic on malformed vendor input.
