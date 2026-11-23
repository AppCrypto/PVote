// Package mobilecore is the on-device boundary to the repository's RB-PVSS.
// It never performs network requests or delegates cryptography to a server.
package mobilecore

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"

	rb "PVote/crypto/RBPVSS"
	bn "github.com/ethereum/go-ethereum/crypto/bn256/cloudflare"
)

const Legacy = "legacy-solidity-v1"
const Bound = "context-v1"

type Parameters struct {
	N     int      `json:"n"`
	T     int      `json:"t"`
	L     int      `json:"l"`
	Min   int64    `json:"min"`
	Max   int64    `json:"max"`
	G0    string   `json:"g0"`
	H0    string   `json:"h0"`
	G1    string   `json:"g1"`
	PKI   string   `json:"pki"`
	PK    []string `json:"pk"`
	Sigma []string `json:"sigma"`
}

type Fixture struct {
	ID         string     `json:"id"`
	Parameters Parameters `json:"parameters"`
}

type Election struct {
	SchemaVersion int        `json:"schemaVersion"`
	ElectionID    string     `json:"electionId"`
	Protocol      string     `json:"protocol"`
	RPCURL        string     `json:"rpcURL"`
	Contract      string     `json:"contract"`
	ChainID       string     `json:"chainId"`
	VoterAddress  string     `json:"voterAddress"`
	VoterIndex    uint64     `json:"voterIndex"`
	Labels        []string   `json:"labels"`
	Parameters    Parameters `json:"parameters"`
}

func Export(pp *rb.PublicParameters) Parameters {
	p := Parameters{N: len(pp.N), T: pp.Threshold, L: len(pp.D), Min: pp.Range.Min, Max: pp.Range.Max,
		G0: hex.EncodeToString(pp.PP1.G0.Marshal()), H0: hex.EncodeToString(pp.PP1.H0.Marshal()),
		G1: hex.EncodeToString(pp.PP2.G1.Marshal()), PKI: hex.EncodeToString(pp.PP2.PKI.Marshal())}
	for _, x := range pp.PK {
		p.PK = append(p.PK, hex.EncodeToString(x.Marshal()))
	}
	for _, x := range pp.PP2.Sigma {
		p.Sigma = append(p.Sigma, hex.EncodeToString(x.Marshal()))
	}
	return p
}

func (p Parameters) Decode() (*rb.PublicParameters, error) {
	if p.N < 1 || p.N > 200 || p.L < 1 || p.L > 64 || p.L >= p.T || p.T > p.N || p.Min > p.Max ||
		p.Min < -1000000 || p.Max > 1000000 || p.Max-p.Min > 255 || len(p.PK) != p.N || len(p.Sigma) != int(p.Max-p.Min+1) {
		return nil, fmt.Errorf("invalid or oversized mobile parameters")
	}
	g1decode := func(s string) (*bn.G1, error) {
		raw, e := hex.DecodeString(strings.TrimPrefix(s, "0x"))
		if e != nil || len(raw) != 64 {
			return nil, fmt.Errorf("invalid G1 encoding")
		}
		point := new(bn.G1)
		if _, e = point.Unmarshal(raw); e != nil {
			return nil, e
		}
		return point, nil
	}
	g2decode := func(s string) (*bn.G2, error) {
		raw, e := hex.DecodeString(strings.TrimPrefix(s, "0x"))
		if e != nil || len(raw) != 128 {
			return nil, fmt.Errorf("invalid G2 encoding")
		}
		point := new(bn.G2)
		if _, e = point.Unmarshal(raw); e != nil {
			return nil, e
		}
		return point, nil
	}
	g0, e := g1decode(p.G0)
	if e != nil {
		return nil, e
	}
	h0, e := g1decode(p.H0)
	if e != nil {
		return nil, e
	}
	g1, e := g2decode(p.G1)
	if e != nil {
		return nil, e
	}
	pki, e := g2decode(p.PKI)
	if e != nil {
		return nil, e
	}
	pp := &rb.PublicParameters{PP1: rb.PVSSParameters{G0: g0, H0: h0, N: p.N, T: p.T, L: p.L},
		PP2: rb.RangeParameters{G0: g0, H0: h0, G1: g1, PKI: pki}, Threshold: p.T, Range: rb.Interval{Min: p.Min, Max: p.Max}}
	for _, s := range p.PK {
		x, e := g1decode(s)
		if e != nil {
			return nil, e
		}
		pp.PK = append(pp.PK, x)
	}
	for _, s := range p.Sigma {
		x, e := g1decode(s)
		if e != nil {
			return nil, e
		}
		pp.PP2.Sigma = append(pp.PP2.Sigma, x)
	}
	for i := 0; i < p.N; i++ {
		pp.N = append(pp.N, i+1)
	}
	for i := 0; i < p.L; i++ {
		pp.D = append(pp.D, int64(i-p.L+1))
	}
	if e := rb.ValidatePublicParameters(pp); e != nil {
		return nil, e
	}
	return pp, nil
}

func (p Parameters) Fingerprint() string {
	data, _ := json.Marshal(p)
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:])
}
