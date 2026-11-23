import copy
import importlib.util
import json
from pathlib import Path
import tempfile
import unittest

spec=importlib.util.spec_from_file_location('comparison_summary',Path(__file__).parents[1]/'summarize_comparison.py')
module=importlib.util.module_from_spec(spec);spec.loader.exec_module(module)

class ComparisonSummaryTest(unittest.TestCase):
    def reports(self):
        out=[]
        for round,order in enumerate(['ascending','descending','rotated']):
            env={'thermalState':'nominal','lowPowerMode':False}
            ids=[f'{s}-k{k}' for s in module.SCHEMES for k in range(1,9)]
            records=[dict(kappa=k,scores=[(1<<k)-1,0,0,0,0],fingerprints={s:f'{s}-{k}' for s in module.SCHEMES}) for k in range(1,9)]
            samples=[]
            for s,phase in module.SCHEMES.items():
                for k in range(1,9):
                    for i in range(1,11):
                        samples.append(dict(scheme=s,phase=phase,kappa=k,n=10,l=5,t=7,budget=(1<<k)-1,priVoters=3,sampleIndex=i,
                            parameterFingerprint=f'{s}-{k}',transcriptHash=f'{round}-{s}-{k}-{i}',generationMs=round*10+i,
                            serializationMs=1,transcriptBytes=100,nativeCallSampledPeakBytes=1000,environmentBefore=env,environmentAfter=env))
            out.append(dict(runId=f'synthetic-{round}',configurationOrderPolicy=order,experiment='score_domain_generation',status='complete',runsPerConfiguration=10,warmupsPerConfiguration=3,
                metadata=dict(sourceSHA256='test-only',hardwareIdentifier='synthetic',iOSVersion='test',isSimulator=False,buildConfiguration='Release',goOS='ios'),
                fixtureIds=ids,verifiedWarmupCases=ids,setupRecords=records,samples=samples))
        return out

    def analyze(self,reports):
        with tempfile.TemporaryDirectory() as directory:
            paths=[]
            for i,r in enumerate(reports):
                p=Path(directory)/f'{i}.json';p.write_text(json.dumps(r));paths.append(p)
            return module.analyze(paths)

    def test_pooled_statistics_and_rejections(self):
        base=self.reports();raw,rows=self.analyze(base)
        self.assertEqual(len(raw),720);self.assertEqual(len(rows),24)
        self.assertEqual(rows[0]['generationMeanMs'],15.5)
        self.assertAlmostEqual(rows[0]['generationSDMs'],77.5**.5)
        self.assertEqual(rows[0]['generationP95Ms'],29)
        for mutate in [lambda r:r[0].update(status='interrupted'),lambda r:r[0]['metadata'].update(isSimulator=True),
                       lambda r:r[1]['samples'][0].update(transcriptHash=r[0]['samples'][0]['transcriptHash']),
                       lambda r:r[0]['samples'][0]['environmentAfter'].update(thermalState='serious')]:
            bad=copy.deepcopy(base);mutate(bad)
            with self.assertRaises(ValueError):self.analyze(bad)

if __name__=='__main__':unittest.main()
