from web3 import Web3
import sys

w3 = Web3(Web3.HTTPProvider('http://127.0.0.1:7545'))
from solcx import compile_standard, install_solc

install_solc("0.8.0")
import json  # to save the output in a JSON file
import sympy  # consider removing this dependency, only needed for mod_inverse
import time
import ZKRP
import PVSS
import util
import random
from py_ecc.bn128 import G1, G2
from py_ecc.bn128 import add, multiply
from py_ecc.bn128 import curve_order as CURVE_ORDER
from py_ecc.bn128 import field_modulus as FIELD_MODULUS
from py_ecc.fields import bn128_FQ as FQ

with open("contracts/CandidatesVote.sol", "r") as file:
    contact_list_file = file.read()

compiled_sol = compile_standard(
    {
        "language": "Solidity",
        "sources": {"CandidatesVote.sol": {"content": contact_list_file}},
        "settings": {
            "outputSelection": {
                "*": {
                    "*": ["abi", "metadata", "evm.bytecode", "evm.bytecode.sourceMap"]
                    # output needed to interact with and deploy contract
                }
            }
        },
    },
    solc_version="0.8.0",
)

# print(compiled_sol)
with open("compiled_code.json", "w") as file:
    json.dump(compiled_sol, file)
# get bytecode
bytecode = compiled_sol["contracts"]["CandidatesVote.sol"]["CandidatesVote"]["evm"]["bytecode"]["object"]
# get abi
abi = json.loads(compiled_sol["contracts"]["CandidatesVote.sol"]["CandidatesVote"]["metadata"])["output"]["abi"]
# Create the contract in Python
contract = w3.eth.contract(abi=abi, bytecode=bytecode)

chain_id = 5777
accounts0 = w3.eth.accounts[0]
transaction_hash = contract.constructor().transact({'from': accounts0})
# Wait for the contract to be deployed
transaction_receipt = w3.eth.wait_for_transaction_receipt(transaction_hash)
# Get the deployed contract address
contract_address = transaction_receipt['contractAddress']

Contract = w3.eth.contract(address=contract_address, abi=abi)

sk_I = 15872232885142738667420701097223674108720232256552480080547895231827275416057
# Trusted third party's private key, where here is the vote initiator's private key

pk_I = multiply(G2, sk_I)
# Trusted third party's public key, here is the vote initiator's public key

H1 = multiply(G1, 9868996996480530350723936346388037348513707152826932716320380442065450531909)  # Generator H1

pks = []  # store for talliers public key


def Vj_VoteForCandidates(w_j: tuple, n: int, t: int):
    """
        input w_j is vote value  range in (a,b)
        The function defines the transaction that a voter V_j should complete,
        (n, t) is the number of participants and the threshold value of secret sharing

        In the voting function,
        n denotes the number of talliers,
        t denotes the number of talliers that need to be satisfied for the calculation of the voting result
    """
    s_j = PVSS.random_scalar()
    # Generate random numbers

    shares = PVSS.ShareForCandidates(s_j, H1, pk, n, t)
    # Call PVSS.Share

    U_j = []
    U_j.extend([add(multiply(H1, w_j[0]), multiply(G1, s_j))])
    for i in range(1, len(w_j)):
        U_j.extend([add(multiply(H1, w_j[i]), multiply(G1, shares["P_j"][i]))])
    # Compute U_jn  for every Candidater i

    zkrp_proof = []
    zkrp_proof.extend([ZKRP.Prove(s_j, w_j[0], U_j[0], GPK["sigam_k"][(w_j[0]) % (b + 1)])])
    for i in range(1, len(w_j)):
        zkrp_proof.extend([ZKRP.Prove(shares["P_j"][i], w_j[i], U_j[i], GPK["sigam_k"][(w_j[i]) % (b + 1)])])
    # Call ZKRP.Prove to generate proof for the generated data for every U_ji

    dleq_proof = []
    for i in range(0, n):
        temp = util.Point2IntArr(shares["DLEQ_Proof"][i])
        dleq_proof.extend([temp])
    # Convert for DLEQ Proof data format

    agg = util.Dataconvert(shares)
    ugg = util.U_jdataconvert(U_j)
    # Data transformation

    Contract.functions.PVSStoSC(agg["c1"], agg["c2"], agg["v1"], agg["v2"], ugg["U_j1"], ugg["U_j2"],
                                dleq_proof, len(w_j)).transact({'from': w3.eth.accounts[0]})
    # Transfer v, c array, dleq_proof array, U_j of PVSS.Share generated by the voter to the smart contract

    for i in range(0, len(w_j)):

        # Transfer Proof generated by ZKRP.Prove generated by the voter to the smart contract for every ZKRP Proof
        Contract.functions.ZKRPtoSC(util.Point2IntArr(zkrp_proof[i][0]), util.Point2IntArr(zkrp_proof[i][1]),
                                    util.Point2IntArr(zkrp_proof[i][2]), util.Point2IntArr(zkrp_proof[i][3]),
                                    util.Point2IntArr(zkrp_proof[i][4]), zkrp_proof[i][5], zkrp_proof[i][6],
                                    zkrp_proof[i][7],
                                    zkrp_proof[i][8], int(U_j[i][0]), int(U_j[i][1])).transact(
            {'from': w3.eth.accounts[0]})

        if (Contract.functions.ZKRP_verify(i, n).call() != True):  # ZKRP.Verify
            print("ZKRP_Verify failue", i, Contract.functions.ZKRP_verify(i, n).call())
            return -1

    print("Vote value:", w_j)
    return True


def Ti_Tally(No: int, pk_i, sk_i):
    """
        The function defines the transactions that a tallier T_i should complete
        No denotes the number of tallier, for example: No 1 is the tallier1 and (pk_1, sk_1) is the key of tallier1
    """

    aggCV = Contract.functions.DownloadAGGVC(No).call()
    # Download the accumulated V,C data from the chain

    C_i = (FQ(aggCV[0][0]), FQ(aggCV[0][1]))
    # Data transformation

    sh1 = PVSS.Decrypt(C_i, sk_i)
    # Call PVSS.Decrypt function to decrypt the cumulative share C*

    proof = PVSS.DLEQ(G1, pk_i, sh1, C_i, sk_i)
    # Generate DLEQ P_Proof, proof is the share c decrypted by this tallier T_i

    Contract.functions.Decrypted_SharetoSC(No, util.Point2IntArr(sh1), util.Point2IntArr(proof)).transact(
        {'from': w3.eth.accounts[0]})
    # upload the decrypted share and P_Proof to the chain.
    # Once verified, keep the decrypted share in the DecryptedShare array on the chain
    print("Tallier", No, "done")


def Aggreagate():
    """
        Execute the on-chain Aggregate function once to aggregate the V and C data onto the chain.
        No arguments are needed since the v and c data are already on the chain
    """
    Contract.functions.Aggregate().transact({'from': w3.eth.accounts[0]})
    print("Aggregate done.")


def Tally(l, m):
    """
        For on-chain voting, the parameter m (number of votes) is entered,
        because we want to determine the range of votes (a*m,b*m).
    """
    score = []
    AllResult = {}
    for i in range(a * m, b * m + 1):
        AllResult[i] = multiply(H1, i)
    # Calculate the likelihood of all the vote results

    for i in range(0, l):
        result = Contract.functions.Tally(i).call()
        # Get the candidates vote result

        added = False
        for i in range(a * m, b * m + 1):
            if (AllResult[i] != None and result[0] == AllResult[i][0] and result[1] == AllResult[i][1]):
                # print("The vote score is",i)
                score.append(i)
                added = True
                break
        if (not add):
            score.append(-1)
            # Compare the result of the vote with the traversal of possible vote results
    return score


if __name__ == '__main__':

    n = int(sys.argv[1])
    # Set the number n of talliers

    t = int(n / 2)
    # Set the threshold value t
    print("...........................................Setup phase.............................................", n, t)
    key = PVSS.Setup(n, t)
    # PVSS Key Generation

    pk = key["pk"]  # Set public key array
    sk = key["sk"]  # Set private key array
    pks = [util.Point2IntArr(pk[i]) for i in range(n)]  # Data transformation

    Contract.functions.setTalliersPK(pks).transact({'from': w3.eth.accounts[0]})
    # Upload the public key to the smart contract

    a = 0  # Vote minimum range a
    b = 5  # Vote Max  range b
    m = 5  # Set the number of voters m
    GPK = ZKRP.Setup(a, b)
    # ZKRP initialization

    l = int(n / 2) + 1  # Number of candidates for this vote
    # l=int(sys.argv[2])
    print("............................................Voting phase...........................................")

    ballot = [0] * l
    for i in range(0, m):  # Generate m voters
        w_j = [int(random.random() * (b - a + 1) + a) for i in range(l)]
        # Randomly generate voting values w_j ∈ (a,b)

        ballot = [b + w for b, w in zip(ballot, w_j)]
        proof = Vj_VoteForCandidates(w_j, n, t)
        # Voter voting function

        if (Contract.functions.PVSS_DVerify().call() and proof):
            Aggreagate()
            # Aggregate the vote uploads after passing both PVSS_DVierfy and ZKRP_Verify verifications
        else:
            print("Invalid vote value w_j......did not aggreagate in SC")

    print("expected ballot value:", ballot)

    print("..........................................tallying phase...........................................")

    # When a player completes the Ti_Tally function,
    # it will successfully upload a decryption share on the chain.
    # When t shares are met, the Tally function can be called to perform the vote
    temp_t = t + 1
    # temp_t = int(sys.argv[2])

    for i in range(0, temp_t):  # Set the number of tallier to be completed
        Ti_Tally(i + 1, pk[i], sk[i])

    print("The tallying result is :", Tally(l, m))  # On-chain tally
    print("............................................Reward phase...........................................")