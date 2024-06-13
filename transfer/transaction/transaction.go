// transaction/transaction.go

package transaction

import (
	"context"
	"crypto/ecdsa"
	"fmt"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"go-ethereum-wallet/transfer/logger"
)

func SendTransaction(client *ethclient.Client, tx *types.Transaction, privateKey *ecdsa.PrivateKey, baseURL string) error {
	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		return fmt.Errorf("failed to get network ID: %w", err)
	}

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		return fmt.Errorf("failed to sign transaction: %w", err)
	}

	if err := client.SendTransaction(context.Background(), signedTx); err != nil {
		return fmt.Errorf("failed to send transaction: %w", err)
	}

	txHash := signedTx.Hash().Hex()
	logger.Info.Printf("Transaction sent: %s\n", txHash)
	logger.Info.Printf("Check the transaction at: %s/tx/%s\n", baseURL, txHash)

	return nil
}
