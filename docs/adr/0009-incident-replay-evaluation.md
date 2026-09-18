# ADR 0009: Incident Replay/Evaluation

**Status**
Accepted

**Context**
Detector changes require objective labeled comparisons.

**Decision**
Replay reports confusion counts, precision, recall, F1, detection latency, false-positive rate, and remediation safety.

**Alternatives**
Online-only validation or accuracy alone.

**Consequences**
Changes are regression-testable.

**Explicit Non-goals**
No ground-truth generation.
