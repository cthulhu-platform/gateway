#!/usr/bin/env bash
# Development script that runs air with proper cleanup on Ctrl+C

# Add Go bin directory to PATH if air was installed via go install
export PATH="${PATH}:${HOME}/go/bin:$(go env GOPATH)/bin"

# Function to cleanup resources
cleanup() {
    echo ""
    echo "Cleaning up resources..."
    # Kill air and its child processes
    pkill -P $$ 2>/dev/null || true
    # Kill any remaining api processes
    pkill -f "tmp/bin/api" 2>/dev/null || true
    # Kill any jq processes related to our pipeline
    pkill -f "jq.*fromjson" 2>/dev/null || true
    # Clean up air's temporary files
    air -c .air.toml -s stop 2>/dev/null || true
    # Add any other cleanup commands here
    echo "Cleanup complete."
    exit 0
}

# Set up signal handlers for cleanup
trap cleanup SIGINT SIGTERM

# Make run.sh executable
chmod +x ./scripts/run.sh

# Run air with the configuration
# Use exec to ensure signals are properly forwarded
air -c .air.toml