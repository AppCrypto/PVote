#!/usr/bin/env bash
set -euo pipefail

cd "$(dirname "$0")"
rm -f ./*.bin ./*.abi ./*.go
Name=StakeManager
solc --evm-version paris --optimize --via-ir --abi "./$Name.sol" -o . --overwrite
solc --evm-version paris --optimize --via-ir --bin "./$Name.sol" -o . --overwrite
abigen --abi="./$Name.abi" --bin="./$Name.bin" --pkg=stakecontract --out="./$Name.go"
