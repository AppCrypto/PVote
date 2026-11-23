#!/usr/bin/env python3
"""Summarize exported iPhone runs; refuse simulator/Debug/incomplete runs by default."""
import argparse
import csv
import json
import math
from pathlib import Path
import statistics


def analyze(paths):
    rows = []
    for path in paths:
        report = json.loads(Path(path).read_text())
        meta = report['metadata']
        if meta.get('isSimulator') is not False or meta.get('buildConfiguration') != 'Release':
            raise ValueError(f'{path}: publication analysis requires a physical-device Release build')
        if report.get('status') != 'complete':
            raise ValueError(f'{path}: interrupted or incomplete run')
        if report.get('schemaVersion') != 1 or not meta.get('sourceSHA256'):
            raise ValueError(f'{path}: missing schema/provenance')
        samples = report['samples']
        if len(samples) != report['runsPerConfiguration'] * len(report['fixtureIds']):
            raise ValueError(f'{path}: missing measured samples')
        hashes = set()
        for sample in samples:
            if sample['transcriptHash'] in hashes:
                raise ValueError(f'{path}: repeated ballot randomness/transcript')
            hashes.add(sample['transcriptHash'])
            for env in [sample['environmentBefore'], sample['environmentAfter']]:
                if env['thermalState'] not in ('nominal', 'fair') or env['lowPowerMode']:
                    raise ValueError(f'{path}: unsuitable thermal/low-power state; repeat the run')
            row = {k: v for k, v in sample.items() if not isinstance(v, (dict, list))}
            row.update(runId=report['runId'], hardware=meta['hardwareIdentifier'],
                       iOSVersion=meta['iOSVersion'], sourceSHA256=meta['sourceSHA256'])
            row['thermalBefore'] = sample['environmentBefore']['thermalState']
            row['thermalAfter'] = sample['environmentAfter']['thermalState']
            rows.append(row)
    groups = {}
    keys = ('hardware', 'iOSVersion', 'sourceSHA256', 'protocol', 'parameterFingerprint', 'n', 'l', 't')
    for row in rows:
        groups.setdefault(tuple(row[k] for k in keys), []).append(row)
    summary = []
    for key, group in sorted(groups.items()):
        record = dict(zip(keys, key))
        values = sorted(float(row['shareMs']) for row in group)
        record.update(samples=len(values), shareMeanMs=statistics.mean(values),
                      shareSDMs=statistics.stdev(values) if len(values) > 1 else 0,
                      shareMedianMs=statistics.median(values), shareP95Ms=values[math.ceil(len(values)*0.95)-1],
                      encodeMeanMs=statistics.mean(float(row['encodeMs']) for row in group),
                      sampledAppFootprintMaxBytes=max(row['nativeCallSampledPeakBytes'] for row in group),
                      transcriptBytes=group[0]['transcriptBytes'], calldataBytes=group[0]['calldataBytes'])
        summary.append(record)
    return rows, summary


def write_csv(path, rows):
    fields = sorted(set().union(*(row.keys() for row in rows)))
    with path.open('w', newline='') as file:
        writer = csv.DictWriter(file, fields)
        writer.writeheader()
        writer.writerows(rows)


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument('reports', nargs='+', type=Path)
    parser.add_argument('--out', type=Path, default=Path('mobile/results'))
    args = parser.parse_args()
    rows, summary = analyze(args.reports)
    args.out.mkdir(parents=True, exist_ok=True)
    write_csv(args.out/'raw.csv', rows)
    write_csv(args.out/'summary.csv', summary)
    for row in summary:
        print(f"{row['hardware']} {row['protocol']} n={row['n']} l={row['l']} t={row['t']}: "
              f"Share {row['shareMeanMs']:.3f} ± {row['shareSDMs']:.3f} ms, P95 {row['shareP95Ms']:.3f} ms")
    print('Memory is sampled whole-app footprint, not an exact Share peak. No energy estimate is produced.')


if __name__ == '__main__':
    main()
