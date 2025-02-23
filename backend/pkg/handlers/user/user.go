package user

import (
	"backend/pkg/auth"
	"backend/pkg/models"
	"errors"
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

// @Summary Register user
// @Description Both the username and the sha256 of the api key should be unique
// @Tags userOp
// @Security APIKeyAuth
// @Accept json
// @Produce json
// @Success 201 {object} models.SuccessMsg
// @Failure 400 {object} models.ErrorMsg "Invalid JSON"
// @Failure 409 {object} models.ErrorMsg "Duplicate user or username"
// @Failure 500 {object} models.ErrorMsg "Database error"
// @Router /api/user/register [post]
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

// @Summary Change the username
// @Description The sha256 of the api key is used as the identifier. 
// @Tags userOp
// @Security APIKeyAuth
// @Accept json
// @Produce json
// @Success 201 {object} models.SuccessMsg
// @Failure 400 {object} models.ErrorMsg "Invalid JSON"
// @Failure 404 {object} models.ErrorMsg "User not found"
// @Failure 409 {object} models.ErrorMsg "Duplicate user or username"
// @Failure 500 {object} models.ErrorMsg "Database error"
// @Router /api/user/update [post]
func (h* UserHandler) ChangeUserName(c *gin.Context) {
	providedKey := c.GetHeader("X-API-Key")
	keyhash := auth.Sha256hex(providedKey)
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(400, models.ErrorMsg{Error: "Invalid JSON"})
		return
	}

	err := h.db.Transaction(func(tx *gorm.DB) error {
		var originalUser models.User
		if err := tx.Where("keyhash = ?", keyhash).First(&originalUser).Error; err != nil {
			return err
		}
		if originalUser.Username == user.Username {
			return fmt.Errorf("same username")
		}
		if err := tx.Model(&originalUser).Update("username", user.Username).Error; err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(404, models.ErrorMsg{Error: "User not found"})
		} else if strings.Contains(err.Error(), "same username") {
			c.JSON(409, models.ErrorMsg{Error: "Same username"})
		} else {
			c.JSON(500, models.ErrorMsg{Error: "Database error"})
		}
		return
	}
	
	c.JSON(200, models.SuccessMsg{Message: "Username changed"})
}

// @Summary Remove user
// @Description The sha256 of the api key is used as the identifier.
// @Tags userOp
// @Security APIKeyAuth
// @Produce json
// @Success 201 {object} models.SuccessMsg
// @Failure 400 {object} models.ErrorMsg "Invalid JSON"
// @Failure 409 {object} models.ErrorMsg "Duplicate user or username"
// @Failure 500 {object} models.ErrorMsg "Database error"
// @Router /api/user/remove [get]
func (h* UserHandler) RemoveUser(c *gin.Context) {
	providedKey := c.GetHeader("X-API-Key")
	keyhash := auth.Sha256hex(providedKey)

	err := h.db.Transaction(func(tx *gorm.DB) error {
		var user models.User
		if err := tx.Where("keyhash = ?", keyhash).First(&user).Error; err != nil {
			return err
		}
		if err := tx.Delete(&user).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(404, models.ErrorMsg{Error: "User not found"})
		} else {
			c.JSON(500, models.ErrorMsg{Error: "Database error"})
		}
		return
	}
	c.JSON(200, models.SuccessMsg{Message: "User removed"})
}