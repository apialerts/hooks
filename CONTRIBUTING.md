# Contributing to Hooks

Thanks for your interest in contributing! This guide covers local development setup for working on the Go server and templates.

If you just want to run Hooks, see the [README](README.md) -- it's a single `docker compose up --build`.

## Prerequisites

- [Go 1.26+](https://go.dev/dl/)
- [Docker](https://docs.docker.com/get-docker/) and Docker Compose (for Postgres)

## Setup

1. Fork and clone the repository:

```bash
git clone https://github.com/<your-username>/hooks.git
cd hooks
```

2. Download the Tailwind CSS standalone CLI and start the database:

```bash
make setup
docker compose up db -d
```

3. Run the server:

```bash
make dev
```

This builds CSS, starts a Tailwind file watcher in the background, and runs the Go server. The watcher rebuilds CSS automatically when you change templates.

4. Open [http://localhost:8080](http://localhost:8080).

### Using your own Postgres

You don't have to use the Docker Postgres. Copy the example env file and configure it:

```bash
cp .env.example .env
```

Then edit `.env` with your connection details:

```bash
# Option 1: Full connection string
DATABASE_URL=postgres://user:pass@localhost:5432/hooks?sslmode=disable

# Option 2: Individual variables
DB_HOST=localhost
DB_USER=myuser
DB_PASSWORD=mypass
DB_NAME=hooks
DB_SSLMODE=disable
```

The Makefile loads `.env` automatically. The server runs migrations on startup, so you just need an empty database.

## Running Tests

```bash
make test
```

Tests run against unit-testable code and do not require a running database.

## Database & Migrations

Migrations live in `internal/db/migrations/` and use [Goose](https://github.com/pressly/goose). They run automatically every time the server starts.

To check migration status or roll back manually:

```bash
make migrate-status
make migrate-down
```

To reset the database entirely (e.g. after changing an existing migration during development):

```bash
docker compose down -v
docker compose up db -d
```

Then restart the server and migrations will re-run from scratch.

## Project Structure

```
hooks/
├── cmd/server/          # Entry point, router, embedded static assets
├── internal/
│   ├── handler/         # HTTP handlers and HTML templates (inline in Go files)
│   ├── db/              # Postgres queries and migrations
│   ├── cleanup/         # Background expiry job
│   └── middleware/       # Rate limiting, CORS
├── input.css            # Tailwind v4 source
├── docker-compose.yml
├── Dockerfile
└── Makefile
```

## Makefile Reference

| Command            | Description                                               |
|--------------------|-----------------------------------------------------------|
| `make setup`       | Download the Tailwind CSS standalone CLI for your platform |
| `make dev`         | Build CSS, start Tailwind watcher, run the Go server       |
| `make build`       | Compile binary to `bin/hooks`                              |
| `make test`        | Run all tests                                              |
| `make css`         | One-off Tailwind CSS build (no minification)               |
| `make css-watch`   | Tailwind CSS watch mode (no minification)                  |
| `make update-htmx` | Download latest vendored HTMX                              |

## Making Changes

1. Create a branch from `main`:

```bash
git checkout -b my-change
```

2. Make your changes and add tests where applicable.

3. Run the full check:

```bash
make test
go build ./...
```

4. Commit with a clear message describing **what** and **why**.

5. Open a pull request against `main`.

## Code Style

- Follow standard Go conventions (`gofmt`).
- HTML templates are inline in handler Go files -- no separate template files.
- Tailwind CSS v4 compiled via the standalone CLI. The Docker build handles minification for production.

## What to Contribute

- Bug fixes
- Performance improvements
- New response presets or inspector features
- Documentation improvements
- Test coverage

If you're planning a large change, open an issue first to discuss the approach.

## License

By contributing, you agree that your contributions will be licensed under the [MIT License](LICENSE).
