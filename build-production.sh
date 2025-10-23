#!/bin/bash

# Production Docker Build Script
# This script builds the production Docker image with proper versioning

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${GREEN}Building Production Docker Image${NC}"
echo "=================================="

# Get version info
VERSION=${VERSION:-$(git describe --tags --always --dirty 2>/dev/null || echo "dev")}
BUILD_TIME=$(date -u +"%Y-%m-%dT%H:%M:%SZ")
GIT_COMMIT=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")

echo -e "${YELLOW}Version:${NC} $VERSION"
echo -e "${YELLOW}Build Time:${NC} $BUILD_TIME"
echo -e "${YELLOW}Git Commit:${NC} $GIT_COMMIT"
echo ""

# Build the image
echo -e "${GREEN}Building Docker image...${NC}"
docker build \
  --build-arg VERSION="$VERSION" \
  --build-arg BUILD_TIME="$BUILD_TIME" \
  --build-arg GIT_COMMIT="$GIT_COMMIT" \
  --target production \
  -t lexicon-bo-api:latest \
  -t lexicon-bo-api:"$VERSION" \
  .

echo ""
echo -e "${GREEN}Build completed successfully!${NC}"
echo ""
echo "Tagged images:"
echo "  - lexicon-bo-api:latest"
echo "  - lexicon-bo-api:$VERSION"
echo ""
echo "To run the image locally:"
echo "  docker run --rm -p 8080:8080 --env-file .env lexicon-bo-api:latest"
echo ""
echo "To push to registry:"
echo "  docker tag lexicon-bo-api:latest your-registry/lexicon-bo-api:$VERSION"
echo "  docker push your-registry/lexicon-bo-api:$VERSION"
