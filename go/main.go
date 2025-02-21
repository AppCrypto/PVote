package main

import (
	"PVote/crypto/Convert"
	"PVote/crypto/PVSS"
	"PVote/crypto/ZKRP"
	"bytes"
	"context"
	"crypto/rand"
	"fmt"

	"PVote/compile/contract"
	"PVote/utils"
	"log"
	"math/big"

	//"time"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	bn256 "github.com/ethereum/go-ethereum/crypto/bn256/cloudflare"
	"github.com/ethereum/go-ethereum/ethclient"
)

type AggregatedValue struct {
	C []*bn256.G1
	U []*bn256.G1
}

func main() {
	//Deploy the smart contract
	//Functions:
	contract_name := "Verification"
	client, err := ethclient.Dial("http://127.0.0.1:8545")
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}

	privatekey := utils.GetENV("PRIVATE_KEY_1")

	auth := utils.Transact(client, privatekey, big.NewInt(0))

	address, tx := utils.Deploy(client, contract_name, auth)

	receipt, _ := bind.WaitMined(context.Background(), client, tx)
	if err != nil {
		log.Fatalf("Tx receipt failed: %v", err)
	}
	fmt.Printf("Deploy Gas used: %d\n", receipt.GasUsed)

	Contract, err := contract.NewContract(common.HexToAddress(address.Hex()), client)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%v\n", Contract)

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
	bInt := b.Int64()
	aInt := a.Int64()
	sigmak := make([]contract.VerificationG1Point, bInt-aInt+1)
	for i := aInt; i < bInt-aInt+1; i++ {
		//Convert *bn256.G1 to G1Point
		sigmak[i-aInt] = Convert.G1ToG1Point(PP.Sigma_k[i-aInt])
	}

	auth0 := utils.Transact(client, privatekey, big.NewInt(0))
	tx0, _ := Contract.UploadParameters(auth0, Convert.G1ToG1Point(PP.G0), Convert.G1ToG1Point(PP.H0), Convert.G2ToG2Point(PP.G1), Convert.G2ToG2Point(PP.PKI), sigmak)

	receipt0, err := bind.WaitMined(context.Background(), client, tx0)
	if err != nil {
		log.Fatalf("Tx receipt failed: %v", err)
	}
	fmt.Printf("UploadParameters Gas used: %d\n", receipt0.GasUsed)

	//All talliers register own key pairs
	SKs, PKs := PVSS.Setup(numTalliers, PP.G0)
	//Publish PKs onto smart contract.
	//TODO:
	pks := make([]contract.VerificationG1Point, numTalliers)
	for i := 0; i < numTalliers; i++ {
		//Convert *bn256.G1 to G1Point
		pks[i] = Convert.G1ToG1Point(PKs[i])
	}

	auth1 := utils.Transact(client, privatekey, big.NewInt(0))
	tx1, _ := Contract.UploadPublicKey(auth1, pks)

	receipt1, err := bind.WaitMined(context.Background(), client, tx1)
	if err != nil {
		log.Fatalf("Tx receipt failed: %v", err)
	}
	fmt.Printf("UploadPublicKey Gas used: %d\n", receipt1.GasUsed)

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

	//Verify the PVSS Shares
	for j := 0; j < numVoters; j++ {
		result := PVSS.DVerify(PVSSShares[j], PP.H0, PKs)
		if result {
			fmt.Printf("The %vth voter is honest\n", j)
		} else {
			fmt.Printf("The %vth voter is not honest\n", j)
		}
	}

	//Upload PVSSShares.V, PVSSShares.C, PVSSShares.Proofs
	//TODO:
	vSet := make([][]contract.VerificationG1Point, numVoters)
	cSet := make([][]contract.VerificationG1Point, numVoters)
	a1Set := make([][]contract.VerificationG1Point, numVoters)
	a2Set := make([][]contract.VerificationG1Point, numVoters)
	challengeSet := make([][]*big.Int, numVoters)
	zSet := make([][]*big.Int, numVoters)
	for j := 0; j < numVoters; j++ {
		vSet[j] = make([]contract.VerificationG1Point, numTalliers)
		cSet[j] = make([]contract.VerificationG1Point, numTalliers)
		a1Set[j] = make([]contract.VerificationG1Point, numTalliers)
		a2Set[j] = make([]contract.VerificationG1Point, numTalliers)
		challengeSet[j] = make([]*big.Int, numTalliers)
		zSet[j] = make([]*big.Int, numTalliers)
	}

	for j := 0; j < numVoters; j++ {
		for i := 0; i < numTalliers; i++ {
			vSet[j][i] = Convert.G1ToG1Point(PVSSShares[j].V[i])
			cSet[j][i] = Convert.G1ToG1Point(PVSSShares[j].C[i])
			a1Set[j][i] = Convert.G1ToG1Point(PVSSShares[j].Proofs[i].RG)
			a2Set[j][i] = Convert.G1ToG1Point(PVSSShares[j].Proofs[i].RH)
			challengeSet[j][i] = PVSSShares[j].Proofs[i].C
			zSet[j][i] = PVSSShares[j].Proofs[i].Z
		}
	}
	sumGasUsed := uint64(0)
	for j := 0; j < numVoters; j++ {
		auth2 := utils.Transact(client, privatekey, big.NewInt(0))
		tx2, _ := Contract.UploadPVSSShares(auth2, big.NewInt(int64(j)), vSet[j], cSet[j], a1Set[j], a2Set[j], challengeSet[j], zSet[j])

		receipt2, err := bind.WaitMined(context.Background(), client, tx2)
		if err != nil {
			log.Fatalf("Tx receipt failed: %v", err)
		}
		sumGasUsed = sumGasUsed + receipt2.GasUsed
	}
	fmt.Printf("UploadPVSSShares Gas used: %d\n", sumGasUsed)

	//Verify the PVSSShares via DVerify smart contract
	//TODO:
	DVerifyResult, _ := Contract.GetDVerifyResult(&bind.CallOpts{})
	fmt.Printf("The Verification results of DVerify is %v\n", DVerifyResult)

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

	// sumGasUsed = uint64(0)
	// for j := 0; j < numVoters; j++ {
	// 	auth4 := utils.Transact(client, privatekey, big.NewInt(0))
	// 	tx4, _ := Contract.UploadRangeProofs(auth4, ESet[j], F1Set[j], F2Set[j], _USet[j], _CSet[j], RPcSet[j], z1Set[j], z2Set[j], z3Set[j])

	// 	receipt4, err := bind.WaitMined(context.Background(), client, tx4)
	// 	if err != nil {
	// 		log.Fatalf("Tx receipt failed: %v", err)
	// 	}
	// 	sumGasUsed = sumGasUsed + receipt4.GasUsed
	// }
	// fmt.Printf("UploadRangeProofs Gas used: %d\n", sumGasUsed)

	//Verify the range proof via ZKRPVerify smart contract
	//TODO:
	selectedIndices := make([]*big.Int, threshold)
	x := make([]*big.Int, numVoters)
	for i := 0; i < threshold; i++ {
		selectedIndices[i] = big.NewInt(int64(i + 1))
	}
	for j := 0; j < numVoters; j++ {
		selectedShares := PVSSShares[j].V[:threshold]
		for d := 0; d < numCandidates; d++ {
			x[d] = new(big.Int).Neg(big.NewInt(int64(d)))
			x[d].Mod(x[d], bn256.Order)
			if ZKRP.Verify(PP.G0, PP.H0, PP.G1, PP.PKI, rangeProofs[j][d], U[j][d], x[d], selectedShares, selectedIndices, threshold) {
				fmt.Printf("The %vth ZKRP proof from %vth votor verification pass!\n", d, j)
			}
		}
	}

	//Upload voting ciphertexts U, the range proof of voting ciphertexts rangProofs onto smart contract
	//TODO:
	USet := make([][]contract.VerificationG1Point, numVoters)
	ESet := make([][]contract.VerificationG1Point, numVoters)
	F1Set := make([][]contract.VerificationG1Point, numVoters)
	F2Set := make([][]contract.VerificationG1Point, numVoters)
	_USet := make([][]contract.VerificationG1Point, numVoters)
	_CSet := make([][]contract.VerificationG1Point, numVoters)
	RPcSet := make([][]*big.Int, numVoters)
	z1Set := make([][]*big.Int, numVoters)
	z2Set := make([][]*big.Int, numVoters)
	z3Set := make([][]*big.Int, numVoters)
	selectedV := make([][]contract.VerificationG1Point, numVoters)

	for j := 0; j < numVoters; j++ {
		USet[j] = make([]contract.VerificationG1Point, numCandidates)
		ESet[j] = make([]contract.VerificationG1Point, numCandidates)
		F1Set[j] = make([]contract.VerificationG1Point, numCandidates)
		F2Set[j] = make([]contract.VerificationG1Point, numCandidates)
		_USet[j] = make([]contract.VerificationG1Point, numCandidates)
		_CSet[j] = make([]contract.VerificationG1Point, numCandidates)
		RPcSet[j] = make([]*big.Int, numCandidates)
		z1Set[j] = make([]*big.Int, numCandidates)
		z2Set[j] = make([]*big.Int, numCandidates)
		z3Set[j] = make([]*big.Int, numCandidates)
		selectedV[j] = make([]contract.VerificationG1Point, threshold)

		for i := 0; i < numCandidates; i++ {
			USet[j][i] = Convert.G1ToG1Point(U[j][i])
			ESet[j][i] = Convert.G1ToG1Point(rangeProofs[j][i].Ej)
			F1Set[j][i] = Convert.G1ToG1Point(rangeProofs[j][i].Fj1)
			F2Set[j][i] = Convert.G1ToG1Point(rangeProofs[j][i].Fj2)
			_USet[j][i] = Convert.G1ToG1Point(rangeProofs[j][i].Uj)
			_CSet[j][i] = Convert.G1ToG1Point(rangeProofs[j][i].Cj)
			RPcSet[j][i] = rangeProofs[j][i].C
			z1Set[j][i] = rangeProofs[j][i].Z1
			z2Set[j][i] = rangeProofs[j][i].Z2
			z3Set[j][i] = rangeProofs[j][i].Z3
		}

		for i := 0; i < threshold; i++ {
			selectedV[j][i] = Convert.G1ToG1Point(PVSSShares[j].V[i])
		}
	}

	sumGasUsed = uint64(0)
	for j := 0; j < numVoters; j++ {
		auth3 := utils.Transact(client, privatekey, big.NewInt(0))
		tx3, _ := Contract.UploadVotingCipher(auth3, ESet[j], F1Set[j], F2Set[j], _USet[j], _CSet[j], RPcSet[j], z1Set[j], z2Set[j], z3Set[j], USet[j], x, selectedV[j], selectedIndices, big.NewInt(int64(threshold)))

		receipt3, err := bind.WaitMined(context.Background(), client, tx3)
		if err != nil {
			log.Fatalf("Tx receipt failed: %v", err)
		}
		sumGasUsed = sumGasUsed + receipt3.GasUsed
	}
	fmt.Printf("UploadVotingvalue Gas used: %d\n", sumGasUsed)
	ZKPRVerifyResult, _ := Contract.GetDZKRPResult(&bind.CallOpts{})
	fmt.Printf("The Verification results of ZKRPVerify is %v\n", ZKPRVerifyResult)

	//Aggregate all voters corresponding information including PVSSShares.V, PVSSShare.C
	aggregatedValueC := make([]*bn256.G1, numTalliers)
	aggregatedValueU := make([]*bn256.G1, numCandidates)
	for i := 0; i < numTalliers; i++ {
		aggregatedValueC[i] = new(bn256.G1).ScalarBaseMult(big.NewInt(0))
		aggregatedValueU[i] = new(bn256.G1).ScalarBaseMult(big.NewInt(0))
		for j := 0; j < numVoters; j++ {
			aggregatedValueC[i] = new(bn256.G1).Add(aggregatedValueC[i], PVSSShares[j].C[i])
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
	//Aggrate smart contract executes this operation
	//TODO:
	auth5 := utils.Transact(client, privatekey, big.NewInt(0))
	tx5, _ := Contract.Aggregate(auth5)

	receipt5, err := bind.WaitMined(context.Background(), client, tx5)
	if err != nil {
		log.Fatalf("Tx receipt failed: %v", err)
	}
	fmt.Printf("Aggregate Gas used: %d\n", receipt5.GasUsed)

	fmt.Printf("=================================Finish the Voting phase==========================================\n")
	//Each Tallier decrypts the aggreated shares
	//AggregateResultC, _ := Contract.GetAggregateValue(&bind.CallOpts{})

	sh := make([]*bn256.G1, numTalliers)
	shProof := make([]PVSS.Proof, numTalliers)
	for i := 0; i < numTalliers; i++ {
		//sh[i], shProof[i] = PVSS.Decrypt(PP.G0, PKs[i], aggregatedValue.C[i], SKs[i])
		//sh[i], shProof[i] = PVSS.Decrypt(PP.G0, PKs[i], Convert.G1PointToG1(AggregateResultC[i]), SKs[i])
		sh[i], shProof[i] = PVSS.Decrypt(PP.G0, PKs[i], aggregatedValue.C[i], SKs[i])
	}

	//PVerify smart contracts verifies the correctness of each decrypted shares
	//TODO:
	//Upload decrypted shares and corresponding DLEQ proofs
	sumGasUsed = uint64(0)
	for i := 0; i < numTalliers; i++ {
		auth6 := utils.Transact(client, privatekey, big.NewInt(0))
		tx6, _ := Contract.PVerify(auth6, big.NewInt(int64(i)), Convert.G1ToG1Point(sh[i]), Convert.G1ToG1Point(shProof[i].RG), Convert.G1ToG1Point(shProof[i].RH), shProof[i].C, shProof[i].Z)

		receipt6, err := bind.WaitMined(context.Background(), client, tx6)
		if err != nil {
			log.Fatalf("Tx receipt failed: %v", err)
		}
		PVerifyResult, _ := Contract.GetPVerifyResult(&bind.CallOpts{})
		fmt.Printf("The Verification results of PVerifyResult is %v\n", PVerifyResult)
		if PVerifyResult[i] {
			fmt.Printf("The %vth tallier is honest!\n", i)
		} else {
			fmt.Printf("The %vth tallier is not honest!\n", i)
		}
		sumGasUsed = sumGasUsed + receipt6.GasUsed
	}
	fmt.Printf("UploadVotingvalue Gas used: %d\n", sumGasUsed)

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
	//_selectedShares := make([]*bn256.G1, threshold)
	for i := 0; i < threshold; i++ {
		_selectedIndices[i] = big.NewInt(int64(i + 1))
	}
	//for j := 0; j < numVoters; j++ {
	_selectedShares := sh[:threshold]
	// for d := 0; d < numCandidates; d++ {
	// 	x := new(big.Int).Neg(big.NewInt(int64(d)))
	// 	x.Mod(x, bn256.Order)
	// 	if ZKRP.Verify(PP.G0, PP.H0, PP.G1, PP.PKI, rangeProofs[j][d], U[j][d], x, _selectedShares, _selectedIndices, threshold) {
	// 		fmt.Printf("The %vth ZKRP proof from %vth votor verification pass!\n", d, j)
	// 	}
	// }
	//}
	S := make([]*bn256.G1, numCandidates)
	for d := 0; d < numCandidates; d++ {
		x := new(big.Int).Neg(big.NewInt(int64(d)))
		x.Mod(x, bn256.Order)
		S[d] = ZKRP.Interpolation(x, _selectedShares, _selectedIndices, threshold)
		aggregatedValue.U[d] = new(bn256.G1).Add(aggregatedValue.U[d], new(bn256.G1).Neg(S[d]))
	}

	//Count the final voting value
	//TODO:
	for d := 0; d < numCandidates; d++ {
		for k := 0; k < int(bInt)*numCandidates; k++ {
			if bytes.Equal(aggregatedValue.U[d].Marshal(), new(bn256.G1).ScalarMult(PP.H0, big.NewInt(int64(k))).Marshal()) {
				fmt.Printf("The %dth voting value is %v\n", d, k)
				break
			}
		}
	}

	auth7 := utils.Transact(client, privatekey, big.NewInt(0))
	tx7, _ := Contract.Tally(auth7, big.NewInt(int64(threshold)), big.NewInt(int64(numCandidates)), a, b)

	receipt7, err := bind.WaitMined(context.Background(), client, tx7)
	if err != nil {
		log.Fatalf("Tx receipt failed: %v", err)
	}
	fmt.Printf("Tally Gas used: %d\n", receipt7.GasUsed)
	TallyResult, _ := Contract.GetTallyValue(&bind.CallOpts{})
	fmt.Printf("The tally results are %v\n", TallyResult)

}
