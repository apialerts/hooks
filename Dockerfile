FROM golang:1.26-alpine AS builder

RUN apk add --no-cache libstdc++ libgcc

WORKDIR /app

# Download platform-appropriate Tailwind CSS standalone CLI
ARG TARGETARCH
RUN if [ "$TARGETARCH" = "arm64" ]; then TWARCH="arm64"; else TWARCH="x64"; fi && \
    wget -qO /usr/local/bin/tailwindcss "https://github.com/tailwindlabs/tailwindcss/releases/download/v4.2.2/tailwindcss-linux-${TWARCH}-musl" && \
    chmod +x /usr/local/bin/tailwindcss

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build CSS from templates
RUN tailwindcss -i input.css -o cmd/server/static/styles.css --minify

RUN CGO_ENABLED=0 GOOS=linux go build -o /hooks ./cmd/server

FROM alpine:3.21

RUN apk --no-cache add ca-certificates
COPY --from=builder /hooks /hooks

EXPOSE 8080
CMD ["/hooks"]
