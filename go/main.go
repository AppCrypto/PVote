package main

import (
	//"context"
	"PVote/crypto/PVSS"
	"PVote/crypto/ZKRP"
	"bytes"
	"crypto/rand"
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

type AggregatedValue struct {
	C []*bn256.G1
	V []*bn256.G1
	U []*bn256.G1
}

func main() {
	numTalliers := 3
	threshold := numTalliers/2 + 1

	//ballot range
	a := big.NewInt(0)
	b := big.NewInt(2)

	//Candidates
	numCandidates := 3
	numVoters := 3

	w := make([][]*big.Int, numVoters)
	for j := 0; j < numVoters; j++ {
		w[j] = make([]*big.Int, numCandidates)
	}

	//Init the public parameters
	_, PP := ZKRP.Setup(int(a.Int64()), int(b.Int64()))
	// Publish all public parameters onto smart contract.
	//TODO:

	//All talliers register own key pairs
	SKs, PKs := PVSS.Setup(numTalliers, PP.G0)
	//Publish PKs onto smart contract.
	//TODO:

	fmt.Printf("=================================Finish the Setup phase==========================================\n")

	//The voter Vj votes for each candidate C(d)
	for j := 0; j < numVoters; j++ {
		for i := 0; i < numCandidates; i++ {
			temp, _ := rand.Int(rand.Reader, new(big.Int).Add(new(big.Int).Sub(b, a), big.NewInt(1)))
			w[j][i] = new(big.Int).Add(a, temp)

		}
	}
	fmt.Printf("All voting value is %v\n", w)
	/*
		If we need to vote for l candidates, we should generate (n+l) part g_0^{p_j(i)}
		The first l shares are used to blind the voting value
	*/
	//Generate the PVSS shares
	secret := make([]*big.Int, numVoters)
	PVSSShares := make([]*PVSS.SecretSharing, numVoters)
	for j := 0; j < numVoters; j++ {
		secret[j], _ = rand.Int(rand.Reader, bn256.Order)
		PVSSShares[j] = PVSS.Share(secret[j], PP.H0, PKs, threshold, numTalliers, numCandidates)
	}

	//Generate the voting ciphertext set U and corresponding range proofs
	//U_d=g_0^{p_j(d)}*h0^{w_d}, d in [0,l-1]
	U := make([][]*bn256.G1, numVoters)

	rangeProofs := make([][]*ZKRP.Proof, numVoters)
	for j := 0; j < numVoters; j++ {
		U[j] = make([]*bn256.G1, numCandidates)
		rangeProofs[j] = make([]*ZKRP.Proof, numCandidates)

		_coefficients := make([]*big.Int, threshold)
		_coefficients[0], _ = rand.Int(rand.Reader, bn256.Order)
		for i := 1; i < threshold; i++ {
			_coefficients[i], _ = rand.Int(rand.Reader, bn256.Order)
		}
		for d := 0; d < numCandidates; d++ {
			U[j][d] = new(bn256.G1).Add(new(bn256.G1).ScalarMult(PP.G0, PVSSShares[j].BindValue[d]), new(bn256.G1).ScalarMult(PP.H0, w[j][d]))
			x := new(big.Int).Neg(big.NewInt(int64(d)))
			x.Mod(x, bn256.Order)
			rangeProofs[j][d] = ZKRP.GenProof(PP.G0, PP.H0, PP.G1, PVSSShares[j].BindValue[d], w[j][d], U[j][d], PP.Sigma_k[new(big.Int).Sub(w[j][d], a).Int64()], x, _coefficients)
		}
	}

	//Upload PVSSShares.V, PVSSShares.C, PVSSShares.Proofs, voting ciphertexts U, the range proof of voting ciphertexts rangProofs onto smart contract
	//TODO:

	//Verify the range proof via ZKRPVerify smart contract
	//TODO:
	selectedIndices := make([]*big.Int, threshold)
	for i := 0; i < threshold; i++ {
		selectedIndices[i] = big.NewInt(int64(i + 1))
	}
	for j := 0; j < numVoters; j++ {
		selectedShares := PVSSShares[j].V[:threshold]
		for d := 0; d < numCandidates; d++ {
			x := new(big.Int).Neg(big.NewInt(int64(d)))
			x.Mod(x, bn256.Order)
			if ZKRP.Verify(PP.G0, PP.H0, PP.G1, PP.PKI, rangeProofs[j][d], U[j][d], x, selectedShares, selectedIndices, threshold) {
				fmt.Printf("The %vth ZKRP proof from %vth votor verification pass!\n", d, j)
			}
		}
	}

	//Aggregate all voters corresponding information including PVSSShares.V, PVSSShare.C
	//Aggrate smart contract executes this operation
	//TODO:
	aggregatedValueC := make([]*bn256.G1, numTalliers)
	//aggregatedValueV := make([]*bn256.G1, numTalliers)
	aggregatedValueU := make([]*bn256.G1, numCandidates)
	for i := 0; i < numTalliers; i++ {
		aggregatedValueC[i] = new(bn256.G1).ScalarBaseMult(big.NewInt(0))
		//aggregatedValueV[i] = new(bn256.G1).ScalarBaseMult(big.NewInt(0))
		aggregatedValueU[i] = new(bn256.G1).ScalarBaseMult(big.NewInt(0))
		for j := 0; j < numVoters; j++ {
			aggregatedValueC[i] = new(bn256.G1).Add(aggregatedValueC[i], PVSSShares[j].C[i])
			//aggregatedValueV[i] = new(bn256.G1).Add(aggregatedValueV[i], PVSSShares[j].V[i])
		}
	}

	for i := 0; i < numCandidates; i++ {
		aggregatedValueU[i] = new(bn256.G1).ScalarBaseMult(big.NewInt(0))
		for j := 0; j < numVoters; j++ {
			aggregatedValueU[i] = new(bn256.G1).Add(aggregatedValueU[i], U[j][i])
		}
	}

	aggregatedValue := new(AggregatedValue)
	aggregatedValue.C = aggregatedValueC
	//aggregatedValue.V = aggregatedValueV
	aggregatedValue.U = aggregatedValueU

	//Each Tallier decrypts the aggreated shares
	sh := make([]*bn256.G1, numTalliers)
	shProof := make([]PVSS.Proof, numTalliers)
	for i := 0; i < numTalliers; i++ {
		sh[i], shProof[i] = PVSS.Decrypt(PP.G0, PKs[i], aggregatedValue.C[i], SKs[i])
	}

	//PVerify smart contracts verifies the correctness of each decrypted shares
	//TODO:
	for i := 0; i < numTalliers; i++ {
		if PVSS.PVerify(PP.G0, PKs[i], aggregatedValue.C[i], sh[i], shProof[i]) {
			fmt.Printf("The %vth tallier is honest!\n", i)
		} else {
			fmt.Printf("The %vth tallier is not honest!\n", i)
		}
	}

	//Decrypt to get voting value
	//TODO:
	_selectedIndices := make([]*big.Int, threshold)
	_selectedShares := make([]*bn256.G1, threshold)
	for i := 0; i < threshold; i++ {
		_selectedIndices[i] = big.NewInt(int64(i + 1))
	}
	for j := 0; j < numVoters; j++ {
		_selectedShares = sh[:threshold]
		for d := 0; d < numCandidates; d++ {
			x := new(big.Int).Neg(big.NewInt(int64(d)))
			x.Mod(x, bn256.Order)
			if ZKRP.Verify(PP.G0, PP.H0, PP.G1, PP.PKI, rangeProofs[j][d], U[j][d], x, _selectedShares, _selectedIndices, threshold) {
				fmt.Printf("The %vth ZKRP proof from %vth votor verification pass!\n", d, j)
			}
		}
	}
	S := make([]*bn256.G1, numCandidates)
	for d := 0; d < numCandidates; d++ {
		x := new(big.Int).Neg(big.NewInt(int64(d)))
		x.Mod(x, bn256.Order)
		S[d] = ZKRP.Interpolation(x, _selectedShares, _selectedIndices, threshold)
		aggregatedValue.U[d] = new(bn256.G1).Add(aggregatedValue.U[d], new(bn256.G1).Neg(S[d]))
	}

	//Count the final voting value
	//TODO:
	bInt := b.Int64()
	for d := 0; d < numCandidates; d++ {
		for k := 0; k < int(bInt)*numCandidates; k++ {
			if bytes.Equal(aggregatedValue.U[d].Marshal(), new(bn256.G1).ScalarMult(PP.H0, big.NewInt(int64(k))).Marshal()) {
				fmt.Printf("The %dth voting value is %v\n", d, k)
				break
			}
		}
	}
}
