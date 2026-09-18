# ADR 0007: Remediation Safety and Approval Gates

## Status
Accepted. Remediation recommendations are policy decisions without side effects; execution is outside this boundary.

## Context
A detector or correlation hypothesis can be wrong, stale, or incomplete. Automated actions such as rollback, traffic shifting, scaling, or disabling a feature can amplify an outage, violate change controls, or affect the wrong tenant. Safety must be enforced before any executor receives a request, and policy outcomes must be explainable. The current remediation package is therefore a pure decision boundary rather than an action runner.

## Decision
Evaluate a typed remediation request containing tenant/resource scope, proposed action, evidence IDs, severity, confidence, blast radius, reversibility, freshness, and requested mode against an explicit policy. Return one of deny, require-approval, or allow-recommendation with reason codes, policy version, expiry, and required approvers. Default-deny applies to missing evidence, stale findings, unknown resources, ambiguous tenant scope, destructive/non-reversible actions, and policy evaluation errors. Even an allowed recommendation remains inert until an external executor performs independent authorization, idempotency, concurrency, and precondition checks.

Policies are versioned and auditable. Approval is explicit, scoped, time-bound, and cannot be inferred from a detector score. The recommendation includes a dry-run or rollback description where available, but never credentials or executable commands in an untrusted payload.

## Alternatives Considered
- **Unreviewed destructive automation:** rejected because detection uncertainty and blast radius are not bounded.
- **Boolean allow/deny only:** rejected because operators need a distinct approval state and actionable reason codes.
- **Policy inside each executor:** rejected because it duplicates safety rules and creates inconsistent gates.
- **Passing credentials/commands through findings:** rejected because it increases privilege and injection risk.

## Consequences
Safety policy is testable, centralized, and visible in audit trails. Incident response may be slower when approval is required, and policy configuration becomes a critical dependency. Executors must implement their own final safeguards; this ADR intentionally does not make a recommendation equivalent to authorization.

## Explicit Non-goals
No action execution, credential storage, change-management replacement, approval identity provider, or guarantee that an approved action is operationally successful is defined here.

## Verification/Exit Criteria
- Tests cover default-deny, stale/expired evidence, ambiguous scope, reversible versus destructive actions, approval expiry, and policy errors.
- Every decision has policy version, evidence references, reason codes, tenant/resource scope, and evaluation time.
- A test proves evaluation performs no network, process, database mutation, or executor call.
- Approved recommendations are bounded by explicit TTL and cannot be replayed as authorization after expiry.
- Security review confirms untrusted evidence cannot inject commands or expand scope.
