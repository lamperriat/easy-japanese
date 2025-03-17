package editor

import (
	"backend/pkg/auth"
	"backend/pkg/models"
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// @Summary Add a reading material to user's dictionary
// @Description
// @Tags userDictOp
// @Security APIKeyAuth
// @Accept json
// @Produce json
// @Success 201 {object} models.SuccessMsg
// @Failure 400 {object} models.ErrorMsg "Invalid JSON"
// @Failure 404 {object} models.ErrorMsg "User not found"
// @Failure 500 {object} models.ErrorMsg "Database error"
// @Router /api/user/reading-material/add [post]
func (h* WordHandler) AddReadingMaterialUser(c *gin.Context) {
    providedKey := c.GetHeader("X-API-Key")
    keyhash := auth.Sha256hex(providedKey)
    var user models.User
    if err := h.db.Where("keyhash = ?", keyhash).
        First(&user).Error; err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            c.JSON(404, models.ErrorMsg{Error: "User not found"})
        } else {
            c.JSON(500, models.ErrorMsg{Error: "Database error"})
        } 
        return
    }

	var newReadingMaterial models.UserReadingMaterial
	if err := c.ShouldBindJSON(&newReadingMaterial); err != nil {
		c.AbortWithStatusJSON(400, models.ErrorMsg{Error: "Invalid JSON"})
		return
	}
	newReadingMaterial.UserID = user.ID
	err := h.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&newReadingMaterial).Error; err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		c.JSON(500, gin.H{"error": "Database operation failed"})
		return
	}
	c.JSON(201, models.SuccessMsg{Message: "Reading material added"})
}

// @Summary Edit a reading material in user's dictionary
// @Description
// @Tags userDictOp
// @Security APIKeyAuth
// @Accept json
// @Produce json
// @Success 200 {object} models.SuccessMsg
// @Failure 400 {object} models.ErrorMsg "Invalid JSON"
// @Failure 404 {object} models.ErrorMsg "User not found"
// @Failure 500 {object} models.ErrorMsg "Database error"
// @Router /api/user/reading-material/edit [post]
func (h* WordHandler) EditReadingMaterialUser(c *gin.Context) {
	providedKey := c.GetHeader("X-API-Key")
	keyhash := auth.Sha256hex(providedKey)
	var user models.User
	if err := h.db.Where("keyhash = ?", keyhash).
		First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(404, models.ErrorMsg{Error: "User not found"})
		} else {
			c.JSON(500, models.ErrorMsg{Error: "Database error"})
		}
		return
	}
	var updatedReadingMaterial models.UserReadingMaterial
	if err := c.ShouldBindJSON(&updatedReadingMaterial); err != nil {
		c.AbortWithStatusJSON(400, models.ErrorMsg{Error: "Invalid JSON"})
		return
	}
	updatedReadingMaterial.UserID = user.ID
	err := h.db.Transaction(func(tx *gorm.DB) error {
		var existing models.UserReadingMaterial
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

	c.JSON(200, models.SuccessMsg{Message: "Reading material updated"})
}

// @Summary Delete a reading material from user's dictionary
// @Description
// @Tags userDictOp
// @Security APIKeyAuth
// @Accept json
// @Produce json
// @Success 200 {object} models.SuccessMsg
// @Failure 400 {object} models.ErrorMsg "Invalid JSON"
// @Failure 404 {object} models.ErrorMsg "User not found"
// @Failure 500 {object} models.ErrorMsg "Database error"
// @Router /api/user/reading-material/delete [post]
func (h* WordHandler) DeleteReadingMaterialUser(c *gin.Context) {
	providedKey := c.GetHeader("X-API-Key")
	keyhash := auth.Sha256hex(providedKey)
	var user models.User
	if err := h.db.Where("keyhash = ?", keyhash).
		First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(404, models.ErrorMsg{Error: "User not found"})
		} else {
			c.JSON(500, models.ErrorMsg{Error: "Database error"})
		}
		return
	}
	var toDelete models.UserReadingMaterial
	if err := c.ShouldBindJSON(&toDelete); err != nil {
		c.AbortWithStatusJSON(400, models.ErrorMsg{Error: "Invalid JSON"})
		return
	}
	toDelete.UserID = user.ID
	err := h.db.Transaction(func(tx *gorm.DB) error {
		var existing models.UserReadingMaterial
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

// @Summary Browse all reading materials from user's dictionary
// @Description
// @Tags userDictOp
// @Security APIKeyAuth
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param RPP query int false "Results per page"
// @Success 200 {object} []models.UserReadingMaterial
// @Failure 404 {object} models.ErrorMsg "User not found"
// @Failure 500 {object} models.ErrorMsg "Database error"
// @Router /api/user/reading-material/get [get]
func (h* WordHandler) GetReadingMaterialUser(c *gin.Context) {
	providedKey := c.GetHeader("X-API-Key")
	keyhash := auth.Sha256hex(providedKey)
	var user models.User
	if err := h.db.Where("keyhash = ?", keyhash).
		First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(404, models.ErrorMsg{Error: "User not found"})
		} else {
			c.JSON(500, models.ErrorMsg{Error: "Database error"})
		}
		return
	}
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
	if err != nil || resultPerPage < 1 {
		resultPerPage = defaultResultPerPage
	}

	var readingMaterials []models.UserReadingMaterial
	if err := h.db.
		Where("user_id = ?", user.ID).
		Limit(resultPerPage).
		Offset((pageInt - 1) * resultPerPage).
		Find(&readingMaterials).Error; err != nil {
		c.JSON(500, models.ErrorMsg{Error: "Database error"})
		return
	}
	c.JSON(200, readingMaterials)
}

// @Summary Fuzzy search in all reading materials from user's dictionary
// @Description
// @Tags userDictOp
// @Security APIKeyAuth
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param RPP query int false "Results per page"
// @Param query query string true "Search query"
// @Success 200 {object} models.SearchResult[models.UserReadingMaterial]
// @Failure 404 {object} models.ErrorMsg "User not found"
// @Failure 500 {object} models.ErrorMsg "Database error"
// @Router /api/user/reading-material/search [get]
func (h* WordHandler) FuzzySearchReadingMaterialUser(c *gin.Context) {
	providedKey := c.GetHeader("X-API-Key")
	keyhash := auth.Sha256hex(providedKey)
	var user models.User
	if err := h.db.Where("keyhash = ?", keyhash).
		First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(404, models.ErrorMsg{Error: "User not found"})
		} else {
			c.JSON(500, models.ErrorMsg{Error: "Database error"})
		}
		return
	}
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
	if err != nil || resultPerPage < 1 {
		resultPerPage = defaultResultPerPage
	}
	query := c.Query("query")
	var readingMaterials []models.UserReadingMaterial
	var count int64
	if err := h.db.
		Model(&models.UserReadingMaterial{}).
		Where("user_id = ? AND (title LIKE ? OR content LIKE ?)",
			user.ID, "%"+query+"%", "%"+query+"%").
		Count(&count).Error; err != nil {
		c.JSON(500, models.ErrorMsg{Error: "Database error"})
		return
	}
	if err := h.db.
		Where("user_id = ? AND (title LIKE ? OR content LIKE ?)",
			user.ID, "%"+query+"%", "%"+query+"%").
		Limit(resultPerPage).
		Offset((pageInt - 1) * resultPerPage).
		Find(&readingMaterials).Error; err != nil {
		c.JSON(500, models.ErrorMsg{Error: "Database error"})
		return
	}
	c.JSON(200, models.SearchResult[models.UserReadingMaterial]{
		Count:    count, 
		Page:     pageInt,
		PageSize: resultPerPage,
		Results:  readingMaterials,
	})
}