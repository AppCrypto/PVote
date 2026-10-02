package priscore

import (
	bn256 "github.com/ethereum/go-ethereum/crypto/bn256/cloudflare"
	"testing"
)

func TestCompleteProtocol(t *testing.T) {
	pp, err := Setup(3, 3, 5)
	if err != nil {
		t.Fatal(err)
	}
	keys := make([]*VoterKey, pp.Voters)
	for i := range keys {
		keys[i], err = KeyGen(pp)
		if err != nil {
			t.Fatal(err)
		}
		for j:=0;j<pp.Candidates;j++{if _,err=KeyDerive(pp,keys[i],j);err!=nil{t.Fatal(err)}}
	}
	if err = CheckKeys(pp, keys); err != nil {
		t.Fatal(err)
	}
	scores := [][]int{{2, 1, 2}, {0, 3, 2}, {1, 4, 0}}
	commits := make([]*Commitment, pp.Voters)
	ballots := make([]*Ballot, pp.Voters)
	for i := range keys {
		commits[i], err = Commit(pp, i, keys[i], keys, scores[i])
		if err != nil {
			t.Fatalf("commit %d: %v", i, err)
		}
		if !VerifyCommit(pp, i, keys, commits[i]) {
			t.Fatalf("commit proof %d", i)
		}
		ballots[i], err = Vote(pp, i, keys[i], keys, commits[i], scores[i])
		if err != nil {
			t.Fatalf("vote %d: %v", i, err)
		}
		if !VerifyVote(pp, i, keys, commits[i], ballots[i]) {
			ys := make([]*bn256.G1, len(keys))
			for q, k := range keys {
				ys[q] = k.Y
			}
			w := signedProduct(ys, i)
			a := otherProduct(keys, i)
			for j := 0; j < pp.Candidates; j++ {
				d := make([]*bn256.G1, len(keys))
				for q, k := range keys {
					d[q] = k.DerivedY[j]
				}
				st := &voteStatement{keys[i].Y, keys[i].DerivedY[j], ballots[i].Ciphertexts[j], w, signedProduct(d, i)}
				if !verifyRange(pp, "priscore/vote-range", commits[i].Pairs[j], a, st, nil, ballots[i].Proof.RangeBinding[j]) {
					t.Fatalf("vote range proof %d/%d", i, j)
				}
			}
			t.Fatalf("vote proof %d", i)
		}
	}
	got, err := SelfTally(pp, ballots)
	if err != nil {
		t.Fatal(err)
	}
	want := []int{3, 8, 4}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("tally=%v want=%v", got, want)
		}
	}
	absent := 1
	partial := append([]*Ballot(nil), ballots...)
	partial[absent] = nil
	recovery := make([]*RecoveryShare, 0, pp.Voters-1)
	for i := range keys {
		if i != absent {
			r, e := RecoverForAbsent(pp, absent, i, keys, commits[absent])
			if e != nil {
				t.Fatal(e)
			}
			recovery = append(recovery, r)
		}
	}
	recovered, e := TallyWithAbort(pp, absent, keys, commits[absent], partial, recovery)
	if e != nil {
		t.Fatal(e)
	}
	for i := range want {
		if recovered[i] != want[i] {
			t.Fatalf("abort tally=%v want=%v", recovered, want)
		}
	}
	ballots[0].Ciphertexts[0].Gamma = pp.G
	if VerifyVote(pp, 0, keys, commits[0], ballots[0]) {
		t.Fatal("accepted modified ballot")
	}
}
