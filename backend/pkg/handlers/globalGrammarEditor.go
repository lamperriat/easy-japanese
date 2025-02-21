package handlers

import (
	"backend/pkg/models"
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func (h* WordHandler) AddGrammar(c *gin.Context) {
	var newGrammar models.Grammar
	if err := c.ShouldBindJSON(&newGrammar); err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": "Invalid JSON"})
		return
	}

	err := h.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&newGrammar).Error; err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		c.JSON(500, gin.H{"error": "Database operation failed"})
		return
	}

	c.JSON(201, gin.H{
		"id": newGrammar.ID,
		"description": newGrammar.Description,
	})
}

func (h* WordHandler) EditGrammar(c *gin.Context) {
	var editedGrammar models.Grammar
	if err := c.ShouldBindJSON(&editedGrammar); err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": "Invalid JSON"})
		return
	}

	err := h.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&models.Grammar{}).
			Where("id = ?", editedGrammar.ID).
			Omit("Examples").
			Updates(&editedGrammar).Error; err != nil {
			return err
		}

		if err := tx.Where("grammar_id = ?", editedGrammar.ID).
			Delete(&models.GrammarExample{}).Error; err != nil {
			return nil
		}

		if len(editedGrammar.Examples) > 0 {
			for i := range editedGrammar.Examples {
				editedGrammar.Examples[i].GrammarID = editedGrammar.ID
				editedGrammar.Examples[i].ID = 0
			}
			if err := tx.Create(&editedGrammar.Examples).Error; err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(404, gin.H{"error": "Grammar not found"})
		} else {
			c.JSON(500, gin.H{"error": "Database operation failed"})
		}
		return
	}

	var updatedGrammar models.Grammar
	if err := h.db.Preload("Examples").First(&updatedGrammar, editedGrammar.ID).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to retrieve updated grammar"})
		return
	}
	c.JSON(200, gin.H{
		"data": updatedGrammar,
		"message": "Grammar updated",
	})
}

func (h *WordHandler) DeleteGrammar(c *gin.Context) {
	var toDelete models.Grammar

	if err := c.ShouldBindJSON(&toDelete); err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": "Invalid JSON format"})
		return
	}

	err := h.db.Transaction(func(tx *gorm.DB) error {
		var existing models.Grammar
		if err := tx.First(&existing, toDelete.ID).Error; err != nil {
			return err
		}
		if err := tx.Where("grammar_id = ?", toDelete.ID).
			Delete(&models.GrammarExample{}).
			Error; err != nil {
			return err
		}
		if err := tx.Delete(&toDelete).Error; err != nil {
			return err
		}
		
		return nil
	})

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(404, gin.H{"error": "Grammar not found"})
		} else {
			c.JSON(500, gin.H{"error": "Database operation error"})
		}
		return
	}
}