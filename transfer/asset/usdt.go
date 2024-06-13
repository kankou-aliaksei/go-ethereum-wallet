package asset

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

// Contract address for USDT token
const UsdtContractAddress = "0xdAC17F958D2ee523a2206206994597C13D831ec7"

// Usdt represents the USDT asset
type Usdt struct {
	tokenContract common.Address
	contractABI   abi.ABI
}

func NewUsdt() *Usdt {
	contractAddress := common.HexToAddress(UsdtContractAddress)
	contractABI, _ := abi.JSON(strings.NewReader(`[{"constant":false,"inputs":[{"name":"_to","type":"address"},{"name":"_value","type":"uint256"}],"name":"transfer","outputs":[{"name":"","type":"bool"}],"payable":false,"stateMutability":"nonpayable","type":"function"}]`))
	return &Usdt{
		tokenContract: contractAddress,
		contractABI:   contractABI,
	}
}

func (u *Usdt) Name() string {
	return "USDT"
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

	// Convert USD amount to USDT (Tether)
	// 1 USDT = 10^6 Tether units (Wei)
	amountUSDT := new(big.Int).Mul(big.NewInt(int64(input.Amount*1000000)), big.NewInt(1))

	// Pack the transfer data
	data, err := u.contractABI.Pack("transfer", toAddress, amountUSDT)
	if err != nil {
		return nil, fmt.Errorf("failed to pack transfer data: %v", err)
	}

	// Estimate gas limit for token transfer
	msg := ethereum.CallMsg{
		From: fromAddress,
		To:   &tokenAddress,
		Data: data,
	}
	gasLimit, err := client.EstimateGas(context.Background(), msg)
	if err != nil {
		return nil, fmt.Errorf("failed to estimate gas limit: %v", err)
	}

	// Create the transaction
	tx := types.NewTx(&types.LegacyTx{
		Nonce:    input.Nonce,
		To:       &tokenAddress,
		Value:    big.NewInt(0), // Set Value to 0 since we are not transferring Ether
		Gas:      gasLimit,
		GasPrice: input.GasPrice,
		Data:     data,
	})

	return tx, nil
}
