# Spec 001: API Versioning Baseline + Unversioned Deprecation

## Goal
Introduce versioned routing so all new APIs live under `/v1`, while keeping a clearly deprecated unversioned root route for transition.

## Scope
- Mount Gin API routes under a `/v1` prefix during app startup.
- Add a versioned health route: `GET /v1/health`.
- Keep unversioned `GET /` temporarily, but mark it deprecated in behavior (message + deprecation header).
- Update existing tests that currently expect unversioned baseline behavior.

## Out of scope
- Business resources (products/orders/categories/promotions).
- `swaggo/swag` API docs generation.

## Acceptance criteria
- `GET /v1/health` returns `200` with JSON payload (`status: "ok"`).
- `GET /` still returns `200`, but includes deprecation header and migration message.
- No routes outside `/v1` are introduced except deprecated root fallback.
- E2E/integration test suite updated to match new behavior.

## Verification
- `go test ./internal/http/routes/... ./test/integration/...`
- `curl -i http://localhost:8080/v1/health`
- `curl -i http://localhost:8080/`
