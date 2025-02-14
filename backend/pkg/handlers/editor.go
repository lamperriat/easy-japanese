package handlers

import (
	"backend/pkg/models"
	"encoding/json"
	"log"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)


func CheckSimilarWords(c *gin.Context) {
	dictName := c.Param("dictName")
	words, err := loadDict(dictName)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to load dictionary"})
		return
	}
	var uploadedWord models.JapaneseWord
	if err := c.ShouldBindJSON(&uploadedWord); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request body"})
		return
	}
	var similarWords []string
	for _, word := range words {
		if word.Kanji != "" && word.Kanji == uploadedWord.Kanji {
			similarWords = append(similarWords, word.Kanji)
		} else if word.Katakana != "" && word.Katakana == uploadedWord.Katakana {
			similarWords = append(similarWords, word.Katakana)
		}
	}
	msg := "Found similar words: "
	for _, word := range similarWords {
		msg += word + ", "
	}
	c.JSON(200, gin.H{"message": msg})
}

func AddWord(c *gin.Context) {
	dictName := c.Param("dictName")
	words, err := loadDict(dictName)
	log.Println(dictName)
	if err != nil {
		log.Println(err)
		c.JSON(500, gin.H{"error": "Failed to load dictionary"})
		return
	}
	var newWord models.JapaneseWord
	if err := c.ShouldBindJSON(&newWord); err != nil {
		log.Println(err)
		c.JSON(400, gin.H{"error": "Invalid request body"})
		return
	}
	newWord.ID, err = models.GetNextID()
	if err != nil {
		log.Println(err)
		c.JSON(500, gin.H{"error": "Failed to get next ID"})
		return
	}
	words = append(words, newWord)
	if err := saveDict(dictName, words); err != nil {
		log.Println(err)
		c.JSON(500, gin.H{"error": "Failed to save dictionary"})
		return
	}
	saveDict(dictName, words)
	c.JSON(200, gin.H{"message": "Word added successfully"})
}

func loadDict(dictName string) ([]models.JapaneseWord, error) {
	dictName = dictName + ".json"
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
	dictName = dictName + ".json"
	path := filepath.Join("data", "japanese", dictName)
	file, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, file, 0644)
}