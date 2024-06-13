// src/transaction

package transfer

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"transfer/logger"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

func sendTransaction(client *ethclient.Client, tx *types.Transaction, privateKey *ecdsa.PrivateKey) error {
	// Sign and send transaction
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(tx.ChainId()), privateKey)
	if err != nil {
		return fmt.Errorf("failed to sign transaction: %v", err)
	}

	if err := client.SendTransaction(context.Background(), signedTx); err != nil {
		return fmt.Errorf("failed to send transaction: %v", err)
	}

	txHash := signedTx.Hash().Hex()
	logger.Info.Printf("Transaction sent: %s\n", txHash)
	logger.Info.Printf("Check the transaction at: https://etherscan.io/tx/%s\n", txHash)

	return nil
}
