#!/usr/bin/env python3
"""Summarize three native score-domain runs on one physical iPhone or Mac."""
import argparse
import collections
import csv
import json
import math
from pathlib import Path
import statistics

SCHEMES = {'PVote':'Share','PriScore':'Commit+Vote','3PVSS-bit':'kappa Share'}

def analyze(paths):
    if len(paths) != 3: raise ValueError('Three complete comparison reports are required')
    reports = [json.loads(Path(p).read_text()) for p in paths]
    if len({r['runId'] for r in reports}) != 3: raise ValueError('Duplicate run')
    if {r['configurationOrderPolicy'] for r in reports} != {'ascending','descending','rotated'}: raise ValueError('Missing comparison order')
    if len({r['metadata']['sourceSHA256'] for r in reports}) != 1: raise ValueError('Different measured source builds')
    if len({(r['metadata']['hardwareIdentifier'],r['metadata'].get('iOSVersion') or r['metadata'].get('osVersion'),r['metadata']['goOS']) for r in reports}) != 1: raise ValueError('Different devices/software')
    hashes=set(); raw=[]; groups=collections.defaultdict(list)
    for r in reports:
        m=r['metadata']
        if r.get('experiment')!='score_domain_generation' or r.get('status')!='complete': raise ValueError('Incomplete or different experiment')
        device = m.get('goOS')=='ios' and m.get('buildConfiguration')=='Release'
        mac = m.get('goOS')=='darwin' and m.get('buildConfiguration')=='Go optimized'
        if m.get('isSimulator') is not False or not (device or mac): raise ValueError('Native optimized device/Mac build required')
        if r.get('runsPerConfiguration')!=10 or r.get('warmupsPerConfiguration')!=3 or len(r['samples'])!=240: raise ValueError('Wrong sample count')
        expected={f'{scheme}-k{k}' for scheme in SCHEMES for k in range(1,9)}
        if set(r['fixtureIds'])!=expected or set(r['verifiedWarmupCases'])!=expected: raise ValueError('Missing verified warmup/configuration')
        local=collections.defaultdict(list)
        setups={x['kappa']:x for x in r['setupRecords']}
        for s in r['samples']:
            scheme=s['scheme'];k=s['kappa']
            if scheme not in SCHEMES or k not in range(1,9) or s['phase']!=SCHEMES[scheme]: raise ValueError('Wrong phase')
            if (s['n'],s['l'],s['t'],s['budget'],s['priVoters'])!=(10,5,7,(1<<k)-1,3): raise ValueError('Wrong comparison parameters')
            setup=setups[k]
            if sum(setup['scores'])!=s['budget'] or len(setup['scores'])!=5: raise ValueError('Invalid common score vector')
            if setup['fingerprints'][scheme]!=s['parameterFingerprint']: raise ValueError('Setup fingerprint mismatch')
            if s['transcriptHash'] in hashes: raise ValueError('Repeated transcript')
            hashes.add(s['transcriptHash'])
            if not math.isfinite(s['generationMs']) or s['generationMs']<=0 or 'verifyMs' in s: raise ValueError('Invalid generation measurement')
            if device:
                for env in [s['environmentBefore'],s['environmentAfter']]:
                    if env['thermalState'] not in ['nominal','fair'] or env['lowPowerMode']: raise ValueError('Unusable environment')
            key=(scheme,k);local[key].append(s);groups[key].append(s)
            row={k:v for k,v in s.items() if not isinstance(v,(dict,list))}
            row.update(runId=r['runId'],order=r['configurationOrderPolicy'],hardware=m['hardwareIdentifier'],osVersion=m.get('iOSVersion') or m.get('osVersion'),goOS=m['goOS'],sourceSHA256=m['sourceSHA256'])
            if device: row.update(thermalBefore=s['environmentBefore']['thermalState'],thermalAfter=s['environmentAfter']['thermalState'])
            raw.append(row)
        if len(local)!=24 or any(len(v)!=10 or {s['sampleIndex'] for s in v}!=set(range(1,11)) for v in local.values()): raise ValueError('Incomplete configuration')
    summary=[]
    metadata=reports[0]['metadata']
    for (scheme,k),samples in sorted(groups.items()):
        if len(samples)!=30 or len({s['parameterFingerprint'] for s in samples})!=1: raise ValueError('Inconsistent pooled configuration')
        values=sorted(s['generationMs'] for s in samples)
        summary.append(dict(scheme=scheme,phase=SCHEMES[scheme],kappa=k,n=10,l=5,t=7,budget=(1<<k)-1,priVoters=3,samples=30,
            hardware=metadata['hardwareIdentifier'],cpu=metadata.get('cpu'),osVersion=metadata.get('iOSVersion') or metadata.get('osVersion'),goOS=metadata['goOS'],sourceSHA256=metadata['sourceSHA256'],
            generationMeanMs=statistics.mean(values),generationSDMs=statistics.stdev(values),generationMedianMs=statistics.median(values),generationP95Ms=values[math.ceil(.95*len(values))-1],
            serializationMeanMs=statistics.mean(s['serializationMs'] for s in samples),transcriptBytes=samples[0]['transcriptBytes'],parameterFingerprint=samples[0]['parameterFingerprint'],
            sampledAppFootprintMaxBytes=max(s['nativeCallSampledPeakBytes'] for s in samples) if all('nativeCallSampledPeakBytes' in s for s in samples) else None))
    return raw,summary

def write_csv(path,rows):
    with path.open('w',newline='') as f:
        writer=csv.DictWriter(f,fieldnames=sorted(set().union(*(r.keys() for r in rows))))
        writer.writeheader();writer.writerows(rows)

def main():
    p=argparse.ArgumentParser(description=__doc__);p.add_argument('reports',nargs='+',type=Path);p.add_argument('--out',required=True,type=Path);p.add_argument('--hardware-metadata',type=Path);a=p.parse_args()
    raw,summary=analyze(a.reports);a.out.mkdir(parents=True,exist_ok=True)
    if a.hardware_metadata:
        hardware=json.loads(a.hardware_metadata.read_text())
        if any(r['goOS']!='darwin' for r in raw): raise ValueError('Mac hardware supplement requires Mac reports')
        for row in raw+summary:
            row.update(verifiedHardware=hardware['hardwareIdentifier'],cpu=hardware['cpu'],hardwareMetadataSource=str(a.hardware_metadata))
    write_csv(a.out/'raw.csv',raw);write_csv(a.out/'summary.csv',summary)
    for r in summary: print(f"{r['scheme']:10s} k={r['kappa']}: {r['generationMeanMs']:.3f} ± {r['generationSDMs']:.3f} ms, n={r['samples']}")

if __name__=='__main__': main()
