FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /hooks ./cmd/server

FROM alpine:3.21

RUN apk --no-cache add ca-certificates
COPY --from=builder /hooks /hooks

EXPOSE 3080
CMD ["/hooks"]
