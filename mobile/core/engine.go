package mobilecore

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"runtime"
	"sync"
	"syscall"
	"time"

	encoding "PVote/chain/rbpvss"
	contract "PVote/compile/contract"
	rb "PVote/crypto/RBPVSS"
)

// Revision is set by the reproducible native build script.
var Revision = "development"
var SourceSHA = "development"

type Engine struct {
	mu          sync.Mutex
	pp          *rb.PublicParameters
	fingerprint string
	comparisons map[int]*comparisonState
}

type Request struct {
	Action          string       `json:"action"`
	Scheme          string       `json:"scheme,omitempty"`
	Kappa           int          `json:"kappa,omitempty"`
	Parameters      *Parameters  `json:"parameters,omitempty"`
	Scores          []int64      `json:"scores,omitempty"`
	Protocol        string       `json:"protocol,omitempty"`
	ElectionID      string       `json:"electionId,omitempty"`
	VoterIndex      uint64       `json:"voterIndex,omitempty"`
	Verify          bool         `json:"verify,omitempty"`
	IncludeCalldata bool         `json:"includeCalldata,omitempty"`
	PrivateKey      string       `json:"privateKey,omitempty"`
	Transaction     *SignRequest `json:"transaction,omitempty"`
}

type Measurement struct {
	Protocol            string   `json:"protocol"`
	Fingerprint         string   `json:"parameterFingerprint"`
	N                   int      `json:"n"`
	T                   int      `json:"t"`
	L                   int      `json:"l"`
	ShareMS             float64  `json:"shareMs"`
	ShareCPUMS          float64  `json:"shareCpuMs"`
	EncodeMS            float64  `json:"encodeMs"`
	VerifyMS            *float64 `json:"verifyMs,omitempty"`
	TranscriptBytes     int      `json:"transcriptBytes"`
	CalldataBytes       int      `json:"calldataBytes"`
	ShareAllocatedBytes uint64   `json:"shareAllocatedBytes"`
	ShareMallocs        uint64   `json:"shareMallocs"`
	ShareGCs            uint32   `json:"shareGCs"`
	GoHeapBefore        uint64   `json:"goHeapBeforeBytes"`
	GoHeapAfter         uint64   `json:"goHeapAfterBytes"`
	Calldata            string   `json:"calldata,omitempty"`
	TranscriptHash      string   `json:"transcriptHash"`
	AcceptanceTopic     string   `json:"acceptanceTopic"`
}

func millis(d time.Duration) float64 { return float64(d) / float64(time.Millisecond) }
func processCPU() int64 {
	var r syscall.Rusage
	if syscall.Getrusage(syscall.RUSAGE_SELF, &r) != nil {
		return 0
	}
	return r.Utime.Nano() + r.Stime.Nano()
}

func (e *Engine) Prepare(p Parameters) (map[string]any, error) {
	pp, err := p.Decode()
	if err != nil {
		return nil, err
	}
	e.mu.Lock()
	defer e.mu.Unlock()
	e.pp = pp
	e.fingerprint = p.Fingerprint()
	return map[string]any{"parameterFingerprint": e.fingerprint, "n": p.N, "t": p.T, "l": p.L}, nil
}

func (e *Engine) Generate(r Request) (*Measurement, error) {
	e.mu.Lock()
	defer e.mu.Unlock()
	pp := e.pp
	if pp == nil {
		return nil, fmt.Errorf("load public parameters first")
	}
	if len(r.Scores) != len(pp.D) {
		return nil, fmt.Errorf("expected %d scores", len(pp.D))
	}
	for _, w := range r.Scores {
		if w < pp.Range.Min || w > pp.Range.Max {
			return nil, fmt.Errorf("score outside public interval")
		}
	}
	if r.Protocol != Legacy && r.Protocol != Bound {
		return nil, fmt.Errorf("unsupported proof protocol")
	}
	ctx := rb.BallotContext{ElectionID: r.ElectionID, VoterIndex: r.VoterIndex}
	if r.Protocol == Bound && (ctx.ElectionID == "" || ctx.VoterIndex == 0) {
		return nil, fmt.Errorf("context-bound proofs require election ID and voter index")
	}
	var before, after runtime.MemStats
	runtime.ReadMemStats(&before)
	cpu := processCPU()
	start := time.Now()
	var instance *rb.Instance
	var err error
	if r.Protocol == Bound {
		instance, err = rb.ShareWithContext(pp, r.Scores, ctx)
	} else {
		instance, err = rb.Share(pp, r.Scores)
	}
	elapsed := time.Since(start)
	cpuElapsed := processCPU() - cpu
	runtime.ReadMemStats(&after)
	if err != nil {
		return nil, err
	}
	start = time.Now()
	dealer, ranges, err := encoding.Transcript(instance)
	if err != nil {
		return nil, err
	}
	abi, err := contract.ContractMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	data, err := abi.Pack("submitRB", dealer, ranges)
	if err != nil {
		return nil, err
	}
	encodingTime := time.Since(start)
	hash := sha256.Sum256(data)
	m := &Measurement{Protocol: r.Protocol, Fingerprint: e.fingerprint, N: len(pp.N), T: pp.Threshold, L: len(pp.D),
		ShareMS: millis(elapsed), ShareCPUMS: float64(cpuElapsed) / 1e6, EncodeMS: millis(encodingTime),
		TranscriptBytes: (6*len(pp.D)+4*len(pp.N))*64 + (4*len(pp.D)+len(pp.N)+1)*32,
		CalldataBytes:   len(data), ShareAllocatedBytes: after.TotalAlloc - before.TotalAlloc, ShareMallocs: after.Mallocs - before.Mallocs,
		ShareGCs: after.NumGC - before.NumGC, GoHeapBefore: before.HeapAlloc, GoHeapAfter: after.HeapAlloc,
		TranscriptHash: hex.EncodeToString(hash[:]), AcceptanceTopic: abi.Events["TranscriptAccepted"].ID.Hex()}
	if r.IncludeCalldata {
		m.Calldata = "0x" + hex.EncodeToString(data)
	}
	if r.Verify {
		start = time.Now()
		valid := false
		if r.Protocol == Bound {
			valid = rb.DVerifyWithContext(pp, instance, ctx)
		} else {
			valid = rb.DVerify(pp, instance)
		}
		duration := millis(time.Since(start))
		m.VerifyMS = &duration
		if !valid {
			return nil, fmt.Errorf("locally generated ballot failed verification")
		}
	}
	return m, nil
}

// Call serializes errors, never private request data, across the C boundary.
func (e *Engine) Call(data []byte) (result []byte) {
	defer func() {
		if recover() != nil {
			result = []byte(`{"error":"native operation failed"}`)
		}
	}()
	var r Request
	if len(data) > 1024*1024 {
		return []byte(`{"error":"request too large"}`)
	}
	if err := json.Unmarshal(data, &r); err != nil {
		return []byte(`{"error":"invalid request JSON"}`)
	}
	var value any
	var err error
	switch r.Action {
	case "version":
		value = map[string]any{"revision": Revision, "sourceSHA256": SourceSHA, "goVersion": runtime.Version(), "goOS": runtime.GOOS, "goArch": runtime.GOARCH, "gomaxprocs": runtime.GOMAXPROCS(0), "gcPolicy": "Go default (enabled)", "protocols": []string{Legacy, Bound}}
	case "prepare":
		if r.Parameters == nil {
			err = fmt.Errorf("parameters missing")
		} else {
			value, err = e.Prepare(*r.Parameters)
		}
	case "generate":
		value, err = e.Generate(r)
	case "prepareComparison":
		value, err = e.PrepareComparison(r.Kappa)
	case "generateComparison":
		value, err = e.GenerateComparison(r)
	case "wallet":
		value, err = NewWallet()
	case "address":
		value, err = WalletAddress(r.PrivateKey)
	case "sign":
		if r.Transaction == nil {
			err = fmt.Errorf("transaction missing")
		} else {
			value, err = Sign(r.PrivateKey, *r.Transaction)
		}
	default:
		err = fmt.Errorf("unknown native action")
	}
	if err != nil {
		result, _ = json.Marshal(map[string]string{"error": err.Error()})
		return result
	}
	result, err = json.Marshal(value)
	if err != nil {
		return []byte(`{"error":"cannot encode response"}`)
	}
	return result
}
