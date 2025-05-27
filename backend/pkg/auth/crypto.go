package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"

	"golang.org/x/crypto/bcrypt"
)

func Sha256hex(t string) string {
	hasher := sha256.New()
	hasher.Write([]byte(t))
	return hex.EncodeToString(hasher.Sum(nil))
}

func SafeHash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// Returns a random string of n_bits length, base64 encoded
// If n_bits is less than 128, it will be set to 128
func SafeRandom(n_bits int) (string, error) {
	if n_bits < 128 {
		n_bits = 128
	}
	bytes := make([]byte, n_bits/8)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(bytes), nil
}

func SafeCompare(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}