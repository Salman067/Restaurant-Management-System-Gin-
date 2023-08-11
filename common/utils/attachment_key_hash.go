package utils

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
)

func GenerateKeyHash() (string, error) {
	// Generate random bytes
	randomBytes := make([]byte, 32)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", err
	}

	// Compute hash of the random bytes
	hash := sha256.Sum256(randomBytes)

	// Convert the hash to a hexadecimal string
	hashString := hex.EncodeToString(hash[:])

	return hashString, nil
}
