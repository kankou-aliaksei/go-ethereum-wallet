// ethereum_client/ethereum_client.go

package ethereum_client

import (
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"go-ethereum-wallet/transfer/logger"
)

const (
	ethPriceURL    = "https://api.coinbase.com/v2/prices/ETH-USD/spot"
	gasPriceFactor = 3
)

func GetETHUSDPrice() (float64, error) {
	resp, err := http.Get(ethPriceURL)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	var result struct {
		Data struct {
			Amount string `json:"amount"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, err
	}

	priceFloat, ok := new(big.Float).SetString(result.Data.Amount)
	if !ok {
		return 0, fmt.Errorf("invalid price value")
	}
	price, _ := priceFloat.Float64()
	return price, nil
}

func GetAddressFromPrivateKey(privateKeyHex string) (common.Address, *ecdsa.PrivateKey, error) {
	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		return common.Address{}, nil, fmt.Errorf("failed to load private key: %w", err)
	}

	publicKey, ok := privateKey.Public().(*ecdsa.PublicKey)
	if !ok {
		return common.Address{}, nil, fmt.Errorf("failed to cast public key to ECDSA")
	}

	return crypto.PubkeyToAddress(*publicKey), privateKey, nil
}

func CalculateGasPrice(client *ethclient.Client, ethPrice float64) (*big.Int, error) {
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to suggest gas price: %w", err)
	}

	increasedGasPrice := new(big.Int).Mul(gasPrice, big.NewInt(int64(gasPriceFactor*10)))
	return new(big.Int).Div(increasedGasPrice, big.NewInt(10)), nil
}

func DisplayGasPrices(client *ethclient.Client, ethPrice float64, increasedGasPrice *big.Int) {
	suggestedGasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		logger.Error.Fatalf("Failed to suggest gas price: %v", err)
	}
	suggestedGasPriceGwei := new(big.Float).Quo(new(big.Float).SetInt(suggestedGasPrice), big.NewFloat(1e9))
	suggestedGasPriceUSD := new(big.Float).Quo(new(big.Float).Mul(new(big.Float).SetInt(suggestedGasPrice), big.NewFloat(ethPrice)), big.NewFloat(1e18))

	logger.Info.Printf("Suggested Gas Price: %s Gwei\n", suggestedGasPriceGwei.String())
	logger.Info.Printf("Suggested Gas Price: $%.6f\n", suggestedGasPriceUSD)

	increasedGasPriceGwei := new(big.Float).Quo(new(big.Float).SetInt(increasedGasPrice), big.NewFloat(1e9))
	increasedGasPriceUSD := new(big.Float).Quo(new(big.Float).Mul(new(big.Float).SetInt(increasedGasPrice), big.NewFloat(ethPrice)), big.NewFloat(1e18))

	logger.Info.Printf("Increased Gas Price: %s Gwei\n", increasedGasPriceGwei.String())
	logger.Info.Printf("Increased Gas Price: $%.6f\n", increasedGasPriceUSD)
}

func CalculateTransactionFee(increasedGasPrice *big.Int, gasLimit uint64, ethPrice float64) float64 {
	transactionFeeWei := new(big.Int).Mul(increasedGasPrice, big.NewInt(int64(gasLimit)))
	transactionFeeUSD := new(big.Float).Quo(new(big.Float).Mul(new(big.Float).SetInt(transactionFeeWei), big.NewFloat(ethPrice)), big.NewFloat(1e18))

	transactionFeeUSDValue, _ := transactionFeeUSD.Float64()
	return transactionFeeUSDValue
}
