package handlers

import (
	"backend/pkg/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

var resultPerPage = 30

func (h *WordHandler) GetDict(c *gin.Context) {
	dictName := c.Param("dictName")
	page := c.Query("page")
	var pageInt int
	var err error
	pageInt, err = strconv.Atoi(page)
	if err != nil || pageInt < 1 {
		pageInt = 1
	}
	query := h.db.Model(&models.JapaneseWord{})
	if dictName != "all" {
		query = query.Where("dict_name = ?", dictName)
	}

	var words []models.JapaneseWord
	if err := query.
		Limit(resultPerPage).
		Offset((pageInt - 1) * resultPerPage).
		Find(&words).Error; 
		err != nil {
		c.AbortWithStatusJSON(500, gin.H{"error": "Database error"})
		return
	}

	c.JSON(200, gin.H{"words": words})


}