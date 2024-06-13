package main

import (
	"fmt"
)

func main() {
	for {
		fmt.Println("Menu:")
		fmt.Println("1. Create account")
		fmt.Println("2. Get address for account")
		fmt.Println("3. Get private key for account")
		fmt.Println("4. Save account with existing private key")
		fmt.Println("5. Exit")
		fmt.Print("Enter your choice: ")

		var choice int
		fmt.Scanln(&choice)

		switch choice {
		case 1:
			CreateAccount()
		case 2:
			GetAddressForAccount()
		case 3:
			GetPrivateKeyForAccount()
		case 4:
			SaveAccountWithPrivateKey()
		case 5:
			return
		default:
			fmt.Println("Invalid choice. Please try again.")
		}
	}
}
