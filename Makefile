.PHONY: dev build run clean test migrate

# Development
dev:
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

HTMX_VERSION ?= 2.0.4

# Update vendored HTMX
update-htmx:
	curl -sL https://unpkg.com/htmx.org@$(HTMX_VERSION)/dist/htmx.min.js -o cmd/server/static/htmx.min.js

# Tailwind (requires tailwindcss standalone CLI)
css:
	./tailwindcss -i input.css -o cmd/server/static/styles.css --minify

css-watch:
	./tailwindcss -i input.css -o cmd/server/static/styles.css --watch
