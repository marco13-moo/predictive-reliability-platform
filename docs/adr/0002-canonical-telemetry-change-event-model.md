# ADR 0002: Canonical Telemetry/Change-Event Model

**Status**
Accepted

**Context**
Signals differ in identity, time, schema, tenancy, and provenance.

**Decision**
Canonical Events carry identity, schema version, tenant, UTC time, labels, and provenance and are validated at ingestion.

**Alternatives**
Vendor structs or implicit tenant context.

**Consequences**
Downstream processing is portable and deterministic.

**Explicit Non-goals**
Not every vendor field or semantic source validation.
