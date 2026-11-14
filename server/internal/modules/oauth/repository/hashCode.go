package repository

import (
	"crypto/sha256"
	"encoding/hex"
)

// hashCode — в базе лежит только SHA-256 кода: утечка таблицы не дает обменять действующий код
func hashCode(code string) string {
	sum := sha256.Sum256([]byte(code))
	return hex.EncodeToString(sum[:])
}
