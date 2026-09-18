# Contributing

Use a focused branch and small, tested changes. Run `gofmt -w .`, `go test ./...`, and `go vet ./...` before opening a PR. New behavior should include focused tests and an ADR when it changes a contract or operating model. Keep integrations behind interfaces and avoid credentials or customer data in fixtures.
