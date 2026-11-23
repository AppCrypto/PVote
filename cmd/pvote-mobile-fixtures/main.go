// Generates public-only fixtures on the Mac. Setup is excluded from voter time.
package main

import (
	rb "PVote/crypto/RBPVSS"
	core "PVote/mobile/core"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	out := flag.String("out", "mobile/fixtures/fixtures.json", "public fixture output")
	flag.Parse()
	pairs := [][3]int{{4, 2, 3}, {10, 5, 7}, {50, 10, 30}, {100, 19, 59}}
	for _, l := range []int{1, 10, 19} {
		for n := 10; n <= 100; n += 10 {
			if l <= n/3 {
				pairs = append(pairs, [3]int{n, l, (n + l) / 2})
			}
		}
	}
	seen := map[[3]int]bool{}
	fixtures := []core.Fixture{}
	for _, p := range pairs {
		if seen[p] {
			continue
		}
		seen[p] = true
		pp, _, err := rb.Setup(128, p[0], p[2], p[1], rb.Interval{Min: 0, Max: 10})
		if err != nil {
			panic(err)
		}
		fixtures = append(fixtures, core.Fixture{ID: fmt.Sprintf("n%d-l%d-t%d", p[0], p[1], p[2]), Parameters: core.Export(pp)})
	}
	data, err := json.MarshalIndent(fixtures, "", "  ")
	if err != nil {
		panic(err)
	}
	if err = os.MkdirAll(filepath.Dir(*out), 0755); err != nil {
		panic(err)
	}
	if err = os.WriteFile(*out, append(data, '\n'), 0644); err != nil {
		panic(err)
	}
	fmt.Printf("Wrote %d public-only fixtures to %s\n", len(fixtures), *out)
}
