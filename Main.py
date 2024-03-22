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
from py_ecc.bn128 import add, multiply, neg, pairing, is_on_curve
from py_ecc.bn128 import curve_order as CURVE_ORDER
from py_ecc.bn128 import field_modulus as FIELD_MODULUS
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
H1 = multiply(G1, 9868996996480530350723936346388037348513707152826932716320380442065450531909)  # 生成元H1
pks = []


def Vj_Vote(w_j: int, n: int, t: int):  # w_j 为 vote value  #函数定义了一个投票者V_j应该完成的事务，（n，t）为秘密分享参与人数和门限值
    s_j = PVSS.random_scalar()
    # 生成随机数
    starttime = time.time()
    shares = PVSS.Share(s_j, H1, pk, n, t)
    t1 = time.time() - starttime
    print("PVSS_Share spent time:", t1)
    # 调用PVSS.Share

    U_j = add(multiply(H1, w_j), multiply(G1, s_j))
    # 链下计算U_j
    starttime2 = time.time()
    zkrp_proof = ZKRP.Prove(s_j, w_j, U_j, GPK["sigam_k"][(w_j) % (b + 1)])
    # 为所生成的数据调用ZKRP.Prove生成对应的承诺
    t2 = time.time() - starttime2
    print("ZKRP.Prove spent time:", t2)

    dleq_proof = []
    for i in range(0, n):
        temp = util.Point2IntArr(shares["DLEQ_Proof"][i])
        dleq_proof.extend([temp])
    # 为DLEQ Proof数据格式转换
    agg = PVSS.Dateconvert(shares, n)  # Data transformation  数据转换
    # 将投票者生成的 PVSS.Share的v，c数组，dleq_proof数组，U_j传输到智能合约上
    Contract.functions.PVSStoSC(agg["c1"], agg["c2"], agg["v1"], agg["v2"], int(U_j[0]), int(U_j[1]),
                                dleq_proof).transact({'from': w3.eth.accounts[0]})
    # 将投票者生成的ZKRP.Prove生成的Proof传输到智能合约上
    Contract.functions.ZKRPtoSC(util.Point2IntArr(zkrp_proof[0]), util.Point2IntArr(zkrp_proof[1]),
                                util.Point2IntArr(zkrp_proof[2]), util.Point2IntArr(zkrp_proof[3]),
                                util.Point2IntArr(zkrp_proof[4]), zkrp_proof[5], zkrp_proof[6], zkrp_proof[7],
                                zkrp_proof[8]).transact({'from': w3.eth.accounts[0]})
    print("Vote value:", w_j)
    # return shares["v"][1:]  # 给ZKRP.Verify提供V_j ,因为生成的v数组第一位为无效0，智能合约上没有像Python方便的操作


def Ti_Tally(No: int, pk_i, sk_i):  # 函数定义了一个唱票者Tallier T_i应该完成的事务
    # 从链上下载累积好的V,C数据
    aggCV = Contract.functions.DownloadAGGVC(No).call()

    # 数据转换
    C_i = (FQ(aggCV[0][0]), FQ(aggCV[0][1]))

    # 从链上获取属于唱票者T_i的公钥
    # pk_i = Contract.functions.ReturnPKi(No).call()
    # pk_i = (FQ(pk_i[0]), FQ(pk_i[1]))

    # 调用PVSS.Decrypt函数为累积份额C解密
    starttime = time.time()
    sh1 = PVSS.Decrypt(C_i, sk_i)
    t1 = time.time() - starttime
    print("PVSS_Decrypt spent time:", t1)
    # 生成DLEQ P_Proof,证明是该唱票者T_i所解密的份额c
    proof = PVSS.DLEQ(G1, pk_i, sh1, C_i, sk_i)

    # 把解密份额和P_Proof上传到链上，通过验证后则将解密份额保留在链上DecryptedShare数组中
    Contract.functions.Decrypted_SharetoSC(No, util.Point2IntArr(sh1), util.Point2IntArr(proof)).transact(
        {'from': w3.eth.accounts[0]})

    gas_estimate = Contract.functions.Decrypted_SharetoSC(No, util.Point2IntArr(sh1),
                                                          util.Point2IntArr(proof)).estimateGas()
    print("PVSS_PVerify gas cost:", gas_estimate)
    print("Tallier", No, "done")


def Aggreagate():  # 执行一次链上的Aggregate函数，将V,C数据聚合到链上，因为v，c数据已经保留在链上了，所以无需参数输入

    Contract.functions.Aggregate().transact({'from': w3.eth.accounts[0]})
    gas_estimate = Contract.functions.Aggregate().estimateGas()
    print("Aggregate gas cost:", gas_estimate)
    print("Aggregate done.")


def Tally(m):  # 链上唱票，输入参数m（投票人数）是因为要确定投票数值范围（a*m,b*m)

    # 得到投票结果
    result = Contract.functions.Tally().call()
    gas_estimate = Contract.functions.Tally().estimateGas()
    print("Tally gas cost:", gas_estimate)
    # 计算所有投票值的可能
    AllResult = {}
    for i in range(a * m, b * m + 1):
        AllResult[i] = multiply(H1, i)
    # 将投票结果和投票可能值遍历做比对
    for i in range(a * m, b * m + 1):
        if (AllResult[i] != None and result[0] == AllResult[i][0] and result[1] == AllResult[i][1]):
            # print("The vote score is",i)
            return i 

    print("No vote result")
    return -100000


def ReturnDate():  # 返回当前所聚合的AGG的数据，测试所用
    agg = Contract.functions.ReturnData().call({'from': w3.eth.accounts[0]})
    print("Agg C:" + str(agg[0]))
    print("Agg V:" + str(agg[1]))
    print("Agg U:" + str(agg[2]))


def ReturnPK():  # 返回链上公钥，测试所用
    res = Contract.functions.ReturnPK().call({'from': w3.eth.accounts[0]})
    print(res)
    # print("pk-onchain:"+str(res))


def ZKRP_verify2(n, t):  # ZKRP的链上验证

    # print(V)
    result1 = Contract.functions.ZKRP_verify1(t).call()  # ZKRP.Verify的第一个等式的验证
    result2 = Contract.functions.ZKRP_verify2().call()  # ZKRP.Verify的第二个等式的验证
    result3 = Contract.functions.ZKRP_verify3().call()  # ZKRP.Verify的第三个等式的验证
    """
    if result1 and result2 and result3:  #三个验证等式全true，ZKRP.Verify才会返回true
        return (True)
    else:
        return (False)
    """
    return (result1, result2, result3)


if __name__ == '__main__':

    n = int(sys.argv[1])  # 唱票者人数n
    t = int(n/2+1)  # 门限值t
    print("...........................................Setup phase.............................................",n,t)
    starttime=time.time()
    key = PVSS.Setup(n, t)  # PVSS Key Generation
    print("PVSS.setup of each tallier average time:", "%.2f"%((time.time()- starttime)/n*1000),"ms", "a public key size:", "%.2f"%(len(str(key["pk"]))/n),"B")
    pk = key["pk"]  # Public key array
    sk = key["sk"]  # Private key array
    pks = [util.Point2IntArr(pk[i]) for i in range(n)]  # 公钥数据格式转换
    # 将公钥上传到智能合约

    Contract.functions.setTalliresPK(pks).transact({'from': w3.eth.accounts[0]})
    gas_estimate = Contract.functions.setTalliresPK(pks).estimateGas()
    print("Initiator setup gas cost:",gas_estimate)
    print("Initiator setup output size:","%.2f" %(len(str(pks))/1024),"kB")
    # exit()
    a = 0  # 投票最小范围a
    b = 5  # 投票最大范围b
    m = 3  # 参与投票人数
    GPK = ZKRP.Setup(a, b)  # ZKRP初始化

    print("............................................Voting phase...........................................")

    # 目前用的是ZKRP_verify2，还有V_ji的这部分需要继续优化
    # 第一个投票者
    ballot = 0
    for i in range(0, m):
        w_j = int(random.random() * (b - a + 1) + a)
        ballot += w_j
        Vj_Vote(w_j, n, t)  # 投票者投票函数
        if (Contract.functions.PVSS_DVerify().call() and Contract.functions.ZKRP_verify(t + 1).call()):
            print("Both PVSS_DVerify result and ZKRP_Verify result return true")
            gas_estimate1 = Contract.functions.PVSS_DVerify().estimateGas()
            print("PVSS.DVerify gas cost:", gas_estimate1)

            # gas_estimate2 = Contract.functions.ZKRP_verify( t+1 ).estimateGas()
            # print("ZKRP.Verify gas cost:",gas_estimate2)
            Aggreagate()  # 通过两个验证后将投票上传的数据聚合
        else:
            print("Invalid vote value w_j......did not aggreagate in SC")
        # print("PVSS_DVerify result:", Contract.functions.PVSS_DVerify().call())  # 链上PVSS.DVerify
        # print("ZKRP_Verify result:", Contract.functions.ZKRP_verify( t+1 ).call()) # 链上ZKRP.Verify

    print("expected ballot value:", ballot)

    w_j = b+1
    Vj_Vote(w_j, n, t)  # 投票者投票函数
    if (Contract.functions.PVSS_DVerify().call() and Contract.functions.ZKRP_verify(t + 1).call()):
        print("Both PVSS_DVerify result and ZKRP_Verify result return true")
        Aggreagate()  # 通过两个验证后将投票上传的数据聚合
    else:
        print("Invalid vote value......did not aggreagate")

    print("..........................................tallying phase...........................................")

    temp_t = t+1   # for循环的方式生成唱票者，每当有一个唱票者完成Ti_Tally函数，则会在链上成功上传一份解密份额，当满足t个份额时可以调用Tally函数进行唱票
    for i in range(0, temp_t):
        Ti_Tally(i+1, pk[i], sk[i])
    # 用完成任务的唱票者数量来代替t的恢复门限值，比如想控制9个份额参与秘密恢复，则使得9个唱票者完成任务，即调用9个唱票者函数
    # 也可以全部列出，以表示t个唱票者完成任务  ，以下例子为10个唱票者完成任务

    # Tally(temp_t ,m)  #链上唱票
    tally=Tally(m)  # 链上唱票
    if tally==ballot:
        print("The tallying result is correct:",tally)
    print("............................................Reward phase...........................................")
