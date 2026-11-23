package main

import (
	"crypto/ecdsa"
	"encoding/json"
	"math/big"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
)

func developmentWallet(t *testing.T, count int) (string, []*ecdsa.PrivateKey) {
	t.Helper()
	keys := make([]*ecdsa.PrivateKey, count)
	wallet := map[string]map[string]string{"private_keys": {}}
	for i := range keys {
		key, err := crypto.ToECDSA(big.NewInt(int64(i + 1)).FillBytes(make([]byte, 32)))
		if err != nil {
			t.Fatal(err)
		}
		keys[i] = key
		wallet["private_keys"][crypto.PubkeyToAddress(key.PublicKey).Hex()] = "0x" + strings.Repeat("0", 64-len(key.D.Text(16))) + key.D.Text(16)
	}
	data, err := json.Marshal(wallet)
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(t.TempDir(), "accounts.json")
	if err := os.WriteFile(path, data, 0600); err != nil {
		t.Fatal(err)
	}
	return path, keys
}

func TestVoterScalingLoads1000DistinctRegisteredAccounts(t *testing.T) {
	path, _ := developmentWallet(t, 1000)
	keys, err := voterScalingKeys("", path, 1000)
	if err != nil || len(keys) != 1000 {
		t.Fatalf("accounts=%d err=%v", len(keys), err)
	}
	if _, err := voterScalingKeys("", path, 1001); err == nil {
		t.Fatal("insufficient account count was silently accepted")
	}
}

func TestVoterScalingRejectsDuplicateOrSubstitutedAccounts(t *testing.T) {
	path, keys := developmentWallet(t, 2)
	key := strings.Repeat("0", 64-len(keys[0].D.Text(16))) + keys[0].D.Text(16)
	if _, err := voterScalingKeys(key+","+key, "", 2); err == nil {
		t.Fatal("duplicate voters were accepted")
	}
	data := []byte(`{"private_keys":{"0x0000000000000000000000000000000000000001":"` + key + `"}}`)
	if err := os.WriteFile(path, data, 0600); err != nil {
		t.Fatal(err)
	}
	if _, err := voterScalingKeys("", path, 1); err == nil {
		t.Fatal("a key registered under the wrong address was accepted")
	}
}
