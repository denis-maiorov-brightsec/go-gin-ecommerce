# Spec 002: Global Validation Middleware + Consistent Error Envelope

## Goal
Standardize request validation and error responses before ecommerce resource APIs expand.

## Scope
- Register Gin app-level validation/binding error handling.
- Add centralized error handling that normalizes all errors to one response shape.
- Ensure Gin binding or validator errors are mapped into structured `details`.

## Out of scope
- Request ID injection (added later in spec 011).
- Request logging middleware.

## Response contract
Error response format:

```json
{
  "timestamp": "2026-01-01T00:00:00.000Z",
  "path": "/v1/products",
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Request validation failed",
    "details": [
      { "field": "name", "constraints": ["name must not be empty"] }
    ]
  }
}
```

## Acceptance criteria
- Invalid payloads return `400` with the envelope above (same top-level shape for all errors).
- Missing routes return `404` with same envelope format.
- Unknown runtime errors return `500` with sanitized message.
- Existing happy-path routes still work unchanged.

## Verification
- `go test ./internal/http/... ./internal/platform/middleware/... ./test/integration/...`
- `go vet ./...`
- Manual/API checks should cover validation failure, not-found, and generic internal error mapping paths.
