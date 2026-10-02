package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"math"
	"math/big"
	"os"
	"runtime"
	"runtime/debug"
	"strconv"
	"time"

	rbpvss "PVote/crypto/RBPVSS"
	bn256 "github.com/ethereum/go-ethereum/crypto/bn256/cloudflare"
)

type measurement struct {
	n, l, t       int
	setup         stats
	share         stats
	dverify       stats
	decrypt       stats
	pverify       stats
	recon         stats
	transcriptKiB float64
}

type stats struct {
	meanMS float64
	stdMS  float64
}

func main() {
	out := flag.String("out", "experiments/rbpvss_results.csv", "CSV output path")
	runs := flag.Int("runs", 20, "timed repetitions per operation")
	setupRuns := flag.Int("setup-runs", 5, "timed Setup repetitions per parameter pair")
	onlyN := flag.Int("only-n", 0, "run only this tallier count (0 selects all)")
	onlyL := flag.Int("only-l", 0, "run only this score-vector length (0 selects all)")
	shareOnly := flag.Bool("share-only", false, "measure only transcript generation and size")
	flag.Parse()
	if *runs < 2 || *setupRuns < 2 {
		panic("runs and setup-runs must both be at least 2")
	}

	runtime.GOMAXPROCS(1)
	previousGC := debug.SetGCPercent(-1)
	defer debug.SetGCPercent(previousGC)

	var rows []measurement
	for _, l := range []int{1, 10, 19} {
		if *onlyL != 0 && l != *onlyL {
			continue
		}
		for n := 10; n <= 100; n += 10 {
			if *onlyN != 0 && n != *onlyN {
				continue
			}
			if l > n/3 {
				continue
			}
			row, err := benchmarkPair(n, l, *runs, *setupRuns, *shareOnly)
			if err != nil {
				panic(err)
			}
			rows = append(rows, row)
			fmt.Printf("n=%d l=%d t=%d Share=%.3fms DVerify=%.3fms size=%.3fKiB\n", n, l, row.t, row.share.meanMS, row.dverify.meanMS, row.transcriptKiB)
		}
	}
	if err := writeCSV(*out, rows, *runs, *setupRuns); err != nil {
		panic(err)
	}
}

func benchmarkPair(n, l, runs, setupRuns int, shareOnly bool) (measurement, error) {
	t := (n + l) / 2
	domain := rbpvss.Interval{Min: 0, Max: 10}
	setupSamples := make([]time.Duration, setupRuns)
	var pp *rbpvss.PublicParameters
	var secretKeys []*big.Int
	for i := range setupSamples {
		started := time.Now()
		var err error
		pp, secretKeys, err = rbpvss.Setup(128, n, t, l, domain)
		setupSamples[i] = time.Since(started)
		if err != nil {
			return measurement{}, err
		}
	}

	scores := make([]int64, l)
	for i := range scores {
		scores[i] = 5
	}
	instances := make([]*rbpvss.Instance, runs)
	shareSamples := make([]time.Duration, runs)
	for i := range shareSamples {
		started := time.Now()
		var err error
		instances[i], err = rbpvss.Share(pp, scores)
		shareSamples[i] = time.Since(started)
		if err != nil {
			return measurement{}, err
		}
	}
	if shareOnly {
		return measurement{
			n:             n,
			l:             l,
			t:             t,
			setup:         summarize(setupSamples),
			share:         summarize(shareSamples),
			transcriptKiB: float64(transcriptSize(instances[0])) / 1024,
		}, nil
	}

	dverifySamples := make([]time.Duration, runs)
	for i := range dverifySamples {
		started := time.Now()
		ok := rbpvss.DVerify(pp, instances[i])
		dverifySamples[i] = time.Since(started)
		if !ok {
			return measurement{}, fmt.Errorf("DVerify rejected honest transcript for n=%d,l=%d", n, l)
		}
	}

	aggregate, err := rbpvss.NewAggregate(pp)
	if err != nil {
		return measurement{}, err
	}
	if err := rbpvss.AggregateVerified(pp, aggregate, instances[0]); err != nil {
		return measurement{}, err
	}

	decryptSamples := make([]time.Duration, runs)
	decryptShares := make([]*bn256.G1, runs)
	decryptProofs := make([]*rbpvss.DecryptionProof, runs)
	for i := range decryptSamples {
		started := time.Now()
		decryptShares[i], decryptProofs[i], err = rbpvss.Decrypt(pp, aggregate.C[0], secretKeys[0])
		decryptSamples[i] = time.Since(started)
		if err != nil {
			return measurement{}, err
		}
	}
	pverifySamples := make([]time.Duration, runs)
	for i := range pverifySamples {
		started := time.Now()
		ok := rbpvss.PVerify(pp, 1, decryptShares[i], aggregate.C[0], decryptProofs[i])
		pverifySamples[i] = time.Since(started)
		if !ok {
			return measurement{}, fmt.Errorf("PVerify rejected honest share for n=%d,l=%d", n, l)
		}
	}

	shareSet := make([]rbpvss.IndexedShare, t)
	for i := range shareSet {
		share, proof, err := rbpvss.Decrypt(pp, aggregate.C[i], secretKeys[i])
		if err != nil || !rbpvss.PVerify(pp, i+1, share, aggregate.C[i], proof) {
			return measurement{}, fmt.Errorf("prepare reconstruction share %d: %v", i+1, err)
		}
		shareSet[i] = rbpvss.IndexedShare{Index: i + 1, Share: share}
	}
	reconSamples := make([]time.Duration, runs)
	for i := range reconSamples {
		started := time.Now()
		result, err := rbpvss.Recon(pp, shareSet, aggregate.U, domain)
		reconSamples[i] = time.Since(started)
		if err != nil || len(result) != l {
			return measurement{}, fmt.Errorf("Recon failed for n=%d,l=%d: %v", n, l, err)
		}
	}

	return measurement{
		n:             n,
		l:             l,
		t:             t,
		setup:         summarize(setupSamples),
		share:         summarize(shareSamples),
		dverify:       summarize(dverifySamples),
		decrypt:       summarize(decryptSamples),
		pverify:       summarize(pverifySamples),
		recon:         summarize(reconSamples),
		transcriptKiB: float64(transcriptSize(instances[0])) / 1024,
	}, nil
}

func transcriptSize(instance *rbpvss.Instance) int {
	size := 0
	for _, point := range instance.V {
		size += len(point.Marshal())
	}
	for _, point := range instance.C {
		size += len(point.Marshal())
	}
	size += scalarSize(instance.Phi.Chi)
	for i := range instance.Phi.A {
		size += len(instance.Phi.A[i].Marshal())
		size += len(instance.Phi.B[i].Marshal())
		size += scalarSize(instance.Phi.Z[i])
	}
	for i, proof := range instance.Pi {
		fBase := proof.FEncoding()
		size += len(instance.U[i].Marshal())
		size += len(proof.E.Marshal()) + len(fBase.Marshal())
		size += len(proof.UPrime.Marshal()) + len(proof.CPrime.Marshal())
		size += scalarSize(proof.Chi) + scalarSize(proof.Z1) + scalarSize(proof.ZBeta) + scalarSize(proof.Z3)
	}
	return size
}

func scalarSize(_ *big.Int) int { return 32 }

func summarize(samples []time.Duration) stats {
	values := make([]float64, len(samples))
	mean := 0.0
	for i, sample := range samples {
		values[i] = float64(sample) / float64(time.Millisecond)
		mean += values[i]
	}
	mean /= float64(len(values))
	variance := 0.0
	for _, value := range values {
		delta := value - mean
		variance += delta * delta
	}
	variance /= float64(len(values) - 1)
	return stats{meanMS: mean, stdMS: math.Sqrt(variance)}
}

func writeCSV(path string, rows []measurement, runs, setupRuns int) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	writer := csv.NewWriter(file)
	defer writer.Flush()
	if err := writer.Write([]string{
		"n", "l", "t", "runs", "setup_runs", "go_version", "goos", "goarch",
		"setup_mean_ms", "setup_std_ms", "share_mean_ms", "share_std_ms",
		"dverify_mean_ms", "dverify_std_ms", "transcript_kib",
		"decrypt_mean_ms", "decrypt_std_ms", "pverify_mean_ms", "pverify_std_ms",
		"recon_mean_ms", "recon_std_ms",
	}); err != nil {
		return err
	}
	for _, row := range rows {
		values := []string{
			strconv.Itoa(row.n), strconv.Itoa(row.l), strconv.Itoa(row.t), strconv.Itoa(runs), strconv.Itoa(setupRuns),
			runtime.Version(), runtime.GOOS, runtime.GOARCH,
			decimal(row.setup.meanMS), decimal(row.setup.stdMS),
			decimal(row.share.meanMS), decimal(row.share.stdMS),
			decimal(row.dverify.meanMS), decimal(row.dverify.stdMS),
			decimal(row.transcriptKiB),
			decimal(row.decrypt.meanMS), decimal(row.decrypt.stdMS),
			decimal(row.pverify.meanMS), decimal(row.pverify.stdMS),
			decimal(row.recon.meanMS), decimal(row.recon.stdMS),
		}
		if err := writer.Write(values); err != nil {
			return err
		}
	}
	return writer.Error()
}

func decimal(value float64) string { return strconv.FormatFloat(value, 'f', 6, 64) }
