package userinput

import (
	"fmt"
	"go-ethereum-wallet/transfer/asset"
)

func SelectAsset(assets map[string]asset.Asset) string {
	var assetChoice string
	fmt.Println("Select the asset to transfer:")
	for key, value := range assets {
		fmt.Printf("%s: %s\n", key, value.Name())
	}
	fmt.Println("Enter the number of your choice: ")
	fmt.Scanln(&assetChoice)
	return assetChoice
}

func GetAccountName() string {
	var accountName string
	fmt.Print("Enter your account name: ")
	fmt.Scanln(&accountName)
	return accountName
}

func GetReceiverAddress() string {
	var receiverAddress string
	fmt.Print("Enter the receiver's address: ")
	fmt.Scanln(&receiverAddress)
	return receiverAddress
}

func GetTransferAmount() float64 {
	var amountInDollars float64
	fmt.Print("Enter the amount to transfer (in USD): ")
	fmt.Scanf("%f", &amountInDollars)
	return amountInDollars
}

func ConfirmTransaction() bool {
	var confirmation string
	fmt.Print("Are you okay with this increased gas price and transaction fee? (yes/no): ")
	fmt.Scanln(&confirmation)
	return confirmation == "yes"
}
