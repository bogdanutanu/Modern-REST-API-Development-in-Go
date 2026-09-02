#!/usr/bin/env bash

# api-login
# Run this once per session to set TOKEN in your current shell.
# Usage: eval "$(./api-login)"

set -eo pipefail

API_PROTOCOL="${API_PROTOCOL:-https}"
API_URL="${API_URL:-${API_PROTOCOL}://localhost:8888}"
USERNAME="${USERNAME:-admin}"
PASSWORD="${PASSWORD:-password}"

TOKEN=$(curl -s -X POST \
  -d "{\"username\": \"${USERNAME}\", \"password\": \"${PASSWORD}\"}" \
  -H "Content-Type: application/json" \
  "${API_URL}/login" | jq -r '.token')

if [ -z "$TOKEN" ] || [ "$TOKEN" = "null" ]; then
  echo "echo 'Failed to obtain token' >&2" 
  exit 1
fi

echo "export TOKEN='$TOKEN'"