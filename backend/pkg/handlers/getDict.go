package handlers

import "github.com/gin-gonic/gin"



func GetDict(c *gin.Context) {
	dictName := c.Param("dictName")
	words, err := loadDict(dictName)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to load dictionary"})
		return
	}
	c.JSON(200, words)
}