package authmid

import (
	"backend/pkg/auth"
	"backend/pkg/logger"
	"backend/pkg/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// @Summary Get authentication token
// @Description Exchange API key for a JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param X-API-KEY header string true "API Key"
// @Success 200 {object} models.TokenResponse
// @Failure 401 {object} models.ErrorMsg
// @Router /api/auth/token [post]
func GetToken(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.GetHeader("X-API-KEY")
		println(apiKey)
		valid, err := auth.ValidateApiKey(db, apiKey)
		println(valid)
		if err != nil {
			logger.Errorf("GetToken: Error validating API key: %v", err)
		}
		if valid && err == nil{
			token, err := auth.GenerateToken(apiKey)
			if err != nil {
				c.JSON(500, models.ErrorMsg{Error: "Failed to generate token"})
				return
			}
			c.JSON(200, models.TokenResponse{
				Token:     token,
				ExpiresIn: int(auth.TokenExpiration.Seconds()),
			})
		} else {
			c.JSON(401, models.ErrorMsg{Error: "Invalid API Key"})
		}
	}
}

