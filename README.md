# hooks.apialerts.com

Free, open-source webhook testing tool. Generate a unique URL, inspect incoming requests, and toggle between success/failure responses to test your retry logic.

**Live at [hooks.apialerts.com](https://hooks.apialerts.com)** | Built by [API Alerts](https://apialerts.com)

## Quick Start

```bash
docker compose up --build
```

Open [http://localhost:8080](http://localhost:8080). The app and Postgres run together, migrations run automatically, and data persists in a Docker volume.

## Features

- **Request inspector** -- headers, body, query params, source IP for every request
- **Response mode toggle** -- 200, 201, 400, 401, 403, 404, 500, 503, or 35s timeout
- **Custom response bodies** -- each status code saves its own body override independently
- **Transaction log** -- plain text request/response view, copy or download as `.txt`
- **Auto-refresh** -- 30-second polling with countdown bar and manual refresh
- **Test button** -- send a sample request without leaving the browser
- **Dark mode** -- automatic, based on system preference
- **No sign-up** -- fully anonymous, endpoints expire after 7 days of inactivity

## How It's Built

The goal is a single binary with zero runtime dependencies (besides Postgres) and zero JavaScript build tooling.

- **Go + embedded assets** -- HTML templates are defined inline in Go handler files using the standard library `html/template` package. Static assets (CSS, HTMX, favicon) are compiled into the binary via `go:embed`. No third-party template engine, no separate template files, no asset pipeline. Build it, ship it, run it.
- **HTMX instead of a JS framework** -- the request list polls via HTMX, request details load as HTML fragments, config saves via `fetch`. The entire frontend is server-rendered HTML with a few lines of vanilla JS. No React, no bundler, no `node_modules`.
- **Tailwind CSS standalone CLI** -- CSS is compiled by a single static binary, not a Node.js toolchain. The Docker build handles minification. In development, `make dev` runs a file watcher that rebuilds on template changes.
- **Postgres + Goose** -- pgx for queries, Goose for migrations. Migrations run automatically on every server start, so there's no manual migration step for self-hosters or contributors.
- **Stateless** -- no sessions, no auth, no cookies. Endpoints are identified by their slug and associated with the creator's IP for the sibling list. This makes it trivial to deploy on Cloud Run, Fly, or any container platform.

## Self-Hosting

### Docker Compose (recommended)

```bash
docker compose up --build -d
```

Runs the Go app and Postgres together. Data persists in a Docker volume. Migrations run automatically on startup.

To change the port:

```bash
PORT=8080 docker compose up --build -d
```

Update the `ports` mapping in `docker-compose.yml` to match: `"8080:8080"`.

### Docker (bring your own Postgres)

If you already have a Postgres instance:

```bash
docker build -t hooks .

docker run -p 8080:8080 \
  -e DATABASE_URL="postgres://user:pass@your-db:5432/hooks?sslmode=require" \
  -e BASE_URL="https://hooks.yourdomain.com" \
  hooks
```

The app creates its tables automatically on startup via Goose migrations. You just need an empty database.

### VPS (DigitalOcean, Hetzner, etc.)

```bash
git clone https://github.com/apialerts/hooks.git
cd hooks
BASE_URL=https://hooks.yourdomain.com docker compose up --build -d
```

Point your domain's DNS to the server IP. Use a reverse proxy like Caddy or nginx for HTTPS. Caddy handles SSL certificates automatically:

```
# Caddyfile
hooks.yourdomain.com {
    reverse_proxy localhost:8080
}
```

### Cloud Run

```bash
# Build and push
gcloud builds submit --tag gcr.io/YOUR_PROJECT/hooks --project YOUR_PROJECT

# Deploy
gcloud run deploy hooks \
  --image gcr.io/YOUR_PROJECT/hooks \
  --region us-central1 \
  --project YOUR_PROJECT \
  --allow-unauthenticated \
  --set-env-vars "DATABASE_URL=postgres://...,BASE_URL=https://hooks.yourdomain.com"
```

Cloud Run sets the `PORT` env var automatically. The app is stateless and scales to zero.

**Connecting to Cloud SQL:** Use the built-in Cloud SQL connector rather than a public IP:

```bash
gcloud run deploy hooks \
  --image gcr.io/YOUR_PROJECT/hooks \
  --region us-central1 \
  --project YOUR_PROJECT \
  --allow-unauthenticated \
  --add-cloudsql-instances YOUR_PROJECT:us-central1:YOUR_INSTANCE \
  --set-env-vars "DATABASE_URL=postgres://user:pass@/hooks?host=/cloudsql/YOUR_PROJECT:us-central1:YOUR_INSTANCE,BASE_URL=https://hooks.yourdomain.com"
```

## Environment Variables

| Variable       | Default                    | Description                                                                    |
|----------------|----------------------------|--------------------------------------------------------------------------------|
| `PORT`         | `8080`                     | Port the server listens on                                                     |
| `DATABASE_URL` | *(see below)*              | Full Postgres connection string                                                |
| `DB_PORT`      | `5432`                     | Postgres port (only used when `DATABASE_URL` is not set)                       |
| `DB_USER`      | --                         | Postgres user (only used when `DATABASE_URL` is not set)                       |
| `DB_PASSWORD`  | --                         | Postgres password (only used when `DATABASE_URL` is not set)                   |
| `DB_HOST`      | --                         | Postgres host (only used when `DATABASE_URL` is not set)                       |
| `DB_NAME`      | --                         | Postgres database name (only used when `DATABASE_URL` is not set)              |
| `DB_SSLMODE`   | `require`                  | Postgres SSL mode (only used when `DATABASE_URL` is not set)                   |
| `BASE_URL`     | `http://localhost:{PORT}`  | Public-facing URL (used in templates and transaction logs)                     |
| `KILL_SWITCH`  | `false`                    | Serve a maintenance page without connecting to the database (for abuse/outages) |

**How the connection string is built:**

- If `DATABASE_URL` is set, it is used directly.
- If `DB_USER`, `DB_PASSWORD`, `DB_HOST`, and `DB_NAME` are all set, the connection string is built as: `postgres://{user}:{password}@{host}:{DB_PORT}/{name}?sslmode={DB_SSLMODE}`
- Otherwise, it defaults to: `postgres://hooks:hooks@localhost:{DB_PORT}/hooks?sslmode=disable`

## Limits

- 5 endpoints per IP address
- 50 requests stored per endpoint
- 60 requests/minute rate limit per endpoint
- 256KB max request payload
- 7-day expiry from last activity

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md).

## Attribution

Built with:

- [Chi](https://github.com/go-chi/chi) -- HTTP router
- [pgx](https://github.com/jackc/pgx) -- PostgreSQL driver
- [Goose](https://github.com/pressly/goose) -- database migrations
- [HTMX](https://htmx.org) -- frontend interactivity
- [Tailwind CSS](https://tailwindcss.com) -- styling
- [Petname](https://github.com/dustinkirkland/golang-petname) -- URL slug generation

## License

MIT
