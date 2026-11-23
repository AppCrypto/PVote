// rbpvss-gas-audit executes signed transactions on an in-memory Geth EVM chain.
// It uses the checked-in Solidity bytecode without changing verifier logic.
package main

import (
	"compress/gzip"
	"crypto/ecdsa"
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	chain "PVote/chain/rbpvss"
	contract "PVote/compile/contract"
	protocol "PVote/crypto/RBPVSS"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus/ethash"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	vmruntime "github.com/ethereum/go-ethereum/core/vm/runtime"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/params"
)

const blockGasLimit = uint64(90_071_992_547)
const txGasLimit = uint64(180_000_000)

type runMetadata struct {
	StartedUTC     string              `json:"started_utc"`
	FinishedUTC    string              `json:"finished_utc,omitempty"`
	Completed      bool                `json:"completed"`
	Backend        string              `json:"backend"`
	Execution      string              `json:"execution"`
	Config         *params.ChainConfig `json:"chain_config"`
	BlockGasLimit  uint64              `json:"block_gas_limit"`
	TxGasLimit     uint64              `json:"transaction_gas_limit"`
	Voters         int                 `json:"voters_per_configuration_per_run"`
	Runs           int                 `json:"runs_per_configuration"`
	Configurations [][2]int            `json:"n_l_configurations"`
	GoVersion      string              `json:"go_version"`
	Platform       string              `json:"platform"`
	BytecodeHash   string              `json:"deployment_bytecode_keccak256"`
	Scope          string              `json:"scope"`
}

func main() {
	voters := flag.Int("voters", 2, "actual consecutive distinct accepted voters per run")
	runs := flag.Int("runs", 3, "independent deployments per configuration")
	configString := flag.String("configs", "10:5,50:10,100:19", "comma-separated n:l configurations")
	out := flag.String("out", "experiments/gas-audit", "fresh output directory")
	flag.Parse()
	if err := execute(*out, *voters, *runs, *configString); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func configurations(value string) ([][2]int, error) {
	var result [][2]int
	for _, raw := range strings.Split(value, ",") {
		var n, l int
		parts := strings.Split(raw, ":")
		if len(parts) != 2 {
			return nil, fmt.Errorf("expected n:l, got %q", raw)
		}
		n, err := strconv.Atoi(parts[0])
		if err != nil {
			return nil, err
		}
		l, err = strconv.Atoi(parts[1])
		if err != nil {
			return nil, err
		}
		if l <= 0 || l >= (n+l)/2 {
			return nil, fmt.Errorf("invalid n:l %q", raw)
		}
		result = append(result, [2]int{n, l})
	}
	return result, nil
}

func execute(out string, voters, runs int, configString string) error {
	if voters < 2 || runs < 1 {
		return fmt.Errorf("require voters>=2 and runs>=1")
	}
	configs, err := configurations(configString)
	if err != nil {
		return err
	}
	if _, err := os.Stat(filepath.Join(out, "metadata.json")); err == nil {
		return fmt.Errorf("output exists: choose a new directory")
	}
	if err := os.MkdirAll(out, 0755); err != nil {
		return err
	}
	config := *params.AllEthashProtocolChanges
	config.ChainID = big.NewInt(1337)
	metadata := runMetadata{
		StartedUTC: time.Now().UTC().Format(time.RFC3339),
		Backend:    "go-ethereum v1.14.6 native EVM / core.GenerateChain",
		Execution:  "in-memory London-rule chain; actual signed transactions, execution receipts and state transitions; fake PoW consensus, no live network",
		Config:     &config, BlockGasLimit: blockGasLimit, TxGasLimit: txGasLimit, Voters: voters, Runs: runs, Configurations: configs,
		GoVersion: runtime.Version(), Platform: runtime.GOOS + "/" + runtime.GOARCH,
		BytecodeHash: crypto.Keccak256Hash(common.FromHex(contract.ContractMetaData.Bin)).Hex(),
		Scope:        "legacy Solidity transcript format; gas audit, not context-binding security validation; permissive local block gas limit, not mainnet deployment feasibility",
	}
	writeMetadata := func() error {
		data, e := json.MarshalIndent(metadata, "", "  ")
		if e != nil {
			return e
		}
		return os.WriteFile(filepath.Join(out, "metadata.json"), append(data, '\n'), 0644)
	}
	if err := writeMetadata(); err != nil {
		return err
	}
	file, err := os.Create(filepath.Join(out, "measurements.csv"))
	if err != nil {
		return err
	}
	defer file.Close()
	csvWriter := csv.NewWriter(file)
	if err := csvWriter.Write([]string{"n", "l", "t", "run", "accepted_count", "gas_used", "transaction_hash", "block_number", "contract_address", "voter_address", "calldata_bytes", "calldata_zero_bytes", "calldata_nonzero_bytes", "traced"}); err != nil {
		return err
	}
	archiveFile, err := os.Create(filepath.Join(out, "transactions.jsonl.gz"))
	if err != nil {
		return err
	}
	defer archiveFile.Close()
	archive := gzip.NewWriter(archiveFile)
	defer archive.Close()
	archiveWriter := json.NewEncoder(archive)
	detailFile, err := os.Create(filepath.Join(out, "breakdown.jsonl"))
	if err != nil {
		return err
	}
	defer detailFile.Close()
	detailWriter := json.NewEncoder(detailFile)
	for _, c := range configs {
		for run := 1; run <= runs; run++ {
			if err := executeConfiguration(&config, c[0], c[1], voters, run, csvWriter, archiveWriter, detailWriter); err != nil {
				return err
			}
		}
	}
	csvWriter.Flush()
	if err := csvWriter.Error(); err != nil {
		return err
	}
	if err := archive.Close(); err != nil {
		return err
	}
	metadata.Completed = true
	metadata.FinishedUTC = time.Now().UTC().Format(time.RFC3339)
	return writeMetadata()
}

func executeConfiguration(config *params.ChainConfig, n, l, voters, run int, writer *csv.Writer, archive, details *json.Encoder) error {
	t := (n + l) / 2
	pp, _, err := protocol.Setup(128, n, t, l, protocol.Interval{Min: 0, Max: 10})
	if err != nil {
		return err
	}
	keys := make([]*ecdsa.PrivateKey, voters)
	accounts := make([]common.Address, voters)
	alloc := make(types.GenesisAlloc, voters)
	for i := range keys {
		// Public deterministic keys for this disposable local fixture only.
		key, e := crypto.ToECDSA(big.NewInt(int64(i + 1)).FillBytes(make([]byte, 32)))
		if e != nil {
			return e
		}
		keys[i] = key
		accounts[i] = crypto.PubkeyToAddress(key.PublicKey)
		alloc[accounts[i]] = types.Account{Balance: new(big.Int).Mul(big.NewInt(1000), big.NewInt(params.Ether))}
	}
	genesis := &core.Genesis{Config: config, Alloc: alloc, GasLimit: blockGasLimit, BaseFee: big.NewInt(params.InitialBaseFee)}
	engine := ethash.NewFaker()
	defer engine.Close()
	db, blocks, _ := core.GenerateChainWithGenesis(genesis, engine, 1, nil)
	defer db.Close()
	parent := blocks[0]
	parsed, err := contract.ContractMetaData.GetAbi()
	if err != nil {
		return err
	}
	g0, h0, g1, pkI, sigma, pks, threshold, coordinates, min, max, err := chain.ConstructorArguments(pp)
	if err != nil {
		return err
	}
	args, err := parsed.Pack("", g0, h0, g1, pkI, sigma, pks, threshold, coordinates, min, max, accounts)
	if err != nil {
		return err
	}
	nonces := make(map[common.Address]uint64)
	submit := func(key *ecdsa.PrivateKey, to *common.Address, data []byte, tracer *gasAudit) (*types.Transaction, *types.Receipt, error) {
		from := crypto.PubkeyToAddress(key.PublicKey)
		tx := types.NewTx(&types.LegacyTx{Nonce: nonces[from], To: to, Gas: txGasLimit, GasPrice: big.NewInt(2_000_000_000), Value: new(big.Int), Data: data})
		tx, e := types.SignTx(tx, types.LatestSigner(config), key)
		if e != nil {
			return nil, nil, e
		}
		var vmConfig vm.Config
		if tracer != nil {
			vmConfig.Tracer = tracer.hooks()
		}
		next, receipts := core.GenerateChain(config, parent, engine, db, 1, func(_ int, b *core.BlockGen) { b.AddTxWithVMConfig(tx, vmConfig) })
		parent = next[0]
		receipt := receipts[0][0]
		nonces[from]++
		if receipt.Status != types.ReceiptStatusSuccessful {
			return tx, receipt, fmt.Errorf("transaction %s reverted, gas=%d", tx.Hash(), receipt.GasUsed)
		}
		return tx, receipt, nil
	}
	deployment := append(common.FromHex(contract.ContractMetaData.Bin), args...)
	tx, receipt, err := submit(keys[0], nil, deployment, nil)
	if err != nil {
		return err
	}
	address := receipt.ContractAddress
	writeArchive := func(phase string, count int, tx *types.Transaction, receipt *types.Receipt) error {
		return archive.Encode(struct {
			Phase         string             `json:"phase"`
			N             int                `json:"n"`
			L             int                `json:"l"`
			Run           int                `json:"run"`
			AcceptedCount int                `json:"accepted_count"`
			Transaction   *types.Transaction `json:"transaction"`
			Receipt       *types.Receipt     `json:"receipt"`
		}{phase, n, l, run, count, tx, receipt})
	}
	if err := writeArchive("deployment", 0, tx, receipt); err != nil {
		return err
	}
	started := time.Now()
	fmt.Printf("start native-EVM n=%d l=%d t=%d run=%d voters=%d\n", n, l, t, run, voters)
	for i, key := range keys {
		instance, e := protocol.Share(pp, make([]int64, l))
		if e != nil {
			return e
		}
		dealer, ranges, e := chain.Transcript(instance)
		if e != nil {
			return e
		}
		data, e := parsed.Pack("submitRB", dealer, ranges)
		if e != nil {
			return e
		}
		var tracer *gasAudit
		if i < 2 || i+1 == voters {
			tracer = newGasAudit(n, l, t, run, i+1, data)
		}
		tx, receipt, e := submit(key, &address, data, tracer)
		if e != nil {
			return e
		}
		if e := checkAcceptedEvent(parsed, address, receipt, i+1); e != nil {
			return e
		}
		if e := writeArchive("submitRB", i+1, tx, receipt); e != nil {
			return e
		}
		if tracer != nil {
			if e := tracer.finish(receipt); e != nil {
				return e
			}
			if e := details.Encode(tracer); e != nil {
				return e
			}
		}
		zero := 0
		for _, v := range data {
			if v == 0 {
				zero++
			}
		}
		row := []string{strconv.Itoa(n), strconv.Itoa(l), strconv.Itoa(t), strconv.Itoa(run), strconv.Itoa(i + 1), strconv.FormatUint(receipt.GasUsed, 10), tx.Hash().Hex(), receipt.BlockNumber.String(), address.Hex(), accounts[i].Hex(), strconv.Itoa(len(data)), strconv.Itoa(zero), strconv.Itoa(len(data) - zero), strconv.FormatBool(tracer != nil)}
		if e := writer.Write(row); e != nil {
			return e
		}
		writer.Flush()
		if e := writer.Error(); e != nil {
			return e
		}
		if i == 0 || (i+1)%50 == 0 || i+1 == voters {
			fmt.Printf("n=%d l=%d run=%d accepted=%d/%d gas=%d elapsed=%s\n", n, l, run, i+1, voters, receipt.GasUsed, time.Since(started).Round(time.Second))
		}
	}
	// Read the actual getter from the final EVM state, independently of events.
	st, err := state.New(parent.Root(), state.NewDatabase(db), nil)
	if err != nil {
		return err
	}
	getter, err := parsed.Pack("acceptedCount")
	if err != nil {
		return err
	}
	output, _, err := vmruntime.Call(address, getter, &vmruntime.Config{State: st, ChainConfig: config, BlockNumber: parent.Number(), GasLimit: txGasLimit})
	if err != nil {
		return err
	}
	values, err := parsed.Unpack("acceptedCount", output)
	if err != nil {
		return err
	}
	if values[0].(*big.Int).Cmp(big.NewInt(int64(voters))) != 0 {
		return fmt.Errorf("final acceptedCount %v, want %d", values[0], voters)
	}
	fmt.Printf("verified native-EVM n=%d l=%d run=%d acceptedCount=%d\n", n, l, run, voters)
	return nil
}

func checkAcceptedEvent(parsed *abi.ABI, address common.Address, receipt *types.Receipt, count int) error {
	for _, log := range receipt.Logs {
		if log.Address == address && len(log.Topics) == 2 && log.Topics[0] == parsed.Events["TranscriptAccepted"].ID && new(big.Int).SetBytes(log.Topics[1][:]).Cmp(big.NewInt(int64(count))) == 0 {
			return nil
		}
	}
	return fmt.Errorf("receipt %s missing accepted-count event %d", receipt.TxHash, count)
}
