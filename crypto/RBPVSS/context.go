package rbpvss

import (
	"crypto/sha256"
	"encoding/binary"
	"hash"
	"math/big"

	bn256 "github.com/ethereum/go-ethereum/crypto/bn256/cloudflare"
)

// BallotContext is supplied by the authenticated application, never taken from
// untrusted ballot fields. ElectionID must be unique across all elections sharing
// these parameters. Applications spanning chains/contracts must namespace it.
type BallotContext struct {
	ElectionID string
	VoterIndex uint64
}

// ShareWithContext generates the paper's context-bound voting transcript.
// Share remains the standalone primitive for compatibility with existing tools.
func ShareWithContext(pp *PublicParameters, scores []int64, context BallotContext) (*Instance, error) {
	binding, err := ballotBinding(pp, context)
	if err != nil {
		return nil, err
	}
	return share(pp, scores, binding)
}

// DVerifyWithContext recomputes both challenges using the expected election and
// authenticated voter index, including the full public-parameter fingerprint.
func DVerifyWithContext(pp *PublicParameters, instance *Instance, context BallotContext) bool {
	binding, err := ballotBinding(pp, context)
	return err == nil && dVerify(pp, instance, binding)
}

// AggregateVerifiedWithContext verifies before changing any aggregate field.
// The application separately enforces eligibility and one acceptance per voter.
func AggregateVerifiedWithContext(pp *PublicParameters, aggregate *Aggregate, instance *Instance, context BallotContext) error {
	binding, err := ballotBinding(pp, context)
	if err != nil {
		return err
	}
	return aggregateVerified(pp, aggregate, instance, binding)
}

func writeField(h hash.Hash, value []byte) {
	var size [8]byte
	binary.BigEndian.PutUint64(size[:], uint64(len(value)))
	h.Write(size[:])
	h.Write(value)
}

func writeNumber(h hash.Hash, value uint64) {
	var encoded [8]byte
	binary.BigEndian.PutUint64(encoded[:], value)
	writeField(h, encoded[:])
}

func ballotBinding(pp *PublicParameters, context BallotContext) ([]byte, error) {
	if validatePublicParameters(pp) != nil || context.ElectionID == "" || context.VoterIndex == 0 {
		return nil, ErrInvalidParameters
	}
	h := sha256.New()
	writeField(h, []byte("RB-PVSS/context/v1"))
	writeField(h, []byte(context.ElectionID))
	writeNumber(h, context.VoterIndex)
	writeNumber(h, uint64(pp.Threshold))
	writeNumber(h, uint64(pp.Range.Min))
	writeNumber(h, uint64(pp.Range.Max))
	writeNumber(h, uint64(len(pp.N)))
	writeNumber(h, uint64(len(pp.D)))
	for _, point := range []*bn256.G1{pp.PP1.G0, pp.PP1.H0} {
		writeField(h, point.Marshal())
	}
	writeField(h, pp.PP2.G1.Marshal())
	writeField(h, pp.PP2.PKI.Marshal())
	for _, point := range pp.PK {
		writeField(h, point.Marshal())
	}
	writeNumber(h, uint64(len(pp.PP2.Sigma)))
	for _, point := range pp.PP2.Sigma {
		writeField(h, point.Marshal())
	}
	for _, index := range pp.N {
		writeNumber(h, uint64(index))
	}
	for _, index := range pp.D {
		writeNumber(h, uint64(index))
	}
	return h.Sum(nil), nil
}

func hashBoundDealerTranscript(instance *Instance, packedCount int, binding []byte) *big.Int {
	if len(binding) == 0 {
		return hashDealerTranscript(instance, packedCount)
	}
	h := sha256.New()
	writeField(h, []byte("RB-PVSS/dealer/v1"))
	writeField(h, binding)
	for _, point := range instance.V { // Includes all packed commitments.
		writeField(h, point.Marshal())
	}
	for i := range instance.C {
		writeField(h, instance.C[i].Marshal())
		writeField(h, instance.Phi.A[i].Marshal())
		writeField(h, instance.Phi.B[i].Marshal())
	}
	return boundHashToScalar(h.Sum(nil))
}

func hashBoundRangeTranscript(coordinate int64, binding []byte, points ...*bn256.G1) *big.Int {
	if len(binding) == 0 {
		return hashRangeTranscript(points[0], points[1], points[2], points[3], points[4], points[5])
	}
	h := sha256.New()
	writeField(h, []byte("RB-PVSS/range/base/v1"))
	writeField(h, binding)
	writeNumber(h, uint64(coordinate))
	for _, point := range points {
		writeField(h, point.Marshal())
	}
	return boundHashToScalar(h.Sum(nil))
}

func verifyBoundDualRS(pp *PublicParameters, instance *Instance, binding []byte) bool {
	if len(binding) == 0 {
		return verifyRandomDualRS(pp, instance.V)
	}
	return verifyDualRS(pp, instance.V, boundDualCoefficients(pp, instance, binding))
}

func boundDualCoefficients(pp *PublicParameters, instance *Instance, binding []byte) []*big.Int {
	h := sha256.New()
	writeField(h, []byte("RB-PVSS/dual/transcript/v1"))
	writeField(h, binding)
	for _, point := range instance.V {
		writeField(h, point.Marshal())
	}
	for i := range instance.C {
		writeField(h, instance.C[i].Marshal())
		writeField(h, instance.Phi.A[i].Marshal())
		writeField(h, instance.Phi.B[i].Marshal())
		writeField(h, normalizeScalar(instance.Phi.Z[i]).Bytes())
	}
	writeField(h, normalizeScalar(instance.Phi.Chi).Bytes())
	seed := h.Sum(nil)
	coefficients := make([]*big.Int, len(pp.N)+len(pp.D)-pp.Threshold)
	for i := range coefficients {
		counter := sha256.New()
		writeField(counter, []byte("RB-PVSS/dual/coefficient/v1"))
		writeField(counter, seed)
		writeNumber(counter, uint64(i))
		coefficients[i] = boundHashToScalar(counter.Sum(nil))
	}
	return coefficients
}

// Rejection sampling avoids modulo bias in the new context-bound hash domains.
func boundHashToScalar(seed []byte) *big.Int {
	space := new(big.Int).Lsh(big.NewInt(1), 256)
	limit := new(big.Int).Sub(space, new(big.Int).Mod(new(big.Int).Set(space), bn256.Order))
	for attempt := uint64(0); ; attempt++ {
		h := sha256.New()
		writeField(h, []byte("RB-PVSS/scalar/v1"))
		writeField(h, seed)
		writeNumber(h, attempt)
		value := new(big.Int).SetBytes(h.Sum(nil))
		if value.Cmp(limit) < 0 {
			return value.Mod(value, bn256.Order)
		}
	}
}

func canonicalBallotScalars(instance *Instance) bool {
	valid := func(x *big.Int) bool {
		return x != nil && x.Sign() >= 0 && x.Cmp(bn256.Order) < 0
	}
	if !valid(instance.Phi.Chi) {
		return false
	}
	for _, z := range instance.Phi.Z {
		if !valid(z) {
			return false
		}
	}
	for _, proof := range instance.Pi {
		if !valid(proof.Chi) || !valid(proof.Z1) || !valid(proof.ZBeta) || !valid(proof.Z3) {
			return false
		}
	}
	return true
}
