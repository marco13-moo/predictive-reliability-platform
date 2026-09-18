# ADR 0005: Anomaly Detection Baseline-vs-ML Evaluation

## Status
Accepted. The deterministic baseline is the release gate and ML is an optional, comparable adapter.

## Context
The platform needs early detection before a model-training pipeline and needs to demonstrate that a model improves outcomes rather than merely producing more alerts. Reliability signals are noisy, non-stationary, and often sparsely labeled. Accuracy is misleading when incidents are rare, and opaque scores are difficult to explain during an incident. A reproducible baseline gives engineers a safety reference and a rollback comparator.

## Decision
Implement a deterministic, explainable baseline (including a configured z-score/robust deviation method, minimum sample count, warm-up behavior, and threshold) behind the detector interface. Candidate ML detectors must consume the same canonical events, emit the same finding contract (score, threshold, evidence window, detector/version, and confidence), and be evaluated offline against versioned labeled replay fixtures. Compare precision, recall, F1, false-positive rate per service/window, detection latency, alert volume, and resource cost; report confidence intervals or sample counts where meaningful. No model may bypass validation, tenant scoping, audit, or remediation gates.

Thresholds and feature windows are configuration with provenance. Training and evaluation data are separated by time or incident to avoid leakage. A model can be deployed only when it meets baseline safety thresholds and has an explicit rollback path; “improvement” must be demonstrated on representative holdout incidents, not just aggregate accuracy.

## Alternatives Considered
- **ML-only detection:** rejected because it lacks a transparent fallback and is difficult to validate with sparse labels.
- **Accuracy as the sole metric:** rejected because class imbalance hides missed incidents and alert fatigue.
- **Per-detector incomparable outputs:** rejected because operators need one finding and audit contract.
- **Unbounded adaptive thresholds:** rejected because drift can silently change behavior without review.

## Consequences
The baseline remains understandable and operationally dependable while ML experimentation is possible. Evaluation requires curated labels, replay infrastructure, and explicit drift monitoring. False positives may increase during warm-up or topology changes; those states must be visible. Model artifacts, feature definitions, and evaluation datasets become governed dependencies.

## Explicit Non-goals
This ADR does not select an ML framework, claim stationarity, promise zero false positives, or authorize autonomous remediation. It does not define a universal incident-labeling process or replace human review of model changes.

## Verification/Exit Criteria
- Baseline tests cover warm-up, constant/zero variance, missing data, outliers, and deterministic threshold boundaries.
- Every detector output identifies algorithm/model version, threshold, feature/evidence window, and input provenance.
- A holdout replay report includes confusion counts, precision, recall, F1, latency, false-positive rate, and alert volume for baseline and candidate.
- Release criteria and rollback thresholds are recorded with the model/configuration artifact.
- A tenant-isolation and malformed-input test proves detectors cannot leak data or panic.
