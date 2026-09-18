# ADR 0010: Storage, Retention, and Privacy

## Status
Accepted. Storage and lifecycle controls are deployment-specific adapters constrained by this data-minimization contract.

## Context
Telemetry and audit evidence can contain tenant identifiers, resource names, actor metadata, traces, and accidental secrets. Retention obligations, deletion requests, residency, access controls, and incident-investigation needs vary by deployment. An indefinite raw-data store would increase breach impact and make deletion impossible to reason about. Core packages must not assume a particular database or silently retain sensitive payloads.

## Decision
Keep storage behind explicit interfaces with tenant-scoped reads/writes, retention metadata, deletion/purge operations, and authorization context supplied by the deployment. Minimize data at normalization: collect only fields needed for detection, correlation, replay, and audit; redact or hash configured sensitive values; bound payload and label sizes; and never store credentials. Raw source extensions are opt-in, access-controlled, and subject to the shortest applicable retention. Configure separate retention classes for canonical events, derived findings, audit chains, and replay fixtures, with documented legal-hold/exception handling and deletion outcomes. Encryption in transit/at rest, residency, backups, and access logging are adapter controls but must be testable at integration boundaries.

Tenant identity is mandatory on stored records and every query is scoped. Deletion or redaction creates an auditable lifecycle record without pretending that cryptographic hashes reveal the deleted payload. Metrics and reports use aggregation or pseudonymization where individual identity is unnecessary.

## Alternatives Considered
- **Indefinite retention of all raw payloads:** rejected because risk and compliance scope grow without bound.
- **One centralized unrestricted store:** rejected because it violates least privilege and deployment-specific residency/lifecycle needs.
- **Drop all provenance immediately:** rejected because investigations and parser verification need bounded provenance.
- **Rely on application convention for tenant filtering:** rejected because storage interfaces must enforce scope structurally.

## Consequences
Deployments can meet different lifecycle and residency needs while core logic remains portable. Operators must configure and periodically verify retention jobs, deletion propagation, backup handling, and access controls. Minimization may limit later investigations; exceptions require documented approval rather than silent expansion. Hash-linked audit evidence needs careful treatment when referenced payloads are purged.

## Explicit Non-goals
This ADR does not define legal advice, a universal retention duration, a centralized data lake, a specific database/cloud, or a complete identity and access-management product. It does not guarantee deletion from systems outside the platform’s control.

## Verification/Exit Criteria
- Integration tests prove every storage operation requires tenant scope and unauthorized cross-tenant reads fail.
- Retention tests delete each data class at its configured deadline and report failures; backup and legal-hold behavior is documented.
- Secret/redaction tests prove configured sensitive fields are not persisted in canonical or audit payloads.
- Records expose provenance, retention class, and lifecycle status sufficient for an operator to explain deletion.
- A periodic operational check verifies encryption, access logging, purge lag, and sampled deletion across adapters.
