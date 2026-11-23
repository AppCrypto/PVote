// Re-execute the identical signed Ganache fixture transactions on native Geth.
package main

import (
	"compress/gzip"
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"math/big"
	"os"
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus/ethash"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/params"
)

type record struct {
	Phase         string             `json:"phase"`
	N             int                `json:"n"`
	L             int                `json:"l"`
	AcceptedCount int                `json:"accepted_count"`
	Transaction   *types.Transaction `json:"transaction"`
	Receipt       *types.Receipt     `json:"receipt"`
}

func main() {
	in := flag.String("archive", "", "Ganache transaction/receipt archive from chain-bench")
	out := flag.String("out", "experiments/gas-breakdown/cross_client.csv", "cross-client gas comparison CSV")
	flag.Parse()
	if err := run(*in, *out); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(input, output string) error {
	file, err := os.Open(input)
	if err != nil {
		return err
	}
	defer file.Close()
	z, err := gzip.NewReader(file)
	if err != nil {
		return err
	}
	defer z.Close()
	decoder := json.NewDecoder(z)
	var records []record
	for {
		var r record
		if err := decoder.Decode(&r); err == io.EOF {
			break
		} else if err != nil {
			return err
		}
		records = append(records, r)
	}
	if len(records) < 2 {
		return fmt.Errorf("need deployment and submission receipts")
	}
	config := *params.AllEthashProtocolChanges
	config.ChainID = big.NewInt(1337)
	alloc := make(types.GenesisAlloc)
	for _, r := range records {
		if r.Transaction == nil || r.Receipt == nil {
			return fmt.Errorf("incomplete archive")
		}
		from, err := types.Sender(types.LatestSigner(&config), r.Transaction)
		if err != nil {
			return err
		}
		alloc[from] = types.Account{Balance: new(big.Int).Mul(big.NewInt(1000), big.NewInt(params.Ether))}
	}
	genesis := &core.Genesis{Config: &config, Alloc: alloc, GasLimit: 90_071_992_547, BaseFee: big.NewInt(params.InitialBaseFee)}
	engine := ethash.NewFaker()
	defer engine.Close()
	db, blocks, _ := core.GenerateChainWithGenesis(genesis, engine, 1, nil)
	defer db.Close()
	parent := blocks[0]
	result, err := os.Create(output)
	if err != nil {
		return err
	}
	defer result.Close()
	w := csv.NewWriter(result)
	defer w.Flush()
	if err := w.Write([]string{"phase", "n", "l", "accepted_count", "transaction_hash", "ganache_gas", "geth_gas", "difference", "receipt_status"}); err != nil {
		return err
	}
	for _, r := range records {
		next, receipts := core.GenerateChain(&config, parent, engine, db, 1, func(_ int, b *core.BlockGen) { b.AddTx(r.Transaction) })
		parent = next[0]
		receipt := receipts[0][0]
		delta := int64(receipt.GasUsed) - int64(r.Receipt.GasUsed)
		if err := w.Write([]string{r.Phase, strconv.Itoa(r.N), strconv.Itoa(r.L), strconv.Itoa(r.AcceptedCount), r.Transaction.Hash().Hex(), strconv.FormatUint(r.Receipt.GasUsed, 10), strconv.FormatUint(receipt.GasUsed, 10), strconv.FormatInt(delta, 10), strconv.FormatUint(receipt.Status, 10)}); err != nil {
			return err
		}
		fmt.Printf("cross-client %s m=%d Ganache=%d Geth=%d delta=%d\n", r.Phase, r.AcceptedCount, r.Receipt.GasUsed, receipt.GasUsed, delta)
		if receipt.Status != r.Receipt.Status || delta != 0 {
			return fmt.Errorf("cross-client gas/status mismatch for %s", r.Transaction.Hash())
		}
		if r.Phase == "deployment" && receipt.ContractAddress == (common.Address{}) {
			return fmt.Errorf("missing deployment address")
		}
	}
	w.Flush()
	return w.Error()
}
