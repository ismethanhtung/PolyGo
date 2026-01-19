# Build stage
FROM golang:1.22-alpine AS builder

# Install build dependencies (git needed for version info)
RUN apk add --no-cache git ca-certificates tzdata

# Set working directory
WORKDIR /app

# Copy go.mod first for better caching
COPY go.mod ./

# Copy all source code
COPY . .

# Tidy dependencies, download and verify
# go mod tidy ensures go.sum has all entries based on actual imports
RUN GOPRIVATE=github.com/polygo go mod tidy && \
    go mod download && \
    go mod verify

# Build binary with optimizations
# -w -s: strip debug info and symbol table
# -trimpath: remove file system paths from binary
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-w -s -X main.version=$(git describe --tags --always --dirty 2>/dev/null || echo 'dev')" \
    -trimpath \
    -o /app/polygo \
    ./cmd/server

# Verify binary exists and is executable
RUN test -x /app/polygo || exit 1

# ============================================
# Runtime stage - minimal Alpine image
# ============================================
FROM alpine:3.19

# Metadata labels
LABEL maintainer="PolyGo Team" \
      org.opencontainers.image.title="PolyGo" \
      org.opencontainers.image.description="High-performance Polymarket API proxy with caching and WebSocket support" \
      org.opencontainers.image.vendor="PolyGo" \
      org.opencontainers.image.licenses="MIT"

# Install runtime dependencies and setup user in one layer
RUN apk add --no-cache ca-certificates tzdata wget && \
    addgroup -g 1000 polygo && \
    adduser -u 1000 -G polygo -s /bin/sh -D polygo && \
    rm -rf /var/cache/apk/* /tmp/*

# Set working directory
WORKDIR /app

# Copy binary from builder with correct ownership
COPY --from=builder --chown=polygo:polygo /app/polygo /app/polygo

# Switch to non-root user for security
USER polygo

# Expose application port
EXPOSE 8080

# Health check endpoint
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# Run application
CMD ["/app/polygo"]
