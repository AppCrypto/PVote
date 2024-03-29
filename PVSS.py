from web3 import Web3

w3 = Web3(Web3.HTTPProvider('http://127.0.0.1:7545'))
from solcx import compile_standard, install_solc

install_solc("0.8.0")
import json  # to save the output in a JSON file
import random
import util
import secrets
import sympy  # consider removing this dependency, only needed for mod_inverse
from py_ecc.bn128 import G1, G2
from py_ecc.bn128 import add, multiply, neg, pairing, is_on_curve
from py_ecc.bn128 import curve_order as CURVE_ORDER
from py_ecc.bn128 import field_modulus as FIELD_MODULUS
from typing import Tuple, Dict, List, Iterable, Union

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

keccak_256 = Web3.solidityKeccak
"""
H0=(9727523064272218541460723335320998459488975639302513747055235660443850046724,5031696974169251245229961296941447383441169981934237515842977230762345915487)
H1=(5031696974169251245229961296941447383441169981934237515842977230762345915487,9727523064272218541460723335320998459488975639302513747055235660443850046724)
"""
H1 = multiply(G1, 9868996996480530350723936346388037348513707152826932716320380442065450531909)  # 生成元H1


def random_scalar() -> int:  # Generate random numbers
    """ Returns a random exponent for the BN128 curve, i.e. a random element from Zq.
    """
    return secrets.randbelow(CURVE_ORDER)


def Setup(n, t):  # PVSS Key Generation   #PVSS的公私钥产生，注意有效密钥从0开始索引
    """ Generates a random keypair on the BN128 curve.
        The public key is an element of the group G1.
        This key is used for deriving the encryption keys used to secure the shares.
        This is NOT a BLS key pair used for signing messages.
    """
    sk = []
    pk = []
    sk.extend([random_scalar() for i in range(0, n)])
    pk.extend([multiply(G1, sk[i]) for i in range(0, n)])  # only need 1-n

    return {"pk": pk, "sk": sk}


def share_secret(secret: int, sharenum: int, threshold: int) -> Dict[
    int, int]:  # PVSS.Share(s,n,t) n is sharenum value .t is threshold value
    coefficients = [secret] + [random_scalar() for j in range(threshold)]

    # Shamir 分享
    def f(x: int) -> int:
        """ evaluation function for secret polynomial
        """
        return (
                sum(coef * pow(x, j, CURVE_ORDER) for j, coef in enumerate(coefficients)) % CURVE_ORDER
        )

    indices = [i for i in range(1, sharenum + 1)]
    shares = {x: f(x) for x in indices}
    return shares


def Share(s_j, H1, pk, n, t):
    SSShare = share_secret(s_j, n, t)  # voter PVSS.share=(v,c)
    # 注意数组第一位为0，v，c数组长度为n+1
    v = [0]
    c = [0]
    DLEQ_Proof = [0]
    v.extend([multiply(H1, SSShare[i + 1]) for i in range(0, n)])  # v_i=H1^s_i
    c.extend([multiply(pk[i], SSShare[i + 1]) for i in range(0, n)])  # c_i=pk_i^s_i

    DLEQ_Proof.extend([DLEQ(H1, v[i + 1], pk[i], c[i + 1], SSShare[i + 1]) for i in range(0, n)])
    # DLEQ的Proof为证明v，c确实是由该多项式f(x)所生成,例如s_i=f(i)。
    res = {"v": v[1:], "c": c[1:], "DLEQ_Proof": DLEQ_Proof[1:]}
    return res


def DLEQ(x1, y1, x2, y2, alpha: int) -> Tuple[int, int]:  # 生成DLEQ承诺 alpha为需要承诺的零知识证明的值
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


def DLEQ_verify(x1, y1, x2, y2, challenge: int, response: int) -> bool:  # DLEQ_Verify的链下验证函数
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


# 删除PVSS.DVerify和PVSS.PVerify，将此部分转移到链上进行

def Decrypt(c_ji, sk_i):  # PVSS.Decrypt，解密函数
    sh1_ji = multiply(c_ji, sympy.mod_inverse((sk_i) % CURVE_ORDER, CURVE_ORDER))
    return sh1_ji


def Reconstruct(res, n, t):  # PVSS.Reconstruct  秘密恢复函数
    recIndex = [i + 1 for i in range(0, t + 1)]
    print(recIndex)
    sum = multiply(H1, 0)

    for i in recIndex:
        print("i", i, util.lagrange_coefficient(i))
        sum = add(sum, multiply(res["v"][i], util.lagrange_coefficient(i)))
    return sum


def ShareForCandidates(s_j, H1, pk, n, t):
    SSShare = share_secret(s_j, n, t)  # voter PVSS.share=(v,c)
    # 注意数组第一位为0，v，c数组长度为n+1
    v = [0]
    c = [0]
    DLEQ_Proof = [0]
    v.extend([multiply(H1, SSShare[i + 1]) for i in range(0, n)])  # v_i=H1^s_i
    c.extend([multiply(pk[i], SSShare[i + 1]) for i in range(0, n)])  # c_i=pk_i^s_i

    DLEQ_Proof.extend([DLEQ(H1, v[i + 1], pk[i], c[i + 1], SSShare[i + 1]) for i in range(0, n)])
    # DLEQ的Proof为证明v，c确实是由该多项式f(x)所生成,例如s_i=f(i)。
    res = {"v": v[1:], "c": c[1:], "DLEQ_Proof": DLEQ_Proof[1:], "P_j": SSShare}
    return res

    # reed solomon check


def RScodeVerify(res):
    # print(len(res))
    def coefficient(i: int) -> int:
        result = 1
        for j in range(1, len(res) + 1):
            # print(j)
            # j=j-1
            if i != j:
                # print(i - j) #% CURVE_ORDER)
                result *= sympy.mod_inverse((i - j) % CURVE_ORDER, CURVE_ORDER)
                result %= CURVE_ORDER
        return result

    codeword = []
    for i in range(1, len(res) + 1):
        codeword.append(coefficient(i) % CURVE_ORDER)
    sum = multiply(H1, 1)
    for i in range(0, len(res)):
        sum = add(sum, multiply(res[i], codeword[i]))
    if (sum == H1):
        return 1
    else:
        return 0
    #   sum=add(sum, multiply(res["v"][i], lagrange_coefficient(i)))


# 留作测试所用
"""
if __name__ == '__main__':
    key = Setup(10, 5)

    n = 10
    t = 5
    print("answer")
    print(multiply(H1, 233333))
    print(".........")

    res1 = Share(233333, H1, key["pk"], n, t)
    #print(res1)
    #print(RScodeVerify(res1["v"]))
    print(Contract.functions.inv(0,CURVE_ORDER).call())
    print(sympy.mod_inverse((1) % CURVE_ORDER, CURVE_ORDER))
    #sum = Reconstruct(res1, 10, 5)
    #print(sum)
    # Reconstruct2(res1,10,5)
    def test():
        for i in range((CURVE_ORDER-30),(CURVE_ORDER+31)):

            print(i,"..........",sympy.mod_inverse((i) % CURVE_ORDER, CURVE_ORDER))

    #test()
"""








