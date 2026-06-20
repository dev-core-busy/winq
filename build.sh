#!/usr/bin/env bash
# Cross-compile winq for Windows from Linux
set -e
GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -o winq-windows-amd64.exe .
echo "✓ winq-windows-amd64.exe"
