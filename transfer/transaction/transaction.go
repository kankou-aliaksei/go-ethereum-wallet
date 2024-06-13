package transaction

import (
	"context"
	"crypto/ecdsa"
	"fmt"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/kankou-aliaksei/go-ethereum-wallet/transfer/logger"
)

func SendTransaction(client *ethclient.Client, tx *types.Transaction, privateKey *ecdsa.PrivateKey) error {
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
