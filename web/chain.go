package main

// This file connects the web workflow to both the exact RB-PVSS verifier and
// the optional settlement contract. Both consume the same canonical Go data.

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"os"
	"strings"

	chainencoding "PVote/chain/rbpvss"
	rbpvsscontract "PVote/compile/contract"
	rbpvss "PVote/crypto/RBPVSS"
	stakecontract "PVote/web/contract"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	bn256 "github.com/ethereum/go-ethereum/crypto/bn256/cloudflare"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/joho/godotenv"
)

const ganacheURL = "http://127.0.0.1:8545"
const stakeTxGasLimit uint64 = 8_000_000

type StakeChain struct {
	URL             string
	Client          *ethclient.Client
	Contract        *stakecontract.Stakecontract
	ContractAddress common.Address
	Verifier        *rbpvsscontract.Contract
	VerifierAddress common.Address
	Initiator       RoleAccount
	Talliers        []RoleAccount
	Voters          []RoleAccount
	Config          StakeConfig
}

type StakeConfig struct {
	InitiatorEscrowWei *big.Int
	VoterStakeWei      *big.Int
	TallierStakeWei    *big.Int
	InitiatorRewardPct uint8
	VoterRewardPct     uint8
	TallierRewardPct   uint8
}

type RoleAccount struct {
	PrivateKey string
	Address    common.Address
}

type ChainOverview struct {
	EscrowFunded    bool
	Settled         bool
	TotalEscrow     *big.Int
	RewardPool      *big.Int
	ContractBalance *big.Int
	VoterCount      *big.Int
	TallierCount    *big.Int
	SettledAt       *big.Int
}

type ChainParticipant struct {
	Address       common.Address
	Deposited     *big.Int
	Claimable     *big.Int
	WalletBalance *big.Int
	Staked        bool
	Honest        bool
	Withdrawn     bool
}

func NewStakeChain(cfg DemoConfig, pp *rbpvss.PublicParameters) (*StakeChain, error) {
	stakeCfg, err := newStakeConfig(cfg)
	if err != nil {
		return nil, err
	}
	privateKeys, err := loadDemoPrivateKeys()
	if err != nil {
		return nil, err
	}
	if len(privateKeys) < cfg.NumTalliers+2 {
		return nil, fmt.Errorf("need at least %d private keys in .env for initiator, talliers, and voters", cfg.NumTalliers+2)
	}
	client, err := ethclient.Dial(ganacheURL)
	if err != nil {
		return nil, fmt.Errorf("connect to ganache: %w", err)
	}
	initiator, err := newRoleAccount(privateKeys[0])
	if err != nil {
		return nil, err
	}
	talliers := make([]RoleAccount, cfg.NumTalliers)
	for i := range talliers {
		talliers[i], err = newRoleAccount(privateKeys[i+1])
		if err != nil {
			return nil, err
		}
	}
	voterKeys := privateKeys[cfg.NumTalliers+1:]
	if len(voterKeys) == 0 {
		return nil, errorsf("no voter accounts left in .env after assigning %d talliers", cfg.NumTalliers)
	}
	voters := make([]RoleAccount, len(voterKeys))
	for i := range voters {
		voters[i], err = newRoleAccount(voterKeys[i])
		if err != nil {
			return nil, err
		}
	}
	auth, err := newAuth(client, initiator.PrivateKey, big.NewInt(0))
	if err != nil {
		return nil, err
	}
	address, tx, contract, err := stakecontract.DeployStakecontract(
		auth,
		client,
		stakeCfg.InitiatorEscrowWei,
		stakeCfg.VoterStakeWei,
		stakeCfg.TallierStakeWei,
		stakeCfg.InitiatorRewardPct,
		stakeCfg.VoterRewardPct,
		stakeCfg.TallierRewardPct,
	)
	if err != nil {
		return nil, fmt.Errorf("deploy stake manager: %w", err)
	}
	if _, err := waitForTxReceipt(client, tx); err != nil {
		return nil, fmt.Errorf("wait for stake manager deployment: %w", err)
	}
	verifierAuth, err := newAuth(client, initiator.PrivateKey, big.NewInt(0))
	if err != nil {
		return nil, err
	}
	g0, h0, g1, pkI, sigma, pks, threshold, coordinates, minimum, maximum, err := chainencoding.ConstructorArguments(pp)
	if err != nil {
		return nil, fmt.Errorf("encode RB-PVSS parameters: %w", err)
	}
	verifierAddress, verifierTx, verifier, err := rbpvsscontract.DeployContract(
		verifierAuth, client, g0, h0, g1, pkI, sigma, pks, threshold, coordinates, minimum, maximum, roleAddresses(voters),
	)
	if err != nil {
		return nil, fmt.Errorf("deploy RB-PVSS verifier: %w", err)
	}
	if _, err := waitForTxReceipt(client, verifierTx); err != nil {
		return nil, fmt.Errorf("wait for RB-PVSS verifier deployment: %w", err)
	}
	return &StakeChain{
		URL:             ganacheURL,
		Client:          client,
		Contract:        contract,
		ContractAddress: address,
		Verifier:        verifier,
		VerifierAddress: verifierAddress,
		Initiator:       initiator,
		Talliers:        talliers,
		Voters:          voters,
		Config:          stakeCfg,
	}, nil
}

func roleAddresses(accounts []RoleAccount) []common.Address {
	addresses := make([]common.Address, len(accounts))
	for i := range accounts {
		addresses[i] = accounts[i].Address
	}
	return addresses
}

func (c *StakeChain) SubmitRB(voterID int, instance *rbpvss.Instance) error {
	if voterID < 1 || voterID > len(c.Voters) {
		return errorsf("voter id is out of range")
	}
	dealer, rangeBinding, err := chainencoding.Transcript(instance)
	if err != nil {
		return err
	}
	auth, err := newAuth(c.Client, c.Voters[voterID-1].PrivateKey, big.NewInt(0))
	if err != nil {
		return err
	}
	tx, err := c.Verifier.SubmitRB(auth, dealer, rangeBinding)
	if err != nil {
		return fmt.Errorf("submit atomic RB-PVSS transcript: %w", err)
	}
	return waitForTx(c.Client, tx)
}

func (c *StakeChain) SubmitDecryptionShare(tallierID int, share *bn256.G1, proof *rbpvss.DecryptionProof) error {
	if tallierID < 1 || tallierID > len(c.Talliers) {
		return errorsf("tallier id is out of range")
	}
	auth, err := newAuth(c.Client, c.Talliers[tallierID-1].PrivateKey, big.NewInt(0))
	if err != nil {
		return err
	}
	tx, err := c.Verifier.SubmitDecryptionShare(auth, big.NewInt(int64(tallierID)), chainencoding.G1Point(share), chainencoding.G1Point(proof.A), chainencoding.G1Point(proof.B), proof.Chi, proof.Z)
	if err != nil {
		return fmt.Errorf("submit decryption share: %w", err)
	}
	return waitForTx(c.Client, tx)
}

func (c *StakeChain) Reconstruct() ([]*bn256.G1, error) {
	encoded, err := c.Verifier.Reconstruct(&bind.CallOpts{Context: context.Background()})
	if err != nil {
		return nil, err
	}
	points := make([]*bn256.G1, len(encoded))
	for i := range encoded {
		points[i], err = chainencoding.G1(encoded[i])
		if err != nil {
			return nil, err
		}
	}
	return points, nil
}

func (c *StakeChain) FundInitiatorEscrow() error {
	auth, err := newAuth(c.Client, c.Initiator.PrivateKey, c.Config.InitiatorEscrowWei)
	if err != nil {
		return err
	}
	tx, err := c.Contract.FundInitiatorEscrow(auth)
	if err != nil {
		return fmt.Errorf("fund initiator escrow: %w", err)
	}
	return waitForTx(c.Client, tx)
}

func (c *StakeChain) StakeTallier(tallierID int) error {
	if tallierID < 1 || tallierID > len(c.Talliers) {
		return errorsf("tallier id is out of range")
	}
	auth, err := newAuth(c.Client, c.Talliers[tallierID-1].PrivateKey, c.Config.TallierStakeWei)
	if err != nil {
		return err
	}
	tx, err := c.Contract.DepositTallierStake(auth, big.NewInt(int64(tallierID)))
	if err != nil {
		return fmt.Errorf("stake tallier %d: %w", tallierID, err)
	}
	return waitForTx(c.Client, tx)
}

func (c *StakeChain) StakeVoter(voterID int) (RoleAccount, error) {
	if voterID < 1 || voterID > len(c.Voters) {
		return RoleAccount{}, errorsf("no ganache-backed voter account is available for voter %d", voterID)
	}
	role := c.Voters[voterID-1]
	auth, err := newAuth(c.Client, role.PrivateKey, c.Config.VoterStakeWei)
	if err != nil {
		return RoleAccount{}, err
	}
	tx, err := c.Contract.DepositVoterStake(auth, big.NewInt(int64(voterID)))
	if err != nil {
		return RoleAccount{}, fmt.Errorf("stake voter %d: %w", voterID, err)
	}
	if err := waitForTx(c.Client, tx); err != nil {
		return RoleAccount{}, err
	}
	return role, nil
}

func (c *StakeChain) SettleRewards(voterIDs, tallierIDs []int) error {
	auth, err := newAuth(c.Client, c.Initiator.PrivateKey, big.NewInt(0))
	if err != nil {
		return err
	}
	tx, err := c.Contract.SettleRewards(auth, intSliceToBigInts(voterIDs), intSliceToBigInts(tallierIDs))
	if err != nil {
		return fmt.Errorf("settle rewards: %w", err)
	}
	return waitForTx(c.Client, tx)
}

func (c *StakeChain) WithdrawInitiator() error {
	auth, err := newAuth(c.Client, c.Initiator.PrivateKey, big.NewInt(0))
	if err != nil {
		return err
	}
	tx, err := c.Contract.WithdrawInitiator(auth)
	if err != nil {
		return fmt.Errorf("withdraw initiator: %w", err)
	}
	return waitForTx(c.Client, tx)
}

func (c *StakeChain) WithdrawVoter(voterID int) error {
	if voterID < 1 || voterID > len(c.Voters) {
		return errorsf("voter id is out of range")
	}
	auth, err := newAuth(c.Client, c.Voters[voterID-1].PrivateKey, big.NewInt(0))
	if err != nil {
		return err
	}
	tx, err := c.Contract.WithdrawVoter(auth, big.NewInt(int64(voterID)))
	if err != nil {
		return fmt.Errorf("withdraw voter %d: %w", voterID, err)
	}
	return waitForTx(c.Client, tx)
}

func (c *StakeChain) WithdrawTallier(tallierID int) error {
	if tallierID < 1 || tallierID > len(c.Talliers) {
		return errorsf("tallier id is out of range")
	}
	auth, err := newAuth(c.Client, c.Talliers[tallierID-1].PrivateKey, big.NewInt(0))
	if err != nil {
		return err
	}
	tx, err := c.Contract.WithdrawTallier(auth, big.NewInt(int64(tallierID)))
	if err != nil {
		return fmt.Errorf("withdraw tallier %d: %w", tallierID, err)
	}
	return waitForTx(c.Client, tx)
}

func (c *StakeChain) ReadOverview() (*ChainOverview, error) {
	data, err := c.Contract.GetEscrowOverview(&bind.CallOpts{Context: context.Background()})
	if err != nil {
		return nil, err
	}
	return &ChainOverview{
		EscrowFunded:    data.EscrowFunded,
		Settled:         data.Settled,
		TotalEscrow:     data.TotalEscrow,
		RewardPool:      data.RewardPool,
		ContractBalance: data.ContractBalance,
		VoterCount:      data.VoterCount,
		TallierCount:    data.TallierCount,
		SettledAt:       data.SettledTimestamp,
	}, nil
}

func (c *StakeChain) ReadInitiator() (*ChainParticipant, error) {
	data, err := c.Contract.GetInitiatorState(&bind.CallOpts{Context: context.Background()})
	if err != nil {
		return nil, err
	}
	address := data.Account
	if address == (common.Address{}) {
		address = c.Initiator.Address
	}
	balance, err := c.Client.BalanceAt(context.Background(), c.Initiator.Address, nil)
	if err != nil {
		return nil, err
	}
	return &ChainParticipant{Address: address, Deposited: data.Deposited, Claimable: data.Claimable, WalletBalance: balance, Staked: data.Staked, Honest: true, Withdrawn: data.Withdrawn}, nil
}

func (c *StakeChain) ReadVoter(voterID int) (*ChainParticipant, error) {
	data, err := c.Contract.GetVoter(&bind.CallOpts{Context: context.Background()}, big.NewInt(int64(voterID)))
	if err != nil {
		return nil, err
	}
	balance := big.NewInt(0)
	if data.Account != (common.Address{}) {
		balance, err = c.Client.BalanceAt(context.Background(), data.Account, nil)
		if err != nil {
			return nil, err
		}
	}
	return &ChainParticipant{Address: data.Account, Deposited: data.Deposited, Claimable: data.Claimable, WalletBalance: balance, Staked: data.Staked, Honest: data.Honest, Withdrawn: data.Withdrawn}, nil
}

func (c *StakeChain) ReadTallier(tallierID int) (*ChainParticipant, error) {
	data, err := c.Contract.GetTallier(&bind.CallOpts{Context: context.Background()}, big.NewInt(int64(tallierID)))
	if err != nil {
		return nil, err
	}
	address := data.Account
	if address == (common.Address{}) && tallierID >= 1 && tallierID <= len(c.Talliers) {
		address = c.Talliers[tallierID-1].Address
	}
	balance := big.NewInt(0)
	if address != (common.Address{}) {
		balance, err = c.Client.BalanceAt(context.Background(), address, nil)
		if err != nil {
			return nil, err
		}
	}
	return &ChainParticipant{Address: address, Deposited: data.Deposited, Claimable: data.Claimable, WalletBalance: balance, Staked: data.Staked, Honest: data.Honest, Withdrawn: data.Withdrawn}, nil
}

func newStakeConfig(cfg DemoConfig) (StakeConfig, error) {
	initiatorEscrow, err := parseETHToWei(cfg.InitiatorEscrowETH)
	if err != nil {
		return StakeConfig{}, fmt.Errorf("parse initiator escrow: %w", err)
	}
	voterStake, err := parseETHToWei(cfg.VoterStakeETH)
	if err != nil {
		return StakeConfig{}, fmt.Errorf("parse voter stake: %w", err)
	}
	tallierStake, err := parseETHToWei(cfg.TallierStakeETH)
	if err != nil {
		return StakeConfig{}, fmt.Errorf("parse tallier stake: %w", err)
	}
	return StakeConfig{InitiatorEscrowWei: initiatorEscrow, VoterStakeWei: voterStake, TallierStakeWei: tallierStake, InitiatorRewardPct: uint8(cfg.InitiatorRewardPercent), VoterRewardPct: uint8(cfg.VoterRewardPercent), TallierRewardPct: uint8(cfg.TallierRewardPercent)}, nil
}

func parseETHToWei(value string) (*big.Int, error) {
	clean := strings.TrimSpace(value)
	if clean == "" {
		return nil, errorsf("missing ETH amount")
	}
	rat, ok := new(big.Rat).SetString(clean)
	if !ok {
		return nil, errorsf("invalid ETH amount %q", value)
	}
	if rat.Sign() < 0 {
		return nil, errorsf("ETH amount must be non-negative")
	}
	weiRat := new(big.Rat).Mul(rat, new(big.Rat).SetInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)))
	if !weiRat.IsInt() {
		return nil, errorsf("ETH amount %q has more than 18 decimals", value)
	}
	return new(big.Int).Set(weiRat.Num()), nil
}

func loadDemoPrivateKeys() ([]string, error) {
	env, err := godotenv.Read(".env")
	if err != nil {
		return nil, fmt.Errorf("read .env: %w", err)
	}
	privateKeys := make([]string, 0, 32)
	for i := 1; i <= 32; i++ {
		if key := strings.TrimSpace(env[fmt.Sprintf("PRIVATE_KEY_%d", i)]); key != "" {
			privateKeys = append(privateKeys, key)
		}
	}
	if len(privateKeys) == 0 {
		return nil, errorsf("no PRIVATE_KEY_* values found in .env")
	}
	return privateKeys, nil
}

func newRoleAccount(privateKeyHex string) (RoleAccount, error) {
	key, err := crypto.HexToECDSA(strings.TrimSpace(privateKeyHex))
	if err != nil {
		return RoleAccount{}, fmt.Errorf("parse private key: %w", err)
	}
	return RoleAccount{PrivateKey: privateKeyHex, Address: crypto.PubkeyToAddress(key.PublicKey)}, nil
}

func newAuth(client *ethclient.Client, privateKeyHex string, value *big.Int) (*bind.TransactOpts, error) {
	key, err := crypto.HexToECDSA(strings.TrimSpace(privateKeyHex))
	if err != nil {
		return nil, err
	}
	publicKey, ok := key.Public().(*ecdsa.PublicKey)
	if !ok {
		return nil, errorsf("failed to parse ECDSA public key")
	}
	from := crypto.PubkeyToAddress(*publicKey)
	nonce, err := client.PendingNonceAt(context.Background(), from)
	if err != nil {
		return nil, err
	}
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		return nil, err
	}
	chainID, err := client.ChainID(context.Background())
	if err != nil {
		return nil, err
	}
	auth, err := bind.NewKeyedTransactorWithChainID(key, chainID)
	if err != nil {
		return nil, err
	}
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = value
	auth.GasLimit = stakeTxGasLimit
	auth.GasPrice = gasPrice
	balance, err := client.BalanceAt(context.Background(), from, nil)
	if err != nil {
		return nil, err
	}
	required := new(big.Int).Add(new(big.Int).Mul(new(big.Int).SetUint64(stakeTxGasLimit), gasPrice), value)
	if balance.Cmp(required) < 0 {
		return nil, fmt.Errorf("ganache account %s has %s ETH, needs at least %s ETH", from.Hex(), formatWeiToETH(balance), formatWeiToETH(required))
	}
	return auth, nil
}

func waitForTx(client *ethclient.Client, tx *types.Transaction) error {
	_, err := waitForTxReceipt(client, tx)
	return err
}

func waitForTxReceipt(client *ethclient.Client, tx *types.Transaction) (*types.Receipt, error) {
	receipt, err := bind.WaitMined(context.Background(), client, tx)
	if err != nil {
		return nil, err
	}
	if receipt.Status != types.ReceiptStatusSuccessful {
		return nil, errorsf("transaction %s reverted", tx.Hash().Hex())
	}
	return receipt, nil
}

func intSliceToBigInts(values []int) []*big.Int {
	result := make([]*big.Int, len(values))
	for i, value := range values {
		result[i] = big.NewInt(int64(value))
	}
	return result
}

func errorsf(format string, args ...any) error { return fmt.Errorf(format, args...) }

func formatWeiToETH(value *big.Int) string {
	if value == nil {
		return "0.000"
	}
	rat := new(big.Rat).SetFrac(value, new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil))
	return rat.FloatString(3)
}

func envOrDefault(key, fallback string) string {
	if value := strings.TrimSpace(os.Getenv(key)); value != "" {
		return value
	}
	return fallback
}
