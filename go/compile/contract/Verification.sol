// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

// https://www.iacr.org/cryptodb/archive/2002/ASIACRYPT/50/50.pdf
contract Verification
{
    uint256 constant FIELD_ORDER = 0x30644e72e131a029b85045b68181585d97816a916871ca8d3c208c16d87cfd47;
    uint256 constant CURVE_ORDER = 0x30644e72e131a029b85045b68181585d2833e84879b9709143e1f593f0000001;
    uint256 constant GROUP_ORDER   = 21888242871839275222246405745257275088548364400416034343698204186575808495617;
    uint256 constant CURVE_B = 3;

    // a = (p+1) / 4
    uint256 constant CURVE_A = 0xc19139cb84c680a6e14116da060561765e05aa45a1c72a34f082305b61f3f52;

    struct G1Point {
        uint X;
        uint Y;
    }

    // Encoding of field elements is: X[0] * z + X[1]
    struct G2Point {
        uint[2] X;
        uint[2] Y;
    }

    // (P+1) / 4
    function A() pure internal returns (uint256) {
        return CURVE_A;
    }

    function P() pure internal returns (uint256) {
        return FIELD_ORDER;
    }

    function N() pure internal returns (uint256) {
        return CURVE_ORDER;
    }

    /// return the generator of G1
    function P1() pure internal returns (G1Point memory) {
        return G1Point(1, 2);
    }

	/// return the generator of G2
	function P2() pure internal returns (G2Point memory) {
		return G2Point(
			[11559732032986387107991004021392285783925812861821192530917403151452391805634,
			 10857046999023057135944570762232829481370756359578518086990519993285655852781],
			[4082367875863433681332203403145435568316851327593401208105741076214120093531,
			 8495653923123431417604973247489272438418190587263600148770280649306958101930]
		);
	}
    /// return the negation of p, i.e. p.add(p.negate()) should be zero.
	function g1neg(G1Point memory p) pure internal returns (G1Point memory) {
		// The prime q in the base field F_q for G1
		uint q = 21888242871839275222246405745257275088696311157297823662689037894645226208583;
		if (p.X == 0 && p.Y == 0)
			return G1Point(0, 0);
		return G1Point(p.X, q - (p.Y % q));
	}

    /// return the sum of two points of G1
    function g1add(G1Point memory p1, G1Point memory p2) view internal returns (G1Point memory r) {
        uint[4] memory input;
        input[0] = p1.X;
        input[1] = p1.Y;
        input[2] = p2.X;
        input[3] = p2.Y;
        bool success;
        assembly {
            success := staticcall(sub(gas(), 2000), 6, input, 0xc0, r, 0x60)
        // Use "invalid" to make gas estimation work
        //switch success case 0 { invalid }
        }
        require(success);
    }

    /// return the product of a point on G1 and a scalar, i.e.
    /// p == p.mul(1) and p.add(p) == p.mul(2) for all points p.
    function g1mul(G1Point memory p, uint s) view internal returns (G1Point memory r) {
        uint[3] memory input;
        input[0] = p.X;
        input[1] = p.Y;
        input[2] = s;
        bool success;
        assembly {
            success := staticcall(sub(gas(), 2000), 7, input, 0x80, r, 0x60)
        // Use "invalid" to make gas estimation work
        //switch success case 0 { invalid }
        }
        require(success);
    }



	function pairing(G1Point[] memory p1, G2Point[] memory p2) view internal returns (bool) {
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


	/// Convenience method for a pairing check for two pairs.
	function pairingProd2(G1Point memory a1, G2Point memory a2, G1Point memory b1, G2Point memory b2) view internal returns (bool) {
		G1Point[] memory p1 = new G1Point[](2);
		G2Point[] memory p2 = new G2Point[](2);
		p1[0] = a1;
		p1[1] = b1;
		p2[0] = a2;
		p2[1] = b2;
		return pairing(p1, p2);
	}

	/// Convenience method for a pairing check for three pairs.
	function pairingProd3(
			G1Point memory a1, G2Point memory a2,
			G1Point memory b1, G2Point memory b2,
			G1Point memory c1, G2Point memory c2
	) view internal returns (bool) {
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

	/// Convenience method for a pairing check for four pairs.
	function pairingProd4(
			G1Point memory a1, G2Point memory a2,
			G1Point memory b1, G2Point memory b2,
			G1Point memory c1, G2Point memory c2,
			G1Point memory d1, G2Point memory d2
	) view internal returns (bool) {
		G1Point[] memory p1 = new G1Point[](4);
		G2Point[] memory p2 = new G2Point[](4);
		p1[0] = a1;
		p1[1] = b1;
		p1[2] = c1;
		p1[3] = d1;
		p2[0] = a2;
		p2[1] = b2;
		p2[2] = c2;
		p2[3] = d2;
		return pairing(p1, p2);
	}

    function pairingProd5(
			G1Point memory a1, G2Point memory a2,
			G1Point memory b1, G2Point memory b2,
			G1Point memory c1, G2Point memory c2,
			G1Point memory d1, G2Point memory d2,
            G1Point memory e1, G2Point memory e2
	) view internal returns (bool) {
		G1Point[] memory p1 = new G1Point[](5);
		G2Point[] memory p2 = new G2Point[](5);
		p1[0] = a1;
		p1[1] = b1;
		p1[2] = c1;
		p1[3] = d1;
        p1[4] = e1;
		p2[0] = a2;
		p2[1] = b2;
		p2[2] = c2;
		p2[3] = d2;
        p2[4] = e2;
		return pairing(p1, p2);
	}

    function submod(uint a, uint b) internal pure returns (uint){
        uint a_nn;

        if(a>b) {
            a_nn = a;
        } else {
            a_nn = a+CURVE_ORDER;
        }

        return addmod(a_nn - b, 0, CURVE_ORDER);
    }

    // Invert function, invert in group
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

    struct Parameters{
        G1Point G0;
        G1Point H0;
        G2Point G1;
        G2Point PKI;
        G1Point[] SigmaK;
        uint256 a;
        uint256 b;
        uint256 numCandidates;
    }

    struct DLEQProof {
        G1Point a1;
        G1Point a2;
        uint256 c;
        uint256 z;
    }

    struct PVSS{
        G1Point[] C;
        G1Point[] V;
        DLEQProof[] Proof;
    }


    Parameters PP;
    G1Point[] PKs;
    G1Point[][] EncShares;
    G1Point[] DecShares;
    G1Point[][] U;
    G1Point[] AggregatedC;
    G1Point[] AggregatedU;
    uint256[] HonestTalliers;
    uint256[] TallyValue;
    bool[] ZKRPVerifyResult;
    bool[] DVerifyResult;
    bool[] PVerifyResult;
    

    function UploadParameters(G1Point memory g0, G1Point memory h0, G2Point memory g1, G2Point memory pkI, G1Point[] memory sigmak,uint256 a, uint256 b,uint256 numCandidates) public {
        PP.G0 = g0;
        PP.H0 = h0;
        PP.G1 = g1;
        PP.PKI = pkI;
        PP.a=a;
        PP.b=b;
        PP.numCandidates=numCandidates;
    	for (uint i = 0; i < sigmak.length; i++) {
            PP.SigmaK.push(sigmak[i]);
        }
    }

    function UploadPublicKey(G1Point[] memory pks) public{
        for (uint i=0;i<pks.length;i++){
            PKs.push(pks[i]);
        }
    }

    function DLEQVerify(G1Point memory g, G1Point memory y1, G1Point memory a1, G1Point memory h, 
                          G1Point memory y2, G1Point memory a2, uint256 c, uint256 z) public payable returns (bool)
    {
        G1Point memory gG = g1mul(g, z);
        G1Point memory y1G = g1mul(y1, c);

        G1Point memory hG = g1mul(h, z);
        G1Point memory y2G = g1mul(y2, c);

        G1Point memory pt1 =  g1add(gG, y1G);
        G1Point memory pt2 =  g1add(hG, y2G);
        if ((a1.X != pt1.X) || (a1.Y != pt1.Y) || (a2.X != pt2.X) || (a2.Y != pt2.Y))
        {
            return false;
        }
        return true;
    }

    function coefficient(uint256 i,uint256 n) public returns (uint256) {
        uint256 result=1;
        for (uint256 j=1;j<n+1;j++){
            if (i!=j){
                result=mulmod(result,inv(submod(i,j),GROUP_ORDER),GROUP_ORDER);
                result=addmod(result,0, GROUP_ORDER);
            }
        }
        return result;
    }

    function RScodeVerify(G1Point[] memory v) public returns (bool)
    {
        uint256[] memory codeword;
        codeword = new uint256[](v.length);
        for (uint i=1;i<=v.length;i++){
            codeword[i-1]= addmod(coefficient(i,v.length),0,GROUP_ORDER);
        }
        G1Point memory sum=g1mul(PP.H0,1);
        for (uint i=0;i<v.length;i++){
            sum = g1add(sum,g1mul(v[i],codeword[i]));
        }

        if ((sum.X != PP.H0.X) || (sum.Y != PP.H0.Y))
        {
            return false;
        }
        return true;
    }

    function DVerify(G1Point memory g, G1Point[] memory y1, G1Point[] memory a1, G1Point[] memory h, G1Point[] memory y2, G1Point[] memory a2, uint256[] memory c, uint256[] memory z) public returns (bool){
        if (RScodeVerify(y1)== false){
            return false;
        }else{
            for(uint i=0;i<PKs.length;i++)
            {
                if (DLEQVerify(g, y1[i], a1[i], h[i],  y2[i], a2[i], c[i], z[i]) ==false)
                {
                    return false;
                } 
            }
        }
        return true;
    }


    function GetDVerifyResult() public view returns (bool[] memory){
        return DVerifyResult;
    }

    function UploadPVSSShares(G1Point[] memory v, G1Point[] memory c, G1Point[] memory a1,
                              G1Point[] memory a2, uint256[] memory challenge, uint256[] memory z)public{
        for (uint i = 0; i < v.length; i++) {
            if (DVerify(PP.H0, v, a1, PKs, c, a2, challenge, z)==false){
                DVerifyResult.push(false);
                return;
            }
        }
        DVerifyResult.push(true);
        EncShares.push(c);
        return;
    }

    function UploadBallotCipher(G1Point[] memory Ej, G1Point[] memory Fj1, G1Point[] memory Fj2,G1Point[] memory _Uj,
                        G1Point[] memory _Cj,uint256[] memory c,uint256[] memory z1,uint256[] memory z2,
                        uint256[] memory z3, G1Point[] memory Uj, G1Point[] memory v, uint threshold)public{
        uint256[] memory selectIndices;
        selectIndices=new uint256[](threshold);
        for (uint i=0;i<threshold;i++){
            selectIndices[i]=i+1;
        }
        uint256[] memory x;
        x=new uint256[](PP.numCandidates);
        for (uint d=0;d<PP.numCandidates;d++){
            x[d]=submod(0,d);
        }
        if (ZKRPVerify(Ej, Fj1, Fj2, _Uj, _Cj, c, z1, z2, z3, Uj, x, v, selectIndices, threshold) == true){
            U.push(Uj);
        }
    }

    function Interpolation(uint256 d, G1Point[] memory v, uint256[] memory indices, uint256 threshold) public returns(G1Point memory){
        uint256[] memory coefficients;
        coefficients = new uint256[](threshold);
        for (uint i=0;i<threshold;i++){
            uint256 num=1;
            uint256 den=1;

            for (uint j=0;j<threshold;j++){
                if (i!=j){
                    num=mulmod(num,submod(d,indices[j]),GROUP_ORDER);
                    den=mulmod(den,submod(indices[i],indices[j]),GROUP_ORDER);
                }
            }

            den=inv(den,GROUP_ORDER);
            uint256 term=mulmod(num,den,GROUP_ORDER);
            coefficients[i]=term;
        }
        G1Point memory secret=g1mul(PP.G0,0);
        for(uint i=0;i<coefficients.length;i++){
            secret=g1add(secret,g1mul(v[i],coefficients[i]));
        }
        return secret;
    }

    function ZKRPVerify(G1Point[] memory Ej, G1Point[] memory Fj1, G1Point[] memory Fj2,G1Point[] memory _Uj,
                        G1Point[] memory _Cj,uint256[] memory c,uint256[] memory z1,uint256[] memory z2,
                        uint256[] memory z3, G1Point[] memory Uj, uint256[] memory d, G1Point[] memory v, 
                        uint256[] memory indices, uint threshold)public returns (bool){
        for(uint i=0;i<d.length;i++){
            G1Point memory temp=Interpolation(d[i],v,indices,threshold);
            temp= g1mul(temp, c[i]);
            temp=g1add(temp,g1mul(PP.H0,z3[i]));
            if ((_Cj[i].X != temp.X) || (_Cj[i].Y != temp.Y)){
                ZKRPVerifyResult.push(false);
                return false;
            }else{
                G1Point memory temp1=g1mul(Uj[i],c[i]);
                temp1=g1add(temp1,g1mul(PP.G0,z3[i]));
                temp1=g1add(temp1,g1mul(PP.H0,z1[i]));
                if((_Uj[i].X != temp1.X) || (_Uj[i].Y != temp1.Y)){
                    ZKRPVerifyResult.push(false);
                    return false;
                }else{
                    G1Point memory right1=g1mul(Ej[i],c[i]);
                    G1Point memory right2=g1neg(g1mul(Ej[i],z1[i]));
                    G1Point memory right3=g1mul(PP.G0,z2[i]);
                    if (pairingProd5(g1neg(Fj1[i]),PP.G1,g1neg(Fj2[i]),PP.G1,right1,PP.PKI,right2,PP.G1,right3,PP.G1)==false){
                        ZKRPVerifyResult.push(false);
                        return false;
                    }
                }
            } 
        }
        ZKRPVerifyResult.push(true);
        return true;
    }
  
    function GetZKRPResult() public view returns (bool[] memory){
        return ZKRPVerifyResult;
    }

    function Aggregate() public{
        for (uint256 i=0;i<EncShares[0].length;i++){
            G1Point memory aggregateC=g1mul(PP.G0,0);
            for (uint256 j=0;j<EncShares.length;j++){
                aggregateC=g1add(aggregateC,EncShares[j][i]);
            }
            AggregatedC.push(aggregateC);
        }

        for(uint256 i=0;i<U[0].length;i++){
            G1Point memory aggregateU=g1mul(PP.G0,0);
            for(uint256 j=0;j<U.length;j++){
                aggregateU=g1add(aggregateU,U[j][i]);
            }
            AggregatedU.push(aggregateU);
        }
    }

    function GetAggregateValue() public view returns (G1Point[] memory){
        return AggregatedC;
    }    

    function PVerify(uint indexTallier, G1Point memory DecShare, G1Point memory a1,G1Point memory a2,
                     uint256 challenge, uint256 z)public{
        if (DLEQVerify(PP.G0, PKs[indexTallier], a1, DecShare, AggregatedC[indexTallier], a2, challenge, z)==true){
            DecShares.push(DecShare);
            HonestTalliers.push(indexTallier);
            PVerifyResult.push(true);
        }else{
            PVerifyResult.push(false);
        }

    }

    function GetPVerifyResult() public view returns (bool[] memory){
        return PVerifyResult;
    }

    function Tally(uint threshold,uint numCandidates, uint a, uint b)public{
        uint256[] memory selectIndices;
        G1Point[] memory selectShares;
        selectIndices=new uint256[](threshold);
        selectShares=new G1Point[](threshold);
        //Generate the index set of honest talliers
        for (uint i=0;i<threshold;i++){
            selectIndices[i]=HonestTalliers[i]+1;
            selectShares[i]=DecShares[i];
        }
        G1Point[] memory S;
        S=new G1Point[](numCandidates);
        for (uint d=0;d<numCandidates;d++){
            uint x=submod(0,d);
            S[d]=Interpolation(x, selectShares, selectIndices,threshold);
            AggregatedU[d]=g1add(AggregatedU[d],g1neg(S[d]));
        }

        for (uint d=0;d<numCandidates;d++){
            for (uint k=(d+1)*a;k<b*numCandidates;k++){
                if (AggregatedU[d].X==g1mul(PP.H0,k).X&&AggregatedU[d].Y==g1mul(PP.H0,k).Y){
                    TallyValue.push(k);
                    break;
                }
            }
        }
    }

    function GetTallyValue() public view returns (uint256[] memory){
        return TallyValue;
    }

}
