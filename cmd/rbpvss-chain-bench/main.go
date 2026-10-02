package main

import (
	"context"
	"crypto/ecdsa"
	"encoding/csv"
	"flag"
	"fmt"
	"math"
	"math/big"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	chain "PVote/chain/rbpvss"
	contract "PVote/compile/contract"
	core "PVote/crypto/RBPVSS"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

const transactionGasLimit uint64 = 180_000_000

type result struct {
	n, l, t                       int
	deployment, submit, shareMean uint64
	shareStd                      float64
	reconstruct                   uint64
}

type voterScalingResult struct {
	n, l, acceptedCount int
	gas                 uint64
}

func main() {
	url := flag.String("url", envOr("PVOTE_EVM_URL", "http://127.0.0.1:8545"), "EVM JSON-RPC URL")
	privateKeyHex := flag.String("private-key", envOr("PVOTE_EVM_PRIVATE_KEY", ""), "funded development account key")
	privateKeysCSV := flag.String("private-keys", envOr("PVOTE_EVM_PRIVATE_KEYS", ""), "comma-separated funded development account keys for voter-count validation")
	out := flag.String("out", "experiments/rbpvss_chain_results.csv", "CSV output path")
	voterScalingOut := flag.String("voter-scaling-out", "experiments/voter_gas_by_m.csv", "voter-count validation CSV output path")
	onlyVoterScaling := flag.Bool("only-voter-scaling", false, "run only the fixed-(n,l) voter-count validation")
	voterScalingConfigs := flag.String("voter-scaling-configs", "10:1", "comma-separated n:l configurations for voter-count validation")
	flag.Parse()
	if *privateKeyHex == "" {
		panic("provide -private-key or PVOTE_EVM_PRIVATE_KEY")
	}
	client, err := ethclient.Dial(*url)
	if err != nil {
		panic(err)
	}
	defer client.Close()
	key, err := crypto.HexToECDSA(strings.TrimPrefix(*privateKeyHex, "0x"))
	if err != nil {
		panic(err)
	}
	if !*onlyVoterScaling {
		var rows []result
		for _, l := range []int{1, 10, 19} {
			for n := 10; n <= 100; n += 10 {
				if l > n/3 {
					continue
				}
				row, err := benchmarkPair(client, key, n, l)
				if err != nil {
					panic(fmt.Errorf("n=%d,l=%d: %w", n, l, err))
				}
				rows = append(rows, row)
				fmt.Printf("n=%d l=%d deploy=%d submitRB=%d share=%.0f reconstruct=%d\n", n, l, row.deployment, row.submit, float64(row.shareMean), row.reconstruct)
			}
		}
		if err := writeCSV(*out, rows); err != nil {
			panic(err)
		}
	}
	if strings.TrimSpace(*privateKeysCSV) != "" {
		keys, err := parsePrivateKeys(*privateKeysCSV)
		if err != nil {
			panic(err)
		}
		configs, err := parseVoterScalingConfigs(*voterScalingConfigs)
		if err != nil {
			panic(err)
		}
		var rows []voterScalingResult
		for _, config := range configs {
			configRows, err := benchmarkVoterScaling(client, keys, config[0], config[1])
			if err != nil {
				panic(fmt.Errorf("voter scaling n=%d,l=%d: %w", config[0], config[1], err))
			}
			rows = append(rows, configRows...)
		}
		if err := writeVoterScalingCSV(*voterScalingOut, rows); err != nil {
			panic(err)
		}
		for _, row := range rows {
			fmt.Printf("fixed n=%d l=%d accepted_count=%d voter_gas=%d\n", row.n, row.l, row.acceptedCount, row.gas)
		}
	} else if *onlyVoterScaling {
		panic("-only-voter-scaling requires -private-keys or PVOTE_EVM_PRIVATE_KEYS")
	}
}

func parseVoterScalingConfigs(value string) ([][2]int, error) {
	var configs [][2]int
	for _, raw := range strings.Split(value, ",") {
		parts := strings.Split(strings.TrimSpace(raw), ":")
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid voter-scaling configuration %q; want n:l", raw)
		}
		n, err := strconv.Atoi(parts[0])
		if err != nil {
			return nil, fmt.Errorf("invalid n in voter-scaling configuration %q: %w", raw, err)
		}
		l, err := strconv.Atoi(parts[1])
		if err != nil {
			return nil, fmt.Errorf("invalid l in voter-scaling configuration %q: %w", raw, err)
		}
		if n <= 0 || l <= 0 || l >= (n+l)/2 {
			return nil, fmt.Errorf("invalid voter-scaling configuration %q: require n>0 and 0<l<t", raw)
		}
		configs = append(configs, [2]int{n, l})
	}
	if len(configs) == 0 {
		return nil, fmt.Errorf("at least one voter-scaling configuration is required")
	}
	return configs, nil
}

func parsePrivateKeys(csvValue string) ([]*ecdsa.PrivateKey, error) {
	parts := strings.Split(csvValue, ",")
	keys := make([]*ecdsa.PrivateKey, 0, len(parts))
	for i, part := range parts {
		value := strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(part), "0x"))
		if value == "" {
			continue
		}
		key, err := crypto.HexToECDSA(value)
		if err != nil {
			return nil, fmt.Errorf("private key %d: %w", i+1, err)
		}
		keys = append(keys, key)
	}
	if len(keys) < 2 {
		return nil, fmt.Errorf("voter-count validation requires at least two private keys")
	}
	return keys, nil
}

func benchmarkVoterScaling(client *ethclient.Client, keys []*ecdsa.PrivateKey, n, l int) ([]voterScalingResult, error) {
	t := (n + l) / 2
	pp, _, err := core.Setup(128, n, t, l, core.Interval{Min: 0, Max: 10})
	if err != nil {
		return nil, err
	}
	voterAccounts := make([]common.Address, len(keys))
	for i, key := range keys {
		voterAccounts[i] = crypto.PubkeyToAddress(key.PublicKey)
	}
	deployer, err := bind.NewKeyedTransactorWithChainID(keys[0], big.NewInt(1337))
	if err != nil {
		return nil, err
	}
	deployer.GasLimit = transactionGasLimit
	g0, h0, g1, pkI, sigma, pks, threshold, coordinates, minimum, maximum, err := chain.ConstructorArguments(pp)
	if err != nil {
		return nil, err
	}
	_, tx, verifier, err := contract.DeployContract(deployer, client, g0, h0, g1, pkI, sigma, pks, threshold, coordinates, minimum, maximum, voterAccounts)
	if err != nil {
		return nil, err
	}
	if _, err := successfulReceipt(client, tx); err != nil {
		return nil, err
	}
	rows := make([]voterScalingResult, 0, len(keys))
	for i, key := range keys {
		instance, err := core.Share(pp, make([]int64, l))
		if err != nil {
			return nil, err
		}
		dealer, rangeBinding, err := chain.Transcript(instance)
		if err != nil {
			return nil, err
		}
		auth, err := bind.NewKeyedTransactorWithChainID(key, big.NewInt(1337))
		if err != nil {
			return nil, err
		}
		auth.GasLimit = transactionGasLimit
		tx, err := verifier.SubmitRB(auth, dealer, rangeBinding)
		if err != nil {
			return nil, fmt.Errorf("submit voter %d: %w", i+1, err)
		}
		gas, err := successfulReceipt(client, tx)
		if err != nil {
			return nil, fmt.Errorf("submit voter %d: %w", i+1, err)
		}
		rows = append(rows, voterScalingResult{n: n, l: l, acceptedCount: i + 1, gas: gas})
	}
	return rows, nil
}

func benchmarkPair(client *ethclient.Client, key *ecdsa.PrivateKey, n, l int) (result, error) {
	t := (n + l) / 2
	pp, secretKeys, err := core.Setup(128, n, t, l, core.Interval{Min: 0, Max: 10})
	if err != nil {
		return result{}, err
	}
	auth, err := bind.NewKeyedTransactorWithChainID(key, big.NewInt(1337))
	if err != nil {
		return result{}, err
	}
	auth.GasLimit = transactionGasLimit
	g0, h0, g1, pkI, sigma, pks, threshold, coordinates, minimum, maximum, err := chain.ConstructorArguments(pp)
	if err != nil {
		return result{}, err
	}
	address, tx, verifier, err := contract.DeployContract(auth, client, g0, h0, g1, pkI, sigma, pks, threshold, coordinates, minimum, maximum, []common.Address{auth.From})
	if err != nil {
		return result{}, err
	}
	deployment, err := successfulReceipt(client, tx)
	if err != nil {
		return result{}, err
	}
	scores := make([]int64, l)
	for i := range scores {
		scores[i] = 5
	}
	instance, err := core.Share(pp, scores)
	if err != nil {
		return result{}, err
	}
	dealer, rangeBinding, err := chain.Transcript(instance)
	if err != nil {
		return result{}, err
	}
	tx, err = verifier.SubmitRB(auth, dealer, rangeBinding)
	if err != nil {
		return result{}, err
	}
	submit, err := successfulReceipt(client, tx)
	if err != nil {
		return result{}, err
	}
	shareGas := make([]uint64, t)
	for i := 0; i < t; i++ {
		share, proof, err := core.Decrypt(pp, instance.C[i], secretKeys[i])
		if err != nil {
			return result{}, err
		}
		tx, err = verifier.SubmitDecryptionShare(auth, big.NewInt(int64(i+1)), chain.G1Point(share), chain.G1Point(proof.A), chain.G1Point(proof.B), proof.Chi, proof.Z)
		if err != nil {
			return result{}, err
		}
		shareGas[i], err = successfulReceipt(client, tx)
		if err != nil {
			return result{}, err
		}
	}
	abi, err := contract.ContractMetaData.GetAbi()
	if err != nil {
		return result{}, err
	}
	data, err := abi.Pack("reconstruct")
	if err != nil {
		return result{}, err
	}
	reconstruct, err := client.EstimateGas(context.Background(), ethereum.CallMsg{From: auth.From, To: &address, Data: data, Gas: transactionGasLimit})
	if err != nil {
		return result{}, err
	}
	mean, std := gasStats(shareGas)
	return result{n: n, l: l, t: t, deployment: deployment, submit: submit, shareMean: uint64(math.Round(mean)), shareStd: std, reconstruct: reconstruct}, nil
}

func successfulReceipt(client *ethclient.Client, tx *types.Transaction) (uint64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	receipt, err := bind.WaitMined(ctx, client, tx)
	if err != nil {
		return 0, err
	}
	if receipt.Status != types.ReceiptStatusSuccessful {
		return 0, fmt.Errorf("transaction %s reverted after %d gas", tx.Hash(), receipt.GasUsed)
	}
	return receipt.GasUsed, nil
}

func gasStats(values []uint64) (float64, float64) {
	var mean float64
	for _, value := range values {
		mean += float64(value)
	}
	mean /= float64(len(values))
	if len(values) == 1 {
		return mean, 0
	}
	var variance float64
	for _, value := range values {
		delta := float64(value) - mean
		variance += delta * delta
	}
	return mean, math.Sqrt(variance / float64(len(values)-1))
}

func writeCSV(path string, rows []result) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	writer := csv.NewWriter(file)
	defer writer.Flush()
	if err := writer.Write([]string{"n", "l", "t", "chain_id", "client", "go_version", "deployment_gas", "submit_rb_gas", "decryption_share_mean_gas", "decryption_share_std_gas", "reconstruct_estimated_gas"}); err != nil {
		return err
	}
	for _, row := range rows {
		if err := writer.Write([]string{
			strconv.Itoa(row.n), strconv.Itoa(row.l), strconv.Itoa(row.t), "1337", "ganache-v7.9.2", runtime.Version(),
			strconv.FormatUint(row.deployment, 10), strconv.FormatUint(row.submit, 10), strconv.FormatUint(row.shareMean, 10),
			strconv.FormatFloat(row.shareStd, 'f', 3, 64), strconv.FormatUint(row.reconstruct, 10),
		}); err != nil {
			return err
		}
	}
	return writer.Error()
}

func writeVoterScalingCSV(path string, rows []voterScalingResult) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	writer := csv.NewWriter(file)
	defer writer.Flush()
	if err := writer.Write([]string{"n", "l", "accepted_count", "chain_id", "client", "go_version", "submit_rb_gas"}); err != nil {
		return err
	}
	for _, row := range rows {
		if err := writer.Write([]string{
			strconv.Itoa(row.n), strconv.Itoa(row.l), strconv.Itoa(row.acceptedCount), "1337", "ganache-v7.9.2", runtime.Version(), strconv.FormatUint(row.gas, 10),
		}); err != nil {
			return err
		}
	}
	return writer.Error()
}

func envOr(name, fallback string) string {
	if value := os.Getenv(name); value != "" {
		return value
	}
	return fallback
}
