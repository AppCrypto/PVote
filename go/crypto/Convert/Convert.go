package Convert

import (
	"PVote/compile/contract"
	"crypto/sha256"
	"encoding/base64"
	"math/big"

	bn256 "github.com/ethereum/go-ethereum/crypto/bn256/cloudflare"
)

func G1ToG1Point(bn256Point *bn256.G1) contract.VerificationG1Point {
	// Marshal the G1 point to get the X and Y coordinates as bytes
	point := bn256Point.Marshal()

	// Create big.Int for X and Y coordinates
	x := new(big.Int).SetBytes(point[:32])
	y := new(big.Int).SetBytes(point[32:64])

	g1Point := contract.VerificationG1Point{
		X: x,
		Y: y,
	}
	return g1Point
}

func G1PointToG1(g1point contract.VerificationG1Point) *bn256.G1 {
	xBytes := g1point.X.Bytes()
	yBytes := g1point.Y.Bytes()

	if len(xBytes) != 32 {
		xBytes = append([]byte{0x00}, xBytes...)
	}
	if len(yBytes) != 32 {
		yBytes = append([]byte{0x00}, yBytes...)
	}

	decodedBytes := append(xBytes, yBytes...)

	g1 := new(bn256.G1)
	g1.Unmarshal(decodedBytes)

	return g1
}

func G2ToG2Point(point *bn256.G2) contract.VerificationG2Point {
	// Marshal the G1 point to get the X and Y coordinates as bytes
	pointBytes := point.Marshal()

	// Create big.Int for X and Y coordinates
	a1 := new(big.Int).SetBytes(pointBytes[:32])
	a2 := new(big.Int).SetBytes(pointBytes[32:64])
	b1 := new(big.Int).SetBytes(pointBytes[64:96])
	b2 := new(big.Int).SetBytes(pointBytes[96:128])

	g2Point := contract.VerificationG2Point{
		X: [2]*big.Int{a1, a2},
		Y: [2]*big.Int{b1, b2},
	}
	return g2Point
}

func G2PointToG2(g2point contract.VerificationG2Point) *bn256.G2 {
	x1Bytes := g2point.X[0].Bytes()
	x2Bytes := g2point.X[1].Bytes()
	y1Bytes := g2point.Y[0].Bytes()
	y2Bytes := g2point.Y[1].Bytes()
	if len(x1Bytes) != 32 {
		x1Bytes = append([]byte{0x00}, x1Bytes...)
	}
	if len(x1Bytes) != 32 {
		x2Bytes = append([]byte{0x00}, x2Bytes...)
	}
	if len(y1Bytes) != 32 {
		y1Bytes = append([]byte{0x00}, y1Bytes...)
	}
	if len(y2Bytes) != 32 {
		y2Bytes = append([]byte{0x00}, y2Bytes...)
	}

	decodedBytes := append(x1Bytes, x2Bytes...)
	decodedBytes = append(decodedBytes, y1Bytes...)
	decodedBytes = append(decodedBytes, y2Bytes...)

	g2 := new(bn256.G2)
	g2.Unmarshal(decodedBytes)

	return g2
}

func GTToString(gt *bn256.GT) string {

	gtBytes := gt.Marshal()

	encoded := base64.StdEncoding.EncodeToString(gtBytes)
	return encoded
}

func StringToGT(encoded string) *bn256.GT {
	decodedBytes, _ := base64.StdEncoding.DecodeString(encoded)

	gt := new(bn256.GT)
	gt.Unmarshal(decodedBytes)
	return gt
}

func StringToBigInt(input string) *big.Int {

	hash := sha256.Sum256([]byte(input))
	bigIntValue := new(big.Int).SetBytes(hash[:])
	return bigIntValue
}

func G1ToBigIntArray(point *bn256.G1) [2]*big.Int {
	// Marshal the G1 point to get the X and Y coordinates as bytes
	pointBytes := point.Marshal()

	// Create big.Int for X and Y coordinates
	x := new(big.Int).SetBytes(pointBytes[:32])
	y := new(big.Int).SetBytes(pointBytes[32:64])

	return [2]*big.Int{x, y}
}
