# ADR 0006: Change Correlation as Hypothesis Evidence

## Status
Accepted. Correlation produces ranked evidence, not a causal verdict.

## Context
A deployment, configuration update, feature-flag change, or infrastructure event near a reliability anomaly is useful investigative evidence. Temporal proximity alone cannot distinguish cause from coincidence, delayed impact, common upstream failure, or an unrelated simultaneous change. Teams need useful ranking without false certainty or blame. Correlation must also handle clock skew, late events, service aliases, and multiple changes in one window.

## Decision
Correlate findings with normalized change events using explicit, configurable before/after windows, normalized service/resource identity, tenant scope, and event-time semantics. Emit a hypothesis record containing candidate changes, time offsets, identity match, window/configuration version, evidence links, and a confidence/ranking score whose factors are inspectable. Multiple candidates remain visible; absence of a candidate is “no matching evidence,” not evidence of no cause. Late or corrected events may revise a hypothesis while preserving prior audit records.

The correlator must be deterministic for identical inputs and must not label a change “root cause.” Clock quality, ingestion delay, and topology uncertainty are included as evidence limitations. A downstream investigator or policy may choose how to act on the hypothesis.

## Alternatives Considered
- **Treat nearest change as cause:** rejected because proximity is not causality and creates unsafe remediation pressure.
- **Unbounded time windows:** rejected because they produce noisy, non-actionable candidate sets.
- **Discard all but the top candidate:** rejected because ties and competing changes are operationally important.
- **Infer team blame from actor metadata:** rejected because identity is evidence, not accountability.

## Consequences
Investigations start with ranked, explainable candidates and can be replayed. Confidence remains intentionally bounded, so operators must perform causal validation. Window configuration and identity mapping become important operational inputs; poor clocks or aliases reduce usefulness. Correlation output must be phrased carefully in UI and incident reports.

## Explicit Non-goals
This ADR does not perform causal inference, root-cause analysis, blame assignment, change approval, or automatic rollback. It does not establish that a correlated service is the failing component.

## Verification/Exit Criteria
- Tests cover before/after boundaries, clock skew/late arrival, multiple candidates, duplicate changes, service aliases, and tenant isolation.
- Every hypothesis includes exact event IDs, timestamps, window configuration, and score factors.
- Identical event sets produce stable candidate ordering and scores.
- A review of representative incidents demonstrates that wording distinguishes evidence from causation.
- Replay preserves prior hypotheses and records revisions rather than overwriting history.
