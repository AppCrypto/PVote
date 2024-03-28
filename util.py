def Point2IntArr(x):  # tuple/list transform to int[]
    ints = [int(num) for num in x]  # 数据转换
    return ints

def Dataconvert(res, n):  # Data conversion functions for bilinear pairing on-chain
    c1 = []
    c2 = []
    # 数据格式转换，将c(x,y)分开放入单组的数组中
    v1 = []
    v2 = []
    # 数据格式转换，将v(x,y)分开放入单组的数组中
    # #注意，因为Share的c，v从1开始为有效数据，所以输出的数组第一位都为0
    c1.extend(int(res["c"][i][0]) for i in range(0, n ))
    c2.extend(int(res["c"][i][1]) for i in range(0, n ))
    v1.extend(int(res["v"][i][0]) for i in range(0, n ))
    v2.extend(int(res["v"][i][1]) for i in range(0, n ))
    return {"c1": c1, "c2": c2, "v1": v1, "v2": v2}  # c1 is x of c, c2 is y of c. And v1,v2,s1,s2 so on...

def U_jdataconvert(U_j):
    U_j1 = []
    U_j2 = []
    U_j1.extend(int(U_j[i][0]) for i in range(0, len(U_j)))
    U_j2.extend(int(U_j[i][1]) for i in range(0, len(U_j)))

    return{"U_j1":U_j1,"U_j2":U_j2}

