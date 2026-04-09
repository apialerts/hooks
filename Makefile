# Load .env file if it exists
-include .env
export

.PHONY: dev build run clean test migrate

# Development
dev: css
	./tailwindcss -i input.css -o cmd/server/static/styles.css --watch &
	go run ./cmd/server

# Build
build:
	go build -o bin/hooks ./cmd/server

# Run built binary
run: build
	./bin/hooks

# Test
test:
	go test -v ./...

# Clean
clean:
	rm -rf bin/

DATABASE_URL ?= postgres://hooks:hooks@localhost:5432/hooks?sslmode=disable

# Database migrations
migrate-up:
	cd internal/db && goose postgres "$(DATABASE_URL)" up

migrate-down:
	cd internal/db && goose postgres "$(DATABASE_URL)" down

migrate-status:
	cd internal/db && goose postgres "$(DATABASE_URL)" status

# Docker
docker-up:
	docker compose up -d

docker-down:
	docker compose down

TAILWIND_VERSION ?= 4.2.2

# Download Tailwind CSS standalone CLI for the current platform
setup:
	@if [ -f ./tailwindcss ]; then echo "tailwindcss already exists"; exit 0; fi; \
	OS=$$(uname -s | tr '[:upper:]' '[:lower:]'); \
	ARCH=$$(uname -m); \
	if [ "$$ARCH" = "arm64" ] || [ "$$ARCH" = "aarch64" ]; then TWARCH="arm64"; else TWARCH="x64"; fi; \
	if [ "$$OS" = "darwin" ]; then TWOS="macos"; else TWOS="linux"; fi; \
	URL="https://github.com/tailwindlabs/tailwindcss/releases/download/v$(TAILWIND_VERSION)/tailwindcss-$$TWOS-$$TWARCH"; \
	echo "Downloading tailwindcss v$(TAILWIND_VERSION) ($$TWOS-$$TWARCH)..."; \
	curl -sL "$$URL" -o ./tailwindcss && chmod +x ./tailwindcss; \
	echo "Done: ./tailwindcss"

HTMX_VERSION ?= 2.0.4

# Update vendored HTMX
update-htmx:
	curl -sL https://unpkg.com/htmx.org@$(HTMX_VERSION)/dist/htmx.min.js -o cmd/server/static/htmx.min.js

# Tailwind (requires tailwindcss standalone CLI)
# CSS is rebuilt automatically by `make dev` and during Docker builds.
css:
	./tailwindcss -i input.css -o cmd/server/static/styles.css

css-watch:
	./tailwindcss -i input.css -o cmd/server/static/styles.css --watch
