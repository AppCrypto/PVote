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

    uint256 constant H1x  = 1368015179489954701390400359078579693043519447331113978918064868415326638035;
    uint256 constant H1y  = 9918110051302171585080402603319702774565515993150576347155970296011118125764;

    struct G1Point {
		uint X;
		uint Y;
	}

	// Encoding of field elements is: X[0] * z + X[1]
	struct G2Point {
		uint[2] X;
		uint[2] Y;
	}

    function P1() pure internal returns (G1Point memory) {
	    return G1Point(1, 2);
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

    function pairing(G1Point[] memory p1, G2Point[] memory p2) view internal returns (bool)
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

    /// Convenience method for a pairing check for three pairs.
	/// Convenience method for a pairing check for three pairs.
	function pairingProd3(
			G1Point memory a1, G2Point memory a2,
			G1Point memory b1, G2Point memory b2,
			G1Point memory c1, G2Point memory c2
	) internal returns (bool) {
		G1Point[] memory p1 = new G1Point[](3);
		G2Point[] memory p2 = new G2Point[](3);
		p1[0] = a1;
		p1[1] = b1;
		p1[2] = c1;
		p2[0] = a2;
		p2[1] = b2;
		p2[2] = c2;
		return pairing(p1, p2);
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

    function ZKRP_verify(
        uint256[2] memory E_j, uint256[2] memory F_j,
        uint256[2] memory U1_j, uint256[2] memory C1_j,
        uint256 c, uint256 z1,
        uint256 z2, uint256 z3,
        uint256[2][] memory V, uint256[] memory lagrange_coefficient,
        uint256[2] memory U_j
    )
    public returns (bool proof1, bool proof2)
    {
        uint256[2] memory  C_j = Interpolate(V, lagrange_coefficient);
        proof1 = ZKRP_verify1(C1_j, C_j, c, z3);
        proof2 = ZKRP_verify2(U1_j, U_j, c, z3, z1);
        //proof_is_valid = proof1 && proof2; 
    }

    function  Interpolate(
        uint256[2][] memory V, uint256[] memory lagrange_coefficient
    )
    public returns (uint256[2] memory)
    {
        uint256[2] memory a1; //acc=g
        uint256[2] memory temp;
        a1[0] = H1x;
        a1[1] = H1y;
        uint elements=lagrange_coefficient.length;//to get the array length 
        for(uint i=0;i<elements;i++)
        {
            temp = bn128_multiply([V[i][0], V[i][1],lagrange_coefficient[i]]);
            a1 = bn128_add([a1[0], a1[1], temp[0], temp[1]]);
        }
        a1=bn128_add([a1[0], a1[1], 1368015179489954701390400359078579693043519447331113978918064868415326638035, 11970132820537103637166003141937572314130795164147247315533067598634108082819]); // add neg(H1)
        return a1;
    }
    
    function ZKRP_verify1(
        uint256[2] memory C1_j, uint256[2] memory C_j,
        uint256 c, uint256 z3
    )
    public returns (bool proof_is_valid)
    {
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
        uint256 c, uint256 z3, uint256 z1
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
        uint256[2] memory F_j, uint256[2] memory E_j,
        uint256[2] memory U1_j, uint256[2] memory C1_j,
        uint256 c, uint256 z1,
        uint256 z2, uint256 z3,
        uint256[2][] memory V, uint256[] memory lagrange_coefficient,
        uint256[2] memory U_j
    )
    public returns (bool proof_is_valid)
    {
        //数据转换在这里
        //pairingProd3();
        
    }

}
