package auth

import (
	"crypto/sha256"
	"encoding/hex"
)

func Sha256hex(t string) string {
	hasher := sha256.New()
	hasher.Write([]byte(t))
	return hex.EncodeToString(hasher.Sum(nil))
}