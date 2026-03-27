# Stack Profile

This repository implements an ecommerce backoffice API in Go.
Use this file for stack/tooling decisions only; API behavior/contracts are defined by specs.

## Product Context
- Project name: `go-gin-ecommerce`
- Domain: `ecommerce-backoffice`
- API style: `REST`

## Core Tech Choices
- Language: `Go`
- Runtime: `Go 1.23+`
- Framework: `Gin`
- ORM / Data Mapper: `GORM`
- Database: `PostgreSQL`

## Repository Conventions
- Package/dependency manager: `Go modules`
- Migration strategy: `SQL-first migrations stored in db/migrations/, applied in order with golang-migrate; schema changes must be committed as forward-only up/down migration pairs`
- Configuration style: `environment variables with optional .env.local for local development; load into a typed config package at startup and fail fast on missing required settings`

## Repository Topology Contract
- Source root path: `cmd/`, `internal/`, `pkg/`
- Module path pattern: `internal/<module>/{http,service,repository,model,dto}`
- Shared/common code path: `internal/platform/` for app wiring and infrastructure, `internal/common/` for reusable business helpers, `pkg/` only for genuinely reusable public utilities
- DB/migrations path: `db/migrations/`
- Test path strategy: `co-located` + `internal/**/**/*_test.go`, plus `separate` integration tests in `test/integration/**/*.go`
- API docs artifact path (if generated): `docs/openapi/openapi.yaml`
- Prohibited top-level paths: `src/`, `app/`, `lib/`, `misc/`, `controllers/`, `services/`, `repositories/`

Concrete repository layout to follow:
- App entrypoint: `cmd/api/main.go`
- Bootstrap and infrastructure wiring: `internal/platform/{config,db,httpserver,logger,middleware}/`
- Feature modules: `internal/<module>/{http,service,repository,model,dto}/`
- Route registration: `internal/http/routes/`
- Shared API response helpers: `internal/common/api/`
- Database migrations: `db/migrations/`
- Seeds or local fixtures: `db/seeds/`
- Integration test helpers: `test/integration/testutil/`

## Quality Gates
- Lint command: `go vet ./...`
- Unit test command: `go test ./...`
- Integration/e2e test command: `go test ./test/integration/...`
- Type-check/static-analysis command: `go test ./...`

## Implementation Preferences (Optional)
- Validation library preference: `Gin binding + go-playground/validator, wrapped so error responses stay spec-driven and consistent`
- Logging library preference: `slog`
- API docs tool preference: `swaggo/swag`, with generated artifacts committed to `docs/openapi/` only when a spec requires API docs; do not maintain handwritten OpenAPI files`
- Auth library preference: `custom Gin middleware with bearer-token parsing; keep auth provider integration behind internal/platform/auth/`

## Additional Constraints
- Performance/security/compliance requirements: `use request-scoped contexts for DB calls, connection pooling via pgx-backed PostgreSQL driver, structured logging, no in-memory runtime repositories for feature behavior, and avoid N+1 query patterns on list endpoints`
- Deployment/runtime environment: `container-friendly Linux service, stateless API process, PostgreSQL as the only required stateful dependency, local development via Docker Compose when added`
- Backward-compatibility rules: `preserve existing routes and response contracts unless the active spec explicitly changes them; version API changes under /v1 as directed by specs; deprecations must be explicit in docs or code comments when required by the spec`

## Precedence Rules
- Specs are the source of truth for API behavior and contracts.
- If this profile conflicts with a spec, follow the spec.
