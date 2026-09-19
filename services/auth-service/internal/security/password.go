package security

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
)

// HashPassword generates a salted SHA-256 hash for secure password storage.
func HashPassword(password string) (string, error) {
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return "", fmt.Errorf("failed to generate salt: %w", err)
	}

	saltHex := hex.EncodeToString(salt)
	hash := sha256.Sum256([]byte(password + saltHex))
	hashHex := hex.EncodeToString(hash[:])

	return fmt.Sprintf("%s:%s", saltHex, hashHex), nil
}

// VerifyPassword checks if a plain password matches a stored salt:hash string.
func VerifyPassword(password, storedHash string) bool {
	parts := strings.Split(storedHash, ":")
	if len(parts) != 2 {
		return false
	}

	saltHex := parts[0]
	expectedHash := parts[1]

	hash := sha256.Sum256([]byte(password + saltHex))
	computedHash := hex.EncodeToString(hash[:])

	return computedHash == expectedHash
}
