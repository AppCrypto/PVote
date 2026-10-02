package threepvss

import (
	"math/big"
	"testing"

	"PVote/baselines/internal/ec"
)

func TestNativeReconstructions(t *testing.T) {
	pp, err := Setup(10, 2, 5)
	if err != nil {
		t.Fatal(err)
	}
	secrets := []*big.Int{big.NewInt(2), big.NewInt(4), big.NewInt(6), big.NewInt(8), big.NewInt(10)}
	tr, err := Share(pp, secrets)
	if err != nil {
		t.Fatal(err)
	}
	if !Verify(pp, tr) {
		t.Fatal("native transcript rejected")
	}
	optimistic, err := OptimisticReconstruct(pp, tr, secrets)
	if err != nil {
		t.Fatal(err)
	}
	pessimistic, err := PessimisticReconstruct(pp, tr, []int{1, 3, 4, 6, 7, 8, 10})
	if err != nil {
		t.Fatal(err)
	}
	for i, secret := range secrets {
		want := ec.Mul(pp.G, secret)
		if !ec.Equal(optimistic[i], want) || !ec.Equal(pessimistic[i], want) {
			t.Fatalf("secret %d reconstruction mismatch", i)
		}
	}
}

func TestScoreRoundTrip(t *testing.T) {
	pp, err := Setup(10, 2, 5)
	if err != nil {
		t.Fatal(err)
	}
	scores := []int{0, 1, 3, 5, 7}
	ballots := make([]*ScoreBallot, 4)
	for i := range ballots {
		ballots[i], err = ShareScore(pp, scores, 3)
		if err != nil {
			t.Fatal(err)
		}
		if !VerifyScore(pp, ballots[i]) {
			t.Fatal("valid ballot rejected")
		}
	}
	totals, err := Tally(pp, ballots)
	if err != nil {
		t.Fatal(err)
	}
	for i, got := range totals {
		if want := len(ballots) * scores[i]; got != want {
			t.Fatalf("coordinate %d: got %d want %d", i, got, want)
		}
	}
}

func TestRejectsTampering(t *testing.T) {
	pp, err := Setup(10, 2, 5)
	if err != nil {
		t.Fatal(err)
	}
	b, err := ShareScore(pp, []int{0, 1, 2, 3, 4}, 3)
	if err != nil {
		t.Fatal(err)
	}
	b.Planes[0].U[0] = b.Planes[0].U[1]
	if VerifyScore(pp, b) {
		t.Fatal("tampered ballot accepted")
	}
}
