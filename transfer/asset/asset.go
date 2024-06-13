package asset

import (
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"math/big"
)

type Asset interface {
	Name() string
	CreateTransferTransaction(client *ethclient.Client, input *TransferInput) (*types.Transaction, error)
}

// TransferInput encapsulates the input parameters for creating a transfer transaction
type TransferInput struct {
	From     string
	To       string
	Amount   float64
	EthPrice float64
	Nonce    uint64
	GasLimit uint64
	GasPrice *big.Int
}
