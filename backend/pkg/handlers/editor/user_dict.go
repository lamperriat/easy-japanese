package editor

import (
	"backend/pkg/auth"
	"backend/pkg/models"
	"errors"
	"fmt"
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


// @Summary Insert word into user's dictionary
// @Description 
// @Tags userDictOp
// @Security APIKeyAuth
// @Accept json
// @Produce json
// @Success 201 {object} models.SuccessMsg
// @Failure 400 {object} models.ErrorMsg "Invalid JSON"
// @Failure 404 {object} models.ErrorMsg "User not found"
// @Failure 409 {object} models.ErrorMsg "Duplicate word"
// @Failure 500 {object} models.ErrorMsg "Database error"
// @Router /api/user/words/add [post]
func (h* WordHandler) AddWordUser(c *gin.Context) {
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

    var newWord models.UserWord

    if err := c.ShouldBindJSON(&newWord); err != nil {
        c.JSON(400, models.ErrorMsg{Error: "Invalid JSON"})
        return
    }

    err := h.db.Transaction(func(tx *gorm.DB) error {
        var existCount int64
        if err := tx.Model(&models.UserWord{}).
            Where("(kanji != '' AND kanji = ?) OR (katakana != '' AND katakana = ?)",
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
        return nil
    })
   
    if err != nil {
        if strings.Contains(err.Error(), "duplicate") {
            c.JSON(409, models.ErrorMsg{Error: "Word already exists in this dictionary"})
        } else {
            c.JSON(500, models.ErrorMsg{Error: "Database operation failed"})
        }
        return
    }

    c.JSON(201, models.SuccessMsg{Message: "Word added"})
}

// @Summary Edit a word in user's dictionary
// @Description 
// @Tags userDictOp
// @Security APIKeyAuth
// @Accept json
// @Produce json
// @Success 200 {object} models.SuccessMsg
// @Failure 400 {object} models.ErrorMsg "Invalid JSON"
// @Failure 404 {object} models.ErrorMsg "User not found"
// @Failure 500 {object} models.ErrorMsg "Database error"
// @Router /api/user/words/edit [post]
func (h* WordHandler) EditWordUser(c *gin.Context) {
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

    var newWord models.UserWord

    if err := c.ShouldBindJSON(&newWord); err != nil {
        c.JSON(400, models.ErrorMsg{Error: "Invalid JSON"})
        return
    }

    err := h.db.Transaction(func(tx* gorm.DB) error {
        var existing models.UserWord
        if err := tx.Model(&models.UserWord{}).
            Where("id = ? AND user_id = ?", newWord.ID, user.ID).
            First(&existing).Error; err != nil {
            return err
        }

        if err := tx.Model(&existing).
            Omit("id", "user_id", "Examples").
            Updates(&newWord).Error; err != nil {
            return err
        }
        if err := tx.Where("user_word_id = ?", newWord.ID).
            Delete(&models.UserWordExample{}).Error; err != nil {
            return err
        }
        if len(newWord.Examples) > 0 {
            for i := range newWord.Examples {
                newWord.Examples[i].UserWordID = newWord.ID
                newWord.Examples[i].ID = 0
            }
            if err := tx.Create(&newWord.Examples).Error; err != nil {
                return err
            }
        }
        return nil
    })

    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            c.JSON(404, models.ErrorMsg{Error: "Word not found"})
        } else {
            c.JSON(500, models.ErrorMsg{Error: "Database operation failed"})
        }
        return
    }

    c.JSON(200, models.SuccessMsg{Message: "Word updated"})
}


// @Summary Delete a word in user's dictionary
// @Description 
// @Tags userDictOp
// @Security APIKeyAuth
// @Accept json
// @Produce json
// @Success 200 {object} models.SuccessMsg
// @Failure 400 {object} models.ErrorMsg "Invalid JSON"
// @Failure 404 {object} models.ErrorMsg "User not found"
// @Failure 500 {object} models.ErrorMsg "Database error"
// @Router /api/user/words/delete [post]
func (h *WordHandler) DeleteWordUser(c *gin.Context) {
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
    var toDelete models.UserWord

    if err := c.ShouldBindJSON(&toDelete); err != nil {
        c.JSON(400, models.ErrorMsg{Error: "Invalid JSON"})
        return
    }

    err := h.db.Transaction(func(tx *gorm.DB) error {
        var existing models.UserWord
        if err := tx.Model(&models.UserWord{}).
            Where("id = ? AND user_id = ?", toDelete.ID, user.ID).
            First(&existing).Error; err != nil {
            return err
        }
        if err := tx.Where("user_word_id = ?", toDelete.ID).
            Delete(&models.UserWordExample{}).Error; err != nil {
            return err
        }
        if err := tx.Delete(&existing).Error; err != nil {
            return err
        }
        return nil
    })

    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            c.JSON(404, models.ErrorMsg{Error: "Word not found"})
        } else {
            c.JSON(500, models.ErrorMsg{Error: "Database operation failed"})
        }
        return
    }
    c.JSON(200, models.SuccessMsg{Message: "Word deleted"})
}