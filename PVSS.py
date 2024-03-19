from web3 import Web3

w3 = Web3(Web3.HTTPProvider('http://127.0.0.1:7545'))
from solcx import compile_standard, install_solc

install_solc("0.8.0")
import json  # to save the output in a JSON file
import random
import secrets
import sympy  # consider removing this dependency, only needed for mod_inverse
import re
import numpy as np
import hashlib
import datetime
import sys
import time
import ZKRP
from py_ecc.bn128 import G1, G2
from py_ecc.bn128 import add, multiply, neg, pairing, is_on_curve
from py_ecc.bn128 import curve_order as CURVE_ORDER
from py_ecc.bn128 import field_modulus as FIELD_MODULUS
from typing import Tuple, Dict, List, Iterable, Union
from py_ecc.fields import bn128_FQ as FQ

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

global pk, sk
sk = []
pk = []

keccak_256 = Web3.solidityKeccak
"""
H0=(9727523064272218541460723335320998459488975639302513747055235660443850046724,5031696974169251245229961296941447383441169981934237515842977230762345915487)
H1=(5031696974169251245229961296941447383441169981934237515842977230762345915487,9727523064272218541460723335320998459488975639302513747055235660443850046724)
"""
H1 = multiply(G1, 9868996996480530350723936346388037348513707152826932716320380442065450531909)


def random_scalar() -> int:  # Generate random numbers
    """ Returns a random exponent for the BN128 curve, i.e. a random element from Zq.
    """
    return secrets.randbelow(CURVE_ORDER)


def IntsTransform(x):  # tuple/list transform to int[]
    ints = [int(num) for num in x]
    return ints


def Setup(n, t):  # PVSS Key Generation
    """ Generates a random keypair on the BN128 curve.
        The public key is an element of the group G1.
        This key is used for deriving the encryption keys used to secure the shares.
        This is NOT a BLS key pair used for signing messages.
    """
    global sk, pk
    sk.extend([random_scalar() for i in range(0, n)])
    pk.extend([multiply(G1, sk[i]) for i in range(0, n)])  # only need 1-n
    return {"pk": pk, "sk": sk}


def share_secret(secret: int, sharenum: int, threshold: int) -> Dict[
    int, int]:  # PVSS.Share(s,n,t) n is sharenum value .t is threshold value
    coefficients = [secret] + [random_scalar() for j in range(threshold)]

    def f(x: int) -> int:
        """ evaluation function for secret polynomial
        """
        return (
                sum(coef * pow(x, j, CURVE_ORDER) for j, coef in enumerate(coefficients)) % CURVE_ORDER
        )

    indices = [i for i in range(1, sharenum + 1)]
    shares = {x: f(x) for x in indices}

    # coefficients = [secret] + [random_scalar() for j in range(threshold-1)]
    # #The polynomial coefficients
    # def f(x: int) -> int:
    #     """ evaluation function for secret polynomial
    #     """
    #     return (
    #         sum(coef * pow(x, j, CURVE_ORDER) for j, coef in enumerate(coefficients)) % CURVE_ORDER
    #     )
    # shares = { x:f(x) for x in range(1,sharenum+1)}
    # print(indices,shares,[multiply(H1, shares[i]) for i in indices])
    return shares


def Dateconvert(res, n):  # Data conversion functions for bilinear pairing on-chain
    c1 = []
    c2 = []

    v1 = []
    v2 = []

    c1.extend(int(res["c"][i][0]) for i in range(0, n))
    c2.extend(int(res["c"][i][1]) for i in range(0, n))
    v1.extend(int(res["v"][i][0]) for i in range(0, n))
    v2.extend(int(res["v"][i][1]) for i in range(0, n))

    return {"c1": c1, "c2": c2, "v1": v1, "v2": v2}  # c1 is x of c, c2 is y of c. And v1,v2,s1,s2 so on...


def DLEQ(x1, y1, x2, y2, alpha: int) -> Tuple[int, int]:
    """ DLEQ... discrete logarithm equality
        Proofs that the caller knows alpha such that y1 = x1**alpha and y2 = x2**alpha
        without revealing alpha.
    """
    w = random_scalar()
    a1 = multiply(x1, w)
    a2 = multiply(x2, w)
    c = keccak_256(
        abi_types=["uint256"] * 12,
        values=[
            int(v)
            for v in (a1)
                     + (a2)
                     + (x1)
                     + (y1)
                     + (x2)
                     + (y2)
        ],
    )
    c = int.from_bytes(c, "big")
    r = (w - alpha * c) % CURVE_ORDER
    return c, r


def DLEQ_verify(x1, y1, x2, y2, challenge: int, response: int) -> bool:
    a1 = add(multiply(x1, response), multiply(y1, challenge))
    a2 = add(multiply(x2, response), multiply(y2, challenge))
    c = keccak_256(  # pylint: disable=E1120
        abi_types=["uint256"] * 12,  # 12,
        values=[
            int(v)
            for v in (a1)
                     + (a2)
                     + (x1)
                     + (y1)
                     + (x2)
                     + (y2)
        ],
    )
    c = int.from_bytes(c, "big")
    return c == challenge


def DVerify(h0, V_j, pk, C_j, sh_j, n, t):
    for i in range(0, n):
        proof = DLEQ(h0, V_j[i], pk[i], C_j[i], sh_j[i + 1])
        result = Contract.functions.DLEQ_verify(IntsTransform(h0), IntsTransform(V_j[i]), IntsTransform(pk[i]),
                                                IntsTransform(C_j[i]), IntsTransform(proof)).call()
        # gas=Contract.functions.DLEQ_verify(IntsTransform(h0),IntsTransform(V_j[i]),IntsTransform(pk[i]),IntsTransform(C_j[i]),IntsTransform(proof)).estimateGas()
        # print("init ","  ...",gas)
        if (result != True):
            return 0
    return 1


def Decrypt(c_ji, sk_i):
    sh1_ji = multiply(c_ji, sympy.mod_inverse((sk_i) % CURVE_ORDER, CURVE_ORDER))
    return sh1_ji


def DecryptShares(c, sk, n, t):
    sh_j = []
    for i in range(0, n):
        temp = Decrypt(c[i], sk[i])
        sh_j.extend([temp])
    return sh_j


def PVerify(g0, pk, sh1_j, C_j, sk, n, t):
    for i in range(0, n):
        proof = DLEQ(g0, pk[i], sh1_j[i], C_j[i], sk[i])
        result = Contract.functions.DLEQ_verify(IntsTransform(g0), IntsTransform(pk[i]), IntsTransform(sh1_j[i]),
                                                IntsTransform(C_j[i]), IntsTransform(proof)).call()
        if (result != True):
            return 0
    return 1


def Share(s_j, H1, pk, n, t):
    SSShare = share_secret(s_j, n, t)  # voter PVSS.share=(v,c)
    v = [0]
    c = [0]

    v.extend([multiply(H1, SSShare[i + 1]) for i in range(0, n)])  # v_i=g2^s_i
    c.extend([multiply(pk[i], SSShare[i + 1]) for i in range(0, n)])  # c_i=pk_i^s_i
    res = {"v": v, "c": c, "raw": SSShare, "s": s_j}
    # print(SSShare)
    return res


def Reconstruct(res, n, t):
    # print(res)
    sum = multiply(H1, 0)

    def lagrange_coefficient(i: int) -> int:
        result = 1
        for j in res["raw"]:
            # j=j-1
            if i != j:
                result *= j * sympy.mod_inverse((j - i) % CURVE_ORDER, CURVE_ORDER)
                result %= CURVE_ORDER
        return result

    # print(res["raw"].items())
    for i, _ in res["raw"].items():
        # print(i)
        sum = add(sum, multiply(res["v"][i], lagrange_coefficient(i)))
        # print(lagrange_coefficient(i))

    return sum


def Reconstruct2(share, n, t):
    res_list = [[int(x) for x in tup] for tup in share["v"]]
    # print(res_list)
    # print(share["v"])
    result = Contract.functions.Interpolate(res_list, LagrangeCoefficient(share["raw"]), int(H1[0]), int(H1[1]),
                                            t).call()
    return result


def Share2(s_j, H1, pk, n, t):
    PVSSshare = share_secret(s_j, n, t)  # voter PVSS.share=(v,c)
    v = []
    c = []
    DLEQ_Proof = []
    """
    print(PVSSshare)
    for i in range(0,n):
        print(PVSSshare[i+1])
        tmp=multiply(G1,PVSSshare[i+1])
        v.extend(tmp)
    """
    v.extend([multiply(H1, PVSSshare[i + 1]) for i in range(0, n)])  # v_i=g2^s_i
    c.extend([multiply(pk[i], PVSSshare[i + 1]) for i in range(0, n)])  # c_i=pk_i^s_i
    DLEQ_Proof.extend([DLEQ(H1, v[i], pk[i], c[i], PVSSshare[i + 1]) for i in range(0, n)])
    res = {"v": v, "c": c, "DLEQ_Proof": DLEQ_Proof}
    return res


if __name__ == '__main__':
    n = 10
    t = 5
    s = 233333
    key = Setup(n, t)

    print("answer")
    print(multiply(H1, s))
    print(".........")
    share1 = Share(s, H1, pk, n, t)
    # print(share1)
    sum = Reconstruct(share1, n, t)
    print(sum)

    # res=Reconstruct2(share1,2,1)
    # print(res)

"""
print(key)
print(pk)
print(sk)
print(G1)
print(type(G1))
print(H0)
print(is_on_curve(H0,FQ(3)))    
shares=Share(random_scalar(),H0,pk)
print(shares)


y1=multiply(G1,shares["s"])
y2=multiply(H0,shares["s"])

proof1=DLEQ(H0,shares["v"][0],pk[0],shares["c"][0],shares["raw"][1])
print(DLEQ_verify(H0,shares["v"][0],pk[0],shares["c"][0],proof1[0],proof1[1]))

result=Contract.functions.DLEQ_verify(IntsTransform(H0),IntsTransform(shares["v"][0]),IntsTransform(pk[0]),IntsTransform(shares["c"][0]),IntsTransform(proof1)).call()
if(result==True):
    print("ok")
"""

"""
starttime = time.time()
SS=DVerify(H0,shares["v"],pk,shares["c"],shares["raw"])
print("SCRAPE DDH verification cost ",time.time()- starttime)    
print(SS)


sh1_j=DecryptShares(shares["c"],sk)
SS=PVerify(G1,pk,sh1_j,shares["c"],sk)
print(SS)

"""

