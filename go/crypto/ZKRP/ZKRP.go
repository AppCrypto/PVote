package ZKRP

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"math/big"

	bn256 "github.com/ethereum/go-ethereum/crypto/bn256/cloudflare"
)

type PP struct {
	G0      *bn256.G1
	H0      *bn256.G1
	G1      *bn256.G2
	PKI     *bn256.G2
	Sigma_k []*bn256.G1
}

type Proof struct {
	Ej  *bn256.G1
	Fj  *bn256.GT
	Fj1 *bn256.G1
	Fj2 *bn256.G1
	Uj  *bn256.G1
	Cj  *bn256.G1
	C   *big.Int
	Z1  *big.Int
	Z2  *big.Int
	Z3  *big.Int
}

func Setup(a int, b int) (*big.Int, PP) {
	g0 := new(bn256.G1).ScalarBaseMult(big.NewInt(int64(1)))
	hScalar, _ := new(big.Int).SetString("9868996996480530350723936346388037348513707152826932716320380442065450531909", 10)
	h0 := new(bn256.G1).ScalarBaseMult(hScalar)
	g1 := new(bn256.G2).ScalarBaseMult(big.NewInt(int64(1)))

	sk, _ := rand.Int(rand.Reader, bn256.Order)
	pk := new(bn256.G2).ScalarMult(g1, sk)

	sigma_k := make([]*bn256.G1, b-a+1)
	for i := 0; i <= b-a; i++ {
		den := new(big.Int).Add(big.NewInt(int64(i+a)), sk)
		den.ModInverse(den, bn256.Order)
		temp := new(bn256.G1).ScalarMult(g0, den)
		sigma_k[i] = temp
	}
	return sk, PP{
		G0:      g0,
		H0:      h0,
		G1:      g1,
		PKI:     pk,
		Sigma_k: sigma_k,
	}
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

func GenProof(g0 *bn256.G1, h0 *bn256.G1, g1 *bn256.G2, sj *big.Int, wj *big.Int, Uj *bn256.G1, sigma_wj *bn256.G1, d *big.Int, coefficients []*big.Int) *Proof {
	b, _ := rand.Int(rand.Reader, bn256.Order)
	_wj, _ := rand.Int(rand.Reader, bn256.Order)
	_b, _ := rand.Int(rand.Reader, bn256.Order)
	_sj := EvaluatePolynomial(coefficients, d, bn256.Order)

	Ej := new(bn256.G1).ScalarMult(sigma_wj, b)
	Fj1 := bn256.Pair(new(bn256.G1).Neg(new(bn256.G1).ScalarMult(Ej, _wj)), g1)
	Fj2 := new(bn256.GT).ScalarMult(bn256.Pair(g0, g1), _b)
	Fj := new(bn256.GT).Add(Fj1, Fj2)
	fj1 := new(bn256.G1).Neg(new(bn256.G1).ScalarMult(Ej, _wj))
	fj2 := new(bn256.G1).ScalarMult(g0, _b)

	_Uj := new(bn256.G1).Add(new(bn256.G1).ScalarMult(h0, _wj), new(bn256.G1).ScalarMult(g0, _sj))

	_Cj := new(bn256.G1).ScalarMult(h0, _sj)

	new_hash := sha256.New()
	new_hash.Write(Ej.Marshal())
	new_hash.Write(Uj.Marshal())
	new_hash.Write(Fj.Marshal())
	new_hash.Write(_Uj.Marshal())
	cb := new_hash.Sum(nil)
	c := new(big.Int).SetBytes(cb)
	c.Mod(c, bn256.Order)

	z1 := new(big.Int).Mul(c, wj)
	z1.Sub(_wj, z1)
	z1.Mod(z1, bn256.Order)

	z2 := new(big.Int).Mul(c, b)
	z2.Sub(_b, z2)
	z2.Mod(z2, bn256.Order)

	z3 := new(big.Int).Mul(c, sj)
	z3.Sub(_sj, z3)
	z3.Mod(z3, bn256.Order)

	return &Proof{
		Ej:  Ej,
		Fj:  Fj,
		Fj1: fj1,
		Fj2: fj2,
		Uj:  _Uj,
		Cj:  _Cj,
		C:   c,
		Z1:  z1,
		Z2:  z2,
		Z3:  z3,
	}
}

// Calculate the lagrange coefficient at d
func Interpolation(d *big.Int, v []*bn256.G1, indices []*big.Int, threshold int) *bn256.G1 {

	coefficient := make([]*big.Int, threshold)

	for i := 0; i < threshold; i++ {
		num := big.NewInt(1)
		den := big.NewInt(1)

		for j := 0; j < threshold; j++ {
			if i != j {
				num.Mul(num, new(big.Int).Sub(d, indices[j]))
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
	secret := new(bn256.G1).ScalarBaseMult(big.NewInt(0))
	for i := 0; i < len(coefficient); i++ {
		secret = new(bn256.G1).Add(secret, new(bn256.G1).ScalarMult(v[i], coefficient[i]))
	}
	return secret
}

func Verify(g0 *bn256.G1, h0 *bn256.G1, g1 *bn256.G2, pk *bn256.G2, proof *Proof, Uj *bn256.G1, d *big.Int, v []*bn256.G1, indices []*big.Int, threshold int) bool {
	temp := Interpolation(d, v, indices, threshold)
	temp = new(bn256.G1).ScalarMult(temp, proof.C)
	temp = new(bn256.G1).Add(temp, new(bn256.G1).ScalarMult(h0, proof.Z3))
	if bytes.Equal(proof.Cj.Marshal(), temp.Marshal()) {
		//fmt.Printf("The first verification Pass!\n")
		temp1 := new(bn256.G1).ScalarMult(Uj, proof.C)
		temp1 = new(bn256.G1).Add(temp1, new(bn256.G1).ScalarMult(g0, proof.Z3))
		temp1 = new(bn256.G1).Add(temp1, new(bn256.G1).ScalarMult(h0, proof.Z1))
		if bytes.Equal(proof.Uj.Marshal(), temp1.Marshal()) {
			//fmt.Printf("The second verification Pass!\n")
			right1 := bn256.Pair(new(bn256.G1).ScalarMult(proof.Ej, proof.C), pk)
			right2 := bn256.Pair(new(bn256.G1).Neg(new(bn256.G1).ScalarMult(proof.Ej, proof.Z1)), g1)
			right3 := bn256.Pair(new(bn256.G1).ScalarMult(g0, proof.Z2), g1)
			right := new(bn256.GT).Add(right1, new(bn256.GT).Add(right2, right3))
			if bytes.Equal(proof.Fj.Marshal(), right.Marshal()) {
				//fmt.Printf("The third verification Pass!\n")
				return true
			}
		}
	}
	return false
}
