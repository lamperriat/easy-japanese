package handlers

import (
	"math/rand"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetRandomNumber(c *gin.Context) {
	min := 1
	max := 100
	randomNumber := rand.Intn(max-min+1) + min

	c.JSON(http.StatusOK, gin.H{
		"random": randomNumber,
	})
}