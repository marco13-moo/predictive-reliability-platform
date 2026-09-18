# ADR 0004: Slo/Error-Budget Semantics

**Status**
Accepted

**Context**
Operators need consistent burn and exhaustion decisions.

**Decision**
Compute target-based bad-event allowance, proportional consumption, remaining budget, burn rate, and exhaustion.

**Alternatives**
Unweighted averages or hidden rounding.

**Consequences**
Results are explainable and replayable.

**Explicit Non-goals**
Does not choose targets or page thresholds.
