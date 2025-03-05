package main

import (
	"PVote/crypto/Convert"
	"PVote/crypto/PVSS"
	"PVote/crypto/ZKRP"
	"bytes"
	"context"
	"crypto/rand"
	"fmt"
	"unsafe"

	"PVote/compile/contract"
	"PVote/utils"
	"log"
	"math/big"

	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	bn256 "github.com/ethereum/go-ethereum/crypto/bn256/cloudflare"
	"github.com/ethereum/go-ethereum/ethclient"
)

type AggregatedValue struct {
	C []*bn256.G1
	U []*bn256.G1
}

// Compute G1Point size
func sizeOfG1Point(p contract.VerificationG1Point) int {

	structSize := int(unsafe.Sizeof(p))
	dataSize := len(p.X.Bytes()) + len(p.Y.Bytes())

	return structSize + dataSize
}

// Compute *big.Int size
func sizeOfBigInt(n *big.Int) int {
	return len(n.Bytes())
}

func main() {
	//Setup Phase
	n := int64(1)
	//Init the public parameters
	//The algorithms in Setup phase: PVSS.Setup and ZKRP.Setup

	//Talliers:[4,6,8,10,12,14,16,18,20,22,24,26,28,30]
	numTalliers := 10
	//Candidates
	numCandidates := 1
	//Voters:[30,60,90,120,150,180,210,240,270,300]
	numVoters := 1
	threshold := (numTalliers + numCandidates) / 2
	fmt.Printf("numTalliers=%v\nnumCandidates=%v\nnumVoters=%v\nThreshold=%v\n", numTalliers, numCandidates, numVoters, threshold)

	//Deploy the smart contract
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

	//ballot range
	a := big.NewInt(0)
	b := big.NewInt(10)

	//w is the set of all ballots
	//w[i][j] denotes the ballot of i-th voter for j-th candidate
	w := make([][]*big.Int, numVoters)
	for j := 0; j < numVoters; j++ {
		w[j] = make([]*big.Int, numCandidates)
	}
	_, PP := ZKRP.Setup(int(a.Int64()), int(b.Int64()))
	// Publish all public parameters(PP, the number of candidates, and voting range [a,b]) onto smart contract.
	bInt := b.Int64()
	aInt := a.Int64()
	sigmak := make([]contract.VerificationG1Point, bInt-aInt+1)
	for i := 0; i < int(bInt-aInt+1); i++ {
		//Convert *bn256.G1 to G1Point
		sigmak[i] = Convert.G1ToG1Point(PP.Sigma_k[i])
	}

	auth0 := utils.Transact(client, privatekey, big.NewInt(0))
	tx0, _ := Contract.UploadParameters(auth0, Convert.G1ToG1Point(PP.G0), Convert.G1ToG1Point(PP.H0), Convert.G2ToG2Point(PP.G1), Convert.G2ToG2Point(PP.PKI), sigmak, a, b, big.NewInt(int64(numCandidates)), big.NewInt(int64(numTalliers)))
	receipt0, err := bind.WaitMined(context.Background(), client, tx0)
	if err != nil {
		log.Fatalf("Tx receipt failed: %v", err)
	}
	fmt.Printf("UploadParameters Gas used: %d\n", receipt0.GasUsed)
	//All talliers register own key pairs
	SKs, PKs := PVSS.Setup(numTalliers, PP.G0)
	//Publish the PKs onto smart contract.
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

	fmt.Printf("====================================Finish the Setup phase==========================================\n")

	/*
		Voting Phase
		1)voters vote;
		2)Smart contracts check voters' honesty;
		3)Smart contracts aggregates ecrypted shares;
	*/

	//The voter Vj votes for each candidate C(d)
	for j := 0; j < numVoters; j++ {
		for i := 0; i < numCandidates; i++ {
			temp, _ := rand.Int(rand.Reader, new(big.Int).Add(new(big.Int).Sub(b, a), big.NewInt(1)))
			w[j][i] = new(big.Int).Add(a, temp)

		}
	}
	fmt.Printf("All voting value is %v\n", w)

	/*
		Note:
		  If we need to vote for l candidates, we should generate (n+l) part g_0^{p_j(i)}
		  The first l shares are used to blind the voting value
	*/

	//Algorithm1(each voter Vj): 1 PVSS.Share+l ZKRP.GenProof
	//Each voter generates the PVSS shares
	secret := make([]*big.Int, numVoters)
	PVSSShares := make([]*PVSS.SecretSharing, numVoters)
	U := make([][]*bn256.G1, numVoters)
	rangeProofs := make([][]*ZKRP.Proof, numVoters)

	starttime := time.Now().UnixMicro()
	for k := 0; k < int(n); k++ {
		for j := 0; j < numVoters; j++ {
			secret[j], _ = rand.Int(rand.Reader, bn256.Order)
			PVSSShares[j] = PVSS.Share(secret[j], PP.H0, PKs, threshold, numTalliers, numCandidates)
			//Generate the voting ciphertext set U and corresponding range proofs
			//U_d=g_0^{p_j(d)}*h0^{w_d}, d in [0,l-1]
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
	}
	endtime := time.Now().UnixMicro()
	fmt.Printf("Algorithm1 Time Used is %v us\n", (endtime-starttime)/n)

	//Convert off-chain information into a format that can be stored on the chain
	vSet := make([][]contract.VerificationG1Point, numVoters)
	cSet := make([][]contract.VerificationG1Point, numVoters)
	a1Set := make([][]contract.VerificationG1Point, numVoters)
	a2Set := make([][]contract.VerificationG1Point, numVoters)
	challengeSet := make([][]*big.Int, numVoters)
	var PVSSComCost float64 = 0
	zSet := make([][]*big.Int, numVoters)
	for j := 0; j < numVoters; j++ {
		vSet[j] = make([]contract.VerificationG1Point, numTalliers+numCandidates)
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
			//PVSSShare Communication Cost
			PVSSComCost = PVSSComCost + float64(sizeOfG1Point(vSet[j][i])) + float64(sizeOfG1Point(cSet[j][i])) + float64(sizeOfG1Point(a1Set[j][i])) + float64(sizeOfG1Point(a2Set[j][i])) + float64(sizeOfBigInt(challengeSet[j][i])) + float64(sizeOfBigInt(zSet[j][i]))

		}
		for i := numTalliers; i < numTalliers+numCandidates; i++ {
			vSet[j][i] = Convert.G1ToG1Point(PVSSShares[j].V[i])
		}
	}

	//Test basic data structure
	//G1Point Size
	sizeG1Point := sizeOfG1Point(vSet[0][0])
	fmt.Printf("Size of G1Point struct: %.6f KB\n", float64(sizeG1Point)/1024)
	//big.Int Size
	sizebigInt := sizeOfBigInt(challengeSet[0][0])
	fmt.Printf("Size of big.Int: %.6f KB\n", float64(sizebigInt)/1024)

	//Each voter's pvss communication cost
	fmt.Printf("Size of PVSSShares of each voter: %.6f KB\n", PVSSComCost/1024)

	//Each voter Vj uploads PVSSShares[j].V, PVSSShares[j].C, PVSSShares[j].Proofs
	/*
		Note:
		  1)Only the encrypted shares PVSSShares[j].C are stored on the blockchain;
		  2)Other messages including PVSSShares[j].V and PVSSShares[j].Proofs can be used to verify the
		    correctess of PVSSShares[j].C and are not stored on the blokchain
		  3)After uploading all information, the UploadPVSSShares algorithm will invoke the DVerify algorithm to
		  	verified PVSSShares[j].C.
		  4)Only verified PVSSShares[j].C will be stored.
	*/

	//Algorithm2(each voter Vj): 1 PVSS.DVerify + l ZKRP.Verify
	sumGasPVSS := uint64(0)
	for j := 0; j < numVoters; j++ {
		auth2 := utils.Transact(client, privatekey, big.NewInt(0))
		tx2, _ := Contract.UploadPVSSShares(auth2, vSet[j], cSet[j], a1Set[j], a2Set[j], challengeSet[j], zSet[j])

		receipt2, err := bind.WaitMined(context.Background(), client, tx2)
		if err != nil {
			log.Fatalf("Tx receipt failed: %v", err)
		}
		sumGasPVSS = sumGasPVSS + receipt2.GasUsed
	}
	fmt.Printf("UploadPVSSShares Gas used(aggregate): %d\n", sumGasPVSS)

	//Algorithm2(each voter Vj): 1 PVSS.DVerify + l ZKRP.Verify+AggregatedC+AggregatedU
	sumGasPVSS1 := uint64(0)
	for j := 0; j < numVoters; j++ {
		auth2 := utils.Transact(client, privatekey, big.NewInt(0))
		tx2, _ := Contract.UploadPVSSShares1(auth2, vSet[j], cSet[j], a1Set[j], a2Set[j], challengeSet[j], zSet[j])

		receipt2, err := bind.WaitMined(context.Background(), client, tx2)
		if err != nil {
			log.Fatalf("Tx receipt failed: %v", err)
		}
		sumGasPVSS1 = sumGasPVSS1 + receipt2.GasUsed
	}
	fmt.Printf("UploadPVSSShares Gas used(no aggregate): %d\n", sumGasPVSS1)

	//Get the result of DVerify algorithm
	DVerifyResult, _ := Contract.GetDVerifyResult(&bind.CallOpts{})
	fmt.Printf("The Verification results of DVerify is %v\n", DVerifyResult)

	//Convert off-chain information into a format that can be stored on the chain
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
	var ZKRPComCost float64 = 0

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
		selectedV[j] = make([]contract.VerificationG1Point, numCandidates)

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
			selectedV[j][i] = Convert.G1ToG1Point(PVSSShares[j].V[i+numTalliers])
			ZKRPComCost = ZKRPComCost + float64(sizeOfG1Point(USet[j][i])) + float64(sizeOfG1Point(ESet[j][i])) + float64(sizeOfG1Point(F1Set[j][i])) + float64(sizeOfG1Point(F2Set[j][i])) + float64(sizeOfG1Point(_USet[j][i])) + float64(sizeOfG1Point(_CSet[j][i])) + float64(sizeOfBigInt(RPcSet[j][i])) + float64(sizeOfBigInt(z1Set[j][i])) + float64(sizeOfBigInt(z2Set[j][i])) + float64(sizeOfBigInt(z3Set[j][i]))
		}
	}

	auth111 := utils.Transact(client, privatekey, big.NewInt(0))
	tx111, _ := Contract.TestAggregateUSet(auth111, USet[0])
	receipt111, err := bind.WaitMined(context.Background(), client, tx111)
	if err != nil {
		log.Fatalf("Tx receipt failed: %v", err)
	}
	fmt.Printf("One AggregateU Gas used: %d\n", receipt111.GasUsed)

	//Each voter's pvss communication cost
	fmt.Printf("Size of ZKRPShares of each voter: %.6f KB\n", ZKRPComCost/1024)

	//All communication cost of Algorithm1
	fmt.Printf("Size of Algorithm1 of each voter: %.6f KB\n", (PVSSComCost+ZKRPComCost)/1024)

	//Upload ballot ciphertexts U, the range proof of voting ciphertexts rangProofs onto smart contract.
	/*
		Note:
		  1)Only the ballot ciphertexts U is stored on the blockchain;
		  2)The range proofs only are used to verify the correctess of ballot ciphertext and are not stored on the blokchain
		  3)After uploading all information, the UploadBallotCipher algorithm will invoke the ZKRPVerify algorithm to
		  	verified U.
		  4)Only ballot ciphertexts which are passed the check will be stored.
	*/
	sumGasZKRP := uint64(0)
	for j := 0; j < numVoters; j++ {
		auth3 := utils.Transact(client, privatekey, big.NewInt(0))
		tx3, _ := Contract.UploadBallotCipher(auth3, ESet[j], F1Set[j], F2Set[j], _USet[j], _CSet[j], RPcSet[j], z1Set[j], z2Set[j], z3Set[j], USet[j], selectedV[j], big.NewInt(int64(threshold)))

		receipt3, err := bind.WaitMined(context.Background(), client, tx3)
		if err != nil {
			log.Fatalf("Tx receipt failed: %v", err)
		}
		sumGasZKRP = sumGasZKRP + receipt3.GasUsed
	}
	fmt.Printf("UploadBallotCipher(Aggregate) Gas used: %d\n", sumGasZKRP)

	sumGasZKRP1 := uint64(0)
	for j := 0; j < numVoters; j++ {
		auth3 := utils.Transact(client, privatekey, big.NewInt(0))
		tx3, _ := Contract.UploadBallotCipher1(auth3, ESet[j], F1Set[j], F2Set[j], _USet[j], _CSet[j], RPcSet[j], z1Set[j], z2Set[j], z3Set[j], USet[j], selectedV[j], big.NewInt(int64(threshold)))

		receipt3, err := bind.WaitMined(context.Background(), client, tx3)
		if err != nil {
			log.Fatalf("Tx receipt failed: %v", err)
		}
		sumGasZKRP1 = sumGasZKRP1 + receipt3.GasUsed
	}
	fmt.Printf("UploadBallotCipher(no Aggregate) Gas used: %d\n", sumGasZKRP1)

	ZKPRVerifyResult, _ := Contract.GetZKRPResult(&bind.CallOpts{})
	fmt.Printf("The Verification results of ZKRPVerify is %v\n", ZKPRVerifyResult[0])

	fmt.Printf("Algorithm2 Gas used: %d\n", sumGasPVSS+sumGasZKRP)

	//Aggregate encrypted shares PVSSShares[j].C
	//Aggrate smart contract executes this operation
	// auth5 := utils.Transact(client, privatekey, big.NewInt(0))
	// tx5, _ := Contract.AggregateC(auth5)

	// receipt5, err := bind.WaitMined(context.Background(), client, tx5)
	// if err != nil {
	// 	log.Fatalf("Tx receipt failed: %v", err)
	// }
	// fmt.Printf("AggregateC Gas used: %d\n", receipt5.GasUsed)

	// auth55 := utils.Transact(client, privatekey, big.NewInt(0))
	// tx55, _ := Contract.AggregateCiphertext(auth55)

	// receipt55, err := bind.WaitMined(context.Background(), client, tx55)
	// if err != nil {
	// 	log.Fatalf("Tx receipt failed: %v", err)
	// }
	// fmt.Printf("AggregateU Gas used: %d\n", receipt55.GasUsed)

	//Off-chain algorithm2(each voter Vj): 1 PVSS.DVerify + l ZKRP.Verify
	Algorithm2Result := make([]bool, numVoters)
	for i := 0; i < numVoters; i++ {
		Algorithm2Result[i] = true
	}
	selectedIndices := make([]*big.Int, threshold)
	for i := 0; i < threshold; i++ {
		selectedIndices[i] = big.NewInt(int64(i + 1))
	}
	x := make([]*big.Int, numCandidates)
	for d := 0; d < numCandidates; d++ {
		x[d] = new(big.Int).Neg(big.NewInt(int64(d)))
		x[d].Mod(x[d], bn256.Order)
	}
	starttime = time.Now().UnixMicro()
	for k := 0; k < int(n); k++ {
		for j := 0; j < numVoters; j++ {
			selectedShares := PVSSShares[j].V[:threshold]
			if !PVSS.DVerify(PVSSShares[j], PP.H0, PKs, numTalliers, numCandidates) {
				Algorithm2Result[j] = false
			} else {
				for d := 0; d < numCandidates; d++ {
					if !ZKRP.Verify(PP.G0, PP.H0, PP.G1, PP.PKI, rangeProofs[j][d], U[j][d], x[d], selectedShares, selectedIndices, threshold) {
						Algorithm2Result[j] = false
						break
					}
				}
			}
		}
	}
	endtime = time.Now().UnixMicro()
	fmt.Printf("Algorithm2 Time Used is %v us\n", (endtime-starttime)/n)
	fmt.Printf("Algorithm2 Result is %v\n", Algorithm2Result[0])

	fmt.Printf("=================================Finish the Voting phase==========================================\n")
	//Each Tallier decrypts the aggreated shares
	//Get the aggregated result from the blockchain
	AggregateResultC, _ := Contract.GetAggregateValue(&bind.CallOpts{})
	//fmt.Printf("The aggregated C is %v\n", AggregateResultC)
	//Each tallier decrypts corresponding aggregated encrypted shares and generates a DLEQ proof.
	sh := make([]*bn256.G1, numTalliers)
	shProof := make([]PVSS.Proof, numTalliers)
	for i := 0; i < numTalliers; i++ {
		//sh[i], shProof[i] = PVSS.Decrypt(PP.G0, PKs[i], aggregatedValue.C[i], SKs[i])
		sh[i], shProof[i] = PVSS.Decrypt(PP.G0, PKs[i], Convert.G1PointToG1(AggregateResultC[i]), SKs[i])
	}

	//TODO:
	//Upload decrypted shares and corresponding DLEQ proofs
	//PVerify smart contracts verifies the correctness of each decrypted shares
	sumGasUsed := uint64(0)
	var DecComCost float64 = 0
	for i := 0; i < numTalliers; i++ {
		//PVSS.Decrypt cpmmunication cost
		DecComCost = DecComCost + float64(sizeOfG1Point(Convert.G1ToG1Point(sh[i]))) + float64(sizeOfG1Point(Convert.G1ToG1Point(shProof[i].RG))) + float64(sizeOfG1Point(Convert.G1ToG1Point(shProof[i].RH))) + float64(sizeOfBigInt(shProof[i].C)) + float64(sizeOfBigInt(shProof[i].Z))
		auth6 := utils.Transact(client, privatekey, big.NewInt(0))
		tx6, _ := Contract.PVerify(auth6, big.NewInt(int64(i)), Convert.G1ToG1Point(sh[i]), Convert.G1ToG1Point(shProof[i].RG), Convert.G1ToG1Point(shProof[i].RH), shProof[i].C, shProof[i].Z)

		receipt6, err := bind.WaitMined(context.Background(), client, tx6)
		if err != nil {
			log.Fatalf("Tx receipt failed: %v", err)
		}
		PVerifyResult, _ := Contract.GetPVerifyResult(&bind.CallOpts{})
		// if PVerifyResult[i] {
		// 	fmt.Printf("The %vth tallier is honest!\n", i)
		// } else {
		// 	fmt.Printf("The %vth tallier is not honest!\n", i)
		// }
		if i == numTalliers-1 {
			fmt.Printf("The Verification results of PVerifyResult is %v\n", PVerifyResult)
		}
		sumGasUsed = sumGasUsed + receipt6.GasUsed
	}
	fmt.Printf("Size of PVSS.Decrypt of all talliers: %.6f KB\n", DecComCost/1024)
	fmt.Printf("PVerify Gas used: %d\n", sumGasUsed)

	//Count the final voting value via Tally algorithm
	//TODO:
	auth7 := utils.Transact(client, privatekey, big.NewInt(0))
	tx7, _ := Contract.Tally(auth7, big.NewInt(int64(threshold)), big.NewInt(int64(numCandidates)), a, b)

	receipt7, err := bind.WaitMined(context.Background(), client, tx7)
	if err != nil {
		log.Fatalf("Tx receipt failed: %v", err)
	}
	fmt.Printf("Tally Gas used: %d\n", receipt7.GasUsed)
	TallyResult, _ := Contract.GetTallyValue(&bind.CallOpts{})
	//fmt.Printf("The tally results are %v\n", TallyResult)
	fmt.Printf("=================================Finish the Tallying phase==========================================\n")

	//Get final tally value
	//finalBallot := make([]int, numCandidates)
	for d := 0; d < numCandidates; d++ {
		for k := (d + 1) * int(aInt); k <= int(bInt)*numVoters; k++ {
			if bytes.Equal(Convert.G1PointToG1(TallyResult[d]).Marshal(), new(bn256.G1).ScalarMult(PP.H0, big.NewInt(int64(k))).Marshal()) {
				fmt.Printf("The %dth voting value is %v\n", d, k)
				break
			}
		}
	}
}
