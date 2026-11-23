#!/bin/bash
set -euo pipefail

cd "$(dirname "$0")/.."
output_dir="${1:-experiments/voter-scaling-1000}"
voter_count="${2:-${PVOTE_M_COUNT:-1000}}"
configs="${3:-${PVOTE_M_CONFIGS:-10:5,50:10,100:19}}"
port="${4:-${PVOTE_M_PORT:-28545}}"
if [[ -e "$output_dir/measurements.csv" ]]; then
  echo "Refusing to overwrite $output_dir/measurements.csv; choose a fresh output directory." >&2
  exit 1
fi
if nc -z 127.0.0.1 "$port" >/dev/null 2>&1; then
  echo "Port $port is already in use; select PVOTE_M_PORT." >&2
  exit 1
fi
mkdir -p "$output_dir"
task_dir=$(mktemp -d /private/tmp/pvote-m1000.XXXXXX)
umask 077
server_pid=""
cleanup() {
  if [[ -n "$server_pid" ]]; then kill "$server_pid" 2>/dev/null || true; fi
  rm -f "$task_dir/accounts.json"
}
trap cleanup EXIT

GOCACHE="${GOCACHE:-/private/tmp/pvote-go-cache}" go build -o "$task_dir/chain-bench" ./cmd/rbpvss-chain-bench
# Deliberately reproduce the permissive local benchmark gas regime. These
# measurements do not establish feasibility under mainnet transaction gas caps.
# Merge matches the Paris-targeted bytecode and permits the constructor payload
# containing the full 1000-account registry; no contract checks are bypassed.
ganache --server.host 127.0.0.1 --server.port "$port" \
  --chain.chainId 1337 --chain.networkId 1337 --chain.hardfork merge \
  --miner.blockGasLimit 90071992547 --miner.callGasLimit 180000000 \
  --wallet.mnemonic PVote --wallet.totalAccounts "$voter_count" \
  --wallet.defaultBalance 1000 --wallet.accountKeysPath "$task_dir/accounts.json" \
  --logging.quiet >"$task_dir/ganache.log" 2>&1 &
server_pid=$!
for ((attempt=0; attempt<200; attempt++)); do
  if ! kill -0 "$server_pid" 2>/dev/null; then
    echo "Ganache stopped during startup; inspect $task_dir/ganache.log locally (contains development keys)." >&2
    exit 1
  fi
  if nc -z 127.0.0.1 "$port" >/dev/null 2>&1 && [[ -s "$task_dir/accounts.json" ]]; then break; fi
  sleep 0.1
done
if ! nc -z 127.0.0.1 "$port" >/dev/null 2>&1; then
  echo "Ganache startup timed out." >&2
  exit 1
fi
python3 - "$output_dir/source_hashes.json" <<'PY'
import hashlib, json, pathlib, sys
paths = ["cmd/rbpvss-chain-bench/main.go", "cmd/rbpvss-chain-bench/voter_scaling.go",
         "crypto/RBPVSS/rbpvss.go", "compile/contract/RBPVSSVerifier.sol",
         "compile/contract/RBPVSSVerifier.go", "scripts/run_voter_scaling_1000.sh"]
pathlib.Path(sys.argv[1]).write_text(json.dumps({p: hashlib.sha256(pathlib.Path(p).read_bytes()).hexdigest() for p in paths}, indent=2)+"\n")
PY
"$task_dir/chain-bench" -url "http://127.0.0.1:$port" \
  -only-voter-scaling -account-keys-file "$task_dir/accounts.json" \
  -voter-scaling-count "$voter_count" -voter-scaling-configs "$configs" \
  -execution-hardfork merge -voter-scaling-out "$output_dir/measurements.csv" \
  | tee "$output_dir/run.log"
