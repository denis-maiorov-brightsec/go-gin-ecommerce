GO ?= go

.PHONY: fmt lint test test-integration vet run swagger migrate-up migrate-down

fmt:
	gofmt -w ./cmd ./internal ./test

lint:
	$(GO) vet ./...

test:
	$(GO) test ./...

test-integration:
	$(GO) test ./test/integration/...

vet:
	$(GO) vet ./...

run:
	$(GO) run ./cmd/api

swagger:
	swag init -g cmd/api/main.go -o docs/openapi

migrate-up:
	migrate -path db/migrations -database "$$DATABASE_URL" up

migrate-down:
	migrate -path db/migrations -database "$$DATABASE_URL" down 1
