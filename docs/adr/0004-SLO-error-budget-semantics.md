# ADR 0004: SLO and Error-Budget Semantics

## Status
Accepted. These semantics are the normative calculation contract for SLO reports and replay.

## Context
Operators need consistent answers to “how much budget remains?” across availability, latency, and other indicators. Unweighted averages and hidden rounding can disagree with dashboards and make policy decisions irreproducible. Measurements may be counts or durations and can include invalid, missing, or partial windows. The SLO package must provide explainable arithmetic without choosing organizational targets or paging policy.

## Decision
For a reporting window, accept an explicit target in [0,1], total eligible events, and bad events (or an equivalent typed measurement). Compute allowed bad events as `total × (1-target)`, consumed budget as `bad/allowed` when allowed is nonzero, remaining budget as `max(0, allowed-bad)`, and burn rate against the configured window allowance. Preserve decimal precision internally and round only at presentation boundaries. Reject negative counts, bad > total, invalid targets, zero/ambiguous denominators, and incompatible units; report empty windows explicitly rather than calling them healthy.

A report includes window start/end, target, totals, bad count, allowance, remaining amount/ratio, burn rate, and an exhaustion indicator with the exact rule used. All calculations are deterministic and use supplied timestamps/configuration. Multi-window aggregation must be weighted by eligible event count or duration as declared by the indicator, never by averaging percentages without weights.

## Alternatives Considered
- **Unweighted average of per-service percentages:** rejected because low-volume services distort fleet results.
- **Implicit rounding before calculation:** rejected because boundary decisions change with display precision.
- **Treat missing data as success:** rejected because it masks collection failures; missingness is reported separately.
- **Let the SLO package choose page thresholds:** rejected because thresholds are operational policy and vary by service.

## Consequences
Reports are comparable, replayable, and auditable; callers can explain every number. Callers must provide correct eligibility semantics and choose indicator-specific aggregation. Zero-volume and partial windows require explicit UI treatment. Existing dashboards may change where they previously averaged or rounded incorrectly, so migration comparisons are required.

## Explicit Non-goals
This ADR does not set SLO targets, alert/page thresholds, maintenance-window policy, service ownership, or error classification. It does not define a storage schema or guarantee statistical confidence intervals.

## Verification/Exit Criteria
- Table-driven tests cover perfect, exhausted, over-budget, empty, partial, invalid, and boundary target cases.
- A hand-calculated reference fixture matches implementation values before presentation rounding.
- Aggregation tests prove weighting by declared eligible count/duration and reject mixed units.
- Reports include enough inputs and formula version to reproduce every output.
- Replay of identical SLO events produces identical reports across runs and machines.
