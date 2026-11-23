# PVote iPhone app

Native SwiftUI research client for iOS 17+. The arm64 Go library generates
ballots on the phone. Wallet keys stay in Keychain, and transactions are signed
locally. Public fixtures contain no issuer or tallier secret keys.

## Build

Open `mobile/ios/PVote.xcodeproj`, select a signing team and connected iPhone,
and build Release. The build phase compiles Go automatically. Set `PVOTE_GO`
only if Go is outside the usual installation paths. For measurements, stop the
debugger and launch the installed app from the phone.

Unsigned build check:

```sh
xcodebuild -project mobile/ios/PVote.xcodeproj -scheme PVote \
  -configuration Release -sdk iphoneos -destination 'generic/platform=iOS' \
  -derivedDataPath /private/tmp/pvote-ios-derived CODE_SIGNING_ALLOWED=NO build
```

## Measure

In **Benchmark**, select the three paper configurations or the full 23-point
curve. Parameters use `(n,l,t)` order. Default runs use 3 warmups and 30 samples
per configuration. Keep the app in the foreground, disable Low Power Mode,
and wait for nominal temperature. Interrupted samples are saved and labeled.

Choose `context-v1` for context-bound proofs or `legacy-solidity-v1` for the
existing contract format. Reports keep the two formats separate.

| Field | Measurement |
| --- | --- |
| `shareMs` | Share wall time, including random sampling |
| `shareCpuMs` | Process user + system CPU time in the same window |
| `encodeMs` / `verifyMs` | Separate encoding / optional verification time |
| `transcriptBytes` / `calldataBytes` | Fixed-width transcript / ABI call size |
| `shareAllocatedBytes` / `shareGCs` | Go allocations / GC count during Share |
| `nativeCallSampledPeakBytes` | Whole-app footprint sampled every 2 ms during the native call |

Setup, loading, warmups, network requests, and report writes are excluded from
Share timing. The first warmup is verified. Go uses `GOMAXPROCS=1` with default
GC. Reports record device, OS, build, source hash, parameters, and thermal and
battery metadata. Sampled footprint is a lower bound, not exact Share memory.
Battery level does not measure energy.

Export JSON/CSV from **Results** and summarize complete physical-device Release
runs. Simulator, interrupted, hot, and Low Power Mode runs are rejected.

```sh
python3 mobile/summarize.py /path/to/benchmark-*.json --out mobile/results
```

## HTTP export

Connect Mac and iPhone to the same network. Allow PVote's local-network access
in iOS Settings. To export an existing report without rerunning measurements:

```sh
python3 -u mobile/receive_http.py --host LAN_IP --port 8766 \
  --out experiments/iphone/raw --expected-sha256 REPORT_SHA256
```

Enable the Xcode launch argument `--export-http=http://LAN_IP:8766/upload`.
Both endpoints verify the digest. Disable the argument after export.

For three fresh curve runs, use a random token file containing at least 16
characters:

```sh
python3 -u mobile/receive_http.py --host LAN_IP --port 8766 \
  --out experiments/iphone/recheck-new/raw --token-file /path/to/token --reports 3
```

Launch with `--benchmark-recheck=http://LAN_IP:8766/upload#TOKEN` in Release
without the debugger. Runs use ascending, descending, and shifted orders,
3 warmups and 30 samples per point, and at least 60 s cooling between runs.
The app waits for nominal temperature and stops after a 10-minute cooling
timeout. Each report is exported separately, giving 2,070 measurements.
Remove the launch argument and token afterward. `--allow-interrupted` permits
diagnostic import only, not inclusion in complete-run statistics.

For the three-scheme generation comparison, use
`--benchmark-comparison=http://LAN_IP:8766/upload#TOKEN`. Three rounds use
3 warmups and 10 samples per point. Summarize with
`mobile/summarize_comparison.py`. Setup, key derivation, encoding, and verification
are excluded. The 720 iPhone samples are in
`experiments/iphone/comparison-20261121/`. Mac comparison data are separate in
`experiments/mac-generation-20261121/` and can be regenerated with
`scripts/run_score_generation.sh`.

## Submit a test ballot

The local-chain path uses the existing `legacy-solidity-v1` verifier.
Context-bound proofs are available offline but are not sent to this bytecode.
The permissive test-chain gas limits do not establish mainnet feasibility.

1. Copy the phone's account address from **Voting**.
2. Run `./mobile/prepare_election.sh PHONE_ADDRESS LAN_IP` on the Mac. It starts
   a Ganache London development chain, funds the phone, registers its address,
   and deploys a `(4,2,3)` election. Keep the terminal running.
3. Transfer `mobile/election.json` through a trusted channel and import it.
4. Generate the ballot, then sign and submit it. Acceptance requires a successful
   receipt and the expected acceptance event. Reports are saved as `vote-*.json`.

The client retains a pending transaction hash across RPC errors and restarts,
then resumes polling without creating a replacement. It does not implement a
complete dropped-transaction or reorganization policy. Receipt polling is every
second, so the reported observation delay includes polling. Each test deployment
accepts the registered account once. Restart the preparation script for another
election. These are public development keys and a research prototype.

## Checks

```sh
GOCACHE=/private/tmp/pvote-go-cache go test ./mobile/... ./cmd/pvote-mobile-fixtures ./cmd/pvote-mobile-election
./mobile/test_bridge.sh
python3 -m unittest discover -s mobile/tests -p 'test_*.py' -v
```

`core/` implements encoding, measurements, and signing. `bridge/` exports the
Go C ABI. `ios/` handles the UI, Keychain, RPC, sampling, and exports. Host and
simulator checks are not iPhone measurements. Regenerate public fixtures only
when needed with `go run ./cmd/pvote-mobile-fixtures` and synchronize their
fingerprints before comparing devices.
