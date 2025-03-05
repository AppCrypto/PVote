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
	ABI: "[{\"inputs\":[],\"name\":\"AggregateCiphertext\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"AggregateEncShares\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"X\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"Y\",\"type\":\"uint256\"}],\"internalType\":\"structVerification.G1Point\",\"name\":\"g\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"X\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"Y\",\"type\":\"uint256\"}],\"internalType\":\"structVerification.G1Point\",\"name\":\"y1\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"X\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"Y\",\"type\":\"uint256\"}],\"internalType\":\"structVerification.G1Point\",\"name\":\"a1\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"X\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"Y\",\"type\":\"uint256\"}],\"internalType\":\"structVerification.G1Point\",\"name\":\"h\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"X\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"Y\",\"type\":\"uint256\"}],\"internalType\":\"structVerification.G1Point\",\"name\":\"y2\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"X\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"Y\",\"type\":\"uint256\"}],\"internalType\":\"structVerification.G1Point\",\"name\":\"a2\",\"type\":\"tuple\"},{\"internalType\":\"uint256\",\"name\":\"c\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"z\",\"type\":\"uint256\"}],\"name\":\"DLEQVerify\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"X\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"Y\",\"type\":\"uint256\"}],\"internalType\":\"structVerification.G1Point\",\"name\":\"g\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"X\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"Y\",\"type\":\"uint256\"}],\"internalType\":\"structVerification.G1Point[]\",\"name\":\"y1\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"X\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"Y\",\"type\":\"uint256\"}],\"internalType\":\"structVerification.G1Point[]\",\"name\":\"a1\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"X\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"Y\",\"type\":\"uint256\"}],\"internalType\":\"structVerification.G1Point[]\",\"name\":\"h\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"X\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"Y\",\"type\":\"uint256\"}],\"internalType\":\"structVerification.G1Point[]\",\"name\":\"y2\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"X\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"Y\",\"type\":\"uint256\"}],\"internalType\":\"structVerification.G1Point[]\",\"name\":\"a2\",\"type\":\"tuple[]\"},{\"internalType\":\"uint256[]\",\"name\":\"c\",\"type\":\"uint256[]\"},{\"internalType\":\"uint256[]\",\"name\":\"z\",\"type\":\"uint256[]\"}],\"name\":\"DVerify\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"GetAggregateValue\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"X\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"Y\",\"type\":\"uint256\"}],\"internalType\":\"structVerification.G1Point[]\",\"name\":\"\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"GetDVerifyResult\",\"outputs\":[{\"internalType\":\"bool[]\",\"name\":\"\",\"type\":\"bool[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"GetPVerifyResult\",\"outputs\":[{\"internalType\":\"bool[]\",\"name\":\"\",\"type\":\"bool[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"GetTallyValue\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"X\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"Y\",\"type\":\"uint256\"}],\"internalType\":\"structVerification.G1Point[]\",\"name\":\"\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"GetZKRPResult\",\"outputs\":[{\"internalType\":\"bool[]\",\"name\":\"\",\"type\":\"bool[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"d\",\"type\":\"uint256\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"X\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"Y\",\"type\":\"uint256\"}],\"internalType\":\"structVerification.G1Point[]\",\"name\":\"v\",\"type\":\"tuple[]\"},{\"internalType\":\"uint256[]\",\"name\":\"indices\",\"type\":\"uint256[]\"},{\"internalType\":\"uint256\",\"name\":\"threshold\",\"type\":\"uint256\"}],\"name\":\"Interpolation\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"X\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"Y\",\"type\":\"uint256\"}],\"internalType\":\"structVerification.G1Point\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"indexTallier\",\"type\":\"uint256\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"X\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"Y\",\"type\":\"uint256\"}],\"internalType\":\"structVerification.G1Point\",\"name\":\"DecShare\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"X\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"Y\",\"type\":\"uint256\"}],\"internalType\":\"structVerification.G1Point\",\"name\":\"a1\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"X\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"Y\",\"type\":\"uint256\"}],\"internalType\":\"structVerification.G1Point\",\"name\":\"a2\",\"type\":\"tuple\"},{\"internalType\":\"uint256\",\"name\":\"challenge\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"z\",\"type\":\"uint256\"}],\"name\":\"PVerify\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"X\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"Y\",\"type\":\"uint256\"}],\"internalType\":\"structVerification.G1Point[]\",\"name\":\"v\",\"type\":\"tuple[]\"}],\"name\":\"RScodeVerify\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"threshold\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"numCandidates\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"a\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"b\",\"type\":\"uint256\"}],\"name\":\"Tally\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"X\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"Y\",\"type\":\"uint256\"}],\"internalType\":\"structVerification.G1Point\",\"name\":\"c\",\"type\":\"tuple\"}],\"name\":\"TestAggregateC\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"X\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"Y\",\"type\":\"uint256\"}],\"internalType\":\"structVerification.G1Point[]\",\"name\":\"c\",\"type\":\"tuple[]\"}],\"name\":\"TestAggregateCSet\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"X\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"Y\",\"type\":\"uint256\"}],\"internalType\":\"structVerification.G1Point[]\",\"name\":\"U\",\"type\":\"tuple[]\"}],\"name\":\"TestAggregateUSet\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"X\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"Y\",\"type\":\"uint256\"}],\"internalType\":\"structVerification.G1Point[]\",\"name\":\"Ej\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"X\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"Y\",\"type\":\"uint256\"}],\"internalType\":\"structVerification.G1Point[]\",\"name\":\"Fj1\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"X\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"Y\",\"type\":\"uint256\"}],\"internalType\":\"structVerification.G1Point[]\",\"name\":\"Fj2\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"X\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"Y\",\"type\":\"uint256\"}],\"internalType\":\"structVerification.G1Point[]\",\"name\":\"_Uj\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"X\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"Y\",\"type\":\"uint256\"}],\"internalType\":\"structVerification.G1Point[]\",\"name\":\"_Cj\",\"type\":\"tuple[]\"},{\"internalType\":\"uint256[]\",\"name\":\"c\",\"type\":\"uint256[]\"},{\"internalType\":\"uint256[]\",\"name\":\"z1\",\"type\":\"uint256[]\"},{\"internalType\":\"uint256[]\",\"name\":\"z2\",\"type\":\"uint256[]\"},{\"internalType\":\"uint256[]\",\"name\":\"z3\",\"type\":\"uint256[]\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"X\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"Y\",\"type\":\"uint256\"}],\"internalType\":\"structVerification.G1Point[]\",\"name\":\"Uj\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"X\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"Y\",\"type\":\"uint256\"}],\"internalType\":\"structVerification.G1Point[]\",\"name\":\"v\",\"type\":\"tuple[]\"},{\"internalType\":\"uint256\",\"name\":\"threshold\",\"type\":\"uint256\"}],\"name\":\"UploadBallotCipher\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"X\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"Y\",\"type\":\"uint256\"}],\"internalType\":\"structVerification.G1Point[]\",\"name\":\"Ej\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"X\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"Y\",\"type\":\"uint256\"}],\"internalType\":\"structVerification.G1Point[]\",\"name\":\"Fj1\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"X\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"Y\",\"type\":\"uint256\"}],\"internalType\":\"structVerification.G1Point[]\",\"name\":\"Fj2\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"X\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"Y\",\"type\":\"uint256\"}],\"internalType\":\"structVerification.G1Point[]\",\"name\":\"_Uj\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"X\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"Y\",\"type\":\"uint256\"}],\"internalType\":\"structVerification.G1Point[]\",\"name\":\"_Cj\",\"type\":\"tuple[]\"},{\"internalType\":\"uint256[]\",\"name\":\"c\",\"type\":\"uint256[]\"},{\"internalType\":\"uint256[]\",\"name\":\"z1\",\"type\":\"uint256[]\"},{\"internalType\":\"uint256[]\",\"name\":\"z2\",\"type\":\"uint256[]\"},{\"internalType\":\"uint256[]\",\"name\":\"z3\",\"type\":\"uint256[]\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"X\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"Y\",\"type\":\"uint256\"}],\"internalType\":\"structVerification.G1Point[]\",\"name\":\"Uj\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"X\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"Y\",\"type\":\"uint256\"}],\"internalType\":\"structVerification.G1Point[]\",\"name\":\"v\",\"type\":\"tuple[]\"},{\"internalType\":\"uint256\",\"name\":\"threshold\",\"type\":\"uint256\"}],\"name\":\"UploadBallotCipher1\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"X\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"Y\",\"type\":\"uint256\"}],\"internalType\":\"structVerification.G1Point[]\",\"name\":\"v\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"X\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"Y\",\"type\":\"uint256\"}],\"internalType\":\"structVerification.G1Point[]\",\"name\":\"c\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"X\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"Y\",\"type\":\"uint256\"}],\"internalType\":\"structVerification.G1Point[]\",\"name\":\"a1\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"X\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"Y\",\"type\":\"uint256\"}],\"internalType\":\"structVerification.G1Point[]\",\"name\":\"a2\",\"type\":\"tuple[]\"},{\"internalType\":\"uint256[]\",\"name\":\"challenge\",\"type\":\"uint256[]\"},{\"internalType\":\"uint256[]\",\"name\":\"z\",\"type\":\"uint256[]\"}],\"name\":\"UploadPVSSShares\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"X\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"Y\",\"type\":\"uint256\"}],\"internalType\":\"structVerification.G1Point[]\",\"name\":\"v\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"X\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"Y\",\"type\":\"uint256\"}],\"internalType\":\"structVerification.G1Point[]\",\"name\":\"c\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"X\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"Y\",\"type\":\"uint256\"}],\"internalType\":\"structVerification.G1Point[]\",\"name\":\"a1\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"X\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"Y\",\"type\":\"uint256\"}],\"internalType\":\"structVerification.G1Point[]\",\"name\":\"a2\",\"type\":\"tuple[]\"},{\"internalType\":\"uint256[]\",\"name\":\"challenge\",\"type\":\"uint256[]\"},{\"internalType\":\"uint256[]\",\"name\":\"z\",\"type\":\"uint256[]\"}],\"name\":\"UploadPVSSShares1\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"X\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"Y\",\"type\":\"uint256\"}],\"internalType\":\"structVerification.G1Point\",\"name\":\"g0\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"X\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"Y\",\"type\":\"uint256\"}],\"internalType\":\"structVerification.G1Point\",\"name\":\"h0\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"uint256[2]\",\"name\":\"X\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256[2]\",\"name\":\"Y\",\"type\":\"uint256[2]\"}],\"internalType\":\"structVerification.G2Point\",\"name\":\"g1\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"uint256[2]\",\"name\":\"X\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256[2]\",\"name\":\"Y\",\"type\":\"uint256[2]\"}],\"internalType\":\"structVerification.G2Point\",\"name\":\"pkI\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"X\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"Y\",\"type\":\"uint256\"}],\"internalType\":\"structVerification.G1Point[]\",\"name\":\"sigmak\",\"type\":\"tuple[]\"},{\"internalType\":\"uint256\",\"name\":\"a\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"b\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"numCandidates\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"numTalliers\",\"type\":\"uint256\"}],\"name\":\"UploadParameters\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"X\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"Y\",\"type\":\"uint256\"}],\"internalType\":\"structVerification.G1Point[]\",\"name\":\"pks\",\"type\":\"tuple[]\"}],\"name\":\"UploadPublicKey\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"X\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"Y\",\"type\":\"uint256\"}],\"internalType\":\"structVerification.G1Point[]\",\"name\":\"Ej\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"X\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"Y\",\"type\":\"uint256\"}],\"internalType\":\"structVerification.G1Point[]\",\"name\":\"Fj1\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"X\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"Y\",\"type\":\"uint256\"}],\"internalType\":\"structVerification.G1Point[]\",\"name\":\"Fj2\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"X\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"Y\",\"type\":\"uint256\"}],\"internalType\":\"structVerification.G1Point[]\",\"name\":\"_Uj\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"X\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"Y\",\"type\":\"uint256\"}],\"internalType\":\"structVerification.G1Point[]\",\"name\":\"_Cj\",\"type\":\"tuple[]\"},{\"internalType\":\"uint256[]\",\"name\":\"c\",\"type\":\"uint256[]\"},{\"internalType\":\"uint256[]\",\"name\":\"z1\",\"type\":\"uint256[]\"},{\"internalType\":\"uint256[]\",\"name\":\"z2\",\"type\":\"uint256[]\"},{\"internalType\":\"uint256[]\",\"name\":\"z3\",\"type\":\"uint256[]\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"X\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"Y\",\"type\":\"uint256\"}],\"internalType\":\"structVerification.G1Point[]\",\"name\":\"Uj\",\"type\":\"tuple[]\"},{\"internalType\":\"uint256[]\",\"name\":\"d\",\"type\":\"uint256[]\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"X\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"Y\",\"type\":\"uint256\"}],\"internalType\":\"structVerification.G1Point[]\",\"name\":\"v\",\"type\":\"tuple[]\"},{\"internalType\":\"uint256[]\",\"name\":\"indices\",\"type\":\"uint256[]\"},{\"internalType\":\"uint256\",\"name\":\"threshold\",\"type\":\"uint256\"}],\"name\":\"ZKRPVerify\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"int256\",\"name\":\"i\",\"type\":\"int256\"},{\"internalType\":\"uint256\",\"name\":\"n\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"l\",\"type\":\"uint256\"}],\"name\":\"coefficient\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"a\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"prime\",\"type\":\"uint256\"}],\"name\":\"inv\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x608080604052346015576134e6908161001b8239f35b600080fdfe6080604052600436101561001257600080fd5b60003560e01c80630c65d6e2146101a7578063338255f3146101a257806343ac6ba41461019d578063446bf71c146101985780634dea579f146101935780634ea9637d1461018e5780635b4af2b8146101895780635e08ddbf1461018457806365d285541461017f57806373c033941461017a57806375457b93146101755780637e3d3acb146101705780638823f8a51461016b5780639e134b5914610166578063a5ac309114610161578063b002450b1461015c578063b52c5ff214610157578063b562c4ab14610152578063beaed4a01461014d578063c3ceb16b14610148578063db301a7f14610143578063e283f1301461013e578063e7db3bb214610139578063ea315a94146101345763f0934fe91461012f57600080fd5b611d50565b611d33565b611c32565b611c15565b6119e1565b611953565b6118ed565b611730565b6115c2565b611528565b6114f5565b611457565b6113a6565b61128d565b610aae565b610a38565b610a12565b610990565b610964565b61090c565b61083d565b610740565b61071a565b610426565b610352565b634e487b7160e01b600052604160045260246000fd5b604081019081106001600160401b038211176101dd57604052565b6101ac565b60c081019081106001600160401b038211176101dd57604052565b90601f801991011681019081106001600160401b038211176101dd57604052565b6040519061022d6040836101fd565b565b60409060231901126102575760405190610248826101c2565b60243582526044356020830152565b600080fd5b60409060631901126102575760405190610275826101c2565b60643582526084356020830152565b60409060a3190112610257576040519061029d826101c2565b60a435825260c4356020830152565b604090600319011261025757604051906102c5826101c2565b60043582526024356020830152565b604090604319011261025757604051906102ed826101c2565b60443582526064356020830152565b6040906101031901126102575760405190610316826101c2565b610104358252610124356020830152565b6040906101431901126102575760405190610341826101c2565b610144358252610164356020830152565b3461025757610120366003190112610257576004356103703661022f565b60016103c661037e3661025c565b61038736610284565b60e4359061010435928661039a89611df4565b50916103a58a611e13565b506103c06103ba6103b4611ebf565b95611ef7565b91611ef7565b936128f2565b15150361041c576103d6906120d6565b60175490600160401b8210156101dd576103f9826001610412940160175561216a565b90919082549060031b91821b91600019901b1916179055565b61041a611f55565b005b505061041a611f15565b3461025757604036600319011261025757600435602435906001198201918083116104af5760c09160209160009460405192610461846101e2565b84845284808501528460408501526060840152608083015260a082015260405192839161048e84846101fd565b83368437856005600019f1156104ab575160405190815260209150f35b5080fd5b6121bb565b6001600160401b0381116101dd5760051b60200190565b81601f82011215610257578035906104e2826104b4565b926104f060405194856101fd565b82845260208085019360061b8301019181831161025757602001925b82841061051a575050505090565b600060408584031261054e575060206040918251610537816101c2565b86358152828701358382015281520193019261050c565b80fd5b9080601f83011215610257578135610568816104b4565b9261057660405194856101fd565b81845260208085019260051b82010192831161025757602001905b82821061059e5750505090565b8135815260209182019101610591565b610180600319820112610257576004356001600160401b03811161025757816105d9916004016104cb565b916024356001600160401b03811161025757826105f8916004016104cb565b916044356001600160401b0381116102575781610617916004016104cb565b916064356001600160401b0381116102575782610636916004016104cb565b916084356001600160401b0381116102575781610655916004016104cb565b9160a4356001600160401b038111610257578261067491600401610551565b9160c4356001600160401b038111610257578161069391600401610551565b9160e4356001600160401b03811161025757826106b291600401610551565b91610104356001600160401b03811161025757816106d291600401610551565b91610124356001600160401b03811161025757826106f2916004016104cb565b9161014435906001600160401b03821161025757610712916004016104cb565b906101643590565b346102575761041a61072b366105ae565b9a9990999891989792979693969594956122d3565b34610257576080366003190112610257576004356024356001600160401b038111610257576107739036906004016104cb565b90604435906001600160401b0382116102575760409261079a6107a4933690600401610551565b906064359261239c565b6107ba8251809260208091805184520151910152565bf35b9080601f830112156102575760408051926107d782856101fd565b8391810192831161025757905b8282106107f15750505090565b81358152602091820191016107e4565b906080610103198301126102575760405161081b816101c2565b6020610838829461082e816101046107bc565b84526101446107bc565b910152565b346102575761022036600319011261025757610858366102ac565b610861366102d4565b906080366083190112610257576040519161087b836101c2565b6108863660846107bc565b83526108933660c46107bc565b60208401526108a136610801565b9261018435936001600160401b038511610257576108c661041a9536906004016104cb565b6101a435916101c435936101e4359561020435976125a8565b602060031982011261025757600435906001600160401b03821161025757610909916004016104cb565b90565b346102575761091a366108df565b60005b815181101561041a5761093081836122bf565b5190601154600160401b8110156101dd57600192610958828561095e94016011556011611ea3565b906120aa565b0161091d565b346102575760603660031901126102575760206109886044356024356004356126a7565b604051908152f35b346102575761099e366108df565b60005b81518110156109e557806109df6109d66109bc600194611e13565b506109d16109ca85886122bf565b5191611ef7565b613029565b61095883611e13565b016109a1565b5060005b815181101561041a5780610a0c610a02600193856122bf565b5161095883611e13565b016109e9565b346102575761041a610a23366105ae565b9a99909998919897929796939695949561273e565b3461025757610a46366108df565b60005b815181101561041a5780610a6c610a62600193856122bf565b5161095883611e32565b01610a49565b602060408183019282815284518094520192019060005b818110610a965750505090565b82511515845260209384019390920191600101610a89565b3461025757600036600319011261025757604051601954808252906020810160196000527f944998273e477b495144fb8794c914197f3ccb46be2900f4698fd0ef743c9695926000935b81601f860110610ff25791610c21948492610c15945491818110610fd6575b818110610fb7575b818110610f98575b818110610f79575b818110610f5b575b818110610f3c575b818110610f1d575b818110610efe575b818110610edf575b818110610ec0575b818110610ea1575b818110610e82575b818110610e63575b818110610e44575b818110610e25575b818110610e06575b818110610de7575b818110610dc8575b818110610da9575b818110610d8a575b818110610d6b575b818110610d4c575b818110610d2d575b818110610d0e575b818110610cef575b818110610cd0575b818110610cb1575b818110610c92575b818110610c73575b818110610c54575b818110610c35575b10610c25575b5003826101fd565b60405191829182610a72565b0390f35b60f81c1515815260200138610c0d565b92602081610c4c60019360ff8760f01c1615159052565b019301610c07565b92602081610c6b60019360ff8760e81c1615159052565b019301610bff565b92602081610c8a60019360ff8760e01c1615159052565b019301610bf7565b92602081610ca960019360ff8760d81c1615159052565b019301610bef565b92602081610cc860019360ff8760d01c1615159052565b019301610be7565b92602081610ce760019360ff8760c81c1615159052565b019301610bdf565b92602081610d0660019360ff8760c01c1615159052565b019301610bd7565b92602081610d2560019360ff8760b81c1615159052565b019301610bcf565b92602081610d4460019360ff8760b01c1615159052565b019301610bc7565b92602081610d6360019360ff8760a81c1615159052565b019301610bbf565b92602081610d8260019360ff8760a01c1615159052565b019301610bb7565b92602081610da160019360ff8760981c1615159052565b019301610baf565b92602081610dc060019360ff8760901c1615159052565b019301610ba7565b92602081610ddf60019360ff8760881c1615159052565b019301610b9f565b92602081610dfe60019360ff8760801c1615159052565b019301610b97565b92602081610e1d60019360ff8760781c1615159052565b019301610b8f565b92602081610e3c60019360ff8760701c1615159052565b019301610b87565b92602081610e5b60019360ff8760681c1615159052565b019301610b7f565b92602081610e7a60019360ff8760601c1615159052565b019301610b77565b92602081610e9960019360ff8760581c1615159052565b019301610b6f565b92602081610eb860019360ff8760501c1615159052565b019301610b67565b92602081610ed760019360ff8760481c1615159052565b019301610b5f565b92602081610ef660019360ff8760401c1615159052565b019301610b57565b92602081610f1560019360ff8760381c1615159052565b019301610b4f565b92602081610f3460019360ff8760301c1615159052565b019301610b47565b92602081610f5360019360ff8760281c1615159052565b019301610b3f565b92602081610f7160019360ff87851c1615159052565b019301610b37565b92602081610f9060019360ff8760181c1615159052565b019301610b2f565b92602081610faf60019360ff8760101c1615159052565b019301610b27565b92602081610fce60019360ff8760081c1615159052565b019301610b1f565b92602081610fea60019360ff871615159052565b019301610b17565b916001610400602092611282865461100e8360ff831615159052565b61102186840160ff8360081c1615159052565b6110356040840160ff8360101c1615159052565b6110496060840160ff8360181c1615159052565b80861c60ff161515608084015261106a60a0840160ff8360281c1615159052565b61107e60c0840160ff8360301c1615159052565b61109260e0840160ff8360381c1615159052565b6110a7610100840160ff8360401c1615159052565b6110bc610120840160ff8360481c1615159052565b6110d1610140840160ff8360501c1615159052565b6110e6610160840160ff8360581c1615159052565b6110fb610180840160ff8360601c1615159052565b6111106101a0840160ff8360681c1615159052565b6111256101c0840160ff8360701c1615159052565b61113a6101e0840160ff8360781c1615159052565b61114f610200840160ff8360801c1615159052565b611164610220840160ff8360881c1615159052565b611179610240840160ff8360901c1615159052565b61118e610260840160ff8360981c1615159052565b6111a3610280840160ff8360a01c1615159052565b6111b86102a0840160ff8360a81c1615159052565b6111cd6102c0840160ff8360b01c1615159052565b6111e26102e0840160ff8360b81c1615159052565b6111f7610300840160ff8360c01c1615159052565b61120c610320840160ff8360c81c1615159052565b611221610340840160ff8360d01c1615159052565b611236610360840160ff8360d81c1615159052565b61124b610380840160ff8360e01c1615159052565b6112606103a0840160ff8360e81c1615159052565b6112756103c0840160ff8360f01c1615159052565b60f81c15156103e0830152565b019301940193610af8565b3461025757610120366003190112610257576112a8366102ac565b6044356001600160401b038111610257576112c79036906004016104cb565b906064356001600160401b038111610257576112e79036906004016104cb565b906084356001600160401b038111610257576113079036906004016104cb565b60a4356001600160401b038111610257576113269036906004016104cb565b60c4356001600160401b038111610257576113459036906004016104cb565b9060e4356001600160401b03811161025757611365903690600401610551565b9261010435956001600160401b03871161025757610c219761138e611394983690600401610551565b9661283e565b60405190151581529081906020820190565b346102575760003660031901126102575760005b60145480156114525760146000527fce6d7b5282bd9a3661ae061feed1dbda4e52ab073b1f9285be6e155d9c38d4ec5482101561041a576001916113fc61021e565b600081526000602082015260005b838110611423575061141d9192506120fb565b016113ba565b9061144b859161144561143f8661143987612185565b50611ea3565b50611ef7565b90613029565b910161140a565b611dde565b6101c03660031901126102575761146d366102ac565b611476366102d4565b9060403660831901126102575760405190611490826101c2565b608435825260a435602083015260403660c3190112610257576020926114eb926040516114bc816101c2565b60c435815260e435868201526114d1366102fc565b906114db36610327565b9261018435946101a435966128f2565b6040519015158152f35b346102575760403660031901126102575761041a611520611515366102ac565b6109d161143f611e51565b610958611e51565b346102575760003660031901126102575760005b60125480156114525760126000527fbb8a6a4669ba250d26cd7a459eca9d215f8307e33aebe50379bc5a3617ec34445482101561041a5760019161157e61021e565b600081526000602082015260005b8381106115a5575061159f919250612120565b0161153c565b906115bb859161144561143f86611439876121a0565b910161158c565b3461025757600036600319011261025757604051601b548082529060208101601b6000527f3ad8aa4f87544323a9d1e5dd902f40c356527a7955687113db5f9a85ad579dc1926000935b81601f8601106117095791610c21948492610c15945491818110610fd657818110610fb757818110610f9857818110610f7957818110610f5b57818110610f3c57818110610f1d57818110610efe57818110610edf57818110610ec057818110610ea157818110610e8257818110610e6357818110610e4457818110610e2557818110610e0657818110610de757818110610dc857818110610da957818110610d8a57818110610d6b57818110610d4c57818110610d2d57818110610d0e57818110610cef57818110610cd057818110610cb157818110610c9257818110610c7357818110610c5457818110610c355710610c25575003826101fd565b916001610400602092611725865461100e8360ff831615159052565b01930194019361160c565b34610257576101c0366003190112610257576004356001600160401b038111610257576117619036906004016104cb565b6024356001600160401b038111610257576117809036906004016104cb565b6044356001600160401b0381116102575761179f9036906004016104cb565b6064356001600160401b038111610257576117be9036906004016104cb565b6084356001600160401b038111610257576117dd9036906004016104cb565b60a4356001600160401b038111610257576117fc903690600401610551565b60c4356001600160401b0381116102575761181b903690600401610551565b60e4356001600160401b0381116102575761183a903690600401610551565b90610104356001600160401b0381116102575761185b903690600401610551565b92610124356001600160401b0381116102575761187c9036906004016104cb565b94610144356001600160401b0381116102575761189d903690600401610551565b96610164356001600160401b038111610257576118be9036906004016104cb565b98610184359a6001600160401b038c1161025757610c219c6118e76113949d3690600401610551565b50612a16565b346102575760206114eb611900366108df565b612c56565b602060408183019282815284518094520192019060005b8181106119295750505090565b9091926020604082611948600194885160208091805184520151910152565b01940192910161191c565b3461025757600036600319011261025757601554611970816104b4565b9061197e60405192836101fd565b808252601560009081527f55f448fdea98c4d29eb340757ef0a66cd03dbb9538908a6a81d96026b71ec475602084015b8383106119c35760405180610c218782611905565b600260206001926119d385611ef7565b8152019201920191906119ae565b3461025757600036600319011261025757604051601a548082529060208101601a6000527f057c384a7d1c54f3a1b2e5e67b2617b8224fdfd1ea7234eea573a6ff665ff63e926000935b81601f860110611b285791610c21948492610c15945491818110610fd657818110610fb757818110610f9857818110610f7957818110610f5b57818110610f3c57818110610f1d57818110610efe57818110610edf57818110610ec057818110610ea157818110610e8257818110610e6357818110610e4457818110610e2557818110610e0657818110610de757818110610dc857818110610da957818110610d8a57818110610d6b57818110610d4c57818110610d2d57818110610d0e57818110610cef57818110610cd057818110610cb157818110610c9257818110610c7357818110610c5457818110610c355710610c25575003826101fd565b916001610400602092611b44865461100e8360ff831615159052565b019301940193611a2b565b9060c0600319830112610257576004356001600160401b0381116102575782611b7a916004016104cb565b916024356001600160401b0381116102575781611b99916004016104cb565b916044356001600160401b0381116102575782611bb8916004016104cb565b916064356001600160401b0381116102575781611bd7916004016104cb565b916084356001600160401b0381116102575782611bf691600401610551565b9160a435906001600160401b0382116102575761090991600401610551565b346102575761041a611c2636611b4f565b94939093929192612de1565b3461025757608036600319011261025757600435602435611c52826121ed565b91611c5c81612e3e565b9260005b828110611cdf5750611c7183612e3e565b60005b848110611c7d57005b80611c94858589611c8f600196612f2e565b61239c565b611c9e82856122bf565b52611ca981846122bf565b50611cd9611cd0611cb983611e32565b506109d16103ba611cca86896122bf565b516130cc565b61095883611e32565b01611c74565b80611d00611cfb611cf160019461216a565b90549060031b1c90565b61221f565b611d0a82856122bf565b52611d1761143f82611e84565b611d2182886122bf565b52611d2c81876122bf565b5001611c60565b346102575761041a611d4436611b4f565b94939093929192612e8e565b3461025757600036600319011261025757601654611d6d816104b4565b90611d7b60405192836101fd565b808252601660009081527fd833147d7dc355ba459fc788f669e58cfaf9dc25ddcd0702e87d69c7b5124289602084015b838310611dc05760405180610c218782611905565b60026020600192611dd085611ef7565b815201920192019190611dab565b634e487b7160e01b600052603260045260246000fd5b60115481101561145257601160005260206000209060011b0190600090565b60155481101561145257601560005260206000209060011b0190600090565b60165481101561145257601660005260206000209060011b0190600090565b6015541561145257601560009081527f55f448fdea98c4d29eb340757ef0a66cd03dbb9538908a6a81d96026b71ec47591565b60135481101561145257601360005260206000209060011b0190600090565b80548210156114525760005260206000209060011b0190600090565b60405190611ecc826101c2565b60005482526001546020830152565b60405190611ee8826101c2565b60025482526003546020830152565b90604051611f04816101c2565b602060018294805484520154910152565b601b54600160401b8110156101dd576001810180601b5581101561145257601b60005260206000208160051c019060ff60f883549260031b161b19169055565b601b54600160401b8110156101dd576001810180601b5581101561145257601b60005260206000208160051c019060f882549160031b169060ff6001831b921b1916179055565b601954600160401b8110156101dd57600181018060195581101561145257601960005260206000208160051c019060ff60f883549260031b161b19169055565b601954600160401b8110156101dd57600181018060195581101561145257601960005260206000208160051c019060f882549160031b169060ff6001831b921b1916179055565b601a54600160401b8110156101dd576001810180601a5581101561145257601a60005260206000208160051c019060ff60f883549260031b161b19169055565b601a54600160401b8110156101dd576001810180601a5581101561145257601a60005260206000208160051c019060f882549160031b169060ff6001831b921b1916179055565b91906120c0576020816001925184550151910155565b634e487b7160e01b600052600060045260246000fd5b60135490600160401b8210156101dd5761095882600161022d94016013556013611ea3565b60165490600160401b8210156101dd5761095882600161022d94016016556016611ea3565b60155490600160401b8210156101dd5761095882600161022d94016015556015611ea3565b600c5490600160401b8210156101dd5761095882600161022d9401600c55600c611ea3565b60175481101561145257601760005260206000200190600090565b60145481101561145257601460005260206000200190600090565b60125481101561145257601260005260206000200190600090565b634e487b7160e01b600052601160045260246000fd5b6000198101919082116104af57565b919082039182116104af57565b906121f7826104b4565b61220460405191826101fd565b8281528092612215601f19916104b4565b0190602036910137565b90600182018092116104af57565b90600282018092116104af57565b90600382018092116104af57565b90600482018092116104af57565b90600582018092116104af57565b919082018092116104af57565b8051156114525760200190565b8051600110156114525760400190565b8051600210156114525760600190565b8051600310156114525760800190565b8051600410156114525760a00190565b80518210156114525760209160051b010190565b9b9a989796959493929190986122e8816121ed565b9060005b81811061234d57505050600f54986123038a6121ed565b9960005b81811061232f57505061231a9b9c612a16565b156123275761022d611fdc565b61022d611f9c565b808c61234682612340600195612f2e565b926122bf565b5201612307565b60018101908181116104af5760019161236682866122bf565b52016122ec565b6040519061237a826101c2565b60006020838281520152565b634e487b7160e01b600052601260045260246000fd5b90926123a661236d565b506123b0816121ed565b92600080937f30644e72e131a029b85045b68181585d2833e84879b9709143e1f593efffffff915b848110612437575050505050506123f56123f0611ebf565b612f9c565b916000925b82518410156124305761242860019161144561241687866122bf565b5161242188886122bf565b5190612fe2565b9301926123fa565b9250505090565b85906001848160005b89811061248157509061245291612ec7565b9261247c576001926000805160206134918339815191529109612475828a6122bf565b52016123d8565b612386565b91509192935080840361249d575b600101889392918691612440565b9190976124b46124ad84896122bf565b5186612f4e565b9061247c5760008051602061349183398151915291099060016000805160206134918339815191526124fb6124e9868a6122bf565b516124f4858b6122bf565b5190612f4e565b60009a0991905061248f565b906006820291808304600614901517156104af57565b9060005b6002811061252e57505050565b600190602083519301928185015501612521565b805160005b60028110612560575050602061022d910151600661251d565b60019060208351930192816004015501612547565b805160005b60028110612593575050602061022d910151600a61251d565b6001906020835193019281600801550161257a565b926125e66125f596936125e16125eb946125d36125f0989e9c9b999e60209080516000550151600155565b805160025560200151600355565b612542565b612575565b600d55565b600e55565b6125fe81600f55565b61260783601055565b60005b845181101561262f5780612629612623600193886122bf565b51612145565b0161260a565b509190925060005b81811061267257505060005b81811061264e575050565b60019061266c61265c61021e565b60008152600060208201526120fb565b01612643565b60019061269061268061021e565b6000815260006020820152612120565b01612637565b600160ff1b81146104af5760000390565b600192600019810190811384166104af576126c190612696565b828113156126cf5750505090565b8082036126ed575b6001600160ff1b0381146104af576001016126c1565b926000600080516020613491833981519152809261273561270e8887613079565b7f30644e72e131a029b85045b68181585d2833e84879b9709143e1f593efffffff90612ec7565b900908926126d7565b97959391899c9b9a99979593919a612755816121ed565b60005b82811061280657505050600f549a61276f8c6121ed565b9960005b8d81106127d357506127859c50612a16565b156127c957612792611fdc565b60005b82518110156127c457806127be611cd06127b0600194611e32565b506109d16109ca85896122bf565b01612795565b509050565b905061022d611f9c565b8b6001929496989a9c9d9e939597999b506127f18261234081612f2e565b5201918e99979593919c9b9a98969492612773565b60019192939597999b9c9d9496989a5061281f8161221f565b61282982856122bf565b5201908e999795939c9b9a9896949291612758565b9491969390959261284e87612c56565b61285f575050505050505050600090565b60005b82518110156128e4576128c98787878787878f8f896128bb816128b4816128ad816128a68161289f81612898816128c29d6122bf565b519c6122bf565b519b6122bf565b519a6122bf565b51996122bf565b51986122bf565b51976122bf565b51966128f2565b156128d657600101612862565b505050505050505050600090565b505050505050505050600190565b8661291f89612919612931989a61291361292b989d612925979a989a612fe2565b98612fe2565b94612fe2565b97612fe2565b92613029565b93613029565b91835181511493841594612981575b50508215612974575b8215612960575b505061295b57600190565b600090565b602091925081015191015114153880612950565b8051825114159250612949565b602090810151910151141592503880612940565b60405191906000835b600282106129b45750505061022d6040836101fd565b600160208192855481520193019101909161299e565b604051906129d7826101c2565b816129e26004612995565b815260206108386006612995565b604051906129fd826101c2565b81612a086008612995565b81526020610838600a612995565b96999a98919597929390949860005b8c51811015612c2d57612a65612a4a8c61242184612a4381886122bf565b51926122bf565b611445612a5784896122bf565b51612a60611edb565b612fe2565b612a6f82856122bf565b515181511490811591612c10575b5015612a98575050505050505050505050505061295b611f9c565b612ad18a611445612a57848f8a611445612ac884612ac2612ac2958f83612a4381612421936122bf565b936122bf565b51612a60611ebf565b612adb82886122bf565b515181511490811591612bf3575b5015612b04575050505050505050505050505061295b611f9c565b80808d818d8f828f91612b1782846122bf565b5191612b22916122bf565b51612b2c91612fe2565b94612b36916122bf565b5191612b41916122bf565b51612b4b91612fe2565b612b54906130cc565b92612b5e916122bf565b51612b67611ebf565b90612b7191612fe2565b612b7b848b6122bf565b51612b85906130cc565b92612b90858d6122bf565b51612b9a906130cc565b90612ba36129ca565b93612bac6129ca565b90612bb56129f0565b91612bbe6129ca565b94612bc76129ca565b97612bd199613208565b15612bde57600101612a25565b5050505050505050505050505061295b611f9c565b9050602080612c02848a6122bf565b510151910151141538612ae9565b9050602080612c1f84876122bf565b510151910151141538612a7d565b50505050505050505050505050612c42611fdc565b600190565b60001981146104af5760010190565b60105491600f5492612c6883516121ed565b9160015b82811115612d2d5750815b612c818684612265565b811015612cca57806000805160206134918339815191526000612cb88987612cb3612cae826001996121e0565b612696565b6126a7565b08612cc382876122bf565b5201612c77565b50935050612cd6611edb565b916000925b8151841015612cff57612cf760019161144561241687866122bf565b930192612cdb565b9250505080516002541490811591612d1c575b5061295b57600190565b602091500151600354141538612d12565b806000805160206134918339815191526000612d4d8987612d66966126a7565b08612d60612d5a836121d1565b876122bf565b52612c47565b612c6c565b60115490612d78826104b4565b91612d8660405193846101fd565b808352601160009081527f31ecc21a745e3968a04e9570e4425bc18fa8019c68028196b546d1669c200c68602085015b838310612dc35750505050565b60026020600192612dd385611ef7565b815201920192019190612db6565b908096959391612e039593612df4611edb565b91612dfd612d6b565b9261283e565b15612e3457612e10612063565b60005b82518110156127c45780612e2e6109d66127b0600194611e13565b01612e13565b905061022d612023565b90612e48826104b4565b612e5560405191826101fd565b8281528092612e66601f19916104b4565b019060005b828110612e7757505050565b602090612e8261236d565b82828501015201612e6b565b612e9d95949392612df4611edb565b15612eaa5761022d612063565b61022d612023565b604051906020612ec281846101fd565b368337565b60c09160405191612ed7836101e2565b60208352602080840152602060408401526060830152608082015260008051602061349183398151915260a0820152602090604051928391612f1984846101fd565b8336843760006005600019f115610257575190565b6000805160206134918339815191529081038181116104af576000900890565b81811115612f7757905b81039081116104af576000600080516020613491833981519152910890565b60008051602061349183398151915281018091116104af5790612f58565b1561025757565b90612fa561236d565b91826080606092602060405191612fbc86846101fd565b8536843780518352015160208201526000604082015260076107cf195a01fa1561025757565b9060809291612fef61236d565b93849160609360206040519261300587856101fd565b863685378051845201516020830152604082015260076107cf195a01fa1561025757565b6020929160c060609261303a61236d565b958693816040519361304d6080866101fd565b6080368637805185520151828401528051604084015201518482015260066107cf195a01fa1561025757565b600082820392128183128116918313901516176104af576000811261309b5790565b6130a490612696565b6000805160206134918339815191520360008051602061349183398151915281116104af5790565b6130d461236d565b5080511580613186575b61316b577f30644e72e131a029b85045b68181585d97816a916871ca8d3c208c16d87cfd4760208251920151067f30644e72e131a029b85045b68181585d97816a916871ca8d3c208c16d87cfd47037f30644e72e131a029b85045b68181585d97816a916871ca8d3c208c16d87cfd4781116104af5760405191613161836101c2565b8252602082015290565b50604051613178816101c2565b600081526000602082015290565b506020810151156130de565b604051906131a160c0836101fd565b6005825281601f196131b360056104b4565b019060409060005b8381106131c85750505050565b60209083516131d6816101c2565b84516131e286826101fd565b85368237815284516131f486826101fd565b8536823783820152828285010152016131bb565b9790919998959694939661321c60056104b4565b9761322a604051998a6101fd565b60058952601f1961323b60056104b4565b0160005b8181106133235750506109099a9b613255613192565b9a61325f8b612272565b526132698a612272565b506132738a61227f565b5261327d8961227f565b506132878961228f565b526132918861228f565b5061329b8861229f565b526132a58761229f565b506132af876122af565b526132b9866122af565b506132c387612272565b526132cd86612272565b506132d78661227f565b526132e18561227f565b506132eb8561228f565b526132f58461228f565b506132ff8461229f565b526133098361229f565b50613313836122af565b5261331d826122af565b5061333b565b808b6020809361333161236d565b920101520161323f565b805161334983518214612f95565b61335281612507565b9261335c846121ed565b9260005b8381106133975750505050602080926133919261337b612eb2565b94859260051b910160086107cf195a01fa612f95565b51151590565b806133a4600192846122bf565b51516133af82612507565b906133bc600092896122bf565b5260206133c983866122bf565b5101516133e16133db611cfb85612507565b896122bf565b526133ec82866122bf565b5151516134036133db6133fe85612507565b61222d565b5261341961341183876122bf565b515160200190565b5161342e6133db61342985612507565b61223b565b52602061343b83876122bf565b51015190505161345b61345561345084612507565b612249565b886122bf565b52613474602061346b83876122bf565b51015160200190565b5161348961345561348484612507565b612257565b520161336056fe30644e72e131a029b85045b68181585d2833e84879b9709143e1f593f0000001a2646970667358221220123b25051f4585c8836fab03e23a7ec5267e4ae0d25a8102ba82551c1af47d2764736f6c634300081c0033",
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

// AggregateCiphertext is a paid mutator transaction binding the contract method 0x8823f8a5.
//
// Solidity: function AggregateCiphertext() returns()
func (_Contract *ContractTransactor) AggregateCiphertext(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Contract.contract.Transact(opts, "AggregateCiphertext")
}

// AggregateCiphertext is a paid mutator transaction binding the contract method 0x8823f8a5.
//
// Solidity: function AggregateCiphertext() returns()
func (_Contract *ContractSession) AggregateCiphertext() (*types.Transaction, error) {
	return _Contract.Contract.AggregateCiphertext(&_Contract.TransactOpts)
}

// AggregateCiphertext is a paid mutator transaction binding the contract method 0x8823f8a5.
//
// Solidity: function AggregateCiphertext() returns()
func (_Contract *ContractTransactorSession) AggregateCiphertext() (*types.Transaction, error) {
	return _Contract.Contract.AggregateCiphertext(&_Contract.TransactOpts)
}

// AggregateEncShares is a paid mutator transaction binding the contract method 0xb002450b.
//
// Solidity: function AggregateEncShares() returns()
func (_Contract *ContractTransactor) AggregateEncShares(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Contract.contract.Transact(opts, "AggregateEncShares")
}

// AggregateEncShares is a paid mutator transaction binding the contract method 0xb002450b.
//
// Solidity: function AggregateEncShares() returns()
func (_Contract *ContractSession) AggregateEncShares() (*types.Transaction, error) {
	return _Contract.Contract.AggregateEncShares(&_Contract.TransactOpts)
}

// AggregateEncShares is a paid mutator transaction binding the contract method 0xb002450b.
//
// Solidity: function AggregateEncShares() returns()
func (_Contract *ContractTransactorSession) AggregateEncShares() (*types.Transaction, error) {
	return _Contract.Contract.AggregateEncShares(&_Contract.TransactOpts)
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

// PVerify is a paid mutator transaction binding the contract method 0x0c65d6e2.
//
// Solidity: function PVerify(uint256 indexTallier, (uint256,uint256) DecShare, (uint256,uint256) a1, (uint256,uint256) a2, uint256 challenge, uint256 z) returns()
func (_Contract *ContractTransactor) PVerify(opts *bind.TransactOpts, indexTallier *big.Int, DecShare VerificationG1Point, a1 VerificationG1Point, a2 VerificationG1Point, challenge *big.Int, z *big.Int) (*types.Transaction, error) {
	return _Contract.contract.Transact(opts, "PVerify", indexTallier, DecShare, a1, a2, challenge, z)
}

// PVerify is a paid mutator transaction binding the contract method 0x0c65d6e2.
//
// Solidity: function PVerify(uint256 indexTallier, (uint256,uint256) DecShare, (uint256,uint256) a1, (uint256,uint256) a2, uint256 challenge, uint256 z) returns()
func (_Contract *ContractSession) PVerify(indexTallier *big.Int, DecShare VerificationG1Point, a1 VerificationG1Point, a2 VerificationG1Point, challenge *big.Int, z *big.Int) (*types.Transaction, error) {
	return _Contract.Contract.PVerify(&_Contract.TransactOpts, indexTallier, DecShare, a1, a2, challenge, z)
}

// PVerify is a paid mutator transaction binding the contract method 0x0c65d6e2.
//
// Solidity: function PVerify(uint256 indexTallier, (uint256,uint256) DecShare, (uint256,uint256) a1, (uint256,uint256) a2, uint256 challenge, uint256 z) returns()
func (_Contract *ContractTransactorSession) PVerify(indexTallier *big.Int, DecShare VerificationG1Point, a1 VerificationG1Point, a2 VerificationG1Point, challenge *big.Int, z *big.Int) (*types.Transaction, error) {
	return _Contract.Contract.PVerify(&_Contract.TransactOpts, indexTallier, DecShare, a1, a2, challenge, z)
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

// Tally is a paid mutator transaction binding the contract method 0xe7db3bb2.
//
// Solidity: function Tally(uint256 threshold, uint256 numCandidates, uint256 a, uint256 b) returns()
func (_Contract *ContractTransactor) Tally(opts *bind.TransactOpts, threshold *big.Int, numCandidates *big.Int, a *big.Int, b *big.Int) (*types.Transaction, error) {
	return _Contract.contract.Transact(opts, "Tally", threshold, numCandidates, a, b)
}

// Tally is a paid mutator transaction binding the contract method 0xe7db3bb2.
//
// Solidity: function Tally(uint256 threshold, uint256 numCandidates, uint256 a, uint256 b) returns()
func (_Contract *ContractSession) Tally(threshold *big.Int, numCandidates *big.Int, a *big.Int, b *big.Int) (*types.Transaction, error) {
	return _Contract.Contract.Tally(&_Contract.TransactOpts, threshold, numCandidates, a, b)
}

// Tally is a paid mutator transaction binding the contract method 0xe7db3bb2.
//
// Solidity: function Tally(uint256 threshold, uint256 numCandidates, uint256 a, uint256 b) returns()
func (_Contract *ContractTransactorSession) Tally(threshold *big.Int, numCandidates *big.Int, a *big.Int, b *big.Int) (*types.Transaction, error) {
	return _Contract.Contract.Tally(&_Contract.TransactOpts, threshold, numCandidates, a, b)
}

// TestAggregateC is a paid mutator transaction binding the contract method 0xa5ac3091.
//
// Solidity: function TestAggregateC((uint256,uint256) c) returns()
func (_Contract *ContractTransactor) TestAggregateC(opts *bind.TransactOpts, c VerificationG1Point) (*types.Transaction, error) {
	return _Contract.contract.Transact(opts, "TestAggregateC", c)
}

// TestAggregateC is a paid mutator transaction binding the contract method 0xa5ac3091.
//
// Solidity: function TestAggregateC((uint256,uint256) c) returns()
func (_Contract *ContractSession) TestAggregateC(c VerificationG1Point) (*types.Transaction, error) {
	return _Contract.Contract.TestAggregateC(&_Contract.TransactOpts, c)
}

// TestAggregateC is a paid mutator transaction binding the contract method 0xa5ac3091.
//
// Solidity: function TestAggregateC((uint256,uint256) c) returns()
func (_Contract *ContractTransactorSession) TestAggregateC(c VerificationG1Point) (*types.Transaction, error) {
	return _Contract.Contract.TestAggregateC(&_Contract.TransactOpts, c)
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

// TestAggregateUSet is a paid mutator transaction binding the contract method 0x73c03394.
//
// Solidity: function TestAggregateUSet((uint256,uint256)[] U) returns()
func (_Contract *ContractTransactor) TestAggregateUSet(opts *bind.TransactOpts, U []VerificationG1Point) (*types.Transaction, error) {
	return _Contract.contract.Transact(opts, "TestAggregateUSet", U)
}

// TestAggregateUSet is a paid mutator transaction binding the contract method 0x73c03394.
//
// Solidity: function TestAggregateUSet((uint256,uint256)[] U) returns()
func (_Contract *ContractSession) TestAggregateUSet(U []VerificationG1Point) (*types.Transaction, error) {
	return _Contract.Contract.TestAggregateUSet(&_Contract.TransactOpts, U)
}

// TestAggregateUSet is a paid mutator transaction binding the contract method 0x73c03394.
//
// Solidity: function TestAggregateUSet((uint256,uint256)[] U) returns()
func (_Contract *ContractTransactorSession) TestAggregateUSet(U []VerificationG1Point) (*types.Transaction, error) {
	return _Contract.Contract.TestAggregateUSet(&_Contract.TransactOpts, U)
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
