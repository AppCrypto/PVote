# Local validation — 2026-11-20

- iPhone arm64 / iOS 17 minimum / Release / signing disabled: **build passed**.
- iOS Simulator arm64 / Release / signing disabled: **build passed**. Simulator UI was not launched.
- `go test ./... -count=1`: **passed**. Existing external-RPC integration tests retain their environment-dependent skip behavior; the new mobile EVM test executes locally without a network service.
- `go vet ./mobile/... ./cmd/pvote-mobile-fixtures ./cmd/pvote-mobile-election`: **passed**.
- `mobile/test_bridge.sh`: **passed** on macOS, including six real proofs across both proof formats, the largest preset, local wallet handling, parameter provenance, RPC success/error/pending cases, receipt checks, and memory probing.
- `python3 -m unittest discover -s mobile/tests -p 'test_*.py' -v`: **passed**.
- Mobile EVM test: a generated, encoded, locally signed ballot was **accepted by the actual checked-in Solidity bytecode**; a second submission from the same identity was **reverted**.
- The 25 bundled fixtures were checked to contain only public parameter fields.

No physical iPhone was connected, no device provisioning was changed, and no iPhone runtime, latency, memory, energy, or network-acceptance results have been collected. The unsigned device build still needs the user's signing Team before installation. Tomorrow's procedure is in `README.md`.

Build logs on this Mac: `/private/tmp/pvote-ios-build-final.log`, `/private/tmp/pvote-ios-sim-build.log`, `/private/tmp/pvote-mobile-regression.log`. Build products are in `/private/tmp/pvote-ios-derived` and `/private/tmp/pvote-ios-sim-derived`.

## Physical-device run — 2026-11-21

- Installed the Release app on an iPhone 13 (`iPhone14,5`), iOS 26.6.2. Debug executable was disabled for the performance run.
- The app completed the full 23-configuration `context-v1` curve: 3 warmups and 30 measured ballots per configuration, 690 measured ballots in total.
- Report saved on the phone: `benchmark-2026-11-20T23-28-25Z-CE89002F.json`, with SHA-256 `59b94972457cea7d953239c5ceedca3936b49cbd77fec19607d412beecdd6db6`.
- After local-network access was enabled, HTTP export succeeded with status 200. The 981,322-byte raw JSON was received unchanged, its SHA-256 matched the value above, and the receiver shut down automatically.
- `mobile/summarize.py` accepted all 690 samples: physical-device Release metadata, complete run, distinct transcripts, nominal thermal state, and low-power mode off. Raw and summary CSV files are in `experiments/iphone/`.
- At `(n,l,t)=(100,19,59)`, Share mean/sample SD is **62.818/0.305 ms**, P95 is **63.603 ms**, and mean encoding time is **0.651 ms**. Across the run, the maximum sampled whole-app footprint is **37.782 MiB**; it is not an exact Share-only peak.
- The paper's iPhone plot, methods, results, and abstract now use this run. The plot includes all 23 observed means and sample SDs inline in LaTeX. The device was USB-powered and these are 30 repetitions per configuration in one session.
- These are local ballot-generation measurements. No physical-device on-chain acceptance, degraded-connectivity latency, or energy measurement is claimed.
- Data analysis tests and `git diff --check` passed. The updated manuscript compiled in two passes with no undefined references; the abstract and evaluation pages were rendered and visually checked. The PDF is `output/pdf/PVote.R1.iphone.pdf`. The original `paper/main.tex` and all reviewer quotations were preserved.
- The one-shot HTTP receiver is closed. Xcode is restored to Release without a debugger and with the export launch argument disabled.

## Figure 6 recheck — 2026-11-21

- Repeated all 23 curve configurations in three prescribed orders: ascending `(l,n)`, descending, and a half-list cyclic shift. Each configuration has 3 warmups and 30 fresh measured ballots per run: **2,070 measurements, 90 per plotted point**.
- Each run started at nominal thermal state, with at least 60 seconds between runs plus any additional cooling. All 4,140 before/after thermal readings are nominal, all samples have low-power mode off, and the phone ran on battery power. Release and no-debugger settings were retained.
- Three complete reports were exported by HTTP with a one-time receiver token and SHA-256 verification. They are in `experiments/iphone/recheck-20261121/cooled/raw/`; `analysis_manifest.json` records all input checksums, run orders, cooling waits, and the measured source fingerprint.
- Pooled n=100 Share mean/sample SD: l=1 **40.345/1.835 ms**, l=10 **53.705/1.405 ms**, l=19 **68.280/4.005 ms**. For l=19, pooled n=90 is **64.373 ms**; n=100 exceeds n=90 in each of the three runs. The l=19, n=100 P95 is **73.036 ms**, with mean encoding **0.688 ms**.
- The maximum sampled whole-app footprint across the current dataset is **47.970 MiB**. The paper's inline plot, methods, results, and abstract use this three-run dataset; `experiments/iphone/raw.csv` and `summary.csv` contain the current measurements and summary.
- Earlier CSVs are archived in `experiments/iphone/initial-20261121/`, and the original raw JSON is retained unchanged. An earlier fair-state 690-sample recheck and its overheated 426-sample interrupted successor are retained separately for diagnostics and excluded from the final cooled dataset.
- Data analysis tests passed. All 23 inline mean/SD pairs match the 90-sample summaries. The updated manuscript compiled in two passes without undefined references; the abstract, methods, and Figure 6 pages were rendered and visually checked. Prior phone icons, grammar corrections, and reviewer quotations were preserved.
- The HTTP receiver closed after the third report. Xcode's automatic recheck argument and its token were removed; Release remains configured without a debugger. The temporary token file was removed.

## Figure 12 generation comparison — 2026-11-21

- Added real RB-PVSS, PriScore Commit+Vote, and 3PVSS-bit generation interfaces using the same Go BN254 backend. Setup/key derivation and public transcript serialization are outside generation timers. Public hashes encode curve points explicitly and omit unexported openings; private setup keys are not exported.
- Go comparison tests passed at kappa 1, 4, and 8, including verification, parameter reuse, and fresh transcripts. The Swift/Go bridge passed additional comparison proofs at kappa 1 and 8. The complete unsigned iOS Release build passed (`/private/tmp/pvote-comparison-ios-build.log`).
- The iPhone was unavailable for connection. The restored Figure 12 panel therefore uses **native Mac measurements**, explicitly identified in the caption and methods. No iPhone comparison timings are claimed.
- The native Mac runner completed three orders, 24 configurations per order, 3 warmups and 10 measured ballots per configuration: **720 distinct transcripts, 30 samples per point**. Every first warmup was verified. Hardware is Apple M4 / Mac16,12 / 10 logical CPUs / 16 GiB, macOS 26.5.2, Go 1.25.4, `GOMAXPROCS=1`, default GC enabled.
- Published kappa=8 generation mean/sample SD: RB-PVSS **16.196/2.250 ms**, PriScore **3742.748/15.452 ms**, 3PVSS-bit **47.510/0.365 ms**. `experiments/mac-generation-20261121/` retains raw reports, 720 sample rows, 24 summary rows, checksums, separately verified hardware metadata, and reproduction instructions.
- Python comparison-summary tests passed. Figure 12's generation panel, caption, methods, and comparison text were restored/updated with blue highlighting; data are inline in LaTeX. The earlier verification and ballot-size panels are retained. The manuscript compiled without undefined references, and the changed pages were rendered and visually checked.
- The unused phone HTTP receiver was stopped, and its one-time token and Xcode launch argument were removed.

## Three-scheme iPhone comparison — 2026-11-21 (later run)

- The same physical iPhone 13 reconnected, and the signed Release app ran the complete three-scheme comparison without a debugger. Native Go 1.25.4, arm64, `GOMAXPROCS=1`, default GC enabled; measured source SHA-256 is `8e478489c748a03c5153e152505ffa8495d4c8346b1a8bc217c201e490fc151e`.
- Completed ascending, descending, and rotated orders, each containing 24 cases with three warmups and ten measured ballots per case: **720 distinct transcripts, 30 observations per point**. All first warmups were verified. Setup, key derivation, cooling, serialization, and verification are outside the generation timer.
- All 1,440 before/after thermal observations are nominal. Low-power mode is off, every sample is on battery power, and zero samples were excluded. Cases start at nominal state, with at least 60 seconds between runs.
- All three raw reports were received by HTTP with token authentication and SHA-256 checks. `experiments/iphone/comparison-20261121/analysis_manifest.json` records input hashes, timestamps, source provenance, and environment counts.
- At kappa=8, generation mean/sample SD is **12.230/0.579 ms** for RB-PVSS, **2793.358/69.040 ms** for PriScore, and **36.533/1.545 ms** for 3PVSS-bit. The maximum sampled whole-app footprint across native calls is **46.048 MiB**, including generation and serialization.
- The shared comparison analyzer accepted all 24 groups and all 720 measurements; its unit tests passed. `summary.csv`, `raw.csv`, `comparison_with_mac.csv`, and the dataset README retain the supplementary phone results separately from the published Mac plot.
- The HTTP receiver closed after the third report. The phone test process was stopped and the automatic launch argument removed; Release remains configured without a debugger.

## Figure 12 adoption of iPhone data — 2026-11-21

- Figure 12(a) now uses the complete three-scheme iPhone dataset: all 24 means and sample SDs are embedded directly in LaTeX, with 30 observations per point. The earlier Mac generation dataset remains available as a reference.
- The methods, generation subcaption, platform caption, and comparison paragraph were synchronized with blue highlighting. At kappa=8, RB-PVSS generation is **12.23 ± 0.58 ms** for `(n,l,t)=(10,5,7)`. Figure 6 and its abstract result of **68.28 ms** use the larger `(100,19,59)` configuration and were preserved.
- The verification and ballot-size panels, as well as all reviewer quotations, were checked against the pre-edit manuscript and are unchanged. All 24 generation coordinates match the iPhone summary CSV after rounding.
- The final saved manuscript snapshot compiled in two passes with stable references and no compilation errors. The methodology and Figure 12 pages were rendered and visually checked. The PDF is `output/pdf/PVote.R1.iphone.pdf`.
