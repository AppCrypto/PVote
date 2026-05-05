# PVote

PVote is a Go implementation and demo workspace for a privacy-preserving voting protocol. It includes the core cryptographic modules, a command-line demo, Solidity contract bindings, and a single-page web demo that simulates initiator, voter, and tallier roles.

## Requirements

- Go 1.22+
- Ganache or ganache-cli
- Solidity compiler, if you need to regenerate contract artifacts
- `abigen`, if you need to regenerate Go bindings from Solidity contracts

Install `abigen`:

```bash
go install github.com/ethereum/go-ethereum/cmd/abigen@v1.14.3
```

## Project Layout

- `main.go`: command-line demo entry point.
- `crypto/PVSS`: PVSS implementation used for encrypted share generation and verification.
- `crypto/ZKRP`: zero-knowledge range proof implementation.
- `crypto/Convert`: helper conversions for cryptographic data.
- `utils`: shared utility code.
- `compile`: Solidity verification contract, ABI/bin artifacts, and generated Go binding for PVSS/ZKRP/PVerifyTally.
- `web`: Go backend and static frontend for the browser-based PVote demo.
- `web/static`: HTML, CSS, and JavaScript for the single-page role simulation UI.
- `web/contract`: stake manager Solidity contract, ABI/bin artifacts, and generated Go binding.
- `paper`: protocol notes and paper-related files.
- `genPrvKey_Linux.sh`, `genPrvKey_Mac.sh`: scripts for generating local demo accounts into `.env`.

## Setup

Generate or refresh `.env` accounts:

```bash
bash genPrvKey_Mac.sh
```

On Linux:

```bash
bash genPrvKey_Linux.sh
```

Start Ganache with the mnemonic expected by the demo accounts:

```bash
ganache --mnemonic "PVote" -l 90071992547 -e 1000
```

If you start Ganache with different accounts, regenerate `.env` so the private keys match the Ganache wallets.

## Run the Command-Line Demo

```bash
go run main.go
```

## Run the Web Demo

Start Ganache first, then run:

```bash
go run ./web
```

Open the printed local URL in a browser, usually:

```text
http://localhost:8080
```

At startup, the web backend deploys both the stake manager contract and `compile/contract/Verification.sol` to Ganache, then uploads the paper protocol public parameters and tallier public keys.

The web page provides tabs for:

- Initiator: configure protocol parameters, deploy/fund Ganache escrow, and inspect public parameters plus on-chain verification status.
- Voter: submit ballots, lock voter stake, and upload PVSS/ZKRP proofs to `Verification.sol`.
- Tallier: lock tallier stake, publish decryption shares from the on-chain aggregate, run `PVerifyTally`, and withdraw settled rewards.

## Regenerate Contract Bindings

Only run this when `web/contract/StakeManager.sol` changes:

```bash
cd web/contract
bash compile.sh
```

Only run this when `compile/contract/Verification.sol` changes:

```bash
cd compile
bash compile.sh
```

Then rerun the Go tests:

```bash
go test ./web ./crypto/... ./utils/...
```
