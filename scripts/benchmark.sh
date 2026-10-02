#!/bin/sh
set -eu

log_file="${TMPDIR:-/tmp}/pvote-ganache-benchmark.log"

# Reproduce the Ganache configuration already committed in README.md.
ganache --mnemonic "PVote" -l 90071992547 -e 1000 --logging.quiet >"$log_file" 2>&1 &
ganache_pid=$!
trap 'kill "$ganache_pid" 2>/dev/null || true' EXIT INT TERM

i=0
until nc -z 127.0.0.1 8545 >/dev/null 2>&1; do
  i=$((i + 1))
  if [ "$i" -ge 50 ]; then
    echo "Ganache did not start; see $log_file" >&2
    exit 1
  fi
  sleep 0.1
done

private_key=$(sed -n 's/^PRIVATE_KEY_1=//p' .env | head -n 1)
test -n "$private_key" || { echo "PRIVATE_KEY_1 is missing from .env" >&2; exit 1; }
private_keys=$(sed -n 's/^PRIVATE_KEY_[0-9][0-9]*=//p' .env | paste -sd, -)
test -n "$private_keys" || { echo "No PRIVATE_KEY_n entries found in .env" >&2; exit 1; }
mkdir -p experiments
if [ "${PVOTE_SKIP_OFFCHAIN:-0}" != "1" ]; then
  GOCACHE="${TMPDIR:-/tmp}/pvote-go-cache" go run ./cmd/rbpvss-bench -runs 20 -setup-runs 5
  GOCACHE="${TMPDIR:-/tmp}/pvote-go-cache" go run ./cmd/comparison-bench -runs 200
  GOCACHE="${TMPDIR:-/tmp}/pvote-go-cache" go run ./cmd/full-comparison-bench -runs 10
fi
PVOTE_EVM_URL="http://127.0.0.1:8545" PVOTE_EVM_PRIVATE_KEY="$private_key" PVOTE_EVM_PRIVATE_KEYS="$private_keys" \
  GOCACHE="${TMPDIR:-/tmp}/pvote-go-cache" go run ./cmd/rbpvss-chain-bench
