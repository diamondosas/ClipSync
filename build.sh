#!/usr/bin/env bash

set -e

# Read version from VERSION file if present, otherwise default to "dev"
if [ -f "VERSION" ]; then
    VERSION=$(tr -d '[:space:]' < VERSION)
else
    VERSION="dev"
fi

echo "=========================================="
echo " Building ClipSync (Version: ${VERSION})"
echo "=========================================="

# Build for Linux
echo "==> Building for Linux (amd64)..."
GOOS=linux GOARCH=amd64 go build -ldflags="-X main.Version=${VERSION} -s -w" -o clipsync .
echo "    [✔] Linux build complete: clipsync"

# Build for Windows
echo "==> Building for Windows (amd64)..."
GOOS=windows GOARCH=amd64 go build -ldflags="-X main.Version=${VERSION} -s -w -H windowsgui" -o clipsync.exe .
echo "    [✔] Windows build complete: clipsync.exe"

echo "=========================================="
echo " Build successful!"
echo " Outputs:"
echo "   - clipsync     (Linux)"
echo "   - clipsync.exe (Windows)"  
echo "=========================================="
