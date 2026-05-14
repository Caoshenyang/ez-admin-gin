# Contributing to EZ Admin Gin

Thank you for your interest in contributing! This guide covers the development setup, code conventions, and PR process.

## Prerequisites

| Tool | Version | Notes |
|------|---------|-------|
| Go | 1.26+ | Backend |
| Node.js | 20+ | Frontend & docs |
| pnpm | 9+ | Package manager |
| Docker | 20+ | Local PostgreSQL + Redis |
| make | Any | Build automation (optional on Windows) |

## Quick Start

```bash
# 1. Start local databases
make docker-up

# 2. Start backend + frontend
make dev

# 3. Run tests
make test-contract          # No DB needed
make test-integration       # Needs DB + Redis
make test-e2e               # Needs running backend + frontend
```

## Project Structure

```
server/               # Go backend
  internal/
    bootstrap/        # App startup and wiring
    modules/          # Business modules (auth, system, ...)
    platform/         # Shared infrastructure (middleware, logger, authn, ...)
    pkg/              # Pure utility packages
  migrations/         # Database migrations (mysql/ + postgres/)
  tests/              # All backend tests (centralized)
    api/              # API black-box tests
    rbac/             # Permission & data scope tests
    contract/         # OpenAPI contract tests
    testutil/         # Test helpers
admin/                # Vue 3 frontend
  src/
    api/              # HTTP client & generated types
    modules/          # Feature modules
    layouts/          # App layouts
  e2e/                # Playwright E2E tests
docs/                 # VitePress documentation site
deploy/               # Docker Compose configs
```

## Branch Naming

| Type | Format | Example |
|------|--------|---------|
| Feature | `feat/short-description` | `feat/dark-mode` |
| Bug fix | `fix/short-description` | `fix/login-redirect` |
| Docs | `docs/short-description` | `docs/deploy-guide` |

## Commit Messages

Follow [Conventional Commits](https://www.conventionalcommits.org/):

```
feat: add user import from CSV
fix: resolve menu tree rendering for empty children
docs: update deployment guide with HTTPS setup
test: add role permission assignment E2E test
refactor: extract pagination helper from list handlers
chore: update Go dependencies
```

## Code Style

### Backend (Go)

- Follow standard Go conventions (`gofmt`, `go vet`)
- Package naming: lowercase, no underscores
- Error handling: use `errorsx` package for consistent error codes
- No co-located `*_test.go` files in business directories — tests go in `server/tests/`

### Frontend (Vue 3 + TypeScript)

- Use `<script setup>` with Composition API
- TypeScript strict mode
- Pages: `admin/src/modules/{module}/pages/{Name}View.vue`
- API types: auto-generated via `make generate-types` (do not edit `admin/src/api/generated.ts` manually)
- Lint: `pnpm lint` must pass

### Documentation

- Docs live in `docs/` (VitePress)
- Primary language: Chinese
- Run `cd docs && pnpm docs:build` to verify before submitting

## Testing

Tests are centralized, not co-located with business code:

| Command | What it runs | Requires DB |
|---------|--------------|-------------|
| `make test-contract` | OpenAPI contract tests | No |
| `make test-api` | API black-box tests | Yes |
| `make test-rbac` | Permission & data scope tests | Yes |
| `make test-integration` | All integration tests (API + RBAC) | Yes |
| `make test-e2e` | Playwright E2E tests | Running server |
| `make lint` | go vet + frontend type check + API type sync | No |

**Rules:**
- Do not create `*_test.go` files in business code directories.
- Do not use `t.Skip` to pass incomplete tests.
- Do not connect to production databases.
- Integration tests use real DB/Redis, not mocks.

## Pull Request Process

1. **Fork** the repository and create a feature branch.
2. **Make changes** following the code style above.
3. **Run checks** before pushing:
   ```bash
   make lint
   make test-contract
   # If your changes affect API/RBAC:
   make test-integration
   ```
4. **Open a PR** against the `main` branch.
5. Fill in the PR template completely.
6. Wait for CI to pass and code review.

### PR Checklist

- [ ] `make lint` passes
- [ ] Tests added or updated for the changed behavior
- [ ] No unrelated changes or refactoring
- [ ] Documentation updated if applicable
- [ ] `CHANGELOG.md` updated (under `[Unreleased]`)

## Key Files to Read

Before making significant changes, read these files:

- `CLAUDE.md` — Project rules and constraints
- `docs/project/PHASE_STATUS.md` — Current phase and progress
- `docs/project/QUALITY_ROADMAP.md` — Long-term roadmap
- `docs/project/TESTING_STRATEGY.md` — Testing conventions
- `docs/project/DECISION_LOG.md` — Architecture decisions

## Questions?

Feel free to open an issue with the `question` label if you need clarification on anything.
