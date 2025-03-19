package editor

import (
	"math/rand"
	"net/http"

	"github.com/gin-gonic/gin"
)

// @Summary Get a random number
// @Description Test use 
// @Produce json
// @Success 200 {object} map[string]int
// @Router /api/random [get]
func GetRandomNumber(c *gin.Context) {
	min := 1
	max := 100
	randomNumber := rand.Intn(max-min+1) + min

	c.JSON(http.StatusOK, gin.H{
		"random": randomNumber,
	})
}