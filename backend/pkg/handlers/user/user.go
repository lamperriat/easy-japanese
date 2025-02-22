package user

import (
	"backend/pkg/auth"
	"backend/pkg/models"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserHandler struct {
    db         *gorm.DB
}

func NewWordHandler(db *gorm.DB) *UserHandler {
    return &UserHandler{
        db:         db,
    }
}

func (h *UserHandler) RegisterUser(c *gin.Context) {
	providedKey := c.GetHeader("X-API-Key")
	keyhash := auth.Sha256hex(providedKey)
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(400, models.ErrorMsg{Error: "Invalid JSON"})
		return
	}
	err := h.db.Transaction(func(tx *gorm.DB) error {
		var existCount int64
		if err := tx.Model(&models.User{}).Where("keyhash = ?", keyhash).Count(&existCount).Error; err != nil {
			return err
		}
		if existCount > 0 {
			return fmt.Errorf("user already exists")
		}
		// check uniqueness of username
		if err := tx.Model(&models.User{}).Where("username = ?", user.Username).Count(&existCount).Error; err != nil {
			return err
		}
		if existCount > 0 {
			return fmt.Errorf("username already exists")
		}
		user.Keyhash = keyhash
		user.ID = 0
		user.ReviewCount = 0
		if err := tx.Create(&user).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		if strings.Contains(err.Error(), "exists") {
			c.JSON(409, models.ErrorMsg{Error: err.Error()})
			return
		} else {
			c.JSON(500, models.ErrorMsg{Error: "Database error"})
			return
		}
	}

	c.JSON(201, models.SuccessMsg{Message: "User registered"})
}