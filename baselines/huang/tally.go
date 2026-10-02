// Package huang implements the tally-relevant exponential-ElGamal path of
// Huang et al. (IEEE TNSE 2022): one ciphertext pair per score coordinate,
// public aggregation of all active ballots, decryption with the election
// secret, and bounded discrete-log decoding. It intentionally does not claim
// to reproduce the paper's ballot proofs or abstention-recovery protocol.
package huang

import (
	"errors"
	"math/big"

	"PVote/baselines/internal/ec"
	bn256 "github.com/ethereum/go-ethereum/crypto/bn256/cloudflare"
)

var ErrInvalid = errors.New("invalid Huang tally input")

type Parameters struct {
	G           *bn256.G1
	PublicKey   *bn256.G1
	SecretKey   *big.Int
	Coordinates int
	MaxScore    int
}

type Ciphertext struct{ C1, C2 *bn256.G1 }
type Ballot []Ciphertext

func Setup(coordinates, maxScore int) (*Parameters, error) {
	if coordinates < 1 || maxScore < 0 {
		return nil, ErrInvalid
	}
	sk, err := ec.Scalar()
	if err != nil {
		return nil, err
	}
	g := new(bn256.G1).ScalarBaseMult(big.NewInt(1))
	return &Parameters{G: g, PublicKey: ec.Mul(g, sk), SecretKey: sk, Coordinates: coordinates, MaxScore: maxScore}, nil
}

func Encrypt(pp *Parameters, scores []int) (Ballot, error) {
	if pp == nil || len(scores) != pp.Coordinates {
		return nil, ErrInvalid
	}
	out := make(Ballot, len(scores))
	for j, score := range scores {
		if score < 0 || score > pp.MaxScore {
			return nil, ErrInvalid
		}
		r, err := ec.Scalar()
		if err != nil {
			return nil, err
		}
		out[j] = Ciphertext{C1: ec.Mul(pp.G, r), C2: ec.Add(ec.Mul(pp.G, big.NewInt(int64(score))), ec.Mul(pp.PublicKey, r))}
	}
	return out, nil
}

// Tally performs the public ciphertext products, election-key decryption,
// and bounded decoding for every score coordinate.
func Tally(pp *Parameters, ballots []Ballot) ([]int, error) {
	if pp == nil || len(ballots) == 0 {
		return nil, ErrInvalid
	}
	totals := make([]int, pp.Coordinates)
	for j := 0; j < pp.Coordinates; j++ {
		var a1, a2 *bn256.G1
		for _, ballot := range ballots {
			if len(ballot) != pp.Coordinates || ballot[j].C1 == nil || ballot[j].C2 == nil {
				return nil, ErrInvalid
			}
			if a1 == nil {
				a1, a2 = ballot[j].C1, ballot[j].C2
			} else {
				a1, a2 = ec.Add(a1, ballot[j].C1), ec.Add(a2, ballot[j].C2)
			}
		}
		plain := ec.Sub(a2, ec.Mul(a1, pp.SecretKey))
		decoded, ok := boundedDecode(pp.G, plain, len(ballots)*pp.MaxScore)
		if !ok {
			return nil, ErrInvalid
		}
		totals[j] = decoded
	}
	return totals, nil
}

func boundedDecode(g, target *bn256.G1, max int) (int, bool) {
	if max > 128 {
		m := 1
		for m*m <= max {
			m++
		}
		table := make(map[string]int, m)
		cur := new(bn256.G1).ScalarBaseMult(big.NewInt(0))
		for j := 0; j < m; j++ {
			table[string(cur.Marshal())] = j
			cur = ec.Add(cur, g)
		}
		step := ec.Neg(ec.Mul(g, big.NewInt(int64(m))))
		giant := target
		for i := 0; i <= max/m; i++ {
			if j, ok := table[string(giant.Marshal())]; ok {
				x := i*m + j
				if x <= max {
					return x, true
				}
			}
			giant = ec.Add(giant, step)
		}
		return 0, false
	}
	cur := new(bn256.G1).ScalarBaseMult(big.NewInt(0))
	for value := 0; value <= max; value++ {
		if ec.Equal(cur, target) {
			return value, true
		}
		cur = ec.Add(cur, g)
	}
	return 0, false
}
