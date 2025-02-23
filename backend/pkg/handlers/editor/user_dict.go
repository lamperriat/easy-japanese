package editor

import (
	"backend/pkg/auth"
	"backend/pkg/models"
	"errors"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// @Summary Check for similar words in the dictionary for that user
// @Description Only ``kanji" and ``katakana" fields are used for comparison
// @Tags userDictOp
// @Security APIKeyAuth
// @Accept json
// @Produce json
// @Success 200 {object} []models.JapaneseWord
// @Failure 404 {object} models.ErrorMsg "User not found"
// @Failure 400 {object} models.ErrorMsg "Invalid JSON or Invalid dict name"
// @Failure 500 {object} models.ErrorMsg "Database error"
// @Router /api/user/words/accurate-search [post]
func (h* WordHandler) AccurateSearchWordUser(c *gin.Context) {
    providedKey := c.GetHeader("X-API-Key")
    keyhash := auth.Sha256hex(providedKey)
    var user models.User
    if err := h.db.Where("keyhash = ?", keyhash).
        First(&user).Error; err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            c.JSON(404, models.ErrorMsg{Error: "User not found"})
        } else {
            c.JSON(500, models.ErrorMsg{Error: "Database error"})
        } 
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

// @Summary Check (fuzzy search) for similar words in the dictionary for that user
// @Description 
// @Tags userDictOp
// @Security APIKeyAuth
// @Accept json
// @Produce json
// @Success 200 {object} []models.JapaneseWord
// @Failure 404 {object} models.ErrorMsg "User not found"
// @Failure 400 {object} models.ErrorMsg "Invalid JSON or Invalid dict name"
// @Failure 500 {object} models.ErrorMsg "Database error"
// @Router /api/user/words/fuzzy-search [get]
func (h* WordHandler) FuzzySearchWordUser(c *gin.Context) {
    providedKey := c.GetHeader("X-API-Key")
    keyhash := auth.Sha256hex(providedKey)
    var user models.User
    if err := h.db.Where("keyhash = ?", keyhash).
        First(&user).Error; err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            c.JSON(404, models.ErrorMsg{Error: "User not found"})
        } else {
            c.JSON(500, models.ErrorMsg{Error: "Database error"})
        } 
        return
    }

    query, page := c.Query("query"), c.Query("page")
    resultPerPageStr := c.Query("RPP")
    pageInt, err := strconv.Atoi(page)
    if err != nil || pageInt < 1 {
        pageInt = 1
    }
    resultPerPage, err := strconv.Atoi(resultPerPageStr)
    if err != nil || resultPerPage < 1 {
        resultPerPage = defaultResultPerPage
    }

    var words []models.UserWord
    var count int64
    q := h.db.Preload("Examples").
        Model(&models.UserWord{}).
        Where("user_id = ?", user.ID)

    queryStr := "%" + query + "%"
    q = q.Where("kanji LIKE ? OR katakana LIKE ? OR hiragana LIKE ? OR chinese LIKE ?", 
        queryStr, queryStr, queryStr, queryStr)
    if err := q.Count(&count).Error; err != nil {
        c.JSON(500, models.ErrorMsg{Error: "Database error"})
        return
    }

    if err := q.
        Limit(resultPerPage).
        Offset((pageInt - 1) * resultPerPage).
        Find(&words).Error; err != nil {
        c.JSON(500, models.ErrorMsg{Error: "Database error"})
        return
    }

    c.JSON(200, models.SearchResult[models.UserWord]{
        Count: count,
        Page: pageInt,
        PageSize: resultPerPage,
        Results: words,
    })

}