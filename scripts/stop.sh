#!/usr/bin/env bash
# Stop script for development server

# Kill air and its child processes
pkill -f "air.*\.air\.toml" 2>/dev/null || true
# Kill any remaining api processes
pkill -f "tmp/bin/api" 2>/dev/null || true
# Kill any jq processes related to our pipeline
pkill -f "jq.*fromjson" 2>/dev/null || true
# Clean up air's temporary files
air -c .air.toml -s stop 2>/dev/null || true

echo "Development server stopped."

