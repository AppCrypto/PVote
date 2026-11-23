# Aggregate bounded discrete-log decoding — 2026-11-22

This experiment measures **one decoded aggregate coordinate** in the inclusive interval `[0, 25,500,000]`: the possible sum of 100,000 scores in `[0,255]`. It samples ten distinct integers uniformly from this interval, constructs their BN254 group elements, and verifies exact recovery of every integer. It does not submit or generate 100,000 individual ballots.

The primary measurement calls the existing `boundedDLog` function used by `Recon` in `crypto/RBPVSS/rbpvss.go`. Its baby-step/giant-step (BSGS) table has **5,050 entries** for the **25,500,001** possible integers. Each primary observation includes fresh table construction and the search. Target sampling, group-element generation, warmups, result checking, and report writing are outside the timer.

| Measurement | Observations | Mean ± sample SD (ms) | Minimum–maximum (ms) |
|---|---:|---:|---:|
| Production decoder, including table construction | 10 | 271.563 ± 50.402 | 186.538–334.182 |
| Query with one reused baby-step table | 10 | 92.851 ± 50.038 | 10.405–158.992 |
| One-time construction of that reusable table | 1 | 175.286 (one observation) | — |

The reusable-table path is a measurement-only variant using the same group operations and BSGS search; the production reconstruction algorithm is unchanged and still constructs its table on each call. The reported SD is variation across ten different random targets, not a confidence interval or repeated timing of one fixed target. The cached table's recorded allocation counter includes transient allocations and is not a retained-memory or peak-footprint measurement.

## Platform and validation

- Apple M4 / Mac16,12, 10 logical CPUs, 16 GiB RAM, macOS 26.5.2, on AC power.
- Go 1.25.4, optimized **native darwin/arm64** executable, `GOMAXPROCS=1`, default garbage collection, no forced GC in the measured loop.
- Three untimed warmups per path verify 0, the midpoint, and the upper bound; both paths also reject the upper bound plus one. All ten random target decodings pass in both paths. Matched cold/cached measurements alternate their order.
- The generator, all target integers and group encodings, and the ChaCha8 seed initialized from `crypto/rand` are preserved in `raw.json`. There are no excluded observations.
- Hardware queries launched inside the measured executable were sandbox-blocked. The original failures in `raw.json` are preserved. `hardware_metadata.json` separately records the successful host query, while the executable itself records its native `arm64` architecture and macOS version.
- Thermal state, energy, network latency, interpolation, and mask removal are not measured here.

## Files and reproduction

`raw.json` is the original report. `samples.csv` retains all ten observations; `summary.csv` and `summary.json` give the statistics above. `manifest.json` records source and executable hashes, commands, build environment, and data checksums. The manuscript embeds the primary result directly in LaTeX and does not load these files.

Run from the repository root, choosing a new directory to preserve existing observations:

```sh
python3 scripts/run_bounded_dlog.py --out experiments/bounded-dlog-repeat
```

To replay exactly the same targets and generator, pass the recorded `seedHex` and `baseHex` values with `--seed` and `--base`. The harness is opt-in and skipped in ordinary Go test runs.
