# Three-scheme generation on iPhone 13

Figure 12(a) uses these three complete runs on the same physical iPhone 13 (`iPhone14,5`), iOS 26.6.2 (23G90), Release without a debugger. All 720 accepted measurements use fresh transcripts; each of the 24 scheme/score-width points has 30 samples.

- Native arm64 Go 1.25.4, `GOMAXPROCS=1`, default garbage collection enabled.
- Parameters `(n,l,t)=(10,5,7)`, score widths kappa=1,...,8, budget B=2^kappa-1. PriScore has three registered voters. Every scheme uses the same cyclic distribution of B over five coordinates.
- Actual phases are context-bound RB-PVSS Share, PriScore Commit+Vote, and kappa binary Share instances for 3PVSS-bit. Setup and key derivation run on the phone before timing; the same parameters are reused across runs.
- Three predefined orders: ascending (kappa,scheme), reverse, and a half-list cyclic shift. Each case has three warmups and ten measured ballots per run; the first warmup is verified. Generation, serialization, and verification have separate timers. Setup, cooling, storage, and HTTP export are outside generation timing.
- Each case starts at nominal thermal state. Runs have at least 60 seconds of cooling between them. All 1,440 before/after thermal readings are nominal, low-power mode is off throughout, and all samples are on battery power. No sample was excluded.
- Raw JSON bytes were uploaded over local HTTP with a one-time token and SHA-256 verification; input checksums, timestamps and measured-source fingerprint are in `analysis_manifest.json`.

## Generation results (mean ± sample SD, ms)

| kappa | RB-PVSS | PriScore | 3PVSS-bit |
|---:|---:|---:|---:|
| 1 | 11.756 ± 0.674 | 23.198 ± 1.211 | 4.772 ± 0.290 |
| 2 | 11.815 ± 0.679 | 45.780 ± 2.583 | 9.563 ± 0.654 |
| 3 | 11.731 ± 0.599 | 89.532 ± 4.105 | 14.265 ± 0.752 |
| 4 | 11.750 ± 0.614 | 177.156 ± 6.961 | 18.687 ± 1.030 |
| 5 | 11.218 ± 0.503 | 341.991 ± 12.012 | 22.447 ± 1.083 |
| 6 | 11.252 ± 0.423 | 684.591 ± 19.447 | 27.593 ± 0.993 |
| 7 | 11.567 ± 0.316 | 1391.314 ± 17.571 | 32.229 ± 0.508 |
| 8 | 12.230 ± 0.579 | 2793.358 ± 69.040 | 36.533 ± 1.545 |

The maximum sampled whole-app physical footprint across generation/serialization calls is **46.048 MiB**, sampled every 2 ms. This includes the shared app/runtime and cached parameters; it is not an isolated per-algorithm memory peak. Energy and network acceptance latency were not measured.

## Reproduce the summaries

```sh
python3 mobile/summarize_comparison.py \
  experiments/iphone/comparison-20261121/raw/benchmark-*.json \
  --out experiments/iphone/comparison-20261121
```

`raw.csv` contains all 720 samples; `summary.csv` contains the 24 pooled groups. `comparison_with_mac.csv` presents the retained native Mac observations alongside the phone observations. Timing ratios describe these runs and do not establish a general hardware ranking. The same implemented generation functions, numerical parameters and Go runtime policy are used; the operating systems and benchmark harnesses differ.

All 24 phone means and SDs are embedded inline in Figure 12(a). Figure 12(b) retains the earlier Mac verification measurements, and its platform is identified in the caption. The Mac generation dataset is retained separately for reference. RB-PVSS independent ranges, PriScore fixed total budget, and the explicitly labeled 3PVSS bit-plane adaptation retain their different ballot semantics.
