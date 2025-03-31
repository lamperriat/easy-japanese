package auth

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	JwtSecret = []byte(getEnvOrDefault("JWT_SECRET", "test-use-secret-key"))
	TokenExpiration = 12 * time.Hour
)

type UserClaim struct {
	KeyHash string `json:"key_hash"`
	jwt.RegisteredClaims
}


func getEnvOrDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		fmt.Printf("Environment variable %s not set, using default value: %s\n", key, defaultValue)
		fmt.Printf("Warning: This is not safe in production\n")
		return defaultValue
	}
	return value
}

func GenerateToken(apiKey string) (string, error) {
	keyHash := Sha256hex(apiKey)
	expiration := time.Now().Add(TokenExpiration)
	claims := UserClaim{
		keyHash,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiration),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	s, err := t.SignedString(JwtSecret)
	return s, err
}

func ParseJwt(token string) (*UserClaim, error) {
	t, err := jwt.ParseWithClaims(token, &UserClaim{}, func(t *jwt.Token) (interface{}, error) {
		return JwtSecret, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := t.Claims.(*UserClaim); ok && t.Valid {
		return claims, nil
	}
	return nil, fmt.Errorf("invalid token")
}