# Contributing to Hooks

Thanks for your interest in contributing! This guide will help you get set up and submit your first pull request.

## Prerequisites

- [Go 1.26+](https://go.dev/dl/)
- [Docker](https://docs.docker.com/get-docker/) and Docker Compose

## Setup

1. Fork and clone the repository:

```bash
git clone https://github.com/<your-username>/hooks.git
cd hooks
```

2. Start the Postgres database:

```bash
docker compose up db -d
```

3. Run the server:

```bash
make dev
```

4. Open [http://localhost:8080](http://localhost:8080) in your browser.

## Running Tests

```bash
make test
```

Tests run against unit-testable code and do not require a running database.

## Project Structure

```
hooks/
├── cmd/server/          # Entry point, router, embedded static assets
├── internal/
│   ├── handler/         # HTTP handlers and HTML templates
│   ├── db/              # Postgres queries and migrations
│   ├── cleanup/         # Background expiry job
│   └── middleware/       # Rate limiting
├── docker-compose.yml
├── Dockerfile
└── Makefile
```

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
- Keep HTML templates inline in handler files — no template engine.
- Tailwind CSS via CDN in development, standalone CLI for production builds.

## What to Contribute

- Bug fixes
- Performance improvements
- New response presets or inspector features
- Documentation improvements
- Test coverage

If you're planning a large change, open an issue first to discuss the approach.

## License

By contributing, you agree that your contributions will be licensed under the [MIT License](LICENSE).
