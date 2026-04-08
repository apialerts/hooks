# hooks.apialerts.com

Free, open-source webhook testing tool. Generate a unique URL, inspect incoming requests, and toggle between success/failure responses to test your retry logic.

Built by [API Alerts](https://apialerts.com).

## Quick Start

```bash
docker compose up --build
```

Open [http://localhost:8080](http://localhost:8080). Click "Create Endpoint" to get a unique URL like `http://localhost:8080/brave-lazy-fox`. Send webhooks to it and watch them appear in the request log.

## How It Works

Each endpoint has a **dual-purpose URL**:
- **Browser** (GET with `Accept: text/html`) — serves the endpoint UI
- **Any other request** (POST, PUT, DELETE, etc.) — webhook receiver, logs the request and returns the configured response

### Features

- **Request inspector** — view method, headers, body (formatted JSON), query params, source IP
- **Response mode dropdown** — toggle between 200, 201, 400, 401, 403, 404, 500, 503, and timeout (30s)
- **Custom response body** — each preset has a default JSON body, editable per endpoint
- **Transaction log** — Chucker-style plain text view with Copy and Download buttons
- **Auto-refresh** — 30-second polling with countdown bar and "Refresh Now" button
- **Test button** — send a sample request without leaving the browser
- **Delete button** — instantly delete endpoint and all data
- **Dark mode** — toggle with moon/sun icon, matches API Alerts branding
- **No sign-up** — fully anonymous, endpoints persist for 14 days from last activity

### Limits

- 500 requests stored per endpoint (oldest trimmed)
- 60 requests/minute rate limit per endpoint
- 256KB max payload size
- 5 endpoints per IP address
- 14-day expiry from last activity

## Tech Stack

- **Go** — Chi router, single binary
- **Postgres** — via pgx, migrations via Goose
- **HTMX** — polling and partial page updates
- **Tailwind CSS** — via CDN (dev), standalone CLI (production)
- **No Node.js, no React, no SPA**

## Project Structure

```
hooks/
├── cmd/server/
│   ├── main.go                 # Entry point, router, config, graceful shutdown
│   └── static/                 # Embedded assets (favicon, HTMX, CSS)
├── internal/
│   ├── handler/
│   │   ├── handler.go          # Response presets, reserved paths
│   │   ├── layout.go           # Shared HTML layout (header, footer, dark mode, branding)
│   │   ├── home.go             # GET / (landing page), POST /endpoints (create)
│   │   ├── endpoint.go         # GET /{id} (UI), PUT /{id}/config, DELETE /{id}/delete, POST /{id}/test
│   │   ├── webhook.go          # ANY /{id} (webhook receiver)
│   │   ├── poll.go             # GET /{id}/requests (list), GET /{id}/requests/{seq} (detail)
│   │   ├── static.go           # Privacy page, 404, robots.txt, sitemap.xml
│   │   └── slug.go             # Random word URL generator (petname)
│   ├── db/
│   │   ├── db.go               # Postgres connection + Goose migration runner
│   │   ├── endpoints.go        # Endpoint CRUD queries
│   │   ├── requests.go         # Request CRUD queries
│   │   └── migrations/         # SQL migrations (Goose)
│   ├── cleanup/
│   │   └── cleanup.go          # Background goroutine, purges endpoints inactive 14+ days
│   └── middleware/
│       └── ratelimit.go        # Per-endpoint + per-IP rate limiting
├── docker-compose.yml          # Go app + Postgres (one command setup)
├── Dockerfile                  # Multi-stage build (~30MB image)
└── Makefile                    # dev, build, migrate, docker commands
```

## Routes

| Method | Path | Description |
|--------|------|-------------|
| GET | `/` | Landing page |
| POST | `/endpoints` | Create new endpoint, redirect to `/{id}` |
| GET | `/privacy` | Privacy policy |
| GET | `/health` | Health check |
| GET | `/{id}` | Endpoint UI (browser) or webhook receiver (API) |
| POST/PUT/PATCH/DELETE | `/{id}` | Webhook receiver — logs request, returns configured response |
| PUT | `/{id}/config` | Update response mode and body |
| POST | `/{id}/test` | Send a test request to the endpoint |
| DELETE | `/{id}/delete` | Delete endpoint and all data |
| GET | `/{id}/requests` | HTMX fragment: request list (polled every 30s) |
| GET | `/{id}/requests/{seq}` | HTMX fragment: request detail (transaction view) |

## Database

Two tables, managed by Goose migrations:

**endpoints** — webhook endpoint configuration
- `id` (text PK) — human-readable slug (e.g. `brave-lazy-fox`)
- `creator_ip` — for per-IP rate limiting
- `response_status` — HTTP status to return (default 200)
- `response_delay_ms` — delay before responding (for timeout testing)
- `response_body` — custom JSON response body
- `request_count` — atomic counter, used for per-endpoint sequence numbers
- `last_activity_at` — updated on each webhook, used for 14-day expiry

**requests** — logged webhook requests
- `seq` — per-endpoint sequence number (1, 2, 3...)
- `method`, `path`, `query_params`, `headers`, `body` — incoming request data
- `response_status`, `response_body` — what was sent back
- `source_ip`, `content_type`, `body_size` — metadata

## Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | `8080` | Server port |
| `DATABASE_URL` | `postgres://hooks:hooks@localhost:5432/hooks?sslmode=disable` | Postgres connection string |
| `BASE_URL` | `http://localhost:{PORT}` | Public URL (for display and transaction logs) |

## Development

```bash
# Start Postgres
docker compose up db -d

# Run the server
make dev

# Or build and run
make build
./bin/hooks
```

## Deployment (Cloud Run)

Single Go binary, connects to Cloud SQL Postgres. Stateless, scales to zero.

```bash
docker build -t hooks .
# Deploy to Cloud Run with DATABASE_URL and BASE_URL env vars
```

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for setup instructions and guidelines.

## Attribution

Built with these open-source projects:

- [Chi](https://github.com/go-chi/chi) — HTTP router
- [pgx](https://github.com/jackc/pgx) — PostgreSQL driver
- [Goose](https://github.com/pressly/goose) — database migrations
- [HTMX](https://htmx.org) — frontend interactivity
- [Tailwind CSS](https://tailwindcss.com) — styling
- [Petname](https://github.com/dustinkirkland/golang-petname) — URL slug generation

## License

MIT
