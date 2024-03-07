from web3 import Web3
w3 = Web3(Web3.HTTPProvider('http://127.0.0.1:7545'))
from solcx import compile_standard,install_solc
install_solc("0.8.0")
import json  #to save the output in a JSON file
import random
import secrets
import sympy # consider removing this dependency, only needed for mod_inverse
import re
import numpy as np
import hashlib
import datetime
import sys
import time
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
                     "*": ["abi", "metadata", "evm.bytecode", "evm.bytecode.sourceMap"] # output needed to interact with and deploy contract 
                }
            }
        },
    },
    solc_version="0.8.0",
)

#print(compiled_sol)
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
#print("合约已部署，地址：", contract_address)
Contract = w3.eth.contract(address=contract_address, abi=abi)

global pk,sk,n,t
sk=[]
pk=[] 
n=10   #PVSS Total distribution 
t=5    #PVSS threshold value

keccak_256 = Web3.solidityKeccak
"""
H0=(9727523064272218541460723335320998459488975639302513747055235660443850046724,5031696974169251245229961296941447383441169981934237515842977230762345915487)
H1=(5031696974169251245229961296941447383441169981934237515842977230762345915487,9727523064272218541460723335320998459488975639302513747055235660443850046724)
"""
H0=multiply(G1,2)



def random_scalar() -> int: #Generate random numbers
    """ Returns a random exponent for the BN128 curve, i.e. a random element from Zq.
    """
    return secrets.randbelow(CURVE_ORDER)

def IntsTransform(x):
    ints=[int(num) for num in x]
    return ints

def Setup() : #PVSS Key Generation
    """ Generates a random keypair on the BN128 curve.
        The public key is an element of the group G1.
        This key is used for deriving the encryption keys used to secure the shares.
        This is NOT a BLS key pair used for signing messages.
    """
    global sk,pk
    sk.extend([random_scalar() for i in range(0,n)])  
    pk.extend([multiply(G1, sk[i]) for i in range(0,n)]) #only need 1-n
    return {"pk":pk,"sk":sk}

def share_secret(secret:int ,sharenum:int ,threshold: int)-> Dict[int, int]: # PVSS.Share(s,n,t) n is sharenum value .t is threshold value

    coefficients = [secret] + [random_scalar() for j in range(threshold-1)]
    #The polynomial coefficients 
    def f(x: int) -> int:
        """ evaluation function for secret polynomial
        """
        return (
            sum(coef * pow(x, j, CURVE_ORDER) for j, coef in enumerate(coefficients)) % CURVE_ORDER
        )
    shares = { x:f(x) for x in range(1,sharenum+1) }
    #print(shares)
    return shares

def Share(s_j,h0,pk):
    PVSSshare=share_secret(s_j,n,t)  #voter PVSS.share=(v,c) 
    v=[]
    c=[]
    """
    print(PVSSshare)
    for i in range(0,n):
        print(PVSSshare[i+1])
        tmp=multiply(G1,PVSSshare[i+1])
        v.extend(tmp)
    """
    v.extend([multiply(h0,PVSSshare[i+1]) for i in range(0,n)])  #v_i=g2^s_i 
    c.extend([multiply(pk[i],PVSSshare[i+1]) for i in range(0,n)]) #c_i=pk_i^s_i
    res={"v":v,"c":c,"raw":PVSSshare,"s":s_j}
    return res 

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
    return c,r


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


def DVerify(h0,V_j,pk,C_j,sh_j):
    for i in range(0,n):
        proof=DLEQ(h0,V_j[i],pk[i],C_j[i],sh_j[i+1])
        result=Contract.functions.DLEQ_verify(IntsTransform(h0),IntsTransform(V_j[i]),IntsTransform(pk[i]),IntsTransform(C_j[i]),IntsTransform(proof)).call()
        if(result!=True):
            return 0
    return 1


def Decrypt(c_ji,sk_i):
    sh1_ji=multiply(c_ji,sympy.mod_inverse((sk_i) % CURVE_ORDER, CURVE_ORDER))
    return sh1_ji

def DecryptShares(c,sk):
    sh_j=[]
    for i in range(0,n):
        temp=Decrypt(c[i],sk[i])
        sh_j.extend([temp])
    return sh_j

def PVerify(g0,pk,sh1_j,C_j,sk):
    for i in range(0,n):
        proof=DLEQ(g0,pk[i],sh1_j[i],C_j[i],sk[i])
        result=Contract.functions.DLEQ_verify(IntsTransform(g0),IntsTransform(pk[i]),IntsTransform(sh1_j[i]),IntsTransform(C_j[i]),IntsTransform(proof)).call()
        if(result!=True):
            return 0
    return 1

#def Reconstruct():

"""
key=Setup()
shares=Share(random_scalar(),H0,pk)


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

