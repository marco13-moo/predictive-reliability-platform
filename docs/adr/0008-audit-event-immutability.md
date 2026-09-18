# ADR 0008: Audit Event Immutability

## Status
Accepted. Audit records are append-only evidence; any correction is a new linked record.

## Context
Reliability findings, correlation hypotheses, SLO decisions, and remediation gates must be explainable after an incident. Mutable history permits accidental or malicious rewriting and prevents reliable replay. Audit storage may be implemented by different adapters, so integrity must be established at the record contract rather than assumed from a database. Serialization failures must not silently create unverifiable records.

## Decision
Represent each audit event with immutable fields, event ID, tenant/scope, event type, actor/source, UTC timestamp, payload/version, and predecessor hash. Compute a SHA-256 hash over a canonical, versioned serialization of the event plus predecessor hash; append-only stores reject mutation or invalid predecessor links. APIs return defensive copies and surface marshal/hash/persistence errors. Corrections, redactions, and supersessions append new events that reference the original, preserving the chain and reason. Verification can walk a chain and report the first broken link.

Canonical serialization excludes nondeterministic map ordering and documents encoding/version. Hash chaining provides tamper evidence, not confidentiality; access and retention remain adapter responsibilities. Audit emission is part of decision processing and failures are observable rather than swallowed.

## Alternatives Considered
- **Mutable audit rows:** rejected because history can be rewritten and replay loses provenance.
- **Best-effort serialization:** rejected because a missing hash undermines the evidence contract.
- **Single final digest only:** rejected because localized verification and append validation are needed.
- **Rely solely on database immutability:** rejected because adapters differ and application-level evidence must travel with records.

## Consequences
Changes are detectable and decision history is reproducible across storage adapters. Chains require careful canonical encoding, key/chain management, and handling of legitimate privacy redactions. Hashes do not prove who wrote a record or prevent deletion; operational controls and access logs remain necessary.

## Explicit Non-goals
No public notarization, encryption, key management, authorization system, legal evidentiary guarantee, or immutable external ledger is selected.

## Verification/Exit Criteria
- Tests prove any payload, predecessor, ordering, or canonicalization change breaks verification.
- Repeated serialization of equivalent events produces identical bytes and hashes across runs.
- APIs cannot mutate stored event values through returned references; marshal and persistence errors are returned.
- Correction/redaction tests append linked records without rewriting originals.
- A failure-injection test demonstrates audit write failure is visible to the caller and operational telemetry.
