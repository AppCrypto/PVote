package ZKRP

import (
	"fmt"

	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"math/big"

	bn256 "github.com/ethereum/go-ethereum/crypto/bn256/google"
)

type PP struct {
	G0      *bn256.G1
	H0      *bn256.G1
	G1      *bn256.G2
	PKI     *bn256.G2
	Sigma_k []*bn256.G1
}

type ZKRP struct {
	Ej  *bn256.G1
	Fj  *bn256.GT
	_Uj *bn256.G1
	C   *big.Int
	Z1  *big.Int
	Z2  *big.Int
	Z3  *big.Int
}

func Setup(a int, b int) PP {
	g0 := new(bn256.G1).ScalarBaseMult(big.NewInt(int64(1)))
	hScalar, _ := new(big.Int).SetString("9868996996480530350723936346388037348513707152826932716320380442065450531909", 10)

	h0 := new(bn256.G1).ScalarBaseMult(hScalar)
	g1 := new(bn256.G2).ScalarBaseMult(big.NewInt(int64(1)))

	sk, _ := rand.Int(rand.Reader, bn256.Order)
	pk := new(bn256.G2).ScalarMult(g1, sk)

	sigma_k := make([]*bn256.G1, b-a+1)
	for i := a; i <= b; i++ {
		den := new(big.Int).Add(big.NewInt(int64(i)), sk)
		den.ModInverse(den, bn256.Order)
		fmt.Printf("den=%v\n", den)
		temp := new(bn256.G1).ScalarMult(g0, den)
		sigma_k[i-a] = temp
	}
	return PP{
		G0:      g0,
		H0:      h0,
		G1:      g1,
		PKI:     pk,
		Sigma_k: sigma_k,
	}
}

func Prove(g0 *bn256.G1, h0 *bn256.G1, g1 *bn256.G2, sj *big.Int, wj *big.Int, Uj *bn256.G1, sigma_wj *bn256.G1) *ZKRP {
	l, _ := rand.Int(rand.Reader, bn256.Order)
	_wj, _ := rand.Int(rand.Reader, bn256.Order)
	_l, _ := rand.Int(rand.Reader, bn256.Order)
	_sj, _ := rand.Int(rand.Reader, bn256.Order)

	Ej := new(bn256.G1).ScalarMult(sigma_wj, l)
	Fj1 := bn256.Pair(new(bn256.G1).Neg(new(bn256.G1).ScalarMult(Ej, _wj)), g1)
	Fj2 := new(bn256.GT).ScalarMult(bn256.Pair(g0, g1), _l)
	Fj := new(bn256.GT).Add(Fj1, Fj2)

	_Uj := new(bn256.G1).Add(new(bn256.G1).ScalarMult(h0, _wj), new(bn256.G1).ScalarMult(g0, _sj))

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

	z2 := new(big.Int).Mul(c, l)
	z2.Sub(_l, z2)
	z2.Mod(z2, bn256.Order)

	z3 := new(big.Int).Mul(c, sj)
	z3.Sub(_sj, z3)
	z3.Mod(z3, bn256.Order)

	return &ZKRP{
		Ej:  Ej,
		Fj:  Fj,
		_Uj: _Uj,
		C:   c,
		Z1:  z1,
		Z2:  z2,
		Z3:  z3,
	}
}

func Verify(g0 *bn256.G1, h0 *bn256.G1, g1 *bn256.G2, pk *bn256.G2, proof *ZKRP, Uj *bn256.G1) bool {
	temp1 := new(bn256.G1).ScalarMult(Uj, proof.C)
	temp1 = new(bn256.G1).Add(temp1, new(bn256.G1).ScalarMult(g0, proof.Z3))
	temp1 = new(bn256.G1).Add(temp1, new(bn256.G1).ScalarMult(h0, proof.Z1))
	if bytes.Equal(proof._Uj.Marshal(), temp1.Marshal()) {
		right1 := bn256.Pair(new(bn256.G1).ScalarMult(proof.Ej, proof.C), pk)
		right2 := bn256.Pair(new(bn256.G1).Neg(new(bn256.G1).ScalarMult(proof.Ej, proof.Z1)), g1)
		right3 := bn256.Pair(new(bn256.G1).ScalarMult(g0, proof.Z2), g1)
		right := new(bn256.GT).Add(right1, new(bn256.GT).Add(right2, right3))
		if bytes.Equal(proof.Fj.Marshal(), right.Marshal()) {
			return true
		}
	}
	return false
}
