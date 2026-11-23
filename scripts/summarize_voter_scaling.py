#!/usr/bin/env python3
"""Check every stored transaction/receipt and summarize actual m=1..1000 data."""
import argparse
import collections
import csv
import gzip
import json
import pathlib
import statistics


def main():
    parser = argparse.ArgumentParser(); parser.add_argument("directory", type=pathlib.Path); parser.add_argument("--plot", action="store_true")
    parser.add_argument("--gas-audit", type=pathlib.Path, help="independent component experiment used to predict every receipt")
    args = parser.parse_args(); directory = args.directory
    metadata = json.loads((directory / "metadata.json").read_text())
    assert metadata["completed"] and metadata["runs_per_configuration"] == 1
    count = metadata["voters_per_configuration_per_run"]
    measurements = list(csv.DictReader((directory / "measurements.csv").open()))
    by_hash = {r["transaction_hash"].lower(): r for r in measurements}
    assert len(by_hash) == len(measurements)
    seen = set(); deployments = 0
    with gzip.open(directory / "transactions.jsonl.gz", "rt") as archive:
        for line in archive:
            record = json.loads(line); receipt, tx = record["receipt"], record["transaction"]
            assert int(receipt["status"], 16) == 1
            assert tx["hash"].lower() == receipt["transactionHash"].lower()
            if record["phase"] == "deployment":
                deployments += 1; continue
            key = receipt["transactionHash"].lower(); row = by_hash[key]
            assert key not in seen; seen.add(key)
            assert int(receipt["gasUsed"], 16) == int(row["gas_used"])
            m = int(row["accepted_count"])
            assert record["accepted_count"] == m
            assert len(receipt["logs"]) == 1 and int(receipt["logs"][0]["topics"][1], 16) == m
            assert receipt["logs"][0]["address"].lower() == row["contract_address"].lower()
            data = bytes.fromhex(tx["input"][2:])
            zero = data.count(0)
            assert len(data) == int(row["calldata_bytes"]) and zero == int(row["calldata_zero_bytes"])
    assert len(seen) == len(measurements)
    grouped = collections.defaultdict(list)
    for row in measurements: grouped[(int(row["n"]), int(row["l"]), int(row["t"]))].append(row)
    assert len(grouped) == deployments == len(metadata["n_l_configurations"])
    summaries = []
    for (n, l, t), rows in sorted(grouped.items()):
        assert len(rows) == count
        assert [int(r["accepted_count"]) for r in rows] == list(range(1, count + 1))
        assert len({r["voter_address"].lower() for r in rows}) == count
        values = [int(r["gas_used"]) for r in rows]
        calldata = [4 * int(r["calldata_zero_bytes"]) + 16 * int(r["calldata_nonzero_bytes"]) for r in rows]
        adjusted = [g - c for g, c in zip(values, calldata)]
        assert len(set(adjusted[1:])) == 1, "unexplained m-dependent execution gas"
        storage_difference = (2 * (n + l) + 1) * 17100
        assert adjusted[0] - adjusted[1] == storage_difference
        summaries.append(dict(n=n, l=l, t=t, accepted_transactions=count,
                              first_gas=values[0], last_gas=values[-1],
                              steady_mean=statistics.mean(values[1:]),
                              steady_std=statistics.stdev(values[1:]),
                              steady_min=min(values[1:]), steady_max=max(values[1:]),
                              steady_execution_without_calldata=adjusted[1],
                              predicted_first_write_premium=storage_difference))
    with (directory / "summary.csv").open("w", newline="") as f:
        writer = csv.DictWriter(f, fieldnames=list(summaries[0])); writer.writeheader(); writer.writerows(summaries)
    reconciliation_count = 0
    if args.gas_audit:
        baseline_metadata = json.loads((args.gas_audit / "metadata.json").read_text())
        assert baseline_metadata["completed"]
        assert baseline_metadata["deployment_bytecode_keccak256"] == metadata["deployment_bytecode_keccak256"]
        baselines = {}
        for line in (args.gas_audit / "breakdown.jsonl").read_text().splitlines():
            r = json.loads(line)
            if r["accepted_count"] != 2: continue
            key = (r["n"], r["l"], r["t"])
            value = r["receipt_gas_used"] - r["components_gas"]["calldata"]
            if key in baselines: assert baselines[key] == value
            baselines[key] = value
        with (directory / "receipt_reconciliation.csv").open("w", newline="") as f:
            w = csv.writer(f)
            w.writerow(["n", "l", "t", "accepted_count", "receipt_gas", "calldata_gas", "first_write_premium", "model_predicted_gas", "difference"])
            for r in measurements:
                n, l, t, m = (int(r[k]) for k in ("n", "l", "t", "accepted_count"))
                data_gas = 4 * int(r["calldata_zero_bytes"]) + 16 * int(r["calldata_nonzero_bytes"])
                premium = (2 * (n + l) + 1) * 17100 if m == 1 else 0
                predicted = baselines[(n, l, t)] + data_gas + premium
                delta = predicted - int(r["gas_used"])
                assert delta == 0, "independent component model does not explain this receipt"
                w.writerow([n, l, t, m, r["gas_used"], data_gas, premium, predicted, delta])
                reconciliation_count += 1
    report = ["ACTUAL CONSECUTIVE VOTER-COUNT MEASUREMENTS", "",
              f"Validated {len(measurements)} successful signed submissions and {deployments} deployments.",
              f"Each configuration executes m=1..{count}, with {count} distinct registered accounts.",
              "Every stored receipt, acceptance event, transaction hash, and calldata length was checked.",
              "After removing calldata gas, execution gas for every m>=2 is exactly constant.",
              "The first-vote premium agrees with [2(n+l)+1]*17,100 gas of initial SSTORE cost.",
              "The final contract getter was independently checked from EVM state by the runner.", ""]
    if reconciliation_count:
        report += [f"The EIP/opcode model from the separate component experiment predicts all {reconciliation_count}",
                   "receipts exactly after inserting actual calldata byte counts and the EIP-derived",
                   "first-write premium. See receipt_reconciliation.csv (all differences are zero).", ""]
    for r in summaries:
        report.append(f"({r['n']},{r['l']},{r['t']}): first={r['first_gas']:,}; m={count}: {r['last_gas']:,}; steady mean={r['steady_mean']:,.2f}, std={r['steady_std']:.2f}, range=[{r['steady_min']:,},{r['steady_max']:,}]")
    report += ["", "Backend: go-ethereum v1.14.6 native in-memory EVM chain, London gas rules,",
               "fake PoW block generation and permissive local gas limits. These are EVM",
               "executions with real transaction receipts and cumulative contract state; they",
               "are not mainnet deployments. The checked-in benchmark uses the legacy proof",
               "format, so this experiment does not measure the revised context-binding path.",
               "The manuscript and its plotted coordinates have not been changed."]
    (directory / "report.txt").write_text("\n".join(report)+"\n")
    parameter_list = ",".join(f"({r['n']},{r['l']},{r['t']})" for r in summaries)
    reply = r"""\AuthorResponse{8}{\color{blue}We thank the reviewer for identifying the discrepancy between the original sample count and plotted range. The original experiment measured ten consecutive ballots per configuration. We have now conducted an additional experiment with %s distinct registered voters at each of \((n,l,t)=%s\), recording all %s successful submission receipts and acceptance events. The final EVM-state accepted count is verified to be %s in each case. The supplementary plots show every measured point without extrapolation and separate the first-write cost from subsequent acceptances. After accounting for calldata byte composition, execution gas is identical for all measured \(m\geq2\). The additional experiment uses the native Geth EVM with London gas rules and a permissive local gas limit; raw transactions, receipts, metadata, and reproducible code are provided.}
""" % (f"{count:,}", parameter_list, f"{len(measurements):,}", f"{count:,}")
    (directory / "response_draft.tex").write_text(reply)
    if args.plot:
        import matplotlib
        matplotlib.use("Agg")
        import matplotlib.pyplot as plt
        fig, axes = plt.subplots(2, len(grouped), figsize=(12, 6), squeeze=False)
        for col, ((n, l, t), rows) in enumerate(sorted(grouped.items())):
            x = [int(r["accepted_count"]) for r in rows]; y = [int(r["gas_used"]) for r in rows]
            axes[0, col].scatter(x[1:], [v/1e6 for v in y[1:]], s=2, alpha=.55)
            axes[0, col].scatter(x[:1], [y[0]/1e6], s=20, color="#c33", label="first acceptance")
            axes[0, col].set_title(f"(n,l,t)=({n},{l},{t})")
            mean = statistics.mean(y[1:])
            axes[1, col].scatter(x[1:], [v-mean for v in y[1:]], s=2, alpha=.6)
            axes[1, col].axhline(0, lw=.7, color="#555")
            for axis in axes[:, col]: axis.set_xlabel("Accepted voters m"); axis.grid(alpha=.15)
        axes[0, 0].set_ylabel("Receipt gas (million)")
        axes[1, 0].set_ylabel("Steady gas minus mean (gas)")
        fig.suptitle("1000 actual sequential acceptances per configuration; no extrapolation")
        fig.tight_layout(rect=(0, 0, 1, .95))
        fig.savefig(directory / "voter_scaling.pdf"); fig.savefig(directory / "voter_scaling.png", dpi=180)
        plt.close(fig)
    print(f"Validated {len(measurements)} submissions; summary written to {directory}")


if __name__ == "__main__": main()
