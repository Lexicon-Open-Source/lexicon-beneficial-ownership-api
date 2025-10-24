# Build stage
FROM golang:1.21-bookworm AS builder

# Install build dependencies
RUN apt-get update && apt-get install -y --no-install-recommends git ca-certificates tzdata && rm -rf /var/lib/apt/lists/*

WORKDIR /build

# Copy go mod files and download dependencies (better layer caching)
COPY go.mod go.sum ./
RUN go mod download && go mod verify

# Copy source code
COPY . .

# Build arguments for version info
ARG VERSION=dev
ARG BUILD_TIME
ARG GIT_COMMIT

# Build the application with optimizations
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s -X main.Version=${VERSION} -X main.BuildTime=${BUILD_TIME} -X main.GitCommit=${GIT_COMMIT}" \
    -trimpath \
    -o /app/bo-api \
    .

# Production stage - using Debian slim for better compatibility
FROM debian:12-slim AS production

# Install runtime dependencies
RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates \
    tzdata \
    curl \
    && groupadd -r appuser -g 1000 \
    && useradd -r -m -u 1000 -g appuser appuser \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/bo-api .

# Use non-root user for security
USER appuser:appuser

# Expose port (default 8080, can be overridden)
EXPOSE 8080

# Add metadata labels
LABEL org.opencontainers.image.title="Lexicon Beneficial Ownership API"
LABEL org.opencontainers.image.description="Production deployment for Beneficial Ownership API"
LABEL org.opencontainers.image.vendor="Lexicon Open Source"

# Health check (adjust the endpoint if needed)
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD curl -f http://localhost:8080/health || exit 1

# Set entrypoint
ENTRYPOINT ["/app/bo-api"]
