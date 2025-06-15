package auth

import (
	"backend/pkg/logger"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

func APIKeyValidate(apiKey string) bool {
	println("APIKeyValidate: This function is deprecated, use ValidateApiKey instead")
	if apiKey == "" {
		return false
	}
	validKeys := strings.Split(os.Getenv("API_KEYS"), ",")
	// fmt.Printf("Valid API Keys: %v\n", validKeys)
	for _, key := range validKeys {
		if key == apiKey {
			return true
		}
	}
	return false
}

func APIKeyAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		providedKey := c.GetHeader("X-API-Key")
		if providedKey == "" {
			c.AbortWithStatusJSON(401, gin.H{"error": "API Key Required"})
			return
		}
		valid := APIKeyValidate(providedKey)

		if !valid {
			c.AbortWithStatusJSON(401, gin.H{"error": "Invalid API Key"})
			return
		}
		c.Next()
	}
}

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		claims, err := ParseJwt(token)
		if err != nil {
			logger.Errorf("JWTAuth: Error parsing token: %v", false, err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "token parse error"})
			return
		}
		c.Set("keyhash", claims.KeyHash)
		c.Next()
	}
}