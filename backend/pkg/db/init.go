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
<<<<<<< HEAD
		&models.UserWord{},
		&models.ReadingMaterial{},
		&models.Grammar{},
		&models.GrammarExample{}, 
=======
		&models.UserWord{}, 
>>>>>>> 81d02e8 (merge: Update main with sqlite features (#6))
	)
    return db, err
}