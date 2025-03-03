package utils

import (
	"crypto/sha256"
	"encoding/hex"
)

// HashPassword gera um hash SHA-256 da senha
func HashPassword(password string) string {
	hash := sha256.Sum256([]byte(password))
	return hex.EncodeToString(hash[:])
}
