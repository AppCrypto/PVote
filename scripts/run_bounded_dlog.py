#!/usr/bin/env python3
"""Measure ten bounded discrete logs using the voting parameters."""

import argparse
import csv
import hashlib
import json
import os
from pathlib import Path
import statistics
import subprocess
import tempfile


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--out", type=Path, default=Path("experiments/bounded-dlog-20261122"))
    out = parser.parse_args().out.resolve()
    root = Path(__file__).resolve().parents[1]
    if (out / "raw.json").exists():
        raise SystemExit("Choose a new output directory.")
    out.mkdir(parents=True, exist_ok=True)
    fixture = root / "mobile/fixtures/fixtures.json"
    files = [root / "go.mod", root / "go.sum", Path(__file__).resolve(), fixture]
    files += list((root / "crypto/RBPVSS").glob("*.go"))
    digest = lambda p: hashlib.sha256(p.read_bytes()).hexdigest()
    sources = lambda: {str(p.relative_to(root)): digest(p) for p in sorted(files)}
    before = sources()
    env = dict(os.environ, GOOS="darwin", GOARCH="arm64", CGO_ENABLED="0",
               GOCACHE="/private/tmp/pvote-go-cache", PVOTE_DLOG_BENCH_OUT=str(out / "raw.json"),
               PVOTE_DLOG_FIXTURES=str(fixture))
    with tempfile.TemporaryDirectory(prefix="pvote-dlog-") as build_dir:
        binary = Path(build_dir) / "decode.test"
        subprocess.run(["go", "test", "-c", "-trimpath", "-o", str(binary), "./crypto/RBPVSS"],
                       cwd=root, env=env, check=True)
        binary_hash = digest(binary)
        subprocess.run([str(binary), "-test.run=^TestBoundedDLogMeasurement$", "-test.count=1", "-test.v"],
                       cwd=root, env=env, check=True)
    assert before == sources(), "Source changed during measurement."
    report = json.loads((out / "raw.json").read_text())
    rows = report["samples"]
    assert report["goArch"] == "arm64"
    assert len(rows) == len({r["exponent"] for r in rows}) == 10
    assert all(r["verified"] and 0 <= r["exponent"] <= 25500000 for r in rows)
    values = [r["decodeNS"] / 1e6 for r in rows]
    summary = dict(samples=10, mean_ms=statistics.mean(values), sd_ms=statistics.stdev(values),
                   min_ms=min(values), max_ms=max(values))
    for name, data in [("samples.csv", rows), ("summary.csv", [summary])]:
        with (out / name).open("w", newline="") as f:
            writer = csv.DictWriter(f, fieldnames=list(data[0]))
            writer.writeheader()
            writer.writerows(data)
    (out / "summary.json").write_text(json.dumps(summary, indent=2) + "\n")
    manifest = dict(sourceFiles=before, binarySHA256=binary_hash,
                    outputs={name: digest(out / name) for name in ["raw.json", "samples.csv", "summary.csv", "summary.json"]})
    (out / "manifest.json").write_text(json.dumps(manifest, indent=2) + "\n")
    print(f"Decode: {summary['mean_ms']:.3f} ± {summary['sd_ms']:.3f} ms, n=10")


if __name__ == "__main__":
    main()
