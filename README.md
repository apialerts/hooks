# hooks.apialerts.com

Free, open-source webhook testing tool. Generate a unique URL, inspect incoming requests, and toggle between success/failure responses to test your retry logic.

**Live at [hooks.apialerts.com](https://hooks.apialerts.com)** | Built by [API Alerts](https://apialerts.com)

## Quick Start (Docker)

```bash
docker compose up --build
```

Open [http://localhost:3080](http://localhost:3080).

## Quick Start (Local)

Requires Go 1.26+ and a Postgres instance.

```bash
# Start just the database
docker compose up db -d

# Run the server
make dev
```

Open [http://localhost:3080](http://localhost:3080).

## Features

- **Request inspector** — headers, body, query params, source IP for every request
- **Response mode toggle** — 200, 201, 400, 401, 403, 404, 500, 503, or 30s timeout
- **Custom response body** — each preset has a sensible JSON default, fully editable
- **Transaction log** — plain text request/response view, copy or download as `.txt`
- **Auto-refresh** — 30-second polling with countdown bar and manual refresh
- **Test button** — send a sample request without leaving the browser
- **Dark mode** — toggle between light and dark themes
- **No sign-up** — fully anonymous, endpoints expire after 7 days of inactivity

## Tech Stack

- **Go** — Chi router, single binary with embedded assets
- **Postgres** — via pgx, migrations via Goose (run automatically on startup)
- **HTMX** — polling and partial page updates
- **Tailwind CSS** — compiled via standalone CLI (no Node.js required)

## Project Structure

```
hooks/
├── cmd/server/
│   ├── main.go                 # Entry point, router, config, graceful shutdown
│   └── static/                 # Embedded assets (favicon, HTMX, compiled CSS)
├── internal/
│   ├── handler/                # HTTP handlers and HTML templates
│   ├── db/                     # Postgres queries and Goose migrations
│   ├── cleanup/                # Background expired endpoint purge
│   └── middleware/             # Rate limiting, CORS
├── input.css                   # Tailwind source (v4 syntax)
├── tailwind.config.js          # Tailwind config
├── docker-compose.yml          # Local dev (Go app + Postgres)
├── Dockerfile                  # Multi-stage production build
└── Makefile
```

## Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | `3080` | Server port |
| `DATABASE_URL` | `postgres://hooks:hooks@localhost:5432/hooks?sslmode=disable` | Postgres connection string |
| `DB_USER` | — | Alternative: Postgres user (used if `DATABASE_URL` is not set) |
| `DB_PASSWORD` | — | Alternative: Postgres password |
| `DB_HOST` | — | Alternative: Postgres host |
| `DB_NAME` | — | Alternative: Postgres database name |
| `BASE_URL` | `http://localhost:{PORT}` | Public URL (used in templates and transaction logs) |

When `DATABASE_URL` is not set but `DB_USER`, `DB_PASSWORD`, `DB_HOST`, and `DB_NAME` are all provided, the connection string is built as `postgres://{user}:{password}@{host}:5432/{name}?sslmode=require`.

## Development

```bash
# Start Postgres
docker compose up db -d

# Run the server (auto-reloads on restart)
make dev

# Run tests
make test

# Rebuild Tailwind CSS after template changes
make css

# Watch Tailwind CSS during development
make css-watch

# Update vendored HTMX
make update-htmx HTMX_VERSION=2.0.4
```

## Self-Hosting

### Docker Compose (recommended)

```bash
docker compose up --build -d
```

Runs the Go app and Postgres together on port 3080. Data persists in a Docker volume.

### Docker (bring your own Postgres)

```bash
docker build -t hooks .

docker run -p 3080:3080 \
  -e DATABASE_URL="postgres://user:pass@your-db:5432/hooks?sslmode=require" \
  -e BASE_URL="https://hooks.yourdomain.com" \
  hooks
```

### VPS (DigitalOcean, Hetzner, etc.)

SSH into your server, clone the repo, and run Docker Compose:

```bash
git clone https://github.com/apialerts/hooks.git
cd hooks
BASE_URL=https://hooks.yourdomain.com docker compose up --build -d
```

Point your domain's DNS to the server IP. Use a reverse proxy like Caddy or nginx for HTTPS — Caddy handles SSL certificates automatically:

```
# Caddyfile
hooks.yourdomain.com {
    reverse_proxy localhost:3080
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

**Connecting to Cloud SQL:** Use the built-in Cloud SQL connector rather than a public IP. Add the `--add-cloudsql-instances` flag and use the Unix socket path as the host:

```bash
gcloud run deploy hooks \
  --image gcr.io/YOUR_PROJECT/hooks \
  --region us-central1 \
  --project YOUR_PROJECT \
  --allow-unauthenticated \
  --add-cloudsql-instances YOUR_PROJECT:us-central1:YOUR_INSTANCE \
  --set-env-vars "DATABASE_URL=postgres://user:pass@/hooks?host=/cloudsql/YOUR_PROJECT:us-central1:YOUR_INSTANCE,BASE_URL=https://hooks.yourdomain.com"
```

Alternatively, enable a [VPC connector](https://cloud.google.com/vpc/docs/configure-serverless-vpc-access) on the Cloud Run service and use the Cloud SQL private IP directly in `DATABASE_URL`.

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

- [Chi](https://github.com/go-chi/chi) — HTTP router
- [pgx](https://github.com/jackc/pgx) — PostgreSQL driver
- [Goose](https://github.com/pressly/goose) — database migrations
- [HTMX](https://htmx.org) — frontend interactivity
- [Tailwind CSS](https://tailwindcss.com) — styling
- [Petname](https://github.com/dustinkirkland/golang-petname) — URL slug generation

## License

MIT
