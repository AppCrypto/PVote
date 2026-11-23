# iPhone voter-generation measurements

Figure 6 and the current paper text use the three completed cooled rechecks below: **2,070 fresh ballots, 90 samples per configuration**. The earlier measurements are retained as historical records.

## Current data

- Device: iPhone 13 (`iPhone14,5`), iOS 26.6.2 (23G90), native arm64, Release; debugger disabled.
- Go 1.25.4, `GOMAXPROCS=1`, default garbage collection enabled; `context-v1` proofs and unchanged public parameter fingerprints.
- Three serial runs in one app process, in ascending `(l,n)` order, descending order, and ascending order cyclically shifted by half the list length. Each run has 23 configurations, 3 warmups and 30 measured ballots per configuration. The first warmup is locally verified; measured ballots use fresh randomness and omit local verification.
- Each run starts at nominal thermal state; runs are separated by at least 60 seconds of cooling and any additional wait for nominal state. All 2,070 samples record nominal state before and after measurement, low-power mode off, and operation on battery power.
- Setup, parameter loading, warmups, cooling, HTTP export, and disk writes are outside the Share timer. Encoding has a separate timer.
- The original JSON bytes were received by HTTP and SHA-256 checked. Full report names, checksums, source fingerprint, timestamps, and cooling waits are in `recheck-20261121/cooled/analysis_manifest.json`.

Reproduce the analysis from the repository root:

```sh
python3 mobile/summarize.py \
  experiments/iphone/recheck-20261121/cooled/raw/benchmark-*.json \
  --out experiments/iphone
```

`raw.csv` contains all 2,070 measurements. `summary.csv` contains 23 groups of 90 samples. SD is the sample standard deviation over all 90 observations, including variation between runs; P95 uses the nearest-rank definition. All means and SDs are embedded directly in `paper/main.R1.tex`; the paper does not load CSV files. The observations are not smoothed or adjusted to enforce monotonicity.

### Rechecked n=100 points

| l | t | Earlier mean ± SD (30 samples, ms) | Recheck mean ± SD (90 samples, ms) |
|---|---|---|---|
| 1 | 50 | 38.771 ± 0.596 | 40.345 ± 1.835 |
| 10 | 55 | 55.464 ± 0.910 | 53.705 ± 1.405 |
| 19 | 59 | 62.818 ± 0.305 | 68.280 ± 4.005 |

The corresponding l=19 run means are:

| Order | n=90 (ms) | n=100 (ms) |
|---|---|---|
| ascending | 67.854 | 72.050 |
| descending | 60.004 | 62.900 |
| rotated | 65.261 | 69.890 |

At `(100,19,59)`, pooled P95 is **73.036 ms**, with mean encoding **0.688 ms**. Across configurations, means range from **5.468 to 68.280 ms**. The maximum whole-app physical footprint sample is **47.970 MiB**, sampled at 2 ms during generation and encoding. This is a sampled lower bound on the whole-app peak.

## Earlier reports and exclusions

- Original report: `raw/benchmark-2026-11-20T23-28-25Z-CE89002F.json`, SHA-256 `59b94972457cea7d953239c5ceedca3936b49cbd77fec19607d412beecdd6db6`. Its 690-sample raw and summary CSVs are archived in `initial-20261121/`. It used one fixed fixture order, placing `(100,19,59)` second.
- The initial recheck before the cooling policy produced one complete 690-sample run in `recheck-20261121/raw/`, with fair thermal state, followed by an interrupted descending-order run in `recheck-20261121/diagnostic/`. The interrupted report contains 426 samples and records the overheat interruption. These are diagnostic records and are excluded from the current three-run dataset.
- Current measured-build source SHA-256: `d29d3d56bdc51a77965eb09afadf5a9afed4e3e810bad32a4ed340728927bd63`.

These experiments measure local generation and encoding, not energy, submission-to-acceptance latency, or degraded-connectivity behavior. The existing on-chain figures use the earlier Solidity verifier. These rechecks show variation between runs but do not isolate its causes.
