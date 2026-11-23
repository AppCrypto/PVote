package mobilecore

import (
	"encoding/json"
	"testing"
)

func TestComparisonRealProtocolsAndFreshTranscripts(t *testing.T) {
	var engine Engine
	if _, err := engine.PrepareComparison(9); err == nil {
		t.Fatal("oversized comparison accepted")
	}
	if _, err := engine.GenerateComparison(Request{Kappa: 1, Scheme: "PVote"}); err == nil {
		t.Fatal("missing setup accepted")
	}
	for _, kappa := range []int{1, 4, 8} {
		first, err := engine.PrepareComparison(kappa)
		if err != nil {
			t.Fatal(err)
		}
		again, err := engine.PrepareComparison(kappa)
		if err != nil {
			t.Fatal(err)
		}
		a, _ := json.Marshal(first)
		b, _ := json.Marshal(again)
		if string(a) != string(b) {
			t.Fatal("parameters not reused")
		}
		for _, scheme := range []string{"PVote", "PriScore", "3PVSS-bit"} {
			request := Request{Kappa: kappa, Scheme: scheme, ElectionID: "comparison-test", VoterIndex: 1, Verify: true}
			m, err := engine.GenerateComparison(request)
			if err != nil {
				t.Fatalf("%s k=%d: %v", scheme, kappa, err)
			}
			if m["verifyMs"] == nil || m["generationMs"].(float64) <= 0 || len(m["transcriptHash"].(string)) != 64 {
				t.Fatal("invalid verified measurement")
			}
			next, err := engine.GenerateComparison(request)
			if err != nil || next["transcriptHash"] == m["transcriptHash"] {
				t.Fatal("reused transcript")
			}
			if m["parameterFingerprint"] != next["parameterFingerprint"] {
				t.Fatal("parameter changed within case")
			}
		}
	}
}
