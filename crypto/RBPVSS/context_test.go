package rbpvss

import (
	"bytes"
	"errors"
	"math/big"
	"testing"

	bn256 "github.com/ethereum/go-ethereum/crypto/bn256/cloudflare"
)

func replayParameters(t *testing.T) (*PublicParameters, []*big.Int) {
	t.Helper()
	pp, keys, err := Setup(128, 4, 3, 2, Interval{Min: 0, Max: 5})
	if err != nil {
		t.Fatal(err)
	}
	return pp, keys
}

func contextBallot(t *testing.T, pp *PublicParameters, ctx BallotContext, scores []int64) *Instance {
	t.Helper()
	ballot, err := ShareWithContext(pp, scores, ctx)
	if err != nil {
		t.Fatal(err)
	}
	if !DVerifyWithContext(pp, ballot, ctx) {
		t.Fatal("honest context-bound ballot rejected")
	}
	return ballot
}

func TestReplayAcrossVotersAndElections(t *testing.T) {
	pp, _ := replayParameters(t)
	source := BallotContext{ElectionID: "election-A", VoterIndex: 1}
	ballot := contextBallot(t, pp, source, []int64{2, 4})
	for name, destination := range map[string]BallotContext{
		"another voter":    {ElectionID: "election-A", VoterIndex: 2},
		"another election": {ElectionID: "election-B", VoterIndex: 1},
		"both changed":     {ElectionID: "election-B", VoterIndex: 2},
	} {
		t.Run(name, func(t *testing.T) {
			if DVerifyWithContext(pp, ballot, destination) {
				t.Fatal("accepted copied ballot in a different authenticated context")
			}
		})
	}
	if DVerify(pp, ballot) {
		t.Fatal("context-bound proof accepted by the standalone verifier")
	}
	standalone, err := Share(pp, []int64{2, 4})
	if err != nil {
		t.Fatal(err)
	}
	if DVerifyWithContext(pp, standalone, source) {
		t.Fatal("standalone proof bypassed context-bound verification")
	}
}

func TestReplayRangeProofUnderAnotherContext(t *testing.T) {
	pp, _ := replayParameters(t)
	alice := BallotContext{ElectionID: "A", VoterIndex: 1}
	ballot := contextBallot(t, pp, alice, []int64{1, 4})
	binding, _ := ballotBinding(pp, alice)
	other, _ := ballotBinding(pp, BallotContext{ElectionID: "A", VoterIndex: 2})
	if !verifyRangeBinding(pp, 0, ballot.Pi[0], ballot.U[0], ballot.V[0], binding) {
		t.Fatal("original range proof rejected")
	}
	if verifyRangeBinding(pp, 0, ballot.Pi[0], ballot.U[0], ballot.V[0], other) {
		t.Fatal("range proof itself is not context-bound")
	}
	if verifyRangeBinding(pp, 1, ballot.Pi[0], ballot.U[0], ballot.V[0], binding) {
		t.Fatal("range proof accepted at another coordinate")
	}
}

// This harness models only the paper's trusted identity->index boundary and
// acceptance state machine; it is not an EVM integration or wallet-auth test.
type replayElection struct {
	pp       *PublicParameters
	id       string
	accepted map[uint64]bool
	tally    *Aggregate
}

func newReplayElection(t *testing.T, pp *PublicParameters, id string) *replayElection {
	t.Helper()
	state, err := NewAggregate(pp)
	if err != nil {
		t.Fatal(err)
	}
	return &replayElection{pp: pp, id: id, accepted: make(map[uint64]bool), tally: state}
}

func (e *replayElection) submit(authenticatedIndex uint64, ballot *Instance) error {
	if e.accepted[authenticatedIndex] {
		return errors.New("voter already accepted")
	}
	ctx := BallotContext{ElectionID: e.id, VoterIndex: authenticatedIndex}
	if err := AggregateVerifiedWithContext(e.pp, e.tally, ballot, ctx); err != nil {
		return err
	}
	e.accepted[authenticatedIndex] = true
	return nil
}

func aggregateEncoding(a *Aggregate) []byte {
	var out []byte
	for _, point := range a.C {
		out = append(out, point.Marshal()...)
	}
	for _, point := range a.U {
		out = append(out, point.Marshal()...)
	}
	return out
}

func TestReplayAcceptanceStateAndEqualScores(t *testing.T) {
	pp, keys := replayParameters(t)
	e := newReplayElection(t, pp, "A")
	alice := contextBallot(t, pp, BallotContext{ElectionID: "A", VoterIndex: 1}, []int64{2, 4})
	if err := e.submit(1, alice); err != nil {
		t.Fatal(err)
	}
	before := aggregateEncoding(e.tally)
	for _, index := range []uint64{1, 2} {
		if err := e.submit(index, alice); err == nil {
			t.Fatal("accepted duplicate or cross-voter replay")
		}
		if e.tally.Count != 1 || !bytes.Equal(before, aggregateEncoding(e.tally)) || e.accepted[2] {
			t.Fatal("rejected replay changed acceptance or aggregate state")
		}
	}
	// Equal plaintext scores remain legitimate; equality of scores is not a
	// reason to reject independently generated, correctly labeled ballots.
	bob := contextBallot(t, pp, BallotContext{ElectionID: "A", VoterIndex: 2}, []int64{2, 4})
	if err := e.submit(2, bob); err != nil {
		t.Fatal(err)
	}
	assertReplayTally(t, pp, keys, e.tally, []int64{4, 8})
}

func TestReplayConcurrentElectionsReuseTallierKeys(t *testing.T) {
	pp, keys := replayParameters(t)
	a := newReplayElection(t, pp, "open-A")
	b := newReplayElection(t, pp, "early-B")
	target := contextBallot(t, pp, BallotContext{ElectionID: a.id, VoterIndex: 1}, []int64{3, 5})
	if err := a.submit(1, target); err != nil {
		t.Fatal(err)
	}
	before := aggregateEncoding(b.tally)
	if err := b.submit(2, target); err == nil {
		t.Fatal("accepted open election's ballot into earlier closing election")
	}
	if b.tally.Count != 0 || !bytes.Equal(before, aggregateEncoding(b.tally)) || b.accepted[2] {
		t.Fatal("cross-election replay changed destination state")
	}
	known := contextBallot(t, pp, BallotContext{ElectionID: b.id, VoterIndex: 2}, []int64{1, 2})
	if err := b.submit(2, known); err != nil {
		t.Fatal(err)
	}
	assertReplayTally(t, pp, keys, b.tally, []int64{1, 2})
	if a.tally.Count != 1 || !a.accepted[1] {
		t.Fatal("finishing the other election changed the source election")
	}
}

func assertReplayTally(t *testing.T, pp *PublicParameters, keys []*big.Int, a *Aggregate, want []int64) {
	t.Helper()
	shares := make([]IndexedShare, pp.Threshold)
	for i := range shares {
		share, proof, err := Decrypt(pp, a.C[i], keys[i])
		if err != nil || !PVerify(pp, i+1, share, a.C[i], proof) {
			t.Fatalf("honest decryption share rejected: %v", err)
		}
		shares[i] = IndexedShare{Index: i + 1, Share: share}
	}
	got, err := Recon(pp, shares, a.U, Interval{Min: int64(a.Count) * pp.Range.Min, Max: int64(a.Count) * pp.Range.Max})
	if err != nil || !equalInt64s(got, want) {
		t.Fatalf("tally=%v, want %v, err=%v", got, want, err)
	}
}

func TestContextParametersAndCanonicalResponses(t *testing.T) {
	pp, _ := replayParameters(t)
	ctx := BallotContext{ElectionID: "A", VoterIndex: 1}
	ballot := contextBallot(t, pp, ctx, []int64{1, 2})
	original, _ := ballotBinding(pp, ctx)
	changed := *pp
	changed.PP2 = pp.PP2
	changed.PP2.PKI = new(bn256.G2).ScalarMult(pp.PP2.G1, big.NewInt(7))
	modified, err := ballotBinding(&changed, ctx)
	if err != nil || bytes.Equal(original, modified) || DVerifyWithContext(&changed, ballot, ctx) {
		t.Fatal("issuer-key substitution was not bound")
	}
	ballot.Phi.Z[0] = new(big.Int).Add(ballot.Phi.Z[0], bn256.Order)
	if DVerifyWithContext(pp, ballot, ctx) {
		t.Fatal("accepted noncanonical response encoding")
	}
	for _, bad := range []BallotContext{{ElectionID: "A"}, {VoterIndex: 1}} {
		if _, err := ShareWithContext(pp, []int64{1, 2}, bad); err == nil || DVerifyWithContext(pp, ballot, bad) {
			t.Fatal("accepted an incomplete authenticated context")
		}
	}
}

func TestTranscriptBoundDualRejectsFixedChallengeOffsets(t *testing.T) {
	pp, _ := replayParameters(t)
	ctx := BallotContext{ElectionID: "A", VoterIndex: 1}
	ballot := contextBallot(t, pp, ctx, []int64{1, 2})
	binding, _ := ballotBinding(pp, ctx)
	coefficients := boundDualCoefficients(pp, ballot, binding)
	points := make([]*big.Int, 0, len(pp.D)+len(pp.N))
	for _, d := range pp.D {
		points = append(points, scalarFromInt64(d))
	}
	for _, i := range pp.N {
		points = append(points, big.NewInt(int64(i)))
	}
	weights := make([]*big.Int, 2)
	for i := range weights {
		denominator := big.NewInt(1)
		for j := range points {
			if i != j {
				denominator = mulMod(denominator, subMod(points[i], points[j]))
			}
		}
		weights[i] = mulMod(evaluatePolynomial(coefficients, points[i]), new(big.Int).ModInverse(denominator, bn256.Order))
	}
	if weights[0].Sign() == 0 && weights[1].Sign() == 0 {
		t.Fatal("degenerate zero packed weights")
	}
	// The reviewer's offsets satisfy y_0*delta_0+y_1*delta_1=0.
	ballot.V[0] = addG1(ballot.V[0], new(bn256.G1).ScalarMult(pp.PP1.H0, weights[1]))
	ballot.V[1] = addG1(ballot.V[1], new(bn256.G1).ScalarMult(pp.PP1.H0, negMod(weights[0])))
	if !verifyDualRS(pp, ballot.V, coefficients) {
		t.Fatal("fixed-challenge offset construction did not preserve the old check")
	}
	if verifyBoundDualRS(pp, ballot, binding) || DVerifyWithContext(pp, ballot, ctx) {
		t.Fatal("accepted packed offsets after recomputing the transcript-bound challenge")
	}
}
