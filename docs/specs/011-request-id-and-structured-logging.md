# Spec 011: Request ID + Structured Request Logging Middleware

## Goal
Improve observability with per-request correlation IDs and structured logs.

## Scope
- Inject request ID (accept incoming header or generate one).
- Include request ID in response headers.
- Emit structured request logs with status, latency, method/path, and request ID.

## Out of scope
- Full distributed tracing integration.
- Log shipping pipeline setup.

## Acceptance criteria
- Every request has a request ID available in handlers and logs.
- Response includes request ID header.
- Logs use structured format aligned with stack conventions.
- Tests verify request ID propagation behavior.

## Verification
- `go test ./internal/http/... ./internal/platform/... ./test/integration/...`
- `curl -i http://localhost:8080/v1/health`
- Confirm the response includes the request ID header and logs include the same request ID.
