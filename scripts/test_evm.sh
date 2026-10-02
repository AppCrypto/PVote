#!/bin/sh
set -eu

log_file="${TMPDIR:-/tmp}/pvote-ganache-test.log"

# Keep the exact Ganache convention committed in the repository README. The
# checked-in .env was generated from the same "PVote" mnemonic.
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
PVOTE_EVM_URL="http://127.0.0.1:8545" PVOTE_EVM_PRIVATE_KEY="$private_key" \
  GOCACHE="${TMPDIR:-/tmp}/pvote-go-cache" go test -p 1 -v ./chain/rbpvss ./web
