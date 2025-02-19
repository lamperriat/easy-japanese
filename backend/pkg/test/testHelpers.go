package test

import (
	"backend/pkg/handlers"
	"backend/pkg/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func initDB() (*gorm.DB, error) {
    db, err := gorm.Open(sqlite.Open("../../data/japanese.db"), &gorm.Config{})
    if err != nil {
        return nil, err
    }

    err = db.AutoMigrate(
		&models.JapaneseWord{}, 
		&models.ExampleSentence{}, 
		&models.User{}, 
		&models.UserWord{}, 
	)
    return db, err
}
func GetTestDB() *gorm.DB {
	db, err := initDB()
	if err != nil {
		panic(err)
	}
	return db
}

func GetTestWordHandler() *handlers.WordHandler {
	return handlers.NewWordHandler(GetTestDB())
}
