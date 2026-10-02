// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

/// @notice Atomic on-chain verifier for the RB-PVSS transcript in PVote.tex.
/// The target-group element F is encoded by its unique base-group preimage
/// FBase relative to fixed g1, i.e., F=e(FBase,g1), matching the Go verifier.
contract RBPVSSVerifier {
    uint256 internal constant Q = 21888242871839275222246405745257275088548364400416034343698204186575808495617;
    uint256 internal constant FP = 21888242871839275222246405745257275088696311157297823662689037894645226208583;

    struct G1Point { uint256 X; uint256 Y; }
    struct G2Point { uint256[2] X; uint256[2] Y; }

    struct DealerTranscript {
        G1Point[] V;       // canonical order D followed by N
        G1Point[] C;
        G1Point[] PhiA;
        G1Point[] PhiB;
        uint256 PhiChi;
        uint256[] PhiZ;
    }

    struct RangeBindingTranscript {
        G1Point[] U;
        G1Point[] E;
        G1Point[] FBase;
        G1Point[] UPrime;
        G1Point[] CPrime;
        uint256[] Chi;
        uint256[] Z1;
        uint256[] ZBeta;
        uint256[] Z3;
    }

    G1Point public g0;
    G1Point public h0;
    G2Point internal g1;
    G2Point internal pkI;
    uint256 public immutable n;
    uint256 public immutable t;
    uint256 public immutable l;
    int256 public immutable a;
    int256 public immutable b;

    G1Point[] internal publicKeys;
    G1Point[] internal rangeSignatures;
    G1Point[] internal aggregateC;
    G1Point[] internal aggregateU;
    uint256 public acceptedCount;
    mapping(address => uint256) public eligibleVoterIndex;
    mapping(uint256 => bool) public acceptedVoterIndex;

    mapping(uint256 => bool) public submittedShare;
    uint256[] internal shareIndices;
    G1Point[] internal decryptedShares;

    event TranscriptAccepted(uint256 indexed acceptedCount);
    event DecryptionShareAccepted(uint256 indexed shareholderIndex);

    constructor(
        G1Point memory g0_, G1Point memory h0_, G2Point memory g1_,
        G2Point memory pkI_, G1Point[] memory sigma_, G1Point[] memory pks_,
        uint256 t_, uint256 l_, int256 a_, int256 b_, address[] memory voterAccounts_
    ) {
        require(pks_.length > 0 && l_ > 0 && l_ < t_ && t_ <= pks_.length, "parameters");
        require(a_ <= b_ && sigma_.length == uint256(b_ - a_ + 1), "range");
        require(voterAccounts_.length > 0, "voters");
        g0 = g0_; h0 = h0_; g1 = g1_; pkI = pkI_;
        n = pks_.length; t = t_; l = l_; a = a_; b = b_;
        for (uint256 i; i < pks_.length; ++i) {
            publicKeys.push(pks_[i]);
            aggregateC.push(G1Point(0, 0));
        }
        for (uint256 d; d < l_; ++d) aggregateU.push(G1Point(0, 0));
        for (uint256 k; k < sigma_.length; ++k) {
            require(_verifyRangeSignature(sigma_[k], a_ + int256(k)), "range signature");
            rangeSignatures.push(sigma_[k]);
        }
        for (uint256 j; j < voterAccounts_.length; ++j) {
            require(voterAccounts_[j] != address(0) && eligibleVoterIndex[voterAccounts_[j]] == 0, "voter mapping");
            eligibleVoterIndex[voterAccounts_[j]] = j + 1;
        }
    }

    /// @notice Performs the complete RB-PVSS.DVerify and only then updates
    /// both homomorphic aggregates in the same transaction.
    function submitRB(DealerTranscript calldata dealer, RangeBindingTranscript calldata rangeProof) external {
        uint256 voterIndex = eligibleVoterIndex[msg.sender];
        require(voterIndex != 0 && !acceptedVoterIndex[voterIndex], "voter");
        require(_validShape(dealer, rangeProof), "shape");
        require(_verifyDealer(dealer), "PVSS dealer proof");
        require(_verifyRS(dealer.V), "dual RS check");
        for (uint256 d; d < l; ++d) {
            require(_verifyRange(dealer.V[d], rangeProof, d), "range binding");
        }
        for (uint256 i; i < n; ++i) aggregateC[i] = _add(aggregateC[i], dealer.C[i]);
        for (uint256 d; d < l; ++d) aggregateU[d] = _add(aggregateU[d], rangeProof.U[d]);
        acceptedVoterIndex[voterIndex] = true;
        ++acceptedCount;
        emit TranscriptAccepted(acceptedCount);
    }

    /// @notice Implements RB-PVSS.PVerify before recording a decryption share.
    function submitDecryptionShare(
        uint256 index, G1Point calldata share, G1Point calldata A,
        G1Point calldata B, uint256 chi, uint256 z
    ) external {
        require(index >= 1 && index <= n && !submittedShare[index], "index");
        G1Point memory pk = publicKeys[index - 1];
        G1Point memory ciphertext = aggregateC[index - 1];
        require(_equal(A, _add(_mul(g0, z), _mul(pk, chi))), "A");
        require(_equal(B, _add(_mul(share, z), _mul(ciphertext, chi))), "B");
        require(chi == _hash4(pk, ciphertext, A, B), "challenge");
        submittedShare[index] = true;
        shareIndices.push(index);
        decryptedShares.push(share);
        emit DecryptionShareAccepted(index);
    }

    /// @notice On-chain mask reconstruction/removal portion of RB-PVSS.Recon.
    /// The returned points are h0^{W_d}; bounded DLog remains a public
    /// off-chain decoding step, exactly as stated in the paper.
    function reconstruct() external view returns (G1Point[] memory unmasked) {
        require(decryptedShares.length >= t, "threshold");
        unmasked = new G1Point[](l);
        for (uint256 d; d < l; ++d) {
            uint256 target = _dPoint(d);
            G1Point memory mask = G1Point(0, 0);
            for (uint256 i; i < t; ++i) {
                uint256 numerator = 1;
                uint256 denominator = 1;
                for (uint256 j; j < t; ++j) if (i != j) {
                    numerator = mulmod(numerator, _sub(target, shareIndices[j]), Q);
                    denominator = mulmod(denominator, _sub(shareIndices[i], shareIndices[j]), Q);
                }
                uint256 mu = mulmod(numerator, _inverse(denominator), Q);
                mask = _add(mask, _mul(decryptedShares[i], mu));
            }
            unmasked[d] = _add(aggregateU[d], _negate(mask));
        }
    }

    function getAggregateC() external view returns (G1Point[] memory) { return aggregateC; }
    function getAggregateU() external view returns (G1Point[] memory) { return aggregateU; }
    function decryptionShareCount() external view returns (uint256) { return decryptedShares.length; }

    function _validShape(DealerTranscript calldata d, RangeBindingTranscript calldata r) internal view returns (bool) {
        return d.V.length == n + l && d.C.length == n && d.PhiA.length == n &&
            d.PhiB.length == n && d.PhiZ.length == n && r.U.length == l &&
            r.E.length == l && r.FBase.length == l &&
            r.UPrime.length == l && r.CPrime.length == l && r.Chi.length == l &&
            r.Z1.length == l && r.ZBeta.length == l && r.Z3.length == l;
    }

    function _verifyDealer(DealerTranscript calldata d) internal view returns (bool) {
        for (uint256 i; i < n; ++i) {
            if (!_equal(d.PhiA[i], _add(_mul(h0, d.PhiZ[i]), _mul(d.V[l + i], d.PhiChi)))) return false;
            if (!_equal(d.PhiB[i], _add(_mul(publicKeys[i], d.PhiZ[i]), _mul(d.C[i], d.PhiChi)))) return false;
        }
        bytes memory transcript;
        for (uint256 i; i < n; ++i) {
            transcript = abi.encodePacked(transcript, d.V[l+i].X, d.V[l+i].Y, d.C[i].X, d.C[i].Y,
                d.PhiA[i].X, d.PhiA[i].Y, d.PhiB[i].X, d.PhiB[i].Y);
        }
        return d.PhiChi == uint256(sha256(transcript)) % Q;
    }

    function _verifyRange(G1Point calldata vd, RangeBindingTranscript calldata r, uint256 d) internal view returns (bool) {
        if (r.E[d].X == 0 && r.E[d].Y == 0) return false;
        bytes32 digest = sha256(abi.encodePacked(
            r.E[d].X, r.E[d].Y, r.U[d].X, r.U[d].Y, vd.X, vd.Y,
            r.FBase[d].X, r.FBase[d].Y,
            r.UPrime[d].X, r.UPrime[d].Y, r.CPrime[d].X, r.CPrime[d].Y
        ));
        if (r.Chi[d] != uint256(digest) % Q) return false;
        if (!_equal(r.CPrime[d], _add(_mul(vd, r.Chi[d]), _mul(h0, r.Z3[d])))) return false;
        if (!_equal(r.UPrime[d], _add(_add(_mul(r.U[d], r.Chi[d]), _mul(g0, r.Z3[d])), _mul(h0, r.Z1[d])))) return false;
        G1Point[] memory p1 = new G1Point[](4);
        G2Point[] memory p2 = new G2Point[](4);
        p1[0] = _negate(r.FBase[d]); p2[0] = g1;
        p1[1] = _mul(r.E[d], r.Chi[d]); p2[1] = pkI;
        p1[2] = _negate(_mul(r.E[d], r.Z1[d])); p2[2] = g1;
        p1[3] = _mul(g0, r.ZBeta[d]); p2[3] = g1;
        return _pairing(p1, p2);
    }

    /// A Fiat-Shamir-derived random dual-RS word makes the public check
    /// deterministic on chain while retaining the paper's random-codeword check.
    function _verifyRS(G1Point[] calldata V) internal view returns (bool) {
        bytes memory encoded;
        for (uint256 i; i < V.length; ++i) encoded = abi.encodePacked(encoded, V[i].X, V[i].Y);
        bytes32 seed = sha256(encoded);
        uint256 degreeBound = n + l - t;
        uint256[] memory qCoeff = new uint256[](degreeBound);
        for (uint256 k; k < degreeBound; ++k) qCoeff[k] = uint256(sha256(abi.encodePacked(seed, k))) % Q;
        G1Point memory sum = G1Point(0, 0);
        for (uint256 i; i < n + l; ++i) {
            uint256 x = i < l ? _dPoint(i) : i - l + 1;
            uint256 denominator = 1;
            for (uint256 j; j < n + l; ++j) if (i != j) {
                uint256 other = j < l ? _dPoint(j) : j - l + 1;
                denominator = mulmod(denominator, _sub(x, other), Q);
            }
            uint256 qx;
            for (uint256 k = degreeBound; k > 0; --k) qx = addmod(mulmod(qx, x, Q), qCoeff[k-1], Q);
            uint256 y = mulmod(qx, _inverse(denominator), Q);
            sum = _add(sum, _mul(V[i], y));
        }
        return sum.X == 0 && sum.Y == 0;
    }

    function _verifyRangeSignature(G1Point memory sigma, int256 w) internal view returns (bool) {
        G1Point[] memory p1 = new G1Point[](3);
        G2Point[] memory p2 = new G2Point[](3);
        p1[0] = sigma; p2[0] = pkI;
        p1[1] = _mul(sigma, _iota(w)); p2[1] = g1;
        p1[2] = _negate(g0); p2[2] = g1;
        return _pairing(p1, p2);
    }

    function _dPoint(uint256 offset) internal view returns (uint256) {
        int256 value = int256(offset) - int256(l) + 1;
        return _iota(value);
    }

    function _iota(int256 value) internal pure returns (uint256) {
        if (value >= 0) return uint256(value) % Q;
        uint256 magnitude = uint256(-value) % Q;
        return magnitude == 0 ? 0 : Q - magnitude;
    }

    function _hash4(G1Point memory x1, G1Point memory x2, G1Point memory x3, G1Point memory x4) internal pure returns (uint256) {
        return uint256(sha256(abi.encodePacked(x1.X,x1.Y,x2.X,x2.Y,x3.X,x3.Y,x4.X,x4.Y))) % Q;
    }

    function _equal(G1Point memory x, G1Point memory y) internal pure returns (bool) { return x.X == y.X && x.Y == y.Y; }
    function _sub(uint256 x, uint256 y) internal pure returns (uint256) { return addmod(x, Q - (y % Q), Q); }
    function _negate(G1Point memory p) internal pure returns (G1Point memory) {
        if (p.X == 0 && p.Y == 0) return p;
        return G1Point(p.X, FP - (p.Y % FP));
    }

    function _add(G1Point memory x, G1Point memory y) internal view returns (G1Point memory r) {
        uint256[4] memory input = [x.X, x.Y, y.X, y.Y];
        bool ok;
        assembly ("memory-safe") { ok := staticcall(gas(), 6, input, 0x80, r, 0x40) }
        require(ok, "ecadd");
    }

    function _mul(G1Point memory p, uint256 scalar) internal view returns (G1Point memory r) {
        uint256[3] memory input = [p.X, p.Y, scalar];
        bool ok;
        assembly ("memory-safe") { ok := staticcall(gas(), 7, input, 0x60, r, 0x40) }
        require(ok, "ecmul");
    }

    function _pairing(G1Point[] memory p1, G2Point[] memory p2) internal view returns (bool) {
        require(p1.length == p2.length, "pairing length");
        uint256[] memory input = new uint256[](p1.length * 6);
        for (uint256 i; i < p1.length; ++i) {
            uint256 j = i * 6;
            input[j] = p1[i].X; input[j+1] = p1[i].Y;
            input[j+2] = p2[i].X[0]; input[j+3] = p2[i].X[1];
            input[j+4] = p2[i].Y[0]; input[j+5] = p2[i].Y[1];
        }
        uint256[1] memory output;
        bool ok;
        assembly ("memory-safe") { ok := staticcall(gas(), 8, add(input, 0x20), mul(mload(input), 0x20), output, 0x20) }
        require(ok, "pairing");
        return output[0] == 1;
    }

    function _inverse(uint256 value) internal view returns (uint256 result) {
        uint256[6] memory input = [uint256(32),32,32,value,Q-2,Q];
        uint256[1] memory output;
        bool ok;
        assembly ("memory-safe") { ok := staticcall(gas(), 5, input, 0xc0, output, 0x20) }
        require(ok, "inverse");
        return output[0];
    }
}
