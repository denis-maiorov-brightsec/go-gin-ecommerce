# Go Gin Ecommerce Spec-Driven Scaffold

This repository is a spec-driven scaffold for building an ecommerce backoffice API in Go with Gin, GORM, PostgreSQL, and `swaggo/swag`.

It keeps an ecommerce-focused spec progression and drives implementation from `docs/STACK_PROFILE.md`.

## What this template includes

- `AGENTS.md`: execution protocol and guardrails for spec-by-spec implementation
- `docs/STACK_PROFILE.md`: Go/Gin/GORM/PostgreSQL stack contract and quality-gate commands (not API contracts)
- `docs/SPECS_INDEX.md`: dependency-ordered backlog with statuses
- `docs/specs/*.md`: implementation-agnostic spec seed set
- `scripts/run-specs-harness.mjs`: implementer + reviewer two-pass automation harness
- `prompts/*.md`: copy/paste prompts to adapt specs and scaffold a fresh project

## Expected usage

1. Fill `docs/STACK_PROFILE.md` with your target stack, commands, and exact repository topology paths.
2. Use `prompts/01-adapt-specs-and-scaffold.md` as the first Codex prompt.
3. Run the harness per spec (or range).

## Local development

```bash
cp .env.example .env.local
docker compose up -d postgres
go run ./cmd/api
```

## Common commands

```bash
make fmt
make lint
make test
make test-integration
make swagger
```

## Harness quick start

```bash
node scripts/run-specs-harness.mjs --dry-run
```

Run up to 3 specs:

```bash
node scripts/run-specs-harness.mjs \
  --max-specs 3
```

## Notes

- Harness expects `docs/SPECS_INDEX.md` table format and `docs/specs/<id>-*.md` files.
- By default, harness enforces clean git state and branch `main`.
- Implementer must create a commit; reviewer commits only when fixes are required.
