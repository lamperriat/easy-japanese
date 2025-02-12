package handlers

import (
	"backend/pkg/models"
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)


func CheckSimilarWords(c *gin.Context) {

}

func loadDict(dictName string) ([]models.JapaneseWord, error) {
	path := filepath.Join("data", "japanese", dictName)
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var words []models.JapaneseWord
	if err := json.Unmarshal(file, &words); err != nil {
		return nil, err
	}
	return words, nil
}

func saveDict(dictName string, data []models.JapaneseWord) error {
	path := filepath.Join("data", "japanese", dictName)
	file, err := json.Marshal(data)
	if err != nil {
		return err
	}

	return os.WriteFile(path, file, 0644)
}