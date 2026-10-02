// Package threepvss implements the individual-commitment 3PVSS binary-voting
// construction and the explicit k-instance score adapter used by the paper.
package threepvss

import (
	"errors"
	"math/big"

	"PVote/baselines/internal/ec"
	bn256 "github.com/ethereum/go-ethereum/crypto/bn256/cloudflare"
)

var ErrInvalid = errors.New("threepvss: invalid input")

type Parameters struct {
	N, F, L, R int
	G, H       *bn256.G1
	PK         []*bn256.G1
	SK         []*big.Int
}

type PDLProof struct {
	Challenge *big.Int
	Z         []*big.Int
}

// Transcript is the broadcast output of Lambda_RO^indv.Share.
type Transcript struct {
	Commitments []*bn256.G1
	Encrypted   []*bn256.G1
	PDL         PDLProof
}

type BinaryProof struct {
	A         [2]*bn256.G1
	B         [2]*bn256.G1
	Challenge [2]*big.Int
	Response  [2]*big.Int
}

type BinaryBallot struct {
	Commitments []*bn256.G1
	Encrypted   []*bn256.G1
	PDL         PDLProof
	U           []*bn256.G1
	Binary      []BinaryProof
}

type ScoreBallot struct {
	Kappa  int
	Planes []*BinaryBallot
}

type DecryptionProof struct {
	A, B *bn256.G1
	C, Z *big.Int
}

// Share implements the native pre-constructed packed-PVSS sharing interface.
func Share(pp *Parameters, secrets []*big.Int) (*Transcript, error) {
	if !validPP(pp) || len(secrets) != pp.L {
		return nil, ErrInvalid
	}
	f, err := constrainedPolynomial(pp, secrets)
	if err != nil {
		return nil, err
	}
	tr := &Transcript{Commitments: make([]*bn256.G1, pp.L), Encrypted: make([]*bn256.G1, pp.N)}
	for d := 0; d < pp.L; d++ {
		tr.Commitments[d] = ec.Mul(pp.H, eval(f, packedX(pp.L, d)))
	}
	for i := 0; i < pp.N; i++ {
		tr.Encrypted[i] = ec.Mul(pp.PK[i], eval(f, big.NewInt(int64(i+1))))
	}
	tr.PDL, err = provePDL(pp, f, tr.Commitments, tr.Encrypted)
	if err != nil {
		return nil, err
	}
	return tr, nil
}

// Verify implements Lambda_RO^indv.Verify.
func Verify(pp *Parameters, tr *Transcript) bool {
	return validPP(pp) && tr != nil && len(tr.Commitments) == pp.L && len(tr.Encrypted) == pp.N && verifyPDL(pp, tr.Commitments, tr.Encrypted, tr.PDL)
}

// OptimisticReconstruct verifies dealer-revealed scalars against the individual
// commitments and returns the corresponding encodings under G.
func OptimisticReconstruct(pp *Parameters, tr *Transcript, secrets []*big.Int) ([]*bn256.G1, error) {
	if !Verify(pp, tr) || len(secrets) != pp.L {
		return nil, ErrInvalid
	}
	out := make([]*bn256.G1, pp.L)
	for d, secret := range secrets {
		if secret == nil || !ec.Equal(tr.Commitments[d], ec.Mul(pp.H, secret)) {
			return nil, ErrInvalid
		}
		out[d] = ec.Mul(pp.G, secret)
	}
	return out, nil
}

// PessimisticReconstruct decrypts a qualified set of r=f+l shares, creates and
// verifies each DLEQ proof, and interpolates all packed secret encodings.
func PessimisticReconstruct(pp *Parameters, tr *Transcript, indices []int) ([]*bn256.G1, error) {
	if !Verify(pp, tr) || len(indices) < pp.R {
		return nil, ErrInvalid
	}
	indices = indices[:pp.R]
	shares := make([]*bn256.G1, pp.R)
	seen := make(map[int]bool, pp.R)
	for k, index := range indices {
		if index < 1 || index > pp.N || seen[index] {
			return nil, ErrInvalid
		}
		seen[index] = true
		i := index - 1
		shares[k] = ec.Mul(tr.Encrypted[i], ec.Inv(pp.SK[i]))
		proof, err := proveDecrypt(pp, i, tr.Encrypted[i], shares[k])
		if err != nil || !verifyDecrypt(pp, i, tr.Encrypted[i], shares[k], proof) {
			return nil, ErrInvalid
		}
	}
	out := make([]*bn256.G1, pp.L)
	for d := 0; d < pp.L; d++ {
		out[d] = ec.Zero()
		x := packedX(pp.L, d)
		for k, index := range indices {
			out[d] = ec.Add(out[d], ec.Mul(shares[k], lagrangeForIndices(index, indices, x)))
		}
	}
	return out, nil
}

func Setup(n, f, l int) (*Parameters, error) {
	if n < 1 || f < 0 || l < 1 || n < 2*f+l {
		return nil, ErrInvalid
	}
	h, err := ec.Scalar()
	if err != nil {
		return nil, err
	}
	pp := &Parameters{N: n, F: f, L: l, R: f + l, G: new(bn256.G1).ScalarBaseMult(big.NewInt(1)), H: new(bn256.G1).ScalarBaseMult(h), PK: make([]*bn256.G1, n), SK: make([]*big.Int, n)}
	for i := 0; i < n; i++ {
		pp.SK[i], err = ec.Scalar()
		if err != nil {
			return nil, err
		}
		pp.PK[i] = ec.Mul(pp.G, pp.SK[i])
	}
	return pp, nil
}

func ShareScore(pp *Parameters, scores []int, kappa int) (*ScoreBallot, error) {
	if !validPP(pp) || len(scores) != pp.L || kappa < 1 || kappa > 30 {
		return nil, ErrInvalid
	}
	limit := 1 << kappa
	for _, s := range scores {
		if s < 0 || s >= limit {
			return nil, ErrInvalid
		}
	}
	out := &ScoreBallot{Kappa: kappa, Planes: make([]*BinaryBallot, kappa)}
	for b := 0; b < kappa; b++ {
		bits := make([]int, pp.L)
		for d := range bits {
			bits[d] = (scores[d] >> b) & 1
		}
		var err error
		out.Planes[b], err = ShareBinary(pp, bits)
		if err != nil {
			return nil, err
		}
	}
	return out, nil
}

func VerifyScore(pp *Parameters, ballot *ScoreBallot) bool {
	if ballot == nil || ballot.Kappa < 1 || len(ballot.Planes) != ballot.Kappa {
		return false
	}
	for _, plane := range ballot.Planes {
		if !VerifyBinary(pp, plane) {
			return false
		}
	}
	return true
}

func ShareBinary(pp *Parameters, bits []int) (*BinaryBallot, error) {
	if !validPP(pp) || len(bits) != pp.L {
		return nil, ErrInvalid
	}
	masks := make([]*big.Int, pp.L)
	for i, bit := range bits {
		if bit != 0 && bit != 1 {
			return nil, ErrInvalid
		}
		var err error
		masks[i], err = ec.Scalar()
		if err != nil {
			return nil, err
		}
	}
	f, err := constrainedPolynomial(pp, masks)
	if err != nil {
		return nil, err
	}
	ballot := &BinaryBallot{Commitments: make([]*bn256.G1, pp.L), Encrypted: make([]*bn256.G1, pp.N), U: make([]*bn256.G1, pp.L), Binary: make([]BinaryProof, pp.L)}
	for d := 0; d < pp.L; d++ {
		x := packedX(pp.L, d)
		s := eval(f, x)
		ballot.Commitments[d] = ec.Mul(pp.H, s)
		ballot.U[d] = ec.Mul(pp.G, ec.AddScalar(s, big.NewInt(int64(bits[d]))))
		ballot.Binary[d], err = proveBinary(pp, ballot.Commitments[d], ballot.U[d], s, bits[d])
		if err != nil {
			return nil, err
		}
	}
	for i := 0; i < pp.N; i++ {
		ballot.Encrypted[i] = ec.Mul(pp.PK[i], eval(f, big.NewInt(int64(i+1))))
	}
	ballot.PDL, err = provePDL(pp, f, ballot.Commitments, ballot.Encrypted)
	if err != nil {
		return nil, err
	}
	return ballot, nil
}

func VerifyBinary(pp *Parameters, ballot *BinaryBallot) bool {
	if !validPP(pp) || ballot == nil || len(ballot.Commitments) != pp.L || len(ballot.Encrypted) != pp.N || len(ballot.U) != pp.L || len(ballot.Binary) != pp.L || len(ballot.PDL.Z) != pp.R {
		return false
	}
	if !verifyPDL(pp, ballot.Commitments, ballot.Encrypted, ballot.PDL) {
		return false
	}
	for d := 0; d < pp.L; d++ {
		if !verifyBinary(pp, ballot.Commitments[d], ballot.U[d], ballot.Binary[d]) {
			return false
		}
	}
	return true
}

// PublicBallotBytes returns the paper's fixed-width communication accounting:
// (n+5l)|G|+(f+5l+1)|Zq| per binary plane on BN254.
func PublicBallotBytes(pp *Parameters, kappa int) int {
	if !validPP(pp) || kappa < 1 {
		return 0
	}
	return kappa * ((pp.N+5*pp.L)*64 + (pp.F+5*pp.L+1)*32)
}

// Tally reconstructs each aggregate bit plane using r=f+l valid talliers.
func Tally(pp *Parameters, ballots []*ScoreBallot) ([]int, error) {
	if !validPP(pp) || len(ballots) == 0 {
		return nil, ErrInvalid
	}
	kappa := ballots[0].Kappa
	for _, b := range ballots {
		if b == nil || b.Kappa != kappa || !VerifyScore(pp, b) {
			return nil, ErrInvalid
		}
	}
	totals := make([]int, pp.L)
	for plane := 0; plane < kappa; plane++ {
		enc := make([]*bn256.G1, pp.N)
		u := make([]*bn256.G1, pp.L)
		for i := range enc {
			enc[i] = ec.Zero()
		}
		for d := range u {
			u[d] = ec.Zero()
		}
		for _, ballot := range ballots {
			for i := range enc {
				enc[i] = ec.Add(enc[i], ballot.Planes[plane].Encrypted[i])
			}
			for d := range u {
				u[d] = ec.Add(u[d], ballot.Planes[plane].U[d])
			}
		}
		shares := make([]*bn256.G1, pp.R)
		for i := 0; i < pp.R; i++ {
			shares[i] = ec.Mul(enc[i], ec.Inv(pp.SK[i]))
			proof, err := proveDecrypt(pp, i, enc[i], shares[i])
			if err != nil || !verifyDecrypt(pp, i, enc[i], shares[i], proof) {
				return nil, ErrInvalid
			}
		}
		for d := 0; d < pp.L; d++ {
			mask := ec.Zero()
			x := packedX(pp.L, d)
			for i := 0; i < pp.R; i++ {
				mask = ec.Add(mask, ec.Mul(shares[i], lagrangeAt(i+1, pp.R, x)))
			}
			count, err := ec.DLog(pp.G, ec.Sub(u[d], mask), len(ballots))
			if err != nil {
				return nil, err
			}
			totals[d] += count << plane
		}
	}
	return totals, nil
}

func provePDL(pp *Parameters, f []*big.Int, commitments, encrypted []*bn256.G1) (PDLProof, error) {
	r := make([]*big.Int, pp.R)
	for i := range r {
		var err error
		r[i], err = ec.Scalar()
		if err != nil {
			return PDLProof{}, err
		}
	}
	y := append(append([]*bn256.G1{}, commitments...), encrypted...)
	c := make([]*bn256.G1, 0, pp.L+pp.N)
	for d := 0; d < pp.L; d++ {
		c = append(c, ec.Mul(pp.H, eval(r, packedX(pp.L, d))))
	}
	for i := 0; i < pp.N; i++ {
		c = append(c, ec.Mul(pp.PK[i], eval(r, big.NewInt(int64(i+1)))))
	}
	d := ec.HashScalar("threepvss/pdl", ec.PointBytes(y...), ec.PointBytes(c...))
	z := make([]*big.Int, pp.R)
	for i := range z {
		z[i] = ec.AddScalar(r[i], ec.MulScalar(d, f[i]))
	}
	return PDLProof{Challenge: d, Z: z}, nil
}

func verifyPDL(pp *Parameters, commitments, encrypted []*bn256.G1, proof PDLProof) bool {
	if proof.Challenge == nil || len(proof.Z) != pp.R {
		return false
	}
	y := append(append([]*bn256.G1{}, commitments...), encrypted...)
	c := make([]*bn256.G1, 0, pp.L+pp.N)
	for d := 0; d < pp.L; d++ {
		c = append(c, ec.Sub(ec.Mul(pp.H, eval(proof.Z, packedX(pp.L, d))), ec.Mul(commitments[d], proof.Challenge)))
	}
	for i := 0; i < pp.N; i++ {
		c = append(c, ec.Sub(ec.Mul(pp.PK[i], eval(proof.Z, big.NewInt(int64(i+1)))), ec.Mul(encrypted[i], proof.Challenge)))
	}
	return proof.Challenge.Cmp(ec.HashScalar("threepvss/pdl", ec.PointBytes(y...), ec.PointBytes(c...))) == 0
}

func proveBinary(pp *Parameters, commitment, u *bn256.G1, secret *big.Int, bit int) (BinaryProof, error) {
	var p BinaryProof
	other := 1 - bit
	w, err := ec.Scalar()
	if err != nil {
		return p, err
	}
	p.Challenge[other], err = ec.Scalar()
	if err != nil {
		return p, err
	}
	p.Response[other], err = ec.Scalar()
	if err != nil {
		return p, err
	}
	p.A[bit], p.B[bit] = ec.Mul(pp.H, w), ec.Mul(pp.G, w)
	p.A[other] = ec.Add(ec.Mul(pp.H, p.Response[other]), ec.Mul(commitment, p.Challenge[other]))
	target := ec.Sub(u, ec.Mul(pp.G, big.NewInt(int64(other))))
	p.B[other] = ec.Add(ec.Mul(pp.G, p.Response[other]), ec.Mul(target, p.Challenge[other]))
	c := ec.HashScalar("threepvss/binary", ec.PointBytes(commitment, u, p.A[0], p.B[0], p.A[1], p.B[1]))
	p.Challenge[bit] = ec.SubScalar(c, p.Challenge[other])
	p.Response[bit] = ec.SubScalar(w, ec.MulScalar(secret, p.Challenge[bit]))
	return p, nil
}

func verifyBinary(pp *Parameters, commitment, u *bn256.G1, p BinaryProof) bool {
	for b := 0; b < 2; b++ {
		if p.A[b] == nil || p.B[b] == nil || p.Challenge[b] == nil || p.Response[b] == nil {
			return false
		}
		if !ec.Equal(p.A[b], ec.Add(ec.Mul(pp.H, p.Response[b]), ec.Mul(commitment, p.Challenge[b]))) {
			return false
		}
		target := ec.Sub(u, ec.Mul(pp.G, big.NewInt(int64(b))))
		if !ec.Equal(p.B[b], ec.Add(ec.Mul(pp.G, p.Response[b]), ec.Mul(target, p.Challenge[b]))) {
			return false
		}
	}
	c := ec.HashScalar("threepvss/binary", ec.PointBytes(commitment, u, p.A[0], p.B[0], p.A[1], p.B[1]))
	return ec.AddScalar(p.Challenge[0], p.Challenge[1]).Cmp(c) == 0
}

func proveDecrypt(pp *Parameters, i int, encrypted, share *bn256.G1) (DecryptionProof, error) {
	w, err := ec.Scalar()
	if err != nil {
		return DecryptionProof{}, err
	}
	a, b := ec.Mul(pp.G, w), ec.Mul(share, w)
	c := ec.HashScalar("threepvss/decrypt", ec.PointBytes(pp.PK[i], encrypted, share, a, b))
	return DecryptionProof{A: a, B: b, C: c, Z: ec.AddScalar(w, ec.MulScalar(c, pp.SK[i]))}, nil
}

func verifyDecrypt(pp *Parameters, i int, encrypted, share *bn256.G1, p DecryptionProof) bool {
	if p.A == nil || p.B == nil || p.C == nil || p.Z == nil {
		return false
	}
	if p.C.Cmp(ec.HashScalar("threepvss/decrypt", ec.PointBytes(pp.PK[i], encrypted, share, p.A, p.B))) != 0 {
		return false
	}
	return ec.Equal(ec.Mul(pp.G, p.Z), ec.Add(p.A, ec.Mul(pp.PK[i], p.C))) && ec.Equal(ec.Mul(share, p.Z), ec.Add(p.B, ec.Mul(encrypted, p.C)))
}

func constrainedPolynomial(pp *Parameters, values []*big.Int) ([]*big.Int, error) {
	base := make([]*big.Int, 1)
	base[0] = new(big.Int)
	for d, value := range values {
		xd := packedX(pp.L, d)
		basis := []*big.Int{big.NewInt(1)}
		den := big.NewInt(1)
		for e := range values {
			if e != d {
				xe := packedX(pp.L, e)
				basis = mulPoly(basis, []*big.Int{ec.Mod(new(big.Int).Neg(xe)), big.NewInt(1)})
				den = ec.MulScalar(den, ec.SubScalar(xd, xe))
			}
		}
		scale := ec.MulScalar(value, ec.Inv(den))
		base = addPoly(base, scalePoly(basis, scale))
	}
	zeroes := []*big.Int{big.NewInt(1)}
	for d := range values {
		zeroes = mulPoly(zeroes, []*big.Int{ec.Mod(new(big.Int).Neg(packedX(pp.L, d))), big.NewInt(1)})
	}
	q := make([]*big.Int, pp.F)
	for i := range q {
		var err error
		q[i], err = ec.Scalar()
		if err != nil {
			return nil, err
		}
	}
	return pad(addPoly(base, mulPoly(zeroes, q)), pp.R), nil
}

func packedX(l, d int) *big.Int { return ec.Mod(big.NewInt(int64(d - (l - 1)))) }
func eval(poly []*big.Int, x *big.Int) *big.Int {
	acc := new(big.Int)
	for i := len(poly) - 1; i >= 0; i-- {
		acc = ec.AddScalar(ec.MulScalar(acc, x), poly[i])
	}
	return acc
}
func addPoly(a, b []*big.Int) []*big.Int {
	n := len(a)
	if len(b) > n {
		n = len(b)
	}
	o := make([]*big.Int, n)
	for i := 0; i < n; i++ {
		x, y := new(big.Int), new(big.Int)
		if i < len(a) {
			x = a[i]
		}
		if i < len(b) {
			y = b[i]
		}
		o[i] = ec.AddScalar(x, y)
	}
	return o
}
func scalePoly(a []*big.Int, s *big.Int) []*big.Int {
	o := make([]*big.Int, len(a))
	for i := range a {
		o[i] = ec.MulScalar(a[i], s)
	}
	return o
}
func mulPoly(a, b []*big.Int) []*big.Int {
	o := make([]*big.Int, len(a)+len(b)-1)
	for i := range o {
		o[i] = new(big.Int)
	}
	for i := range a {
		for j := range b {
			o[i+j] = ec.AddScalar(o[i+j], ec.MulScalar(a[i], b[j]))
		}
	}
	return o
}
func pad(a []*big.Int, n int) []*big.Int {
	o := make([]*big.Int, n)
	for i := range o {
		o[i] = new(big.Int)
		if i < len(a) {
			o[i].Set(a[i])
		}
	}
	return o
}
func lagrangeAt(index, count int, x *big.Int) *big.Int {
	num, den := big.NewInt(1), big.NewInt(1)
	xi := big.NewInt(int64(index))
	for j := 1; j <= count; j++ {
		if j == index {
			continue
		}
		xj := big.NewInt(int64(j))
		num = ec.MulScalar(num, ec.SubScalar(x, xj))
		den = ec.MulScalar(den, ec.SubScalar(xi, xj))
	}
	return ec.MulScalar(num, ec.Inv(den))
}
func lagrangeForIndices(index int, indices []int, x *big.Int) *big.Int {
	num, den := big.NewInt(1), big.NewInt(1)
	xi := big.NewInt(int64(index))
	for _, j := range indices {
		if j == index {
			continue
		}
		xj := big.NewInt(int64(j))
		num = ec.MulScalar(num, ec.SubScalar(x, xj))
		den = ec.MulScalar(den, ec.SubScalar(xi, xj))
	}
	return ec.MulScalar(num, ec.Inv(den))
}
func validPP(pp *Parameters) bool {
	return pp != nil && pp.N >= 1 && pp.L >= 1 && pp.F >= 0 && pp.R == pp.F+pp.L && pp.N >= 2*pp.F+pp.L && len(pp.PK) == pp.N && len(pp.SK) == pp.N
}
