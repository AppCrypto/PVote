// Package rbpvss implements the RB-PVSS construction specified in the paper.
//
// The exported algorithms intentionally mirror the paper's syntax:
// Setup, Share, DVerify, Decrypt, PVerify, and Recon. An Instance is the
// complete atomic dealer transcript (V, C, Phi, U, Pi); callers must never
// verify or aggregate its PVSS and range-proof components separately.
package rbpvss

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"fmt"
	"math/big"

	bn256 "github.com/ethereum/go-ethereum/crypto/bn256/cloudflare"
)

var (
	// ErrInvalidParameters indicates that Setup inputs do not satisfy the
	// syntax in the paper.
	ErrInvalidParameters = errors.New("invalid RB-PVSS parameters")
	// ErrInvalidTranscript indicates that a dealer transcript fails DVerify.
	ErrInvalidTranscript = errors.New("invalid RB-PVSS dealer transcript")
	// ErrNoDiscreteLog is returned when Recon cannot decode a value in S'.
	ErrNoDiscreteLog = errors.New("no bounded discrete logarithm in decoding interval")
)

// Interval is an inclusive integer interval [Min, Max]. It represents S in
// Setup and S' in Recon. Values are encoded in Z_q by Iota.
type Interval struct {
	Min int64
	Max int64
}

// PVSSParameters is pp_1=(g_0,h_0,n,t,l).
type PVSSParameters struct {
	G0 *bn256.G1
	H0 *bn256.G1
	N  int
	T  int
	L  int
}

// RangeParameters is pp_2=(g_0,h_0,g_1,pk_I,{sigma_k}_{k in S}).
type RangeParameters struct {
	G0    *bn256.G1
	H0    *bn256.G1
	G1    *bn256.G2
	PKI   *bn256.G2
	Sigma []*bn256.G1
}

// PublicParameters is pp=(pp_1,pp_2,{pk_i},N,D,t,S). D is stored in the
// paper's order {-l+1,...,0}; N is stored as {1,...,n}.
type PublicParameters struct {
	PP1 PVSSParameters
	PP2 RangeParameters
	PK  []*bn256.G1

	N         []int
	D         []int64
	Threshold int
	Range     Interval
}

// DealerProof is phi_j=(chi_j,{A^phi_ji,B^phi_ji,z^phi_ji}_{i in N}).
// The A, B, and Z field names are scoped by this type and are distinct from
// the elements of a DecryptionProof. Chi is one common Fiat-Shamir challenge
// for every encrypted shareholder evaluation.
type DealerProof struct {
	Chi *big.Int
	A   []*bn256.G1
	B   []*bn256.G1
	Z   []*big.Int
}

// RangeProof is pi_jd=(E_jd,F_jd,U'_jd,C'_jd,chi_jd,z_jd,1,
// z_jd,beta,z_jd,3).
type RangeProof struct {
	E      *bn256.G1
	F      *bn256.GT
	UPrime *bn256.G1
	CPrime *bn256.G1
	Chi    *big.Int
	Z1     *big.Int
	ZBeta  *big.Int
	Z3     *big.Int

	// fBase is the base-group first message used by the EVM verifier, with
	// F=e(fBase,g_1) for the fixed generator g_1.
	fBase *bn256.G1
}

// FEncoding returns the EVM representation of F used by the paper's
// implementation accounting and by the Solidity verifier.
func (p *RangeProof) FEncoding() *bn256.G1 {
	if p == nil {
		return nil
	}
	return cloneG1(p.fBase)
}

// Instance is I_j=({v_ji}_{i in D union N},{C_ji}_{i in N},phi_j,
// {U_jd,pi_jd}_{d in D}). V is ordered as D followed by N.
type Instance struct {
	V   []*bn256.G1
	C   []*bn256.G1
	Phi DealerProof
	U   []*bn256.G1
	Pi  []*RangeProof
}

// DecryptionProof is
// varphi_ji=(A^varphi_ji,B^varphi_ji,chi_ji,z^varphi_ji).
type DecryptionProof struct {
	A   *bn256.G1
	B   *bn256.G1
	Chi *big.Int
	Z   *big.Int
}

// IndexedShare is one pair (i, sh_i) supplied to Recon.
type IndexedShare struct {
	Index int
	Share *bn256.G1
}

// Aggregate is the online aggregation state ({C~_i},{U*_d}); Count is the
// accepted-ballot count m=|A|.
type Aggregate struct {
	C     []*bn256.G1
	U     []*bn256.G1
	Count int
}

// Setup implements RB-PVSS.Setup(1^lambda,n,t,l,S). The range-issuer key is
// deliberately not returned: the paper requires it to be erased after the
// public range-signature table is generated and checked.
func Setup(lambda, n, t, l int, scoreDomain Interval) (*PublicParameters, []*big.Int, error) {
	if lambda <= 0 || n <= 0 || t <= 0 || t > n || l <= 0 || l >= t || scoreDomain.Min > scoreDomain.Max {
		return nil, nil, ErrInvalidParameters
	}

	g0Scalar, err := randomNonZeroScalar()
	if err != nil {
		return nil, nil, err
	}
	h0Scalar, err := randomNonZeroScalar()
	if err != nil {
		return nil, nil, err
	}
	g1Scalar, err := randomNonZeroScalar()
	if err != nil {
		return nil, nil, err
	}
	g0 := new(bn256.G1).ScalarBaseMult(g0Scalar)
	h0 := new(bn256.G1).ScalarBaseMult(h0Scalar)
	g1 := new(bn256.G2).ScalarBaseMult(g1Scalar)

	sks := make([]*big.Int, n)
	pks := make([]*bn256.G1, n)
	for i := range sks {
		sks[i], err = randomNonZeroScalar()
		if err != nil {
			return nil, nil, err
		}
		pks[i] = new(bn256.G1).ScalarMult(g0, sks[i])
	}

	rangeSize, err := intervalSize(scoreDomain)
	if err != nil {
		return nil, nil, fmt.Errorf("%w: %v", ErrInvalidParameters, err)
	}
	rangeSK, err := sampleRangeIssuerKey(scoreDomain)
	if err != nil {
		return nil, nil, err
	}
	rangePK := new(bn256.G2).ScalarMult(g1, rangeSK)
	signatures := make([]*bn256.G1, rangeSize)
	for offset := int64(0); offset < int64(rangeSize); offset++ {
		w := scoreDomain.Min + offset
		denominator := addMod(rangeSK, Iota(w))
		inverse := new(big.Int).ModInverse(denominator, bn256.Order)
		if inverse == nil {
			return nil, nil, fmt.Errorf("%w: non-invertible range-signature denominator", ErrInvalidParameters)
		}
		signatures[offset] = new(bn256.G1).ScalarMult(g0, inverse)
	}
	// Best-effort in-process erasure. Setup never exposes this key; production
	// deployments should additionally isolate setup in a short-lived process.
	rangeSK.SetInt64(0)

	nSet := make([]int, n)
	for i := range nSet {
		nSet[i] = i + 1
	}
	dSet := make([]int64, l)
	for k := range dSet {
		dSet[k] = int64(k - l + 1)
	}

	pp := &PublicParameters{
		PP1:       PVSSParameters{G0: g0, H0: h0, N: n, T: t, L: l},
		PP2:       RangeParameters{G0: g0, H0: h0, G1: g1, PKI: rangePK, Sigma: signatures},
		PK:        pks,
		N:         nSet,
		D:         dSet,
		Threshold: t,
		Range:     scoreDomain,
	}
	if err := ValidatePublicParameters(pp); err != nil {
		return nil, nil, fmt.Errorf("generated public parameters failed validation: %w", err)
	}
	return pp, sks, nil
}

// ValidatePublicParameters performs the auditable one-time setup check,
// including every Boneh-Boyen-style range signature. The issuer key is safe to
// erase only after this function succeeds.
func ValidatePublicParameters(pp *PublicParameters) error {
	if err := validatePublicParameters(pp); err != nil {
		return err
	}
	expected := bn256.Pair(pp.PP1.G0, pp.PP2.G1)
	for offset, signature := range pp.PP2.Sigma {
		w := pp.Range.Min + int64(offset)
		denominatorPoint := new(bn256.G2).Add(
			pp.PP2.PKI,
			new(bn256.G2).ScalarMult(pp.PP2.G1, Iota(w)),
		)
		if !equalGT(bn256.Pair(signature, denominatorPoint), expected) {
			return fmt.Errorf("%w: invalid range signature for %d", ErrInvalidParameters, w)
		}
	}
	return nil
}

// Share implements RB-PVSS.Share(pp,w_j). It samples exactly one polynomial
// p_j and uses p_j(d) both in v_jd and as the mask in U_jd.
func Share(pp *PublicParameters, scores []int64) (*Instance, error) {
	if err := validatePublicParameters(pp); err != nil {
		return nil, err
	}
	if len(scores) != len(pp.D) {
		return nil, fmt.Errorf("%w: got %d scores, want %d", ErrInvalidParameters, len(scores), len(pp.D))
	}
	for _, w := range scores {
		if w < pp.Range.Min || w > pp.Range.Max {
			return nil, fmt.Errorf("score %d is outside [%d,%d]", w, pp.Range.Min, pp.Range.Max)
		}
	}

	polynomial := make([]*big.Int, pp.Threshold)
	for k := range polynomial {
		var err error
		polynomial[k], err = randomScalar()
		if err != nil {
			return nil, err
		}
	}

	l := len(pp.D)
	n := len(pp.N)
	instance := &Instance{
		V:  make([]*bn256.G1, l+n),
		C:  make([]*bn256.G1, n),
		U:  make([]*bn256.G1, l),
		Pi: make([]*RangeProof, l),
		Phi: DealerProof{
			A: make([]*bn256.G1, n),
			B: make([]*bn256.G1, n),
			Z: make([]*big.Int, n),
		},
	}
	packedEvaluations := make([]*big.Int, l)
	for k, d := range pp.D {
		packedEvaluations[k] = evaluatePolynomial(polynomial, scalarFromInt64(d))
		instance.V[k] = new(bn256.G1).ScalarMult(pp.PP1.H0, packedEvaluations[k])
	}
	shareEvaluations := make([]*big.Int, n)
	randomizers := make([]*big.Int, n)
	for k, i := range pp.N {
		shareEvaluations[k] = evaluatePolynomial(polynomial, big.NewInt(int64(i)))
		instance.V[l+k] = new(bn256.G1).ScalarMult(pp.PP1.H0, shareEvaluations[k])
		instance.C[k] = new(bn256.G1).ScalarMult(pp.PK[k], shareEvaluations[k])
		var err error
		randomizers[k], err = randomScalar()
		if err != nil {
			return nil, err
		}
		instance.Phi.A[k] = new(bn256.G1).ScalarMult(pp.PP1.H0, randomizers[k])
		instance.Phi.B[k] = new(bn256.G1).ScalarMult(pp.PK[k], randomizers[k])
	}
	instance.Phi.Chi = hashDealerTranscript(instance, l)
	for k := range pp.N {
		instance.Phi.Z[k] = subMod(randomizers[k], mulMod(instance.Phi.Chi, shareEvaluations[k]))
	}

	for k, w := range scores {
		s := packedEvaluations[k]
		instance.U[k] = addG1(
			new(bn256.G1).ScalarMult(pp.PP1.H0, Iota(w)),
			new(bn256.G1).ScalarMult(pp.PP1.G0, s),
		)
		proof, err := proveRangeBinding(pp, k, s, w, instance.U[k], instance.V[k])
		if err != nil {
			return nil, err
		}
		instance.Pi[k] = proof
	}
	return instance, nil
}

// DVerify implements RB-PVSS.DVerify(pp,I_j). It verifies the underlying
// PVSS transcript and every range-binding proof as one atomic object.
func DVerify(pp *PublicParameters, instance *Instance) bool {
	if validatePublicParameters(pp) != nil || !validInstanceShape(pp, instance) {
		return false
	}
	l := len(pp.D)
	for k := range pp.N {
		expectedA := addG1(
			new(bn256.G1).ScalarMult(pp.PP1.H0, instance.Phi.Z[k]),
			new(bn256.G1).ScalarMult(instance.V[l+k], instance.Phi.Chi),
		)
		expectedB := addG1(
			new(bn256.G1).ScalarMult(pp.PK[k], instance.Phi.Z[k]),
			new(bn256.G1).ScalarMult(instance.C[k], instance.Phi.Chi),
		)
		if !equalG1(expectedA, instance.Phi.A[k]) || !equalG1(expectedB, instance.Phi.B[k]) {
			return false
		}
	}
	if !equalScalar(instance.Phi.Chi, hashDealerTranscript(instance, l)) {
		return false
	}
	if !verifyRandomDualRS(pp, instance.V) {
		return false
	}
	for k := range pp.D {
		if !verifyRangeBinding(pp, instance.Pi[k], instance.U[k], instance.V[k]) {
			return false
		}
	}
	return true
}

// Decrypt implements RB-PVSS.Decrypt(pp,C_ji,sk_i).
func Decrypt(pp *PublicParameters, ciphertext *bn256.G1, secretKey *big.Int) (*bn256.G1, *DecryptionProof, error) {
	if validatePublicParameters(pp) != nil || ciphertext == nil || secretKey == nil || secretKey.Sign() == 0 {
		return nil, nil, ErrInvalidParameters
	}
	inverse := new(big.Int).ModInverse(normalizeScalar(secretKey), bn256.Order)
	if inverse == nil {
		return nil, nil, ErrInvalidParameters
	}
	share := new(bn256.G1).ScalarMult(ciphertext, inverse)
	r, err := randomScalar()
	if err != nil {
		return nil, nil, err
	}
	// The caller's key index is intentionally not inferred here. As in the
	// paper, the proof statement itself contains pk_i; identify it by sk_i.
	pk := new(bn256.G1).ScalarMult(pp.PP1.G0, normalizeScalar(secretKey))
	proof := &DecryptionProof{
		A: new(bn256.G1).ScalarMult(pp.PP1.G0, r),
		B: new(bn256.G1).ScalarMult(share, r),
	}
	proof.Chi = hashG1(pk, ciphertext, proof.A, proof.B)
	proof.Z = subMod(r, mulMod(proof.Chi, secretKey))
	return share, proof, nil
}

// PVerify implements RB-PVSS.PVerify(pp,i,sh_ji,C_ji,varphi_ji). Index is
// one-based, exactly as shareholder positions N={1,...,n} in the paper.
func PVerify(pp *PublicParameters, index int, share, ciphertext *bn256.G1, proof *DecryptionProof) bool {
	if validatePublicParameters(pp) != nil || index < 1 || index > len(pp.N) || share == nil || ciphertext == nil || !validDecryptionProof(proof) {
		return false
	}
	pk := pp.PK[index-1]
	expectedA := addG1(
		new(bn256.G1).ScalarMult(pp.PP1.G0, proof.Z),
		new(bn256.G1).ScalarMult(pk, proof.Chi),
	)
	expectedB := addG1(
		new(bn256.G1).ScalarMult(share, proof.Z),
		new(bn256.G1).ScalarMult(ciphertext, proof.Chi),
	)
	return equalG1(proof.A, expectedA) &&
		equalG1(proof.B, expectedB) &&
		equalScalar(proof.Chi, hashG1(pk, ciphertext, proof.A, proof.B))
}

// Recon implements RB-PVSS.Recon(pp,{(i,sh_i)},{U_d},S'). It reconstructs
// every mask at d in D, removes it from U_d, and performs bounded DLog.
func Recon(pp *PublicParameters, shares []IndexedShare, maskedScores []*bn256.G1, decoding Interval) ([]int64, error) {
	if err := validatePublicParameters(pp); err != nil {
		return nil, err
	}
	if len(shares) < pp.Threshold || len(maskedScores) != len(pp.D) || decoding.Min > decoding.Max {
		return nil, ErrInvalidParameters
	}
	if _, err := intervalSize(decoding); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidParameters, err)
	}
	seen := make(map[int]struct{}, len(shares))
	indices := make([]*big.Int, len(shares))
	for k, item := range shares {
		if item.Index < 1 || item.Index > len(pp.N) || item.Share == nil {
			return nil, ErrInvalidParameters
		}
		if _, exists := seen[item.Index]; exists {
			return nil, fmt.Errorf("%w: duplicate shareholder index %d", ErrInvalidParameters, item.Index)
		}
		seen[item.Index] = struct{}{}
		indices[k] = big.NewInt(int64(item.Index))
	}

	result := make([]int64, len(pp.D))
	for k, d := range pp.D {
		coefficients, ok := lagrangeCoefficients(scalarFromInt64(d), indices)
		if !ok {
			return nil, ErrInvalidParameters
		}
		mask := zeroG1()
		for r, coefficient := range coefficients {
			mask = addG1(mask, new(bn256.G1).ScalarMult(shares[r].Share, coefficient))
		}
		unmasked := addG1(maskedScores[k], new(bn256.G1).Neg(mask))
		value, ok := boundedDLog(pp.PP1.H0, unmasked, decoding)
		if !ok {
			return nil, fmt.Errorf("%w at packed position %d", ErrNoDiscreteLog, d)
		}
		result[k] = value
	}
	return result, nil
}

// NewAggregate returns the identity aggregation state for pp.
func NewAggregate(pp *PublicParameters) (*Aggregate, error) {
	if err := validatePublicParameters(pp); err != nil {
		return nil, err
	}
	aggregate := &Aggregate{C: make([]*bn256.G1, len(pp.N)), U: make([]*bn256.G1, len(pp.D))}
	for i := range aggregate.C {
		aggregate.C[i] = zeroG1()
	}
	for i := range aggregate.U {
		aggregate.U[i] = zeroG1()
	}
	return aggregate, nil
}

// AggregateVerified atomically verifies and accumulates one complete dealer
// transcript. Neither C nor U is changed if DVerify rejects it.
func AggregateVerified(pp *PublicParameters, aggregate *Aggregate, instance *Instance) error {
	if err := validatePublicParameters(pp); err != nil {
		return err
	}
	if aggregate == nil || len(aggregate.C) != len(pp.N) || len(aggregate.U) != len(pp.D) {
		return ErrInvalidParameters
	}
	if !DVerify(pp, instance) {
		return ErrInvalidTranscript
	}
	nextC := make([]*bn256.G1, len(aggregate.C))
	nextU := make([]*bn256.G1, len(aggregate.U))
	for i := range nextC {
		nextC[i] = addG1(aggregate.C[i], instance.C[i])
	}
	for i := range nextU {
		nextU[i] = addG1(aggregate.U[i], instance.U[i])
	}
	aggregate.C = nextC
	aggregate.U = nextU
	aggregate.Count++
	return nil
}

func proveRangeBinding(pp *PublicParameters, coordinate int, s *big.Int, w int64, u, v *bn256.G1) (*RangeProof, error) {
	beta, err := randomNonZeroScalar()
	if err != nil {
		return nil, err
	}
	wPrime, err := randomScalar()
	if err != nil {
		return nil, err
	}
	betaPrime, err := randomScalar()
	if err != nil {
		return nil, err
	}
	sPrime, err := randomScalar()
	if err != nil {
		return nil, err
	}
	signatureIndex := int(w - pp.Range.Min)
	if signatureIndex < 0 || signatureIndex >= len(pp.PP2.Sigma) {
		return nil, ErrInvalidParameters
	}
	e := new(bn256.G1).ScalarMult(pp.PP2.Sigma[signatureIndex], beta)
	fLeft := new(bn256.G1).Neg(new(bn256.G1).ScalarMult(e, wPrime))
	fRight := new(bn256.G1).ScalarMult(pp.PP1.G0, betaPrime)
	fBase := addG1(fLeft, fRight)
	f := bn256.Pair(fBase, pp.PP2.G1)
	uPrime := addG1(new(bn256.G1).ScalarMult(pp.PP1.H0, wPrime), new(bn256.G1).ScalarMult(pp.PP1.G0, sPrime))
	cPrime := new(bn256.G1).ScalarMult(pp.PP1.H0, sPrime)
	chi := hashRangeTranscript(e, u, v, fBase, uPrime, cPrime)
	return &RangeProof{
		E:      e,
		F:      f,
		UPrime: uPrime,
		CPrime: cPrime,
		Chi:    chi,
		Z1:     subMod(wPrime, mulMod(chi, Iota(w))),
		ZBeta:  subMod(betaPrime, mulMod(chi, beta)),
		Z3:     subMod(sPrime, mulMod(chi, s)),
		fBase:  fBase,
	}, nil
}

func verifyRangeBinding(pp *PublicParameters, proof *RangeProof, u, v *bn256.G1) bool {
	if !validRangeProof(proof) || u == nil || v == nil {
		return false
	}
	encodedF := bn256.Pair(proof.fBase, pp.PP2.G1)
	if !equalGT(proof.F, encodedF) {
		return false
	}
	if !equalScalar(proof.Chi, hashRangeTranscript(proof.E, u, v, proof.fBase, proof.UPrime, proof.CPrime)) {
		return false
	}
	expectedCPrime := addG1(
		new(bn256.G1).ScalarMult(v, proof.Chi),
		new(bn256.G1).ScalarMult(pp.PP1.H0, proof.Z3),
	)
	if !equalG1(proof.CPrime, expectedCPrime) {
		return false
	}
	expectedUPrime := addG1(
		new(bn256.G1).ScalarMult(u, proof.Chi),
		new(bn256.G1).ScalarMult(pp.PP1.G0, proof.Z3),
		new(bn256.G1).ScalarMult(pp.PP1.H0, proof.Z1),
	)
	if !equalG1(proof.UPrime, expectedUPrime) {
		return false
	}
	right := new(bn256.GT).ScalarMult(bn256.Pair(proof.E, pp.PP2.PKI), proof.Chi)
	right = new(bn256.GT).Add(right, new(bn256.GT).ScalarMult(bn256.Pair(proof.E, pp.PP2.G1), negMod(proof.Z1)))
	right = new(bn256.GT).Add(right, new(bn256.GT).ScalarMult(bn256.Pair(pp.PP1.G0, pp.PP2.G1), proof.ZBeta))
	return equalGT(proof.F, right)
}

func verifyRandomDualRS(pp *PublicParameters, commitments []*bn256.G1) bool {
	length := len(pp.D) + len(pp.N)
	dualPolynomialLength := length - pp.Threshold
	if dualPolynomialLength <= 0 || len(commitments) != length {
		return false
	}
	q := make([]*big.Int, dualPolynomialLength)
	for {
		nonZero := false
		for i := range q {
			var err error
			q[i], err = randomScalar()
			if err != nil {
				return false
			}
			nonZero = nonZero || q[i].Sign() != 0
		}
		if nonZero {
			break
		}
	}
	points := make([]*big.Int, 0, length)
	for _, d := range pp.D {
		points = append(points, scalarFromInt64(d))
	}
	for _, i := range pp.N {
		points = append(points, big.NewInt(int64(i)))
	}
	sum := zeroG1()
	for i, x := range points {
		denominator := big.NewInt(1)
		for j, other := range points {
			if i != j {
				denominator = mulMod(denominator, subMod(x, other))
			}
		}
		inverse := new(big.Int).ModInverse(denominator, bn256.Order)
		if inverse == nil {
			return false
		}
		y := mulMod(evaluatePolynomial(q, x), inverse)
		sum = addG1(sum, new(bn256.G1).ScalarMult(commitments[i], y))
	}
	return equalG1(sum, zeroG1())
}

func hashDealerTranscript(instance *Instance, packedCount int) *big.Int {
	h := sha256.New()
	for i := range instance.C {
		_, _ = h.Write(instance.V[packedCount+i].Marshal())
		_, _ = h.Write(instance.C[i].Marshal())
		_, _ = h.Write(instance.Phi.A[i].Marshal())
		_, _ = h.Write(instance.Phi.B[i].Marshal())
	}
	return hashToScalar(h.Sum(nil))
}

func hashRangeTranscript(e, u, v, fBase, uPrime, cPrime *bn256.G1) *big.Int {
	h := sha256.New()
	_, _ = h.Write(e.Marshal())
	_, _ = h.Write(u.Marshal())
	_, _ = h.Write(v.Marshal())
	_, _ = h.Write(fBase.Marshal())
	_, _ = h.Write(uPrime.Marshal())
	_, _ = h.Write(cPrime.Marshal())
	return hashToScalar(h.Sum(nil))
}

func hashG1(points ...*bn256.G1) *big.Int {
	h := sha256.New()
	for _, point := range points {
		_, _ = h.Write(point.Marshal())
	}
	return hashToScalar(h.Sum(nil))
}

func hashToScalar(data []byte) *big.Int {
	result := new(big.Int).SetBytes(data)
	return result.Mod(result, bn256.Order)
}

// Iota implements the paper's integer-to-field encoding iota(w)=w mod q.
func Iota(value int64) *big.Int {
	return scalarFromInt64(value)
}

func evaluatePolynomial(coefficients []*big.Int, x *big.Int) *big.Int {
	result := big.NewInt(0)
	for i := len(coefficients) - 1; i >= 0; i-- {
		result = addMod(mulMod(result, x), coefficients[i])
	}
	return result
}

func lagrangeCoefficients(target *big.Int, indices []*big.Int) ([]*big.Int, bool) {
	coefficients := make([]*big.Int, len(indices))
	for i := range indices {
		numerator := big.NewInt(1)
		denominator := big.NewInt(1)
		for j := range indices {
			if i == j {
				continue
			}
			numerator = mulMod(numerator, subMod(target, indices[j]))
			denominator = mulMod(denominator, subMod(indices[i], indices[j]))
		}
		inverse := new(big.Int).ModInverse(denominator, bn256.Order)
		if inverse == nil {
			return nil, false
		}
		coefficients[i] = mulMod(numerator, inverse)
	}
	return coefficients, true
}

func boundedDLog(base, target *bn256.G1, interval Interval) (int64, bool) {
	size, err := intervalSize(interval)
	if err != nil {
		return 0, false
	}
	shifted := addG1(target, new(bn256.G1).Neg(new(bn256.G1).ScalarMult(base, Iota(interval.Min))))
	if size <= 128 {
		current := zeroG1()
		for offset := 0; offset < size; offset++ {
			if equalG1(current, shifted) {
				return interval.Min + int64(offset), true
			}
			current = addG1(current, base)
		}
		return 0, false
	}
	// Baby-step/giant-step keeps public bounded decoding sublinear in the
	// aggregate range while preserving the paper's Recon interface.
	m := 1
	for m*m < size {
		m++
	}
	table := make(map[string]int, m)
	current := zeroG1()
	for j := 0; j < m; j++ {
		table[string(current.Marshal())] = j
		current = addG1(current, base)
	}
	step := new(bn256.G1).Neg(new(bn256.G1).ScalarMult(base, big.NewInt(int64(m))))
	giant := shifted
	for i := 0; i <= size/m; i++ {
		if j, ok := table[string(giant.Marshal())]; ok {
			offset := i*m + j
			if offset < size {
				return interval.Min + int64(offset), true
			}
		}
		giant = addG1(giant, step)
	}
	return 0, false
}

func validatePublicParameters(pp *PublicParameters) error {
	if pp == nil || pp.PP1.G0 == nil || pp.PP1.H0 == nil || pp.PP2.G0 == nil || pp.PP2.H0 == nil || pp.PP2.G1 == nil || pp.PP2.PKI == nil ||
		len(pp.N) == 0 || len(pp.D) == 0 || pp.Threshold <= len(pp.D) || pp.Threshold > len(pp.N) ||
		len(pp.PK) != len(pp.N) || pp.Range.Min > pp.Range.Max {
		return ErrInvalidParameters
	}
	if equalG1(pp.PP1.G0, zeroG1()) || equalG1(pp.PP1.H0, zeroG1()) ||
		equalG2(pp.PP2.G1, zeroG2()) || equalG2(pp.PP2.PKI, zeroG2()) {
		return ErrInvalidParameters
	}
	if pp.PP1.N != len(pp.N) || pp.PP1.T != pp.Threshold || pp.PP1.L != len(pp.D) ||
		!equalG1(pp.PP1.G0, pp.PP2.G0) || !equalG1(pp.PP1.H0, pp.PP2.H0) {
		return ErrInvalidParameters
	}
	rangeSize, err := intervalSize(pp.Range)
	if err != nil || rangeSize != len(pp.PP2.Sigma) {
		return ErrInvalidParameters
	}
	for i, value := range pp.N {
		if value != i+1 || pp.PK[i] == nil || equalG1(pp.PK[i], zeroG1()) {
			return ErrInvalidParameters
		}
	}
	for i, value := range pp.D {
		if value != int64(i-len(pp.D)+1) {
			return ErrInvalidParameters
		}
	}
	for _, signature := range pp.PP2.Sigma {
		if signature == nil {
			return ErrInvalidParameters
		}
	}
	return nil
}

func validInstanceShape(pp *PublicParameters, instance *Instance) bool {
	if instance == nil || len(instance.V) != len(pp.D)+len(pp.N) || len(instance.C) != len(pp.N) ||
		len(instance.U) != len(pp.D) || len(instance.Pi) != len(pp.D) ||
		instance.Phi.Chi == nil || len(instance.Phi.A) != len(pp.N) || len(instance.Phi.B) != len(pp.N) || len(instance.Phi.Z) != len(pp.N) {
		return false
	}
	for _, point := range instance.V {
		if point == nil {
			return false
		}
	}
	for i := range instance.C {
		if instance.C[i] == nil || instance.Phi.A[i] == nil || instance.Phi.B[i] == nil || instance.Phi.Z[i] == nil {
			return false
		}
	}
	for i := range instance.U {
		if instance.U[i] == nil || !validRangeProof(instance.Pi[i]) {
			return false
		}
	}
	return true
}

func validRangeProof(proof *RangeProof) bool {
	return proof != nil && proof.E != nil && proof.F != nil && proof.UPrime != nil && proof.CPrime != nil &&
		proof.Chi != nil && proof.Z1 != nil && proof.ZBeta != nil && proof.Z3 != nil && proof.fBase != nil &&
		!equalG1(proof.E, zeroG1())
}

func validDecryptionProof(proof *DecryptionProof) bool {
	return proof != nil && proof.A != nil && proof.B != nil && proof.Chi != nil && proof.Z != nil
}

func sampleRangeIssuerKey(interval Interval) (*big.Int, error) {
	for {
		candidate, err := randomNonZeroScalar()
		if err != nil {
			return nil, err
		}
		valid := true
		for value := interval.Min; ; value++ {
			if addMod(candidate, Iota(value)).Sign() == 0 {
				valid = false
				break
			}
			if value == interval.Max {
				break
			}
		}
		if valid {
			return candidate, nil
		}
	}
}

func intervalSize(interval Interval) (int, error) {
	if interval.Min > interval.Max {
		return 0, errors.New("empty interval")
	}
	width := new(big.Int).Sub(big.NewInt(interval.Max), big.NewInt(interval.Min))
	width.Add(width, big.NewInt(1))
	if width.Cmp(bn256.Order) >= 0 {
		return 0, errors.New("interval must contain fewer than q values")
	}
	if !width.IsInt64() || width.Int64() > int64(^uint(0)>>1) {
		return 0, errors.New("interval is too large for this implementation")
	}
	return int(width.Int64()), nil
}

func randomScalar() (*big.Int, error) {
	return rand.Int(rand.Reader, bn256.Order)
}

func randomNonZeroScalar() (*big.Int, error) {
	for {
		value, err := randomScalar()
		if err != nil {
			return nil, err
		}
		if value.Sign() != 0 {
			return value, nil
		}
	}
}

func scalarFromInt64(value int64) *big.Int {
	result := big.NewInt(value)
	result.Mod(result, bn256.Order)
	if result.Sign() < 0 {
		result.Add(result, bn256.Order)
	}
	return result
}

func normalizeScalar(value *big.Int) *big.Int {
	result := new(big.Int).Mod(new(big.Int).Set(value), bn256.Order)
	if result.Sign() < 0 {
		result.Add(result, bn256.Order)
	}
	return result
}

func addMod(values ...*big.Int) *big.Int {
	result := big.NewInt(0)
	for _, value := range values {
		result.Add(result, value)
		result.Mod(result, bn256.Order)
	}
	return result
}

func subMod(left, right *big.Int) *big.Int {
	result := new(big.Int).Sub(left, right)
	result.Mod(result, bn256.Order)
	if result.Sign() < 0 {
		result.Add(result, bn256.Order)
	}
	return result
}

func mulMod(left, right *big.Int) *big.Int {
	result := new(big.Int).Mul(left, right)
	return result.Mod(result, bn256.Order)
}

func negMod(value *big.Int) *big.Int {
	if normalizeScalar(value).Sign() == 0 {
		return big.NewInt(0)
	}
	return new(big.Int).Sub(bn256.Order, normalizeScalar(value))
}

func addG1(points ...*bn256.G1) *bn256.G1 {
	result := zeroG1()
	for _, point := range points {
		result = new(bn256.G1).Add(result, point)
	}
	return result
}

func zeroG1() *bn256.G1 {
	return new(bn256.G1).ScalarBaseMult(big.NewInt(0))
}

func zeroG2() *bn256.G2 {
	return new(bn256.G2).ScalarBaseMult(big.NewInt(0))
}

func cloneG1(point *bn256.G1) *bn256.G1 {
	if point == nil {
		return nil
	}
	clone := new(bn256.G1)
	if _, err := clone.Unmarshal(point.Marshal()); err != nil {
		return nil
	}
	return clone
}

func equalScalar(left, right *big.Int) bool {
	return left != nil && right != nil && normalizeScalar(left).Cmp(normalizeScalar(right)) == 0
}

func equalG1(left, right *bn256.G1) bool {
	return left != nil && right != nil && bytes.Equal(left.Marshal(), right.Marshal())
}

func equalG2(left, right *bn256.G2) bool {
	return left != nil && right != nil && bytes.Equal(left.Marshal(), right.Marshal())
}

func equalGT(left, right *bn256.GT) bool {
	return left != nil && right != nil && bytes.Equal(left.Marshal(), right.Marshal())
}
