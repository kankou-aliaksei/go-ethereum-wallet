package asset

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

type Usdt struct {
	tokenContract common.Address
	contractABI   abi.ABI
}

func NewUsdt(contractAddress string) (*Usdt, error) {
	address := common.HexToAddress(contractAddress)
	contractABI, err := abi.JSON(strings.NewReader(`[{"constant":false,"inputs":[{"name":"_to","type":"address"},{"name":"_value","type":"uint256"}],"name":"transfer","outputs":[{"name":"","type":"bool"}],"payable":false,"stateMutability":"nonpayable","type":"function"}]`))
	if err != nil {
		return nil, fmt.Errorf("failed to parse ABI: %v", err)
	}
	return &Usdt{
		tokenContract: address,
		contractABI:   contractABI,
	}, nil
}

func (u *Usdt) Name() string {
	return "Usdt"
}

func (u *Usdt) CreateTransferTransaction(client *ethclient.Client, input *TransferInput) (*types.Transaction, error) {
	if input.From == "" || input.To == "" {
		return nil, errors.New("from and to addresses are required")
	}
	if input.Amount <= 0 {
		return nil, errors.New("amount must be greater than zero")
	}

	fromAddress := common.HexToAddress(input.From)
	toAddress := common.HexToAddress(input.To)
	tokenAddress := u.tokenContract

	amountUsdt := new(big.Int).Mul(big.NewInt(int64(input.Amount*1000000)), big.NewInt(1))

	data, err := u.contractABI.Pack("transfer", toAddress, amountUsdt)
	if err != nil {
		return nil, fmt.Errorf("failed to pack transfer data: %v", err)
	}

	msg := ethereum.CallMsg{
		From: fromAddress,
		To:   &tokenAddress,
		Data: data,
	}
	gasLimit, err := client.EstimateGas(context.Background(), msg)
	if err != nil {
		return nil, fmt.Errorf("failed to estimate gas limit: %v", err)
	}

	tx := types.NewTx(&types.LegacyTx{
		Nonce:    input.Nonce,
		To:       &tokenAddress,
		Value:    big.NewInt(0),
		Gas:      gasLimit,
		GasPrice: input.GasPrice,
		Data:     data,
	})

	return tx, nil
}
