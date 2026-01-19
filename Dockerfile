# Build stage
FROM golang:1.22-alpine AS builder

# Install build dependencies and git (needed for version info)
RUN apk add --no-cache git ca-certificates tzdata

# Set working directory
WORKDIR /app

# Copy go mod files first for better cache layer
COPY go.mod go.sum* ./

# Copy source code (needed for go mod tidy to analyze imports)
COPY . .

# Update go.sum based on actual source code imports and download dependencies
# go mod tidy analyzes the source code and ensures go.sum has all required entries
RUN go mod tidy && \
    go mod download && \
    go mod verify

# Build binary with optimizations
# Remove debug info (-w) and symbol table (-s) for smaller binary
# Set version from git tag or default to 'dev'
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-w -s -X main.version=$(git describe --tags --always --dirty 2>/dev/null || echo 'dev')" \
    -trimpath \
    -o /app/polygo \
    ./cmd/server

# Verify binary was created
RUN ls -lh /app/polygo && \
    file /app/polygo

# Runtime stage - minimal Alpine image
FROM alpine:3.19

# Labels for metadata
LABEL maintainer="PolyGo Team" \
      org.opencontainers.image.title="PolyGo" \
      org.opencontainers.image.description="High-performance Polymarket API proxy" \
      org.opencontainers.image.vendor="PolyGo"

# Install only runtime dependencies and create user in one layer
RUN apk add --no-cache ca-certificates tzdata wget && \
    addgroup -g 1000 polygo && \
    adduser -u 1000 -G polygo -s /bin/sh -D polygo && \
    rm -rf /var/cache/apk/* /tmp/*

# Set working directory
WORKDIR /app

# Copy binary from builder stage
COPY --from=builder --chown=polygo:polygo /app/polygo /app/polygo

# Note: Config files can be mounted as volumes or use environment variables
# If config.yaml is needed, mount it at runtime: -v /path/to/config.yaml:/app/config.yaml

# Switch to non-root user
USER polygo

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# Use CMD instead of ENTRYPOINT for flexibility
CMD ["/app/polygo"]
