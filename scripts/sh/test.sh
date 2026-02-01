#!/usr/bin/env bash
set -euo pipefail

BASE_URL="${BASE_URL:-http://localhost:8080}"
ENDPOINTS=(hello cpu memory)

RATE="${RATE:-20}"
DURATION="${DURATION:-5s}"
SLEEP_BETWEEN="${SLEEP_BETWEEN:-0.1}"

command_exists() {
  command -v "$1" >/dev/null 2>&1
}

error() {
  echo "Error: $*" >&2
}

ensure_vegeta() {
  if command_exists vegeta; then
    return
  fi

  if ! command_exists go; then
    error "vegeta not found and Go is not available for automatic installation."
    error "Please install Go and run this script again."
    exit 1
  fi

  echo "Installing vegeta via go install..." >&2
  go install github.com/tsenart/vegeta@latest

  export PATH="$PATH:$(go env GOPATH)/bin"

  if ! command_exists vegeta; then
    error "Failed to install vegeta automatically. Please check PATH (GOPATH/bin)."
    exit 1
  fi
}

run_attack() {
  local pids=()

  for ep in "${ENDPOINTS[@]}"; do
    printf "GET %s/%s\n" "$BASE_URL" "$ep" \
      | vegeta attack \
          -rate="$RATE" \
          -duration="$DURATION" \
          -name="$ep" \
          >/dev/null &
    pids+=("$!")
  done

  for pid in "${pids[@]}"; do
    wait "$pid"
  done
}

ensure_vegeta

while true; do
  run_attack
  sleep "$SLEEP_BETWEEN"
done
