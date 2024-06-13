package main

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/kankou-aliaksei/go-ethereum-wallet/keygen"
	"github.com/kankou-aliaksei/go-ethereum-wallet/transfer/asset"
	"github.com/kankou-aliaksei/go-ethereum-wallet/transfer/ethereum_client"
	"github.com/kankou-aliaksei/go-ethereum-wallet/transfer/logger"
	"github.com/kankou-aliaksei/go-ethereum-wallet/transfer/transaction"
)

const (
	publicNodeUrl       = "https://rpc.sepolia.org"                    // Mainnet: "https://cloudflare-eth.com", Sepolia: "https://rpc.sepolia.org"
	ethereumExplorerUrl = "https://sepolia.etherscan.io"               // Mainnet: "https://etherscan.io", Sepolia: "https://sepolia.etherscan.io"
	usdtContractAddress = "0xdAC17F958D2ee523a2206206994597C13D831ec7" // Mainnet: "0xdAC17F958D2ee523a2206206994597C13D831ec7"
)

func main() {
	var accountName, receiverAddress, assetChoice string
	var amountInDollars float64

	usdtAsset, err := asset.NewUsdt(usdtContractAddress)
	if err != nil {
		logger.Error.Fatalf("Failed to create Usdt asset: %v", err)
	}

	var assets = map[string]asset.Asset{
		"1": &asset.Ether{},
		"2": usdtAsset,
	}

	fmt.Println("Select the asset to transfer:")
	for key, value := range assets {
		fmt.Printf("%s: %s\n", key, value.Name())
	}
	fmt.Println("Enter the number of your choice: ")
	fmt.Scanln(&assetChoice)

	currentAsset, exists := assets[assetChoice]
	if !exists {
		logger.Error.Fatalf("Invalid asset choice")
	}

	fmt.Print("Enter your account name: ")
	fmt.Scanln(&accountName)

	privateKeyHex, err := keygen.RetrievePrivateKeyHex(accountName)
	if err != nil {
		logger.Error.Fatalf("Failed to retrieve private key: %v", err)
	}

	fmt.Print("Enter the receiver's address: ")
	fmt.Scanln(&receiverAddress)

	ethPrice, err := ethereum_client.GetETHUSDPrice()
	if err != nil {
		logger.Error.Fatalf("Failed to get ETH price: %v", err)
	}
	logger.Info.Printf("Current ETH/USD price: $%.2f\n", ethPrice)

	fmt.Print("Enter the amount to transfer (in USD): ")
	fmt.Scanf("%f", &amountInDollars)

	if amountInDollars <= 0 {
		logger.Error.Fatalf("Invalid amount")
	}

	client, err := ethclient.Dial(publicNodeUrl)
	if err != nil {
		logger.Error.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}

	fromAddress, privateKey, err := ethereum_client.GetAddressFromPrivateKey(privateKeyHex)
	if err != nil {
		logger.Error.Fatalf("Failed to get address from private key: %v", err)
	}

	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		logger.Error.Fatalf("Failed to get nonce: %v", err)
	}
	logger.Info.Printf("Nonce: %d\n", nonce)

	increasedGasPrice, err := ethereum_client.CalculateGasPrice(client, ethPrice)
	if err != nil {
		logger.Error.Fatalf("Failed to calculate gas price: %v", err)
	}

	ethereum_client.DisplayGasPrices(client, ethPrice, increasedGasPrice)

	gasLimit := uint64(21000)
	transactionFeeUSD := ethereum_client.CalculateTransactionFee(increasedGasPrice, gasLimit, ethPrice)

	logger.Info.Printf("Transaction Fee: $%.6f\n", transactionFeeUSD)

	balance, err := client.BalanceAt(context.Background(), fromAddress, nil)
	if err != nil {
		logger.Error.Fatalf("Failed to get balance: %v", err)
	}
	logger.Info.Printf("Sender's balance: %s wei\n", balance.String())

	requiredGasFee := new(big.Int).Mul(big.NewInt(int64(gasLimit)), increasedGasPrice)
	if balance.Cmp(requiredGasFee) < 0 {
		logger.Error.Fatalf("Insufficient balance to cover transaction fee: required %s wei, but only %s wei available", requiredGasFee.String(), balance.String())
	}

	var confirmation string
	fmt.Print("Are you okay with this increased gas price and transaction fee? (yes/no): ")
	fmt.Scanln(&confirmation)

	if confirmation != "yes" {
		logger.Info.Println("Transaction cancelled.")
		return
	}

	input := &asset.TransferInput{
		From:     fromAddress.Hex(),
		To:       receiverAddress,
		Amount:   amountInDollars,
		EthPrice: ethPrice,
		Nonce:    nonce,
		GasLimit: gasLimit,
		GasPrice: increasedGasPrice,
	}

	tx, err := currentAsset.CreateTransferTransaction(client, input)
	if err != nil {
		logger.Error.Fatalf("Failed to create transaction: %v", err)
	}

	if err := transaction.SendTransaction(client, tx, privateKey, ethereumExplorerUrl); err != nil {
		logger.Error.Fatalf("Transaction sending failed: %v", err)
	}
}
