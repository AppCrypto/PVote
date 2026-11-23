package mobilecore

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"reflect"
	"runtime"
	"time"

	pri "PVote/baselines/priscore"
	three "PVote/baselines/threepvss"
	rb "PVote/crypto/RBPVSS"
	bn "github.com/ethereum/go-ethereum/crypto/bn256/cloudflare"
)

type comparisonState struct {
	rb           *rb.PublicParameters
	pri          *pri.Parameters
	keys         []*pri.VoterKey
	three        *three.Parameters
	scores       []int
	fingerprints map[string]string
}

// Comparison setup and key derivation are outside all generation timers. The
// three runs reuse these public parameters, as in the main device benchmark.
func (e *Engine) PrepareComparison(kappa int) (map[string]any, error) {
	e.mu.Lock()
	defer e.mu.Unlock()
	if kappa < 1 || kappa > 8 {
		return nil, fmt.Errorf("comparison kappa must be in [1,8]")
	}
	if e.comparisons == nil {
		e.comparisons = make(map[int]*comparisonState)
	}
	if e.comparisons[kappa] == nil {
		budget := (1 << kappa) - 1
		s := &comparisonState{scores: make([]int, 5), fingerprints: make(map[string]string)}
		for i := 0; i < budget; i++ {
			s.scores[i%5]++
		}
		var err error
		s.rb, _, err = rb.Setup(128, 10, 7, 5, rb.Interval{Min: 0, Max: int64(budget)})
		if err != nil {
			return nil, err
		}
		s.pri, err = pri.Setup(3, 5, budget)
		if err != nil {
			return nil, err
		}
		s.keys = make([]*pri.VoterKey, 3)
		publicKeys := make([]*bn.G1, 3)
		derived := make([][]*bn.G1, 3)
		for i := range s.keys {
			s.keys[i], err = pri.KeyGen(s.pri)
			if err != nil {
				return nil, err
			}
			for d := 0; d < 5; d++ {
				if _, err = pri.KeyDerive(s.pri, s.keys[i], d); err != nil {
					return nil, err
				}
			}
			publicKeys[i] = s.keys[i].Y
			derived[i] = s.keys[i].DerivedY
		}
		s.three, err = three.Setup(10, 2, 5)
		if err != nil {
			return nil, err
		}
		public := map[string]any{
			"PVote": Export(s.rb),
			"PriScore": struct {
				PP         *pri.Parameters
				PublicKeys []*bn.G1
				Derived    [][]*bn.G1
			}{s.pri, publicKeys, derived},
			"3PVSS-bit": struct {
				N, F, L, R int
				G, H       *bn.G1
				PK         []*bn.G1
			}{10, 2, 5, 7, s.three.G, s.three.H, s.three.PK},
		}
		for scheme, value := range public {
			fingerprint, err := comparisonHash(value)
			if err != nil {
				return nil, err
			}
			s.fingerprints[scheme] = fingerprint
		}
		e.comparisons[kappa] = s
	}
	return map[string]any{"kappa": kappa, "n": 10, "l": 5, "t": 7, "budget": (1 << kappa) - 1, "priVoters": 3, "scores": e.comparisons[kappa].scores, "fingerprints": e.comparisons[kappa].fingerprints}, nil
}

func (e *Engine) GenerateComparison(r Request) (map[string]any, error) {
	e.mu.Lock()
	defer e.mu.Unlock()
	s := e.comparisons[r.Kappa]
	if s == nil {
		return nil, fmt.Errorf("prepare comparison parameters first")
	}
	if _, ok := s.fingerprints[r.Scheme]; !ok {
		return nil, fmt.Errorf("unknown comparison scheme")
	}
	if r.Scheme == "PVote" && (r.ElectionID == "" || r.VoterIndex == 0) {
		return nil, fmt.Errorf("context required")
	}
	scores64 := make([]int64, len(s.scores))
	for i, v := range s.scores {
		scores64[i] = int64(v)
	}
	var instance *rb.Instance
	var commit *pri.Commitment
	var ballot *pri.Ballot
	var bitBallot *three.ScoreBallot
	var err error
	var before, after runtime.MemStats
	runtime.ReadMemStats(&before)
	cpu := processCPU()
	start := time.Now()
	phase := "Share"
	switch r.Scheme {
	case "PVote":
		instance, err = rb.ShareWithContext(s.rb, scores64, rb.BallotContext{ElectionID: r.ElectionID, VoterIndex: r.VoterIndex})
	case "PriScore":
		phase = "Commit+Vote"
		commit, err = pri.Commit(s.pri, 0, s.keys[0], s.keys, s.scores)
		if err == nil {
			ballot, err = pri.Vote(s.pri, 0, s.keys[0], s.keys, commit, s.scores)
		}
	case "3PVSS-bit":
		phase = "kappa Share"
		bitBallot, err = three.ShareScore(s.three, s.scores, r.Kappa)
	}
	elapsed := time.Since(start)
	cpuElapsed := processCPU() - cpu
	runtime.ReadMemStats(&after)
	if err != nil {
		return nil, err
	}
	start = time.Now()
	var public any
	size := 0
	switch r.Scheme {
	case "PVote":
		public = instance
		size = (6*5+4*10)*64 + (4*5+10+1)*32
	case "PriScore":
		public = struct {
			Commit *pri.Commitment
			Vote   *pri.Ballot
		}{commit, ballot}
		size = pri.PublicTranscriptBytes(s.pri)
	case "3PVSS-bit":
		public = bitBallot
		size = three.PublicBallotBytes(s.three, r.Kappa)
	}
	digest, err := comparisonHash(public)
	if err != nil {
		return nil, err
	}
	serialization := time.Since(start)
	result := map[string]any{"scheme": r.Scheme, "phase": phase, "kappa": r.Kappa, "budget": (1 << r.Kappa) - 1, "n": 10, "l": 5, "t": 7, "priVoters": 3,
		"generationMs": millis(elapsed), "generationCpuMs": float64(cpuElapsed) / 1e6, "serializationMs": millis(serialization),
		"parameterFingerprint": s.fingerprints[r.Scheme], "transcriptHash": digest, "transcriptBytes": size,
		"generationAllocatedBytes": after.TotalAlloc - before.TotalAlloc, "generationGCs": after.NumGC - before.NumGC}
	if r.Verify {
		start = time.Now()
		valid := false
		switch r.Scheme {
		case "PVote":
			valid = rb.DVerifyWithContext(s.rb, instance, rb.BallotContext{ElectionID: r.ElectionID, VoterIndex: r.VoterIndex})
		case "PriScore":
			valid = pri.VerifyVote(s.pri, 0, s.keys, commit, ballot)
		case "3PVSS-bit":
			valid = three.VerifyScore(s.three, bitBallot)
		}
		result["verifyMs"] = millis(time.Since(start))
		if !valid {
			return nil, fmt.Errorf("comparison transcript rejected")
		}
	}
	return result, nil
}

// Public point/scalar encodings avoid encoding opaque curve objects as {}.
// Unexported commitment openings are excluded; setup fingerprints use only the
// explicitly selected public fields above. No secret keys leave this process.
func comparisonHash(value any) (string, error) {
	data, err := json.Marshal(comparisonPublic(reflect.ValueOf(value)))
	if err != nil {
		return "", err
	}
	h := sha256.Sum256(data)
	return hex.EncodeToString(h[:]), nil
}
func comparisonPublic(v reflect.Value) any {
	if !v.IsValid() {
		return nil
	}
	if v.Kind() == reflect.Interface {
		return comparisonPublic(v.Elem())
	}
	if v.Kind() == reflect.Pointer {
		if v.IsNil() {
			return nil
		}
		switch x := v.Interface().(type) {
		case *bn.G1:
			return hex.EncodeToString(x.Marshal())
		case *bn.G2:
			return hex.EncodeToString(x.Marshal())
		case *bn.GT:
			return hex.EncodeToString(x.Marshal())
		case *big.Int:
			return x.String()
		}
		return comparisonPublic(v.Elem())
	}
	switch v.Kind() {
	case reflect.Struct:
		out := map[string]any{}
		for i := 0; i < v.NumField(); i++ {
			f := v.Type().Field(i)
			if f.PkgPath == "" {
				out[f.Name] = comparisonPublic(v.Field(i))
			}
		}
		return out
	case reflect.Slice, reflect.Array:
		out := make([]any, v.Len())
		for i := range out {
			out[i] = comparisonPublic(v.Index(i))
		}
		return out
	default:
		return v.Interface()
	}
}
