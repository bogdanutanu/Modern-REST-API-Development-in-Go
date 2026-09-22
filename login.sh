#!/usr/bin/env bash

# api-login
# Run this once per session to set TOKEN in your current shell.
# Usage: eval "$(./api-login)"

set -eo pipefail

API_PROTOCOL="${API_PROTOCOL:-http}"
API_URL="${API_URL:-${API_PROTOCOL}://localhost:8888}"
API_PATH="${API_PATH:-/api}"
USERNAME="${USERNAME:-admin}"
PASSWORD="${PASSWORD:-password}"

TOKEN=$(curl -sS -X POST \
  -d "{\"username\": \"${USERNAME}\", \"password\": \"${PASSWORD}\"}" \
  -H "Content-Type: application/json" \
  "${API_URL}${API_PATH}/login" | jq -r '.token')

if [ -z "$TOKEN" ] || [ "$TOKEN" = "null" ]; then
  echo "echo 'Failed to obtain token' >&2" 
  exit 1
fi

echo "export TOKEN='$TOKEN'"