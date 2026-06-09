# Repository Guidelines

## Project Structure & Module Organization

This repository contains a Go API and a Vite/React frontend for Gymotric. Go entry points live in `cmd/`: `cmd/api` runs the API, `cmd/goose` runs migrations, and `cmd/seed` seeds data. Domain packages are under `internal/` with `handler.go`, `service.go`, `repository.go`, `models.go`, and `*_test.go` files. SQL query definitions live in `queries/`; generated sqlc output is in `internal/repository/`. Goose migrations are in `cmd/goose/migrations/`. Frontend source is in `frontend/src`, with routes in `frontend/src/routes`, components in `frontend/src/components`, and assets in `frontend/public`.

## Build, Test, and Development Commands

- `make all`: build the API and run backend tests.
- `make build`: compile `cmd/api/main.go` to `./main`; also renders MJML email templates.
- `make run`: run the Go API and start the frontend Vite dev server.
- `make run-api`: run only the Go API.
- `make test`: run `go test ./... -v`.
- `make test-coverage`: generate `testresults/cover.html`.
- `make docker-run` / `make docker-down`: start or stop Docker Compose.
- `make db-create`, `make db-up`, `make db-reset`: create and migrate SQLite.
- `npm run dev --prefix frontend`: start only the frontend.
- `npm run build --prefix frontend`: type-check and build the frontend.
- `npm run lint --prefix frontend`: run ESLint.

## Coding Style & Naming Conventions

Format Go code with `gofmt`; use tabs for Go indentation and lowercase package names. Prefer explicit error wrapping with `fmt.Errorf("...: %w", err)`. Keep domain logic in `service.go`, HTTP handling in `handler.go`, and database access behind repository interfaces. Do not edit `internal/repository/*.sql.go` by hand; update `queries/` and regenerate with sqlc. Frontend components use PascalCase exports, TypeScript types, and Tailwind utility classes.

## Testing Guidelines

Backend tests use Go’s `testing` package with `stretchr/testify`. Place tests beside the package under test and name files `*_test.go`; use `Test...` function names. Run `make test` before submitting backend changes. For frontend changes, run `npm run lint --prefix frontend` and `npm run build --prefix frontend`; no frontend test suite is configured.

## Commit & Pull Request Guidelines

Recent commits use short conventional prefixes such as `fix:`, `feature:`, and `chore:`. Keep messages concise and behavior-focused, for example `fix: trim whitespace on exercise type names`. Pull requests should describe the change, list test/build commands run, link related issues, and include screenshots for visible frontend changes.

## Security & Configuration Tips

Local setup copies `.default.env` to `.env`; do not commit secrets or environment-specific credentials. Keep authentication, token, and email changes covered by focused tests where practical.
