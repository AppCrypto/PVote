package main

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"testing"

	"github.com/ethereum/go-ethereum/params"
)

// Exercise real signed EVM transactions rather than mirroring pricing helpers.
// This catches double-counted forwarded CALL gas and incorrect storage warming.
func TestTraceReconstructsFirstAndSteadyAcceptance(t *testing.T) {
	var measurements, archive, details bytes.Buffer
	config := *params.AllEthashProtocolChanges
	if err := executeConfiguration(&config, 4, 1, 2, 1, csv.NewWriter(&measurements), json.NewEncoder(&archive), json.NewEncoder(&details)); err != nil {
		t.Fatal(err)
	}
	decoder := json.NewDecoder(&details)
	var first, steady gasAudit
	if err := decoder.Decode(&first); err != nil {
		t.Fatal(err)
	}
	if err := decoder.Decode(&steady); err != nil {
		t.Fatal(err)
	}
	for _, audit := range []gasAudit{first, steady} {
		if audit.Difference != 0 || audit.ReceiptGas != audit.PredictedGas {
			t.Fatalf("gas accounting mismatch: %+v", audit)
		}
	}
	// 2*(n+l) aggregate coordinate slots plus the existing counter switch from
	// 20,000 to 2,900 gas; the new voter acceptance flag remains zero->nonzero.
	wantDifference := uint64((2*(4+1) + 1) * 17100)
	if got := first.Components["sstore"] - steady.Components["sstore"]; got != wantDifference {
		t.Fatalf("first-write difference=%d, want %d", got, wantDifference)
	}
	firstWithoutData := first.ReceiptGas - first.Components["calldata"]
	steadyWithoutData := steady.ReceiptGas - steady.Components["calldata"]
	if firstWithoutData-steadyWithoutData != wantDifference {
		t.Fatal("unexplained gas difference after removing storage and calldata")
	}
}
