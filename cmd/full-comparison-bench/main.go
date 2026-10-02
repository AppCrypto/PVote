// full-comparison-bench executes the implemented protocol phases on one BN254
// backend. Setup and key generation are prepared before timed phases.
package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"runtime"
	"runtime/debug"
	"strconv"
	"time"

	huang "PVote/baselines/huang"
	pri "PVote/baselines/priscore"
	three "PVote/baselines/threepvss"
	rbpvss "PVote/crypto/RBPVSS"
)

type result struct {
	experiment               string
	x                        int
	scheme, phase            string
	meanMS                   float64
	runs, n, l, t, budget, m int
	note                     string
}

func main() {
	out := flag.String("out", "experiments/full_protocol_comparison.csv", "CSV output")
	runs := flag.Int("runs", 5, "timed repetitions")
	scoreDomainOnly := flag.Bool("score-domain-only", false, "run only score-domain generation, verification, and size")
	onlyKappa := flag.Int("only-kappa", 0, "run only this score width (0 selects all)")
	flag.Parse()
	if *runs < 1 {
		panic("runs must be positive")
	}
	runtime.GOMAXPROCS(1)
	old := debug.SetGCPercent(-1)
	defer debug.SetGCPercent(old)
	const n, l, t = 10, 5, 7
	var rows []result
	for _, kappa := range []int{1, 2, 3, 4, 5, 6, 7, 8} {
		if *onlyKappa != 0 && kappa != *onlyKappa {
			continue
		}
		rows = append(rows, scoreRows(*runs, n, l, t, (1<<kappa)-1, kappa)...)
	}
	if !*scoreDomainOnly {
		for _, m := range []int{10, 25, 50, 100} {
			r, e := finalizationRows(*runs, n, l, t, 10, m)
			if e != nil {
				panic(e)
			}
			rows = append(rows, r...)
			rows = append(rows, structuralRows(n, l, t, 10, m)...)
		}
		for _, m := range []int{100, 250, 500, 1000} {
			r, e := tallyScalingRows(*runs, n, l, t, 10, m)
			if e != nil {
				panic(e)
			}
			rows = append(rows, r...)
			if m > 100 {
				rows = append(rows, structuralRows(n, l, t, 10, m)...)
			}
		}
	}
	if e := write(*out, rows); e != nil {
		panic(e)
	}
	for _, r := range rows {
		fmt.Printf("%-24s x=%-3d %-9s %-18s %.3f\n", r.experiment, r.x, r.scheme, r.phase, r.meanMS)
	}
}

func tallyScalingRows(runs, n, l, t, budget, m int) ([]result, error) {
	scores := distributed(l, budget)
	scores64 := make([]int64, l)
	for i := range scores {
		scores64[i] = int64(scores[i])
	}
	pp, sks, err := rbpvss.Setup(128, n, t, l, rbpvss.Interval{Min: 0, Max: int64(budget)})
	if err != nil {
		return nil, err
	}
	agg, err := rbpvss.NewAggregate(pp)
	if err != nil {
		return nil, err
	}
	for i := 0; i < m; i++ {
		inst, e := rbpvss.Share(pp, scores64)
		if e != nil {
			return nil, e
		}
		if e = rbpvss.AggregateVerified(pp, agg, inst); e != nil {
			return nil, e
		}
	}
	pv := mean(runs, func() {
		shares := make([]rbpvss.IndexedShare, t)
		for i := 0; i < t; i++ {
			s, proof, e := rbpvss.Decrypt(pp, agg.C[i], sks[i])
			must(e)
			if !rbpvss.PVerify(pp, i+1, s, agg.C[i], proof) {
				panic("PVerify")
			}
			shares[i] = rbpvss.IndexedShare{Index: i + 1, Share: s}
		}
		_, e := rbpvss.Recon(pp, shares, agg.U, rbpvss.Interval{Min: 0, Max: int64(m * budget)})
		must(e)
	})
	hpp, err := huang.Setup(l, budget)
	if err != nil {
		return nil, err
	}
	ballots := make([]huang.Ballot, m)
	for i := range ballots {
		ballots[i], err = huang.Encrypt(hpp, scores)
		if err != nil {
			return nil, err
		}
	}
	ht := mean(runs, func() { _, e := huang.Tally(hpp, ballots); must(e) })
	note := "tally-path scaling on one BN254 backend with BSGS bounded decoding; Huang scans all m*l ciphertext pairs, PVote consumes t aggregate shares and l aggregate commitments"
	return []result{{"huang_tally_scaling", m, "PVote", "aggregate finalization", pv, runs, n, l, t, budget, m, note}, {"huang_tally_scaling", m, "Huang et al.", "ciphertext scan", ht, runs, n, l, t, budget, m, note}}, nil
}

func structuralRows(n, l, t, budget, m int) []result {
	f := t - l
	tp, _ := three.Setup(n, f, l)
	kappa := 4 // The common score domain is [0,15]=[0,2^kappa-1].
	workingSetBudget := (1 << kappa) - 1
	p, _ := pri.Setup(m, l, workingSetBudget)
	workingSetNote := "fixed-width cryptographic input read by finalization; common domain [0,15] with B=15 and kappa=4; PVote counts only its active aggregate, not immutable transaction history; participation counts distinct required publishers"
	participantNote := "fixed-width cryptographic input read by finalization; PVote counts only its active aggregate, not immutable transaction history; 3PVSS-bit uses kappa=4; participation counts distinct required publishers"
	return []result{
		{"finalization_working_set_kib", m, "PVote", "active aggregate", float64((n+l)*64) / 1024, 1, n, l, t, workingSetBudget, m, workingSetNote},
		{"finalization_working_set_kib", m, "PriScore", "commit+ballot", float64(m*pri.PublicTranscriptBytes(p)) / 1024, 1, n, l, t, workingSetBudget, m, workingSetNote},
		{"finalization_working_set_kib", m, "Huang et al.", "tally ciphertexts", float64(m*l*2*64) / 1024, 1, n, l, t, workingSetBudget, m, workingSetNote},
		{"finalization_working_set_kib", m, "3PVSS-bit", "posted bit planes", float64(m*three.PublicBallotBytes(tp, kappa)) / 1024, 1, n, l, t, workingSetBudget, m, workingSetNote},
		{"post_deadline_participants_normal", m, "PVote", "talliers", float64(t), 1, n, l, t, budget, m, participantNote},
		{"post_deadline_participants_normal", m, "PriScore", "normal self-tally", 0, 1, n, l, t, budget, m, participantNote},
		{"post_deadline_participants_normal", m, "3PVSS", "qualified talliers", float64(f + l), 1, n, l, t, budget, m, participantNote},
		{"post_deadline_participants_abort", m, "PVote", "talliers", float64(t), 1, n, l, t, budget, m, participantNote},
		{"post_deadline_participants_abort", m, "PriScore", "remaining voters", float64(m - 1), 1, n, l, t, budget, m, participantNote},
		{"post_deadline_participants_abort", m, "3PVSS", "qualified talliers", float64(f + l), 1, n, l, t, budget, m, participantNote},
	}
}

func scoreRows(runs, n, l, t, budget, kappa int) []result {
	pp, _, e := rbpvss.Setup(128, n, t, l, rbpvss.Interval{Min: 0, Max: int64(budget)})
	must(e)
	scores64 := make([]int64, l)
	for i, v := range distributed(l, budget) {
		scores64[i] = int64(v)
	}
	inst, e := rbpvss.Share(pp, scores64)
	must(e)
	if !rbpvss.DVerify(pp, inst) {
		panic("PVote transcript rejected")
	}
	pv := mean(runs, func() { _, e := rbpvss.Share(pp, scores64); must(e) })
	pvv := mean(runs, func() {
		if !rbpvss.DVerify(pp, inst) {
			panic("PVote verification")
		}
	})
	priPP, e := pri.Setup(3, l, budget)
	must(e)
	keys := make([]*pri.VoterKey, 3)
	for i := range keys {
		keys[i], e = pri.KeyGen(priPP)
		must(e)
		for j := 0; j < l; j++ {
			_, e = pri.KeyDerive(priPP, keys[i], j)
			must(e)
		}
	}
	scores := distributed(l, budget)
	commit, e := pri.Commit(priPP, 0, keys[0], keys, scores)
	must(e)
	ballot, e := pri.Vote(priPP, 0, keys[0], keys, commit, scores)
	must(e)
	prv := mean(runs, func() {
		c, e := pri.Commit(priPP, 0, keys[0], keys, scores)
		must(e)
		_, e = pri.Vote(priPP, 0, keys[0], keys, c, scores)
		must(e)
	})
	prvv := mean(runs, func() {
		if !pri.VerifyVote(priPP, 0, keys, commit, ballot) {
			panic("PriScore verification")
		}
	})
	f := t - l
	threePP, e := three.Setup(n, f, l)
	must(e)
	threeBallot, e := three.ShareScore(threePP, scores, kappa)
	must(e)
	if !three.VerifyScore(threePP, threeBallot) {
		panic("3PVSS transcript rejected")
	}
	threeVote := mean(runs, func() { _, e := three.ShareScore(threePP, scores, kappa); must(e) })
	threeVerify := mean(runs, func() {
		if !three.VerifyScore(threePP, threeBallot) {
			panic("3PVSS verification")
		}
	})
	pvoteBytes := (6*l+4*n)*64 + (4*l+n+1)*32
	note := "one BN254 backend; B=2^kappa-1; PriScore ballots satisfy sum scores=B; 3PVSS-bit runs kappa native binary instances"
	return []result{
		{"score_domain_generation", kappa, "PVote", "Share", pv, runs, n, l, t, budget, 1, note},
		{"score_domain_generation", kappa, "PriScore", "Commit+Vote", prv, runs, n, l, t, budget, 1, note},
		{"score_domain_generation", kappa, "3PVSS-bit", "kappa Share", threeVote, runs, n, l, t, budget, 1, note},
		{"score_domain_verification", kappa, "PVote", "DVerify", pvv, runs, n, l, t, budget, 1, note},
		{"score_domain_verification", kappa, "PriScore", "VerifyVote", prvv, runs, n, l, t, budget, 1, note},
		{"score_domain_verification", kappa, "3PVSS-bit", "kappa Verify", threeVerify, runs, n, l, t, budget, 1, note},
		{"score_domain_ballot_kib", kappa, "PVote", "serialized ballot", float64(pvoteBytes) / 1024, 1, n, l, t, budget, 1, note},
		{"score_domain_ballot_kib", kappa, "PriScore", "serialized commit+ballot", float64(pri.PublicTranscriptBytes(priPP)) / 1024, 1, n, l, t, budget, 1, note},
		{"score_domain_ballot_kib", kappa, "3PVSS-bit", "paper-accounted ballot", float64(three.PublicBallotBytes(threePP, kappa)) / 1024, 1, n, l, t, budget, 1, note},
	}
}

func finalizationRows(runs, n, l, t, budget, m int) ([]result, error) {
	scores := distributed(l, budget)
	scores64 := make([]int64, l)
	for i, p := range scores {
		scores64[i] = int64(p)
	}
	pp, sks, e := rbpvss.Setup(128, n, t, l, rbpvss.Interval{Min: 0, Max: int64(budget)})
	if e != nil {
		return nil, e
	}
	agg, e := rbpvss.NewAggregate(pp)
	if e != nil {
		return nil, e
	}
	for i := 0; i < m; i++ {
		inst, e := rbpvss.Share(pp, scores64)
		if e != nil {
			return nil, e
		}
		if e = rbpvss.AggregateVerified(pp, agg, inst); e != nil {
			return nil, e
		}
	}
	pv := mean(runs, func() {
		shares := make([]rbpvss.IndexedShare, t)
		for i := 0; i < t; i++ {
			s, p, e := rbpvss.Decrypt(pp, agg.C[i], sks[i])
			must(e)
			if !rbpvss.PVerify(pp, i+1, s, agg.C[i], p) {
				panic("PVerify")
			}
			shares[i] = rbpvss.IndexedShare{Index: i + 1, Share: s}
		}
		_, e := rbpvss.Recon(pp, shares, agg.U, rbpvss.Interval{Min: 0, Max: int64(m * budget)})
		must(e)
	})
	priPP, e := pri.Setup(m, l, budget)
	if e != nil {
		return nil, e
	}
	keys := make([]*pri.VoterKey, m)
	commits := make([]*pri.Commitment, m)
	ballots := make([]*pri.Ballot, m)
	for i := 0; i < m; i++ {
		keys[i], e = pri.KeyGen(priPP)
		if e != nil {
			return nil, e
		}
		for j := 0; j < l; j++ {
			if _, e = pri.KeyDerive(priPP, keys[i], j); e != nil {
				return nil, e
			}
		}
	}
	for i := 0; i < m; i++ {
		commits[i], e = pri.Commit(priPP, i, keys[i], keys, scores)
		if e != nil {
			return nil, e
		}
		ballots[i], e = pri.Vote(priPP, i, keys[i], keys, commits[i], scores)
		if e != nil {
			return nil, e
		}
	}
	pr := mean(runs, func() { _, e := pri.SelfTally(priPP, ballots); must(e) })
	hPP, e := huang.Setup(l, budget)
	if e != nil {
		return nil, e
	}
	hBallots := make([]huang.Ballot, m)
	for i := range hBallots {
		hBallots[i], e = huang.Encrypt(hPP, scores)
		if e != nil {
			return nil, e
		}
	}
	ht := mean(runs, func() { _, e := huang.Tally(hPP, hBallots); must(e) })
	note := "implemented post-deadline phases for PVote and PriScore; Huang row is its tally-relevant ciphertext scan and election-key decryption"
	return []result{{"post_deadline_vs_voters", m, "PVote", "Decrypt+Recon", pv, runs, n, l, t, budget, m, note}, {"post_deadline_vs_voters", m, "PriScore", "Self-Tallying", pr, runs, n, l, t, budget, m, note}, {"post_deadline_vs_voters", m, "Huang et al.", "Tally scan", ht, runs, n, l, t, budget, m, note}}, nil
}

func distributed(l, budget int) []int {
	s := make([]int, l)
	for i := 0; i < budget; i++ {
		s[i%l]++
	}
	return s
}
func mean(runs int, f func()) float64 {
	for i := 0; i < 1; i++ {
		f()
	}
	start := time.Now()
	for i := 0; i < runs; i++ {
		f()
	}
	return float64(time.Since(start).Nanoseconds()) / float64(runs) / 1e6
}
func must(e error) {
	if e != nil {
		panic(e)
	}
}
func write(path string, rows []result) error {
	f, e := os.Create(path)
	if e != nil {
		return e
	}
	defer f.Close()
	w := csv.NewWriter(f)
	defer w.Flush()
	_ = w.Write([]string{"experiment", "x", "scheme", "phase", "mean_ms", "runs", "n", "l", "t", "budget_B", "voters_m", "go_version", "goos", "goarch", "note"})
	for _, r := range rows {
		_ = w.Write([]string{r.experiment, strconv.Itoa(r.x), r.scheme, r.phase, strconv.FormatFloat(r.meanMS, 'f', 6, 64), strconv.Itoa(r.runs), strconv.Itoa(r.n), strconv.Itoa(r.l), strconv.Itoa(r.t), strconv.Itoa(r.budget), strconv.Itoa(r.m), runtime.Version(), runtime.GOOS, runtime.GOARCH, r.note})
	}
	return w.Error()
}
