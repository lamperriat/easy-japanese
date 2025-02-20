package handlers

import (
	"backend/pkg/models"
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)


func (h* WordHandler) AddReadingMaterial(c *gin.Context) {
    var newReadingMaterial models.ReadingMaterial

    if err := c.ShouldBindJSON(&newReadingMaterial); err != nil {
        c.AbortWithStatusJSON(400, gin.H{"error": "Invalid JSON format"})
        return
    }

	err := h.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&newReadingMaterial).Error; err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		c.JSON(500, gin.H{"error": "Database error"})
		return
	}

	c.JSON(201, gin.H{
		"id": newReadingMaterial.ID,
	})
}


func (h* WordHandler) EditReadingMaterial(c *gin.Context) {
	var updatedReadingMaterial models.ReadingMaterial

	if err := c.ShouldBindJSON(&updatedReadingMaterial); err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": "Invalid JSON format"})
		return
	}

	err := h.db.Transaction(func(tx *gorm.DB) error {
		var existing models.ReadingMaterial
		if err := tx.First(&existing, updatedReadingMaterial.ID).Error; err != nil {
			return err
		}
		if err := tx.Model(&existing).
			Updates(&updatedReadingMaterial).Error; err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(404, gin.H{"error": "Reading material not found"})
		} else {
			c.JSON(500, gin.H{"error": "Database error"})
		}
		return
	}

	c.JSON(200, gin.H{
		"data": updatedReadingMaterial, 
		"message": "Reading material updated",
	})
}

func (h* WordHandler) DeleteReadingMaterial(c *gin.Context) {
	var toDelete models.ReadingMaterial

	if err := c.ShouldBindJSON(&toDelete); err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": "Invalid JSON format"})
		return
	}

	err := h.db.Transaction(func(tx *gorm.DB) error {
		var existing models.ReadingMaterial
		if err := tx.First(&existing, toDelete.ID).Error; err != nil {
			return err
		}
		if err := tx.Delete(&toDelete).Error; err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(404, gin.H{"error": "Reading material not found"})
		} else {
			c.JSON(500, gin.H{"error": "Database error"})
		}
		return
	}

	c.JSON(200, gin.H{
		"message": "Reading material deleted",
	})
}

const DefaultResultPerPage = 20

func (h* WordHandler) GetReadingMaterials(c *gin.Context) {
	page := c.Query("page")
	resultPerPageStr := c.Query("RPP")
	var pageInt int
	var err error
	pageInt, err = strconv.Atoi(page)
	if err != nil || pageInt < 1 {
		pageInt = 1
	}
	var resultPerPage int
	resultPerPage, err = strconv.Atoi(resultPerPageStr)
	if err != nil || resultPerPage < 1 || resultPerPage > 100 {
		resultPerPage = DefaultResultPerPage
	}
	
	var readingMaterials []models.ReadingMaterial
	if err := h.db.
		Limit(resultPerPage).
		Offset((pageInt - 1) * resultPerPage).
		Find(&readingMaterials).Error; err != nil {
		c.JSON(500, gin.H{"error": "Database error"})
		return
	}

	c.JSON(200, gin.H{"readingMaterials": readingMaterials})
}

func (h* WordHandler) FuzzySearchReadingMaterials(c *gin.Context) {
	query := c.Query("query")
	page := c.Query("page")
	resultPerPageStr := c.Query("RPP")
	var pageInt int
	var err error
	pageInt, err = strconv.Atoi(page)
	if err != nil || pageInt < 1 {
		pageInt = 1
	}
	var resultPerPage int
	resultPerPage, err = strconv.Atoi(resultPerPageStr)
	if err != nil || resultPerPage < 1 || resultPerPage > 100 {
		resultPerPage = DefaultResultPerPage
	}
	var readingMaterials []models.ReadingMaterial
	if err := h.db.
		Where("content LIKE ?", "%"+query+"%").
		Limit(resultPerPage).
		Offset((pageInt - 1) * resultPerPage).
		Find(&readingMaterials).Error; err != nil {
		c.JSON(500, gin.H{"error": "Database error"})
		return
	}

	c.JSON(200, gin.H{"readingMaterials": readingMaterials})
}