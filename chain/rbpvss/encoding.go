// Package rbpvss provides the lossless boundary between the paper's Go
// implementation and its Solidity verifier. It contains no cryptography.
package rbpvss

import (
	"fmt"
	"math/big"

	contract "PVote/compile/contract"
	core "PVote/crypto/RBPVSS"

	bn256 "github.com/ethereum/go-ethereum/crypto/bn256/cloudflare"
)

// ConstructorArguments converts pp into the exact Solidity constructor order.
func ConstructorArguments(pp *core.PublicParameters) (
	contract.RBPVSSVerifierG1Point,
	contract.RBPVSSVerifierG1Point,
	contract.RBPVSSVerifierG2Point,
	contract.RBPVSSVerifierG2Point,
	[]contract.RBPVSSVerifierG1Point,
	[]contract.RBPVSSVerifierG1Point,
	*big.Int, *big.Int, *big.Int, *big.Int, error,
) {
	if err := core.ValidatePublicParameters(pp); err != nil {
		return contract.RBPVSSVerifierG1Point{}, contract.RBPVSSVerifierG1Point{},
			contract.RBPVSSVerifierG2Point{}, contract.RBPVSSVerifierG2Point{}, nil, nil,
			nil, nil, nil, nil, err
	}
	return G1Point(pp.PP1.G0), G1Point(pp.PP1.H0), G2Point(pp.PP2.G1), G2Point(pp.PP2.PKI),
		G1Points(pp.PP2.Sigma), G1Points(pp.PK), big.NewInt(int64(pp.Threshold)),
		big.NewInt(int64(pp.PP1.L)), big.NewInt(pp.Range.Min), big.NewInt(pp.Range.Max), nil
}

// Transcript converts one complete atomic RB-PVSS instance. The two returned
// tuples must be supplied together to submitRB; neither is independently valid.
func Transcript(instance *core.Instance) (contract.RBPVSSVerifierDealerTranscript, contract.RBPVSSVerifierRangeBindingTranscript, error) {
	if instance == nil {
		return contract.RBPVSSVerifierDealerTranscript{}, contract.RBPVSSVerifierRangeBindingTranscript{}, fmt.Errorf("nil RB-PVSS transcript")
	}
	dealer := contract.RBPVSSVerifierDealerTranscript{
		V: G1Points(instance.V), C: G1Points(instance.C), PhiA: G1Points(instance.Phi.A),
		PhiB: G1Points(instance.Phi.B), PhiChi: cloneInt(instance.Phi.Chi), PhiZ: cloneInts(instance.Phi.Z),
	}
	rangeBinding := contract.RBPVSSVerifierRangeBindingTranscript{
		U: G1Points(instance.U), E: make([]contract.RBPVSSVerifierG1Point, len(instance.Pi)),
		FBase: make([]contract.RBPVSSVerifierG1Point, len(instance.Pi)), UPrime: make([]contract.RBPVSSVerifierG1Point, len(instance.Pi)),
		CPrime: make([]contract.RBPVSSVerifierG1Point, len(instance.Pi)), Chi: make([]*big.Int, len(instance.Pi)),
		Z1: make([]*big.Int, len(instance.Pi)), ZBeta: make([]*big.Int, len(instance.Pi)), Z3: make([]*big.Int, len(instance.Pi)),
	}
	for i, proof := range instance.Pi {
		if proof == nil || proof.FEncoding() == nil {
			return contract.RBPVSSVerifierDealerTranscript{}, contract.RBPVSSVerifierRangeBindingTranscript{}, fmt.Errorf("nil range proof %d", i)
		}
		rangeBinding.E[i] = G1Point(proof.E)
		rangeBinding.FBase[i] = G1Point(proof.FEncoding())
		rangeBinding.UPrime[i] = G1Point(proof.UPrime)
		rangeBinding.CPrime[i] = G1Point(proof.CPrime)
		rangeBinding.Chi[i] = cloneInt(proof.Chi)
		rangeBinding.Z1[i] = cloneInt(proof.Z1)
		rangeBinding.ZBeta[i] = cloneInt(proof.ZBeta)
		rangeBinding.Z3[i] = cloneInt(proof.Z3)
	}
	return dealer, rangeBinding, nil
}

// G1Point converts a Go BN254 point to the Solidity affine representation.
func G1Point(point *bn256.G1) contract.RBPVSSVerifierG1Point {
	if point == nil {
		return contract.RBPVSSVerifierG1Point{}
	}
	encoded := point.Marshal()
	return contract.RBPVSSVerifierG1Point{X: new(big.Int).SetBytes(encoded[:32]), Y: new(big.Int).SetBytes(encoded[32:])}
}

// G1 decodes a Solidity affine point and rejects malformed coordinates.
func G1(point contract.RBPVSSVerifierG1Point) (*bn256.G1, error) {
	encoded := append(leftPad32(point.X), leftPad32(point.Y)...)
	decoded := new(bn256.G1)
	if _, err := decoded.Unmarshal(encoded); err != nil {
		return nil, fmt.Errorf("decode G1: %w", err)
	}
	return decoded, nil
}

// G2Point converts a Go BN254 G2 point in the precompile's coefficient order.
func G2Point(point *bn256.G2) contract.RBPVSSVerifierG2Point {
	if point == nil {
		return contract.RBPVSSVerifierG2Point{}
	}
	encoded := point.Marshal()
	return contract.RBPVSSVerifierG2Point{
		X: [2]*big.Int{new(big.Int).SetBytes(encoded[:32]), new(big.Int).SetBytes(encoded[32:64])},
		Y: [2]*big.Int{new(big.Int).SetBytes(encoded[64:96]), new(big.Int).SetBytes(encoded[96:])},
	}
}

func G1Points(points []*bn256.G1) []contract.RBPVSSVerifierG1Point {
	converted := make([]contract.RBPVSSVerifierG1Point, len(points))
	for i, point := range points {
		converted[i] = G1Point(point)
	}
	return converted
}

func cloneInt(value *big.Int) *big.Int {
	if value == nil {
		return nil
	}
	return new(big.Int).Set(value)
}

func cloneInts(values []*big.Int) []*big.Int {
	cloned := make([]*big.Int, len(values))
	for i, value := range values {
		cloned[i] = cloneInt(value)
	}
	return cloned
}

func leftPad32(value *big.Int) []byte {
	result := make([]byte, 32)
	if value == nil {
		return result
	}
	raw := value.Bytes()
	if len(raw) > len(result) {
		raw = raw[len(raw)-len(result):]
	}
	copy(result[len(result)-len(raw):], raw)
	return result
}
