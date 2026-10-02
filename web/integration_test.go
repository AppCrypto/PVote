package main

import (
	"os"
	"testing"
)

func TestGanacheBackedWebWorkflow(t *testing.T) {
	if os.Getenv("PVOTE_EVM_URL") == "" {
		t.Skip("run through scripts/test_evm.sh")
	}
	workingDirectory, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(".."); err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := os.Chdir(workingDirectory); err != nil {
			t.Errorf("restore working directory: %v", err)
		}
	}()

	state, err := newDemoState(defaultConfig())
	if err != nil {
		t.Fatal(err)
	}
	if state.Chain == nil {
		t.Fatalf("Ganache contracts were not deployed: %s", state.ChainError)
	}
	defer state.Chain.Client.Close()
	if err := state.fundInitiatorEscrow(); err != nil {
		t.Fatal(err)
	}
	if err := state.submitVote("integration voter", []int{2, 3}); err != nil {
		t.Fatal(err)
	}
	for i := 1; i <= state.Config.Threshold; i++ {
		if err := state.stakeTallier(i); err != nil {
			t.Fatal(err)
		}
		if err := state.decryptTallier(i); err != nil {
			t.Fatal(err)
		}
	}
	if err := state.finalizeTally(); err != nil {
		t.Fatal(err)
	}
	if state.Tally == nil || !state.Tally.Verified || len(state.Tally.Results) != 2 || state.Tally.Results[0] != 2 || state.Tally.Results[1] != 3 {
		t.Fatalf("unexpected final tally: %+v", state.Tally)
	}
}
