# ADR 0008: Audit/Event Immutability

**Status**
Accepted

**Context**
Decisions need tamper evidence and reproducibility.

**Decision**
Use append-only SHA-256 predecessor chains, explicit marshal errors, and defensive copies.

**Alternatives**
Mutable history or best-effort serialization.

**Consequences**
Tampering is detectable and persistence is pluggable.

**Explicit Non-goals**
No notarization, encryption, or authorization.
