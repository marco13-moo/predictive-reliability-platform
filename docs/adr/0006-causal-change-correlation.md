# ADR 0006: Causal/Change Correlation

**Status**
Accepted

**Context**
Temporal proximity helps but does not prove causality.

**Decision**
Correlate same-service changes inside explicit windows and retain confidence as hypothesis evidence.

**Alternatives**
Treat correlation as causation or use unbounded windows.

**Consequences**
Uncertainty remains transparent.

**Explicit Non-goals**
No root-cause or team blame inference.
