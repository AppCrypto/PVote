# PVote

PVote is the reference implementation of the RB-PVSS construction defined in
`paper/main.tex`. The cryptographic API and transcript types follow the paper
directly; the EVM contract verifies the same transcript atomically, and the
browser demo adds an optional, separate stake-settlement layer.

## Requirements

- Go 1.22+
- Ganache or ganache-cli for EVM tests, benchmarks, and the web demo
- Solidity 0.8.25 and `abigen` 1.14.x when regenerating contract bindings

## Canonical RB-PVSS implementation

`crypto/RBPVSS` is the only cryptographic implementation. Its exported
algorithms correspond to the paper as follows:

| Paper algorithm | Go function |
| --- | --- |
| `RB-PVSS.Setup(1^lambda,n,t,l,S)` | `rbpvss.Setup(lambda, n, t, l, S)` |
| `RB-PVSS.Share(pp,w_j)` | `rbpvss.Share(pp, scores)` |
| `RB-PVSS.DVerify(pp,I_j)` | `rbpvss.DVerify(pp, instance)` |
| `RB-PVSS.Decrypt(pp,C_ji,sk_i)` | `rbpvss.Decrypt(pp, ciphertext, secretKey)` |
| `RB-PVSS.PVerify(pp,i,sh_ji,C_ji,varphi_ji)` | `rbpvss.PVerify(pp, index, share, ciphertext, proof)` |
| `RB-PVSS.Recon(pp,{(i,sh_i)},{U_d},S')` | `rbpvss.Recon(pp, shares, maskedScores, interval)` |

The dealer output is one atomic `Instance` containing
`(V,C,Phi,U,Pi)`. `AggregateVerified` calls `DVerify` before updating either
aggregate, so an invalid or partially verified transcript cannot affect the
tally.

The implementation includes the details required by the construction:

- one degree-`<t` polynomial shared by PVSS and every packed range proof;
- positions `D={-l+1,...,0}` and `N={1,...,n}`;
- nonzero shareholder keys and an admissible, erased range-issuer key;
- the common PVSS dealer challenge from `Phi`;
- a freshly sampled word from the dual Reed--Solomon code in `DVerify`;
- range-proof challenges binding `E,U,v_d,F,U',C'` and explicit challenge
  recomputation;
- decryption-share challenge recomputation;
- integer encoding `iota(w)=w mod q`, including negative score intervals;
- mask reconstruction, removal, and bounded discrete-log decoding inside
  `Recon`.

`compile/contract/RBPVSSVerifier.sol` implements the corresponding on-chain
state machine. `submitRB` performs the full `DVerify` and updates both
aggregates in one transaction; any failed check reverts before state changes.
The contract also verifies decryption shares and performs interpolation/mask
removal. `chain/rbpvss` is a lossless encoding boundary, not a second protocol
implementation.

## Run and test

Run the complete Setup--Share--DVerify--Aggregate--Decrypt--PVerify--Recon
flow:

```bash
go run . -n 7 -l 2 -m 3 -a 0 -b 5
```

Run all tests:

```bash
go test ./...
```

Run the real JSON-RPC integration test with the Ganache configuration already
used by this repository:

```bash
./scripts/test_evm.sh
```

Reproduce all off-chain timings and on-chain gas results used in the paper:

```bash
./scripts/benchmark.sh
```

The repository also contains independent, paper-interface reproductions used
by the cross-protocol evaluation:

- `baselines/priscore`: Setup, KeyGen/KeyDerive, Commit, Vote, both proof
  verifiers, Self-Tallying, and abort recovery;
- `baselines/threepvss`: the native individual-commitment 3PVSS Setup, Share,
  Verify, optimistic and pessimistic reconstruction, packed PDL proof,
  binary-vote proofs, decryption-share proofs, and the explicitly labeled
  `kappa`-instance score adapter and Tally.


Run their round-trip and tamper tests with `go test ./baselines/...`. Run the
same-machine comparison of the implemented protocol phases with:

```bash
go run ./cmd/full-comparison-bench -runs 10
```

The scripts write raw results to `experiments/rbpvss_results.csv`,
`experiments/rbpvss_chain_results.csv`,
`experiments/voter_gas_by_m.csv`,
`experiments/normalized_comparison_results.csv`, and
`experiments/full_protocol_comparison.csv`. The normalized file applies the
published operation counts to one primitive calibration; the full-protocol
file times actual Setup-prepared protocol calls and records the serialized
working-set and supplementary participation counts; Figures 12--13 use the
corresponding score-domain and finalization-working-set rows. The
`voter_gas_by_m.csv` file fixes `(n, l) = (10, 1)` and records ten sequential
accepted-ballot transactions to validate that steady-state per-voter gas does
not grow with the accepted-ballot count.

### Supplementary gas accounting and 1000-voter measurements

The original ten-ballot CSVs are retained as historical measurements. New
experiments do not modify the manuscript or replace those files. Run:

```bash
go run ./cmd/rbpvss-gas-audit -voters 2 -runs 3 -out experiments/gas-breakdown-new
MPLCONFIGDIR=/private/tmp/pvote-matplotlib python3 scripts/summarize_gas_audit.py experiments/gas-breakdown-new --plot
go run ./cmd/rbpvss-gas-audit -voters 1000 -runs 1 -out experiments/voter-scaling-1000-new
MPLCONFIGDIR=/private/tmp/pvote-matplotlib python3 scripts/summarize_voter_scaling.py experiments/voter-scaling-1000-new --gas-audit experiments/gas-breakdown-new --plot
```

The audit uses actual signed transactions and `core.ApplyTransaction` receipts
on a native go-ethereum v1.14.6 in-memory chain, with London gas rules, fake PoW
consensus and a permissive block gas limit of 90,071,992,547. It executes the
checked-in Solidity bytecode; it is not an operation-count-only simulation or a
live network deployment. The legacy transcript format is retained, so this is
a gas experiment rather than a test of the revised context binding. Voter keys
are deterministic public development fixtures, with independent ballot randomness
and zero scores in `[0,10]`. Every receipt's acceptance event is checked and the
final `acceptedCount` getter is independently evaluated from the final state.

The three `(n,l,t)` configurations are `(10,5,7)`, `(50,10,30)` and
`(100,19,59)`. Component experiments use three independent deployments each,
separating the first acceptance from the second. A 1000-voter run executes all
1000 consecutive ballots per configuration, with tracing at `m=1,2,1000`.
Outputs include all measurement rows, compressed signed transactions/receipts,
fork and bytecode metadata, and JSON opcode/component accounting. The model
checks EIP-1108 precompile prices, EIP-2028 calldata, EIP-2565 modular inverses,
EIP-2200/2929 storage and EIP-3529 refunds. It subtracts forwarded CALL gas to
avoid double counting. Memory expansion and other EVM instruction costs are
reported separately; the decomposition must equal receipt gas exactly.

For a live local Ganache reproduction with 1000 distinct funded accounts:

```bash
./scripts/run_voter_scaling_1000.sh experiments/voter-scaling-1000-ganache
```

This slower path uses Ganache 7.9.2, Merge execution rules and the same
permissive gas limit. Merge permits the constructor containing the complete
1000-account registry; the bytecode targets Paris. It records transactions and
receipts without changing verifier checks. Public transactions archived from
Ganache can also be replayed identically on native Geth:

```bash
go run ./cmd/rbpvss-gas-crosscheck -archive experiments/voter-scaling-1000-ganache/measurements.csv.receipts.jsonl.gz -out experiments/cross-client.csv
```

These experiments measure gas accounting and dependence on the accepted voter
count. They do not establish feasibility under current mainnet transaction gas
caps. See the supplementary `metadata.json`, `report.txt`, `summary.csv`, and
`gas_breakdown.pdf` for the actual measured scope and results.

The recorded supplementary runs are in `experiments/gas-breakdown` (18 traced
acceptances) and `experiments/voter-scaling-1000-native` (3000 actual acceptances).
The latter's `receipt_reconciliation.csv` checks the independent component
model against every scaling receipt after inserting its calldata byte counts
and first-write storage premium. Separate `response_draft.tex` files provide
proposed replies; no manuscript text or original figure data is changed.

The RB-PVSS tests cover single-dealer reconstruction, aggregate-only
reconstruction, nonconsecutive tallier subsets, negative intervals, tampered
dealer/range/decryption challenges, altered masked commitments, altered packed
PVSS commitments, invalid parameters, out-of-range scores, and failed bounded
discrete-log decoding.

## Web demo and optional settlement

The web demo creates proofs through `crypto/RBPVSS`. If Ganache is available,
it deploys `RBPVSSVerifier.sol`, submits each accepted ballot and decryption
share to that verifier, cross-checks its reconstructed points, and separately
uses `web/contract/StakeManager.sol` for deposits and rewards. Settlement never
replaces `DVerify` or `PVerify`.

Generate demo accounts and start Ganache:

```bash
bash genPrvKey_Mac.sh
ganache --mnemonic "PVote" -l 90071992547 -e 1000
```

On Linux, use `bash genPrvKey_Linux.sh`. Then run:

```bash
go run ./web
```

Open the printed local URL, normally `http://localhost:8080`.

Regenerate the settlement binding only after changing its Solidity source:

```bash
cd web/contract
bash compile.sh
```

Regenerate the RB-PVSS verifier binding after changing its source:

```bash
./compile/compile.sh
```

The removed legacy verification contract split one RB-PVSS transcript across
multiple transactions and did not recompute all Fiat--Shamir challenges. It was
therefore not an implementation of the paper's atomic `DVerify` interface.

## iPhone voter app

The native SwiftUI app and on-device Go RB-PVSS benchmarks are in [`mobile/`](mobile/README.md). Open `mobile/ios/PVote.xcodeproj` to build for iPhone. Public fixtures, local-chain election preparation, raw JSON/CSV export, and the physical-device measurement procedure are documented there. Device measurements must be collected on an actual iPhone; host tests and simulator builds are not mobile performance results.
