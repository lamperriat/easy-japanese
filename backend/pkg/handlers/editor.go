package handlers

import (
	"backend/pkg/models"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type WordHandler struct {
    db         *gorm.DB
}

func NewWordHandler(db *gorm.DB) *WordHandler {
    return &WordHandler{
        db:         db,
    }
}


func (h *WordHandler) CheckSimilarWords(c *gin.Context) {
    dictName := c.Param("dictName") // dictName, or "all"
    var uploadedWord models.JapaneseWord

    if err := c.ShouldBindJSON(&uploadedWord); err != nil {
        c.AbortWithStatusJSON(400, gin.H{"error": "Invalid JSON"})
        return
    }

    query := h.db.Model(&models.JapaneseWord{})
    if dictName != "all" {
        query = query.Where("dict_name = ?", dictName)
    }

    orConditions := []string{}
    args := make([]interface{}, 0)
    if uploadedWord.Kanji != "" {
        orConditions = append(orConditions, "kanji = ?")
        args = append(args, uploadedWord.Kanji)
    }
    if uploadedWord.Katakana != "" {
        orConditions = append(orConditions, "katakana = ?")
        args = append(args, uploadedWord.Katakana)
    }

    if len(orConditions) == 0 {
        c.JSON(200, gin.H{"similar": []interface{}{}})
        return
    }

    query = query.Where(strings.Join(orConditions, " OR "), args...)

    var similarWords []models.JapaneseWord
    if err := query.Find(&similarWords).Error; err != nil {
        c.AbortWithStatusJSON(500, gin.H{"error": "Database error"})
        return
    }

    c.JSON(200, gin.H{
        "similar": similarWords,
        "dict":    dictName,
    })
}

func (h *WordHandler) AddWord(c *gin.Context) {
    dictName := c.Param("dictName")
    var newWord models.JapaneseWord

    if err := c.ShouldBindJSON(&newWord); err != nil {
        c.AbortWithStatusJSON(400, gin.H{"error": "Invalid JSON format"})
        return
    }
    newWord.DictName = dictName

    err := h.db.Transaction(func(tx *gorm.DB) error {
        var existCount int64
        if err := tx.Model(&models.JapaneseWord{}).
            Where("dict_name = ? AND ((kanji != '' AND kanji = ?) OR (katakana != '' AND katakana = ?))", 
            dictName, 
            newWord.Kanji, 
            newWord.Katakana).
            Count(&existCount).Error; err != nil {
            return err
        }

        if existCount > 0 {
            return fmt.Errorf("duplicate word in dictionary")
        }

        if err := tx.Create(&newWord).Error; err != nil {
            return err
        }

        // if len(newWord.Examples) > 0 {
        //     log.Printf("Adding %d examples", len(newWord.Examples))
        //     for i := range newWord.Examples {
        //         newWord.Examples[i].JapaneseWordID = newWord.ID
        //         newWord.Examples[i].ID = 0
        //     }
        //     if err := tx.Create(&newWord.Examples).Error; err != nil {
        //         return err
        //     }
        // }
        return nil
    })

    if err != nil {
        if strings.Contains(err.Error(), "duplicate") {
            c.JSON(409, gin.H{"error": "Word already exists in this dictionary"})
        } else {
            c.JSON(500, gin.H{"error": "Database operation failed"})
        }
        return
    }

    c.JSON(201, gin.H{
        "id":   newWord.ID,
        "dict": dictName,
    })
}

func (h *WordHandler) EditWord(c *gin.Context) {
	dictName := c.Param("dictName")
	var editedWord models.JapaneseWord

	if err := c.ShouldBindJSON(&editedWord); err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": "Invalid JSON format"})
		return
	}
	editedWord.DictName = dictName

	err := h.db.Transaction(func(tx *gorm.DB) error {
		var existing models.JapaneseWord
		if err := tx.Where("id = ? AND dict_name = ?", editedWord.ID, dictName).
			First(&existing).Error; err != nil {
			return err
		}
		if err := tx.Model(&existing).
			Omit("id", "dict_name").
			Updates(&editedWord).Error; err != nil {
			return err
		}

		if err := tx.Where("japanese_word_id = ?", editedWord.ID).
			Delete(&models.ExampleSentence{}).Error; err != nil {
			return err
		}

        if len(editedWord.Examples) > 0 {
            for i := range editedWord.Examples {
                editedWord.Examples[i].JapaneseWordID = editedWord.ID
                editedWord.Examples[i].ID = 0
            }
            if err := tx.Create(&editedWord.Examples).Error; err != nil {
                return err
            }
        }
		return nil
	})

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
            c.JSON(http.StatusNotFound, gin.H{
                "error": fmt.Sprintf("Word %d not found in %s", editedWord.ID, dictName),
            })
		} else {
            c.JSON(http.StatusInternalServerError, gin.H{
                "error": "Update failed: " + err.Error(),
            })
		}
		return
	}
	var updatedWord models.JapaneseWord
	h.db.Preload("Examples").First(&updatedWord, editedWord.ID)
	c.JSON(http.StatusOK, gin.H{
		"data": updatedWord,
		"message": "Word updated",
	})
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

// func saveDict(dictName string, data []models.JapaneseWord) error {
// 	dictName = dictName + ".json"
// 	path := filepath.Join("data", "japanese", dictName)
// 	file, err := json.MarshalIndent(data, "", "  ")
// 	if err != nil {
// 		return err
// 	}

// 	return os.WriteFile(path, file, 0644)
// }