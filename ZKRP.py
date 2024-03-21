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
from py_ecc.bn128 import add, multiply, neg, pairing, is_on_curve
from py_ecc.bn128 import curve_order as CURVE_ORDER
from py_ecc.bn128 import field_modulus as FIELD_MODULUS

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

sk_I = 15872232885142738667420701097223674108720232256552480080547895231827275416057  #可信第三方的私钥，其中这里为投票发起者的私钥
pk_I = multiply(G2, sk_I)  #可信第三方公钥，为投票发起者的公钥


def Setup(a, b):   #ZKRP的初始化，a-b为选定的范围，为范围内的整数生成sigam_k
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
    F_j1 = multiply(E_j, s)      #其中链下生成双线性配对十分漫长，将F_j拆分为两个群上的数，分别为F_j1和F_j2,上传到链上进行一次双线性配对
    F_j2 = neg(multiply(G1, t))  #双线性配对起码占用3~5s时间 ，修改为群元素后，整个Prove生成平均只需要0.00085s
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

    return E_j, F_j1, F_j2, U1_j, C1_j, c, z1, z2, z3


def Verify(proof, V_j, U_j, s_j, pk_I):  #ZKRP的链下验证，为测试所用
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


def ZKRP_verify(V_j, n, t):        #ZKRP的链上验证

    recIndex = [i + 1 for i in range(0, t + 1)]  #确定t个份额的下标

    def lagrange_coefficient(i: int) -> int:  #计算拉格朗日多项式系数
        result = 1
        for j in recIndex:
            # print(j)
            # j=j-1
            if i != j:
                result *= j * sympy.mod_inverse((j - i) % CURVE_ORDER, CURVE_ORDER)
                result %= CURVE_ORDER
        return result

    lar = [lagrange_coefficient(i) for i in recIndex]   #转换为int[]
    V = [util.Point2IntArr(V_j[i]) for i in recIndex]  #转换为uint256[2][]
    result1 = Contract.functions.ZKRP_verify1(V, lar).call() #ZKRP.Verify的第一个等式的验证
    result2 = Contract.functions.ZKRP_verify2().call()  #ZKRP.Verify的第二个等式的验证
    result3 = Contract.functions.ZKRP_verify3().call()  #ZKRP.Verify的第三个等式的验证

    if result1 and result2 and result3:  #三个验证等式全true，ZKRP.Verify才会返回true
        return (True)
    else:
        return (False)


def ZKRP_verify2(V_j, n, t):        #ZKRP的链上验证    
    recIndex = [i + 1 for i in range(0, t + 1)]  #确定t个份额的下标

    V = [util.Point2IntArr(V_j[i]) for i in recIndex]  #转换为uint256[2][]
    result1 = Contract.functions.ZKRP_verify1(V).call() #ZKRP.Verify的第一个等式的验证
    result2 = Contract.functions.ZKRP_verify2().call()  #ZKRP.Verify的第二个等式的验证
    result3 = Contract.functions.ZKRP_verify3().call()  #ZKRP.Verify的第三个等式的验证

    if result1 and result2 and result3:  #三个验证等式全true，ZKRP.Verify才会返回true
        return (True)
    else:
        return (False)
