package db

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"backend/pkg/models"
)
var realDB *gorm.DB
func InitDB() (*gorm.DB, error) {
	// singleton
	if realDB != nil {
		return realDB, nil
	}
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
		&models.UserWordExample{},
		&models.UserGrammar{},
		&models.UserGrammarExample{},
		&models.UserReadingMaterial{}, 
		&models.AdminAccount{},
		&models.ApiKey{},
	)
	realDB = db
    return db, err
}

var testDB *gorm.DB

func InitDBTest() (*gorm.DB, error) {
	// singleton
	if testDB != nil {
		return testDB, nil
	}
	db, err := gorm.Open(sqlite.Open("data/japanese_test.db"), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	println("You are using test-only database.")
	err = db.AutoMigrate(
		&models.JapaneseWord{}, 
		&models.ExampleSentence{}, 
		&models.User{}, 
		&models.UserWord{},
		&models.ReadingMaterial{},
		&models.Grammar{},
		&models.GrammarExample{},
		&models.UserWordExample{},
		&models.UserGrammar{},
		&models.UserGrammarExample{}, 
		&models.UserReadingMaterial{}, 
		&models.AdminAccount{},
		&models.ApiKey{},
	)
	testDB = db
	return db, err
}