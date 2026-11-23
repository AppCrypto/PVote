# Native Mac transcript generation: retained reference

This dataset initially restored the transcript-generation panel of Figure 12 and is retained as a reference. Figure 12(a) now uses the physical iPhone data in `experiments/iphone/comparison-20261121/`. All three Mac generation paths here were executed on the same Apple M4 Mac (Mac16,12), 10 logical CPUs, 16 GiB RAM, macOS 26.5.2.

- Native Go 1.25.4, optimized build, `GOMAXPROCS=1`, default garbage collection enabled.
- Fixed `(n,l,t)=(10,5,7)`, `kappa=1,...,8`, `B=2^kappa-1`. PriScore has three registered voters. The same cyclic distribution of B over five coordinates is used in all schemes.
- Actual phases: context-bound RB-PVSS `ShareWithContext`; PriScore `Commit` plus `Vote`; and 3PVSS-bit `ShareScore` executing kappa binary instances. These are the repository reproductions, not operation-count estimates.
- Setup and key derivation execute once before measurement; parameters are cached and reused across the three runs. Each configuration has three warmups and ten measured ballots per run. The first warmup is verified. Serialization/hash calculation and verification are outside the generation timer.
- Three prescribed orders: ascending `(kappa,scheme)` with scheme order PVote/PriScore/3PVSS-bit; the reverse list; and a half-list cyclic shift. All 720 measured transcripts are distinct. Every plotted point pools 30 samples, and SD includes variation across the three runs.
- No generation samples were filtered. Mac thermal telemetry was not collected. This benchmark does not measure phone performance or energy.
- Raw JSON records are retained unchanged. Their nested sysctl calls returned unavailable in the workspace sandbox, so hardware metadata was independently verified using sysctl and stored in `hardware_metadata.json`. CSV columns `verifiedHardware`, `cpu`, and `hardwareMetadataSource` identify that supplement.

## Results at kappa=8

| Scheme | Phase | Mean ± sample SD (ms) |
|---|---|---:|
| PVote | Share | 16.196 ± 2.250 |
| PriScore | Commit+Vote | 3742.748 ± 15.452 |
| 3PVSS-bit | kappa Share | 47.510 ± 0.365 |

## Reproduce

Run a new native Mac experiment from the repository root (use a new output directory):

```sh
./scripts/run_score_generation.sh experiments/mac-generation-rerun
```

Recompute the published summaries from the retained reports:

```sh
python3 mobile/summarize_comparison.py \
  experiments/mac-generation-20261121/raw/benchmark-*.json \
  --out experiments/mac-generation-20261121 \
  --hardware-metadata experiments/mac-generation-20261121/hardware_metadata.json
```

`raw.csv` contains 720 sample measurements; `summary.csv` has 24 scheme/score-width groups. Full input checksums and the measured-source fingerprint are recorded in `analysis_manifest.json` and the raw metadata. The current Figure 12(a) embeds the separate iPhone means and SDs directly in LaTeX; these Mac generation observations remain available for comparison.

RB-PVSS supports independent coordinate ranges, whereas PriScore imposes a fixed total-score budget. The 3PVSS comparison is explicitly a bit-plane adaptation. Matching the numerical score scale does not make their ballot semantics or trust assumptions identical.
