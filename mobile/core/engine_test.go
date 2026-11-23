package mobilecore

import (
	"encoding/hex"
	"encoding/json"
	"math/big"
	"strings"
	"testing"

	contract "PVote/compile/contract"
	rb "PVote/crypto/RBPVSS"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
)

func testParameters(t *testing.T) Parameters {
	t.Helper()
	pp, _, err := rb.Setup(128, 4, 3, 2, rb.Interval{Min: 0, Max: 10})
	if err != nil {
		t.Fatal(err)
	}
	return Export(pp)
}

func TestPublicParameterBoundary(t *testing.T) {
	p := testParameters(t)
	pp, err := p.Decode()
	if err != nil {
		t.Fatal(err)
	}
	if Export(pp).Fingerprint() != p.Fingerprint() {
		t.Fatal("public parameter round trip changed fingerprint")
	}
	p.Sigma[0] = p.Sigma[1]
	if _, err = p.Decode(); err == nil {
		t.Fatal("accepted corrupted range signature table")
	}
	p.N = 1000000
	if _, err = p.Decode(); err == nil {
		t.Fatal("accepted oversized public parameters")
	}
}

func TestOnDeviceGenerationAndEncoding(t *testing.T) {
	var engine Engine
	p := testParameters(t)
	if _, err := engine.Prepare(p); err != nil {
		t.Fatal(err)
	}
	var previous string
	for _, protocol := range []string{Legacy, Bound} {
		request := Request{Protocol: protocol, Scores: []int64{2, 4}, ElectionID: "chain/contract/election", VoterIndex: 1, Verify: true, IncludeCalldata: true}
		m, err := engine.Generate(request)
		if err != nil {
			t.Fatal(err)
		}
		if m.VerifyMS == nil || m.ShareMS <= 0 || m.TranscriptHash == previous {
			t.Fatal("invalid measurement or reused randomness")
		}
		if m.TranscriptBytes != (6*2+4*4)*64+(4*2+4+1)*32 {
			t.Fatal("wrong paper transcript size")
		}
		raw, err := hex.DecodeString(strings.TrimPrefix(m.Calldata, "0x"))
		if err != nil {
			t.Fatal(err)
		}
		abi, _ := contract.ContractMetaData.GetAbi()
		values, err := abi.Methods["submitRB"].Inputs.Unpack(raw[4:])
		if err != nil || len(values) != 2 {
			t.Fatalf("cannot decode complete atomic ballot: %v", err)
		}
		if len(raw) != m.CalldataBytes {
			t.Fatal("calldata size mismatch")
		}
		previous = m.TranscriptHash
	}
	if _, err := engine.Generate(Request{Protocol: Bound, Scores: []int64{2, 4}}); err == nil {
		t.Fatal("missing context accepted")
	}
	if _, err := engine.Generate(Request{Protocol: Legacy, Scores: []int64{2, 11}}); err == nil {
		t.Fatal("out-of-range score accepted")
	}
	var output map[string]any
	if err := json.Unmarshal(engine.Call([]byte(`{"action":"generate"}`)), &output); err != nil || output["error"] == nil {
		t.Fatal("bridge did not return structured error")
	}
}

func TestLocalSigning(t *testing.T) {
	wallet, err := NewWallet()
	if err != nil {
		t.Fatal(err)
	}
	r := SignRequest{ChainID: "1337", Nonce: 7, Gas: 3000000, GasPrice: "2000000000", To: "0x1111111111111111111111111111111111111111", Data: "0x12345678"}
	signed, err := Sign(wallet["privateKey"], r)
	if err != nil {
		t.Fatal(err)
	}
	raw, _ := hex.DecodeString(strings.TrimPrefix(signed["rawTransaction"], "0x"))
	var tx types.Transaction
	if err = tx.UnmarshalBinary(raw); err != nil {
		t.Fatal(err)
	}
	sender, err := types.Sender(types.NewEIP155Signer(big.NewInt(1337)), &tx)
	if err != nil {
		t.Fatal(err)
	}
	key, _ := crypto.HexToECDSA(wallet["privateKey"])
	if sender != crypto.PubkeyToAddress(key.PublicKey) || tx.Nonce() != 7 || tx.Hash().Hex() != signed["transactionHash"] {
		t.Fatal("signed transaction changed sender, nonce, or hash")
	}
	r.ChainID = "0"
	if _, err = Sign(wallet["privateKey"], r); err == nil {
		t.Fatal("invalid signing chain accepted")
	}
}
