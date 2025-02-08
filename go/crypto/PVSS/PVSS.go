package PVSS

import (
	"crypto/rand"
	"crypto/sha256"

	//"fmt"
	"bytes"
	"math/big"

	bn256 "github.com/ethereum/go-ethereum/crypto/bn256/google"
)

type DLEQ struct {
	C  *big.Int
	Z  *big.Int
	RG *bn256.G1
	RH *bn256.G1
}

// 需要运算的参数
type SecretSharing struct {
	//Fi []*big.Int
	V      []*bn256.G1
	C      []*bn256.G1
	Proofs []DLEQ
}

// PVSS.Setup
func Setup(numTalliers int) ([]*big.Int, []*bn256.G1) {
	sks := make([]*big.Int, numTalliers)
	pks := make([]*bn256.G1, numTalliers)
	for i := 0; i < numTalliers; i++ {
		sks[i], _ = rand.Int(rand.Reader, bn256.Order)
		pks[i] = new(bn256.G1).ScalarBaseMult(sks[i])
	}
	return sks, pks
}

func DLEQProof(G, H *bn256.G1, xG, xH *bn256.G1, x *big.Int) DLEQ {
	//生成承诺
	r, _ := rand.Int(rand.Reader, bn256.Order)
	rG := new(bn256.G1).ScalarMult(G, r)
	rH := new(bn256.G1).ScalarMult(H, r)

	// 计算挑战
	new_hash := sha256.New()
	new_hash.Write(xG.Marshal())
	new_hash.Write(xH.Marshal())
	new_hash.Write(rG.Marshal())
	new_hash.Write(rH.Marshal())

	cb := new_hash.Sum(nil)
	c := new(big.Int).SetBytes(cb)
	c.Mod(c, bn256.Order)

	// 生成相应
	z := new(big.Int).Mul(c, x)
	z.Sub(r, z)
	z.Mod(z, bn256.Order)

	return DLEQ{
		C:  c,
		Z:  z,
		RG: rG,
		RH: rH,
	}
}

// PVSS.Share: Generate PVSS Shares({Vj},{Cj},DLEQ proofs)
func Share(secret *big.Int, h *bn256.G1, pks []*bn256.G1, threshold, numShares int) *SecretSharing {
	coefficients := make([]*big.Int, threshold)
	coefficients[0] = secret
	for i := 1; i < threshold; i++ {
		coefficients[i], _ = rand.Int(rand.Reader, bn256.Order)
	}
	shares := make([]*big.Int, numShares)
	for i := 0; i < numShares; i++ {
		x := big.NewInt(int64(i + 1))
		shares[i] = EvaluatePolynomial(coefficients, x, bn256.Order)
	}

	v := make([]*bn256.G1, numShares)
	for i := 0; i < numShares; i++ {
		v[i] = new(bn256.G1).ScalarMult(h, shares[i])
	}

	//Generate the encrypted shares under talliers' public key
	c := make([]*bn256.G1, numShares)
	for i := 0; i < numShares; i++ {
		c[i] = new(bn256.G1).ScalarMult(pks[i], shares[i])
	}

	//Generate the DLEQ proofs for each encrypted shares
	proofs := make([]DLEQ, numShares)
	for i := 0; i < numShares; i++ {
		proofs[i] = DLEQProof(h, pks[i], v[i], c[i], shares[i])
	}

	return &SecretSharing{
		//Shares: shares,
		V:      v,
		C:      c,
		Proofs: proofs,
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
	if !(rG.String() == a.String() && rH.String() == b.String()) {
		return false
	}
	return true
}

// Reed Solomon check
func coefficient(i int, n int) *big.Int {
	result := big.NewInt(1)
	for j := 1; j < n+1; j++ {
		if i != j {
			result = result.Mul(result, new(big.Int).ModInverse(new(big.Int).Sub(big.NewInt(int64(i)), big.NewInt(int64(j))), bn256.Order))
			result = result.Mod(result, bn256.Order)
		}
	}
	return result
}

func RScodeVerify(shares []*bn256.G1, H1 *bn256.G1) bool {
	codeword := make([]*big.Int, len(shares))
	for i := 1; i <= len(shares); i++ {
		codeword[i-1] = new(big.Int).Mod(coefficient(i, len(shares)), bn256.Order)
	}
	sum := new(bn256.G1).ScalarMult(H1, big.NewInt(1))
	for i := 0; i < len(shares); i++ {
		sum = new(bn256.G1).Add(sum, new(bn256.G1).ScalarMult(shares[i], codeword[i]))
	}
	if bytes.Equal(sum.Marshal(), H1.Marshal()) {
		return true // if verify pass return true
	}
	return false // else return false
}

// PVSS.DVerify
func DVerify(secretsharing *SecretSharing, h *bn256.G1, pks []*bn256.G1) bool {
	for i := 0; i < len(secretsharing.V); i++ {
		if !DLEQVerify(secretsharing.Proofs[i].C, secretsharing.Proofs[i].Z, h, pks[i], secretsharing.V[i], secretsharing.C[i], secretsharing.Proofs[i].RG, secretsharing.Proofs[i].RH) {
			return false
		}
	}
	if !RScodeVerify(secretsharing.V, h) {
		return false
	}
	return true
}

// PVSS.Decrypt
func Decrypt(h *bn256.G1, pk *bn256.G1, c *bn256.G1, sk *big.Int) (*bn256.G1, DLEQ) {
	sh := new(bn256.G1).ScalarMult(c, sk)
	proof := DLEQProof(h, pk, sh, c, sk)
	return sh, proof
}

// PVSS.PVerify
func PVerify(h *bn256.G1, pk *bn256.G1, c *bn256.G1, sh *bn256.G1, proof *DLEQ) bool {
	if !DLEQVerify(proof.C, proof.Z, h, pk, sh, c, proof.RG, proof.RH) {
		return false
	}
	return true
}

// evaluatePolynomial 在给定的 x 处计算多项式的值
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

// lagrangeInterpolation 求拉格朗日插值法的系数
func LagrangeCoefficient(l *big.Int, indices []*big.Int, threshold int) []*big.Int {
	// k是分享的数量
	coefficient := make([]*big.Int, threshold)

	// 对于每个分享
	for i := 0; i < threshold; i++ {
		// 初始化分子（num）和分母（den）为1
		num := big.NewInt(1)
		den := big.NewInt(1)

		// 计算拉格朗日基函数的分子和分母
		for j := 0; j < threshold; j++ {
			if i != j {
				// 分子累乘 -indices[j]
				num.Mul(num, new(big.Int).Sub(l, indices[j]))
				num.Mod(num, bn256.Order)

				// 分母累乘 indices[i] - indices[j]
				den.Mul(den, new(big.Int).Sub(indices[i], indices[j]))
				den.Mod(den, bn256.Order)
			}
		}

		// 计算分母的逆元（模order）
		den.ModInverse(den, bn256.Order)
		// 计算每一项的值 shares[i] * num * den
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
