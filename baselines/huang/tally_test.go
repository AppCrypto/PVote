package huang

import "testing"

func TestTally(t *testing.T) {
	pp, err := Setup(3, 10)
	if err != nil {
		t.Fatal(err)
	}
	inputs := [][]int{{1, 2, 3}, {4, 0, 6}, {5, 8, 1}}
	ballots := make([]Ballot, len(inputs))
	for i := range inputs {
		ballots[i], err = Encrypt(pp, inputs[i])
		if err != nil {
			t.Fatal(err)
		}
	}
	got, err := Tally(pp, ballots)
	if err != nil {
		t.Fatal(err)
	}
	want := []int{10, 10, 10}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("coordinate %d: got %d want %d", i, got[i], want[i])
		}
	}
}

func TestRejectsMalformedBallot(t *testing.T) {
	pp, _ := Setup(2, 3)
	if _, err := Tally(pp, []Ballot{{}}); err == nil {
		t.Fatal("malformed ballot accepted")
	}
}
