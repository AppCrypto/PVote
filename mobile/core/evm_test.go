package mobilecore

import (
	"encoding/hex"
	"math/big"
	"strings"
	"testing"

	encoding "PVote/chain/rbpvss"
	contract "PVote/compile/contract"
	rb "PVote/crypto/RBPVSS"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus/ethash"
	gethcore "github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/params"
)

// Real EVM execution without a network service: iPhone payload -> locally signed
// transaction -> acceptance event, then rejection of a repeated voter identity.
func TestMobileBallotAcceptedByContract(t *testing.T) {
	key, _ := crypto.GenerateKey()
	address := crypto.PubkeyToAddress(key.PublicKey)
	config := *params.AllEthashProtocolChanges
	config.ChainID = big.NewInt(1337)
	genesis := &gethcore.Genesis{Config: &config, GasLimit: 150000000, BaseFee: big.NewInt(params.InitialBaseFee),
		Alloc: types.GenesisAlloc{address: {Balance: new(big.Int).Exp(big.NewInt(10), big.NewInt(21), nil)}}}
	consensus := ethash.NewFaker()
	defer consensus.Close()
	db, blocks, _ := gethcore.GenerateChainWithGenesis(genesis, consensus, 1, nil)
	defer db.Close()
	parent := blocks[0]
	submit := func(tx *types.Transaction) *types.Receipt {
		t.Helper()
		next, receipts := gethcore.GenerateChain(&config, parent, consensus, db, 1, func(_ int, b *gethcore.BlockGen) { b.AddTx(tx) })
		parent = next[0]
		return receipts[0][0]
	}
	pp, _, err := rb.Setup(128, 4, 3, 2, rb.Interval{Min: 0, Max: 10})
	if err != nil {
		t.Fatal(err)
	}
	g0, h0, g1, pki, sigma, pks, threshold, l, a, b, err := encoding.ConstructorArguments(pp)
	if err != nil {
		t.Fatal(err)
	}
	abi, _ := contract.ContractMetaData.GetAbi()
	args, err := abi.Pack("", g0, h0, g1, pki, sigma, pks, threshold, l, a, b, []common.Address{address})
	if err != nil {
		t.Fatal(err)
	}
	data := append(common.FromHex(contract.ContractMetaData.Bin), args...)
	deployment, err := types.SignTx(types.NewContractCreation(0, new(big.Int), 100000000, big.NewInt(2000000000), data), types.NewEIP155Signer(config.ChainID), key)
	if err != nil {
		t.Fatal(err)
	}
	if receipt := submit(deployment); receipt.Status != 1 {
		t.Fatal("deployment reverted")
	}
	deployed := crypto.CreateAddress(address, 0)
	var engine Engine
	if _, err = engine.Prepare(Export(pp)); err != nil {
		t.Fatal(err)
	}
	m, err := engine.Generate(Request{Protocol: Legacy, Scores: []int64{2, 4}, Verify: true, IncludeCalldata: true})
	if err != nil {
		t.Fatal(err)
	}
	for nonce := uint64(1); nonce <= 2; nonce++ {
		signed, err := Sign(hex.EncodeToString(crypto.FromECDSA(key)), SignRequest{ChainID: "1337", Nonce: nonce, Gas: 100000000, GasPrice: "2000000000", To: deployed.Hex(), Data: m.Calldata})
		if err != nil {
			t.Fatal(err)
		}
		raw, _ := hex.DecodeString(strings.TrimPrefix(signed["rawTransaction"], "0x"))
		var vote types.Transaction
		if err = vote.UnmarshalBinary(raw); err != nil {
			t.Fatal(err)
		}
		receipt := submit(&vote)
		if nonce == 1 {
			if receipt.Status != 1 || len(receipt.Logs) != 1 || receipt.Logs[0].Topics[0].Hex() != m.AcceptanceTopic || receipt.Logs[0].Topics[1].Big().Int64() != 1 {
				t.Fatal("native ballot was not accepted atomically")
			}
		} else if receipt.Status != 0 || len(receipt.Logs) != 0 {
			t.Fatal("repeat vote was not rejected")
		}
	}
}
