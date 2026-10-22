# ── Build stage ───────────────────────────────────────────────────────────────
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Cache dependency downloads
COPY go.mod go.sum ./
RUN go mod download

# Copy all source
COPY . .

# Build all binaries in one RUN so any source change busts the whole layer
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o api          ./cmd/api && \
    CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o migrate      ./cmd/migrate && \
    CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o seed_admin   ./cmd/seed_admin && \
    CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o seed_products ./cmd/seed_products

# ── Runtime stage ─────────────────────────────────────────────────────────────
FROM alpine:3.19

# ca-certificates required for TLS connection to Nhost PostgreSQL
RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

COPY --from=builder /app/api          ./api
COPY --from=builder /app/migrate      ./migrate
COPY --from=builder /app/seed_admin   ./seed_admin
COPY --from=builder /app/seed_products ./seed_products
COPY --from=builder /app/migrations     ./migrations

RUN mkdir -p /app/uploads

COPY start.sh ./start.sh
RUN chmod +x ./start.sh

EXPOSE 8080

CMD ["./start.sh"]
