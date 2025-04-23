#!/bin/bash

# Get the version from package.json
VERSION=$(node -e "console.log(require('./package.json').version)")

# Set other build information
COMMIT=$(git rev-parse --short HEAD)
BUILD_DATE=$(date -u +"%Y-%m-%dT%H:%M:%SZ")

# Define the output file name with version suffix
VERSIONED_OUTPUT="management-console-${VERSION}"

# Ensure bin directory exists
mkdir -p bin

# Build the Go binary with version information
go build -buildvcs=false \
  -ldflags "-X github.com/unity-sds/unity-management-console/backend/internal/version.Version=$VERSION" \
  -o "./bin/${VERSIONED_OUTPUT}" ./backend/cmd/web

# Create a symlink to the versioned binary
ln -sf "${VERSIONED_OUTPUT}" "./bin/management-console"

echo "Built Unity Management Console v$VERSION (commit: $COMMIT, built: $BUILD_DATE)"
echo "Created binary: bin/${VERSIONED_OUTPUT}"
echo "Created symlink: bin/management-console -> ${VERSIONED_OUTPUT}"