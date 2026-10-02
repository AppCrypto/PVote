#!/usr/bin/env bash
set -euo pipefail
cd "$(dirname "$0")"
solc --evm-version paris --optimize --via-ir --abi contract/RBPVSSVerifier.sol -o contract --overwrite
solc --evm-version paris --optimize --via-ir --bin contract/RBPVSSVerifier.sol -o contract --overwrite
abigen --abi=contract/RBPVSSVerifier.abi --bin=contract/RBPVSSVerifier.bin --pkg=contract --out=contract/RBPVSSVerifier.go
