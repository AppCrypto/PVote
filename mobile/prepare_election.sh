#!/bin/sh
# Disposable development chain, kept alive until Ctrl-C. No personal keys used.
set -eu
mobile_dir=$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)
if [ "$#" -ne 2 ]; then
  echo 'Usage: ./mobile/prepare_election.sh IPHONE_WALLET_ADDRESS MAC_LAN_IP' >&2
  exit 2
fi
mkdir -p "$mobile_dir/build"
wallet_file=$(mktemp "${TMPDIR:-/tmp}/pvote-mobile-wallet.XXXXXX")
chmod 600 "$wallet_file"
ganache --wallet.mnemonic 'PVote-mobile-research' --wallet.accountKeysPath "$wallet_file" \
 --wallet.defaultBalance 1000 --chain.chainId 1337 --chain.hardfork london \
 --server.host 0.0.0.0 --server.port 18545 --miner.blockGasLimit 90071992547 \
 --miner.callGasLimit 180000000 --logging.quiet > "$mobile_dir/build/dev-chain.log" 2>&1 &
chain_pid=$!
trap 'kill "$chain_pid" 2>/dev/null || true; rm -f "$wallet_file"' EXIT INT TERM
attempt=0
until nc -z 127.0.0.1 18545 >/dev/null 2>&1; do
 attempt=$((attempt+1))
 if [ "$attempt" -ge 100 ] || ! kill -0 "$chain_pid" 2>/dev/null; then
  echo "Development chain did not start; see mobile/build/dev-chain.log" >&2; exit 1
 fi
 sleep 0.1
done
deployer_key=$(python3 - "$wallet_file" <<'PY'
import json,sys
wallet=json.load(open(sys.argv[1]))
print(next(iter(wallet['private_keys'].values())).removeprefix('0x'))
PY
)
cd "$mobile_dir/.."
PVOTE_MOBILE_DEPLOYER_KEY="$deployer_key" GOCACHE="${TMPDIR:-/tmp}/pvote-go-cache" \
 go run ./cmd/pvote-mobile-election -voter "$1" -phone-rpc "http://$2:18545" -out "$mobile_dir/election.json"
unset deployer_key
echo 'Transfer mobile/election.json to the iPhone, import it, and keep this terminal running.'
echo 'This is a disposable LAN test chain with test funds; stop it with Ctrl-C.'
wait "$chain_pid"
