package authmid

import (
	"backend/pkg/auth"
	"backend/pkg/models"

	"github.com/gin-gonic/gin"
)


// @Summary Get authentication token
// @Description Exchange API key for a JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param X-API-Key header string true "API Key"
// @Success 200 {object} models.TokenResponse
// @Failure 401 {object} models.ErrorMsg
// @Router /api/auth/token [post]
func GetToken(c *gin.Context) {
	apiKey := c.GetHeader("X-API-Key")
	valid := auth.APIKeyValidate(apiKey)
	if valid {
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