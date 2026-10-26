# ── Build stage ───────────────────────────────────────────────────────────────
FROM golang:1.26-alpine AS builder

WORKDIR /app

# ca-certificates and tzdata for build/downloads
RUN apk add --no-cache git ca-certificates tzdata

# Cache dependency downloads
COPY go.mod go.sum ./
RUN go mod download

# Copy all source
COPY . .

# Build API binary
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/bin/api ./cmd/api

# ── Runtime stage ─────────────────────────────────────────────────────────────
FROM alpine:3.20

# ca-certificates required for TLS connection to PostgreSQL / external services
RUN apk add --no-cache ca-certificates tzdata

# Create non-root system user
RUN addgroup -g 10001 huntiq && \
    adduser -u 10001 -G huntiq -s /bin/sh -D appuser

WORKDIR /app

COPY --from=builder /app/bin/api ./api

# Setup uploads directory with appropriate permissions
RUN mkdir -p /app/uploads && chown -R appuser:huntiq /app

USER appuser

EXPOSE 8080

ENTRYPOINT ["./api"]

