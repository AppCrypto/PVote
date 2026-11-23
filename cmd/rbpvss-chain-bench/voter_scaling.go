package main

import (
	"compress/gzip"
	"context"
	"crypto/ecdsa"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"time"

	chain "PVote/chain/rbpvss"
	contract "PVote/compile/contract"
	core "PVote/crypto/RBPVSS"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

// Ganache's wallet file stays outside the results; only public transactions,
// receipts, account addresses, and chain metadata are archived.
func voterScalingKeys(csvValue, walletPath string, count int) ([]*ecdsa.PrivateKey, error) {
	if count < 0 || (walletPath != "" && strings.TrimSpace(csvValue) != "") {
		return nil, fmt.Errorf("use one account source and a nonnegative voter count")
	}
	var keys []*ecdsa.PrivateKey
	if walletPath == "" {
		var err error
		keys, err = parsePrivateKeys(csvValue)
		if err != nil {
			return nil, err
		}
	} else {
		data, err := os.ReadFile(walletPath)
		if err != nil {
			return nil, err
		}
		var wallet struct {
			PrivateKeys map[string]string `json:"private_keys"`
		}
		if err := json.Unmarshal(data, &wallet); err != nil {
			return nil, fmt.Errorf("read Ganache development wallet: %w", err)
		}
		addresses := make([]string, 0, len(wallet.PrivateKeys))
		for address := range wallet.PrivateKeys {
			addresses = append(addresses, address)
		}
		sort.Strings(addresses)
		for _, address := range addresses {
			key, err := crypto.HexToECDSA(strings.TrimPrefix(wallet.PrivateKeys[address], "0x"))
			if err != nil || !common.IsHexAddress(address) {
				return nil, fmt.Errorf("invalid account in development wallet")
			}
			if crypto.PubkeyToAddress(key.PublicKey) != common.HexToAddress(address) {
				return nil, fmt.Errorf("development wallet key/address mismatch")
			}
			keys = append(keys, key)
		}
	}
	if count == 0 {
		count = len(keys)
	}
	if count == 0 || count > len(keys) {
		return nil, fmt.Errorf("requested %d voters but supplied %d accounts", count, len(keys))
	}
	keys = keys[:count]
	seen := make(map[common.Address]bool, count)
	for _, key := range keys {
		address := crypto.PubkeyToAddress(key.PublicKey)
		if seen[address] {
			return nil, fmt.Errorf("voter accounts must be distinct")
		}
		seen[address] = true
	}
	return keys, nil
}

type scalingMetadata struct {
	StartedUTC          string   `json:"started_utc"`
	FinishedUTC         string   `json:"finished_utc,omitempty"`
	Completed           bool     `json:"completed"`
	ChainID             string   `json:"chain_id"`
	ClientVersion       string   `json:"client_version"`
	ExecutionHardfork   string   `json:"configured_execution_hardfork"`
	BlockGasLimit       uint64   `json:"block_gas_limit"`
	TransactionGasLimit uint64   `json:"transaction_gas_limit"`
	GoVersion           string   `json:"go_version"`
	Platform            string   `json:"platform"`
	VotersPerConfig     int      `json:"voters_per_configuration"`
	Configs             [][2]int `json:"n_l_configurations"`
	BytecodeKeccak      string   `json:"deployment_bytecode_keccak256"`
	ScorePolicy         string   `json:"score_policy"`
	VerifierScope       string   `json:"verifier_scope"`
}

func runVoterScaling(client *ethclient.Client, keys []*ecdsa.PrivateKey, configs [][2]int, path, hardfork string) (err error) {
	ctx := context.Background()
	chainID, err := client.ChainID(ctx)
	if err != nil {
		return err
	}
	header, err := client.HeaderByNumber(ctx, nil)
	if err != nil {
		return err
	}
	var version string
	if err := client.Client().CallContext(ctx, &version, "web3_clientVersion"); err != nil {
		return err
	}
	metadata := scalingMetadata{
		StartedUTC: time.Now().UTC().Format(time.RFC3339), ChainID: chainID.String(), ClientVersion: version,
		ExecutionHardfork: hardfork, BlockGasLimit: header.GasLimit, TransactionGasLimit: transactionGasLimit,
		GoVersion: runtime.Version(), Platform: runtime.GOOS + "/" + runtime.GOARCH,
		VotersPerConfig: len(keys), Configs: configs,
		BytecodeKeccak: crypto.Keccak256Hash(common.FromHex(contract.ContractMetaData.Bin)).Hex(),
		ScorePolicy:    "fresh independent transcript per voter; all scores 0 in [0,10]",
		VerifierScope:  "checked-in Solidity benchmark verifier (legacy transcript format); not a context-binding security test",
	}
	writeMetadata := func() error {
		data, err := json.MarshalIndent(metadata, "", "  ")
		if err != nil {
			return err
		}
		return os.WriteFile(path+".metadata.json", append(data, '\n'), 0644)
	}
	if err := writeMetadata(); err != nil {
		return err
	}
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	writer := csv.NewWriter(file)
	if err := writer.Write([]string{"n", "l", "accepted_count", "chain_id", "client", "go_version", "submit_rb_gas", "t", "contract_address", "voter_address", "transaction_hash", "block_number", "receipt_status", "calldata_bytes", "transaction_gas_limit", "block_gas_limit"}); err != nil {
		return err
	}
	archive, err := os.Create(path + ".receipts.jsonl.gz")
	if err != nil {
		return err
	}
	defer archive.Close()
	compressed := gzip.NewWriter(archive)
	defer func() {
		if closeErr := compressed.Close(); err == nil {
			err = closeErr
		}
	}()
	encoder := json.NewEncoder(compressed)
	for _, config := range configs {
		if err := measureVoterScaling(client, keys, config[0], config[1], chainID, writer, encoder, metadata); err != nil {
			return fmt.Errorf("n=%d,l=%d: %w", config[0], config[1], err)
		}
	}
	metadata.Completed = true
	metadata.FinishedUTC = time.Now().UTC().Format(time.RFC3339)
	return writeMetadata()
}

func measureVoterScaling(client *ethclient.Client, keys []*ecdsa.PrivateKey, n, l int, chainID *big.Int, writer *csv.Writer, archive *json.Encoder, metadata scalingMetadata) error {
	t := (n + l) / 2
	pp, _, err := core.Setup(128, n, t, l, core.Interval{Min: 0, Max: 10})
	if err != nil {
		return err
	}
	accounts := make([]common.Address, len(keys))
	for i, key := range keys {
		accounts[i] = crypto.PubkeyToAddress(key.PublicKey)
	}
	deployer, err := bind.NewKeyedTransactorWithChainID(keys[0], chainID)
	if err != nil {
		return err
	}
	deployer.GasLimit = transactionGasLimit
	g0, h0, g1, pkI, sigma, pks, threshold, coordinates, minimum, maximum, err := chain.ConstructorArguments(pp)
	if err != nil {
		return err
	}
	address, tx, verifier, err := contract.DeployContract(deployer, client, g0, h0, g1, pkI, sigma, pks, threshold, coordinates, minimum, maximum, accounts)
	if err != nil {
		return err
	}
	record := func(phase string, count int, tx *types.Transaction) (*types.Receipt, error) {
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()
		receipt, err := bind.WaitMined(ctx, client, tx)
		if err != nil {
			return nil, err
		}
		if err := archive.Encode(struct {
			Phase         string             `json:"phase"`
			N             int                `json:"n"`
			L             int                `json:"l"`
			AcceptedCount int                `json:"accepted_count"`
			Transaction   *types.Transaction `json:"transaction"`
			Receipt       *types.Receipt     `json:"receipt"`
		}{phase, n, l, count, tx, receipt}); err != nil {
			return nil, err
		}
		if receipt.Status != types.ReceiptStatusSuccessful {
			return nil, fmt.Errorf("%s %s reverted after %d gas", phase, tx.Hash(), receipt.GasUsed)
		}
		return receipt, nil
	}
	if _, err := record("deployment", 0, tx); err != nil {
		return err
	}
	started := time.Now()
	fmt.Printf("start n=%d l=%d t=%d voters=%d contract=%s\n", n, l, t, len(keys), address.Hex())
	for i, key := range keys {
		instance, err := core.Share(pp, make([]int64, l))
		if err != nil {
			return err
		}
		dealer, rangeBinding, err := chain.Transcript(instance)
		if err != nil {
			return err
		}
		auth, err := bind.NewKeyedTransactorWithChainID(key, chainID)
		if err != nil {
			return err
		}
		auth.GasLimit = transactionGasLimit
		tx, err := verifier.SubmitRB(auth, dealer, rangeBinding)
		if err != nil {
			return fmt.Errorf("submit %d: %w", i+1, err)
		}
		receipt, err := record("submitRB", i+1, tx)
		if err != nil {
			return err
		}
		// Validate the emitted on-chain count, not just the loop counter.
		matched := false
		for _, log := range receipt.Logs {
			if log.Address != address {
				continue
			}
			event, err := verifier.ParseTranscriptAccepted(*log)
			if err == nil && event.AcceptedCount.Cmp(big.NewInt(int64(i+1))) == 0 {
				matched = true
			}
		}
		if !matched {
			return fmt.Errorf("receipt %s missing expected accepted count %d", tx.Hash(), i+1)
		}
		row := []string{
			strconv.Itoa(n), strconv.Itoa(l), strconv.Itoa(i + 1), chainID.String(), metadata.ClientVersion, runtime.Version(),
			strconv.FormatUint(receipt.GasUsed, 10), strconv.Itoa(t), address.Hex(), accounts[i].Hex(), tx.Hash().Hex(),
			receipt.BlockNumber.String(), strconv.FormatUint(receipt.Status, 10), strconv.Itoa(len(tx.Data())),
			strconv.FormatUint(tx.Gas(), 10), strconv.FormatUint(metadata.BlockGasLimit, 10),
		}
		if err := writer.Write(row); err != nil {
			return err
		}
		writer.Flush() // Keep actual measurements even if a long run is interrupted.
		if err := writer.Error(); err != nil {
			return err
		}
		if i == 0 || (i+1)%25 == 0 || i+1 == len(keys) {
			fmt.Printf("n=%d l=%d accepted=%d/%d gas=%d elapsed=%s\n", n, l, i+1, len(keys), receipt.GasUsed, time.Since(started).Round(time.Second))
		}
	}
	count, err := verifier.AcceptedCount(nil)
	if err != nil || count == nil || count.Cmp(big.NewInt(int64(len(keys)))) != 0 {
		return fmt.Errorf("final on-chain accepted count %v, expected %d (error %v)", count, len(keys), err)
	}
	fmt.Printf("verified n=%d l=%d acceptedCount=%s\n", n, l, count)
	return nil
}
