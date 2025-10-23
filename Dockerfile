# Build stage
FROM golang:1.21-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata

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
    -o /app/definition \
    .

# Production stage - using Alpine for smaller size and better compatibility
FROM alpine:3.19 AS production

# Install runtime dependencies
RUN apk add --no-cache \
    ca-certificates \
    tzdata \
    && addgroup -g 1000 appuser \
    && adduser -D -u 1000 -G appuser appuser

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/definition .

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
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# Set entrypoint
ENTRYPOINT ["/app/definition"]
