// comparison-bench calibrates the operation counts reported in the paper
// against one common implementation of A0, E0, E1, and P. PriScore rows are
// operation-normalized estimates, not claims of executing the original
// authors' implementation.
package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"math/big"
	"os"
	"runtime"
	"runtime/debug"
	"strconv"
	"time"

	bn256 "github.com/ethereum/go-ethereum/crypto/bn256/cloudflare"
)

type primitiveCosts struct {
	add0NS float64
	exp0NS float64
	exp1NS float64
	pairNS float64
}

type row struct {
	experiment string
	x          int
	n, l, t    int
	budget     int
	batch      int
	pvoteMS    float64
	priscoreMS float64
	note       string
}

func main() {
	out := flag.String("out", "experiments/normalized_comparison_results.csv", "CSV output path")
	runs := flag.Int("runs", 200, "calibration repetitions per primitive")
	flag.Parse()
	if *runs < 20 {
		panic("runs must be at least 20")
	}

	runtime.GOMAXPROCS(1)
	previousGC := debug.SetGCPercent(-1)
	defer debug.SetGCPercent(previousGC)

	costs := calibrate(*runs)
	const n, l, t = 40, 10, 25
	const verifyBatch = 1000
	var rows []row

	for _, budget := range []int{1, 3, 5, 10, 20, 50} {
		rows = append(rows, row{
			experiment: "vote_vs_score_parameter",
			x:          budget,
			n:          n,
			l:          l,
			t:          t,
			budget:     budget,
			batch:      1,
			pvoteMS:    nsToMS(float64(9*l+4*n)*costs.exp0NS + float64(l)*costs.pairNS),
			priscoreMS: nsToMS(float64((14*budget+17)*l+4) * costs.exp0NS),
			note:       "operation-normalized; PriScore uses B=x and PVote uses range width b-a=x",
		})
		rows = append(rows, row{
			experiment: "verify_per_ballot_vs_score_parameter",
			x:          budget,
			n:          n,
			l:          l,
			t:          t,
			budget:     budget,
			batch:      verifyBatch,
			pvoteMS:    nsToMS(float64(5*n+9*l)*costs.exp0NS + float64(4*l)*costs.pairNS),
			priscoreMS: nsToMS(float64(17*budget*l+8) * costs.exp0NS),
			note:       "operation-normalized per ballot",
		})
	}

	for _, voters := range []int{100, 1000, 10000, 100000} {
		rows = append(rows, row{
			experiment: "post_deadline_group_work_vs_voters",
			x:          voters,
			n:          n,
			l:          l,
			t:          t,
			batch:      voters,
			pvoteMS:    nsToMS(float64(7*t+t*l) * costs.exp0NS),
			note:       "operation-normalized group work",
		})
	}

	rows = append(rows, row{
		experiment: "streaming_aggregation_per_ballot",
		x:          n + l,
		n:          n,
		l:          l,
		t:          t,
		batch:      1,
		pvoteMS:    nsToMS(float64(n+l) * costs.add0NS),
		note:       "PVote incremental aggregate update; baselines defer aggregation",
	})

	if err := writeCSV(*out, *runs, costs, rows); err != nil {
		panic(err)
	}
	fmt.Printf("A0=%.0fns E0=%.0fns E1=%.0fns P=%.0fns\n", costs.add0NS, costs.exp0NS, costs.exp1NS, costs.pairNS)
	for _, r := range rows {
		fmt.Printf("%s x=%d PVote=%.3fms PriScore=%.3fms\n", r.experiment, r.x, r.pvoteMS, r.priscoreMS)
	}
}

func calibrate(runs int) primitiveCosts {
	scalar := new(big.Int).Sub(bn256.Order, big.NewInt(17))
	g1a := new(bn256.G1).ScalarBaseMult(big.NewInt(7))
	g1b := new(bn256.G1).ScalarBaseMult(big.NewInt(11))
	g2 := new(bn256.G2).ScalarBaseMult(big.NewInt(13))

	return primitiveCosts{
		add0NS: measureNS(runs, func() { _ = new(bn256.G1).Add(g1a, g1b) }),
		exp0NS: measureNS(runs, func() { _ = new(bn256.G1).ScalarMult(g1a, scalar) }),
		exp1NS: measureNS(runs, func() { _ = new(bn256.G2).ScalarMult(g2, scalar) }),
		pairNS: measureNS(runs, func() { _ = bn256.Pair(g1a, g2) }),
	}
}

func measureNS(runs int, operation func()) float64 {
	for i := 0; i < 10; i++ {
		operation()
	}
	started := time.Now()
	for i := 0; i < runs; i++ {
		operation()
	}
	return float64(time.Since(started).Nanoseconds()) / float64(runs)
}

func nsToMS(value float64) float64 { return value / float64(time.Millisecond) }

func writeCSV(path string, runs int, costs primitiveCosts, rows []row) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	writer := csv.NewWriter(file)
	defer writer.Flush()
	header := []string{
		"experiment", "x", "n", "l", "t", "budget_B", "batch_m",
		"pvote_normalized_ms", "priscore_normalized_ms",
		"calibration_runs", "go_version", "goos", "goarch",
		"a0_mean_ns", "e0_mean_ns", "e1_mean_ns", "pair_mean_ns", "note",
	}
	if err := writer.Write(header); err != nil {
		return err
	}
	for _, r := range rows {
		values := []string{
			r.experiment, strconv.Itoa(r.x), strconv.Itoa(r.n), strconv.Itoa(r.l),
			strconv.Itoa(r.t), strconv.Itoa(r.budget), strconv.Itoa(r.batch),
			decimal(r.pvoteMS), decimal(r.priscoreMS),
			strconv.Itoa(runs), runtime.Version(), runtime.GOOS, runtime.GOARCH,
			decimal(costs.add0NS), decimal(costs.exp0NS), decimal(costs.exp1NS),
			decimal(costs.pairNS), r.note,
		}
		if err := writer.Write(values); err != nil {
			return err
		}
	}
	return writer.Error()
}

func decimal(value float64) string {
	return strconv.FormatFloat(value, 'f', 6, 64)
}
