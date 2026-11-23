package mobilecore

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
)

type SignRequest struct {
	ChainID  string `json:"chainId"`
	Nonce    uint64 `json:"nonce"`
	Gas      uint64 `json:"gas"`
	GasPrice string `json:"gasPrice"`
	To       string `json:"to"`
	Data     string `json:"data"`
}

func NewWallet() (map[string]string, error) {
	key, err := crypto.GenerateKey()
	if err != nil {
		return nil, err
	}
	return map[string]string{"privateKey": hex.EncodeToString(crypto.FromECDSA(key)), "address": crypto.PubkeyToAddress(key.PublicKey).Hex()}, nil
}

func WalletAddress(secret string) (map[string]string, error) {
	key, err := crypto.HexToECDSA(strings.TrimPrefix(secret, "0x"))
	if err != nil {
		return nil, fmt.Errorf("invalid wallet key")
	}
	return map[string]string{"address": crypto.PubkeyToAddress(key.PublicKey).Hex()}, nil
}

func Sign(secret string, r SignRequest) (map[string]string, error) {
	key, err := crypto.HexToECDSA(strings.TrimPrefix(secret, "0x"))
	if err != nil {
		return nil, fmt.Errorf("invalid wallet key")
	}
	id, ok := new(big.Int).SetString(r.ChainID, 10)
	if !ok || id.Sign() <= 0 {
		return nil, fmt.Errorf("invalid chain ID")
	}
	price, ok := new(big.Int).SetString(r.GasPrice, 10)
	if !ok || price.Sign() <= 0 {
		return nil, fmt.Errorf("invalid gas price")
	}
	if !common.IsHexAddress(r.To) || common.HexToAddress(r.To) == (common.Address{}) || r.Gas < 21000 || r.Gas > 180000000 {
		return nil, fmt.Errorf("invalid transaction destination or gas")
	}
	data, err := hex.DecodeString(strings.TrimPrefix(r.Data, "0x"))
	if err != nil || len(data) < 4 || len(data) > 512000 {
		return nil, fmt.Errorf("invalid transaction data")
	}
	to := common.HexToAddress(r.To)
	unsigned := types.NewTx(&types.LegacyTx{Nonce: r.Nonce, To: &to, Gas: r.Gas, GasPrice: price, Data: data})
	signed, err := types.SignTx(unsigned, types.NewEIP155Signer(id), key)
	if err != nil {
		return nil, err
	}
	raw, err := signed.MarshalBinary()
	if err != nil {
		return nil, err
	}
	return map[string]string{"rawTransaction": "0x" + hex.EncodeToString(raw), "transactionHash": signed.Hash().Hex(), "address": crypto.PubkeyToAddress(key.PublicKey).Hex()}, nil
}
