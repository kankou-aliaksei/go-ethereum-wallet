# Ethereum Account Management App

This application allows users to manage Ethereum accounts by creating accounts, retrieving public addresses, and accessing private keys securely. The private keys are encrypted with a user-provided password and saved as `<account>.enc` files.

## Features

- **Create Account**: Generate a new Ethereum account with a private key and public address, and save the encrypted private key.
- **Get Address for Account**: Retrieve the public address for an existing account by decrypting the private key.
- **Get Private Key for Account**: Access the private key for an existing account by decrypting the private key.
- **Save Account with Existing Private Key**: Save an account using an existing private key and password encryption.

## Prerequisites

- Go 1.22.2 or higher
- Go Ethereum package
    - `github.com/ethereum/go-ethereum`
    - `golang.org/x/crypto/ssh/terminal`
    - `golang.org/x/crypto/scrypt`

## Installation

1. Install the necessary Go packages:

    ```bash
    go get -u github.com/ethereum/go-ethereum
    go get -u golang.org/x/crypto/ssh/terminal
    go get -u golang.org/x/crypto/scrypt
    ```

## Usage

1. Run the application:

    ```bash
    go run main.go
    ```

2. Follow the on-screen menu to create an account, get the address for an account, get the private key for an account, or save an account with an existing private key.

### Create Account

- Select option `1` from the menu.
- Enter a name for the account.
- Enter a password to encrypt the private key.
- The encrypted private key is saved as `<account>.enc`.

### Get Address for Account

- Select option `2` from the menu.
- Enter the account name.
- Enter the password used to encrypt the private key.
- The public address is displayed.

### Get Private Key for Account

- Select option `3` from the menu.
- Enter the account name.
- Enter the password used to encrypt the private key.
- The private key is displayed in hexadecimal format.

### Save Account with Existing Private Key

- Select option `4` from the menu.
- Enter the account name.
- Enter the private key in hexadecimal format.
- Enter a password to encrypt the private key.
- The encrypted private key is saved as `<account>.enc`.

### Exit

- Select option `5` to exit the application.

## Security

- Private keys are encrypted using AES encryption with a password derived using scrypt.
- Encrypted private keys are saved with the `.enc` extension and ignored by Git to prevent accidental commits.
