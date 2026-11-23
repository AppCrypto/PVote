package rbpvss

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"math/big"
	"os"
	"runtime"
	"testing"
	"time"

	bn256 "github.com/ethereum/go-ethereum/crypto/bn256/cloudflare"
)

// Run scripts/run_bounded_dlog.py to measure the existing voting decoder.
func TestBoundedDLogMeasurement(t *testing.T) {
	out := os.Getenv("PVOTE_DLOG_BENCH_OUT")
	if out == "" {
		t.Skip("measurement disabled")
	}
	runtime.GOMAXPROCS(1)
	data, err := os.ReadFile(os.Getenv("PVOTE_DLOG_FIXTURES"))
	if err != nil {
		t.Fatal(err)
	}
	var fixtures []struct {
		ID         string `json:"id"`
		Parameters struct {
			H0 string `json:"h0"`
		} `json:"parameters"`
	}
	if err := json.Unmarshal(data, &fixtures); err != nil {
		t.Fatal(err)
	}
	const fixtureID = "n100-l19-t59"
	var h0Hex string
	for _, fixture := range fixtures {
		if fixture.ID == fixtureID {
			h0Hex = fixture.Parameters.H0
		}
	}
	encoded, err := hex.DecodeString(h0Hex)
	if err != nil || len(encoded) != 64 {
		t.Fatal("voting generator not found")
	}
	h0 := new(bn256.G1)
	if _, err := h0.Unmarshal(encoded); err != nil || equalG1(h0, zeroG1()) {
		t.Fatal("invalid voting generator")
	}
	interval := Interval{Min: 0, Max: 255 * 100000}
	for _, x := range []int64{0, interval.Max / 2, interval.Max} {
		y := new(bn256.G1).ScalarMult(h0, big.NewInt(x))
		if got, ok := boundedDLog(h0, y, interval); !ok || got != x {
			t.Fatal("warmup recovery failed")
		}
	}
	type sample struct {
		Index    int    `json:"index"`
		Exponent int64  `json:"exponent"`
		Target   string `json:"targetHex"`
		DecodeNS int64  `json:"decodeNS"`
		Verified bool   `json:"verified"`
	}
	rows := make([]sample, 0, 10)
	seen := make(map[int64]bool)
	for len(rows) < 10 {
		x, err := rand.Int(rand.Reader, big.NewInt(interval.Max+1))
		if err != nil {
			t.Fatal(err)
		}
		if seen[x.Int64()] {
			continue
		}
		seen[x.Int64()] = true
		y := new(bn256.G1).ScalarMult(h0, x)
		start := time.Now()
		got, ok := boundedDLog(h0, y, interval)
		elapsed := time.Since(start).Nanoseconds()
		if !ok || got != x.Int64() {
			t.Fatal("sample recovery failed")
		}
		rows = append(rows, sample{len(rows) + 1, x.Int64(), hex.EncodeToString(y.Marshal()), elapsed, true})
	}
	report := map[string]any{
		"measuredUTC": time.Now().UTC().Format(time.RFC3339Nano),
		"group":       "BN254 G1", "fixtureID": fixtureID, "h0Hex": h0Hex,
		"intervalMin": interval.Min, "intervalMax": interval.Max,
		"goVersion": runtime.Version(), "goArch": runtime.GOARCH, "gomaxprocs": runtime.GOMAXPROCS(0),
		"measurement": "full boundedDLog call; sampling, target generation, warmups, and verification excluded",
		"samples":     rows,
	}
	f, err := os.OpenFile(out, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0644)
	if err != nil {
		t.Fatal(err)
	}
	err = json.NewEncoder(f).Encode(report)
	closeErr := f.Close()
	if err != nil {
		t.Fatal(err)
	}
	if closeErr != nil {
		t.Fatal(closeErr)
	}
	t.Log("10 voting-group targets recovered")
}
