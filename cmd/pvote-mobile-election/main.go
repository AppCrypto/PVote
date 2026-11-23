// Deploy a disposable local-chain election for a phone-owned wallet. Only public
// parameters and the public wallet address are written to the election manifest.
package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"time"

	encoding "PVote/chain/rbpvss"
	contract "PVote/compile/contract"
	rb "PVote/crypto/RBPVSS"
	mobile "PVote/mobile/core"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
func run() error {
	rpc := flag.String("rpc", "http://127.0.0.1:18545", "Mac's development RPC")
	publicRPC := flag.String("phone-rpc", "", "same RPC accessible from the phone, e.g. http://192.168.1.2:18545")
	voter := flag.String("voter", "", "address copied from iPhone app")
	out := flag.String("out", "mobile/election.json", "public manifest output")
	n := flag.Int("n", 4, "talliers")
	l := flag.Int("l", 2, "scores")
	threshold := flag.Int("t", 3, "threshold")
	flag.Parse()
	if !common.IsHexAddress(*voter) || common.HexToAddress(*voter) == (common.Address{}) || *publicRPC == "" {
		return fmt.Errorf("provide -voter and -phone-rpc")
	}
	key, err := crypto.HexToECDSA(os.Getenv("PVOTE_MOBILE_DEPLOYER_KEY"))
	if err != nil {
		return fmt.Errorf("set PVOTE_MOBILE_DEPLOYER_KEY to a funded disposable local-chain key")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	client, err := ethclient.DialContext(ctx, *rpc)
	if err != nil {
		return err
	}
	defer client.Close()
	id, err := client.ChainID(ctx)
	if err != nil {
		return err
	}
	if id.Cmp(big.NewInt(1337)) != 0 && id.Cmp(big.NewInt(31337)) != 0 {
		return fmt.Errorf("this fixture deployer only supports disposable local chains 1337/31337")
	}
	pp, _, err := rb.Setup(128, *n, *threshold, *l, rb.Interval{Min: 0, Max: 10})
	if err != nil {
		return err
	}
	from := crypto.PubkeyToAddress(key.PublicKey)
	phone := common.HexToAddress(*voter)
	nonce, err := client.PendingNonceAt(ctx, from)
	if err != nil {
		return err
	}
	price, err := client.SuggestGasPrice(ctx)
	if err != nil {
		return err
	}
	funding, err := types.SignTx(types.NewTransaction(nonce, phone, new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil), 21000, price, nil), types.NewEIP155Signer(id), key)
	if err != nil {
		return err
	}
	if err = client.SendTransaction(ctx, funding); err != nil {
		return err
	}
	receipt, err := bind.WaitMined(ctx, client, funding)
	if err != nil {
		return err
	}
	if receipt.Status != 1 {
		return fmt.Errorf("funding reverted")
	}
	auth, err := bind.NewKeyedTransactorWithChainID(key, id)
	if err != nil {
		return err
	}
	auth.Context = ctx
	auth.GasLimit = 100000000
	auth.GasPrice = price
	g0, h0, g1, pki, sigma, pks, t, coordinates, a, b, err := encoding.ConstructorArguments(pp)
	if err != nil {
		return err
	}
	address, tx, _, err := contract.DeployContract(auth, client, g0, h0, g1, pki, sigma, pks, t, coordinates, a, b, []common.Address{phone})
	if err != nil {
		return err
	}
	receipt, err = bind.WaitMined(ctx, client, tx)
	if err != nil {
		return err
	}
	if receipt.Status != 1 {
		return fmt.Errorf("deployment reverted")
	}
	labels := make([]string, *l)
	for i := range labels {
		labels[i] = fmt.Sprintf("Option %d", i+1)
	}
	election := mobile.Election{SchemaVersion: 1, ElectionID: fmt.Sprintf("%s/%s", id, address.Hex()), Protocol: mobile.Legacy,
		RPCURL: *publicRPC, Contract: address.Hex(), ChainID: id.String(), VoterAddress: phone.Hex(), VoterIndex: 1, Labels: labels, Parameters: mobile.Export(pp)}
	data, err := json.MarshalIndent(election, "", "  ")
	if err != nil {
		return err
	}
	if err = os.MkdirAll(filepath.Dir(*out), 0755); err != nil {
		return err
	}
	if err = os.WriteFile(*out, append(data, '\n'), 0644); err != nil {
		return err
	}
	fmt.Printf("Registered phone wallet %s and deployed %s\nImport %s on iPhone.\n", phone.Hex(), address.Hex(), *out)
	return nil
}
