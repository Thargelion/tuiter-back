# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project

Tuiter Back is the backend API for a Twitter-clone Android app. Go, Chi router, GORM over MySQL.

## Commands

```shell
# Start local MySQL (required before running the server)
make local.up
make local.down

# Run the server (needs PORT, JWT_SECRET, MYSQL_USER, MYSQL_PASS, MYSQL_HOST, MYSQL_DB env vars set)
go build -o out/tuiter-back ./cmd/tuiter/main.go && ./out/tuiter-back

# Lint (config in .golangci.yaml, v2)
golangci-lint run

# Tests
go test ./...
go test ./internal/domain/tuit/...        # single package
go test -run TestName ./path/to/package   # single test

# Regenerate Swagger docs (cmd/tuiter/docs)
make docs.generate
```

There are currently no `*_test.go` files in the repo — CI (`go-test.yml`) still runs `go build ./...` and `go test ./...` on every PR to `main`/`develop`/`releases/**`, and `go-lint.yml` runs golangci-lint v2.6 the same way. Match those on any change.

## Architecture

Package-oriented design with inbound/outbound ports split across layers:

- **`cmd/tuiter/main.go`** — composition root. Wires config (from env vars), DB connection/migration, repositories, services, error handlers, HTTP handlers, and Chi routers, then starts the server. Read this first to see how any given feature is assembled end-to-end.
- **`internal/domain/<name>`** — core domain packages (`tuit`, `tuitpost`, `user`, `avatar`). Each holds its `model.go` (plain structs), `repository.go` (outbound port interface the domain expects, e.g. `tuit.Repository`), and sometimes `use_cases.go` (business logic operating only on the interfaces, no I/O specifics).
- **`internal/application/handlers`** — inbound HTTP layer (`net/http` handlers). Talk to services/use-cases, not to `mysql` directly. Includes cross-cutting middleware (`RequestTagger`, `ApiValidation`, pagination) and shared response/error helpers.
- **`internal/application/router`** — Chi route registration per domain (`tuit.go`, `users.go`, `login.go`, `likes.go`, `fileserver.go`), each taking a handler interface (defined locally in the router file) so routers stay decoupled from concrete handler types.
- **`internal/application/services`** — orchestration layer between handlers and domain repositories/use-cases (e.g. `UserPostService` composes `tuit.Repository` + a feed repository; `UserAuthenticator`/`UserService` wrap login and user management).
- **`internal/infrastructure/mysql`** — GORM entities and outbound port implementations (`UserRepository`, `TuitRepository`, `FeedRepository`), plus a MySQL-specific error translator (`NewErrorHandler`) that maps DB errors into domain-level errors.
- **`pkg/`** — shared, dependency-free-ish utilities usable outside `internal`: `security` (JWT issuing/validation, auth middleware, bcrypt error handling), `query` (pagination/query param parsing), `logging` (contextual logger), `instant` (injectable clock/timezone), `syserror` (typed system errors).

### Request flow

`chiRouter` → global middleware (recoverer, timeout, request tagger, API validation, request logger, CORS) → `/v1` routes. Public routes: `/login`, `/users`. Everything else is mounted under a sub-router guarded by `securityMiddleware.Middleware` (JWT auth via `pkg/security`), which puts the validated token in the request context under `security.TokenMan` for handlers to read.

### Error handling

Handlers depend on a single `NewErrorsHandler` composed from per-source handlers (`mysql.NewErrorHandler()`, `security.NewBcryptErrorHandler()`) that translate infrastructure errors into consistent API error responses — new infrastructure error sources should register a handler here rather than leaking raw errors to handlers.

### Adding a new domain feature

Follow the existing pattern for `tuit`/`user`: model + repository interface in `internal/domain/<name>`, GORM entity + repository impl in `internal/infrastructure/mysql`, service/use-case wiring in `internal/application/services` (only if orchestration across repos is needed), HTTP handler in `internal/application/handlers`, route registration in `internal/application/router`, then wire it all in `cmd/tuiter/main.go`.
