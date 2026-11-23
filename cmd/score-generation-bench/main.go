// score-generation-bench executes the same three real generation paths on one
// native Mac process. It records individual timings, excluding setup and encoding.
package main

import (
	core "PVote/mobile/core"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

func main() {
	out := flag.String("out", "experiments/mac-generation-20261121/raw", "raw report directory")
	flag.Parse()
	runtime.GOMAXPROCS(1)
	var engine core.Engine
	var metadata map[string]any
	if err := json.Unmarshal(engine.Call([]byte(`{"action":"version"}`)), &metadata); err != nil {
		panic(err)
	}
	read := func(name string, args ...string) string {
		b, e := exec.Command(name, args...).Output()
		if e != nil {
			return "unavailable"
		}
		return strings.TrimSpace(string(b))
	}
	metadata["hardwareIdentifier"] = read("sysctl", "-n", "hw.model")
	metadata["cpu"] = read("sysctl", "-n", "machdep.cpu.brand_string")
	metadata["osVersion"] = read("sw_vers", "-productVersion")
	metadata["physicalMemoryBytes"] = read("sysctl", "-n", "hw.memsize")
	metadata["logicalCPUs"] = read("sysctl", "-n", "hw.logicalcpu")
	metadata["powerStatus"] = read("pmset", "-g", "batt")
	metadata["platform"] = "macOS"
	metadata["buildConfiguration"] = "Go optimized"
	metadata["isSimulator"] = false
	metadata["thermalTelemetry"] = "not collected on Mac"
	if runtime.GOOS != "darwin" || metadata["sourceSHA256"] == "development" {
		panic("requires a provenance-stamped native Mac build")
	}
	setups := make([]map[string]any, 0, 8)
	for k := 1; k <= 8; k++ {
		s, e := engine.PrepareComparison(k)
		must(e)
		setups = append(setups, s)
	}
	type item struct {
		k      int
		scheme string
	}
	var ascending []item
	for k := 1; k <= 8; k++ {
		for _, s := range []string{"PVote", "PriScore", "3PVSS-bit"} {
			ascending = append(ascending, item{k, s})
		}
	}
	must(os.MkdirAll(*out, 0755))
	for round, order := range []string{"ascending", "descending", "rotated"} {
		cases := append([]item{}, ascending...)
		if round == 1 {
			for i, j := 0, len(cases)-1; i < j; i, j = i+1, j-1 {
				cases[i], cases[j] = cases[j], cases[i]
			}
		}
		if round == 2 {
			cases = append(append([]item{}, ascending[12:]...), ascending[:12]...)
		}
		id := time.Now().UTC().Format("2006-01-02T15-04-05.000000000Z") + "-mac"
		ids := make([]string, 0, 24)
		for _, c := range cases {
			ids = append(ids, fmt.Sprintf("%s-k%d", c.scheme, c.k))
		}
		report := map[string]any{"schemaVersion": 1, "runId": id, "experiment": "score_domain_generation", "metadata": metadata,
			"runsPerConfiguration": 10, "warmupsPerConfiguration": 3, "recheckRound": round + 1, "configurationOrderPolicy": order,
			"fixtureIds": ids, "setupRecords": setups, "setupAndParameterLoadingTimed": false, "verifyEachMeasuredBallot": false,
			"scorePolicy": "same cyclic distribution of B=2^kappa-1 across five coordinates for all schemes",
			"startedAt":   time.Now().UTC().Format(time.RFC3339Nano), "status": "running"}
		rows := make([]map[string]any, 0, 240)
		verified := []string{}
		path := filepath.Join(*out, "benchmark-"+id+".json")
		persist := func() {
			report["samples"] = rows
			report["verifiedWarmupCases"] = verified
			b, e := json.MarshalIndent(report, "", "  ")
			must(e)
			must(os.WriteFile(path+".partial", b, 0644))
			must(os.Rename(path+".partial", path))
		}
		for _, c := range cases {
			fixture := fmt.Sprintf("%s-k%d", c.scheme, c.k)
			r := core.Request{Kappa: c.k, Scheme: c.scheme, ElectionID: "comparison-" + id, VoterIndex: 1}
			for warm := 0; warm < 3; warm++ {
				r.Verify = warm == 0
				_, e := engine.GenerateComparison(r)
				must(e)
			}
			verified = append(verified, fixture)
			r.Verify = false
			sum := 0.0
			for sample := 1; sample <= 10; sample++ {
				m, e := engine.GenerateComparison(r)
				must(e)
				m["fixtureId"] = fixture
				m["sampleIndex"] = sample
				m["timestamp"] = time.Now().UTC().Format(time.RFC3339Nano)
				rows = append(rows, m)
				sum += m["generationMs"].(float64)
				persist()
			}
			fmt.Printf("round=%d order=%s %-10s k=%d generation=%.3f ms (10 samples)\n", round+1, order, c.scheme, c.k, sum/10)
		}
		report["status"] = "complete"
		report["completedAt"] = time.Now().UTC().Format(time.RFC3339Nano)
		persist()
		fmt.Println("SAVED", path)
	}
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
