#!/usr/bin/env python3
"""Validate and summarize the native-EVM gas audit; never changes paper files."""
import argparse
import collections
import csv
import json
import math
import pathlib
import statistics


def load(directory):
    metadata = json.loads((directory / "metadata.json").read_text())
    if not metadata["completed"]:
        raise ValueError(f"incomplete experiment: {directory}")
    records = [json.loads(line) for line in (directory / "breakdown.jsonl").read_text().splitlines()]
    for record in records:
        c, p = record["components_gas"], record["precompiles"]
        total = sum(v for k, v in c.items() if k != "refund_subtracted") - c["refund_subtracted"]
        assert total == record["receipt_gas_used"] == record["eip_opcode_reconstructed_gas"]
        assert record["difference"] == 0
        n, l = record["n"], record["l"]
        assert p["ecmul"]["calls"] == 5 * n + 9 * l
        assert p["ecadd"]["calls"] == 4 * n + 5 * l
        assert p["pairing"]["calls"] == l and p["pairing"]["pairing_terms"] == 4 * l
        assert p["modexp"]["calls"] == n + l
        assert c["ecmul"] == 6000 * p["ecmul"]["calls"]
        assert c["ecadd"] == 150 * p["ecadd"]["calls"]
        assert c["pairing"] == 45000 * l + 34000 * 4 * l
        assert c["modexp"] == 1349 * (n + l)
        assert p["sha256"]["input_bytes"] % 32 == 0
        assert c["sha256"] == 60 * p["sha256"]["calls"] + 12 * p["sha256"]["input_bytes"] // 32
        assert c["calldata"] == 4 * record["calldata_zero_bytes"] + 16 * record["calldata_nonzero_bytes"]
        words = math.ceil(record["peak_root_memory_bytes"] / 32)
        assert c["memory_expansion"] == 3 * words + words * words // 512
        expected_store = 20000 * (2 * (n + l) + 2) if record["accepted_count"] == 1 else 20000 + 2900 * (2 * (n + l) + 1)
        assert c["sstore"] == expected_store
        assert c["sload"] == record["eip2929_predicted_sload"]
        assert c["sstore"] == record["eip2200_2929_predicted_sstore"]
    return metadata, records


def write_summary(directory, records):
    grouped = collections.defaultdict(list)
    for r in records:
        grouped[(r["n"], r["l"], r["t"], r["accepted_count"])].append(r)
    rows = []
    for (n, l, t, m), sample in sorted(grouped.items()):
        row = dict(n=n, l=l, t=t, accepted_count=m, repetitions=len(sample),
                   total_mean=statistics.mean(r["receipt_gas_used"] for r in sample),
                   total_std=statistics.stdev(r["receipt_gas_used"] for r in sample) if len(sample) > 1 else 0,
                   peak_memory_bytes=sample[0]["peak_root_memory_bytes"],
                   maximum_model_difference=max(abs(r["difference"]) for r in sample))
        for component in sample[0]["components_gas"]:
            row[component + "_mean"] = statistics.mean(r["components_gas"][component] for r in sample)
        for name, p in sample[0]["precompiles"].items():
            row[name + "_calls"] = p["calls"]
        rows.append(row)
    with (directory / "summary.csv").open("w", newline="") as f:
        w = csv.DictWriter(f, fieldnames=list(rows[0])); w.writeheader(); w.writerows(rows)
    return rows


def plot(directory, rows):
    import matplotlib
    matplotlib.use("Agg")
    import matplotlib.pyplot as plt
    labels = [f"({r['n']},{r['l']},{r['t']})  m={r['accepted_count']}" for r in rows]
    groups = [
        ("Pairing", lambda r: r["pairing_mean"], "#365c8d"),
        ("Scalar multiplication", lambda r: r["ecmul_mean"], "#21918c"),
        ("Storage (SLOAD+SSTORE)", lambda r: r["sload_mean"] + r["sstore_mean"], "#5ec962"),
        ("Memory expansion", lambda r: r["memory_expansion_mean"], "#f9c74f"),
        ("Other EVM execution", lambda r: r["other_execution_mean"], "#f9844a"),
        ("Remaining costs", lambda r: sum(r[k + "_mean"] for k in ["transaction_base", "calldata", "ecadd", "modexp", "sha256"]), "#bcbddc"),
    ]
    fig, axes = plt.subplots(1, 2, figsize=(12, 4.8), sharey=True)
    for axis, normalized in zip(axes, [False, True]):
        left = [0.0] * len(rows)
        for name, fn, color in groups:
            values = [fn(r) / (r["total_mean"] / 100 if normalized else 1e6) for r in rows]
            axis.barh(range(len(rows)), values, left=left, color=color, label=name, height=.7)
            left = [a + b for a, b in zip(left, values)]
        axis.set_xlabel("Share of receipt gas (%)" if normalized else "Receipt gas (million)")
        axis.set_yticks(range(len(rows)), labels)
        axis.grid(axis="x", alpha=.2); axis.set_axisbelow(True)
        axis.set_title("Relative composition" if normalized else "First acceptance versus steady state")
    axes[0].invert_yaxis()
    handles, legend = axes[0].get_legend_handles_labels()
    fig.legend(handles, legend, loc="lower center", ncol=3, frameon=False, fontsize=9)
    fig.suptitle("RB-PVSS submitRB gas audit: 3 independent deployments per configuration", fontsize=12)
    fig.tight_layout(rect=(0, .13, 1, .94))
    fig.savefig(directory / "gas_breakdown.pdf")
    fig.savefig(directory / "gas_breakdown.png", dpi=180)
    plt.close(fig)


def write_report(directory, metadata, records, rows):
    lines = ["RB-PVSS GAS BREAKDOWN: SUPPLEMENTARY EXPERIMENT", "",
             "Scope and method", "----------------",
             "The manuscript has not been changed. This audit executes the checked-in Solidity",
             "benchmark bytecode on a native go-ethereum v1.14.6 in-memory chain using London",
             "gas rules. Signed transactions pass through ApplyTransaction and produce receipts",
             "and real persistent state changes. Fake PoW supplies local blocks; no live network",
             "or mainnet deployment feasibility is claimed. The bytecode uses the legacy proof",
             "format, so these measurements are not security tests of the revised context binding.",
             "Each deployment uses fresh public parameters and independent all-zero score ballots",
             "in S=[0,10], with a distinct registered account per ballot. Three deployments are",
             "measured per (n,l,t); m=1 and m=2 separate initialization from steady-state storage.",
             "The block gas limit is 90,071,992,547; each transaction allows 180,000,000 gas.", "",
             "EIP and opcode accounting", "-------------------------",
             "ECMUL: 6,000 gas; ECADD: 150 gas; pairing check: 45,000+34,000*k (EIP-1108).",
             "Calldata: 4 gas/zero byte and 16 gas/nonzero byte (EIP-2028); base tx: 21,000.",
             "ModExp: EIP-2565, here 32-byte operands and exponent q-2, hence 1,349 per call.",
             "SLOAD: cold 2,100 / warm 100; SSTORE uses EIP-2200 plus EIP-2929 and EIP-3529.",
             "Memory: 3*w+floor(w*w/512), where w is the root frame peak word count.",
             "Other execution is summed from opcode charges and net STATICCALL overhead.",
             "Forwarded child gas is subtracted at STATICCALL and precompile consumption is",
             "added separately, avoiding double counting. Memory expansion is separated from",
             "opcode totals. Raw per-opcode STATICCALL entries omit call gas, which is recorded",
             "separately as net staticcall_overhead_including_call_memory. EIP precompile and",
             "storage predictions are checked independently",
             "of receipt totals. Refunds are subtracted after applying the EIP-3529 cap.", "",
             f"Validation: all {len(records)} traced submissions match receipt gas exactly (0 gas error).",
             "The complete decomposition is an EIP/opcode accounting reconstruction, not a",
             "claim that the paper's abstract EC operation count alone predicts total gas.", "",
             "Results (mean +/- sample standard deviation, three independent deployments)",
             "-------------------------------------------------------------------------"]
    for r in rows:
        total = r["total_mean"]
        lines.append(f"n={r['n']} l={r['l']} t={r['t']} m={r['accepted_count']}: {total:,.0f} +/- {r['total_std']:.2f} gas")
        for name in ["pairing", "ecmul", "sload", "sstore", "memory_expansion", "other_execution", "calldata", "ecadd", "modexp", "sha256", "transaction_base"]:
            value = r[name + "_mean"]
            lines.append(f"  {name:19s} {value:12,.0f} ({value/total*100:6.2f}%)")
    lines += ["", "Consistency with the operation analysis", "---------------------------------------",
              "Traces confirm exactly 5*n+9*l ECMUL calls, l four-term pairing checks,",
              "4*n+5*l ECADD calls (including proof-internal additions), and n+l inversions.",
              "Thus scalar multiplication and pairing counts agree with the analytic verifier",
              "counts. Hashing, finite-field loops, ABI manipulation, memory and storage must",
              "also be charged to explain full receipt gas.", "",
              "First acceptance writes 2*(n+l)+2 zero slots to nonzero. Later acceptances write",
              "one fresh accepted-index slot and update 2*(n+l)+1 existing nonzero slots.",
              "The predicted first-versus-later SSTORE difference is [2*(n+l)+1]*17,100 gas:",
              "530,100; 2,069,100; and 4,086,900 gas for the three configurations.",
              "This matches the measured storage difference exactly. Small differences in",
              "total gas additionally follow from zero/nonzero calldata byte counts.", "",
              "Optimization implications (not measured speedups)", "------------------------------------------------",
              "The largest steady-state case spends about 70.6% on memory expansion plus other",
              "EVM execution; pairing alone is about 10.0%. Source-level candidates include",
              "preallocating transcript buffers instead of repeated abi.encodePacked copying,",
              "and precomputing public dual-RS interpolation weights instead of nested per-ballot",
              "denominator loops and repeated inversions. These require separate implementation",
              "and validation before claiming any gas reduction.", "",
              "Primary references", "------------------",
              "https://eips.ethereum.org/EIPS/eip-1108",
              "https://eips.ethereum.org/EIPS/eip-2028",
              "https://eips.ethereum.org/EIPS/eip-2200",
              "https://eips.ethereum.org/EIPS/eip-2929",
              "https://eips.ethereum.org/EIPS/eip-3529",
              "https://eips.ethereum.org/EIPS/eip-2565",
              "https://github.com/ethereum/go-ethereum/blob/v1.14.6/params/protocol_params.go", ""]
    cross_client = directory / "cross_client.csv"
    if cross_client.exists():
        cross = list(csv.DictReader(cross_client.open()))
        assert all(int(r["difference"]) == 0 for r in cross)
        lines += ["Cross-client execution check", "----------------------------",
                  f"The same {len(cross)} signed Ganache fixture transactions (one deployment and two",
                  "(100,19,59) submissions) were re-executed on native Geth. Receipt gas matches",
                  "exactly for every transaction. The source fixture used Ganache v7.9.2 / Merge;",
                  "the Geth audit uses London rules. Both price the exercised operations identically.", ""]
    scaling_validation = directory.parent / "voter-scaling-1000-native" / "receipt_reconciliation.csv"
    if scaling_validation.exists():
        validation = list(csv.DictReader(scaling_validation.open()))
        assert all(int(r["difference"]) == 0 for r in validation)
        lines += ["Independent 1000-voter validation", "--------------------------------",
                  f"The component model predicts all {len(validation)} receipts in the separate scaling",
                  "experiment exactly, using each transaction's measured calldata byte counts",
                  "and the EIP-derived first-write storage premium. All residuals are zero.", ""]
    (directory / "report.txt").write_text("\n".join(lines))
    largest = next(r for r in rows if r["n"] == 100 and r["l"] == 19 and r["accepted_count"] == 2)
    memory_other = (largest["memory_expansion_mean"] + largest["other_execution_mean"]) / largest["total_mean"] * 100
    reply = r"""\AuthorResponse{9}{\color{blue}We thank the reviewer for requesting a more informative gas analysis. We conducted a supplementary experiment at \((n,l,t)\in\{(10,5,7),(50,10,30),(100,19,59)\}\), with three independent deployments per configuration and separate first and second acceptances (18 transactions). The breakdown reports pairing, scalar multiplication, other precompiles, calldata, storage, memory expansion, and remaining EVM execution. EIP-based prices and traced opcode charges reconstruct all 18 receipt totals with zero gas discrepancy. In particular, the observed \(5n+9l\) scalar multiplications and \(l\) four-term pairing checks agree with the operation analysis. For \((100,19,59)\), memory expansion and other EVM execution account for %.1f\%% of steady-state gas, identifying transcript-buffer allocation and dual-RS arithmetic as optimization candidates. The experiment uses an in-memory Geth EVM chain with London gas rules and a permissive local gas limit; no mainnet feasibility or measured optimization speedup is claimed. Raw transactions, receipts, accounting checks, and plots accompany the supplementary results.}
""" % memory_other
    (directory / "response_draft.tex").write_text(reply)


def main():
    parser = argparse.ArgumentParser(); parser.add_argument("directory", type=pathlib.Path); parser.add_argument("--plot", action="store_true")
    args = parser.parse_args(); metadata, records = load(args.directory)
    rows = write_summary(args.directory, records); write_report(args.directory, metadata, records, rows)
    if args.plot: plot(args.directory, rows)
    print(f"Validated {len(records)} traces; summary and report written to {args.directory}")


if __name__ == "__main__":
    main()
