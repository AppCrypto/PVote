// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contract

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
	_ = abi.ConvertType
)

// VerificationG1Point is an auto generated low-level Go binding around an user-defined struct.
type VerificationG1Point struct {
	X *big.Int
	Y *big.Int
}

// VerificationG2Point is an auto generated low-level Go binding around an user-defined struct.
type VerificationG2Point struct {
	X [2]*big.Int
	Y [2]*big.Int
}

// ContractMetaData contains all meta data concerning the Contract contract.
var ContractMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"X\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"Y\",\"type\":\"uint256\"}],\"internalType\":\"structVerification.G1Point\",\"name\":\"g\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"X\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"Y\",\"type\":\"uint256\"}],\"internalType\":\"structVerification.G1Point\",\"name\":\"y1\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"X\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"Y\",\"type\":\"uint256\"}],\"internalType\":\"structVerification.G1Point\",\"name\":\"a1\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"X\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"Y\",\"type\":\"uint256\"}],\"internalType\":\"structVerification.G1Point\",\"name\":\"h\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"X\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"Y\",\"type\":\"uint256\"}],\"internalType\":\"structVerification.G1Point\",\"name\":\"y2\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"X\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"Y\",\"type\":\"uint256\"}],\"internalType\":\"structVerification.G1Point\",\"name\":\"a2\",\"type\":\"tuple\"},{\"internalType\":\"uint256\",\"name\":\"c\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"z\",\"type\":\"uint256\"}],\"name\":\"DLEQVerify\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"X\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"Y\",\"type\":\"uint256\"}],\"internalType\":\"structVerification.G1Point\",\"name\":\"g\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"X\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"Y\",\"type\":\"uint256\"}],\"internalType\":\"structVerification.G1Point[]\",\"name\":\"y1\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"X\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"Y\",\"type\":\"uint256\"}],\"internalType\":\"structVerification.G1Point[]\",\"name\":\"a1\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"X\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"Y\",\"type\":\"uint256\"}],\"internalType\":\"structVerification.G1Point[]\",\"name\":\"h\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"X\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"Y\",\"type\":\"uint256\"}],\"internalType\":\"structVerification.G1Point[]\",\"name\":\"y2\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"X\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"Y\",\"type\":\"uint256\"}],\"internalType\":\"structVerification.G1Point[]\",\"name\":\"a2\",\"type\":\"tuple[]\"},{\"internalType\":\"uint256[]\",\"name\":\"c\",\"type\":\"uint256[]\"},{\"internalType\":\"uint256[]\",\"name\":\"z\",\"type\":\"uint256[]\"}],\"name\":\"DVerify\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"GetAggregateValue\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"X\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"Y\",\"type\":\"uint256\"}],\"internalType\":\"structVerification.G1Point[]\",\"name\":\"\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"GetDVerifyResult\",\"outputs\":[{\"internalType\":\"bool[]\",\"name\":\"\",\"type\":\"bool[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"GetPVerifyResult\",\"outputs\":[{\"internalType\":\"bool[]\",\"name\":\"\",\"type\":\"bool[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"GetTallyValue\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"X\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"Y\",\"type\":\"uint256\"}],\"internalType\":\"structVerification.G1Point[]\",\"name\":\"\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"GetZKRPResult\",\"outputs\":[{\"internalType\":\"bool[]\",\"name\":\"\",\"type\":\"bool[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"d\",\"type\":\"uint256\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"X\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"Y\",\"type\":\"uint256\"}],\"internalType\":\"structVerification.G1Point[]\",\"name\":\"v\",\"type\":\"tuple[]\"},{\"internalType\":\"uint256[]\",\"name\":\"indices\",\"type\":\"uint256[]\"},{\"internalType\":\"uint256\",\"name\":\"threshold\",\"type\":\"uint256\"}],\"name\":\"Interpolation\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"X\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"Y\",\"type\":\"uint256\"}],\"internalType\":\"structVerification.G1Point\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[]\",\"name\":\"indexTallier\",\"type\":\"uint256[]\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"X\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"Y\",\"type\":\"uint256\"}],\"internalType\":\"structVerification.G1Point[]\",\"name\":\"DecShare\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"X\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"Y\",\"type\":\"uint256\"}],\"internalType\":\"structVerification.G1Point[]\",\"name\":\"a1\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"X\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"Y\",\"type\":\"uint256\"}],\"internalType\":\"structVerification.G1Point[]\",\"name\":\"a2\",\"type\":\"tuple[]\"},{\"internalType\":\"uint256[]\",\"name\":\"challenge\",\"type\":\"uint256[]\"},{\"internalType\":\"uint256[]\",\"name\":\"z\",\"type\":\"uint256[]\"},{\"internalType\":\"uint256\",\"name\":\"threshold\",\"type\":\"uint256\"}],\"name\":\"PVerify\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[]\",\"name\":\"indexTallier\",\"type\":\"uint256[]\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"X\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"Y\",\"type\":\"uint256\"}],\"internalType\":\"structVerification.G1Point[]\",\"name\":\"DecShare\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"X\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"Y\",\"type\":\"uint256\"}],\"internalType\":\"structVerification.G1Point[]\",\"name\":\"a1\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"X\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"Y\",\"type\":\"uint256\"}],\"internalType\":\"structVerification.G1Point[]\",\"name\":\"a2\",\"type\":\"tuple[]\"},{\"internalType\":\"uint256[]\",\"name\":\"challenge\",\"type\":\"uint256[]\"},{\"internalType\":\"uint256[]\",\"name\":\"z\",\"type\":\"uint256[]\"},{\"internalType\":\"uint256\",\"name\":\"threshold\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"numCandidates\",\"type\":\"uint256\"}],\"name\":\"PVerifyTally\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"X\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"Y\",\"type\":\"uint256\"}],\"internalType\":\"structVerification.G1Point[]\",\"name\":\"v\",\"type\":\"tuple[]\"}],\"name\":\"RScodeVerify\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"X\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"Y\",\"type\":\"uint256\"}],\"internalType\":\"structVerification.G1Point[]\",\"name\":\"c\",\"type\":\"tuple[]\"}],\"name\":\"TestAggregateCSet\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"X\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"Y\",\"type\":\"uint256\"}],\"internalType\":\"structVerification.G1Point[]\",\"name\":\"Ej\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"X\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"Y\",\"type\":\"uint256\"}],\"internalType\":\"structVerification.G1Point[]\",\"name\":\"Fj1\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"X\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"Y\",\"type\":\"uint256\"}],\"internalType\":\"structVerification.G1Point[]\",\"name\":\"Fj2\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"X\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"Y\",\"type\":\"uint256\"}],\"internalType\":\"structVerification.G1Point[]\",\"name\":\"_Uj\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"X\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"Y\",\"type\":\"uint256\"}],\"internalType\":\"structVerification.G1Point[]\",\"name\":\"_Cj\",\"type\":\"tuple[]\"},{\"internalType\":\"uint256[]\",\"name\":\"c\",\"type\":\"uint256[]\"},{\"internalType\":\"uint256[]\",\"name\":\"z1\",\"type\":\"uint256[]\"},{\"internalType\":\"uint256[]\",\"name\":\"z2\",\"type\":\"uint256[]\"},{\"internalType\":\"uint256[]\",\"name\":\"z3\",\"type\":\"uint256[]\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"X\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"Y\",\"type\":\"uint256\"}],\"internalType\":\"structVerification.G1Point[]\",\"name\":\"Uj\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"X\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"Y\",\"type\":\"uint256\"}],\"internalType\":\"structVerification.G1Point[]\",\"name\":\"v\",\"type\":\"tuple[]\"},{\"internalType\":\"uint256\",\"name\":\"threshold\",\"type\":\"uint256\"}],\"name\":\"UploadBallotCipher\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"X\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"Y\",\"type\":\"uint256\"}],\"internalType\":\"structVerification.G1Point[]\",\"name\":\"Ej\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"X\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"Y\",\"type\":\"uint256\"}],\"internalType\":\"structVerification.G1Point[]\",\"name\":\"Fj1\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"X\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"Y\",\"type\":\"uint256\"}],\"internalType\":\"structVerification.G1Point[]\",\"name\":\"Fj2\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"X\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"Y\",\"type\":\"uint256\"}],\"internalType\":\"structVerification.G1Point[]\",\"name\":\"_Uj\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"X\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"Y\",\"type\":\"uint256\"}],\"internalType\":\"structVerification.G1Point[]\",\"name\":\"_Cj\",\"type\":\"tuple[]\"},{\"internalType\":\"uint256[]\",\"name\":\"c\",\"type\":\"uint256[]\"},{\"internalType\":\"uint256[]\",\"name\":\"z1\",\"type\":\"uint256[]\"},{\"internalType\":\"uint256[]\",\"name\":\"z2\",\"type\":\"uint256[]\"},{\"internalType\":\"uint256[]\",\"name\":\"z3\",\"type\":\"uint256[]\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"X\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"Y\",\"type\":\"uint256\"}],\"internalType\":\"structVerification.G1Point[]\",\"name\":\"Uj\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"X\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"Y\",\"type\":\"uint256\"}],\"internalType\":\"structVerification.G1Point[]\",\"name\":\"v\",\"type\":\"tuple[]\"},{\"internalType\":\"uint256\",\"name\":\"threshold\",\"type\":\"uint256\"}],\"name\":\"UploadBallotCipher1\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"X\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"Y\",\"type\":\"uint256\"}],\"internalType\":\"structVerification.G1Point[]\",\"name\":\"v\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"X\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"Y\",\"type\":\"uint256\"}],\"internalType\":\"structVerification.G1Point[]\",\"name\":\"c\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"X\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"Y\",\"type\":\"uint256\"}],\"internalType\":\"structVerification.G1Point[]\",\"name\":\"a1\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"X\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"Y\",\"type\":\"uint256\"}],\"internalType\":\"structVerification.G1Point[]\",\"name\":\"a2\",\"type\":\"tuple[]\"},{\"internalType\":\"uint256[]\",\"name\":\"challenge\",\"type\":\"uint256[]\"},{\"internalType\":\"uint256[]\",\"name\":\"z\",\"type\":\"uint256[]\"}],\"name\":\"UploadPVSSShares\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"X\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"Y\",\"type\":\"uint256\"}],\"internalType\":\"structVerification.G1Point[]\",\"name\":\"v\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"X\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"Y\",\"type\":\"uint256\"}],\"internalType\":\"structVerification.G1Point[]\",\"name\":\"c\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"X\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"Y\",\"type\":\"uint256\"}],\"internalType\":\"structVerification.G1Point[]\",\"name\":\"a1\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"X\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"Y\",\"type\":\"uint256\"}],\"internalType\":\"structVerification.G1Point[]\",\"name\":\"a2\",\"type\":\"tuple[]\"},{\"internalType\":\"uint256[]\",\"name\":\"challenge\",\"type\":\"uint256[]\"},{\"internalType\":\"uint256[]\",\"name\":\"z\",\"type\":\"uint256[]\"}],\"name\":\"UploadPVSSShares1\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"X\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"Y\",\"type\":\"uint256\"}],\"internalType\":\"structVerification.G1Point\",\"name\":\"g0\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"X\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"Y\",\"type\":\"uint256\"}],\"internalType\":\"structVerification.G1Point\",\"name\":\"h0\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"uint256[2]\",\"name\":\"X\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256[2]\",\"name\":\"Y\",\"type\":\"uint256[2]\"}],\"internalType\":\"structVerification.G2Point\",\"name\":\"g1\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"uint256[2]\",\"name\":\"X\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256[2]\",\"name\":\"Y\",\"type\":\"uint256[2]\"}],\"internalType\":\"structVerification.G2Point\",\"name\":\"pkI\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"X\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"Y\",\"type\":\"uint256\"}],\"internalType\":\"structVerification.G1Point[]\",\"name\":\"sigmak\",\"type\":\"tuple[]\"},{\"internalType\":\"uint256\",\"name\":\"a\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"b\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"numCandidates\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"numTalliers\",\"type\":\"uint256\"}],\"name\":\"UploadParameters\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"X\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"Y\",\"type\":\"uint256\"}],\"internalType\":\"structVerification.G1Point[]\",\"name\":\"pks\",\"type\":\"tuple[]\"}],\"name\":\"UploadPublicKey\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"X\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"Y\",\"type\":\"uint256\"}],\"internalType\":\"structVerification.G1Point[]\",\"name\":\"Ej\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"X\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"Y\",\"type\":\"uint256\"}],\"internalType\":\"structVerification.G1Point[]\",\"name\":\"Fj1\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"X\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"Y\",\"type\":\"uint256\"}],\"internalType\":\"structVerification.G1Point[]\",\"name\":\"Fj2\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"X\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"Y\",\"type\":\"uint256\"}],\"internalType\":\"structVerification.G1Point[]\",\"name\":\"_Uj\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"X\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"Y\",\"type\":\"uint256\"}],\"internalType\":\"structVerification.G1Point[]\",\"name\":\"_Cj\",\"type\":\"tuple[]\"},{\"internalType\":\"uint256[]\",\"name\":\"c\",\"type\":\"uint256[]\"},{\"internalType\":\"uint256[]\",\"name\":\"z1\",\"type\":\"uint256[]\"},{\"internalType\":\"uint256[]\",\"name\":\"z2\",\"type\":\"uint256[]\"},{\"internalType\":\"uint256[]\",\"name\":\"z3\",\"type\":\"uint256[]\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"X\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"Y\",\"type\":\"uint256\"}],\"internalType\":\"structVerification.G1Point[]\",\"name\":\"Uj\",\"type\":\"tuple[]\"},{\"internalType\":\"uint256[]\",\"name\":\"d\",\"type\":\"uint256[]\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"X\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"Y\",\"type\":\"uint256\"}],\"internalType\":\"structVerification.G1Point[]\",\"name\":\"v\",\"type\":\"tuple[]\"},{\"internalType\":\"uint256[]\",\"name\":\"indices\",\"type\":\"uint256[]\"},{\"internalType\":\"uint256\",\"name\":\"threshold\",\"type\":\"uint256\"}],\"name\":\"ZKRPVerify\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"int256\",\"name\":\"i\",\"type\":\"int256\"},{\"internalType\":\"uint256\",\"name\":\"n\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"l\",\"type\":\"uint256\"}],\"name\":\"coefficient\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"a\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"prime\",\"type\":\"uint256\"}],\"name\":\"inv\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60808060405234601557613448908161001b8239f35b600080fdfe6080604052600436101561001257600080fd5b60003560e01c8063338255f31461016757806343ac6ba414610162578063446bf71c1461015d5780634c4791ec146101585780634dea579f146101535780634ea9637d1461014e5780635b4af2b8146101495780635e08ddbf1461014457806365d285541461013f57806375457b931461013a5780637e3d3acb146101355780639e134b5914610130578063b52c5ff21461012b578063b562c4ab14610126578063beaed4a014610121578063c3ceb16b1461011c578063db301a7f14610117578063e17afe5714610112578063e283f1301461010d578063ea315a94146101085763f0934fe91461010357600080fd5b611a98565b611a7b565b611a5e565b6118b7565b611749565b6116bb565b611655565b611498565b61132a565b61128c565b611173565b610994565b610932565b6108dd565b6108b1565b610859565b61078a565b61062d565b6105b1565b610589565b346101f55760403660031901126101f557600435602435906001198201918083116101f05760c091602091600094604051926101a284610230565b84845284808501528460408501526060840152608083015260a08201526040519283916101cf848461024b565b83368437856005600019f1156101ec575160405190815260209150f35b5080fd5b611b26565b600080fd5b634e487b7160e01b600052604160045260246000fd5b604081019081106001600160401b0382111761022b57604052565b6101fa565b60c081019081106001600160401b0382111761022b57604052565b90601f801991011681019081106001600160401b0382111761022b57604052565b6040519061027b60408361024b565b565b6001600160401b03811161022b5760051b60200190565b60409060031901126101f557604051906102ad82610210565b60043582526024356020830152565b60409060431901126101f557604051906102d582610210565b60443582526064356020830152565b6040906101031901126101f557604051906102fe82610210565b610104358252610124356020830152565b6040906101431901126101f5576040519061032982610210565b610144358252610164356020830152565b81601f820112156101f5578035906103518261027d565b9261035f604051948561024b565b82845260208085019360061b830101918183116101f557602001925b828410610389575050505090565b60006040858403126103bd5750602060409182516103a681610210565b86358152828701358382015281520193019261037b565b80fd5b9080601f830112156101f55781356103d78161027d565b926103e5604051948561024b565b81845260208085019260051b8201019283116101f557602001905b82821061040d5750505090565b8135815260209182019101610400565b6101806003198201126101f5576004356001600160401b0381116101f557816104489160040161033a565b916024356001600160401b0381116101f557826104679160040161033a565b916044356001600160401b0381116101f557816104869160040161033a565b916064356001600160401b0381116101f557826104a59160040161033a565b916084356001600160401b0381116101f557816104c49160040161033a565b9160a4356001600160401b0381116101f557826104e3916004016103c0565b9160c4356001600160401b0381116101f55781610502916004016103c0565b9160e4356001600160401b0381116101f55782610521916004016103c0565b91610104356001600160401b0381116101f55781610541916004016103c0565b91610124356001600160401b0381116101f557826105619160040161033a565b9161014435906001600160401b0382116101f5576105819160040161033a565b906101643590565b346101f5576105af61059a3661041d565b9a999099989198979297969396959495611d6e565b005b346101f55760803660031901126101f5576004356024356001600160401b0381116101f5576105e490369060040161033a565b90604435906001600160401b0382116101f55760409261060b6106159336906004016103c0565b9060643592611e8d565b61062b8251809260208091805184520151910152565bf35b346101f55760e03660031901126101f5576004356001600160401b0381116101f55761065d9036906004016103c0565b6024356001600160401b0381116101f55761067c90369060040161033a565b906044356001600160401b0381116101f55761069c90369060040161033a565b916064356001600160401b0381116101f5576106bc90369060040161033a565b6084356001600160401b0381116101f5576106db9036906004016103c0565b9060a435946001600160401b0386116101f5576106ff6105af9636906004016103c0565b9360c43595612077565b9080601f830112156101f5576040805192610724828561024b565b839181019283116101f557905b82821061073e5750505090565b8135815260209182019101610731565b906080610103198301126101f55760405161076881610210565b6020610785829461077b81610104610709565b8452610144610709565b910152565b346101f5576102203660031901126101f5576107a536610294565b6107ae366102bc565b9060803660831901126101f557604051916107c883610210565b6107d3366084610709565b83526107e03660c4610709565b60208401526107ee3661074e565b9261018435936001600160401b0385116101f5576108136105af95369060040161033a565b6101a435916101c435936101e43595610204359761229e565b60206003198201126101f557600435906001600160401b0382116101f5576108569160040161033a565b90565b346101f5576108673661082c565b60005b81518110156105af5761087d8183611c45565b5190601154600160401b81101561022b576001926108a582856108ab9401601155601161205b565b906121de565b0161086a565b346101f55760603660031901126101f55760206108d56044356024356004356123b8565b604051908152f35b346101f5576108eb3661082c565b60005b81518110156105af578061092c610923610909600194611ffe565b5061091e6109178588611c45565b5191611e6f565b612f86565b6108a583611ffe565b016108ee565b346101f5576105af6109433661041d565b9a99909998919897929796939695949561244f565b602060408183019282815284518094520192019060005b81811061097c5750505090565b8251151584526020938401939092019160010161096f565b346101f55760003660031901126101f557604051601a548082529060208101601a6000527f057c384a7d1c54f3a1b2e5e67b2617b8224fdfd1ea7234eea573a6ff665ff63e926000935b81601f860110610ed85791610b07948492610afb945491818110610ebc575b818110610e9d575b818110610e7e575b818110610e5f575b818110610e41575b818110610e22575b818110610e03575b818110610de4575b818110610dc5575b818110610da6575b818110610d87575b818110610d68575b818110610d49575b818110610d2a575b818110610d0b575b818110610cec575b818110610ccd575b818110610cae575b818110610c8f575b818110610c70575b818110610c51575b818110610c32575b818110610c13575b818110610bf4575b818110610bd5575b818110610bb6575b818110610b97575b818110610b78575b818110610b59575b818110610b3a575b818110610b1b575b10610b0b575b50038261024b565b60405191829182610958565b0390f35b60f81c1515815260200138610af3565b92602081610b3260019360ff8760f01c1615159052565b019301610aed565b92602081610b5160019360ff8760e81c1615159052565b019301610ae5565b92602081610b7060019360ff8760e01c1615159052565b019301610add565b92602081610b8f60019360ff8760d81c1615159052565b019301610ad5565b92602081610bae60019360ff8760d01c1615159052565b019301610acd565b92602081610bcd60019360ff8760c81c1615159052565b019301610ac5565b92602081610bec60019360ff8760c01c1615159052565b019301610abd565b92602081610c0b60019360ff8760b81c1615159052565b019301610ab5565b92602081610c2a60019360ff8760b01c1615159052565b019301610aad565b92602081610c4960019360ff8760a81c1615159052565b019301610aa5565b92602081610c6860019360ff8760a01c1615159052565b019301610a9d565b92602081610c8760019360ff8760981c1615159052565b019301610a95565b92602081610ca660019360ff8760901c1615159052565b019301610a8d565b92602081610cc560019360ff8760881c1615159052565b019301610a85565b92602081610ce460019360ff8760801c1615159052565b019301610a7d565b92602081610d0360019360ff8760781c1615159052565b019301610a75565b92602081610d2260019360ff8760701c1615159052565b019301610a6d565b92602081610d4160019360ff8760681c1615159052565b019301610a65565b92602081610d6060019360ff8760601c1615159052565b019301610a5d565b92602081610d7f60019360ff8760581c1615159052565b019301610a55565b92602081610d9e60019360ff8760501c1615159052565b019301610a4d565b92602081610dbd60019360ff8760481c1615159052565b019301610a45565b92602081610ddc60019360ff8760401c1615159052565b019301610a3d565b92602081610dfb60019360ff8760381c1615159052565b019301610a35565b92602081610e1a60019360ff8760301c1615159052565b019301610a2d565b92602081610e3960019360ff8760281c1615159052565b019301610a25565b92602081610e5760019360ff87851c1615159052565b019301610a1d565b92602081610e7660019360ff8760181c1615159052565b019301610a15565b92602081610e9560019360ff8760101c1615159052565b019301610a0d565b92602081610eb460019360ff8760081c1615159052565b019301610a05565b92602081610ed060019360ff871615159052565b0193016109fd565b9160016104006020926111688654610ef48360ff831615159052565b610f0786840160ff8360081c1615159052565b610f1b6040840160ff8360101c1615159052565b610f2f6060840160ff8360181c1615159052565b80861c60ff1615156080840152610f5060a0840160ff8360281c1615159052565b610f6460c0840160ff8360301c1615159052565b610f7860e0840160ff8360381c1615159052565b610f8d610100840160ff8360401c1615159052565b610fa2610120840160ff8360481c1615159052565b610fb7610140840160ff8360501c1615159052565b610fcc610160840160ff8360581c1615159052565b610fe1610180840160ff8360601c1615159052565b610ff66101a0840160ff8360681c1615159052565b61100b6101c0840160ff8360701c1615159052565b6110206101e0840160ff8360781c1615159052565b611035610200840160ff8360801c1615159052565b61104a610220840160ff8360881c1615159052565b61105f610240840160ff8360901c1615159052565b611074610260840160ff8360981c1615159052565b611089610280840160ff8360a01c1615159052565b61109e6102a0840160ff8360a81c1615159052565b6110b36102c0840160ff8360b01c1615159052565b6110c86102e0840160ff8360b81c1615159052565b6110dd610300840160ff8360c01c1615159052565b6110f2610320840160ff8360c81c1615159052565b611107610340840160ff8360d01c1615159052565b61111c610360840160ff8360d81c1615159052565b611131610380840160ff8360e01c1615159052565b6111466103a0840160ff8360e81c1615159052565b61115b6103c0840160ff8360f01c1615159052565b60f81c15156103e0830152565b0193019401936109de565b346101f5576101203660031901126101f55761118e36610294565b6044356001600160401b0381116101f5576111ad90369060040161033a565b906064356001600160401b0381116101f5576111cd90369060040161033a565b906084356001600160401b0381116101f5576111ed90369060040161033a565b60a4356001600160401b0381116101f55761120c90369060040161033a565b60c4356001600160401b0381116101f55761122b90369060040161033a565b9060e4356001600160401b0381116101f55761124b9036906004016103c0565b9261010435956001600160401b0387116101f557610b079761127461127a9836906004016103c0565b96612558565b60405190151581529081906020820190565b6101c03660031901126101f5576112a236610294565b6112ab366102bc565b9060403660831901126101f557604051906112c582610210565b608435825260a435602083015260403660c31901126101f557602092611320926040516112f181610210565b60c435815260e43586820152611306366102e4565b906113103661030f565b9261018435946101a4359661260c565b6040519015158152f35b346101f55760003660031901126101f557604051601c548082529060208101601c6000527f0e4562a10381dec21b205ed72637e6b1b523bdd0e4d4d50af5cd23dd4500a211926000935b81601f8601106114715791610b07948492610afb945491818110610ebc57818110610e9d57818110610e7e57818110610e5f57818110610e4157818110610e2257818110610e0357818110610de457818110610dc557818110610da657818110610d8757818110610d6857818110610d4957818110610d2a57818110610d0b57818110610cec57818110610ccd57818110610cae57818110610c8f57818110610c7057818110610c5157818110610c3257818110610c1357818110610bf457818110610bd557818110610bb657818110610b9757818110610b7857818110610b5957818110610b3a57818110610b1b5710610b0b5750038261024b565b91600161040060209261148d8654610ef48360ff831615159052565b019301940193611374565b346101f5576101c03660031901126101f5576004356001600160401b0381116101f5576114c990369060040161033a565b6024356001600160401b0381116101f5576114e890369060040161033a565b6044356001600160401b0381116101f55761150790369060040161033a565b6064356001600160401b0381116101f55761152690369060040161033a565b6084356001600160401b0381116101f55761154590369060040161033a565b60a4356001600160401b0381116101f5576115649036906004016103c0565b60c4356001600160401b0381116101f5576115839036906004016103c0565b60e4356001600160401b0381116101f5576115a29036906004016103c0565b90610104356001600160401b0381116101f5576115c39036906004016103c0565b92610124356001600160401b0381116101f5576115e490369060040161033a565b94610144356001600160401b0381116101f5576116059036906004016103c0565b96610164356001600160401b0381116101f55761162690369060040161033a565b98610184359a6001600160401b038c116101f557610b079c61164f61127a9d36906004016103c0565b50612730565b346101f55760206113206116683661082c565b612970565b602060408183019282815284518094520192019060005b8181106116915750505090565b90919260206040826116b0600194885160208091805184520151910152565b019401929101611684565b346101f55760003660031901126101f5576015546116d88161027d565b906116e6604051928361024b565b808252601560009081527f55f448fdea98c4d29eb340757ef0a66cd03dbb9538908a6a81d96026b71ec475602084015b83831061172b5760405180610b07878261166d565b6002602060019261173b85611e6f565b815201920192019190611716565b346101f55760003660031901126101f557604051601b548082529060208101601b6000527f3ad8aa4f87544323a9d1e5dd902f40c356527a7955687113db5f9a85ad579dc1926000935b81601f8601106118905791610b07948492610afb945491818110610ebc57818110610e9d57818110610e7e57818110610e5f57818110610e4157818110610e2257818110610e0357818110610de457818110610dc557818110610da657818110610d8757818110610d6857818110610d4957818110610d2a57818110610d0b57818110610cec57818110610ccd57818110610cae57818110610c8f57818110610c7057818110610c5157818110610c3257818110610c1357818110610bf457818110610bd557818110610bb657818110610b9757818110610b7857818110610b5957818110610b3a57818110610b1b5710610b0b5750038261024b565b9160016104006020926118ac8654610ef48360ff831615159052565b019301940193611793565b346101f5576101003660031901126101f5576004356001600160401b0381116101f5576118e89036906004016103c0565b6024356001600160401b0381116101f55761190790369060040161033a565b906044356001600160401b0381116101f55761192790369060040161033a565b916064356001600160401b0381116101f55761194790369060040161033a565b926084356001600160401b0381116101f5576119679036906004016103c0565b60a435946001600160401b0386116101f55761198a6105af9636906004016103c0565b9260c4359460e43596612b4b565b9060c06003198301126101f5576004356001600160401b0381116101f557826119c39160040161033a565b916024356001600160401b0381116101f557816119e29160040161033a565b916044356001600160401b0381116101f55782611a019160040161033a565b916064356001600160401b0381116101f55781611a209160040161033a565b916084356001600160401b0381116101f55782611a3f916004016103c0565b9160a435906001600160401b0382116101f557610856916004016103c0565b346101f5576105af611a6f36611998565b94939093929192612cdf565b346101f5576105af611a8c36611998565b94939093929192612d74565b346101f55760003660031901126101f557601754611ab58161027d565b90611ac3604051928361024b565b808252601760009081527fc624b66cc0138b8fabc209247f72d758e1cf3343756d543badbf24212bed8c15602084015b838310611b085760405180610b07878261166d565b60026020600192611b1885611e6f565b815201920192019190611af3565b634e487b7160e01b600052601160045260246000fd5b6000198101919082116101f057565b919082039182116101f057565b90611b628261027d565b611b6f604051918261024b565b8281528092611b80601f199161027d565b0190602036910137565b90600182018092116101f057565b90600282018092116101f057565b90600382018092116101f057565b90600482018092116101f057565b90600582018092116101f057565b919082018092116101f057565b634e487b7160e01b600052603260045260246000fd5b805115611c005760200190565b611bdd565b805160011015611c005760400190565b805160021015611c005760600190565b805160031015611c005760800190565b805160041015611c005760a00190565b8051821015611c005760209160051b010190565b601a54600160401b81101561022b576001810180601a55811015611c0057601a60005260206000208160051c019060ff60f883549260031b161b19169055565b601a54600160401b81101561022b576001810180601a55811015611c0057601a60005260206000208160051c019060f882549160031b169060ff6001831b921b1916179055565b601c54600160401b81101561022b576001810180601c55811015611c0057601c60005260206000208160051c019060f882549160031b169060ff6001831b921b1916179055565b601b54600160401b81101561022b576001810180601b55811015611c0057601b60005260206000208160051c019060f882549160031b169060ff6001831b921b1916179055565b9b9a98979695949392919098611d8381611b58565b9060005b818110611de857505050600f5498611d9e8a611b58565b9960005b818110611dca575050611db59b9c612730565b15611dc25761027b611c99565b61027b611c59565b808c611de182611ddb600195612e8b565b92611c45565b5201611da2565b60018101908181116101f057600191611e018286611c45565b5201611d87565b60405190611e1582610210565b60006020838281520152565b634e487b7160e01b600052601260045260246000fd5b60405190611e4482610210565b60005482526001546020830152565b60405190611e6082610210565b60025482526003546020830152565b90604051611e7c81610210565b602060018294805484520154910152565b9092611e97611e08565b50611ea181611b58565b92600080937f30644e72e131a029b85045b68181585d2833e84879b9709143e1f593efffffff915b848110611f2e57505050505050611ee6611ee1611e37565b612ef9565b916000925b8251841015611f2757611f1f600191611f19611f078786611c45565b51611f128888611c45565b5190612f3f565b90612f86565b930192611eeb565b9250505090565b85906001848160005b898110611f78575090611f4991612e24565b92611f73576001926000805160206133f38339815191529109611f6c828a611c45565b5201611ec9565b611e21565b915091929350808403611f94575b600101889392918691611f37565b919097611fab611fa48489611c45565b5186612eab565b90611f73576000805160206133f383398151915291099060016000805160206133f3833981519152611ff2611fe0868a611c45565b51611feb858b611c45565b5190612eab565b60009a09919050611f86565b601554811015611c0057601560005260206000209060011b0190600090565b601154811015611c0057601160005260206000209060011b0190600090565b601754811015611c0057601760005260206000209060011b0190600090565b8054821015611c005760005260206000209060011b0190600090565b939095949160005b868110612090575050505050505050565b61209a8187611c45565b5160001981019081116101f0576120b09061201d565b50906120bc8184611c45565b516120c7828b611c45565b51926120d3838a611c45565b516000198101929083116101f0576001946120f061213694611ffe565b50926120fc868c611c45565b5191612108878b611c45565b5193612114888d611c45565b519561213061212a612124611e37565b95611e6f565b91611e6f565b9361260c565b500161207f565b906006820291808304600614901517156101f057565b9060005b6002811061216457505050565b600190602083519301928185015501612157565b805160005b60028110612196575050602061027b9101516006612153565b6001906020835193019281600401550161217d565b805160005b600281106121c9575050602061027b910151600a612153565b600190602083519301928160080155016121b0565b91906121f4576020816001925184550151910155565b634e487b7160e01b600052600060045260246000fd5b600c5490600160401b82101561022b576108a582600161027b9401600c55600c61205b565b60155490600160401b82101561022b576108a582600161027b9401601555601561205b565b60165490600160401b82101561022b576108a582600161027b9401601655601661205b565b60175490600160401b82101561022b576108a582600161027b9401601755601761205b565b926122dc6122eb96936122d76122e1946122c96122e6989e9c9b999e60209080516000550151600155565b805160025560200151600355565b612178565b6121ab565b600d55565b600e55565b6122f481600f55565b6122fd83601055565b60005b8451811015612325578061231f61231960019388611c45565b5161220a565b01612300565b509190925060005b81811061236857505060005b818110612344575050565b60019061236261235261026c565b6000815260006020820152612279565b01612339565b60019061238661237661026c565b600081526000602082015261222f565b6123a161239161026c565b6000815260006020820152612254565b0161232d565b600160ff1b81146101f05760000390565b600192600019810190811384166101f0576123d2906123a7565b828113156123e05750505090565b8082036123fe575b6001600160ff1b0381146101f0576001016123d2565b9260006000805160206133f3833981519152809261244661241f8887612fd6565b7f30644e72e131a029b85045b68181585d2833e84879b9709143e1f593efffffff90612e24565b900908926123e8565b97959391899c9b9a99979593919a61246681611b58565b60005b82811061252057505050600f549a6124808c611b58565b9960005b8d81106124ed57506124969c50612730565b156124e3576124a3611c99565b60005b82518110156124de57806124d86124cf6124c160019461203c565b5061091e6109178589611c45565b6108a58361203c565b016124a6565b509050565b905061027b611c59565b8b6001929496989a9c9d9e939597999b5061250b82611ddb81612e8b565b5201918e99979593919c9b9a98969492612484565b60019192939597999b9c9d9496989a5061253981611b8a565b6125438285611c45565b5201908e999795939c9b9a9896949291612469565b9491969390959261256887612970565b612579575050505050505050600090565b60005b82518110156125fe576125e38787878787878f8f896125d5816125ce816125c7816125c0816125b9816125b2816125dc9d611c45565b519c611c45565b519b611c45565b519a611c45565b5199611c45565b5198611c45565b5197611c45565b519661260c565b156125f05760010161257c565b505050505050505050600090565b505050505050505050600190565b866126398961263361264b989a61262d612645989d61263f979a989a612f3f565b98612f3f565b94612f3f565b97612f3f565b92612f86565b93612f86565b9183518151149384159461269b575b5050821561268e575b821561267a575b505061267557600190565b600090565b60209192508101519101511415388061266a565b8051825114159250612663565b60209081015191015114159250388061265a565b60405191906000835b600282106126ce5750505061027b60408361024b565b60016020819285548152019301910190916126b8565b604051906126f182610210565b816126fc60046126af565b8152602061078560066126af565b6040519061271782610210565b8161272260086126af565b81526020610785600a6126af565b96999a98919597929390949860005b8c518110156129475761277f6127648c611f128461275d8188611c45565b5192611c45565b611f196127718489611c45565b5161277a611e53565b612f3f565b6127898285611c45565b51518151149081159161292a575b50156127b25750505050505050505050505050612675611c59565b6127eb8a611f19612771848f8a611f196127e2846127dc6127dc958f8361275d81611f1293611c45565b93611c45565b5161277a611e37565b6127f58288611c45565b51518151149081159161290d575b501561281e5750505050505050505050505050612675611c59565b80808d818d8f828f916128318284611c45565b519161283c91611c45565b5161284691612f3f565b9461285091611c45565b519161285b91611c45565b5161286591612f3f565b61286e90613029565b9261287891611c45565b51612881611e37565b9061288b91612f3f565b612895848b611c45565b5161289f90613029565b926128aa858d611c45565b516128b490613029565b906128bd6126e4565b936128c66126e4565b906128cf61270a565b916128d86126e4565b946128e16126e4565b976128eb99613165565b156128f85760010161273f565b50505050505050505050505050612675611c59565b905060208061291c848a611c45565b510151910151141538612803565b90506020806129398487611c45565b510151910151141538612797565b5050505050505050505050505061295c611c99565b600190565b60001981146101f05760010190565b60105491600f54926129828351611b58565b9160015b82811115612a475750815b61299b8684611bd0565b8110156129e457806000805160206133f383398151915260006129d289876129cd6129c882600199611b4b565b6123a7565b6123b8565b086129dd8287611c45565b5201612991565b509350506129f0611e53565b916000925b8151841015612a1957612a11600191611f19611f078786611c45565b9301926129f5565b9250505080516002541490811591612a36575b5061267557600190565b602091500151600354141538612a2c565b806000805160206133f38339815191526000612a678987612a80966123b8565b08612a7a612a7483611b3c565b87611c45565b52612961565b612986565b60115490612a928261027d565b91612aa0604051938461024b565b808352601160009081527f31ecc21a745e3968a04e9570e4425bc18fa8019c68028196b546d1669c200c68602085015b838310612add5750505050565b60026020600192612aed85611e6f565b815201920192019190612ad0565b90612b058261027d565b612b12604051918261024b565b8281528092612b23601f199161027d565b019060005b828110612b3457505050565b602090612b3f611e08565b82828501015201612b28565b909295949395612b5a86611b58565b94612b6487612afb565b9760005b888110612beb5750505050505050612b7f84612afb565b9260005b858110612b9257505050505050565b80612ba9858585612ba4600196612e8b565b611e8d565b612bb38288611c45565b52612bbe8187611c45565b50612be56124cf612bce8361203c565b5061091e61212a612bdf868c611c45565b51613029565b01612b83565b612bf58186611c45565b51612bff90611b3c565b612c089061201d565b50612c138288611c45565b51612c1e838a611c45565b51612c298489611c45565b51612c3390611b3c565b612c3c90611ffe565b50612c478589611c45565b5190612c538688611c45565b5192612c5f878a611c45565b5194612c69611e37565b96612c7390611e6f565b92612c7d90611e6f565b92612c879761260c565b151560011490600191612c9b575b01612b68565b612ca58189611c45565b51612cb0828d611c45565b52612cbb818c611c45565b50612cc68187611c45565b51612cd1828b611c45565b52612cda611ce0565b612c95565b908096959391612d019593612cf2611e53565b91612cfb612a85565b92612558565b15612d3257612d0e611d27565b60005b82518110156124de5780612d2c6109236124c1600194611ffe565b01612d11565b9050601b54600160401b81101561022b576001810180601b55811015611c0057601b60005260206000208160051c019060ff60f883549260031b161b19169055565b612d8395949392612cf2611e53565b15612dcf57601d54600160401b81101561022b576001810180601d55811015611c0057601d60005260206000208160051c019060f882549160031b169060ff6001831b921b1916179055565b601d54600160401b81101561022b576001810180601d55811015611c0057601d60005260206000208160051c019060ff60f883549260031b161b19169055565b604051906020612e1f818461024b565b368337565b60c09160405191612e3483610230565b6020835260208084015260206040840152606083015260808201526000805160206133f383398151915260a0820152602090604051928391612e76848461024b565b8336843760006005600019f1156101f5575190565b6000805160206133f38339815191529081038181116101f0576000900890565b81811115612ed457905b81039081116101f05760006000805160206133f3833981519152910890565b6000805160206133f383398151915281018091116101f05790612eb5565b156101f557565b90612f02611e08565b91826080606092602060405191612f19868461024b565b8536843780518352015160208201526000604082015260076107cf195a01fa156101f557565b9060809291612f4c611e08565b938491606093602060405192612f62878561024b565b863685378051845201516020830152604082015260076107cf195a01fa156101f557565b6020929160c0606092612f97611e08565b9586938160405193612faa60808661024b565b6080368637805185520151828401528051604084015201518482015260066107cf195a01fa156101f557565b600082820392128183128116918313901516176101f05760008112612ff85790565b613001906123a7565b6000805160206133f3833981519152036000805160206133f383398151915281116101f05790565b613031611e08565b50805115806130e3575b6130c8577f30644e72e131a029b85045b68181585d97816a916871ca8d3c208c16d87cfd4760208251920151067f30644e72e131a029b85045b68181585d97816a916871ca8d3c208c16d87cfd47037f30644e72e131a029b85045b68181585d97816a916871ca8d3c208c16d87cfd4781116101f057604051916130be83610210565b8252602082015290565b506040516130d581610210565b600081526000602082015290565b5060208101511561303b565b604051906130fe60c08361024b565b6005825281601f19613110600561027d565b019060409060005b8381106131255750505050565b602090835161313381610210565b845161313f868261024b565b8536823781528451613151868261024b565b853682378382015282828501015201613118565b97909199989596949396613179600561027d565b97613187604051998a61024b565b60058952601f19613198600561027d565b0160005b8181106132805750506108569a9b6131b26130ef565b9a6131bc8b611bf3565b526131c68a611bf3565b506131d08a611c05565b526131da89611c05565b506131e489611c15565b526131ee88611c15565b506131f888611c25565b5261320287611c25565b5061320c87611c35565b5261321686611c35565b5061322087611bf3565b5261322a86611bf3565b5061323486611c05565b5261323e85611c05565b5061324885611c15565b5261325284611c15565b5061325c84611c25565b5261326683611c25565b5061327083611c35565b5261327a82611c35565b50613298565b808b6020809361328e611e08565b920101520161319c565b80516132a683518214612ef2565b6132af8161213d565b926132b984611b58565b9260005b8381106132f45750505050602080926132ee926132d8612e0f565b94859260051b910160086107cf195a01fa612ef2565b51151590565b8061330160019284611c45565b515161330c8261213d565b9061331960009289611c45565b5260206133268386611c45565b51015161334361333d6133388561213d565b611b8a565b89611c45565b5261334e8286611c45565b51515161336561333d6133608561213d565b611b98565b5261337b6133738387611c45565b515160200190565b5161339061333d61338b8561213d565b611ba6565b52602061339d8387611c45565b5101519050516133bd6133b76133b28461213d565b611bb4565b88611c45565b526133d660206133cd8387611c45565b51015160200190565b516133eb6133b76133e68461213d565b611bc2565b52016132bd56fe30644e72e131a029b85045b68181585d2833e84879b9709143e1f593f0000001a26469706673582212204c915711e26609a2d1938401223599fac92dc6cfc0238e0b89dc84cb58ba92bd64736f6c634300081c0033",
}

// ContractABI is the input ABI used to generate the binding from.
// Deprecated: Use ContractMetaData.ABI instead.
var ContractABI = ContractMetaData.ABI

// ContractBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use ContractMetaData.Bin instead.
var ContractBin = ContractMetaData.Bin

// DeployContract deploys a new Ethereum contract, binding an instance of Contract to it.
func DeployContract(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *Contract, error) {
	parsed, err := ContractMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(ContractBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Contract{ContractCaller: ContractCaller{contract: contract}, ContractTransactor: ContractTransactor{contract: contract}, ContractFilterer: ContractFilterer{contract: contract}}, nil
}

// Contract is an auto generated Go binding around an Ethereum contract.
type Contract struct {
	ContractCaller     // Read-only binding to the contract
	ContractTransactor // Write-only binding to the contract
	ContractFilterer   // Log filterer for contract events
}

// ContractCaller is an auto generated read-only Go binding around an Ethereum contract.
type ContractCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ContractTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ContractFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ContractSession struct {
	Contract     *Contract         // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ContractCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ContractCallerSession struct {
	Contract *ContractCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts   // Call options to use throughout this session
}

// ContractTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ContractTransactorSession struct {
	Contract     *ContractTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// ContractRaw is an auto generated low-level Go binding around an Ethereum contract.
type ContractRaw struct {
	Contract *Contract // Generic contract binding to access the raw methods on
}

// ContractCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ContractCallerRaw struct {
	Contract *ContractCaller // Generic read-only contract binding to access the raw methods on
}

// ContractTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ContractTransactorRaw struct {
	Contract *ContractTransactor // Generic write-only contract binding to access the raw methods on
}

// NewContract creates a new instance of Contract, bound to a specific deployed contract.
func NewContract(address common.Address, backend bind.ContractBackend) (*Contract, error) {
	contract, err := bindContract(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Contract{ContractCaller: ContractCaller{contract: contract}, ContractTransactor: ContractTransactor{contract: contract}, ContractFilterer: ContractFilterer{contract: contract}}, nil
}

// NewContractCaller creates a new read-only instance of Contract, bound to a specific deployed contract.
func NewContractCaller(address common.Address, caller bind.ContractCaller) (*ContractCaller, error) {
	contract, err := bindContract(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ContractCaller{contract: contract}, nil
}

// NewContractTransactor creates a new write-only instance of Contract, bound to a specific deployed contract.
func NewContractTransactor(address common.Address, transactor bind.ContractTransactor) (*ContractTransactor, error) {
	contract, err := bindContract(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ContractTransactor{contract: contract}, nil
}

// NewContractFilterer creates a new log filterer instance of Contract, bound to a specific deployed contract.
func NewContractFilterer(address common.Address, filterer bind.ContractFilterer) (*ContractFilterer, error) {
	contract, err := bindContract(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ContractFilterer{contract: contract}, nil
}

// bindContract binds a generic wrapper to an already deployed contract.
func bindContract(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ContractMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Contract *ContractRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Contract.Contract.ContractCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Contract *ContractRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Contract.Contract.ContractTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Contract *ContractRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Contract.Contract.ContractTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Contract *ContractCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Contract.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Contract *ContractTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Contract.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Contract *ContractTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Contract.Contract.contract.Transact(opts, method, params...)
}

// GetAggregateValue is a free data retrieval call binding the contract method 0xc3ceb16b.
//
// Solidity: function GetAggregateValue() view returns((uint256,uint256)[])
func (_Contract *ContractCaller) GetAggregateValue(opts *bind.CallOpts) ([]VerificationG1Point, error) {
	var out []interface{}
	err := _Contract.contract.Call(opts, &out, "GetAggregateValue")

	if err != nil {
		return *new([]VerificationG1Point), err
	}

	out0 := *abi.ConvertType(out[0], new([]VerificationG1Point)).(*[]VerificationG1Point)

	return out0, err

}

// GetAggregateValue is a free data retrieval call binding the contract method 0xc3ceb16b.
//
// Solidity: function GetAggregateValue() view returns((uint256,uint256)[])
func (_Contract *ContractSession) GetAggregateValue() ([]VerificationG1Point, error) {
	return _Contract.Contract.GetAggregateValue(&_Contract.CallOpts)
}

// GetAggregateValue is a free data retrieval call binding the contract method 0xc3ceb16b.
//
// Solidity: function GetAggregateValue() view returns((uint256,uint256)[])
func (_Contract *ContractCallerSession) GetAggregateValue() ([]VerificationG1Point, error) {
	return _Contract.Contract.GetAggregateValue(&_Contract.CallOpts)
}

// GetDVerifyResult is a free data retrieval call binding the contract method 0xdb301a7f.
//
// Solidity: function GetDVerifyResult() view returns(bool[])
func (_Contract *ContractCaller) GetDVerifyResult(opts *bind.CallOpts) ([]bool, error) {
	var out []interface{}
	err := _Contract.contract.Call(opts, &out, "GetDVerifyResult")

	if err != nil {
		return *new([]bool), err
	}

	out0 := *abi.ConvertType(out[0], new([]bool)).(*[]bool)

	return out0, err

}

// GetDVerifyResult is a free data retrieval call binding the contract method 0xdb301a7f.
//
// Solidity: function GetDVerifyResult() view returns(bool[])
func (_Contract *ContractSession) GetDVerifyResult() ([]bool, error) {
	return _Contract.Contract.GetDVerifyResult(&_Contract.CallOpts)
}

// GetDVerifyResult is a free data retrieval call binding the contract method 0xdb301a7f.
//
// Solidity: function GetDVerifyResult() view returns(bool[])
func (_Contract *ContractCallerSession) GetDVerifyResult() ([]bool, error) {
	return _Contract.Contract.GetDVerifyResult(&_Contract.CallOpts)
}

// GetPVerifyResult is a free data retrieval call binding the contract method 0xb52c5ff2.
//
// Solidity: function GetPVerifyResult() view returns(bool[])
func (_Contract *ContractCaller) GetPVerifyResult(opts *bind.CallOpts) ([]bool, error) {
	var out []interface{}
	err := _Contract.contract.Call(opts, &out, "GetPVerifyResult")

	if err != nil {
		return *new([]bool), err
	}

	out0 := *abi.ConvertType(out[0], new([]bool)).(*[]bool)

	return out0, err

}

// GetPVerifyResult is a free data retrieval call binding the contract method 0xb52c5ff2.
//
// Solidity: function GetPVerifyResult() view returns(bool[])
func (_Contract *ContractSession) GetPVerifyResult() ([]bool, error) {
	return _Contract.Contract.GetPVerifyResult(&_Contract.CallOpts)
}

// GetPVerifyResult is a free data retrieval call binding the contract method 0xb52c5ff2.
//
// Solidity: function GetPVerifyResult() view returns(bool[])
func (_Contract *ContractCallerSession) GetPVerifyResult() ([]bool, error) {
	return _Contract.Contract.GetPVerifyResult(&_Contract.CallOpts)
}

// GetTallyValue is a free data retrieval call binding the contract method 0xf0934fe9.
//
// Solidity: function GetTallyValue() view returns((uint256,uint256)[])
func (_Contract *ContractCaller) GetTallyValue(opts *bind.CallOpts) ([]VerificationG1Point, error) {
	var out []interface{}
	err := _Contract.contract.Call(opts, &out, "GetTallyValue")

	if err != nil {
		return *new([]VerificationG1Point), err
	}

	out0 := *abi.ConvertType(out[0], new([]VerificationG1Point)).(*[]VerificationG1Point)

	return out0, err

}

// GetTallyValue is a free data retrieval call binding the contract method 0xf0934fe9.
//
// Solidity: function GetTallyValue() view returns((uint256,uint256)[])
func (_Contract *ContractSession) GetTallyValue() ([]VerificationG1Point, error) {
	return _Contract.Contract.GetTallyValue(&_Contract.CallOpts)
}

// GetTallyValue is a free data retrieval call binding the contract method 0xf0934fe9.
//
// Solidity: function GetTallyValue() view returns((uint256,uint256)[])
func (_Contract *ContractCallerSession) GetTallyValue() ([]VerificationG1Point, error) {
	return _Contract.Contract.GetTallyValue(&_Contract.CallOpts)
}

// GetZKRPResult is a free data retrieval call binding the contract method 0x75457b93.
//
// Solidity: function GetZKRPResult() view returns(bool[])
func (_Contract *ContractCaller) GetZKRPResult(opts *bind.CallOpts) ([]bool, error) {
	var out []interface{}
	err := _Contract.contract.Call(opts, &out, "GetZKRPResult")

	if err != nil {
		return *new([]bool), err
	}

	out0 := *abi.ConvertType(out[0], new([]bool)).(*[]bool)

	return out0, err

}

// GetZKRPResult is a free data retrieval call binding the contract method 0x75457b93.
//
// Solidity: function GetZKRPResult() view returns(bool[])
func (_Contract *ContractSession) GetZKRPResult() ([]bool, error) {
	return _Contract.Contract.GetZKRPResult(&_Contract.CallOpts)
}

// GetZKRPResult is a free data retrieval call binding the contract method 0x75457b93.
//
// Solidity: function GetZKRPResult() view returns(bool[])
func (_Contract *ContractCallerSession) GetZKRPResult() ([]bool, error) {
	return _Contract.Contract.GetZKRPResult(&_Contract.CallOpts)
}

// DLEQVerify is a paid mutator transaction binding the contract method 0x9e134b59.
//
// Solidity: function DLEQVerify((uint256,uint256) g, (uint256,uint256) y1, (uint256,uint256) a1, (uint256,uint256) h, (uint256,uint256) y2, (uint256,uint256) a2, uint256 c, uint256 z) payable returns(bool)
func (_Contract *ContractTransactor) DLEQVerify(opts *bind.TransactOpts, g VerificationG1Point, y1 VerificationG1Point, a1 VerificationG1Point, h VerificationG1Point, y2 VerificationG1Point, a2 VerificationG1Point, c *big.Int, z *big.Int) (*types.Transaction, error) {
	return _Contract.contract.Transact(opts, "DLEQVerify", g, y1, a1, h, y2, a2, c, z)
}

// DLEQVerify is a paid mutator transaction binding the contract method 0x9e134b59.
//
// Solidity: function DLEQVerify((uint256,uint256) g, (uint256,uint256) y1, (uint256,uint256) a1, (uint256,uint256) h, (uint256,uint256) y2, (uint256,uint256) a2, uint256 c, uint256 z) payable returns(bool)
func (_Contract *ContractSession) DLEQVerify(g VerificationG1Point, y1 VerificationG1Point, a1 VerificationG1Point, h VerificationG1Point, y2 VerificationG1Point, a2 VerificationG1Point, c *big.Int, z *big.Int) (*types.Transaction, error) {
	return _Contract.Contract.DLEQVerify(&_Contract.TransactOpts, g, y1, a1, h, y2, a2, c, z)
}

// DLEQVerify is a paid mutator transaction binding the contract method 0x9e134b59.
//
// Solidity: function DLEQVerify((uint256,uint256) g, (uint256,uint256) y1, (uint256,uint256) a1, (uint256,uint256) h, (uint256,uint256) y2, (uint256,uint256) a2, uint256 c, uint256 z) payable returns(bool)
func (_Contract *ContractTransactorSession) DLEQVerify(g VerificationG1Point, y1 VerificationG1Point, a1 VerificationG1Point, h VerificationG1Point, y2 VerificationG1Point, a2 VerificationG1Point, c *big.Int, z *big.Int) (*types.Transaction, error) {
	return _Contract.Contract.DLEQVerify(&_Contract.TransactOpts, g, y1, a1, h, y2, a2, c, z)
}

// DVerify is a paid mutator transaction binding the contract method 0x7e3d3acb.
//
// Solidity: function DVerify((uint256,uint256) g, (uint256,uint256)[] y1, (uint256,uint256)[] a1, (uint256,uint256)[] h, (uint256,uint256)[] y2, (uint256,uint256)[] a2, uint256[] c, uint256[] z) returns(bool)
func (_Contract *ContractTransactor) DVerify(opts *bind.TransactOpts, g VerificationG1Point, y1 []VerificationG1Point, a1 []VerificationG1Point, h []VerificationG1Point, y2 []VerificationG1Point, a2 []VerificationG1Point, c []*big.Int, z []*big.Int) (*types.Transaction, error) {
	return _Contract.contract.Transact(opts, "DVerify", g, y1, a1, h, y2, a2, c, z)
}

// DVerify is a paid mutator transaction binding the contract method 0x7e3d3acb.
//
// Solidity: function DVerify((uint256,uint256) g, (uint256,uint256)[] y1, (uint256,uint256)[] a1, (uint256,uint256)[] h, (uint256,uint256)[] y2, (uint256,uint256)[] a2, uint256[] c, uint256[] z) returns(bool)
func (_Contract *ContractSession) DVerify(g VerificationG1Point, y1 []VerificationG1Point, a1 []VerificationG1Point, h []VerificationG1Point, y2 []VerificationG1Point, a2 []VerificationG1Point, c []*big.Int, z []*big.Int) (*types.Transaction, error) {
	return _Contract.Contract.DVerify(&_Contract.TransactOpts, g, y1, a1, h, y2, a2, c, z)
}

// DVerify is a paid mutator transaction binding the contract method 0x7e3d3acb.
//
// Solidity: function DVerify((uint256,uint256) g, (uint256,uint256)[] y1, (uint256,uint256)[] a1, (uint256,uint256)[] h, (uint256,uint256)[] y2, (uint256,uint256)[] a2, uint256[] c, uint256[] z) returns(bool)
func (_Contract *ContractTransactorSession) DVerify(g VerificationG1Point, y1 []VerificationG1Point, a1 []VerificationG1Point, h []VerificationG1Point, y2 []VerificationG1Point, a2 []VerificationG1Point, c []*big.Int, z []*big.Int) (*types.Transaction, error) {
	return _Contract.Contract.DVerify(&_Contract.TransactOpts, g, y1, a1, h, y2, a2, c, z)
}

// Interpolation is a paid mutator transaction binding the contract method 0x446bf71c.
//
// Solidity: function Interpolation(uint256 d, (uint256,uint256)[] v, uint256[] indices, uint256 threshold) returns((uint256,uint256))
func (_Contract *ContractTransactor) Interpolation(opts *bind.TransactOpts, d *big.Int, v []VerificationG1Point, indices []*big.Int, threshold *big.Int) (*types.Transaction, error) {
	return _Contract.contract.Transact(opts, "Interpolation", d, v, indices, threshold)
}

// Interpolation is a paid mutator transaction binding the contract method 0x446bf71c.
//
// Solidity: function Interpolation(uint256 d, (uint256,uint256)[] v, uint256[] indices, uint256 threshold) returns((uint256,uint256))
func (_Contract *ContractSession) Interpolation(d *big.Int, v []VerificationG1Point, indices []*big.Int, threshold *big.Int) (*types.Transaction, error) {
	return _Contract.Contract.Interpolation(&_Contract.TransactOpts, d, v, indices, threshold)
}

// Interpolation is a paid mutator transaction binding the contract method 0x446bf71c.
//
// Solidity: function Interpolation(uint256 d, (uint256,uint256)[] v, uint256[] indices, uint256 threshold) returns((uint256,uint256))
func (_Contract *ContractTransactorSession) Interpolation(d *big.Int, v []VerificationG1Point, indices []*big.Int, threshold *big.Int) (*types.Transaction, error) {
	return _Contract.Contract.Interpolation(&_Contract.TransactOpts, d, v, indices, threshold)
}

// PVerify is a paid mutator transaction binding the contract method 0x4c4791ec.
//
// Solidity: function PVerify(uint256[] indexTallier, (uint256,uint256)[] DecShare, (uint256,uint256)[] a1, (uint256,uint256)[] a2, uint256[] challenge, uint256[] z, uint256 threshold) returns()
func (_Contract *ContractTransactor) PVerify(opts *bind.TransactOpts, indexTallier []*big.Int, DecShare []VerificationG1Point, a1 []VerificationG1Point, a2 []VerificationG1Point, challenge []*big.Int, z []*big.Int, threshold *big.Int) (*types.Transaction, error) {
	return _Contract.contract.Transact(opts, "PVerify", indexTallier, DecShare, a1, a2, challenge, z, threshold)
}

// PVerify is a paid mutator transaction binding the contract method 0x4c4791ec.
//
// Solidity: function PVerify(uint256[] indexTallier, (uint256,uint256)[] DecShare, (uint256,uint256)[] a1, (uint256,uint256)[] a2, uint256[] challenge, uint256[] z, uint256 threshold) returns()
func (_Contract *ContractSession) PVerify(indexTallier []*big.Int, DecShare []VerificationG1Point, a1 []VerificationG1Point, a2 []VerificationG1Point, challenge []*big.Int, z []*big.Int, threshold *big.Int) (*types.Transaction, error) {
	return _Contract.Contract.PVerify(&_Contract.TransactOpts, indexTallier, DecShare, a1, a2, challenge, z, threshold)
}

// PVerify is a paid mutator transaction binding the contract method 0x4c4791ec.
//
// Solidity: function PVerify(uint256[] indexTallier, (uint256,uint256)[] DecShare, (uint256,uint256)[] a1, (uint256,uint256)[] a2, uint256[] challenge, uint256[] z, uint256 threshold) returns()
func (_Contract *ContractTransactorSession) PVerify(indexTallier []*big.Int, DecShare []VerificationG1Point, a1 []VerificationG1Point, a2 []VerificationG1Point, challenge []*big.Int, z []*big.Int, threshold *big.Int) (*types.Transaction, error) {
	return _Contract.Contract.PVerify(&_Contract.TransactOpts, indexTallier, DecShare, a1, a2, challenge, z, threshold)
}

// PVerifyTally is a paid mutator transaction binding the contract method 0xe17afe57.
//
// Solidity: function PVerifyTally(uint256[] indexTallier, (uint256,uint256)[] DecShare, (uint256,uint256)[] a1, (uint256,uint256)[] a2, uint256[] challenge, uint256[] z, uint256 threshold, uint256 numCandidates) returns()
func (_Contract *ContractTransactor) PVerifyTally(opts *bind.TransactOpts, indexTallier []*big.Int, DecShare []VerificationG1Point, a1 []VerificationG1Point, a2 []VerificationG1Point, challenge []*big.Int, z []*big.Int, threshold *big.Int, numCandidates *big.Int) (*types.Transaction, error) {
	return _Contract.contract.Transact(opts, "PVerifyTally", indexTallier, DecShare, a1, a2, challenge, z, threshold, numCandidates)
}

// PVerifyTally is a paid mutator transaction binding the contract method 0xe17afe57.
//
// Solidity: function PVerifyTally(uint256[] indexTallier, (uint256,uint256)[] DecShare, (uint256,uint256)[] a1, (uint256,uint256)[] a2, uint256[] challenge, uint256[] z, uint256 threshold, uint256 numCandidates) returns()
func (_Contract *ContractSession) PVerifyTally(indexTallier []*big.Int, DecShare []VerificationG1Point, a1 []VerificationG1Point, a2 []VerificationG1Point, challenge []*big.Int, z []*big.Int, threshold *big.Int, numCandidates *big.Int) (*types.Transaction, error) {
	return _Contract.Contract.PVerifyTally(&_Contract.TransactOpts, indexTallier, DecShare, a1, a2, challenge, z, threshold, numCandidates)
}

// PVerifyTally is a paid mutator transaction binding the contract method 0xe17afe57.
//
// Solidity: function PVerifyTally(uint256[] indexTallier, (uint256,uint256)[] DecShare, (uint256,uint256)[] a1, (uint256,uint256)[] a2, uint256[] challenge, uint256[] z, uint256 threshold, uint256 numCandidates) returns()
func (_Contract *ContractTransactorSession) PVerifyTally(indexTallier []*big.Int, DecShare []VerificationG1Point, a1 []VerificationG1Point, a2 []VerificationG1Point, challenge []*big.Int, z []*big.Int, threshold *big.Int, numCandidates *big.Int) (*types.Transaction, error) {
	return _Contract.Contract.PVerifyTally(&_Contract.TransactOpts, indexTallier, DecShare, a1, a2, challenge, z, threshold, numCandidates)
}

// RScodeVerify is a paid mutator transaction binding the contract method 0xbeaed4a0.
//
// Solidity: function RScodeVerify((uint256,uint256)[] v) returns(bool)
func (_Contract *ContractTransactor) RScodeVerify(opts *bind.TransactOpts, v []VerificationG1Point) (*types.Transaction, error) {
	return _Contract.contract.Transact(opts, "RScodeVerify", v)
}

// RScodeVerify is a paid mutator transaction binding the contract method 0xbeaed4a0.
//
// Solidity: function RScodeVerify((uint256,uint256)[] v) returns(bool)
func (_Contract *ContractSession) RScodeVerify(v []VerificationG1Point) (*types.Transaction, error) {
	return _Contract.Contract.RScodeVerify(&_Contract.TransactOpts, v)
}

// RScodeVerify is a paid mutator transaction binding the contract method 0xbeaed4a0.
//
// Solidity: function RScodeVerify((uint256,uint256)[] v) returns(bool)
func (_Contract *ContractTransactorSession) RScodeVerify(v []VerificationG1Point) (*types.Transaction, error) {
	return _Contract.Contract.RScodeVerify(&_Contract.TransactOpts, v)
}

// TestAggregateCSet is a paid mutator transaction binding the contract method 0x5e08ddbf.
//
// Solidity: function TestAggregateCSet((uint256,uint256)[] c) returns()
func (_Contract *ContractTransactor) TestAggregateCSet(opts *bind.TransactOpts, c []VerificationG1Point) (*types.Transaction, error) {
	return _Contract.contract.Transact(opts, "TestAggregateCSet", c)
}

// TestAggregateCSet is a paid mutator transaction binding the contract method 0x5e08ddbf.
//
// Solidity: function TestAggregateCSet((uint256,uint256)[] c) returns()
func (_Contract *ContractSession) TestAggregateCSet(c []VerificationG1Point) (*types.Transaction, error) {
	return _Contract.Contract.TestAggregateCSet(&_Contract.TransactOpts, c)
}

// TestAggregateCSet is a paid mutator transaction binding the contract method 0x5e08ddbf.
//
// Solidity: function TestAggregateCSet((uint256,uint256)[] c) returns()
func (_Contract *ContractTransactorSession) TestAggregateCSet(c []VerificationG1Point) (*types.Transaction, error) {
	return _Contract.Contract.TestAggregateCSet(&_Contract.TransactOpts, c)
}

// UploadBallotCipher is a paid mutator transaction binding the contract method 0x65d28554.
//
// Solidity: function UploadBallotCipher((uint256,uint256)[] Ej, (uint256,uint256)[] Fj1, (uint256,uint256)[] Fj2, (uint256,uint256)[] _Uj, (uint256,uint256)[] _Cj, uint256[] c, uint256[] z1, uint256[] z2, uint256[] z3, (uint256,uint256)[] Uj, (uint256,uint256)[] v, uint256 threshold) returns()
func (_Contract *ContractTransactor) UploadBallotCipher(opts *bind.TransactOpts, Ej []VerificationG1Point, Fj1 []VerificationG1Point, Fj2 []VerificationG1Point, _Uj []VerificationG1Point, _Cj []VerificationG1Point, c []*big.Int, z1 []*big.Int, z2 []*big.Int, z3 []*big.Int, Uj []VerificationG1Point, v []VerificationG1Point, threshold *big.Int) (*types.Transaction, error) {
	return _Contract.contract.Transact(opts, "UploadBallotCipher", Ej, Fj1, Fj2, _Uj, _Cj, c, z1, z2, z3, Uj, v, threshold)
}

// UploadBallotCipher is a paid mutator transaction binding the contract method 0x65d28554.
//
// Solidity: function UploadBallotCipher((uint256,uint256)[] Ej, (uint256,uint256)[] Fj1, (uint256,uint256)[] Fj2, (uint256,uint256)[] _Uj, (uint256,uint256)[] _Cj, uint256[] c, uint256[] z1, uint256[] z2, uint256[] z3, (uint256,uint256)[] Uj, (uint256,uint256)[] v, uint256 threshold) returns()
func (_Contract *ContractSession) UploadBallotCipher(Ej []VerificationG1Point, Fj1 []VerificationG1Point, Fj2 []VerificationG1Point, _Uj []VerificationG1Point, _Cj []VerificationG1Point, c []*big.Int, z1 []*big.Int, z2 []*big.Int, z3 []*big.Int, Uj []VerificationG1Point, v []VerificationG1Point, threshold *big.Int) (*types.Transaction, error) {
	return _Contract.Contract.UploadBallotCipher(&_Contract.TransactOpts, Ej, Fj1, Fj2, _Uj, _Cj, c, z1, z2, z3, Uj, v, threshold)
}

// UploadBallotCipher is a paid mutator transaction binding the contract method 0x65d28554.
//
// Solidity: function UploadBallotCipher((uint256,uint256)[] Ej, (uint256,uint256)[] Fj1, (uint256,uint256)[] Fj2, (uint256,uint256)[] _Uj, (uint256,uint256)[] _Cj, uint256[] c, uint256[] z1, uint256[] z2, uint256[] z3, (uint256,uint256)[] Uj, (uint256,uint256)[] v, uint256 threshold) returns()
func (_Contract *ContractTransactorSession) UploadBallotCipher(Ej []VerificationG1Point, Fj1 []VerificationG1Point, Fj2 []VerificationG1Point, _Uj []VerificationG1Point, _Cj []VerificationG1Point, c []*big.Int, z1 []*big.Int, z2 []*big.Int, z3 []*big.Int, Uj []VerificationG1Point, v []VerificationG1Point, threshold *big.Int) (*types.Transaction, error) {
	return _Contract.Contract.UploadBallotCipher(&_Contract.TransactOpts, Ej, Fj1, Fj2, _Uj, _Cj, c, z1, z2, z3, Uj, v, threshold)
}

// UploadBallotCipher1 is a paid mutator transaction binding the contract method 0x43ac6ba4.
//
// Solidity: function UploadBallotCipher1((uint256,uint256)[] Ej, (uint256,uint256)[] Fj1, (uint256,uint256)[] Fj2, (uint256,uint256)[] _Uj, (uint256,uint256)[] _Cj, uint256[] c, uint256[] z1, uint256[] z2, uint256[] z3, (uint256,uint256)[] Uj, (uint256,uint256)[] v, uint256 threshold) returns()
func (_Contract *ContractTransactor) UploadBallotCipher1(opts *bind.TransactOpts, Ej []VerificationG1Point, Fj1 []VerificationG1Point, Fj2 []VerificationG1Point, _Uj []VerificationG1Point, _Cj []VerificationG1Point, c []*big.Int, z1 []*big.Int, z2 []*big.Int, z3 []*big.Int, Uj []VerificationG1Point, v []VerificationG1Point, threshold *big.Int) (*types.Transaction, error) {
	return _Contract.contract.Transact(opts, "UploadBallotCipher1", Ej, Fj1, Fj2, _Uj, _Cj, c, z1, z2, z3, Uj, v, threshold)
}

// UploadBallotCipher1 is a paid mutator transaction binding the contract method 0x43ac6ba4.
//
// Solidity: function UploadBallotCipher1((uint256,uint256)[] Ej, (uint256,uint256)[] Fj1, (uint256,uint256)[] Fj2, (uint256,uint256)[] _Uj, (uint256,uint256)[] _Cj, uint256[] c, uint256[] z1, uint256[] z2, uint256[] z3, (uint256,uint256)[] Uj, (uint256,uint256)[] v, uint256 threshold) returns()
func (_Contract *ContractSession) UploadBallotCipher1(Ej []VerificationG1Point, Fj1 []VerificationG1Point, Fj2 []VerificationG1Point, _Uj []VerificationG1Point, _Cj []VerificationG1Point, c []*big.Int, z1 []*big.Int, z2 []*big.Int, z3 []*big.Int, Uj []VerificationG1Point, v []VerificationG1Point, threshold *big.Int) (*types.Transaction, error) {
	return _Contract.Contract.UploadBallotCipher1(&_Contract.TransactOpts, Ej, Fj1, Fj2, _Uj, _Cj, c, z1, z2, z3, Uj, v, threshold)
}

// UploadBallotCipher1 is a paid mutator transaction binding the contract method 0x43ac6ba4.
//
// Solidity: function UploadBallotCipher1((uint256,uint256)[] Ej, (uint256,uint256)[] Fj1, (uint256,uint256)[] Fj2, (uint256,uint256)[] _Uj, (uint256,uint256)[] _Cj, uint256[] c, uint256[] z1, uint256[] z2, uint256[] z3, (uint256,uint256)[] Uj, (uint256,uint256)[] v, uint256 threshold) returns()
func (_Contract *ContractTransactorSession) UploadBallotCipher1(Ej []VerificationG1Point, Fj1 []VerificationG1Point, Fj2 []VerificationG1Point, _Uj []VerificationG1Point, _Cj []VerificationG1Point, c []*big.Int, z1 []*big.Int, z2 []*big.Int, z3 []*big.Int, Uj []VerificationG1Point, v []VerificationG1Point, threshold *big.Int) (*types.Transaction, error) {
	return _Contract.Contract.UploadBallotCipher1(&_Contract.TransactOpts, Ej, Fj1, Fj2, _Uj, _Cj, c, z1, z2, z3, Uj, v, threshold)
}

// UploadPVSSShares is a paid mutator transaction binding the contract method 0xe283f130.
//
// Solidity: function UploadPVSSShares((uint256,uint256)[] v, (uint256,uint256)[] c, (uint256,uint256)[] a1, (uint256,uint256)[] a2, uint256[] challenge, uint256[] z) returns()
func (_Contract *ContractTransactor) UploadPVSSShares(opts *bind.TransactOpts, v []VerificationG1Point, c []VerificationG1Point, a1 []VerificationG1Point, a2 []VerificationG1Point, challenge []*big.Int, z []*big.Int) (*types.Transaction, error) {
	return _Contract.contract.Transact(opts, "UploadPVSSShares", v, c, a1, a2, challenge, z)
}

// UploadPVSSShares is a paid mutator transaction binding the contract method 0xe283f130.
//
// Solidity: function UploadPVSSShares((uint256,uint256)[] v, (uint256,uint256)[] c, (uint256,uint256)[] a1, (uint256,uint256)[] a2, uint256[] challenge, uint256[] z) returns()
func (_Contract *ContractSession) UploadPVSSShares(v []VerificationG1Point, c []VerificationG1Point, a1 []VerificationG1Point, a2 []VerificationG1Point, challenge []*big.Int, z []*big.Int) (*types.Transaction, error) {
	return _Contract.Contract.UploadPVSSShares(&_Contract.TransactOpts, v, c, a1, a2, challenge, z)
}

// UploadPVSSShares is a paid mutator transaction binding the contract method 0xe283f130.
//
// Solidity: function UploadPVSSShares((uint256,uint256)[] v, (uint256,uint256)[] c, (uint256,uint256)[] a1, (uint256,uint256)[] a2, uint256[] challenge, uint256[] z) returns()
func (_Contract *ContractTransactorSession) UploadPVSSShares(v []VerificationG1Point, c []VerificationG1Point, a1 []VerificationG1Point, a2 []VerificationG1Point, challenge []*big.Int, z []*big.Int) (*types.Transaction, error) {
	return _Contract.Contract.UploadPVSSShares(&_Contract.TransactOpts, v, c, a1, a2, challenge, z)
}

// UploadPVSSShares1 is a paid mutator transaction binding the contract method 0xea315a94.
//
// Solidity: function UploadPVSSShares1((uint256,uint256)[] v, (uint256,uint256)[] c, (uint256,uint256)[] a1, (uint256,uint256)[] a2, uint256[] challenge, uint256[] z) returns()
func (_Contract *ContractTransactor) UploadPVSSShares1(opts *bind.TransactOpts, v []VerificationG1Point, c []VerificationG1Point, a1 []VerificationG1Point, a2 []VerificationG1Point, challenge []*big.Int, z []*big.Int) (*types.Transaction, error) {
	return _Contract.contract.Transact(opts, "UploadPVSSShares1", v, c, a1, a2, challenge, z)
}

// UploadPVSSShares1 is a paid mutator transaction binding the contract method 0xea315a94.
//
// Solidity: function UploadPVSSShares1((uint256,uint256)[] v, (uint256,uint256)[] c, (uint256,uint256)[] a1, (uint256,uint256)[] a2, uint256[] challenge, uint256[] z) returns()
func (_Contract *ContractSession) UploadPVSSShares1(v []VerificationG1Point, c []VerificationG1Point, a1 []VerificationG1Point, a2 []VerificationG1Point, challenge []*big.Int, z []*big.Int) (*types.Transaction, error) {
	return _Contract.Contract.UploadPVSSShares1(&_Contract.TransactOpts, v, c, a1, a2, challenge, z)
}

// UploadPVSSShares1 is a paid mutator transaction binding the contract method 0xea315a94.
//
// Solidity: function UploadPVSSShares1((uint256,uint256)[] v, (uint256,uint256)[] c, (uint256,uint256)[] a1, (uint256,uint256)[] a2, uint256[] challenge, uint256[] z) returns()
func (_Contract *ContractTransactorSession) UploadPVSSShares1(v []VerificationG1Point, c []VerificationG1Point, a1 []VerificationG1Point, a2 []VerificationG1Point, challenge []*big.Int, z []*big.Int) (*types.Transaction, error) {
	return _Contract.Contract.UploadPVSSShares1(&_Contract.TransactOpts, v, c, a1, a2, challenge, z)
}

// UploadParameters is a paid mutator transaction binding the contract method 0x4dea579f.
//
// Solidity: function UploadParameters((uint256,uint256) g0, (uint256,uint256) h0, (uint256[2],uint256[2]) g1, (uint256[2],uint256[2]) pkI, (uint256,uint256)[] sigmak, uint256 a, uint256 b, uint256 numCandidates, uint256 numTalliers) returns()
func (_Contract *ContractTransactor) UploadParameters(opts *bind.TransactOpts, g0 VerificationG1Point, h0 VerificationG1Point, g1 VerificationG2Point, pkI VerificationG2Point, sigmak []VerificationG1Point, a *big.Int, b *big.Int, numCandidates *big.Int, numTalliers *big.Int) (*types.Transaction, error) {
	return _Contract.contract.Transact(opts, "UploadParameters", g0, h0, g1, pkI, sigmak, a, b, numCandidates, numTalliers)
}

// UploadParameters is a paid mutator transaction binding the contract method 0x4dea579f.
//
// Solidity: function UploadParameters((uint256,uint256) g0, (uint256,uint256) h0, (uint256[2],uint256[2]) g1, (uint256[2],uint256[2]) pkI, (uint256,uint256)[] sigmak, uint256 a, uint256 b, uint256 numCandidates, uint256 numTalliers) returns()
func (_Contract *ContractSession) UploadParameters(g0 VerificationG1Point, h0 VerificationG1Point, g1 VerificationG2Point, pkI VerificationG2Point, sigmak []VerificationG1Point, a *big.Int, b *big.Int, numCandidates *big.Int, numTalliers *big.Int) (*types.Transaction, error) {
	return _Contract.Contract.UploadParameters(&_Contract.TransactOpts, g0, h0, g1, pkI, sigmak, a, b, numCandidates, numTalliers)
}

// UploadParameters is a paid mutator transaction binding the contract method 0x4dea579f.
//
// Solidity: function UploadParameters((uint256,uint256) g0, (uint256,uint256) h0, (uint256[2],uint256[2]) g1, (uint256[2],uint256[2]) pkI, (uint256,uint256)[] sigmak, uint256 a, uint256 b, uint256 numCandidates, uint256 numTalliers) returns()
func (_Contract *ContractTransactorSession) UploadParameters(g0 VerificationG1Point, h0 VerificationG1Point, g1 VerificationG2Point, pkI VerificationG2Point, sigmak []VerificationG1Point, a *big.Int, b *big.Int, numCandidates *big.Int, numTalliers *big.Int) (*types.Transaction, error) {
	return _Contract.Contract.UploadParameters(&_Contract.TransactOpts, g0, h0, g1, pkI, sigmak, a, b, numCandidates, numTalliers)
}

// UploadPublicKey is a paid mutator transaction binding the contract method 0x4ea9637d.
//
// Solidity: function UploadPublicKey((uint256,uint256)[] pks) returns()
func (_Contract *ContractTransactor) UploadPublicKey(opts *bind.TransactOpts, pks []VerificationG1Point) (*types.Transaction, error) {
	return _Contract.contract.Transact(opts, "UploadPublicKey", pks)
}

// UploadPublicKey is a paid mutator transaction binding the contract method 0x4ea9637d.
//
// Solidity: function UploadPublicKey((uint256,uint256)[] pks) returns()
func (_Contract *ContractSession) UploadPublicKey(pks []VerificationG1Point) (*types.Transaction, error) {
	return _Contract.Contract.UploadPublicKey(&_Contract.TransactOpts, pks)
}

// UploadPublicKey is a paid mutator transaction binding the contract method 0x4ea9637d.
//
// Solidity: function UploadPublicKey((uint256,uint256)[] pks) returns()
func (_Contract *ContractTransactorSession) UploadPublicKey(pks []VerificationG1Point) (*types.Transaction, error) {
	return _Contract.Contract.UploadPublicKey(&_Contract.TransactOpts, pks)
}

// ZKRPVerify is a paid mutator transaction binding the contract method 0xb562c4ab.
//
// Solidity: function ZKRPVerify((uint256,uint256)[] Ej, (uint256,uint256)[] Fj1, (uint256,uint256)[] Fj2, (uint256,uint256)[] _Uj, (uint256,uint256)[] _Cj, uint256[] c, uint256[] z1, uint256[] z2, uint256[] z3, (uint256,uint256)[] Uj, uint256[] d, (uint256,uint256)[] v, uint256[] indices, uint256 threshold) returns(bool)
func (_Contract *ContractTransactor) ZKRPVerify(opts *bind.TransactOpts, Ej []VerificationG1Point, Fj1 []VerificationG1Point, Fj2 []VerificationG1Point, _Uj []VerificationG1Point, _Cj []VerificationG1Point, c []*big.Int, z1 []*big.Int, z2 []*big.Int, z3 []*big.Int, Uj []VerificationG1Point, d []*big.Int, v []VerificationG1Point, indices []*big.Int, threshold *big.Int) (*types.Transaction, error) {
	return _Contract.contract.Transact(opts, "ZKRPVerify", Ej, Fj1, Fj2, _Uj, _Cj, c, z1, z2, z3, Uj, d, v, indices, threshold)
}

// ZKRPVerify is a paid mutator transaction binding the contract method 0xb562c4ab.
//
// Solidity: function ZKRPVerify((uint256,uint256)[] Ej, (uint256,uint256)[] Fj1, (uint256,uint256)[] Fj2, (uint256,uint256)[] _Uj, (uint256,uint256)[] _Cj, uint256[] c, uint256[] z1, uint256[] z2, uint256[] z3, (uint256,uint256)[] Uj, uint256[] d, (uint256,uint256)[] v, uint256[] indices, uint256 threshold) returns(bool)
func (_Contract *ContractSession) ZKRPVerify(Ej []VerificationG1Point, Fj1 []VerificationG1Point, Fj2 []VerificationG1Point, _Uj []VerificationG1Point, _Cj []VerificationG1Point, c []*big.Int, z1 []*big.Int, z2 []*big.Int, z3 []*big.Int, Uj []VerificationG1Point, d []*big.Int, v []VerificationG1Point, indices []*big.Int, threshold *big.Int) (*types.Transaction, error) {
	return _Contract.Contract.ZKRPVerify(&_Contract.TransactOpts, Ej, Fj1, Fj2, _Uj, _Cj, c, z1, z2, z3, Uj, d, v, indices, threshold)
}

// ZKRPVerify is a paid mutator transaction binding the contract method 0xb562c4ab.
//
// Solidity: function ZKRPVerify((uint256,uint256)[] Ej, (uint256,uint256)[] Fj1, (uint256,uint256)[] Fj2, (uint256,uint256)[] _Uj, (uint256,uint256)[] _Cj, uint256[] c, uint256[] z1, uint256[] z2, uint256[] z3, (uint256,uint256)[] Uj, uint256[] d, (uint256,uint256)[] v, uint256[] indices, uint256 threshold) returns(bool)
func (_Contract *ContractTransactorSession) ZKRPVerify(Ej []VerificationG1Point, Fj1 []VerificationG1Point, Fj2 []VerificationG1Point, _Uj []VerificationG1Point, _Cj []VerificationG1Point, c []*big.Int, z1 []*big.Int, z2 []*big.Int, z3 []*big.Int, Uj []VerificationG1Point, d []*big.Int, v []VerificationG1Point, indices []*big.Int, threshold *big.Int) (*types.Transaction, error) {
	return _Contract.Contract.ZKRPVerify(&_Contract.TransactOpts, Ej, Fj1, Fj2, _Uj, _Cj, c, z1, z2, z3, Uj, d, v, indices, threshold)
}

// Coefficient is a paid mutator transaction binding the contract method 0x5b4af2b8.
//
// Solidity: function coefficient(int256 i, uint256 n, uint256 l) returns(uint256)
func (_Contract *ContractTransactor) Coefficient(opts *bind.TransactOpts, i *big.Int, n *big.Int, l *big.Int) (*types.Transaction, error) {
	return _Contract.contract.Transact(opts, "coefficient", i, n, l)
}

// Coefficient is a paid mutator transaction binding the contract method 0x5b4af2b8.
//
// Solidity: function coefficient(int256 i, uint256 n, uint256 l) returns(uint256)
func (_Contract *ContractSession) Coefficient(i *big.Int, n *big.Int, l *big.Int) (*types.Transaction, error) {
	return _Contract.Contract.Coefficient(&_Contract.TransactOpts, i, n, l)
}

// Coefficient is a paid mutator transaction binding the contract method 0x5b4af2b8.
//
// Solidity: function coefficient(int256 i, uint256 n, uint256 l) returns(uint256)
func (_Contract *ContractTransactorSession) Coefficient(i *big.Int, n *big.Int, l *big.Int) (*types.Transaction, error) {
	return _Contract.Contract.Coefficient(&_Contract.TransactOpts, i, n, l)
}

// Inv is a paid mutator transaction binding the contract method 0x338255f3.
//
// Solidity: function inv(uint256 a, uint256 prime) returns(uint256)
func (_Contract *ContractTransactor) Inv(opts *bind.TransactOpts, a *big.Int, prime *big.Int) (*types.Transaction, error) {
	return _Contract.contract.Transact(opts, "inv", a, prime)
}

// Inv is a paid mutator transaction binding the contract method 0x338255f3.
//
// Solidity: function inv(uint256 a, uint256 prime) returns(uint256)
func (_Contract *ContractSession) Inv(a *big.Int, prime *big.Int) (*types.Transaction, error) {
	return _Contract.Contract.Inv(&_Contract.TransactOpts, a, prime)
}

// Inv is a paid mutator transaction binding the contract method 0x338255f3.
//
// Solidity: function inv(uint256 a, uint256 prime) returns(uint256)
func (_Contract *ContractTransactorSession) Inv(a *big.Int, prime *big.Int) (*types.Transaction, error) {
	return _Contract.Contract.Inv(&_Contract.TransactOpts, a, prime)
}
