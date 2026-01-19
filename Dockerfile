# Build stage
FROM golang:1.22-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata

# Set working directory
WORKDIR /app

# Copy go mod file
COPY go.mod ./

# Download dependencies (go.sum will be auto-generated if missing)
RUN go mod download

# Copy source code
COPY . .

# Build with optimizations
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s -X main.version=$(git describe --tags --always --dirty 2>/dev/null || echo 'dev')" \
    -o /app/polygo \
    ./cmd/server

# Final stage
FROM alpine:3.19

# Install runtime dependencies
RUN apk add --no-cache ca-certificates tzdata

# Create non-root user
RUN addgroup -g 1000 polygo && \
    adduser -u 1000 -G polygo -s /bin/sh -D polygo

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/polygo /app/polygo

# Copy config if exists
COPY --from=builder /app/config*.yaml /app/ 2>/dev/null || true

# Set ownership
RUN chown -R polygo:polygo /app

# Switch to non-root user
USER polygo

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# Run
ENTRYPOINT ["/app/polygo"]
