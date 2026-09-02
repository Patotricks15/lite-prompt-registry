# Build stage
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Download dependencies
COPY go.mod ./
RUN go mod download

# Copy source code
COPY cmd/ ./cmd/
COPY internal/ ./internal/

# Build static binary
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /app/prompt-registry ./cmd/prompt-registry

# Runtime stage
FROM alpine:3.20

RUN apk add --no-cache curl ca-certificates

WORKDIR /app

COPY --from=builder /app/prompt-registry /app/prompt-registry

EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=5s --start-period=5s --retries=3 \
    CMD curl -f http://localhost:8080/healthz || exit 1

ENTRYPOINT ["/app/prompt-registry"]
