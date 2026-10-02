package rbpvss

import (
	"bytes"
	"context"
	crand "crypto/rand"
	"crypto/sha256"
	"math/big"
	"os"
	"strings"
	"testing"
	"time"

	contract "PVote/compile/contract"
	core "PVote/crypto/RBPVSS"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	bn256 "github.com/ethereum/go-ethereum/crypto/bn256/cloudflare"
	"github.com/ethereum/go-ethereum/ethclient"
)

const integrationGasLimit uint64 = 180_000_000

func TestGoTranscriptAndSolidityVerifierAgree(t *testing.T) {
	url := os.Getenv("PVOTE_EVM_URL")
	privateKey := strings.TrimPrefix(os.Getenv("PVOTE_EVM_PRIVATE_KEY"), "0x")
	if url == "" || privateKey == "" {
		t.Skip("set PVOTE_EVM_URL and PVOTE_EVM_PRIVATE_KEY, or run scripts/test_evm.sh")
	}
	client, err := ethclient.Dial(url)
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()
	key, err := crypto.HexToECDSA(privateKey)
	if err != nil {
		t.Fatal(err)
	}
	auth, err := bind.NewKeyedTransactorWithChainID(key, big.NewInt(1337))
	if err != nil {
		t.Fatal(err)
	}
	auth.GasLimit = integrationGasLimit

	pp, secretKeys, err := core.Setup(128, 4, 3, 1, core.Interval{Min: 0, Max: 3})
	if err != nil {
		t.Fatal(err)
	}
	instance, err := core.Share(pp, []int64{2})
	if err != nil {
		t.Fatal(err)
	}
	dealer, rangeBinding, err := Transcript(instance)
	if err != nil {
		t.Fatal(err)
	}
	verifier := deployVerifier(t, client, auth, pp)

	// Build a transcript that satisfies the A-side dealer equations and its
	// Fiat--Shamir challenge, but deliberately violates every B-side equation.
	// A verifier that merely hashes PhiB, without checking it, accepts this.
	originalReader := crand.Reader
	deterministicRandom := make([]byte, 8192)
	// Share first samples t=3 polynomial coefficients and n=4 dealer-proof
	// randomizers. Keep those zero, but make the following range-proof beta
	// nonzero as required by the protocol.
	deterministicRandom[(3+4)*32+31] = 1
	crand.Reader = bytes.NewReader(deterministicRandom)
	zeroInstance, zeroErr := core.Share(pp, []int64{0})
	crand.Reader = originalReader
	if zeroErr != nil {
		t.Fatalf("construct zero-witness transcript: %v", zeroErr)
	}
	forgedDealer, forgedRange, err := Transcript(zeroInstance)
	if err != nil {
		t.Fatal(err)
	}
	forgedDealer.C[0] = G1Point(pp.PP1.G0)
	forgedDealer.PhiChi = dealerChallenge(forgedDealer, pp.PP1.L)
	if forgedDealer.PhiChi.Sign() == 0 {
		t.Fatal("unexpected zero forged dealer challenge")
	}
	tx, err := verifier.SubmitRB(auth, forgedDealer, forgedRange)
	if err != nil {
		t.Fatalf("submit transcript with invalid PhiB relation: %v", err)
	}
	if receipt := waitReceipt(t, client, tx); receipt.Status != types.ReceiptStatusFailed {
		t.Fatalf("invalid PhiB relation status=%d, want revert", receipt.Status)
	}

	tampered := rangeBinding
	tampered.Chi = cloneInts(rangeBinding.Chi)
	tampered.Chi[0].Add(tampered.Chi[0], big.NewInt(1))
	tx, err = verifier.SubmitRB(auth, dealer, tampered)
	if err != nil {
		t.Fatalf("submit tampered transcript: %v", err)
	}
	if receipt := waitReceipt(t, client, tx); receipt.Status != types.ReceiptStatusFailed {
		t.Fatalf("tampered transcript status=%d, want revert", receipt.Status)
	}
	accepted, err := verifier.AcceptedCount(nil)
	if err != nil || accepted.Sign() != 0 {
		t.Fatalf("accepted count after revert=%v, err=%v", accepted, err)
	}

	tx, err = verifier.SubmitRB(auth, dealer, rangeBinding)
	if err != nil {
		t.Fatalf("submit honest transcript: %v", err)
	}
	if receipt := waitReceipt(t, client, tx); receipt.Status != types.ReceiptStatusSuccessful {
		t.Fatalf("honest transcript reverted; gas=%d", receipt.GasUsed)
	}
	onChainU, err := verifier.GetAggregateU(nil)
	if err != nil || len(onChainU) != 1 {
		t.Fatalf("get aggregate U: len=%d err=%v", len(onChainU), err)
	}
	if got, err := G1(onChainU[0]); err != nil || string(got.Marshal()) != string(instance.U[0].Marshal()) {
		t.Fatal("Solidity aggregate differs from the accepted Go transcript")
	}
	tx, err = verifier.SubmitRB(auth, dealer, rangeBinding)
	if err != nil {
		t.Fatalf("submit replay: %v", err)
	}
	if receipt := waitReceipt(t, client, tx); receipt.Status != types.ReceiptStatusFailed {
		t.Fatal("the same mapped voter was able to submit twice")
	}
	accepted, err = verifier.AcceptedCount(nil)
	if err != nil || accepted.Cmp(big.NewInt(1)) != 0 {
		t.Fatalf("accepted count after replay=%v, err=%v", accepted, err)
	}

	for i := 0; i < pp.Threshold; i++ {
		share, proof, err := core.Decrypt(pp, instance.C[i], secretKeys[i])
		if err != nil {
			t.Fatal(err)
		}
		tx, err = verifier.SubmitDecryptionShare(auth, big.NewInt(int64(i+1)), G1Point(share), G1Point(proof.A), G1Point(proof.B), proof.Chi, proof.Z)
		if err != nil {
			t.Fatalf("submit decryption share %d: %v", i+1, err)
		}
		if receipt := waitReceipt(t, client, tx); receipt.Status != types.ReceiptStatusSuccessful {
			t.Fatalf("decryption share %d reverted", i+1)
		}
	}
	unmasked, err := verifier.Reconstruct(&bind.CallOpts{Context: context.Background()})
	if err != nil || len(unmasked) != 1 {
		t.Fatalf("reconstruct: len=%d err=%v", len(unmasked), err)
	}
	got, err := G1(unmasked[0])
	if err != nil {
		t.Fatal(err)
	}
	want := new(bn256.G1).ScalarMult(pp.PP1.H0, big.NewInt(2))
	if string(got.Marshal()) != string(want.Marshal()) {
		t.Fatal("on-chain reconstruction does not equal h_0^2")
	}
}

func deployVerifier(t *testing.T, client *ethclient.Client, auth *bind.TransactOpts, pp *core.PublicParameters) *contract.Contract {
	t.Helper()
	g0, h0, g1, pkI, sigma, pks, threshold, coordinates, minimum, maximum, err := ConstructorArguments(pp)
	if err != nil {
		t.Fatal(err)
	}
	_, tx, verifier, err := contract.DeployContract(auth, client, g0, h0, g1, pkI, sigma, pks, threshold, coordinates, minimum, maximum, []common.Address{auth.From})
	if err != nil {
		t.Fatalf("deploy verifier: %v", err)
	}
	if receipt := waitReceipt(t, client, tx); receipt.Status != types.ReceiptStatusSuccessful {
		t.Fatalf("verifier deployment reverted; gas=%d", receipt.GasUsed)
	}
	return verifier
}

func waitReceipt(t *testing.T, client *ethclient.Client, tx *types.Transaction) *types.Receipt {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	receipt, err := bind.WaitMined(ctx, client, tx)
	if err != nil {
		t.Fatal(err)
	}
	return receipt
}

func dealerChallenge(dealer contract.RBPVSSVerifierDealerTranscript, packedCount int) *big.Int {
	h := sha256.New()
	for i := range dealer.C {
		for _, point := range []contract.RBPVSSVerifierG1Point{
			dealer.V[packedCount+i], dealer.C[i], dealer.PhiA[i], dealer.PhiB[i],
		} {
			_, _ = h.Write(point.X.FillBytes(make([]byte, 32)))
			_, _ = h.Write(point.Y.FillBytes(make([]byte, 32)))
		}
	}
	return new(big.Int).Mod(new(big.Int).SetBytes(h.Sum(nil)), bn256.Order)
}
