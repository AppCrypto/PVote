package main

import (
	//"context"
	"PVote/crypto/PVSS"
	"fmt"

	//"dttp/utils"k
	//"fmt"
	//"log"
	"math/big"
	//"time"
	//"github.com/ethereum/go-ethereum/accounts/abi/bind"
	//"github.com/ethereum/go-ethereum/common"
	//"github.com/ethereum/go-ethereum/ethclient"
	//"github.com/ethereum/go-ethereum/core/types"
	bn256 "github.com/ethereum/go-ethereum/crypto/bn256/google"
)

func main() {
	numShares := 3
	threshold := 2

	hScalar, _ := new(big.Int).SetString("9868996996480530350723936346388037348513707152826932716320380442065450531909", 10)
	// 计算 H1 = g * scalar
	var h = new(bn256.G1).ScalarBaseMult(hScalar)

	secret := big.NewInt(3)
	fmt.Printf("Test Secret=%v\n", new(bn256.G1).ScalarMult(h, secret))

	//Each tallier registers own key pair
	_, PKs := PVSS.Setup(numShares)
	Shares := PVSS.Share(secret, h, PKs, threshold, numShares)

	//使用 SetString 初始化大数
	selectedIndices := make([]*big.Int, threshold)
	for i := 0; i < threshold; i++ {
		selectedIndices[i] = big.NewInt(int64(i + 1))
	}
	Coefficient := PVSS.LagrangeCoefficient(big.NewInt(0), selectedIndices, threshold)
	fmt.Printf("Coefficient=%v\n", Coefficient)
	selectedShares := Shares.V[:threshold]
	result := PVSS.Reconstruct(Coefficient, selectedShares)
	fmt.Printf("Reconstruct result=%v\n", result)

}
