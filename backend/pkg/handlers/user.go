package handlers

import (
	"backend/pkg/models"
	"backend/pkg/sync"
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
)

type queryType int

const (
	correct queryType = iota
	incorrect
)

func UpdateWordWeightCorrect(c *gin.Context) {
	if err := updateWordWeight(c, correct); err != nil {
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Word weight updated"})
}

func UpdateWordWeightIncorrect(c *gin.Context) {
	if err := updateWordWeight(c, incorrect); err != nil {
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Word weight updated"})
}


func updateWordWeight(c *gin.Context, query queryType) error {
	userID := c.GetHeader("X-API-KEY") // for convenience, we use the API key as the user ID
	wordIDStr := c.Param("wordId")
	wordID, err := strconv.Atoi(wordIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid word ID"})
		return err
	}

	var user models.User
	userIDInt, err := strconv.Atoi(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return err
	}
	sync.GetUserMutex(userIDInt).Lock()
	defer sync.GetUserMutex(userIDInt).Unlock()
	if err := loadUserData(userID, &user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return err
	}

	var word *models.UserWord
	for i, w := range user.Learned {
		if w.ID == wordID {
			word = &user.Learned[i]
			break
		}
	}
	if word == nil {
		user.Learned = append(user.Learned, models.UserWord{
			ID: wordID,
			Weight: models.DefaultWeight,
			UserNote: "",
		})
		word = &user.Learned[len(user.Learned)-1]
	}
	if query == correct {
		word.Weight = max(models.MinWeight, word.Weight - models.ChangeRate)
	} else {
		word.Weight = min(models.MaxWeight, word.Weight + models.ChangeRate)
	}

	if err := saveUserData(userID, &user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return err
	}
	return nil
}

func loadUserData(userID string, target *models.User) error {
	path := filepath.Join("data", "users", userID+".json")
	file, err := os.ReadFile(path)
	// create the file if it doesn't exist
	if os.IsNotExist(err) {
		target.ID = userID
		file, err := os.Create(path)
		if err != nil {
			return err
		}
		defer file.Close()
		return nil
	}
	if err != nil {
		return err
	}
	
	return json.Unmarshal(file, target)
}

func saveUserData(userID string, data *models.User) error {
	path := filepath.Join("data", "users", userID+".json")
	file, err := json.Marshal(data)
	if err != nil {
		return err
	}
	
	return os.WriteFile(path, file, 0644)
}