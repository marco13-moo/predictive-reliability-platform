# ADR 0005: Anomaly Detection Baseline-Vs-Ml Evaluation

**Status**
Accepted

**Context**
Explainable detection is needed before models and models need measurable comparison.

**Decision**
Use deterministic z-score baseline and evaluate alternatives with precision, recall, latency, and false positives.

**Alternatives**
ML-only or aggregate accuracy.

**Consequences**
A reproducible comparator exists and ML remains an adapter.

**Explicit Non-goals**
No stationarity or model-governance claim.
