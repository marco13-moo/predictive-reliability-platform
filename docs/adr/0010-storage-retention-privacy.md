# ADR 0010: Storage/Retention/Privacy

**Status**
Accepted

**Context**
Telemetry can contain tenant-sensitive data and obligations.

**Decision**
Keep storage neutral; adapters configure retention, deletion, access control, and minimization with tenant scope.

**Alternatives**
Indefinite retention or raw secrets.

**Consequences**
Deployments can meet lifecycle obligations.

**Explicit Non-goals**
No legal policy or centralized data store.
