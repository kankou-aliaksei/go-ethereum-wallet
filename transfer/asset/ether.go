package asset

import (
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

// Ether represents the Ether asset
type Ether struct{}

func (e *Ether) Name() string {
	return "Ether"
}

func (e *Ether) CreateTransferTransaction(client *ethclient.Client, input *TransferInput) (*types.Transaction, error) {
	if input.From == "" || input.To == "" {
		return nil, errors.New("from and to addresses are required")
	}
	if input.Amount <= 0 {
		return nil, errors.New("amount must be greater than zero")
	}

	toAddress := common.HexToAddress(input.To)

	// Convert amount from USD to Wei
	amountInWei := new(big.Float).Mul(big.NewFloat(input.Amount), big.NewFloat(1e18))
	amountInWei.Quo(amountInWei, big.NewFloat(input.EthPrice))
	amountBigInt := new(big.Int)
	amountInWei.Int(amountBigInt)

	tx := types.NewTx(&types.LegacyTx{
		Nonce:    input.Nonce,
		To:       &toAddress,
		Value:    amountBigInt,
		Gas:      input.GasLimit,
		GasPrice: input.GasPrice,
		Data:     nil,
	})

	return tx, nil
}
