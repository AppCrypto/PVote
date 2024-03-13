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
import PVSS
import VoteforDAOs
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

sk_I = 15872232885142738667420701097223674108720232256552480080547895231827275416057
pk_I = multiply(G2, sk_I)
H1 = multiply(G1, 9868996996480530350723936346388037348513707152826932716320380442065450531909)


def Vote(w_j: int):  # vote value
    s_j = PVSS.random_scalar()
    shares = PVSS.Share(s_j, H1, pk)

    U_j = add(multiply(H1, w_j), multiply(G1, s_j))
    proof = ZKRP.Prove(s_j, w_j, U_j, GPK["sigam_k"][(w_j - 1) % 5])

    return {"v": shares["v"], "c": shares["c"], "U": U_j, "Proof": proof, "raw": shares["raw"]}


def VoteVerify(shares):
    SS = PVSS.DVerify(H1, shares["v"], pk, shares["c"], shares["raw"])

    result = ZKRP.ZKRP_verify(shares["Proof"], shares["v"], PVSS.LagrangeCoefficient(shares["raw"]), shares["U"])
    print("PVSS DVerify Result:" + str(SS))
    print("ZKRP Verify Result:" + str(result))


def TallierVerify(shares):
    SS = PVSS.PVerify(G1, pk, shares, ReturnPointC(), sk)
    print("PVSS PVerify Result:" + str(SS))


def Aggreagate(shares):
    agg = PVSS.Dateconvert(shares)
    Contract.functions.Aggregate(agg["c1"], agg["c2"], agg["v1"], agg["v2"], int(shares["U"][0]),
                                 int(shares["U"][1])).transact({'from': w3.eth.accounts[0]})
    print("Aggregate done.")
    return 1


def Tally(sh1_j, lar):
    res_list = [[int(x) for x in tup] for tup in sh1_j]
    result = Contract.functions.Tally(res_list, lar).call()
    AllResult = []
    for i in range(a * m, b * m + 1):
        AllResult.extend([multiply(H1, i)])

    for i in range(0, (b - a) * m + 1):
        if (result[0] == AllResult[i][0] and result[1] == AllResult[i][1]):
            print("The vote score is " + str(i + a * m))
    return result


def ReturnDate():
    agg = Contract.functions.ReturnData().call({'from': w3.eth.accounts[0]})
    print("Agg C:" + str(agg[0]))
    print("Agg V:" + str(agg[1]))
    print("Agg U:" + str(agg[2]))


def ReturnPointC():
    agg = Contract.functions.ReturnPointC().call({'from': w3.eth.accounts[0]})
    formatted_data = [(FQ(item[0]), FQ(item[1])) for item in agg]
    return formatted_data


def Init(n: int):
    Contract.functions.Init(n).transact({'from': w3.eth.accounts[0]})
    # print("............................Init Done ...............................")


def VoteAgg(x):  # Multiple voters are selected to vote at once
    for i in range(0, x):
        Vote1 = Vote(4)
        VoteVerify(Vote1)
        Aggreagate(Vote1)


if __name__ == '__main__':
    print("...........................................Setup phase.............................................")
    Init(PVSS.n)
    key = PVSS.Setup()  # PVSS Key Generation
    # pvss setup
    pk = key["pk"]  # Public key array
    sk = key["sk"]  # Private key array

    a = 1
    b = 5
    m = 10
    GPK = ZKRP.Setup(a, b)

    print("............................................Voting phase...........................................")
    VoteAgg(10)
    Vote11 = Vote(4)
    VoteVerify(Vote11)
    Aggreagate(Vote11)

    print("..........................................tallying phase...........................................")
    sh1 = PVSS.DecryptShares(ReturnPointC(), sk)

    TallierVerify(sh1)

    Tally(sh1, PVSS.LagrangeCoefficient(Vote11["raw"]))

    print("............................................Reward phase...........................................")

