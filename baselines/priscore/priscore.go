// Package priscore reproduces the complete PriScore protocol from Algorithms
// 1--4 of Yang et al., including both voting-round proof systems and abort
// recovery. Multiplicative notation in the paper is represented additively on
// BN254 G1.
package priscore

import (
	"errors"
	"fmt"
	"math/big"

	"PVote/baselines/internal/ec"
	bn256 "github.com/ethereum/go-ethereum/crypto/bn256/cloudflare"
)

var ErrInvalid = errors.New("invalid PriScore input or proof")

type Parameters struct {
	G                          *bn256.G1
	Voters, Candidates, Budget int
}
type VoterKey struct {
	X        *big.Int
	Y        *bn256.G1
	DerivedX []*big.Int
	DerivedY []*bn256.G1
}
type CommitmentPair struct{ Zeta, Eta *bn256.G1 }
type Ciphertext struct{ Beta, Gamma *bn256.G1 }

type DLEQProof struct{ C, Z *big.Int }
type ORBranch struct {
	C *big.Int
	Z []*big.Int
}
type ORProof struct{ Branches []ORBranch }
type CommitProof struct {
	Owner DLEQProof
	Range []ORProof
	Sum   DLEQProof
}
type VoteProof struct {
	RangeBinding []ORProof
	SumOwner     LinearProof
}
type Commitment struct {
	Pairs    []CommitmentPair
	Proof    CommitProof
	openings []*big.Int
}
type Ballot struct {
	Ciphertexts []Ciphertext
	Proof       VoteProof
}
type RecoveryShare struct {
	Voter           int
	WBar            *bn256.G1
	ZetaTilde, ZBar []*bn256.G1
	ProofX          DLEQProof
	ProofXZ         []LinearProof
}

// LinearProof proves several linear representation equations with one common
// Fiat--Shamir challenge. Responses correspond to the witness vector.
type LinearProof struct {
	Commitments []*bn256.G1
	C           *big.Int
	Z           []*big.Int
}

func Setup(voters, candidates, budget int) (*Parameters, error) {
	if voters < 2 || candidates < 1 || budget < 0 {
		return nil, ErrInvalid
	}
	gx, err := ec.Scalar()
	if err != nil {
		return nil, err
	}
	return &Parameters{G: new(bn256.G1).ScalarBaseMult(gx), Voters: voters, Candidates: candidates, Budget: budget}, nil
}

func KeyGen(pp *Parameters) (*VoterKey, error) {
	if pp == nil {
		return nil, ErrInvalid
	}
	x, err := ec.Scalar()
	if err != nil {
		return nil, err
	}
	k := &VoterKey{X: x, Y: ec.Mul(pp.G, x), DerivedX: make([]*big.Int, pp.Candidates), DerivedY: make([]*bn256.G1, pp.Candidates)}
	return k, nil
}

// KeyDerive implements x_i,j=H(X_i,C_j,r_j), y_i,j=g^x_i,j.
func KeyDerive(pp *Parameters, key *VoterKey, candidate int) (*bn256.G1, error) {
	if pp == nil || key == nil || candidate < 0 || candidate >= pp.Candidates {
		return nil, ErrInvalid
	}
	nonce, e := ec.Scalar()
	if e != nil {
		return nil, e
	}
	x := ec.HashScalar("priscore/key-derive", ec.ScalarBytes(key.X, nonce, big.NewInt(int64(candidate))))
	if x.Sign() == 0 {
		x.SetInt64(1)
	}
	key.DerivedX[candidate] = x
	key.DerivedY[candidate] = ec.Mul(pp.G, x)
	return key.DerivedY[candidate], nil
}

func validateScores(pp *Parameters, scores []int) bool {
	if pp == nil || len(scores) != pp.Candidates {
		return false
	}
	sum := 0
	for _, p := range scores {
		if p < 0 || p > pp.Budget {
			return false
		}
		sum += p
	}
	return sum == pp.Budget
}

func otherProduct(keys []*VoterKey, skip int) *bn256.G1 {
	points := make([]*bn256.G1, 0, len(keys)-1)
	for i, k := range keys {
		if i != skip {
			points = append(points, k.Y)
		}
	}
	return ec.Product(points)
}
func signedProduct(points []*bn256.G1, index int) *bn256.G1 {
	acc := ec.Zero()
	for i, p := range points {
		if i < index {
			acc = ec.Add(acc, p)
		} else if i > index {
			acc = ec.Sub(acc, p)
		}
	}
	return acc
}

func Commit(pp *Parameters, voter int, key *VoterKey, all []*VoterKey, scores []int) (*Commitment, error) {
	if !validateScores(pp, scores) || key == nil || voter < 0 || voter >= len(all) || len(all) != pp.Voters || !ec.Equal(key.Y, all[voter].Y) {
		return nil, ErrInvalid
	}
	a := otherProduct(all, voter)
	out := &Commitment{Pairs: make([]CommitmentPair, pp.Candidates)}
	secrets := make([]*big.Int, pp.Candidates)
	for j, p := range scores {
		s, e := ec.Scalar()
		if e != nil {
			return nil, e
		}
		secrets[j] = s
		out.Pairs[j] = CommitmentPair{ec.Mul(pp.G, s), ec.Add(ec.Mul(pp.G, big.NewInt(int64(p))), ec.Mul(a, s))}
	}
	out.openings = secrets
	out.Proof.Owner = proveDLEQ("priscore/commit-owner", pp.G, key.Y, pp.G, key.Y, key.X)
	out.Proof.Range = make([]ORProof, pp.Candidates)
	for j, p := range scores {
		out.Proof.Range[j] = proveRange(pp, "priscore/commit-range", out.Pairs[j], a, p, secrets[j], nil, nil)
	}
	zeta := make([]*bn256.G1, pp.Candidates)
	eta := make([]*bn256.G1, pp.Candidates)
	sumS := new(big.Int)
	for j, c := range out.Pairs {
		zeta[j] = c.Zeta
		eta[j] = c.Eta
		sumS = ec.AddScalar(sumS, secrets[j])
	}
	out.Proof.Sum = proveDLEQ("priscore/commit-sum", pp.G, ec.Product(zeta), a, ec.Sub(ec.Product(eta), ec.Mul(pp.G, big.NewInt(int64(pp.Budget)))), sumS)
	return out, nil
}

func VerifyCommit(pp *Parameters, voter int, all []*VoterKey, c *Commitment) bool {
	if pp == nil || c == nil || voter < 0 || voter >= len(all) || len(c.Pairs) != pp.Candidates || len(c.Proof.Range) != pp.Candidates {
		return false
	}
	a := otherProduct(all, voter)
	if !verifyDLEQ("priscore/commit-owner", pp.G, all[voter].Y, pp.G, all[voter].Y, c.Proof.Owner) {
		return false
	}
	zetas := make([]*bn256.G1, len(c.Pairs))
	etas := make([]*bn256.G1, len(c.Pairs))
	for j, pair := range c.Pairs {
		if !verifyRange(pp, "priscore/commit-range", pair, a, nil, nil, c.Proof.Range[j]) {
			return false
		}
		zetas[j] = pair.Zeta
		etas[j] = pair.Eta
	}
	return verifyDLEQ("priscore/commit-sum", pp.G, ec.Product(zetas), a, ec.Sub(ec.Product(etas), ec.Mul(pp.G, big.NewInt(int64(pp.Budget)))), c.Proof.Sum)
}

func Vote(pp *Parameters, voter int, key *VoterKey, all []*VoterKey, commit *Commitment, scores []int) (*Ballot, error) {
	if !validateScores(pp, scores) || key == nil || key.X == nil || key.Y == nil || commit == nil || voter < 0 ||
		voter >= len(all) || len(all) != pp.Voters || len(key.DerivedX) != pp.Candidates ||
		len(commit.Pairs) != pp.Candidates || len(commit.openings) != pp.Candidates ||
		all[voter] == nil || !ec.Equal(key.Y, all[voter].Y) {
		return nil, ErrInvalid
	}
	for _, candidateKey := range all {
		if candidateKey == nil || candidateKey.Y == nil || len(candidateKey.DerivedY) != pp.Candidates {
			return nil, ErrInvalid
		}
	}
	for _, scalar := range key.DerivedX {
		if scalar == nil {
			return nil, ErrInvalid
		}
	}
	for i, pair := range commit.Pairs {
		if pair.Zeta == nil || pair.Eta == nil || commit.openings[i] == nil {
			return nil, ErrInvalid
		}
	}
	ys := make([]*bn256.G1, len(all))
	for i, k := range all {
		ys[i] = k.Y
	}
	w := signedProduct(ys, voter)
	out := &Ballot{Ciphertexts: make([]Ciphertext, pp.Candidates)}
	rs := make([]*big.Int, pp.Candidates)
	for j, p := range scores {
		derived := make([]*bn256.G1, len(all))
		for i, k := range all {
			derived[i] = k.DerivedY[j]
		}
		z := signedProduct(derived, voter)
		r, e := ec.Scalar()
		if e != nil {
			return nil, e
		}
		rs[j] = r
		out.Ciphertexts[j] = Ciphertext{Beta: ec.Add(ec.Mul(z, key.DerivedX[j]), ec.Mul(pp.G, r)), Gamma: ec.Add(ec.Add(ec.Mul(pp.G, big.NewInt(int64(p))), ec.Mul(w, key.X)), ec.Mul(pp.G, r))}
	}
	out.Proof.RangeBinding = make([]ORProof, pp.Candidates)
	for j, p := range scores {
		derived := make([]*bn256.G1, len(all))
		for i, k := range all {
			derived[i] = k.DerivedY[j]
		}
		z := signedProduct(derived, voter)
		out.Proof.RangeBinding[j] = proveRange(pp, "priscore/vote-range", commit.Pairs[j], otherProduct(all, voter), p, commit.openings[j], &voteWitness{key.X, key.DerivedX[j], rs[j]}, &voteStatement{all[voter].Y, key.DerivedY[j], out.Ciphertexts[j], w, z})
	}
	sumR := new(big.Int)
	for _, r := range rs {
		sumR = ec.AddScalar(sumR, r)
	}
	gammas := make([]*bn256.G1, len(out.Ciphertexts))
	for j, c := range out.Ciphertexts {
		gammas[j] = c.Gamma
	}
	target := ec.Sub(ec.Product(gammas), ec.Mul(pp.G, big.NewInt(int64(pp.Budget))))
	ncW := ec.Mul(w, big.NewInt(int64(pp.Candidates)))
	out.Proof.SumOwner = proveLinear("priscore/vote-sum", [][]*bn256.G1{{pp.G}, {ncW, pp.G}}, []*bn256.G1{all[voter].Y, target}, []*big.Int{key.X, sumR}, [][]int{{0}, {0, 1}})
	return out, nil
}

func VerifyVote(pp *Parameters, voter int, all []*VoterKey, commit *Commitment, b *Ballot) bool {
	if !VerifyCommit(pp, voter, all, commit) || b == nil || len(b.Ciphertexts) != pp.Candidates || len(b.Proof.RangeBinding) != pp.Candidates {
		return false
	}
	ys := make([]*bn256.G1, len(all))
	for i, k := range all {
		ys[i] = k.Y
	}
	w := signedProduct(ys, voter)
	a := otherProduct(all, voter)
	for j := 0; j < pp.Candidates; j++ {
		derived := make([]*bn256.G1, len(all))
		for i, k := range all {
			derived[i] = k.DerivedY[j]
		}
		z := signedProduct(derived, voter)
		st := &voteStatement{all[voter].Y, all[voter].DerivedY[j], b.Ciphertexts[j], w, z}
		if !verifyRange(pp, "priscore/vote-range", commit.Pairs[j], a, st, nil, b.Proof.RangeBinding[j]) {
			return false
		}
	}
	gammas := make([]*bn256.G1, len(b.Ciphertexts))
	for j, c := range b.Ciphertexts {
		gammas[j] = c.Gamma
	}
	target := ec.Sub(ec.Product(gammas), ec.Mul(pp.G, big.NewInt(int64(pp.Budget))))
	ncW := ec.Mul(w, big.NewInt(int64(pp.Candidates)))
	return verifyLinear("priscore/vote-sum", [][]*bn256.G1{{pp.G}, {ncW, pp.G}}, []*bn256.G1{all[voter].Y, target}, [][]int{{0}, {0, 1}}, b.Proof.SumOwner)
}

func SelfTally(pp *Parameters, ballots []*Ballot) ([]int, error) {
	if pp == nil || len(ballots) != pp.Voters {
		return nil, ErrInvalid
	}
	totals := make([]int, pp.Candidates)
	for j := 0; j < pp.Candidates; j++ {
		bs := make([]*bn256.G1, len(ballots))
		gs := make([]*bn256.G1, len(ballots))
		for i, b := range ballots {
			if b == nil || len(b.Ciphertexts) != pp.Candidates {
				return nil, ErrInvalid
			}
			bs[i] = b.Ciphertexts[j].Beta
			gs[i] = b.Ciphertexts[j].Gamma
		}
		x, e := ec.DLog(pp.G, ec.Sub(ec.Product(gs), ec.Product(bs)), pp.Voters*pp.Budget)
		if e != nil {
			return nil, e
		}
		totals[j] = x
	}
	return totals, nil
}

// RecoverForAbsent produces one remaining voter's Algorithm 4 correction
// values and their proofs.
func RecoverForAbsent(pp *Parameters, absent, voter int, all []*VoterKey, commit *Commitment) (*RecoveryShare, error) {
	if pp == nil || absent < 0 || absent >= len(all) || voter < 0 || voter >= len(all) || voter == absent || commit == nil {
		return nil, ErrInvalid
	}
	sign := int64(1)
	if absent < voter {
		sign = -1
	}
	baseY := ec.Mul(all[absent].Y, big.NewInt(sign))
	r := &RecoveryShare{Voter: voter, WBar: ec.Mul(baseY, all[voter].X), ZetaTilde: make([]*bn256.G1, pp.Candidates), ZBar: make([]*bn256.G1, pp.Candidates), ProofXZ: make([]LinearProof, pp.Candidates)}
	r.ProofX = proveDLEQ("priscore/recover-X", pp.G, all[voter].Y, baseY, r.WBar, all[voter].X)
	for j := 0; j < pp.Candidates; j++ {
		baseZ := ec.Mul(all[absent].DerivedY[j], big.NewInt(sign))
		r.ZetaTilde[j] = ec.Mul(commit.Pairs[j].Zeta, all[voter].X)
		r.ZBar[j] = ec.Mul(baseZ, all[voter].DerivedX[j])
		r.ProofXZ[j] = proveLinear("priscore/recover-XZ", [][]*bn256.G1{{pp.G}, {pp.G}, {commit.Pairs[j].Zeta}, {baseZ}}, []*bn256.G1{all[voter].Y, all[voter].DerivedY[j], r.ZetaTilde[j], r.ZBar[j]}, []*big.Int{all[voter].X, all[voter].DerivedX[j]}, [][]int{{0}, {1}, {0}, {1}})
	}
	return r, nil
}

func VerifyRecovery(pp *Parameters, absent int, all []*VoterKey, commit *Commitment, r *RecoveryShare) bool {
	if pp == nil || r == nil || r.Voter == absent || r.Voter < 0 || r.Voter >= len(all) || len(r.ZBar) != pp.Candidates || len(r.ProofXZ) != pp.Candidates {
		return false
	}
	sign := int64(1)
	if absent < r.Voter {
		sign = -1
	}
	baseY := ec.Mul(all[absent].Y, big.NewInt(sign))
	if !verifyDLEQ("priscore/recover-X", pp.G, all[r.Voter].Y, baseY, r.WBar, r.ProofX) {
		return false
	}
	for j := 0; j < pp.Candidates; j++ {
		baseZ := ec.Mul(all[absent].DerivedY[j], big.NewInt(sign))
		if !verifyLinear("priscore/recover-XZ", [][]*bn256.G1{{pp.G}, {pp.G}, {commit.Pairs[j].Zeta}, {baseZ}}, []*bn256.G1{all[r.Voter].Y, all[r.Voter].DerivedY[j], r.ZetaTilde[j], r.ZBar[j]}, [][]int{{0}, {1}, {0}, {1}}, r.ProofXZ[j]) {
			return false
		}
	}
	return true
}

// TallyWithAbort implements Algorithm 4. ballots is indexed by voter and has
// nil at absent; one verified correction share is required from every other
// voter.
func TallyWithAbort(pp *Parameters, absent int, all []*VoterKey, commit *Commitment, ballots []*Ballot, recovery []*RecoveryShare) ([]int, error) {
	if pp == nil || absent < 0 || absent >= len(all) || len(ballots) != pp.Voters || ballots[absent] != nil || len(recovery) != pp.Voters-1 {
		return nil, ErrInvalid
	}
	by := map[int]*RecoveryShare{}
	for _, r := range recovery {
		if !VerifyRecovery(pp, absent, all, commit, r) || by[r.Voter] != nil {
			return nil, ErrInvalid
		}
		by[r.Voter] = r
	}
	totals := make([]int, pp.Candidates)
	for j := 0; j < pp.Candidates; j++ {
		zetaParts := make([]*bn256.G1, 0, pp.Voters-1)
		omega1 := make([]*bn256.G1, 0, 2*(pp.Voters-1))
		omega2 := make([]*bn256.G1, 0, 2*(pp.Voters-1))
		for i, b := range ballots {
			if i == absent {
				continue
			}
			r := by[i]
			if b == nil || r == nil {
				return nil, ErrInvalid
			}
			zetaParts = append(zetaParts, r.ZetaTilde[j])
			omega1 = append(omega1, b.Ciphertexts[j].Beta, r.ZBar[j])
			omega2 = append(omega2, b.Ciphertexts[j].Gamma, r.WBar)
		}
		missing, e := ec.DLog(pp.G, ec.Sub(commit.Pairs[j].Eta, ec.Product(zetaParts)), pp.Budget)
		if e != nil {
			return nil, e
		}
		present, e := ec.DLog(pp.G, ec.Sub(ec.Product(omega2), ec.Product(omega1)), (pp.Voters-1)*pp.Budget)
		if e != nil {
			return nil, e
		}
		totals[j] = missing + present
	}
	return totals, nil
}

type voteWitness struct{ X, x, r *big.Int }
type voteStatement struct {
	Y, y   *bn256.G1
	Cipher Ciphertext
	W, Z   *bn256.G1
}

// proveRange implements PriScore's one-out-of-(B+1) proof. For Commit it is
// a two-equation DLEQ OR proof; for Vote it additionally binds the same score
// to the commitment, voter keys, and encrypted ballot.
func proveRange(pp *Parameters, domain string, pair CommitmentPair, a *bn256.G1, actual int, s *big.Int, vw *voteWitness, vs *voteStatement) ORProof {
	branches := make([]ORBranch, pp.Budget+1)
	commits := make([][]*bn256.G1, len(branches))
	randomReal := make([]*big.Int, 0)
	sumFalse := new(big.Int)
	for k := range branches {
		bases, targets, mapIdx := rangeEquations(pp, pair, a, k, vs)
		wc := 1
		if vs != nil {
			wc = 4
		}
		branches[k].Z = make([]*big.Int, wc)
		commits[k] = make([]*bn256.G1, len(targets))
		if k == actual {
			randomReal = make([]*big.Int, wc)
			for z := range randomReal {
				randomReal[z], _ = ec.Scalar()
			}
			for e := range targets {
				acc := ec.Zero()
				for q, b := range bases[e] {
					acc = ec.Add(acc, ec.Mul(b, randomReal[mapIdx[e][q]]))
				}
				commits[k][e] = acc
			}
		} else {
			branches[k].C, _ = ec.Scalar()
			sumFalse = ec.AddScalar(sumFalse, branches[k].C)
			for z := range branches[k].Z {
				branches[k].Z[z], _ = ec.Scalar()
			}
			for e, t := range targets {
				acc := ec.Zero()
				for q, b := range bases[e] {
					acc = ec.Add(acc, ec.Mul(b, branches[k].Z[mapIdx[e][q]]))
				}
				commits[k][e] = ec.Sub(acc, ec.Mul(t, branches[k].C))
			}
		}
	}
	h := rangeChallenge(domain, pair, a, vs, commits)
	branches[actual].C = ec.SubScalar(h, sumFalse)
	w := []*big.Int{s}
	if vs != nil {
		w = []*big.Int{s, vw.X, vw.x, vw.r}
	}
	for z := range branches[actual].Z {
		branches[actual].Z[z] = ec.AddScalar(randomReal[z], ec.MulScalar(branches[actual].C, w[z]))
	}
	return ORProof{branches}
}

func verifyRange(pp *Parameters, domain string, pair CommitmentPair, a *bn256.G1, vs *voteStatement, _ *voteWitness, proof ORProof) bool {
	if len(proof.Branches) != pp.Budget+1 {
		return false
	}
	commits := make([][]*bn256.G1, len(proof.Branches))
	sum := new(big.Int)
	for k, br := range proof.Branches {
		bases, targets, mapIdx := rangeEquations(pp, pair, a, k, vs)
		wc := 1
		if vs != nil {
			wc = 4
		}
		if br.C == nil || len(br.Z) != wc {
			return false
		}
		sum = ec.AddScalar(sum, br.C)
		commits[k] = make([]*bn256.G1, len(targets))
		for e, t := range targets {
			acc := ec.Zero()
			for q, b := range bases[e] {
				acc = ec.Add(acc, ec.Mul(b, br.Z[mapIdx[e][q]]))
			}
			commits[k][e] = ec.Sub(acc, ec.Mul(t, br.C))
		}
	}
	return sum.Cmp(rangeChallenge(domain, pair, a, vs, commits)) == 0
}

func rangeEquations(pp *Parameters, pair CommitmentPair, a *bn256.G1, k int, vs *voteStatement) ([][]*bn256.G1, []*bn256.G1, [][]int) {
	bases := [][]*bn256.G1{{pp.G}, {a}}
	targets := []*bn256.G1{pair.Zeta, ec.Sub(pair.Eta, ec.Mul(pp.G, big.NewInt(int64(k))))}
	idx := [][]int{{0}, {0}}
	if vs != nil {
		bases = append(bases, []*bn256.G1{pp.G}, []*bn256.G1{pp.G}, []*bn256.G1{vs.Z, pp.G}, []*bn256.G1{vs.W, pp.G})
		targets = append(targets, vs.Y, vs.y, vs.Cipher.Beta, ec.Sub(vs.Cipher.Gamma, ec.Mul(pp.G, big.NewInt(int64(k)))))
		idx = append(idx, []int{1}, []int{2}, []int{2, 3}, []int{1, 3})
	}
	return bases, targets, idx
}

func rangeChallenge(domain string, pair CommitmentPair, a *bn256.G1, vs *voteStatement, commits [][]*bn256.G1) *big.Int {
	parts := [][]byte{ec.PointBytes(pair.Zeta, pair.Eta, a)}
	if vs != nil {
		parts = append(parts, ec.PointBytes(vs.Y, vs.y, vs.Cipher.Beta, vs.Cipher.Gamma, vs.W, vs.Z))
	}
	for _, row := range commits {
		parts = append(parts, ec.PointBytes(row...))
	}
	return ec.HashScalar(domain, parts...)
}

func proveDLEQ(domain string, b1, t1, b2, t2 *bn256.G1, w *big.Int) DLEQProof {
	r, _ := ec.Scalar()
	a1 := ec.Mul(b1, r)
	a2 := ec.Mul(b2, r)
	c := ec.HashScalar(domain, ec.PointBytes(b1, t1, b2, t2, a1, a2))
	return DLEQProof{c, ec.AddScalar(r, ec.MulScalar(c, w))}
}
func verifyDLEQ(domain string, b1, t1, b2, t2 *bn256.G1, p DLEQProof) bool {
	if p.C == nil || p.Z == nil {
		return false
	}
	a1 := ec.Sub(ec.Mul(b1, p.Z), ec.Mul(t1, p.C))
	a2 := ec.Sub(ec.Mul(b2, p.Z), ec.Mul(t2, p.C))
	return p.C.Cmp(ec.HashScalar(domain, ec.PointBytes(b1, t1, b2, t2, a1, a2))) == 0
}

func proveLinear(domain string, bases [][]*bn256.G1, targets []*bn256.G1, w []*big.Int, index [][]int) LinearProof {
	r := make([]*big.Int, len(w))
	for i := range r {
		r[i], _ = ec.Scalar()
	}
	a := make([]*bn256.G1, len(targets))
	for e := range targets {
		a[e] = ec.Zero()
		for j, b := range bases[e] {
			a[e] = ec.Add(a[e], ec.Mul(b, r[index[e][j]]))
		}
	}
	c := linearChallenge(domain, bases, targets, a)
	z := make([]*big.Int, len(w))
	for i := range z {
		z[i] = ec.AddScalar(r[i], ec.MulScalar(c, w[i]))
	}
	return LinearProof{a, c, z}
}
func verifyLinear(domain string, bases [][]*bn256.G1, targets []*bn256.G1, index [][]int, p LinearProof) bool {
	if len(p.Commitments) != len(targets) || p.C == nil {
		return false
	}
	for e := range targets {
		acc := ec.Zero()
		for j, b := range bases[e] {
			if index[e][j] >= len(p.Z) {
				return false
			}
			acc = ec.Add(acc, ec.Mul(b, p.Z[index[e][j]]))
		}
		if !ec.Equal(acc, ec.Add(p.Commitments[e], ec.Mul(targets[e], p.C))) {
			return false
		}
	}
	return p.C.Cmp(linearChallenge(domain, bases, targets, p.Commitments)) == 0
}
func linearChallenge(domain string, bases [][]*bn256.G1, targets, commitments []*bn256.G1) *big.Int {
	parts := make([][]byte, 0, len(targets)*3)
	for i := range targets {
		parts = append(parts, ec.PointBytes(bases[i]...), ec.PointBytes(targets[i]), ec.PointBytes(commitments[i]))
	}
	return ec.HashScalar(domain, parts...)
}

func CheckKeys(pp *Parameters, keys []*VoterKey) error {
	if pp == nil || len(keys) != pp.Voters {
		return ErrInvalid
	}
	for i, k := range keys {
		if k == nil || len(k.DerivedY) != pp.Candidates || !ec.Equal(k.Y, ec.Mul(pp.G, k.X)) {
			return fmt.Errorf("%w: voter %d", ErrInvalid, i)
		}
		for j := range k.DerivedY {
			if !ec.Equal(k.DerivedY[j], ec.Mul(pp.G, k.DerivedX[j])) {
				return fmt.Errorf("%w: voter %d candidate %d", ErrInvalid, i, j)
			}
		}
	}
	return nil
}

// PublicTranscriptBytes is the exact fixed-width size of this implementation's
// public Commit and Vote objects (64-byte G1 points and 32-byte scalars).
func PublicTranscriptBytes(pp *Parameters) int {
	if pp == nil {
		return 0
	}
	b := pp.Budget + 1
	commit := 2*pp.Candidates*64 + 2*32 + pp.Candidates*b*2*32 + 2*32
	vote := 2*pp.Candidates*64 + pp.Candidates*b*5*32 + 2*64 + 3*32
	return commit + vote
}
