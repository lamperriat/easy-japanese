package editor

import (
	"backend/pkg/auth"
	"backend/pkg/models"
	"errors"

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