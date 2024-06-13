package main

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/ethclient"

	"go-ethereum-wallet/keygen"
	"go-ethereum-wallet/transfer/asset"
	"go-ethereum-wallet/transfer/config"
	"go-ethereum-wallet/transfer/ethereum_client"
	"go-ethereum-wallet/transfer/logger"
	"go-ethereum-wallet/transfer/transaction"
	"go-ethereum-wallet/transfer/userinput"
)

func main() {
	var accountName, receiverAddress, assetChoice string
	var amountInDollars float64

	// Choose the desired configuration
	cfg := config.SepoliaTestnet

	usdtAsset, err := asset.NewUsdt(cfg.UsdtContractAddress)
	if err != nil {
		logger.Error.Fatalf("Failed to create Usdt asset: %v", err)
	}

	var assets = map[string]asset.Asset{
		"1": &asset.Ether{},
		"2": usdtAsset,
	}

	assetChoice = userinput.SelectAsset(assets)
	currentAsset, exists := assets[assetChoice]
	if !exists {
		logger.Error.Fatalf("Invalid asset choice")
	}

	accountName = userinput.GetAccountName()
	privateKeyHex, err := keygen.RetrievePrivateKeyHex(accountName)
	if err != nil {
		logger.Error.Fatalf("Failed to retrieve private key: %v", err)
	}

	receiverAddress = userinput.GetReceiverAddress()
	ethPrice, err := ethereum_client.GetETHUSDPrice()
	if err != nil {
		logger.Error.Fatalf("Failed to get ETH price: %v", err)
	}
	logger.Info.Printf("Current ETH/USD price: $%.2f\n", ethPrice)

	amountInDollars = userinput.GetTransferAmount()
	if amountInDollars <= 0 {
		logger.Error.Fatalf("Invalid amount")
	}

	client, err := ethclient.Dial(cfg.PublicNodeUrl)
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

	if !userinput.ConfirmTransaction() {
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

	if err := transaction.SendTransaction(client, tx, privateKey, cfg.EthereumExplorerUrl); err != nil {
		logger.Error.Fatalf("Transaction sending failed: %v", err)
	}
}
