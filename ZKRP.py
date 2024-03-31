from web3 import Web3
w3 = Web3(Web3.HTTPProvider('http://127.0.0.1:7545'))
from solcx import compile_standard, install_solc

install_solc("0.8.0")
import json  # to save the output in a JSON file
import PVSS
import util
import time
import re
import sympy  # Needed for mod_inverse
from web3 import Web3
from py_ecc.bn128 import G1, G2
from py_ecc.bn128 import add, multiply, neg, pairing
from py_ecc.bn128 import curve_order as CURVE_ORDER
from py_ecc.bn128 import field_modulus as FIELD_MODULUS

H1 = multiply(G1, 9868996996480530350723936346388037348513707152826932716320380442065450531909)  #Generator H1
keccak_256 = Web3.solidityKeccak

with open("contracts/DAOsForVote.sol", "r") as file:
    contact_list_file = file.read()

compiled_sol = compile_standard(
    {
        "language": "Solidity",
        "sources": {"DAOsForVote.sol": {"content": contact_list_file}},
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

with open("compiled_code.json", "w") as file:
    json.dump(compiled_sol, file)
# get bytecode
bytecode = compiled_sol["contracts"]["DAOsForVote.sol"]["DAOsForVote"]["evm"]["bytecode"]["object"]
# get abi
abi = json.loads(compiled_sol["contracts"]["DAOsForVote.sol"]["DAOsForVote"]["metadata"])["output"]["abi"]
# Create the contract in Python
contract = w3.eth.contract(abi=abi, bytecode=bytecode)

chain_id = 5777
accounts0 = w3.eth.accounts[0]
transaction_hash = contract.constructor().transact({'from': accounts0})
# Wait for the contract to be deployed
transaction_receipt = w3.eth.wait_for_transaction_receipt(transaction_hash)
# Get the deployed contract address
contract_address = transaction_receipt['contractAddress']
# print(" contract deployed, address: ", contract_address)
Contract = w3.eth.contract(address=contract_address, abi=abi)

sk_I = 15872232885142738667420701097223674108720232256552480080547895231827275416057
# Trusted third party's private key, where here is the vote initiator's private key

pk_I = multiply(G2, sk_I)
# Trusted third party's public key, here is the vote initiator's public key

def Setup(a, b):
    # initialization of ZKRP, a-b for the selected range, generate sigam_k for the integers in the range
    sigma_k = []
    for i in range(a, b + 1):
        temp = multiply(G1, sympy.mod_inverse((sk_I + i) % CURVE_ORDER, CURVE_ORDER))
        sigma_k.extend([temp])
    return {"sk_I": sk_I, "pk_I": pk_I, "sigam_k": sigma_k}


def Prove(s_j, w_j, U_j, sigma_wj):  #ZKRP.Prove 生成Proof
    v = PVSS.random_scalar()
    s = PVSS.random_scalar()
    t = PVSS.random_scalar()
    m = PVSS.random_scalar()

    E_j = multiply(sigma_wj, v)
    F_j1 = multiply(E_j, s)
    """
    # where the generation of bilinear pairing off the chain is very long, 
    F_j is split into two groups of numbers, F_j1 and F_j2 respectively, 
    and uploaded to the chain for a bilinear pairing.

    Bilinear pairing takes at least 3~5s time. 
    After modifying to group element, the whole Prove generation only takes 0.00085s on average
    """
    F_j2 = neg(multiply(G1, t))
    U1_j = add(multiply(H1, s), multiply(G1, m))
    C1_j = multiply(H1, m)

    c = keccak_256(
        abi_types=["uint256"] * 12,
        values=[
            int(v)
            for v in (E_j)
                     + (U_j)
                     + (F_j1)
                     + (F_j2)
                     + (U1_j)
                     + (C1_j)
        ],
    )
    c = int.from_bytes(c, "big")
    z1 = (s - w_j * c) % CURVE_ORDER
    z2 = (t - v * c) % CURVE_ORDER
    z3 = (m - s_j * c) % CURVE_ORDER
    return E_j, F_j1, F_j2, U1_j, C1_j, c, z1, z2, z3    #Output of one ZKRP.Prove


def Verify(proof, V_j, U_j, s_j, pk_I):  # off-chain  ZKRP.Verify for testing purposes
    C_j = multiply(H1, s_j)

    if (proof[3] != add(multiply(C_j, proof[4]), multiply(H1, proof[7]))):
        print("no1")
        return 0

    if (proof[2] != add(multiply(U_j, proof[4]), add(multiply(G1, proof[7]), multiply(H1, proof[5])))):
        print("no2")
        return 0

    if (proof[8] != (pairing(pk_I, multiply(proof[0], proof[4])) * pairing(G2, multiply(proof[0], CURVE_ORDER - proof[
        5])) * pairing(G2, multiply(G1, proof[6])))):
        print("no3")
        return 0

    return 1

"""
    On-chain validation of ZKRP, saved in the on-chain solidity DAOsForVote.sol  and CandidatesVote.sol file.

    The specific operation is to verify the following three verification expressions.

    function ZKRP_verify(uint8 t)
    {
        bool  proof1 = ZKRP_verify1(C_j,t);
        bool  proof2 = ZKRP_verify2();
        bool  proof3 = ZKRP_verify3();
        return proof1 && proof2 && proof3;
    }

    function ZKRP_verify1(uint256[2] memory C_j,uint8 t)

    function ZKRP_verify2()

    function ZKRP_verify3()
"""