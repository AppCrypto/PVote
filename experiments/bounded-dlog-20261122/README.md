# Bounded discrete-log benchmark

- Uses the voting code's BN254 G1 group and `h0` from `mobile/fixtures/fixtures.json`, fixture `n100-l19-t59`.
- Recovers ten uniform random targets in `[0, 25,500,000]`, the aggregate range for 100,000 scores in `[0,255]`. All ten recoveries are correct.
- Times the complete existing `boundedDLog` call after three warmups. Sampling, target generation, and verification are excluded.
- Apple M4, native arm64, Go 1.25.4, `GOMAXPROCS=1`, default GC. Mean ± sample SD: **276.659 ± 50.766 ms per coordinate**.
- Tallying decodes each option's aggregate once, giving `l` recoveries per election.

Run `python3 scripts/run_bounded_dlog.py --out experiments/dlog-repeat`.
`raw.json`, `samples.csv`, and `summary.json` contain the data. `manifest.json` records source and data hashes. `previous/` preserves the superseded measurement.
