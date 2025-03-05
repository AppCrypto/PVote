package PVSS

import (
	"crypto/rand"
	"crypto/sha256"

	"bytes"
	"math/big"

	bn256 "github.com/ethereum/go-ethereum/crypto/bn256/cloudflare"
)

type Proof struct {
	C  *big.Int
	Z  *big.Int
	RG *bn256.G1
	RH *bn256.G1
}

// 需要运算的参数
type SecretSharing struct {
	BindValue []*big.Int
	Shares    []*big.Int
	V         []*bn256.G1
	C         []*bn256.G1
	Proofs    []Proof
}

// PVSS.Setup
func Setup(numTalliers int, g0 *bn256.G1) ([]*big.Int, []*bn256.G1) {
	sks := make([]*big.Int, numTalliers)
	pks := make([]*bn256.G1, numTalliers)
	for i := 0; i < numTalliers; i++ {
		sks[i], _ = rand.Int(rand.Reader, bn256.Order)
		pks[i] = new(bn256.G1).ScalarMult(g0, sks[i])
	}
	return sks, pks
}

func DLEQProof(G, H *bn256.G1, xG, xH *bn256.G1, x *big.Int) Proof {
	r, _ := rand.Int(rand.Reader, bn256.Order)
	rG := new(bn256.G1).ScalarMult(G, r)
	rH := new(bn256.G1).ScalarMult(H, r)

	new_hash := sha256.New()
	new_hash.Write(xG.Marshal())
	new_hash.Write(xH.Marshal())
	new_hash.Write(rG.Marshal())
	new_hash.Write(rH.Marshal())

	cb := new_hash.Sum(nil)
	c := new(big.Int).SetBytes(cb)
	c.Mod(c, bn256.Order)

	z := new(big.Int).Mul(c, x)
	z.Sub(r, z)
	z.Mod(z, bn256.Order)

	return Proof{
		C:  c,
		Z:  z,
		RG: rG,
		RH: rH,
	}
}

// PVSS.Share: Generate PVSS Shares({Vj},{Cj},DLEQ proofs)
func Share(secret *big.Int, h *bn256.G1, pks []*bn256.G1, threshold, numTalliers, numCandidates int) *SecretSharing {
	coefficients := make([]*big.Int, threshold)
	coefficients[0] = secret
	for i := 1; i < threshold; i++ {
		coefficients[i], _ = rand.Int(rand.Reader, bn256.Order)
	}
	bindValue := make([]*big.Int, numCandidates)
	for i := 0; i < numCandidates; i++ {
		x := new(big.Int).Neg(big.NewInt(int64(i)))
		x.Mod(x, bn256.Order)
		bindValue[i] = EvaluatePolynomial(coefficients, x, bn256.Order)
	}
	shares := make([]*big.Int, numTalliers)
	for i := 0; i < numTalliers; i++ {
		x := big.NewInt(int64(i + 1))
		shares[i] = EvaluatePolynomial(coefficients, x, bn256.Order)
	}

	v := make([]*bn256.G1, numTalliers+numCandidates)
	for i := 0; i < numTalliers; i++ {
		v[i] = new(bn256.G1).ScalarMult(h, shares[i])
	}
	for i := numTalliers; i < numTalliers+numCandidates; i++ {
		v[i] = new(bn256.G1).ScalarMult(h, bindValue[i-numTalliers])
	}

	//Generate the encrypted shares under talliers' public key
	c := make([]*bn256.G1, numTalliers)
	for i := 0; i < numTalliers; i++ {
		c[i] = new(bn256.G1).ScalarMult(pks[i], shares[i])
	}

	//Generate the DLEQ proofs for each encrypted shares
	proofs := make([]Proof, numTalliers)
	for i := 0; i < numTalliers; i++ {
		proofs[i] = DLEQProof(h, pks[i], v[i], c[i], shares[i])
	}

	return &SecretSharing{
		BindValue: bindValue,
		Shares:    shares,
		V:         v,
		C:         c,
		Proofs:    proofs,
	}
}

// Verify verifies the DLEQ proof
func DLEQVerify(c, z *big.Int, G, H, xG, xH, rG, rH *bn256.G1) bool {
	zG := new(bn256.G1).ScalarMult(G, z)
	zH := new(bn256.G1).ScalarMult(H, z)
	cxG := new(bn256.G1).ScalarMult(xG, c)
	cxH := new(bn256.G1).ScalarMult(xH, c)
	a := new(bn256.G1).Add(zG, cxG)
	b := new(bn256.G1).Add(zH, cxH)
	return (rG.String() == a.String() && rH.String() == b.String())

}

// Reed Solomon check
func Coefficient(i int, n int, l int) *big.Int {
	result := big.NewInt(1)
	for j := -(l - 1); j <= n; j++ {
		if i != j {
			result = result.Mul(result, new(big.Int).ModInverse(new(big.Int).Sub(big.NewInt(int64(i)), big.NewInt(int64(j))), bn256.Order))
			result = result.Mod(result, bn256.Order)

		}
	}
	return result
}

func RScodeVerify(shares []*bn256.G1, H1 *bn256.G1, n int, l int) bool {
	codeword := make([]*big.Int, n+l)
	for i := 1; i <= n; i++ {
		codeword[i-1] = new(big.Int).Mod(Coefficient(i, n, l), bn256.Order)
	}
	for i := n; i < n+l; i++ {
		codeword[i] = new(big.Int).Mod(Coefficient(-(i-n), n, l), bn256.Order)
	}

	sum := new(bn256.G1).ScalarMult(H1, big.NewInt(1))
	for i := 0; i < n+l; i++ {
		sum = new(bn256.G1).Add(sum, new(bn256.G1).ScalarMult(shares[i], codeword[i]))
	}
	return bytes.Equal(sum.Marshal(), H1.Marshal())
}

// PVSS.DVerify
func DVerify(secretsharing *SecretSharing, h *bn256.G1, pks []*bn256.G1, numTalliers int, numCandidates int) bool {
	for i := 0; i < len(pks); i++ {
		if !DLEQVerify(secretsharing.Proofs[i].C, secretsharing.Proofs[i].Z, h, pks[i], secretsharing.V[i], secretsharing.C[i], secretsharing.Proofs[i].RG, secretsharing.Proofs[i].RH) {
			return false
		}
	}
	return RScodeVerify(secretsharing.V, h, numTalliers, numCandidates)
}

// PVSS.Decrypt
func Decrypt(h *bn256.G1, pk *bn256.G1, c *bn256.G1, sk *big.Int) (*bn256.G1, Proof) {
	skInverse := new(big.Int).ModInverse(sk, bn256.Order)
	sh := new(bn256.G1).ScalarMult(c, skInverse)
	proof := DLEQProof(h, sh, pk, c, sk)
	return sh, proof
}

// PVSS.PVerify
func PVerify(h *bn256.G1, pk *bn256.G1, c *bn256.G1, sh *bn256.G1, proof Proof) bool {
	return DLEQVerify(proof.C, proof.Z, h, sh, pk, c, proof.RG, proof.RH)
}

// Evaluate the polynomial at a given x
func EvaluatePolynomial(coefficients []*big.Int, x, order *big.Int) *big.Int {
	result := new(big.Int).Set(coefficients[0])
	xPower := new(big.Int).Set(x)

	for i := 1; i < len(coefficients); i++ {
		term := new(big.Int).Mul(coefficients[i], xPower)
		term.Mod(term, order)
		result.Add(result, term)
		result.Mod(result, order)
		xPower.Mul(xPower, x)
		xPower.Mod(xPower, order)
	}

	return result
}

// Calculate the lagrange coefficient at l
func LagrangeCoefficient(l *big.Int, indices []*big.Int, threshold int) []*big.Int {

	coefficient := make([]*big.Int, threshold)

	for i := 0; i < threshold; i++ {
		num := big.NewInt(1)
		den := big.NewInt(1)

		for j := 0; j < threshold; j++ {
			if i != j {
				num.Mul(num, new(big.Int).Sub(l, indices[j]))
				num.Mod(num, bn256.Order)

				den.Mul(den, new(big.Int).Sub(indices[i], indices[j]))
				den.Mod(den, bn256.Order)
			}
		}

		den.ModInverse(den, bn256.Order)
		term := new(big.Int).Mul(num, den)
		term.Mod(term, bn256.Order)
		coefficient[i] = term
	}
	return coefficient
}

func Reconstruct(coefficient []*big.Int, sh []*bn256.G1) *bn256.G1 {
	secret := new(bn256.G1).ScalarBaseMult(big.NewInt(0))
	for i := 0; i < len(coefficient); i++ {
		secret = new(bn256.G1).Add(secret, new(bn256.G1).ScalarMult(sh[i], coefficient[i]))
	}
	return secret
}
