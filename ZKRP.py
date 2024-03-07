import PVSS
import time
import re
import sympy # Needed for mod_inverse
from web3 import Web3
from py_ecc.bn128 import G1, G2
from py_ecc.bn128 import add, multiply, neg, pairing, is_on_curve
from py_ecc.bn128 import curve_order as CURVE_ORDER
from py_ecc.bn128 import field_modulus as FIELD_MODULUS

global pk,sk,n,t
n=10    #tally_people /Registered Tallier
t=5     #The voting system need tallier to recover secret
H1=multiply(G1,2)
keccak_256 = Web3.solidityKeccak


key=PVSS.Setup()
pk=key["pk"]
sk=key["sk"]

shares=PVSS.Share(PVSS.random_scalar(),H1,pk)


#starttime = time.time()
#SS=PVSS.DVerify(H1,shares["v"],pk,shares["c"],shares["raw"])
#print("SCRAPE DDH verification cost ",time.time()- starttime)    
#print(SS)


def Setup(a,b):
    sk_I=PVSS.random_scalar()
    pk_I=multiply(G2,sk_I)
    sigma_k=[]
    for i in range(a,b+1):
        temp=multiply(G1,sympy.mod_inverse((sk_I+i) % CURVE_ORDER, CURVE_ORDER))
        sigma_k.extend([temp])
    return {"sk_I":sk_I,"pk_I":pk_I,"sigam_k":sigma_k}

def Prove(s_j,w_j,U_j,sigma_wj):
    v = PVSS.random_scalar()
    s = PVSS.random_scalar()
    t = PVSS.random_scalar()
    m = PVSS.random_scalar()
    #m1 = PVSS.random_scalar()

    E_j = multiply(sigma_wj,v)

    par1=pairing(G2,multiply(E_j,CURVE_ORDER-s))*pairing(G2,multiply(G1,t))

    def F_jTransform(par):
    #F_j = pairing(G2,G1)*pairing
        temp = re.findall("\d+",str(par))
        tempF_j=[]
        for i in range (len(temp)):
            tempF_j.append(int(temp[i]))
        return tuple(tempF_j)
    
    F_j=F_jTransform(par1)

    U1_j = add(multiply(H1,s), multiply(G1,m))
    C1_j = multiply(H1,m)
    #D_j =  add(multiply(L1,m),add(multiply(E_j,CURVE_ORDER-s),multiply(G1,t)))
    #G_j = multiply(add(L1,neg(H1)),m1)

    c = keccak_256(
        abi_types=["uint256"] * 20,
        values=[
            int(v)
            for v in (E_j)
            + (U_j)
            + (F_j)
            + (U1_j)
            + (C1_j)
        ],
    )
    c = int.from_bytes(c, "big")
    z1 = (s - w_j * c) % CURVE_ORDER
    z2 = (t - v * c) % CURVE_ORDER
    z3 = (m - s_j * c) % CURVE_ORDER
    #z4 = (m1 - m * c) % CURVE_ORDER
    return E_j,F_j,U1_j,C1_j,c,z1,z2,z3,par1

def Verify(proof,V_j,U_j,s_j,pk_I):
    C_j=multiply(H1,s_j)
    
    if(proof[3]!=add(multiply(C_j,proof[4]),multiply(H1,proof[7]))):
        print("no1")
        return 0
    
    if(proof[2]!=add(multiply(U_j,proof[4]),add(multiply(G1,proof[7]),multiply(H1,proof[5])))):
        print("no2")
        return 0
    
    if(proof[8]!=(pairing(pk_I,multiply(proof[0],proof[4]))*pairing(G2,multiply(proof[0],CURVE_ORDER-proof[5]))*pairing(G2,multiply(G1,proof[6])))):
        print("no3")
        return 0
    
    return 1 

GPK=Setup(1,5)
#print(GPK["sigam_k"][0])
s_j=PVSS.random_scalar()
w_j=2
U_j=add(multiply(H1,w_j), multiply(G1,s_j))
proof = Prove(s_j,w_j,U_j,GPK["sigam_k"][0])
#print(proof)

result=Verify(proof,G1,U_j,s_j,GPK["pk_I"])
print(result)

aa=add(G1,G1)
print(aa)



