# Multi-stage Dockerfile for HUNTIQ Go Backend API
FROM golang:1.24-alpine AS builder

WORKDIR /build

RUN apk add --no-cache git ca-certificates

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /build/huntiq-api ./cmd/api

# Final minimal runtime image
FROM alpine:3.20

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

COPY --from=builder /build/huntiq-api /app/huntiq-api

# Uploads directory
RUN mkdir -p /app/uploads

EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=5s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

ENTRYPOINT ["/app/huntiq-api"]
