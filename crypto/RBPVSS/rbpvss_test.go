package rbpvss

import (
	"errors"
	"math/big"
	"testing"

	bn256 "github.com/ethereum/go-ethereum/crypto/bn256/cloudflare"
)

func TestPaperInterfacesRoundTripAndAggregate(t *testing.T) {
	pp, secretKeys, err := Setup(128, 5, 4, 3, Interval{Min: -2, Max: 3})
	if err != nil {
		t.Fatalf("Setup: %v", err)
	}
	if got, want := pp.D, []int64{-2, -1, 0}; !equalInt64s(got, want) {
		t.Fatalf("D=%v, want %v", got, want)
	}
	if len(secretKeys) != 5 || len(pp.PK) != 5 {
		t.Fatalf("Setup returned wrong key count")
	}

	vectors := [][]int64{{-2, 0, 3}, {1, -1, 2}, {3, 3, -2}}
	aggregate, err := NewAggregate(pp)
	if err != nil {
		t.Fatalf("NewAggregate: %v", err)
	}
	for _, vector := range vectors {
		instance, err := Share(pp, vector)
		if err != nil {
			t.Fatalf("Share(%v): %v", vector, err)
		}
		if !DVerify(pp, instance) {
			t.Fatalf("DVerify rejected honest transcript for %v", vector)
		}
		if err := AggregateVerified(pp, aggregate, instance); err != nil {
			t.Fatalf("AggregateVerified: %v", err)
		}
	}

	chosen := []int{1, 3, 4, 5}
	shares := make([]IndexedShare, len(chosen))
	for k, index := range chosen {
		share, proof, err := Decrypt(pp, aggregate.C[index-1], secretKeys[index-1])
		if err != nil {
			t.Fatalf("Decrypt(%d): %v", index, err)
		}
		if !PVerify(pp, index, share, aggregate.C[index-1], proof) {
			t.Fatalf("PVerify rejected honest share %d", index)
		}
		shares[k] = IndexedShare{Index: index, Share: share}
	}

	got, err := Recon(pp, shares, aggregate.U, Interval{Min: -6, Max: 9})
	if err != nil {
		t.Fatalf("Recon: %v", err)
	}
	want := []int64{2, 2, 3}
	if !equalInt64s(got, want) {
		t.Fatalf("Recon=%v, want %v", got, want)
	}
}

func TestSetupMatchesPP1PP2AndRangeSignatureRelation(t *testing.T) {
	pp, secretKeys, err := Setup(128, 5, 4, 3, Interval{Min: -2, Max: 3})
	if err != nil {
		t.Fatal(err)
	}
	if pp.PP1.N != 5 || pp.PP1.T != 4 || pp.PP1.L != 3 {
		t.Fatalf("pp_1 dimensions are (%d,%d,%d)", pp.PP1.N, pp.PP1.T, pp.PP1.L)
	}
	if !equalG1(pp.PP1.G0, pp.PP2.G0) || !equalG1(pp.PP1.H0, pp.PP2.H0) {
		t.Fatal("pp_1 and pp_2 do not contain the same g_0,h_0")
	}
	if err := ValidatePublicParameters(pp); err != nil {
		t.Fatalf("ValidatePublicParameters: %v", err)
	}
	for i, sk := range secretKeys {
		if sk.Sign() == 0 {
			t.Fatalf("sk_%d is zero", i+1)
		}
		if !equalG1(pp.PK[i], new(bn256.G1).ScalarMult(pp.PP1.G0, sk)) {
			t.Fatalf("pk_%d != g_0^sk_%d", i+1, i+1)
		}
	}
	expectedPairing := bn256.Pair(pp.PP1.G0, pp.PP2.G1)
	for offset, sigma := range pp.PP2.Sigma {
		w := pp.Range.Min + int64(offset)
		denominatorPoint := new(bn256.G2).Add(
			pp.PP2.PKI,
			new(bn256.G2).ScalarMult(pp.PP2.G1, Iota(w)),
		)
		if !equalGT(bn256.Pair(sigma, denominatorPoint), expectedPairing) {
			t.Fatalf("range signature for %d fails e(sigma,pk_I*g_1^w)=e(g_0,g_1)", w)
		}
	}
	pp.PP2.Sigma[0] = new(bn256.G1).ScalarMult(pp.PP1.G0, big.NewInt(7))
	if err := ValidatePublicParameters(pp); err == nil {
		t.Fatal("accepted a corrupted public range-signature table")
	}
}

func TestSingleDealerRecon(t *testing.T) {
	pp, secretKeys, err := Setup(128, 4, 3, 2, Interval{Min: 0, Max: 5})
	if err != nil {
		t.Fatal(err)
	}
	instance, err := Share(pp, []int64{5, 1})
	if err != nil {
		t.Fatal(err)
	}
	shares := make([]IndexedShare, pp.Threshold)
	for i := 0; i < pp.Threshold; i++ {
		share, proof, err := Decrypt(pp, instance.C[i], secretKeys[i])
		if err != nil || !PVerify(pp, i+1, share, instance.C[i], proof) {
			t.Fatalf("decryption share %d failed: %v", i+1, err)
		}
		shares[i] = IndexedShare{Index: i + 1, Share: share}
	}
	got, err := Recon(pp, shares, instance.U, pp.Range)
	if err != nil {
		t.Fatal(err)
	}
	if !equalInt64s(got, []int64{5, 1}) {
		t.Fatalf("got %v", got)
	}
}

func TestDVerifyRecomputesEveryChallengeAndBindsMaskedScore(t *testing.T) {
	pp, _, err := Setup(128, 4, 3, 2, Interval{Min: 0, Max: 5})
	if err != nil {
		t.Fatal(err)
	}

	t.Run("dealer challenge", func(t *testing.T) {
		instance := mustShare(t, pp, []int64{1, 4})
		instance.Phi.Chi = addMod(instance.Phi.Chi, big.NewInt(1))
		if DVerify(pp, instance) {
			t.Fatal("accepted a transcript with a modified dealer challenge")
		}
	})

	t.Run("dealer ciphertext", func(t *testing.T) {
		instance := mustShare(t, pp, []int64{1, 4})
		instance.C[0] = new(bn256.G1).Add(instance.C[0], pp.PP1.G0)
		if DVerify(pp, instance) {
			t.Fatal("accepted a transcript with a modified dealer ciphertext")
		}
	})

	t.Run("dealer ciphertext response", func(t *testing.T) {
		instance := mustShare(t, pp, []int64{1, 4})
		instance.Phi.B[0] = new(bn256.G1).Add(instance.Phi.B[0], pp.PP1.G0)
		if DVerify(pp, instance) {
			t.Fatal("accepted a transcript with a modified PhiB response")
		}
	})

	t.Run("range challenge", func(t *testing.T) {
		instance := mustShare(t, pp, []int64{1, 4})
		instance.Pi[0].Chi = addMod(instance.Pi[0].Chi, big.NewInt(1))
		if DVerify(pp, instance) {
			t.Fatal("accepted a transcript with a modified range challenge")
		}
	})

	t.Run("identity range signature randomization", func(t *testing.T) {
		instance := mustShare(t, pp, []int64{1, 4})
		instance.Pi[0].E = zeroG1()
		if validRangeProof(instance.Pi[0]) {
			t.Fatal("treated identity E_jd as a valid range-proof element")
		}
		if DVerify(pp, instance) {
			t.Fatal("accepted a range proof with identity E_jd")
		}
	})

	t.Run("masked commitment", func(t *testing.T) {
		instance := mustShare(t, pp, []int64{1, 4})
		instance.U[0] = new(bn256.G1).Add(instance.U[0], pp.PP1.H0)
		if DVerify(pp, instance) {
			t.Fatal("accepted a range proof replayed against a modified U_jd")
		}
	})

	t.Run("packed PVSS evaluation", func(t *testing.T) {
		instance := mustShare(t, pp, []int64{1, 4})
		instance.V[0] = new(bn256.G1).Add(instance.V[0], pp.PP1.H0)
		if DVerify(pp, instance) {
			t.Fatal("accepted a proof after changing v_jd")
		}
	})
}

func TestPVerifyRecomputesChallenge(t *testing.T) {
	pp, secretKeys, err := Setup(128, 4, 3, 2, Interval{Min: 0, Max: 5})
	if err != nil {
		t.Fatal(err)
	}
	instance := mustShare(t, pp, []int64{2, 3})
	share, proof, err := Decrypt(pp, instance.C[0], secretKeys[0])
	if err != nil {
		t.Fatal(err)
	}
	proof.Chi = addMod(proof.Chi, big.NewInt(1))
	if PVerify(pp, 1, share, instance.C[0], proof) {
		t.Fatal("accepted a decryption proof with a modified challenge")
	}
}

func TestAggregateVerifiedRejectsAtomically(t *testing.T) {
	pp, _, err := Setup(128, 4, 3, 2, Interval{Min: 0, Max: 5})
	if err != nil {
		t.Fatal(err)
	}
	aggregate, err := NewAggregate(pp)
	if err != nil {
		t.Fatal(err)
	}
	instance := mustShare(t, pp, []int64{2, 3})
	instance.U[0] = new(bn256.G1).Add(instance.U[0], pp.PP1.H0)
	if err := AggregateVerified(pp, aggregate, instance); !errors.Is(err, ErrInvalidTranscript) {
		t.Fatalf("AggregateVerified returned %v", err)
	}
	if aggregate.Count != 0 {
		t.Fatalf("invalid transcript changed count to %d", aggregate.Count)
	}
	for i, point := range aggregate.C {
		if !equalG1(point, zeroG1()) {
			t.Fatalf("invalid transcript changed aggregate C[%d]", i)
		}
	}
	for i, point := range aggregate.U {
		if !equalG1(point, zeroG1()) {
			t.Fatalf("invalid transcript changed aggregate U[%d]", i)
		}
	}
}

func TestRejectsInvalidInputsAndMissingBoundedDLog(t *testing.T) {
	if _, _, err := Setup(128, 3, 2, 2, Interval{Min: 0, Max: 5}); !errors.Is(err, ErrInvalidParameters) {
		t.Fatalf("Setup with l>=t returned %v", err)
	}
	pp, secretKeys, err := Setup(128, 4, 3, 2, Interval{Min: 0, Max: 5})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := Share(pp, []int64{0, 6}); err == nil {
		t.Fatal("Share accepted an out-of-range score")
	}
	instance := mustShare(t, pp, []int64{2, 3})
	shares := make([]IndexedShare, pp.Threshold)
	for i := range shares {
		share, _, err := Decrypt(pp, instance.C[i], secretKeys[i])
		if err != nil {
			t.Fatal(err)
		}
		shares[i] = IndexedShare{Index: i + 1, Share: share}
	}
	tampered := append([]*bn256.G1(nil), instance.U...)
	tampered[0] = new(bn256.G1).Add(tampered[0], new(bn256.G1).ScalarMult(pp.PP1.H0, big.NewInt(20)))
	if _, err := Recon(pp, shares, tampered, pp.Range); !errors.Is(err, ErrNoDiscreteLog) {
		t.Fatalf("Recon returned %v, want ErrNoDiscreteLog", err)
	}
}

func mustShare(t *testing.T, pp *PublicParameters, scores []int64) *Instance {
	t.Helper()
	instance, err := Share(pp, scores)
	if err != nil {
		t.Fatal(err)
	}
	return instance
}

func equalInt64s(left, right []int64) bool {
	if len(left) != len(right) {
		return false
	}
	for i := range left {
		if left[i] != right[i] {
			return false
		}
	}
	return true
}
