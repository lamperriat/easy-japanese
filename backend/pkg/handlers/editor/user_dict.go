package editor

import (
	"backend/pkg/auth"
	"backend/pkg/models"
	"strings"

	"github.com/gin-gonic/gin"
)

// @Summary Check for similar words in the dictionary for that user
// @Description Only ``kanji" and ``katakana" fields are used for comparison
// @Tags userDictOp
// @Security APIKeyAuth
// @Accept json
// @Produce json
// @Success 200 {object} []models.JapaneseWord
// @Failure 400 {object} models.ErrorMsg "Invalid JSON or Invalid dict name"
// @Failure 500 {object} models.ErrorMsg "Database error"
// @Router /api/words/{dictName}/accurate-search [post]
func (h* WordHandler) AccurateSearchWordUser(c *gin.Context) {
    providedKey := c.GetHeader("X-API-Key")
    keyhash := auth.Sha256hex(providedKey)
    var user models.User
    if err := h.db.Where("keyhash = ?", keyhash).
        First(&user).Error; err != nil {
        c.JSON(500, models.ErrorMsg{Error: "Database error"})
        return
    }
    
    var uploadedWord models.UserWord

    if err := c.ShouldBindJSON(&uploadedWord); err != nil {
        c.JSON(400, models.ErrorMsg{Error: "Invalid JSON"})
        return
    }

    query := h.db.
        Preload("Examples").
        Model(&models.UserWord{}).
        Where("user_id = ?", user.ID)

    orConditions := []string{}
    args := []interface{}{}
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

    var similarWords []models.UserWord
    if err := query.Find(&similarWords).Error; err != nil {
        c.JSON(500, models.ErrorMsg{Error: "Database error"})
        return
    }

    c.JSON(200, similarWords)
}