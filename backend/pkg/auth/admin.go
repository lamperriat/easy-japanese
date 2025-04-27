package auth

import (
	"backend/pkg/models"
	"fmt"
	"log"
	"os"

	"gorm.io/gorm"
)

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
    
    log.Println("Admin account created successfully")
    return nil
}