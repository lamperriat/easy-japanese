package handlers

import (
	"backend/pkg/models"
	"errors"

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