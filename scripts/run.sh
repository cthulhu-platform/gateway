#!/usr/bin/env bash
# Wrapper script to run the API binary and pipe output through jq for JSON formatting

# Run the API and pipe through jq for JSON formatting
# If jq fails (non-JSON output), it will still pass through the original line
exec ./tmp/bin/api 2>&1 | jq -R 'try fromjson catch .'
