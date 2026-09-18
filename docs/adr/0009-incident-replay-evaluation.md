# ADR 0009: Incident Replay and Evaluation

## Status
Accepted. Replay is the regression and comparison mechanism for detector, correlation, SLO, and safety changes.

## Context
Online outcomes are affected by traffic, clocks, deployment mix, sampling, and changing topology; they cannot alone prove that a change improved reliability. Incident labels and event timelines are valuable but imperfect, and accuracy alone hides missed incidents and alert fatigue. Engineers need a repeatable way to compare versions without executing production side effects.

## Decision
Replay consumes a versioned, tenant-scoped fixture of canonical events, labels/annotations, topology-as-of data, configuration, and expected policy mode. It runs the same pure detector, SLO, correlation, and remediation-policy interfaces with injected time and deterministic ordering. Reports include confusion counts, precision, recall, F1, false-positive rate, alert volume, detection latency, coverage, policy decisions, and safety violations; denominators and unknown/unlabeled intervals are explicit. Reports identify code/config/schema versions and fixture provenance.

Replay never calls external executors, mutates production state, or treats missing labels as negatives. Dataset splits are by incident/time where model tuning is involved. Golden fixtures cover ordinary, boundary, late, duplicate, and malformed events. Differences are reviewable at event/finding level, not only as aggregate metrics.

## Alternatives Considered
- **Online-only validation:** rejected because it is nondeterministic and unsafe for policy changes.
- **Accuracy alone:** rejected because rare incidents make it misleading.
- **Synthetic ground truth only:** rejected because synthetic data misses operational messiness.
- **Replay that invokes real remediation:** rejected because evaluation must be side-effect free.

## Consequences
Changes become regression-testable and trade-offs are visible before rollout. Curating labels and keeping fixtures representative require ongoing ownership; metrics can be statistically weak for rare incidents. Replay must preserve compatibility with evolving contracts and clearly report unsupported historical data.

## Explicit Non-goals
Replay does not generate ground truth, prove causality, guarantee production performance, or replace live monitoring and canary rollout. It does not authorize remediation or disclose unrestricted tenant data.

## Verification/Exit Criteria
- A fixed fixture produces a stable report and event-level diff across repeated runs.
- Reports include all stated metrics, denominators, versions, and unknown-label handling.
- Tests prove no network, executor, credential, or production-store mutation occurs during replay.
- Regression thresholds fail clearly on missed incidents, unsafe recommendations, or material false-positive increases.
- At least one representative fixture covers each supported event kind and late/duplicate delivery behavior.
