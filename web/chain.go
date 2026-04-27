package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"os"
	"strings"

	stakecontract "PVote/web/contract"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
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

func NewStakeChain(cfg DemoConfig) (*StakeChain, error) {
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
	for i := 0; i < cfg.NumTalliers; i++ {
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
	for i, privateKey := range voterKeys {
		voters[i], err = newRoleAccount(privateKey)
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
	if _, err := bind.WaitMined(context.Background(), client, tx); err != nil {
		return nil, fmt.Errorf("wait for stake manager deployment: %w", err)
	}

	return &StakeChain{
		URL:             ganacheURL,
		Client:          client,
		Contract:        contract,
		ContractAddress: address,
		Initiator:       initiator,
		Talliers:        talliers,
		Voters:          voters,
		Config:          stakeCfg,
	}, nil
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

func (c *StakeChain) SettleRewards(voterIDs []int, tallierIDs []int) error {
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
	call := &bind.CallOpts{Context: context.Background()}
	data, err := c.Contract.GetEscrowOverview(call)
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
	call := &bind.CallOpts{Context: context.Background()}
	data, err := c.Contract.GetInitiatorState(call)
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
	return &ChainParticipant{
		Address:       address,
		Deposited:     data.Deposited,
		Claimable:     data.Claimable,
		WalletBalance: balance,
		Staked:        data.Staked,
		Honest:        true,
		Withdrawn:     data.Withdrawn,
	}, nil
}

func (c *StakeChain) ReadVoter(voterID int) (*ChainParticipant, error) {
	call := &bind.CallOpts{Context: context.Background()}
	data, err := c.Contract.GetVoter(call, big.NewInt(int64(voterID)))
	if err != nil {
		return nil, err
	}

	address := data.Account
	balance := big.NewInt(0)
	if address != (common.Address{}) {
		balance, err = c.Client.BalanceAt(context.Background(), address, nil)
		if err != nil {
			return nil, err
		}
	}

	return &ChainParticipant{
		Address:       address,
		Deposited:     data.Deposited,
		Claimable:     data.Claimable,
		WalletBalance: balance,
		Staked:        data.Staked,
		Honest:        data.Honest,
		Withdrawn:     data.Withdrawn,
	}, nil
}

func (c *StakeChain) ReadTallier(tallierID int) (*ChainParticipant, error) {
	call := &bind.CallOpts{Context: context.Background()}
	data, err := c.Contract.GetTallier(call, big.NewInt(int64(tallierID)))
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

	return &ChainParticipant{
		Address:       address,
		Deposited:     data.Deposited,
		Claimable:     data.Claimable,
		WalletBalance: balance,
		Staked:        data.Staked,
		Honest:        data.Honest,
		Withdrawn:     data.Withdrawn,
	}, nil
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

	return StakeConfig{
		InitiatorEscrowWei: initiatorEscrow,
		VoterStakeWei:      voterStake,
		TallierStakeWei:    tallierStake,
		InitiatorRewardPct: uint8(cfg.InitiatorRewardPercent),
		VoterRewardPct:     uint8(cfg.VoterRewardPercent),
		TallierRewardPct:   uint8(cfg.TallierRewardPercent),
	}, nil
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

	privateKeys := make([]string, 0, 16)
	for i := 1; i <= 32; i++ {
		key := strings.TrimSpace(env[fmt.Sprintf("PRIVATE_KEY_%d", i)])
		if key == "" {
			continue
		}
		privateKeys = append(privateKeys, key)
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
	return RoleAccount{
		PrivateKey: privateKeyHex,
		Address:    crypto.PubkeyToAddress(key.PublicKey),
	}, nil
}

func newAuth(client *ethclient.Client, privateKeyHex string, value *big.Int) (*bind.TransactOpts, error) {
	key, err := crypto.HexToECDSA(strings.TrimSpace(privateKeyHex))
	if err != nil {
		return nil, err
	}
	publicKeyECDSA, ok := key.Public().(*ecdsa.PublicKey)
	if !ok {
		return nil, errorsf("failed to parse ECDSA public key")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
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

	balance, err := client.BalanceAt(context.Background(), fromAddress, nil)
	if err != nil {
		return nil, err
	}
	requiredGas := new(big.Int).Mul(new(big.Int).SetUint64(stakeTxGasLimit), gasPrice)
	required := new(big.Int).Add(requiredGas, value)
	if balance.Cmp(required) < 0 {
		return nil, fmt.Errorf(
			"ganache account %s has %s ETH, needs at least %s ETH for value plus reserved gas; start Ganache with `ganache --mnemonic \"PVote\" -l 90071992547 -e 1000` or regenerate .env for your current Ganache accounts",
			fromAddress.Hex(),
			formatWeiToETH(balance),
			formatWeiToETH(required),
		)
	}
	return auth, nil
}

func waitForTx(client *ethclient.Client, tx *types.Transaction) error {
	receipt, err := bind.WaitMined(context.Background(), client, tx)
	if err != nil {
		return err
	}
	if receipt.Status != types.ReceiptStatusSuccessful {
		return errorsf("transaction %s reverted", tx.Hash().Hex())
	}
	return nil
}

func intSliceToBigInts(values []int) []*big.Int {
	result := make([]*big.Int, len(values))
	for i, value := range values {
		result[i] = big.NewInt(int64(value))
	}
	return result
}

func errorsf(format string, args ...any) error {
	return fmt.Errorf(format, args...)
}

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
