package main

import (
	"PVote/crypto/PVSS"
	"PVote/crypto/ZKRP"
	"crypto/rand"
	"crypto/sha256"
	"embed"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"math/big"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	bn256 "github.com/ethereum/go-ethereum/crypto/bn256/cloudflare"
)

//go:embed static/*
var staticFiles embed.FS

type App struct {
	mu      sync.RWMutex
	state   *DemoState
	assets  fs.FS
	index   []byte
	created time.Time
}

type DemoConfig struct {
	NumTalliers   int `json:"numTalliers"`
	NumCandidates int `json:"numCandidates"`
	Threshold     int `json:"threshold"`
	RangeMin      int `json:"rangeMin"`
	RangeMax      int `json:"rangeMax"`
}

type DemoState struct {
	Config      DemoConfig
	SessionID   string
	PP          ZKRP.PP
	TallierSKs  []*big.Int
	TallierPKs  []*bn256.G1
	AggregateC  []*bn256.G1
	AggregateU  []*bn256.G1
	Votes       []*VoterRecord
	Decryptions map[int]*TallierRecord
	Tally       *TallyRecord
	Events      []EventRecord
	CreatedAt   time.Time
}

type VoterRecord struct {
	ID            int
	Alias         string
	Scores        []int
	Share         *PVSS.SecretSharing
	Ballots       []*bn256.G1
	Proofs        []*ZKRP.Proof
	PVSSVerified  bool
	RangeVerified bool
	SubmittedAt   time.Time
}

type TallierRecord struct {
	ID          int
	Share       *bn256.G1
	Proof       PVSS.Proof
	Verified    bool
	DecryptedAt time.Time
}

type TallyRecord struct {
	Results     []int
	Points      []*bn256.G1
	Verified    bool
	FinalizedAt time.Time
}

type EventRecord struct {
	Role   string
	Title  string
	Detail string
	At     time.Time
}

type SetupRequest struct {
	NumTalliers   int `json:"numTalliers"`
	NumCandidates int `json:"numCandidates"`
	Threshold     int `json:"threshold"`
	RangeMin      int `json:"rangeMin"`
	RangeMax      int `json:"rangeMax"`
}

type VoteRequest struct {
	Alias  string `json:"alias"`
	Scores []int  `json:"scores"`
}

type APIResponse struct {
	Message string        `json:"message,omitempty"`
	Error   string        `json:"error,omitempty"`
	State   StateSnapshot `json:"state,omitempty"`
}

type StateSnapshot struct {
	Meta      MetaSnapshot      `json:"meta"`
	Overview  OverviewSnapshot  `json:"overview"`
	Initiator InitiatorSnapshot `json:"initiator"`
	Voters    []VoterSnapshot   `json:"voters"`
	Talliers  []TallierSnapshot `json:"talliers"`
	Tally     *TallySnapshot    `json:"tally,omitempty"`
	Timeline  []EventSnapshot   `json:"timeline"`
}

type MetaSnapshot struct {
	SessionID   string `json:"sessionId"`
	Phase       string `json:"phase"`
	GeneratedAt string `json:"generatedAt"`
}

type OverviewSnapshot struct {
	VoteCount          int    `json:"voteCount"`
	VerifiedVotes      int    `json:"verifiedVotes"`
	DecryptionCount    int    `json:"decryptionCount"`
	Threshold          int    `json:"threshold"`
	RangeLabel         string `json:"rangeLabel"`
	CanVote            bool   `json:"canVote"`
	CanDecrypt         bool   `json:"canDecrypt"`
	CanFinalize        bool   `json:"canFinalize"`
	PlaintextAvailable bool   `json:"plaintextAvailable"`
}

type InitiatorSnapshot struct {
	Config          DemoConfig           `json:"config"`
	CandidateLabels []string             `json:"candidateLabels"`
	PublicParams    ParameterSnapshot    `json:"publicParams"`
	TallierKeys     []TallierKeySnapshot `json:"tallierKeys"`
	Aggregate       AggregateSnapshot    `json:"aggregate"`
}

type ParameterSnapshot struct {
	G0  string `json:"g0"`
	H0  string `json:"h0"`
	G1  string `json:"g1"`
	PKI string `json:"pki"`
}

type TallierKeySnapshot struct {
	ID        int    `json:"id"`
	PublicKey string `json:"publicKey"`
}

type AggregateSnapshot struct {
	EncryptedShares []string `json:"encryptedShares"`
	BallotCipher    []string `json:"ballotCipher"`
}

type VoterSnapshot struct {
	ID              int      `json:"id"`
	Alias           string   `json:"alias"`
	Scores          []int    `json:"scores"`
	PVSSVerified    bool     `json:"pvssVerified"`
	RangeVerified   bool     `json:"rangeVerified"`
	SubmittedAt     string   `json:"submittedAt"`
	BindCommitments []string `json:"bindCommitments"`
	EncryptedShares []string `json:"encryptedShares"`
	BallotCipher    []string `json:"ballotCipher"`
}

type TallierSnapshot struct {
	ID           int    `json:"id"`
	PublicKey    string `json:"publicKey"`
	HasDecrypted bool   `json:"hasDecrypted"`
	Verified     bool   `json:"verified"`
	Share        string `json:"share"`
	Proof        string `json:"proof"`
	DecryptedAt  string `json:"decryptedAt"`
}

type TallySnapshot struct {
	FinalizedAt string                    `json:"finalizedAt"`
	Verified    bool                      `json:"verified"`
	Results     []CandidateResultSnapshot `json:"results"`
}

type CandidateResultSnapshot struct {
	Label         string `json:"label"`
	Total         int    `json:"total"`
	ExpectedTotal int    `json:"expectedTotal"`
	Commitment    string `json:"commitment"`
}

type EventSnapshot struct {
	Role   string `json:"role"`
	Title  string `json:"title"`
	Detail string `json:"detail"`
	At     string `json:"at"`
}

func NewApp() (*App, error) {
	assets, err := fs.Sub(staticFiles, "static")
	if err != nil {
		return nil, fmt.Errorf("load static assets: %w", err)
	}

	index, err := fs.ReadFile(assets, "index.html")
	if err != nil {
		return nil, fmt.Errorf("load index: %w", err)
	}

	state, err := newDemoState(defaultConfig())
	if err != nil {
		return nil, err
	}

	return &App{
		state:   state,
		assets:  assets,
		index:   index,
		created: time.Now(),
	}, nil
}

func (a *App) routes() http.Handler {
	mux := http.NewServeMux()
	mux.Handle("GET /assets/", http.StripPrefix("/assets/", http.FileServer(http.FS(a.assets))))
	mux.HandleFunc("GET /api/state", a.handleState)
	mux.HandleFunc("POST /api/setup", a.handleSetup)
	mux.HandleFunc("POST /api/voters", a.handleVote)
	mux.HandleFunc("POST /api/talliers/{id}/decrypt", a.handleDecryptTallier)
	mux.HandleFunc("POST /api/tally/finalize", a.handleFinalizeTally)
	mux.HandleFunc("/", a.handleIndex)
	return mux
}

func (a *App) handleIndex(w http.ResponseWriter, r *http.Request) {
	if strings.HasPrefix(r.URL.Path, "/api/") {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write(a.index)
}

func (a *App) handleState(w http.ResponseWriter, _ *http.Request) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	writeJSON(w, http.StatusOK, APIResponse{
		State: a.state.snapshot(),
	})
}

func (a *App) handleSetup(w http.ResponseWriter, r *http.Request) {
	var req SetupRequest
	if err := decodeJSON(r.Body, &req); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	state, err := newDemoState(DemoConfig{
		NumTalliers:   req.NumTalliers,
		NumCandidates: req.NumCandidates,
		Threshold:     req.Threshold,
		RangeMin:      req.RangeMin,
		RangeMax:      req.RangeMax,
	})
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	a.mu.Lock()
	a.state = state
	snapshot := a.state.snapshot()
	a.mu.Unlock()

	writeJSON(w, http.StatusOK, APIResponse{
		Message: "Initiator rebuilt the demo session.",
		State:   snapshot,
	})
}

func (a *App) handleVote(w http.ResponseWriter, r *http.Request) {
	var req VoteRequest
	if err := decodeJSON(r.Body, &req); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	a.mu.Lock()
	err := a.state.submitVote(req.Alias, req.Scores)
	snapshot := a.state.snapshot()
	a.mu.Unlock()

	if err != nil {
		writeErrorWithState(w, http.StatusConflict, err, snapshot)
		return
	}

	writeJSON(w, http.StatusOK, APIResponse{
		Message: "Voter ballot accepted. PVSS and range proofs both verified.",
		State:   snapshot,
	})
}

func (a *App) handleDecryptTallier(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, fmt.Errorf("invalid tallier id"))
		return
	}

	a.mu.Lock()
	err = a.state.decryptTallier(id)
	snapshot := a.state.snapshot()
	a.mu.Unlock()

	if err != nil {
		writeErrorWithState(w, http.StatusConflict, err, snapshot)
		return
	}

	writeJSON(w, http.StatusOK, APIResponse{
		Message: fmt.Sprintf("Tallier %02d published a verified decryption share.", id),
		State:   snapshot,
	})
}

func (a *App) handleFinalizeTally(w http.ResponseWriter, _ *http.Request) {
	a.mu.Lock()
	err := a.state.finalizeTally()
	snapshot := a.state.snapshot()
	a.mu.Unlock()

	if err != nil {
		writeErrorWithState(w, http.StatusConflict, err, snapshot)
		return
	}

	writeJSON(w, http.StatusOK, APIResponse{
		Message: "Tallier quorum reached. Final tally reconstructed from encrypted commitments.",
		State:   snapshot,
	})
}

func defaultConfig() DemoConfig {
	return DemoConfig{
		NumTalliers:   4,
		NumCandidates: 3,
		Threshold:     3,
		RangeMin:      0,
		RangeMax:      5,
	}
}

func newDemoState(cfg DemoConfig) (*DemoState, error) {
	cfg, err := normalizeConfig(cfg)
	if err != nil {
		return nil, err
	}

	_, pp := ZKRP.Setup(cfg.RangeMin, cfg.RangeMax)
	sks, pks := PVSS.Setup(cfg.NumTalliers, pp.G0)

	state := &DemoState{
		Config:      cfg,
		SessionID:   newSessionID(),
		PP:          pp,
		TallierSKs:  sks,
		TallierPKs:  pks,
		AggregateC:  make([]*bn256.G1, cfg.NumTalliers),
		AggregateU:  make([]*bn256.G1, cfg.NumCandidates),
		Decryptions: make(map[int]*TallierRecord),
		CreatedAt:   time.Now(),
	}

	for i := range state.AggregateC {
		state.AggregateC[i] = zeroPoint()
	}
	for i := range state.AggregateU {
		state.AggregateU[i] = zeroPoint()
	}

	state.appendEvent(
		"initiator",
		"Demo session initialized",
		fmt.Sprintf(
			"Setup produced %d tallier key pairs, %d candidates, threshold %d, score range [%d, %d].",
			cfg.NumTalliers,
			cfg.NumCandidates,
			cfg.Threshold,
			cfg.RangeMin,
			cfg.RangeMax,
		),
	)

	return state, nil
}

func normalizeConfig(cfg DemoConfig) (DemoConfig, error) {
	defaults := defaultConfig()

	if cfg.NumTalliers == 0 {
		cfg.NumTalliers = defaults.NumTalliers
	}
	if cfg.NumCandidates == 0 {
		cfg.NumCandidates = defaults.NumCandidates
	}
	if cfg.RangeMax == 0 && cfg.RangeMin == 0 {
		cfg.RangeMin = defaults.RangeMin
		cfg.RangeMax = defaults.RangeMax
	}
	if cfg.Threshold == 0 {
		cfg.Threshold = minInt(cfg.NumTalliers, maxInt(2, (cfg.NumTalliers+cfg.NumCandidates)/2))
	}

	if cfg.NumTalliers < 2 || cfg.NumTalliers > 8 {
		return cfg, errors.New("numTalliers must be between 2 and 8")
	}
	if cfg.NumCandidates < 1 || cfg.NumCandidates > 6 {
		return cfg, errors.New("numCandidates must be between 1 and 6")
	}
	if cfg.RangeMin < 0 {
		return cfg, errors.New("rangeMin must be >= 0")
	}
	if cfg.RangeMax < cfg.RangeMin {
		return cfg, errors.New("rangeMax must be >= rangeMin")
	}
	if cfg.Threshold < 2 || cfg.Threshold > cfg.NumTalliers {
		return cfg, errors.New("threshold must be between 2 and numTalliers")
	}

	return cfg, nil
}

func (s *DemoState) submitVote(alias string, scores []int) error {
	if s.Tally != nil {
		return errors.New("tally already finalized; rebuild the session to vote again")
	}
	if len(s.Decryptions) > 0 {
		return errors.New("tallier decryption has already started; voting is now closed")
	}
	if len(scores) != s.Config.NumCandidates {
		return fmt.Errorf("expected %d candidate scores", s.Config.NumCandidates)
	}

	for _, score := range scores {
		if score < s.Config.RangeMin || score > s.Config.RangeMax {
			return fmt.Errorf("score %d is outside [%d, %d]", score, s.Config.RangeMin, s.Config.RangeMax)
		}
	}

	if strings.TrimSpace(alias) == "" {
		alias = fmt.Sprintf("Voter %02d", len(s.Votes)+1)
	}

	secret, err := rand.Int(rand.Reader, bn256.Order)
	if err != nil {
		return fmt.Errorf("generate voter secret: %w", err)
	}

	share := PVSS.Share(secret, s.PP.H0, s.TallierPKs, s.Config.Threshold, s.Config.NumTalliers, s.Config.NumCandidates)
	if !PVSS.DVerify(share, s.PP.H0, s.TallierPKs, s.Config.NumTalliers, s.Config.NumCandidates) {
		return errors.New("pvss verification failed")
	}

	coefficients, err := randomPolynomial(s.Config.Threshold)
	if err != nil {
		return fmt.Errorf("generate proof polynomial: %w", err)
	}

	indices := sequentialIndices(s.Config.Threshold)
	selectedShares := share.V[:s.Config.Threshold]
	ballots := make([]*bn256.G1, s.Config.NumCandidates)
	proofs := make([]*ZKRP.Proof, s.Config.NumCandidates)

	for i, score := range scores {
		scoreValue := big.NewInt(int64(score))
		ballots[i] = new(bn256.G1).Add(
			new(bn256.G1).ScalarMult(s.PP.G0, share.BindValue[i]),
			new(bn256.G1).ScalarMult(s.PP.H0, scoreValue),
		)

		x := candidateCoordinate(i)
		proofs[i] = ZKRP.GenProof(
			s.PP.G0,
			s.PP.H0,
			s.PP.G1,
			share.BindValue[i],
			scoreValue,
			ballots[i],
			s.PP.Sigma_k[score-s.Config.RangeMin],
			x,
			coefficients,
		)

		if !ZKRP.Verify(s.PP.G0, s.PP.H0, s.PP.G1, s.PP.PKI, proofs[i], ballots[i], x, selectedShares, indices, s.Config.Threshold) {
			return fmt.Errorf("range proof verification failed for candidate %d", i+1)
		}
	}

	for i := range s.AggregateC {
		s.AggregateC[i] = new(bn256.G1).Add(s.AggregateC[i], share.C[i])
	}
	for i := range s.AggregateU {
		s.AggregateU[i] = new(bn256.G1).Add(s.AggregateU[i], ballots[i])
	}

	s.Votes = append(s.Votes, &VoterRecord{
		ID:            len(s.Votes) + 1,
		Alias:         alias,
		Scores:        cloneInts(scores),
		Share:         share,
		Ballots:       ballots,
		Proofs:        proofs,
		PVSSVerified:  true,
		RangeVerified: true,
		SubmittedAt:   time.Now(),
	})

	s.appendEvent(
		"voter",
		fmt.Sprintf("%s cast a ballot", alias),
		fmt.Sprintf("Scores %v were encrypted into PVSS shares and candidate commitments.", scores),
	)

	return nil
}

func (s *DemoState) decryptTallier(id int) error {
	if len(s.Votes) == 0 {
		return errors.New("no ballots available yet")
	}
	if s.Tally != nil {
		return errors.New("tally already finalized")
	}
	if id < 1 || id > s.Config.NumTalliers {
		return errors.New("tallier id is out of range")
	}
	if _, exists := s.Decryptions[id]; exists {
		return errors.New("this tallier already published a decryption share")
	}

	share, proof := PVSS.Decrypt(s.PP.G0, s.TallierPKs[id-1], s.AggregateC[id-1], s.TallierSKs[id-1])
	if !PVSS.PVerify(s.PP.G0, s.TallierPKs[id-1], s.AggregateC[id-1], share, proof) {
		return errors.New("generated decryption proof failed verification")
	}

	s.Decryptions[id] = &TallierRecord{
		ID:          id,
		Share:       share,
		Proof:       proof,
		Verified:    true,
		DecryptedAt: time.Now(),
	}

	s.appendEvent(
		"tallier",
		fmt.Sprintf("Tallier %02d decrypted an aggregate share", id),
		fmt.Sprintf("The decryption proof validated against the aggregated ciphertext for tallier %02d.", id),
	)

	return nil
}

func (s *DemoState) finalizeTally() error {
	if len(s.Votes) == 0 {
		return errors.New("no ballots submitted")
	}
	if s.Tally != nil {
		return errors.New("tally already finalized")
	}

	verified := s.verifiedTalliers()
	if len(verified) < s.Config.Threshold {
		return fmt.Errorf("need %d verified talliers, currently have %d", s.Config.Threshold, len(verified))
	}

	indices := make([]*big.Int, s.Config.Threshold)
	shares := make([]*bn256.G1, s.Config.Threshold)
	for i, rec := range verified[:s.Config.Threshold] {
		indices[i] = big.NewInt(int64(rec.ID))
		shares[i] = rec.Share
	}

	results := make([]int, s.Config.NumCandidates)
	points := make([]*bn256.G1, s.Config.NumCandidates)
	expected := s.expectedTallies()
	minTotal := len(s.Votes) * s.Config.RangeMin
	maxTotal := len(s.Votes) * s.Config.RangeMax
	verifiedPlaintext := true

	for d := 0; d < s.Config.NumCandidates; d++ {
		coefficients := PVSS.LagrangeCoefficient(candidateCoordinate(d), indices, s.Config.Threshold)
		blindCommitment := PVSS.Reconstruct(coefficients, shares)
		points[d] = new(bn256.G1).Add(s.AggregateU[d], new(bn256.G1).Neg(blindCommitment))

		value, err := decodeCommitment(s.PP.H0, points[d], minTotal, maxTotal)
		if err != nil {
			return fmt.Errorf("failed to decode tally for candidate %d: %w", d+1, err)
		}

		results[d] = value
		if value != expected[d] {
			verifiedPlaintext = false
		}
	}

	s.Tally = &TallyRecord{
		Results:     results,
		Points:      points,
		Verified:    verifiedPlaintext,
		FinalizedAt: time.Now(),
	}

	s.appendEvent(
		"tallier",
		"Threshold tally finalized",
		fmt.Sprintf("Recovered %d candidate totals using %d verified tallier shares.", s.Config.NumCandidates, s.Config.Threshold),
	)

	return nil
}

func (s *DemoState) verifiedTalliers() []*TallierRecord {
	records := make([]*TallierRecord, 0, len(s.Decryptions))
	for _, rec := range s.Decryptions {
		if rec.Verified {
			records = append(records, rec)
		}
	}
	sort.Slice(records, func(i, j int) bool {
		return records[i].ID < records[j].ID
	})
	return records
}

func (s *DemoState) expectedTallies() []int {
	sums := make([]int, s.Config.NumCandidates)
	for _, vote := range s.Votes {
		for i, score := range vote.Scores {
			sums[i] += score
		}
	}
	return sums
}

func (s *DemoState) phase() string {
	switch {
	case s.Tally != nil:
		return "completed"
	case len(s.Decryptions) > 0:
		return "tallying"
	case len(s.Votes) > 0:
		return "voting"
	default:
		return "setup"
	}
}

func (s *DemoState) snapshot() StateSnapshot {
	verifiedVotes := 0
	for _, vote := range s.Votes {
		if vote.PVSSVerified && vote.RangeVerified {
			verifiedVotes++
		}
	}

	tallierKeys := make([]TallierKeySnapshot, len(s.TallierPKs))
	for i, pk := range s.TallierPKs {
		tallierKeys[i] = TallierKeySnapshot{
			ID:        i + 1,
			PublicKey: pointFingerprint(pk),
		}
	}

	aggregate := AggregateSnapshot{
		EncryptedShares: make([]string, len(s.AggregateC)),
		BallotCipher:    make([]string, len(s.AggregateU)),
	}
	for i, point := range s.AggregateC {
		aggregate.EncryptedShares[i] = pointFingerprint(point)
	}
	for i, point := range s.AggregateU {
		aggregate.BallotCipher[i] = pointFingerprint(point)
	}

	voters := make([]VoterSnapshot, len(s.Votes))
	for i, vote := range s.Votes {
		binds := make([]string, s.Config.NumCandidates)
		encrypted := make([]string, s.Config.NumTalliers)
		ballots := make([]string, s.Config.NumCandidates)

		for idx := 0; idx < s.Config.NumCandidates; idx++ {
			binds[idx] = pointFingerprint(vote.Share.V[s.Config.NumTalliers+idx])
			ballots[idx] = pointFingerprint(vote.Ballots[idx])
		}
		for idx := 0; idx < s.Config.NumTalliers; idx++ {
			encrypted[idx] = pointFingerprint(vote.Share.C[idx])
		}

		voters[i] = VoterSnapshot{
			ID:              vote.ID,
			Alias:           vote.Alias,
			Scores:          cloneInts(vote.Scores),
			PVSSVerified:    vote.PVSSVerified,
			RangeVerified:   vote.RangeVerified,
			SubmittedAt:     vote.SubmittedAt.Format("15:04:05"),
			BindCommitments: binds,
			EncryptedShares: encrypted,
			BallotCipher:    ballots,
		}
	}

	talliers := make([]TallierSnapshot, s.Config.NumTalliers)
	for i := 0; i < s.Config.NumTalliers; i++ {
		rec, ok := s.Decryptions[i+1]
		talliers[i] = TallierSnapshot{
			ID:           i + 1,
			PublicKey:    pointFingerprint(s.TallierPKs[i]),
			HasDecrypted: ok,
		}
		if ok {
			talliers[i].Verified = rec.Verified
			talliers[i].Share = pointFingerprint(rec.Share)
			talliers[i].Proof = scalarFingerprint(rec.Proof.C)
			talliers[i].DecryptedAt = rec.DecryptedAt.Format("15:04:05")
		}
	}

	timeline := make([]EventSnapshot, len(s.Events))
	for i, event := range s.Events {
		timeline[i] = EventSnapshot{
			Role:   event.Role,
			Title:  event.Title,
			Detail: event.Detail,
			At:     event.At.Format("15:04:05"),
		}
	}

	snapshot := StateSnapshot{
		Meta: MetaSnapshot{
			SessionID:   s.SessionID,
			Phase:       s.phase(),
			GeneratedAt: time.Now().Format("15:04:05"),
		},
		Overview: OverviewSnapshot{
			VoteCount:          len(s.Votes),
			VerifiedVotes:      verifiedVotes,
			DecryptionCount:    len(s.verifiedTalliers()),
			Threshold:          s.Config.Threshold,
			RangeLabel:         fmt.Sprintf("[%d, %d]", s.Config.RangeMin, s.Config.RangeMax),
			CanVote:            s.Tally == nil && len(s.Decryptions) == 0,
			CanDecrypt:         len(s.Votes) > 0 && s.Tally == nil,
			CanFinalize:        s.Tally == nil && len(s.verifiedTalliers()) >= s.Config.Threshold,
			PlaintextAvailable: s.Tally != nil,
		},
		Initiator: InitiatorSnapshot{
			Config:          s.Config,
			CandidateLabels: candidateLabels(s.Config.NumCandidates),
			PublicParams: ParameterSnapshot{
				G0:  pointFingerprint(s.PP.G0),
				H0:  pointFingerprint(s.PP.H0),
				G1:  point2Fingerprint(s.PP.G1),
				PKI: point2Fingerprint(s.PP.PKI),
			},
			TallierKeys: tallierKeys,
			Aggregate:   aggregate,
		},
		Voters:   voters,
		Talliers: talliers,
		Timeline: timeline,
	}

	if s.Tally != nil {
		results := make([]CandidateResultSnapshot, s.Config.NumCandidates)
		expected := s.expectedTallies()
		for i := 0; i < s.Config.NumCandidates; i++ {
			results[i] = CandidateResultSnapshot{
				Label:         candidateLabels(s.Config.NumCandidates)[i],
				Total:         s.Tally.Results[i],
				ExpectedTotal: expected[i],
				Commitment:    pointFingerprint(s.Tally.Points[i]),
			}
		}

		snapshot.Tally = &TallySnapshot{
			FinalizedAt: s.Tally.FinalizedAt.Format("15:04:05"),
			Verified:    s.Tally.Verified,
			Results:     results,
		}
	}

	return snapshot
}

func (s *DemoState) appendEvent(role, title, detail string) {
	s.Events = append([]EventRecord{{
		Role:   role,
		Title:  title,
		Detail: detail,
		At:     time.Now(),
	}}, s.Events...)
	if len(s.Events) > 18 {
		s.Events = s.Events[:18]
	}
}

func decodeCommitment(base *bn256.G1, point *bn256.G1, start, end int) (int, error) {
	target := point.Marshal()
	for value := start; value <= end; value++ {
		candidate := new(bn256.G1).ScalarMult(base, big.NewInt(int64(value)))
		if string(candidate.Marshal()) == string(target) {
			return value, nil
		}
	}
	return 0, fmt.Errorf("commitment did not match any scalar in [%d, %d]", start, end)
}

func randomPolynomial(length int) ([]*big.Int, error) {
	coefficients := make([]*big.Int, length)
	for i := range coefficients {
		value, err := rand.Int(rand.Reader, bn256.Order)
		if err != nil {
			return nil, err
		}
		coefficients[i] = value
	}
	return coefficients, nil
}

func sequentialIndices(threshold int) []*big.Int {
	indices := make([]*big.Int, threshold)
	for i := 0; i < threshold; i++ {
		indices[i] = big.NewInt(int64(i + 1))
	}
	return indices
}

func candidateCoordinate(index int) *big.Int {
	x := new(big.Int).Neg(big.NewInt(int64(index)))
	return x.Mod(x, bn256.Order)
}

func zeroPoint() *bn256.G1 {
	return new(bn256.G1).ScalarBaseMult(big.NewInt(0))
}

func newSessionID() string {
	buf := make([]byte, 4)
	if _, err := rand.Read(buf); err != nil {
		return fmt.Sprintf("session-%d", time.Now().UnixNano())
	}
	return fmt.Sprintf("session-%s", hex.EncodeToString(buf))
}

func pointFingerprint(point *bn256.G1) string {
	if point == nil {
		return ""
	}
	return "G1 " + fingerprintBytes(point.Marshal())
}

func point2Fingerprint(point *bn256.G2) string {
	if point == nil {
		return ""
	}
	return "G2 " + fingerprintBytes(point.Marshal())
}

func scalarFingerprint(value *big.Int) string {
	if value == nil {
		return ""
	}
	return "Fr " + fingerprintBytes(value.Bytes())
}

func fingerprintBytes(input []byte) string {
	sum := sha256.Sum256(input)
	encoded := hex.EncodeToString(sum[:])
	return encoded[:12] + "..." + encoded[len(encoded)-8:]
}

func candidateLabels(count int) []string {
	labels := make([]string, count)
	for i := 0; i < count; i++ {
		labels[i] = fmt.Sprintf("Candidate %c", rune('A'+i))
	}
	return labels
}

func cloneInts(input []int) []int {
	out := make([]int, len(input))
	copy(out, input)
	return out
}

func decodeJSON(body io.ReadCloser, target any) error {
	defer body.Close()
	decoder := json.NewDecoder(body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(target); err != nil {
		return err
	}
	return nil
}

func writeError(w http.ResponseWriter, code int, err error) {
	writeJSON(w, code, APIResponse{Error: err.Error()})
}

func writeErrorWithState(w http.ResponseWriter, code int, err error, state StateSnapshot) {
	writeJSON(w, code, APIResponse{
		Error: err.Error(),
		State: state,
	})
}

func writeJSON(w http.ResponseWriter, code int, payload APIResponse) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(payload)
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}
