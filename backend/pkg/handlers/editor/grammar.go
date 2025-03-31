package editor

import (
	"backend/pkg/models"
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// @Summary Add grammar
// @Description
// @Tags globalDictOp
// @Security JWTAuth
// @Accept json
// @Produce json
// @Success 201 {object} models.Grammar
// @Failure 400 {object} models.ErrorMsg "Invalid JSON"
// @Failure 500 {object} models.ErrorMsg "Database error"
// @Router /api/grammar/add [post]
func (h* WordHandler) AddGrammar(c *gin.Context) {
	var newGrammar models.Grammar
	if err := c.ShouldBindJSON(&newGrammar); err != nil {
		c.AbortWithStatusJSON(400, models.ErrorMsg{Error: "Invalid JSON"})
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

	c.JSON(201, newGrammar)
}

// @Summary Edit grammar 
// @Description 
// @Tags globalDictOp
// @Security JWTAuth
// @Accept json
// @Produce json
// @Success 200 {object} models.Grammar
// @Failure 400 {object} models.ErrorMsg "Invalid JSON"
// @Failure 404 {object} models.ErrorMsg "Not found"
// @Failure 500 {object} models.ErrorMsg "Database error"
// @Router /api/grammar/edit [post]
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
	c.JSON(200, updatedGrammar)
}

// @Summary Delete grammar 
// @Description 
// @Tags globalDictOp
// @Security JWTAuth
// @Accept json
// @Produce json
// @Success 200 {object} models.SuccessMsg "Grammar deleted"
// @Failure 400 {object} models.ErrorMsg "Invalid JSON"
// @Failure 404 {object} models.ErrorMsg "Not found"
// @Failure 500 {object} models.ErrorMsg "Database error"
// @Router /api/grammar/delete [post]
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

	c.JSON(200, models.SuccessMsg{Message: "Grammar deleted"})
}

// @Summary Browse all grammars 
// @Description 
// @Tags globalDictOp
// @Security JWTAuth
// @Produce json
// @Param page query int false "Page number"
// @Param RPP query int false "Results per page"
// @Success 200 {object} []models.Grammar
// @Failure 500 {object} models.ErrorMsg "Database error"
// @Router /api/grammar/get [get]
func (h *WordHandler) GetGrammar(c *gin.Context) {
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
		resultPerPage = defaultResultPerPage
	}
	query := h.db.Preload("Examples").Model(&models.Grammar{})
	var grammars []models.Grammar
	if err := query.
		Limit(resultPerPage).
		Offset((pageInt - 1) * resultPerPage).
		Find(&grammars).Error; err != nil {
		c.AbortWithStatusJSON(500, gin.H{"error": "Database error"})
		return
	}

	c.JSON(200, grammars)
}

// @Summary Search among all grammars 
// @Description 
// @Tags globalDictOp
// @Security JWTAuth
// @Produce json
// @Param query query string true "Search query"
// @Param page query int false "Page number"
// @Param RPP query int false "Results per page"
// @Success 200 {object} models.SearchResult[models.Grammar]
// @Failure 500 {object} models.ErrorMsg "Database error"
// @Router /api/grammar/search [get]
func (h* WordHandler) FuzzySearchGrammar(c *gin.Context) {
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
		resultPerPage = defaultResultPerPage
	}
	var grammars []models.Grammar
	var count int64
	if err := h.db.Preload("Examples").
		Model(&models.Grammar{}).
		Where("description LIKE ?", "%"+query+"%").
		Count(&count).Error; err != nil {
		c.JSON(500, models.ErrorMsg{Error: "Database error"})
		return
	}

	if err := h.db.
		Where("description LIKE ?", "%"+query+"%").
		Limit(resultPerPage).
		Offset((pageInt - 1) * resultPerPage).
		Find(&grammars).Error; err != nil {
		c.JSON(500, models.ErrorMsg{Error: "Database error"})
		return
	}

	c.JSON(200, models.SearchResult[models.Grammar]{
		Count:    count,
		Page:     pageInt,
		PageSize: resultPerPage,
		Results:  grammars,
	})
}