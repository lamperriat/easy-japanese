package editor

import (
	"backend/pkg/auth"
	"backend/pkg/models"
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// @Summary Add a grammar to user's dictionary
// @Description
// @Tags userDictOp
// @Security APIKeyAuth
// @Accept json
// @Produce json
// @Success 201 {object} models.SuccessMsg
// @Failure 400 {object} models.ErrorMsg "Invalid JSON"
// @Failure 404 {object} models.ErrorMsg "User not found"
// @Failure 500 {object} models.ErrorMsg "Database error"
// @Router /api/user/grammar/add [post]
func (h* WordHandler) AddGrammarUser(c *gin.Context) {
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

	var newGrammar models.UserGrammar
	if err := c.ShouldBindJSON(&newGrammar); err != nil {
		c.AbortWithStatusJSON(400, models.ErrorMsg{Error: "Invalid JSON"})
		return
	}

	newGrammar.UserID = user.ID

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

	c.JSON(201, models.SuccessMsg{Message: "Grammar added"})
}

// @Summary Edit a grammar in user's dictionary
// @Description 
// @Tags userDictOp
// @Security APIKeyAuth
// @Accept json
// @Produce json
// @Success 200 {object} models.SuccessMsg
// @Failure 400 {object} models.ErrorMsg "Invalid JSON"
// @Failure 404 {object} models.ErrorMsg "User/Grammar not found"
// @Failure 500 {object} models.ErrorMsg "Database error"
// @Router /api/user/grammar/edit [post]
func (h* WordHandler) EditGrammarUser(c *gin.Context) {
	providedKey := c.GetHeader("X-API-Key")
	keyhash := auth.Sha256hex(providedKey)
	var user models.User
	if err := h.db.Where("keyhash =?", keyhash).
		First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(404, models.ErrorMsg{Error: "User not found"})
		} else {
			c.JSON(500, models.ErrorMsg{Error: "Database error"})
		} 
		return
	}
	var editedGrammar models.UserGrammar
	if err := c.ShouldBindJSON(&editedGrammar); err != nil {
		c.AbortWithStatusJSON(400, models.ErrorMsg{Error: "Invalid JSON"})
		return
	}
	err := h.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&models.UserGrammar{}).
			Where("id = ?", editedGrammar.ID).
			Omit("Examples").
			Updates(&editedGrammar).Error; err != nil {
			return err
		}

		editedGrammar.UserID = user.ID
		if err := tx.Where("grammar_id = ?", editedGrammar.ID).
			Delete(&models.UserGrammarExample{}).Error; err != nil {
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
			c.JSON(404, models.ErrorMsg{Error: "Grammar not found"})
		} else {
			c.JSON(500, models.ErrorMsg{Error: "Database error"})
		}
		return
	}

	c.JSON(200, models.SuccessMsg{Message: "Grammar edited"})
}


// @Summary Delete a grammar from user's dictionary
// @Description 
// @Tags userDictOp
// @Security APIKeyAuth
// @Accept json
// @Produce json
// @Success 200 {object} models.SuccessMsg
// @Failure 400 {object} models.ErrorMsg "Invalid JSON"
// @Failure 404 {object} models.ErrorMsg "User/Grammar not found"
// @Failure 500 {object} models.ErrorMsg "Database error"
// @Router /api/user/grammar/delete [post]
func (h* WordHandler) DeleteGrammarUser(c *gin.Context) {
	providedKey := c.GetHeader("X-API-Key")
	keyhash := auth.Sha256hex(providedKey)
	var user models.User
	if err := h.db.Where("keyhash =?", keyhash).
		First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(404, models.ErrorMsg{Error: "User not found"})
		} else {
			c.JSON(500, models.ErrorMsg{Error: "Database error"})
		} 
		return
	}
	var toDelete models.UserGrammar
	if err := c.ShouldBindJSON(&toDelete); err != nil {
		c.AbortWithStatusJSON(400, models.ErrorMsg{Error: "Invalid JSON"})
		return
	}
	toDelete.UserID = user.ID
	err := h.db.Transaction(func(tx *gorm.DB) error {
		var existing models.UserGrammar
		if err := tx.First(&existing, toDelete.ID).Error; err != nil {
			return err
		}
		if err := tx.Where("grammar_id = ?", toDelete.ID).
			Delete(&models.UserGrammarExample{}).
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
			c.JSON(404, models.ErrorMsg{Error: "Grammar not found"})
		} else {
			c.JSON(500, models.ErrorMsg{Error: "Database error"})
		}
		return
	}

	c.JSON(200, models.SuccessMsg{Message: "Grammar deleted"})
}

// @Summary Browse all grammars in user's dictionary
// @Description 
// @Tags userDictOp
// @Security APIKeyAuth
// @Produce json
// @Param page query int false "Page number"
// @Param RPP query int false "Results per page"
// @Success 200 {object} []models.UserGrammar
// @Failure 500 {object} models.ErrorMsg "Database error"
// @Router /api/user/grammar/get [get]
func (h* WordHandler) GetGrammarUser(c *gin.Context) {
	providedKey := c.GetHeader("X-API-Key")
	keyhash := auth.Sha256hex(providedKey)
	var user models.User
	if err := h.db.Where("keyhash =?", keyhash).
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
	query := h.db.Preload("Examples").Model(&models.UserGrammar{})
	var grammars []models.UserGrammar
	if err := query.
		Limit(resultPerPage).
		Offset((pageInt - 1) * resultPerPage).
		Where("user_id = ?", user.ID).
		Find(&grammars).Error; err != nil {
		c.AbortWithStatusJSON(500, models.ErrorMsg{Error: "Database error"})
		return
	}
}

// @Summary Search among all grammars in user's dictionary
// @Description 
// @Tags userDictOp
// @Security APIKeyAuth
// @Produce json
// @Param query query string true "Search query"
// @Param page query int false "Page number"
// @Param RPP query int false "Results per page"
// @Success 200 {object} models.SearchResult[models.UserGrammar]
// @Failure 500 {object} models.ErrorMsg "Database error"
// @Router /api/user/grammar/search [get]
func (h* WordHandler) SearchGrammarUser(c *gin.Context) {
	providedKey := c.GetHeader("X-API-Key")
	keyhash := auth.Sha256hex(providedKey)
	var user models.User
	if err := h.db.Where("keyhash =?", keyhash).
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
	var grammars []models.UserGrammar
	var count int64
	query := h.db.Preload("Examples").Model(&models.UserGrammar{})
	if err := query.
		Where("user_id = ?", user.ID).
		Where("description LIKE ?", "%"+c.Query("query")+"%").
		Count(&count).Error; err != nil {
		c.AbortWithStatusJSON(500, models.ErrorMsg{Error: "Database error"})
		return
	}
	if err := query.
		Where("user_id = ?", user.ID).
		Where("description LIKE ?", "%"+c.Query("query")+"%").
		Limit(resultPerPage).
		Offset((pageInt - 1) * resultPerPage).
		Find(&grammars).Error; err != nil {
		c.AbortWithStatusJSON(500, models.ErrorMsg{Error: "Database error"})
		return
	}
	c.JSON(200, models.SearchResult[models.UserGrammar]{
		Count: count,
		Page:  pageInt,
		PageSize: resultPerPage,
		Results: grammars,
	})
}