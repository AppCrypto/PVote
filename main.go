package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand/v2"

	rbpvss "PVote/crypto/RBPVSS"
)

// The CLI is a direct executable trace of the six RB-PVSS interfaces in the
// paper. Blockchain settlement is demonstrated separately by ./web; it is not
// part of the RB-PVSS primitive.
func main() {
	n := flag.Int("n", 7, "number of talliers/shareholders")
	l := flag.Int("l", 2, "number of ballot coordinates")
	t := flag.Int("t", 0, "reconstruction threshold (default floor((n+l)/2))")
	voters := flag.Int("m", 3, "number of accepted demo ballots")
	minScore := flag.Int64("a", 0, "minimum per-coordinate score")
	maxScore := flag.Int64("b", 5, "maximum per-coordinate score")
	flag.Parse()

	threshold := *t
	if threshold == 0 {
		threshold = (*n + *l) / 2
	}
	domain := rbpvss.Interval{Min: *minScore, Max: *maxScore}

	// (pp,{sk_i}) <- RB-PVSS.Setup(1^lambda,n,t,l,S)
	pp, secretKeys, err := rbpvss.Setup(128, *n, threshold, *l, domain)
	if err != nil {
		log.Fatal(err)
	}
	aggregate, err := rbpvss.NewAggregate(pp)
	if err != nil {
		log.Fatal(err)
	}

	expected := make([]int64, *l)
	for voter := 1; voter <= *voters; voter++ {
		vector := randomVector(*l, domain)
		for d, score := range vector {
			expected[d] += score
		}

		// I_j <- RB-PVSS.Share(pp,w_j), followed by atomic DVerify and
		// aggregation of only an accepted complete transcript.
		instance, err := rbpvss.Share(pp, vector)
		if err != nil {
			log.Fatalf("voter %d Share: %v", voter, err)
		}
		if err := rbpvss.AggregateVerified(pp, aggregate, instance); err != nil {
			log.Fatalf("voter %d aggregate: %v", voter, err)
		}
		fmt.Printf("accepted voter %d: %v\n", voter, vector)
	}

	// Each selected T_i invokes Decrypt; the public verifier invokes PVerify.
	shareSet := make([]rbpvss.IndexedShare, 0, threshold)
	for i := 1; i <= threshold; i++ {
		share, proof, err := rbpvss.Decrypt(pp, aggregate.C[i-1], secretKeys[i-1])
		if err != nil {
			log.Fatalf("tallier %d Decrypt: %v", i, err)
		}
		if !rbpvss.PVerify(pp, i, share, aggregate.C[i-1], proof) {
			log.Fatalf("tallier %d PVerify rejected an honest share", i)
		}
		shareSet = append(shareSet, rbpvss.IndexedShare{Index: i, Share: share})
	}

	// W <- RB-PVSS.Recon(pp,{(i,sh_i)},{U*_d},S_A).
	decodeDomain := rbpvss.Interval{
		Min: int64(aggregate.Count) * domain.Min,
		Max: int64(aggregate.Count) * domain.Max,
	}
	result, err := rbpvss.Recon(pp, shareSet, aggregate.U, decodeDomain)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("reconstructed aggregate: %v\n", result)
	fmt.Printf("expected aggregate:      %v\n", expected)
}

func randomVector(length int, domain rbpvss.Interval) []int64 {
	result := make([]int64, length)
	width := uint64(domain.Max-domain.Min) + 1
	for i := range result {
		result[i] = domain.Min + int64(rand.Uint64N(width))
	}
	return result
}
