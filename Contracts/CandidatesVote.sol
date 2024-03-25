pragma solidity ^0.8.0;


contract CandidatesVote {
    ////////////////////////////////////////////////////////////////////////////////////////////////
    //// CRYPTOGRAPHIC CONSTANTS

    uint256 constant GROUP_ORDER   = 21888242871839275222246405745257275088548364400416034343698204186575808495617;
    uint256 constant FIELD_MODULUS = 21888242871839275222246405745257275088696311157297823662689037894645226208583;

    // definition of two indepently selected generator for the groups G1 and G2 over
    // the bn128 elliptic curve
    // TODO: maybe swap generators G and H
    uint256 constant G1x  = 1;
    uint256 constant G1y  = 2;

    uint256 constant negG1x = 1;
    uint256 constant negG1y = 21888242871839275222246405745257275088696311157297823662689037894645226208581;

    uint256 constant H1x  = 15264291051155210722230395084766962011373976396997290700295946518477517838363;
    uint256 constant H1y  = 18062169012241050521396281509436922807270932827014386397365657617881670284318;

    uint256 constant negH1x  = 15264291051155210722230395084766962011373976396997290700295946518477517838363;
    uint256 constant negH1y  = 3826073859598224700850124235820352281425378330283437265323380276763555924265;

    uint256 constant G2xx = 10857046999023057135944570762232829481370756359578518086990519993285655852781;
    uint256 constant G2xy = 11559732032986387107991004021392285783925812861821192530917403151452391805634;
    uint256 constant G2yx = 8495653923123431417604973247489272438418190587263600148770280649306958101930;
    uint256 constant G2yy = 4082367875863433681332203403145435568316851327593401208105741076214120093531;

    uint256 constant Tallires = 30;    //取消Init函数，需要提前设置唱票者人数，生成对应的存储空间,这里直接生成30个唱票者的空间

    struct G1Point {
		uint X;
		uint Y;
	}

	// Encoding of field elements is: X[0] * z + X[1]
	struct G2Point {
		uint[2] X;
		uint[2] Y;
	}




    uint256[2][] AGGPointU;   //存储Aggregate之后的数据U*
    uint256[2][] AGGPointC;   //存储Aggregate之后的数据C*
    uint256[2][] AGGPointV;   //存储Aggregate之后的数据V*
    uint256[2][] Tallires_pk;
    uint256[2][] DecryptedShare;

    constructor() {
        //AGGPointU = [0,0]; // 在构造函数中为该数组赋值
        uint256[2] memory temp;
        temp[0]=0;
        temp[1]=0;   // Initialize G1Point

        //构造函数，生成Tallires所需的存储空间

        for (uint i = 0; i < Tallires; i++) {
            uint256[2] memory newStruct = temp;
            AGGPointC.push(newStruct);
            AGGPointV.push(newStruct);

        }
        for (uint i=0 ; i< Tallires/2+1; i++)
        {
            uint256[2] memory newStruct = temp;
            AGGPointU.push(newStruct);
        }
    }


    struct PVSS_Data  //一个保存投票数据的数据结构类型
    {
        uint256[]  c1;
        uint256[]  c2;
        uint256[]  v1;
        uint256[]  v2;
        uint256[]  U1;   //现在进入的是U1数组
        uint256[]  U2;  //U2同理，长度应该为l，表示l个委员会投票承诺
        uint256[2][] D_Proof;
        uint256 ulen;
    }

    struct ZKRP_Proof  //一个保存ZKRP_Proof数据的数据结构类型
    {
        uint256[2] E_j;
        uint256[2] F_j1;
        uint256[2] F_j2;
        uint256[2] U1_j;
        uint256[2] C1_j;
        uint256 c;
        uint256 z1;
        uint256 z2;
        uint256 z3;
        uint256 U1;
        uint256 U2;
    }


    PVSS_Data  public VoteData;   //生成实例
    ZKRP_Proof public ZKRPProof;

    function PVSStoSC(uint256[] memory  _c1 , uint256[] memory _c2, uint256[] memory _v1, uint256[] memory _v2, uint256[] memory _U1, uint256[] memory _U2, uint256[2][] memory _D_Proof, uint256 _ulen)
    public
    {   // 实例化结构体并赋值
        //上传投票数据到链上的函数，为VoteData赋值
        VoteData = PVSS_Data({
            c1: _c1,
            c2: _c2,
            v1: _v1,
            v2: _v2,
            U1: _U1,
            U2: _U2,
            D_Proof: _D_Proof,
            ulen: _ulen
            //lagrangeCoefficient: _lagrangeCoefficient  , uint256[] memory _lagrangeCoefficient
        });


    }

    function ZKRPtoSC(
        uint256[2] memory _E_j,
        uint256[2] memory _F_j1,
        uint256[2] memory _F_j2,
        uint256[2] memory _U1_j,
        uint256[2] memory _C1_j,
        uint256 _c,
        uint256 _z1,
        uint256 _z2,
        uint256 _z3,
        uint256 _u1,
        uint256 _u2
    ) public { ///为ZKRP Proof 赋值，上传Proof到链上的函数
        ZKRPProof = ZKRP_Proof({
            E_j: _E_j,
            F_j1: _F_j1,
            F_j2: _F_j2,
            U1_j: _U1_j,
            C1_j: _C1_j,
            c: _c,
            z1: _z1,
            z2: _z2,
            z3: _z3,
            U1:_u1,
            U2:_u2
        });
    }


    function setTalliresPK(uint256[2][] memory pk) public {  //上传唱票者公钥函数
        Tallires_pk = pk; //保存到Tallires_pk数组
    }

     //下载单个所聚合的c，v的数据，No为第几个数
    function DownloadAGGVC(uint No) public returns (uint256[2] memory, uint256[2] memory)
    {
        return (AGGPointC[No-1],AGGPointV[No-1]);
    }

    //上传解密份额和dleq_Proof数据，调用PVSS_PVerify，通过则保留份额到DecryptedShare数组中
    function Decrypted_SharetoSC(uint No, uint256[2] memory DShare, uint256[2] memory P_Proof) public returns (bool)
    {
        if(PVSS_PVerify(Tallires_pk[No-1], DShare, AGGPointC[No-1], P_Proof))
        {
            DecryptedShare.push(DShare);
            return true;
        }
        return false;
    }

    //返回第i个唱票者的pk,可以优化掉，测试所用
    function ReturnPKi(uint i) public returns(uint256[2] memory)
    {
        return (Tallires_pk[i-1]);
    }

     //返回投票者的v数组的函数
    function ReturnV() public returns(uint256[] memory,uint256[] memory)
    {
        return (VoteData.v1,VoteData.v2);
    }

     //返回DecryptedShare数组和该数组长度的函数
    function ReturnDS() public returns(uint256[2][] memory,uint)
    {
        return (DecryptedShare,DecryptedShare.length);
    }


     //返回目前聚合了的c数组，测试所用
    function ReturnPointC() public returns(uint[2][] memory)
    {
        return AGGPointC;
    }

     //聚合函数，对v，c，U的聚合
    function Aggregate()
    public
    {
        uint elements = VoteData.c1.length;  //get array length
        for(uint i=0; i<elements; i++)
        {
            AGGPointC[i]=bn128_add([VoteData.c1[i], VoteData.c2[i], AGGPointC[i][0], AGGPointC[i][1]]);
            AGGPointV[i]=bn128_add([VoteData.v1[i], VoteData.v2[i], AGGPointV[i][0], AGGPointV[i][1]]);
        }

        for(uint i=0; i<VoteData.ulen; i++)  //对l个U_jk进行聚合
        {
            AGGPointU[i]=bn128_add([VoteData.U1[i],VoteData.U2[i],AGGPointU[i][0],AGGPointU[i][1]]);
        }

    }

    function P1() pure internal returns (G1Point memory) {
	    return G1Point(1, 2);
	}

    function P2() pure internal returns (G2Point memory) {
		return G2Point(
			[11559732032986387107991004021392285783925812861821192530917403151452391805634,
			 10857046999023057135944570762232829481370756359578518086990519993285655852781],
			[4082367875863433681332203403145435568316851327593401208105741076214120093531,
			 8495653923123431417604973247489272438418190587263600148770280649306958101930]
		);
	}

    	/// return the negation of p, i.e. p.add(p.negate()) should be zero.
	function G1neg(uint256 p) pure internal returns (uint r) {
		// The prime q in the base field F_q for G1
		uint256 q = 21888242871839275222246405745257275088696311157297823662689037894645226208583;
            r = (q - (p % q));
	}

    //返回投票发起者的公钥，为ZKRP.Verify所用
    function pk_I() pure internal returns (G2Point memory) {
		return G2Point(
			[410331679793378253270676581922857325978420412939832091293874594317968404676,
            18763800358864019000552773440545754218678238064612111471382832602429099653153],
			[15498345051290780549250954827254971939745244621482817595357857072172484368730,
            18428974775939053033653030813069583546036508418472013925185079794649562797187]
		);
	}

    function bn128_add(uint256[4] memory input)
    public returns (uint256[2] memory result) {
        // computes P + Q
        // input: 4 values of 256 bit each
        //  *) x-coordinate of point P
        //  *) y-coordinate of point P
        //  *) x-coordinate of point Q
        //  *) y-coordinate of point Q

        bool success;
        assembly {
            // 0x06     id of precompiled bn256Add contract
            // 0        number of ether to transfer
            // 128      size of call parameters, i.e. 128 bytes total
            // 64       size of call return value, i.e. 64 bytes / 512 bit for a BN256 curve point
            success := call(not(0), 0x06, 0, input, 128, result, 64)
        }
        require(success, "elliptic curve addition failed");
    }

    function bn128_multiply(uint256[3] memory input)
    public returns (uint256[2] memory result) {
        // computes P*x
        // input: 3 values of 256 bit each
        //  *) x-coordinate of point P
        //  *) y-coordinate of point P
        //  *) scalar x

        bool success;
        assembly {
            // 0x07     id of precompiled bn256ScalarMul contract
            // 0        number of ether to transfer
            // 96       size of call parameters, i.e. 96 bytes total (256 bit for x, 256 bit for y, 256 bit for scalar)
            // 64       size of call return value, i.e. 64 bytes / 512 bit for a BN256 curve point
            success := call(not(0), 0x07, 0, input, 96, result, 64)
        }
        require(success, "elliptic curve multiplication failed");
    }

    /// @return the result of computing the pairing check
	/// e(p1[0], p2[0]) *  .... * e(p1[n], p2[n]) == 1
	/// For example pairing([P1(), P1().negate()], [P2(), P2()]) should
	/// return true.
    function pairing(G1Point[] memory p1, G2Point[] memory p2) public returns (bool)
    {
		require(p1.length == p2.length);
		uint elements = p1.length;
		uint inputSize = elements * 6;
		uint[] memory input = new uint[](inputSize);
		for (uint i = 0; i < elements; i++)
		{
			input[i * 6 + 0] = p1[i].X;
			input[i * 6 + 1] = p1[i].Y;
			input[i * 6 + 2] = p2[i].X[0];
			input[i * 6 + 3] = p2[i].X[1];
			input[i * 6 + 4] = p2[i].Y[0];
			input[i * 6 + 5] = p2[i].Y[1];
		}
		uint[1] memory out;
		bool success;
		assembly {
			success := staticcall(sub(gas()	, 2000), 8, add(input, 0x20), mul(inputSize, 0x20), out, 0x20)
			// Use "invalid" to make gas estimation work
			//switch success case 0 { invalid }
		}
		require(success);
		return out[0] != 0;
	}

    //链上PVSS.DVerify函数
    function PVSS_DVerify() public returns(bool)
    {

        uint elements = VoteData.c1.length;  //get array length
        for (uint i = 0; i < elements; i++)
        {

            if(!DLEQ_verify([H1x,H1y],[VoteData.v1[i], VoteData.v2[i]],[Tallires_pk[i][0],Tallires_pk[i][1]],[VoteData.c1[i], VoteData.c2[i]],[VoteData.D_Proof[i][0],VoteData.D_Proof[i][1]]))
            {
                return false;
            }
        }

        return true;
    }

    //链上PVSS.PVerify函数
    function PVSS_PVerify(uint256[2] memory pk, uint256[2] memory sh_i, uint256[2] memory c_i, uint256[2] memory Proof ) public returns(bool)
    {
        if(!DLEQ_verify([G1x,G1y],[pk[0], pk[1]],[sh_i[0],sh_i[1]],[c_i[0], c_i[1]],[Proof[0],Proof[1]]))
        {
            return false;
        }

        return true;
    }

     //链上DLEQ验证
    function DLEQ_verify(
        uint256[2] memory x1, uint256[2] memory y1,
        uint256[2] memory x2, uint256[2] memory y2,
        uint256[2] memory proof
    )
    public returns (bool proof_is_valid)
    {
        uint256[2] memory tmp1;
        uint256[2] memory tmp2;

        tmp1 = bn128_multiply([x1[0], x1[1], proof[1]]);
        tmp2 = bn128_multiply([y1[0], y1[1], proof[0]]);
        uint256[2] memory a1 = bn128_add([tmp1[0], tmp1[1], tmp2[0], tmp2[1]]);

        tmp1 = bn128_multiply([x2[0], x2[1], proof[1]]);
        tmp2 = bn128_multiply([y2[0], y2[1], proof[0]]);
        uint256[2] memory a2 = bn128_add([tmp1[0], tmp1[1], tmp2[0], tmp2[1]]);

        uint256 challenge = uint256(keccak256(abi.encodePacked(a1, a2, x1, y1, x2, y2)));
        proof_is_valid = challenge == proof[0];
    }


    function ZKRP_verify(uint8 x, uint8 t)
    public returns (bool)
    {
        uint256[2][] memory V;
        V=mergeArrays(VoteData.v1,VoteData.v2);
        uint256[] memory lagrange_coefficient;
        lagrange_coefficient = lagrangeCoefficient(x,t);
        //uint elements=lagrange_coefficient.length;
        uint256[2] memory  C_j = Interpolate(V, lagrange_coefficient);
        bool  proof1 = ZKRP_verify1(C_j,t);
        bool  proof2 = ZKRP_verify2();
        bool  proof3 = ZKRP_verify3();

        return proof1 && proof2 && proof3;
        //return (proof1,proof2,proof3);
    }


     //链上插值函数
    function  Interpolate(
        uint256[2][] memory V, uint256[] memory lagrange_coefficient
    )
    public returns (uint256[2] memory)
    {
        uint256[2] memory a1; //acc=g
        uint256[2] memory temp;
        a1[0] = 0;
        a1[1] = 0;
        //uint elements=lagrange_coefficient.length;//to get the array length
        uint elements=lagrange_coefficient.length;
        for(uint i=0;i<elements;i++)
        {
            temp = bn128_multiply([V[i][0], V[i][1],lagrange_coefficient[i]]);
            a1 = bn128_add([a1[0], a1[1], temp[0], temp[1]]);
        }
        //a1=bn128_add([a1[0], a1[1], G1X, G1neg(G1Y)]); // add neg(H1)
        return a1;
    }

    function mergeArrays(uint[] memory arr1, uint[] memory arr2) public  returns (uint[2][] memory) {
        require(arr1.length == arr2.length, "Input arrays must have the same length");

        uint[2][] memory mergedArray = new uint[2][](arr1.length);

        for (uint i = 0; i < arr1.length; i++) {
            mergedArray[i][0] = arr1[i];
            mergedArray[i][1] = arr2[i];
        }

        return mergedArray;
    }


     //ZKRP第一个等式的链上验证
    function ZKRP_verify1(uint256[2] memory C_j,uint8 t)
    public returns (bool proof_is_valid)
    {

        uint256[2] memory temp;
        uint256[2] memory temp2;
        uint256[2] memory temp3;
        temp2 = bn128_multiply([C_j[0], C_j[1], ZKRPProof.c]);  // NEED C_j to  Interpolate
        temp3 = bn128_multiply([H1x, H1y, ZKRPProof.z3]);
        temp = bn128_add([temp2[0], temp2[1], temp3[0], temp3[1]]);
        proof_is_valid = ZKRPProof.C1_j[0] == temp[0] && ZKRPProof.C1_j[1] == temp[1];
    }

     //ZKRP第二个等式的链上验证
    function ZKRP_verify2()
    public returns (bool proof_is_valid)
    {
        uint256[2] memory temp;
        uint256[2] memory temp2;
        uint256[2] memory temp3;
        uint256[2] memory temp4;

        temp2 = bn128_multiply([ZKRPProof.U1, ZKRPProof.U2, ZKRPProof.c]);
        temp3 = bn128_multiply([G1x, G1y, ZKRPProof.z3]);
        temp4 = bn128_multiply([H1x, H1y, ZKRPProof.z1]);


        temp = bn128_add([temp2[0], temp2[1], temp3[0], temp3[1]]);
        temp = bn128_add([temp[0], temp[1], temp4[0],temp4[1]]);
        proof_is_valid=ZKRPProof.U1_j[0] == temp[0] && ZKRPProof.U1_j[1] == temp[1];
    }

    //ZKRP第三个等式的链上验证，为双线性配对验证
    function ZKRP_verify3()
    public returns (bool)
    {
        G1Point[] memory p1 = new G1Point[](5);
		G2Point[] memory p2 = new G2Point[](5);

        uint256[2] memory temp1;
        uint256[2] memory temp2;
        uint256[2] memory temp3;
        temp1 = bn128_multiply([ZKRPProof.E_j[0], ZKRPProof.E_j[1], ZKRPProof.c]);
        temp2 = bn128_multiply([ZKRPProof.E_j[0], ZKRPProof.E_j[1], (GROUP_ORDER-ZKRPProof.z1)]);
        temp3 = bn128_multiply([G1x, G1y, ZKRPProof.z2]);

        p1[0].X = ZKRPProof.F_j1[0];
        p1[0].Y = ZKRPProof.F_j1[1];

        p1[1].X = ZKRPProof.F_j2[0];
        p1[1].Y = ZKRPProof.F_j2[1];


        p1[2].X = temp1[0];
        p1[2].Y = temp1[1];

        p1[3].X = temp2[0];
        p1[3].Y = temp2[1];

        p1[4].X = temp3[0];
        p1[4].Y = temp3[1];

        p2[0] = P2();
        p2[1] = P2();
        p2[2] = pk_I();
        p2[3] = P2();
        p2[4] = P2();

		return pairing(p1, p2);
    }


    //求逆函数，在群范围内求逆
    function inv(uint256 a, uint256 prime) public returns (uint256){
    	return modPow(a, prime-2, prime);
    }

    function modPow(uint256 base, uint256 exponent, uint256 modulus) internal returns (uint256) {
	    uint256[6] memory input = [32,32,32,base,exponent,modulus];
	    uint256[1] memory result;
	    assembly {
	      if iszero(call(not(0), 0x05, 0, input, 0xc0, result, 0x20)) {
	        revert(0, 0)
	      }
	    }
	    return result[0];
	}

    //目前使用的是这个，我晚上回来继续修改，生成朗日插值系数函数(生成拉格朗日插值并返回), x是恢复f(x),当x取0时为f(0),即s_j
    function lagrangeCoefficient(uint256 x, uint256 t) public returns (uint256[] memory){

        uint256[] memory lar2 = new uint256[](t);
        uint256 result = 1;
        uint256 inverse = 0;
        uint256 intermediate_result = 0;
        for (uint256 i = 1; i< t+1 ; i++)
        {
            result=1;
            for (uint256 j = 1; j < t+1;j++) {
                if (i != j) {
                    //inverse = inv((j.sub(i)).mod(GROUP_ORDER), GROUP_ORDER);
                    inverse = inv(((j+GROUP_ORDER-i)%GROUP_ORDER), GROUP_ORDER);//%GROUP_ORDER
                    intermediate_result = mulmod((j+GROUP_ORDER-x),inverse,GROUP_ORDER);
                    result = mulmod(result,intermediate_result,GROUP_ORDER);
                }
            }
            lar2[i-1]=result;
        }
        return lar2;
    }

    //链上唱票函数,这次不再需要任何参数的链下输入, i为第几位委员会唱票，利用的是第i个聚合U_jk
    function Tally(uint8 i)
    public returns(uint256[2] memory)
    {
        uint256[] memory lagrange_coefficient;
        lagrange_coefficient = lagrangeCoefficient(i,DecryptedShare.length);
        uint256[2] memory G1ACC;
        G1ACC =  Interpolate(DecryptedShare, lagrange_coefficient);
        G1ACC =  bn128_add([AGGPointU[i][0], AGGPointU[i][1], G1ACC[0], G1neg(G1ACC[1])]);
        //G1ACC =  bn128_add([G1ACC[0], G1ACC[1], G1x, G1neg(G1y)]);
        return G1ACC;
    }

    /*
    function ZKRP_ForGasTest(uint8 t) public returns (uint256[2] memory C_j)
    {
        uint256[2][] memory V;
        V=mergeArrays(VoteData.v1,VoteData.v2);
        uint256[] memory lagrange_coefficient;
        lagrange_coefficient = lagrangeCoefficient(t);
        //uint elements=lagrange_coefficient.length;
        C_j = Interpolate(V, lagrange_coefficient);
    }
    */
}