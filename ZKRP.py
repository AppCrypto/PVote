from web3 import Web3

w3 = Web3(Web3.HTTPProvider('http://127.0.0.1:7545'))
from solcx import compile_standard, install_solc

install_solc("0.8.0")
import json  # to save the output in a JSON file
import PVSS
import time
import re
import sympy  # Needed for mod_inverse
from web3 import Web3
from py_ecc.bn128 import G1, G2
from py_ecc.bn128 import add, multiply, neg, pairing, is_on_curve
from py_ecc.bn128 import curve_order as CURVE_ORDER
from py_ecc.bn128 import field_modulus as FIELD_MODULUS

global pk, sk, n, t
n = 10  # tally_people /Registered Tallier
t = 5  # The voting system need tallier to recover secret
H1 = multiply(G1, 9868996996480530350723936346388037348513707152826932716320380442065450531909)
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

# print(compiled_sol)
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
# 等待合约部署完成
transaction_receipt = w3.eth.wait_for_transaction_receipt(transaction_hash)
# 获取部署后的合约地址
contract_address = transaction_receipt['contractAddress']
# print("合约已部署，地址：", contract_address)
Contract = w3.eth.contract(address=contract_address, abi=abi)

sk_I = 15872232885142738667420701097223674108720232256552480080547895231827275416057
pk_I = multiply(G2, sk_I)


# starttime = time.time()
# SS=PVSS.DVerify(H1,shares["v"],pk,shares["c"],shares["raw"])
# print("SCRAPE DDH verification cost ",time.time()- starttime)
# print(SS)


def Setup(a, b):
    sigma_k = []
    for i in range(a, b + 1):
        temp = multiply(G1, sympy.mod_inverse((sk_I + i) % CURVE_ORDER, CURVE_ORDER))
        sigma_k.extend([temp])
    return {"sk_I": sk_I, "pk_I": pk_I, "sigam_k": sigma_k}


def Prove(s_j, w_j, U_j, sigma_wj):
    v = PVSS.random_scalar()
    s = PVSS.random_scalar()
    t = PVSS.random_scalar()
    m = PVSS.random_scalar()
    # m1 = PVSS.random_scalar()

    E_j = multiply(sigma_wj, v)

    """
    par1=pairing(G2,multiply(E_j,CURVE_ORDER-s))*pairing(G2,multiply(G1,t))

    def F_jTransform(par):
    #F_j = pairing(G2,G1)*pairing
        temp = re.findall("\d+",str(par))
        tempF_j=[]
        for i in range (len(temp)):
            tempF_j.append(int(temp[i]))
        return tuple(tempF_j)

    F_j=F_jTransform(par1)
    """
    # F_j = add(multiply(G1,t),neg(multiply(E_j,s)))
    # F_j = neg(F_j)
    F_j1 = multiply(E_j, s)
    F_j2 = neg(multiply(G1, t))
    U1_j = add(multiply(H1, s), multiply(G1, m))
    C1_j = multiply(H1, m)
    # D_j =  add(multiply(L1,m),add(multiply(E_j,CURVE_ORDER-s),multiply(G1,t)))
    # G_j = multiply(add(L1,neg(H1)),m1)

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
    # z4 = (m1 - m * c) % CURVE_ORDER
    return E_j, F_j1, F_j2, U1_j, C1_j, c, z1, z2, z3


def Verify(proof, V_j, U_j, s_j, pk_I):
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


def ZKRP_verify(proof, V_j, lagrangeCoefficient, U_j):
    def IntsTransform(x):  # tuple/list transform to int[]
        ints = [int(num) for num in x]
        return ints

    vv = []
    for i in range(0, n):
        temp = PVSS.IntsTransform(V_j[i])
        vv.extend([temp])

    result1 = Contract.functions.ZKRP_verify1(vv, lagrangeCoefficient, IntsTransform(proof[4]), proof[5],
                                              proof[8]).call()
    result2 = Contract.functions.ZKRP_verify2(IntsTransform(proof[3]), IntsTransform(U_j), proof[5], proof[6],
                                              proof[8]).call()
    result3 = Contract.functions.ZKRP_verify3(IntsTransform(proof[1]), IntsTransform(proof[2]), IntsTransform(proof[0]),
                                              proof[5], proof[6], proof[7]).call()

    return result1, result2, result3



