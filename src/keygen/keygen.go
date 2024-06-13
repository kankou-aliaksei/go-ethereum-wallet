package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdsa"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"github.com/ethereum/go-ethereum/crypto"
	"golang.org/x/crypto/scrypt"
	"golang.org/x/crypto/ssh/terminal"
	"log"
	"os"
	"path/filepath"
)

const AccountPath = "../../account"

func CreateAccount() {
	fmt.Print("Enter account name: ")
	var accountName string
	fmt.Scanln(&accountName)

	// Generate a new private key
	privateKey, err := ecdsa.GenerateKey(crypto.S256(), rand.Reader)
	if err != nil {
		log.Fatalf("Failed to generate private key: %v", err)
	}

	// Convert the private key to bytes
	privateKeyBytes := crypto.FromECDSA(privateKey)
	privateKeyHex := hex.EncodeToString(privateKeyBytes)
	fmt.Printf("Private Key: %s\n", privateKeyHex)

	// Get the corresponding public address
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatalf("Failed to cast public key to ECDSA")
	}

	address := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()
	fmt.Printf("Public Address: %s\n", address)

	// Save the private key to a file
	fmt.Print("Enter a password to encrypt the private key: ")
	password, err := terminal.ReadPassword(0)
	if err != nil {
		log.Fatalf("Failed to read password: %v", err)
	}
	fmt.Println()

	encryptedKey, err := EncryptKey(privateKeyBytes, string(password))
	if err != nil {
		log.Fatalf("Failed to encrypt private key: %v", err)
	}

	filePath := filepath.Join(AccountPath, accountName+".enc")
	err = os.MkdirAll(filepath.Dir(filePath), 0755)
	if err != nil {
		log.Fatalf("Failed to create directories: %v", err)
	}

	err = os.WriteFile(filePath, encryptedKey, 0644)
	if err != nil {
		log.Fatalf("Failed to save private key: %v", err)
	}

	fmt.Printf("Private key successfully saved to '%s'\n", filePath)
}

func GetAddressForAccount() {
	fmt.Print("Enter account name: ")
	var accountName string
	fmt.Scanln(&accountName)

	filePath := filepath.Join(AccountPath, accountName+".enc")
	encryptedKey, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("Failed to read private key file: %v", err)
	}

	fmt.Print("Enter password: ")
	password, err := terminal.ReadPassword(0)
	if err != nil {
		log.Fatalf("Failed to read password: %v", err)
	}
	fmt.Println()

	privateKeyBytes, err := DecryptKey(encryptedKey, string(password))
	if err != nil {
		log.Fatalf("Failed to decrypt private key: %v", err)
	}

	privateKey, err := crypto.ToECDSA(privateKeyBytes)
	if err != nil {
		log.Fatalf("Failed to convert bytes to ECDSA: %v", err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatalf("Failed to cast public key to ECDSA")
	}

	address := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()
	fmt.Printf("Public Address: %s\n", address)
}

func GetPrivateKeyForAccount() {
	fmt.Print("Enter account name: ")
	var accountName string
	fmt.Scanln(&accountName)

	privateKeyHex, err := RetrievePrivateKeyHex(accountName)
	if err != nil {
		log.Fatalf("Failed to retrieve private key: %v", err)
	}

	fmt.Printf("Private Key: %s\n", privateKeyHex)
}

func SaveAccountWithPrivateKey() {
	fmt.Print("Enter account name: ")
	var accountName string
	fmt.Scanln(&accountName)

	fmt.Print("Enter private key (hex format): ")
	var privateKeyHex string
	fmt.Scanln(&privateKeyHex)

	privateKeyBytes, err := hex.DecodeString(privateKeyHex)
	if err != nil {
		log.Fatalf("Failed to decode private key: %v", err)
	}

	privateKey, err := crypto.ToECDSA(privateKeyBytes)
	if err != nil {
		log.Fatalf("Invalid private key: %v", err)
	}

	// Get the corresponding public address
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatalf("Failed to cast public key to ECDSA")
	}

	address := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()
	fmt.Printf("Public Address: %s\n", address)

	// Save the private key to a file
	fmt.Print("Enter a password to encrypt the private key: ")
	password, err := terminal.ReadPassword(0)
	if err != nil {
		log.Fatalf("Failed to read password: %v", err)
	}
	fmt.Println()

	encryptedKey, err := EncryptKey(privateKeyBytes, string(password))
	if err != nil {
		log.Fatalf("Failed to encrypt private key: %v", err)
	}

	filePath := filepath.Join(AccountPath, accountName+".enc")
	err = os.MkdirAll(filepath.Dir(filePath), 0755)
	if err != nil {
		log.Fatalf("Failed to create directories: %v", err)
	}

	err = os.WriteFile(filePath, encryptedKey, 0644)
	if err != nil {
		log.Fatalf("Failed to save private key: %v", err)
	}

	fmt.Printf("Private key successfully saved to '%s'\n", filePath)
}

func RetrievePrivateKeyHex(accountName string) (string, error) {
	filePath := filepath.Join(AccountPath, accountName+".enc")
	encryptedKey, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to read private key file: %v", err)
	}

	fmt.Print("Enter password: ")
	password, err := terminal.ReadPassword(0)
	if err != nil {
		return "", fmt.Errorf("failed to read password: %v", err)
	}
	fmt.Println()

	privateKeyBytes, err := DecryptKey(encryptedKey, string(password))
	if err != nil {
		return "", fmt.Errorf("failed to decrypt private key: %v", err)
	}

	return hex.EncodeToString(privateKeyBytes), nil
}

// Encrypt the key using AES encryption
func EncryptKey(key []byte, passphrase string) ([]byte, error) {
	salt := make([]byte, 16)
	_, err := rand.Read(salt)
	if err != nil {
		return nil, err
	}

	dk, err := scrypt.Key([]byte(passphrase), salt, 32768, 8, 1, 32)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(dk)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	_, err = rand.Read(nonce)
	if err != nil {
		return nil, err
	}

	ciphertext := gcm.Seal(nonce, nonce, key, nil)
	return append(salt, ciphertext...), nil
}

// Decrypt the key using AES decryption
func DecryptKey(encryptedKey []byte, passphrase string) ([]byte, error) {
	salt := encryptedKey[:16]
	ciphertext := encryptedKey[16:]

	dk, err := scrypt.Key([]byte(passphrase), salt, 32768, 8, 1, 32)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(dk)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, fmt.Errorf("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plainText, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return plainText, nil
}
