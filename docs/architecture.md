# Architecture and maintainer process

Contracts are stable at the package boundary. The pipeline is normalize → calculate SLO → detect → correlate → rank → policy → audit. Components are pure where possible, deterministic for replay, and replaceable through small interfaces. Maintainers require tests, review, CI green, and an ADR for externally visible contract changes. Releases are tagged from the default branch after changelog review.
