package auth

import (
	"backend/pkg/models"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func registerAdmin(db *gorm.DB, username, password string) error {
    if username == "" || password == "" {
        return fmt.Errorf("ADMIN_USERNAME and ADMIN_PASSWORD environment variables must be set")
    }

	if len(password) < 12 {
		return fmt.Errorf("password must be at least 12 characters long")
	}

    hashedPassword, err := SafeHash(password)
    if err != nil {
        return fmt.Errorf("failed to hash admin password: %w", err)
    }
    
	admin := models.AdminAccount{
		Username:     username,
		PasswordHash: hashedPassword,
	}
    
    if err := db.Create(&admin).Error; err != nil {
        return fmt.Errorf("failed to create admin account: %w", err)
    }
    return nil
}


func InitAdminAccount(db *gorm.DB) error {
    // Check for existing admin
    var adminCount int64
    if err := db.Model(&models.AdminAccount{}).Count(&adminCount).Error; err != nil {
        return fmt.Errorf("failed to check for existing admin: %w", err)
    }
    
    if adminCount > 0 {
        log.Println("Admin account already exists, skipping initialization")
        return nil
    }
    
    username := os.Getenv("EASYJP_ADMIN_USERNAME")
    password := os.Getenv("EASYJP_ADMIN_PASSWORD")
    
	if err := registerAdmin(db, username, password); err != nil {
		return err
	}

	log.Printf("Default admin username: %s\n", username)
    log.Println("Admin account created successfully")
    return nil
}

func AdminAuth(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		username := c.GetHeader("Admin-Name")
		password := c.GetHeader("Admin-Password")
		if username == "" || password == "" {
			c.JSON(401, models.ErrorMsg{Error: "Admin-Name and Admin-Password headers are required"})
			c.Abort()
			return
		}
		var admin models.AdminAccount
		if err := db.Where("username = ?", username).First(&admin).Error; err != nil {
			c.JSON(401, models.ErrorMsg{Error: "Invalid admin credentials"})
			c.Abort()
			return
		}
		if !SafeCompare(password, admin.PasswordHash) {
			c.JSON(401, models.ErrorMsg{Error: "Invalid admin credentials"})
			c.Abort()
			return
		}
		c.Next()
	}
}

type AdminAccountJSON struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Admin Operations

// @Summary Create admin account
// @Description Create a new admin account
// @Tags admin
// @Security AdminAuth
// @Accept json
// @Produce json
// @Success 200 {object} models.SuccessMsg
// @Failure 400 {object} models.ErrorMsg "Invalid JSON"
// @Failure 500 {object} models.ErrorMsg "Database error"
// @Router /admin-api/account/create [post]
func CreateAdminAccount(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var adminAccount AdminAccountJSON
		if err := c.ShouldBindJSON(&adminAccount); err != nil {
			c.JSON(400, models.ErrorMsg{Error: "Invalid JSON"})
			return
		}
		if err := registerAdmin(db, adminAccount.Username, adminAccount.Password); err != nil {
			c.JSON(500, models.ErrorMsg{Error: err.Error()})
			return
		}
		c.JSON(200, models.SuccessMsg{Message: "Admin account created successfully"})
	}
}

type ApiKeyJSON struct {
	Key string `json:"key"`
}

// @Summary Create new apikey
// @Description 
// @Tags admin
// @Security AdminAuth
// @Produce json
// @Success 200 {object} ApiKeyJSON
// @Failure 400 {object} models.ErrorMsg "Invalid JSON"
// @Failure 500 {object} models.ErrorMsg "Database error"
// @Router /admin-api/apikey/create [GET]
func GenerateApiKey(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		newKey, err := SafeRandom(256)
		if err != nil {
			c.JSON(500, models.ErrorMsg{Error: err.Error()})
			return
		}
		keyHash, err := SafeHash(newKey)
		if err != nil {
			c.JSON(500, models.ErrorMsg{Error: err.Error()})
			return
		}
		apiKey := models.ApiKey{
			KeyHash: keyHash, 
		}
		if err := db.Create(&apiKey).Error; err != nil {
			c.JSON(500, models.ErrorMsg{Error: err.Error()})
			return
		}
		c.JSON(200, ApiKeyJSON{Key: newKey})
	}
}