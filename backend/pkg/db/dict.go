package db

import (
	"gorm.io/gorm"
	"backend/pkg/models"
)

func CreateWord(db *gorm.DB, word *models.JapaneseWord) error {
    return db.Transaction(func(tx *gorm.DB) error {
        if err := tx.Create(word).Error; err != nil {
            return err
        }

        for i := range word.Examples {
            word.Examples[i].JapaneseWordID = word.ID
        }
        return tx.Create(&word.Examples).Error
    })
}

func SearchMatchingWords(db *gorm.DB, word *models.JapaneseWord) ([]*models.JapaneseWord, error) {
	var words []*models.JapaneseWord
	if err := db.Where("kanji = ? OR katakana = ?", word.Kanji, word.Katakana).Find(&words).Error; err != nil {
		return nil, err
	}
	return words, nil
}

