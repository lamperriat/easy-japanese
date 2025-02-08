package auth

import (
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

func APIKeyAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		providedKey := c.GetHeader("X-API-Key")
		if providedKey == "" {
			c.AbortWithStatusJSON(401, gin.H{"error": "API Key Required"})
			return
		}
		validKeys := strings.Split(os.Getenv("API_KEYS"), ",")

		valid := false
		for _, key := range validKeys {
			if key == providedKey {
				valid = true
				break
			}
		}

		if !valid {
			c.AbortWithStatusJSON(401, gin.H{"error": "Invalid API Key"})
			return
		}
		c.Next()
	}
}