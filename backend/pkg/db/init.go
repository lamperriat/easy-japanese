package db

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"backend/pkg/models"
)

func InitDB() (*gorm.DB, error) {
    db, err := gorm.Open(sqlite.Open("data/japanese.db"), &gorm.Config{})
    if err != nil {
        return nil, err
    }

    err = db.AutoMigrate(
		&models.JapaneseWord{}, 
		&models.ExampleSentence{}, 
		&models.User{}, 
		&models.UserWord{},
		&models.ReadingMaterial{},
		&models.Grammar{},
		&models.GrammarExample{}, 
	)
    return db, err
}