pragma solidity ^0.8.0;


contract DAOsForVote {
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

    uint256 constant Tallires = 30;    //取消Init函数，需要提前设置唱票者人数，生成对应的存储空间

    struct G1Point {
		uint X;
		uint Y;
	}

	// Encoding of field elements is: X[0] * z + X[1]
	struct G2Point {
		uint[2] X;
		uint[2] Y;
	}

    uint256[2][] Tallires_pk;


    uint256[2]   AGGPointU;   //存储Aggregate之后的数据U*
    uint256[2][] AGGPointC;   //存储Aggregate之后的数据C*
    uint256[2][] AGGPointV;   //存储Aggregate之后的数据V*

    constructor() {
        AGGPointU = [0,0]; // 在构造函数中为该数组赋值
        uint256[2] memory temp;
        temp[0]=0;
        temp[1]=0;   // Initialize G1Point

        //构造函数，生成Tallires所需的存储空间

        for (uint i = 0; i < Tallires; i++) {
            uint256[2] memory newStruct = temp;
            AGGPointC.push(newStruct);
            AGGPointV.push(newStruct);
        }
    }


    struct Vote_Data
    {
        uint256[]  c1;
        uint256[]  c2;
        uint256[]  v1;
        uint256[]  v2;
        uint256 U1;
        uint256 U2;
        uint256[2][] D_Proof;
        //uint256[2][] P_Proof;
        //uint256[]  lagrangeCoefficient;
    }

    struct ZKRP_Proof
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
    }

    Vote_Data public VoteData;   //生成实例
    ZKRP_Proof public ZKRPProof;

    function PVSStoSC(uint256[] memory  _c1 , uint256[] memory _c2, uint256[] memory _v1, uint256[] memory _v2, uint256 _U1, uint256 _U2, uint256[2][] memory _D_Proof)
    public
    {   // 实例化结构体并赋值

        VoteData = Vote_Data({
            c1: _c1,
            c2: _c2,
            v1: _v1,
            v2: _v2,
            U1: _U1,
            U2: _U2,
            D_Proof: _D_Proof
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
        uint256 _z3
    ) public { //ZKRP Proof 赋值
        ZKRPProof = ZKRP_Proof({
            E_j: _E_j,
            F_j1: _F_j1,
            F_j2: _F_j2,
            U1_j: _U1_j,
            C1_j: _C1_j,
            c: _c,
            z1: _z1,
            z2: _z2,
            z3: _z3
        });
    }

    function PK_Setup(uint256[2][] memory pk) public
    {
        Tallires_pk = pk;
    }

    function ReturnData() public  returns (uint256[2][] memory, uint256[2][] memory, uint256[2] memory) {
        return (AGGPointC, AGGPointV, AGGPointU);
    }


    function ReturnPointC() public returns(uint[2][] memory)
    {
        return AGGPointC;
    }

    function Aggregate()
    public
    {
        uint elements = VoteData.c1.length;  //get array length
        for(uint i=0; i<elements; i++)
        {
            AGGPointC[i]=bn128_add([VoteData.c1[i], VoteData.c2[i], AGGPointC[i][0], AGGPointC[i][1]]);
            AGGPointV[i]=bn128_add([VoteData.v1[i], VoteData.v2[i], AGGPointV[i][0], AGGPointV[i][1]]);
        }
        AGGPointU=bn128_add([VoteData.U1,VoteData.U2,AGGPointU[0],AGGPointU[1]]);
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


    function PVSS_PVerify() public returns(bool)
    {

    }

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

    /*
    function ZKRP_verify(
        uint256[2] memory E_j, uint256[2] memory F_j1, uint256[2] memory F_j2,
        uint256[2] memory U1_j, uint256[2] memory C1_j,
        uint256 c, uint256 z1,
        uint256 z2, uint256 z3,
        uint256[2][] memory V, uint256[] memory lagrange_coefficient, uint256[2] memory U_j
    )
    public returns (bool proof1, bool proof2, bool proof3)
    {
        uint256[2] memory  C_j = Interpolate(V, lagrange_coefficient);
        proof1 = ZKRP_verify1(C1_j, C_j, c, z3);
        proof2 = ZKRP_verify2(U1_j, U_j, c, z3, z1);
        //proof3 = ZKRP_verify3(F_j1, F_j2, E_j, c, z1, z2); //c, z1, z2, pk_I
        //proof_is_valid = proof1 && proof2;
    }
    */

    function  Interpolate(
        uint256[2][] memory V, uint256[] memory lagrange_coefficient, uint256 G1X, uint G1Y, uint t
    )
    public returns (uint256[2] memory)
    {
        uint256[2] memory a1; //acc=g
        uint256[2] memory temp;
        a1[0] = G1X;
        a1[1] = G1Y;
        //uint elements=lagrange_coefficient.length;//to get the array length
        uint elements=t;
        for(uint i=0;i<elements;i++)
        {
            temp = bn128_multiply([V[i][0], V[i][1],lagrange_coefficient[i]]);
            a1 = bn128_add([a1[0], a1[1], temp[0], temp[1]]);
        }
        a1=bn128_add([a1[0], a1[1], G1X, G1neg(G1Y)]); // add neg(H1)
        return a1;
    }


    function ZKRP_verify1(
        uint256[2][] memory V, uint256[] memory lagrange_coefficient, uint256[2] memory C1_j,
        uint256 c, uint256 z3
    )
    public returns (bool proof_is_valid)
    {
        uint elements=lagrange_coefficient.length;
        uint256[2] memory  C_j = Interpolate(V, lagrange_coefficient, H1x, H1y,elements);
        uint256[2] memory temp;
        uint256[2] memory temp2;
        uint256[2] memory temp3;
        temp2 = bn128_multiply([C_j[0], C_j[1], c]);  // NEED C_j to  Interpolate
        temp3 = bn128_multiply([H1x, H1y, z3]);
        temp = bn128_add([temp2[0], temp2[1], temp3[0], temp3[1]]);
        proof_is_valid = C1_j[0] == temp[0] && C1_j[1] == temp[1];
    }

    function ZKRP_verify2(
        uint256[2] memory U1_j, uint256[2] memory U_j,
        uint256 c, uint256 z1, uint256 z3
    )
    public returns (bool proof_is_valid)
    {
        uint256[2] memory temp;
        uint256[2] memory temp2;
        uint256[2] memory temp3;
        uint256[2] memory temp4;

        temp2 = bn128_multiply([U_j[0], U_j[1], c]);
        temp3 = bn128_multiply([G1x, G1y, z3]);
        temp4 = bn128_multiply([H1x, H1y, z1]);


        temp = bn128_add([temp2[0], temp2[1], temp3[0], temp3[1]]);
        temp = bn128_add([temp[0], temp[1], temp4[0],temp4[1]]);
        proof_is_valid=U1_j[0] == temp[0] && U1_j[1] == temp[1];

    }

    function ZKRP_verify3(
        uint256[2] memory F_j1, uint256[2] memory F_j2,
        uint256[2] memory E_j,
        uint256 c, uint256 z1, uint256 z2
    )
    public returns (bool)
    {
        G1Point[] memory p1 = new G1Point[](5);
		G2Point[] memory p2 = new G2Point[](5);

        uint256[2] memory temp1;
        uint256[2] memory temp2;
        uint256[2] memory temp3;
        temp1 = bn128_multiply([E_j[0], E_j[1], c]);
        temp2 = bn128_multiply([E_j[0], E_j[1], (GROUP_ORDER-z1)]);
        temp3 = bn128_multiply([1, 2, z2]);

        p1[0].X = F_j1[0];
        p1[0].Y = F_j1[1];

        p1[1].X = F_j2[0];
        p1[1].Y = F_j2[1];


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
    /*
    function Tally(
        uint256[2][] memory C, uint256[] memory lagrange_coefficient, uint t
    )
    public returns(uint256[2] memory)
    {
        uint256[2] memory G1ACC;
        G1ACC =  Interpolate(C, lagrange_coefficient, G1x, G1y, t);
        G1ACC =  bn128_add([G1PointU[0], G1PointU[1], G1ACC[0], G1neg(G1ACC[1])]);
        G1ACC =  bn128_add([G1ACC[0], G1ACC[1], G1x, G1neg(G1y)]);
        return G1ACC;
    }
    */
}