import copy
import importlib.util
import json
from pathlib import Path
import tempfile
import unittest

spec = importlib.util.spec_from_file_location('summary', Path(__file__).parents[1]/'summarize.py')
module = importlib.util.module_from_spec(spec)
spec.loader.exec_module(module)


class SummaryTest(unittest.TestCase):
    def report(self):
        environment = {'thermalState': 'nominal', 'lowPowerMode': False}
        sample = dict(protocol='context-v1', parameterFingerprint='fixture', n=4, l=2, t=3,
                      encodeMs=1, nativeCallSampledPeakBytes=1000, transcriptBytes=2208, calldataBytes=3000,
                      environmentBefore=environment, environmentAfter=environment)
        return dict(schemaVersion=1, runId='synthetic-unit-test', status='complete', runsPerConfiguration=3,
                    fixtureIds=['fixture'], metadata=dict(isSimulator=False, buildConfiguration='Release',
                    hardwareIdentifier='synthetic', iOSVersion='test', sourceSHA256='test-only'),
                    samples=[dict(sample, shareMs=v, transcriptHash=str(i)) for i,v in enumerate([1,2,3])])

    def analyze(self, report):
        with tempfile.TemporaryDirectory() as directory:
            path = Path(directory)/'unit-test.json'
            path.write_text(json.dumps(report))
            return module.analyze([path])

    def test_statistics_and_guardrails(self):
        base = self.report()
        _, rows = self.analyze(base)
        self.assertEqual(rows[0]['shareMeanMs'], 2)
        self.assertEqual(rows[0]['shareSDMs'], 1)
        for field, value in [('isSimulator', True), ('buildConfiguration', 'Debug')]:
            bad = copy.deepcopy(base); bad['metadata'][field] = value
            with self.assertRaises(ValueError): self.analyze(bad)
        bad = copy.deepcopy(base); bad['status'] = 'interrupted'
        with self.assertRaises(ValueError): self.analyze(bad)
        bad = copy.deepcopy(base); bad['samples'][1]['transcriptHash'] = '0'
        with self.assertRaises(ValueError): self.analyze(bad)


if __name__ == '__main__': unittest.main()
