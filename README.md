# Predictive Reliability Platform

An executable, research-grade foundation for predicting reliability risk and
automating incident response without allowing analysis to mutate production
systems. The first slice is deliberately deterministic: it provides a
reproducible baseline against which future causal and machine-learning
detectors can be measured.

## Problem and personas

Modern services emit metrics, logs, traces, deployments, feature-flag changes,
and cloud events through incompatible interfaces. Operators need a trustworthy
answer to three questions: *what is abnormal, what changed, and what action is
safe?* Reliability engineers need measurable detection quality; service owners
need explainable hypotheses and error-budget context; incident commanders need
bounded, approval-aware actions; researchers need replayable data and baselines.

## Goals and non-goals

The platform will normalize observations and changes into versioned,
tenant-scoped contracts; calculate SLO/error-budget state; detect anomalies;
correlate them with recent changes; rank evidence-backed incident hypotheses;
evaluate detectors on labeled replay data; and produce remediation plans that
fail closed for ambiguous or high-risk actions. Every decision is auditable.

This repository does **not** deploy changes, replace an incident-management
system, claim causal certainty from correlation, provide a hosted control
plane, or pretend that synthetic fixtures represent production behavior.
Vendor adapters (OTLP, Prometheus, logs, traces, cloud events, and feature
flags) are extension points, not fake integrations.

## Architecture

The analysis plane is separated from action execution:

```text
observations/changes -> normalize + validate -> canonical events
                                      |             |
                                      v             v
                              SLO + detector -> correlation -> hypotheses
                                                                  |
                                                                  v
                                                        policy + approval gate
                                                                  |
                                      audit hash chain <--- bounded plan
```

`pkg/contracts` defines identity, schema version, tenant scope, timestamp, and
provenance. `pkg/ingest` normalizes and orders events with duplicate
suppression. `pkg/slo`, `pkg/detector`, `pkg/correlate`, and `pkg/incident`
implement the analysis spine. `pkg/remediation` only decides whether a plan
may proceed; it has no production executor. `pkg/audit` records tamper-evident
decisions, while `pkg/replay` supplies evaluation metrics.

## Threat and safety model

Inputs are untrusted and must validate before analysis. Tenant and service
scope are explicit to prevent cross-scope joins. Event IDs support
idempotency; timestamps are normalized to UTC; provenance is retained for
forensics. The platform is fail-closed: high-risk actions require explicit
approval, ambiguous policy inputs are rejected, and analysis cannot directly
mutate production. Audit records are append-only and hash chained. Production
deployments should additionally enforce least-privilege credentials,
encryption, access logging, retention controls, and independent approval
identity.

## Research questions

1. When do simple robust baselines outperform ML detectors under drift and
   sparse labels?
2. Which temporal, dependency, and change features improve ranking without
   overstating causal confidence?
3. How should detector thresholds trade recall, latency, false positives, and
   operator toil across services?
4. Which approval policies reduce unsafe automation while preserving useful
   time-to-mitigation?
5. How can replay datasets remain representative without exposing sensitive
   telemetry?

## Evaluation criteria

Every detector/replay experiment reports precision, recall, F1, false-positive
rate, detection latency, and coverage of labeled incidents. Remediation
evaluation reports the fraction of executed actions that were explicitly
approved and policy-compliant. Experiments must use fixed fixtures or versioned
datasets, document labeling and matching windows, compare against the
deterministic baseline, and retain enough metadata to reproduce results.

## Local quickstart

Requirements: Go 1.22 or newer.

```sh
gofmt -w .
go test ./...
go vet ./...
go run ./cmd/demo
```

The demo uses deterministic fixtures and no cloud credentials. It prints the
anomaly count, incident severity, remaining budget, approval decision, audit
verification, and replay F1. Synthetic fixtures exercise code paths only; they
are not evidence of production performance.

## Roadmap

1. **Foundation (current):** contracts, deterministic ingestion/analysis,
   replay metrics, safety gates, audit trail, fixtures, and ADRs.
2. **Telemetry adapters:** OTLP, Prometheus, logs, traces, cloud events, and
   feature-flag providers with contract conformance tests.
3. **Operational persistence:** durable event storage, retention enforcement,
   privacy controls, and multi-tenant authorization.
4. **Research plane:** dependency graphs, drift-aware robust baselines, causal
   inference experiments, and gated ML model comparison.
5. **Execution plane:** independently deployed, least-privilege executors with
   dry-run, rollback, approval identity, and post-action verification.

See [docs/architecture.md](docs/architecture.md), the
[ADRs](docs/adr/), [CONTRIBUTING.md](CONTRIBUTING.md), and
[SECURITY.md](SECURITY.md).
