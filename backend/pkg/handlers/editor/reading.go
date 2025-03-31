package editor

import (
	"backend/pkg/models"
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// @Summary Add reading material
// @Description
// @Tags globalDictOp
// @Security JWTAuth
// @Accept json
// @Produce json
// @Success 201 {object} models.ReadingMaterial
// @Failure 400 {object} models.ErrorMsg "Invalid JSON"
// @Failure 500 {object} models.ErrorMsg "Database error"
// @Router /api/reading-material/add [post]
func (h* WordHandler) AddReadingMaterial(c *gin.Context) {
    var newReadingMaterial models.ReadingMaterial

    if err := c.ShouldBindJSON(&newReadingMaterial); err != nil {
        c.AbortWithStatusJSON(400, models.ErrorMsg{Error: "Invalid JSON format"})
        return
    }

	err := h.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&newReadingMaterial).Error; err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		c.JSON(500, models.ErrorMsg{Error: "Database error"})
		return
	}

	c.JSON(201, newReadingMaterial)
}

// @Summary Edit reading material
// @Description 
// @Tags globalDictOp
// @Security JWTAuth
// @Accept json
// @Produce json
// @Failure 400 {object} models.ErrorMsg "Invalid JSON"
// @Failure 404 {object} models.ErrorMsg "Not found"
// @Failure 500 {object} models.ErrorMsg "Database error"
// @Router /api/reading-material/edit [post]
func (h* WordHandler) EditReadingMaterial(c *gin.Context) {
	var updatedReadingMaterial models.ReadingMaterial

	if err := c.ShouldBindJSON(&updatedReadingMaterial); err != nil {
		c.AbortWithStatusJSON(400, models.ErrorMsg{Error: "Invalid JSON format"})
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
			c.JSON(404, models.ErrorMsg{Error: "Reading material not found"})
		} else {
			c.JSON(500, models.ErrorMsg{Error: "Database error"})
		}
		return
	}

	c.JSON(200, updatedReadingMaterial)
}

// @Summary Delete reading material
// @Description 
// @Tags globalDictOp
// @Security JWTAuth
// @Accept json
// @Produce json
// @Failure 400 {object} models.ErrorMsg "Invalid JSON"
// @Failure 404 {object} models.ErrorMsg "Not found"
// @Failure 500 {object} models.ErrorMsg "Database error"
// @Router /api/reading-material/delete [post]
func (h* WordHandler) DeleteReadingMaterial(c *gin.Context) {
	var toDelete models.ReadingMaterial

	if err := c.ShouldBindJSON(&toDelete); err != nil {
		c.AbortWithStatusJSON(400, models.ErrorMsg{Error: "Invalid JSON format"})
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
			c.JSON(404, models.ErrorMsg{Error: "Reading material not found"})
		} else {
			c.JSON(500, models.ErrorMsg{Error: "Database error"})
		}
		return
	}

	c.JSON(200, models.SuccessMsg{Message: "Reading material deleted"})
}


// @Summary Browse all reading materials
// @Description 
// @Tags globalDictOp
// @Security JWTAuth
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param RPP query int false "Results per page"
// @Success 200 {object} []models.ReadingMaterial
// @Failure 500 {object} models.ErrorMsg "Database error"
// @Router /api/reading-material/get [get]
func (h* WordHandler) GetReadingMaterial(c *gin.Context) {
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
	
	var readingMaterials []models.ReadingMaterial
	if err := h.db.
		Limit(resultPerPage).
		Offset((pageInt - 1) * resultPerPage).
		Find(&readingMaterials).Error; err != nil {
		c.JSON(500, models.ErrorMsg{Error: "Database error"})
		return
	}

	c.JSON(200, readingMaterials)
}

// @Summary Fuzzy search in all reading materials
// @Description 
// @Tags globalDictOp
// @Security JWTAuth
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param RPP query int false "Results per page"
// @Param query query string true "Search query"
// @Success 200 {object} models.SearchResult[models.ReadingMaterial]
// @Failure 400 {object} models.ErrorMsg "Database error"
// @Router /api/reading-material/search [get]
func (h* WordHandler) FuzzySearchReadingMaterial(c *gin.Context) {
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
	var readingMaterials []models.ReadingMaterial
	var count int64
	if err := h.db.Model(&models.ReadingMaterial{}).
		Where("content LIKE ?", "%"+query+"%").
		Count(&count).Error; err != nil {
		c.JSON(500, models.ErrorMsg{Error: "Database error"})
		return
	}


	if err := h.db.
		Where("content LIKE ?", "%"+query+"%").
		Limit(resultPerPage).
		Offset((pageInt - 1) * resultPerPage).
		Find(&readingMaterials).Error; err != nil {
		c.JSON(500, models.ErrorMsg{Error: "Database error"})
		return
	}

	c.JSON(200, models.SearchResult[models.ReadingMaterial]{
		Count: count, 
		Page: pageInt, 
		PageSize: resultPerPage, 
		Results: readingMaterials,
	})
}