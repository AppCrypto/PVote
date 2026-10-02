// Package ec contains the small BN254 helpers shared by the independent
// PriScore, 3PVSS, and Huang reproductions. Protocol algorithms remain in their
// respective packages; this package only centralizes constant-time group and
// scalar plumbing.
package ec

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"math/big"

	bn256 "github.com/ethereum/go-ethereum/crypto/bn256/cloudflare"
)

var Order = new(big.Int).Set(bn256.Order)

func Scalar() (*big.Int, error) {
	for {
		x, err := rand.Int(rand.Reader, Order)
		if err != nil {
			return nil, err
		}
		if x.Sign() != 0 {
			return x, nil
		}
	}
}

func Mod(x *big.Int) *big.Int {
	if x == nil {
		return new(big.Int)
	}
	z := new(big.Int).Mod(new(big.Int).Set(x), Order)
	if z.Sign() < 0 {
		z.Add(z, Order)
	}
	return z
}

func AddScalar(a, b *big.Int) *big.Int { return Mod(new(big.Int).Add(a, b)) }
func SubScalar(a, b *big.Int) *big.Int { return Mod(new(big.Int).Sub(a, b)) }
func MulScalar(a, b *big.Int) *big.Int { return Mod(new(big.Int).Mul(a, b)) }
func Inv(a *big.Int) *big.Int          { return new(big.Int).ModInverse(Mod(a), Order) }

func Zero() *bn256.G1                       { return new(bn256.G1).ScalarBaseMult(new(big.Int)) }
func Add(a, b *bn256.G1) *bn256.G1          { return new(bn256.G1).Add(a, b) }
func Mul(a *bn256.G1, s *big.Int) *bn256.G1 { return new(bn256.G1).ScalarMult(a, Mod(s)) }
func Neg(a *bn256.G1) *bn256.G1             { return Mul(a, new(big.Int).Sub(Order, big.NewInt(1))) }
func Sub(a, b *bn256.G1) *bn256.G1          { return Add(a, Neg(b)) }
func Equal(a, b *bn256.G1) bool             { return a != nil && b != nil && bytes.Equal(a.Marshal(), b.Marshal()) }
func EqualGT(a, b *bn256.GT) bool {
	return a != nil && b != nil && bytes.Equal(a.Marshal(), b.Marshal())
}

func Product(points []*bn256.G1) *bn256.G1 {
	acc := Zero()
	for _, p := range points {
		acc = Add(acc, p)
	}
	return acc
}

func HashScalar(domain string, parts ...[]byte) *big.Int {
	h := sha256.New()
	writePart(h, []byte(domain))
	for _, part := range parts {
		writePart(h, part)
	}
	return Mod(new(big.Int).SetBytes(h.Sum(nil)))
}

type byteWriter interface{ Write([]byte) (int, error) }

func writePart(w byteWriter, value []byte) {
	var n [8]byte
	binary.BigEndian.PutUint64(n[:], uint64(len(value)))
	_, _ = w.Write(n[:])
	_, _ = w.Write(value)
}

func PointBytes(points ...*bn256.G1) []byte {
	var out []byte
	for _, p := range points {
		if p == nil {
			out = append(out, make([]byte, 64)...)
			continue
		}
		out = append(out, p.Marshal()...)
	}
	return out
}

func ScalarBytes(values ...*big.Int) []byte {
	var out []byte
	for _, value := range values {
		buf := make([]byte, 32)
		if value != nil {
			value.FillBytes(buf)
		}
		out = append(out, buf...)
	}
	return out
}

// DLog returns x in [0,max] such that target=base*x.
func DLog(base, target *bn256.G1, max int) (int, error) {
	if max < 0 {
		return 0, errors.New("negative discrete-log bound")
	}
	if max > 2048 {
		m := 1
		for m*m <= max {
			m++
		}
		table := make(map[string]int, m)
		cur := Zero()
		for j := 0; j < m; j++ {
			table[string(cur.Marshal())] = j
			cur = Add(cur, base)
		}
		step := Neg(Mul(base, big.NewInt(int64(m))))
		giant := target
		for i := 0; i <= max/m; i++ {
			if j, ok := table[string(giant.Marshal())]; ok {
				x := i*m + j
				if x <= max {
					return x, nil
				}
			}
			giant = Add(giant, step)
		}
		return 0, errors.New("bounded discrete logarithm not found")
	}
	cur := Zero()
	for x := 0; x <= max; x++ {
		if Equal(cur, target) {
			return x, nil
		}
		cur = Add(cur, base)
	}
	return 0, errors.New("bounded discrete logarithm not found")
}
