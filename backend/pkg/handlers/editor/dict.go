package editor

import (
	"backend/pkg/models"
	"errors"
	"fmt"
	"net/http"
	"strconv"
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

var availableDicts = map[string]struct{}{
    "book_1": {}, 
    "book_2": {},
    "book_3": {},
    "book_4": {},
    "book_5": {},
    "book_6": {},
    "all": {},
}

// @Summary Check for similar words in the dictionary
// @Description Only ``kanji" and ``katakana" fields are used for comparison
// @Tags globalDictOp
// @Security JWTAuth
// @Param dictName path string true "Dictionary name"
// @Accept json
// @Produce json
// @Success 200 {object} []models.JapaneseWord
// @Failure 400 {object} models.ErrorMsg "Invalid JSON or Invalid dict name"
// @Failure 500 {object} models.ErrorMsg "Database error"
// @Router /api/words/{dictName}/accurate-search [post]
func (h *WordHandler) AccurateSearchWord(c *gin.Context) {
    dictName := c.Param("dictName") // dictName, or "all"
    if _, ok := availableDicts[dictName]; !ok {
        c.AbortWithStatusJSON(400, models.ErrorMsg{Error: "Invalid dictionary name"})
        return
    }
    var uploadedWord models.JapaneseWord

    if err := c.ShouldBindJSON(&uploadedWord); err != nil {
        c.AbortWithStatusJSON(400, models.ErrorMsg{Error: "Invalid JSON"})
        return
    }

    query := h.db.Preload("Examples").Model(&models.JapaneseWord{})
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
        c.JSON(200, []interface{}{})
        return
    }

    query = query.Where(strings.Join(orConditions, " OR "), args...)

    var similarWords []models.JapaneseWord
    if err := query.Find(&similarWords).Error; err != nil {
        c.AbortWithStatusJSON(500, models.ErrorMsg{Error: "Database error"})
        return
    }

    c.JSON(200, similarWords)
}
// @Summary Check for similar words (fuzzy) in the dictionary
// @Description 
// @Tags globalDictOp
// @Security JWTAuth
// @Param dictName path string true "Dictionary name"
// @Param query query string true "Search query"
// @Param page query int false "Page number"
// @Param RPP query int false "Results per page"
// @Accept json
// @Produce json
// @Success 200 {object} models.SearchResult[models.JapaneseWord]
// @Failure 400 {object} models.ErrorMsg "Invalid JSON or Invalid dict name"
// @Failure 500 {object} models.ErrorMsg "Database error"
// @Router /api/words/{dictName}/fuzzy-search [get]
func (h* WordHandler) FuzzySearchWord(c *gin.Context) {
    dictName := c.Param("dictName")
    if _, ok := availableDicts[dictName]; !ok {
        c.AbortWithStatusJSON(400, models.ErrorMsg{Error: "Invalid dictionary name"})
        return
    }
	query := c.Query("query")
	page := c.Query("page")
	resultPerPageStr := c.Query("RPP")
	var pageInt int
	var err error
	pageInt, err = strconv.Atoi(page)
	if err != nil || pageInt < 1 {
		pageInt = 1
	}
	var resultPerPage int
	resultPerPage, err = strconv.Atoi(resultPerPageStr)
	if err != nil || resultPerPage < 1 || resultPerPage > 100 {
		resultPerPage = defaultResultPerPage
	}

    var words []models.JapaneseWord
    var count int64
    // log.Printf("Searching for %s in %s", query, dictName)
    // TODO: Here when this handler is called for the same query
    // but different page, the count is calculated again
    // Possible solutions: Use cache, or separate it as another handler
    q := h.db.Preload("Examples").Model(&models.JapaneseWord{})
    if dictName != "all" {
        q = q.Where("dict_name = ?", dictName)
    }
    queryStr := "%" + query + "%"
    q = q.Where("kanji LIKE ? OR katakana LIKE ? OR hiragana LIKE ? OR chinese LIKE ?", 
        queryStr, queryStr, queryStr, queryStr)
    if err := q.Count(&count).Error; err != nil {
        c.JSON(500, models.ErrorMsg{Error: "Database error"})
        return
    }

    // TODO: We may use FTS5 for better performance
    if err := q.
        Limit(resultPerPage).
        Offset((pageInt - 1) * resultPerPage).
        Find(&words).Error; err != nil {
        c.JSON(500, models.ErrorMsg{Error: "Database error"})
        return
    }

    c.JSON(200, models.SearchResult[models.JapaneseWord]{
        Count: count,
        Page: pageInt,
        PageSize: resultPerPage,
        Results: words,
    })
}


// @Summary Insert word into dictionary
// @Description 
// @Tags globalDictOp
// @Security JWTAuth
// @Param dictName path string true "Dictionary name"
// @Accept json
// @Produce json
// @Success 201 {object} models.SuccessMsg
// @Failure 400 {object} models.ErrorMsg "Invalid JSON or Invalid dict name"
// @Failure 409 {object} models.ErrorMsg "Duplicate word"
// @Failure 500 {object} models.ErrorMsg "Database error"
// @Router /api/words/{dictName}/add [post]
func (h *WordHandler) AddWord(c *gin.Context) {
    dictName := c.Param("dictName")
    if _, ok := availableDicts[dictName]; !ok {
        c.AbortWithStatusJSON(400, models.ErrorMsg{Error: "Invalid dictionary name"})
        return
    }
    var newWord models.JapaneseWord

    if err := c.ShouldBindJSON(&newWord); err != nil {
        c.AbortWithStatusJSON(400, models.ErrorMsg{Error: "Invalid JSON format"})
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
            c.JSON(409, models.ErrorMsg{Error: "Word already exists in this dictionary"})
        } else {
            c.JSON(500, models.ErrorMsg{Error: "Database operation failed"})
        }
        return
    }

    c.JSON(201, models.SuccessMsg{Message: "Word added"})
}

// @Summary Update word in dictionary
// @Description 
// @Tags globalDictOp
// @Security JWTAuth
// @Param dictName path string true "Dictionary name"
// @Accept json
// @Produce json
// @Success 200 {object} models.JapaneseWord
// @Failure 400 {object} models.ErrorMsg "Invalid JSON or Invalid dict name"
// @Failure 404 {object} models.ErrorMsg "Not found"
// @Failure 500 {object} models.ErrorMsg "Database error"
// @Router /api/words/{dictName}/edit [post]
func (h *WordHandler) EditWord(c *gin.Context) {
	dictName := c.Param("dictName")
    if _, ok := availableDicts[dictName]; !ok {
        c.AbortWithStatusJSON(400, models.ErrorMsg{Error: "Invalid dictionary name"})
        return
    }
	var editedWord models.JapaneseWord

	if err := c.ShouldBindJSON(&editedWord); err != nil {
		c.AbortWithStatusJSON(400, models.ErrorMsg{Error: "Invalid JSON format"})
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
			Omit("id", "dict_name", "Examples").
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
            c.JSON(http.StatusNotFound, models.ErrorMsg{Error:  "Word not found in dictionary"})
		} else {
            c.JSON(http.StatusInternalServerError, models.ErrorMsg{Error: "Update failed: " + err.Error()})
		}
		return
	}
	var updatedWord models.JapaneseWord
	if err := h.db.Preload("Examples").First(&updatedWord, editedWord.ID).Error; err != nil {
        c.JSON(http.StatusInternalServerError, models.ErrorMsg{Error: "Failed to retrieve updated word"})
        return
    }
	c.JSON(http.StatusOK, updatedWord)
}

// @Summary Delete word in dictionary
// @Description 
// @Tags globalDictOp
// @Security JWTAuth
// @Param dictName path string true "Dictionary name"
// @Accept json
// @Produce json
// @Success 200 {object} models.SuccessMsg
// @Failure 400 {object} models.ErrorMsg "Invalid JSON or Invalid dict name"
// @Failure 404 {object} models.ErrorMsg "Not found"
// @Failure 500 {object} models.ErrorMsg "Database error"
// @Router /api/words/{dictName}/delete [post]
func (h *WordHandler) DeleteWord(c *gin.Context) {
	dictName := c.Param("dictName")
    if _, ok := availableDicts[dictName]; !ok {
        c.AbortWithStatusJSON(400, models.ErrorMsg{Error: "Invalid dictionary name"})
        return
    }
	var wordToDelete models.JapaneseWord

	if err := c.ShouldBindJSON(&wordToDelete); err != nil {
		c.AbortWithStatusJSON(400, models.ErrorMsg{Error: "Invalid JSON format"})
		return
	}

	wordToDelete.DictName = dictName    
    
    err := h.db.Transaction(func(tx *gorm.DB) error {
        var existing models.JapaneseWord
        if err := tx.Where("id = ? AND dict_name = ?", wordToDelete.ID, dictName).
            First(&existing).Error; err != nil {
            return err
        }
        if err := tx.Where("japanese_word_id = ?", wordToDelete.ID).
            Delete(&models.ExampleSentence{}).Error; err != nil {
            return err
        }
        if err := tx.Delete(&wordToDelete).Error; err != nil {
            return err
        }
        return nil
    })

    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            c.JSON(http.StatusNotFound, models.ErrorMsg{Error: "Word not found in dictionary"})
        } else {
            c.JSON(http.StatusInternalServerError, models.ErrorMsg{Error: "Delete failed: " + err.Error()})
        }
        return
    }
    c.JSON(http.StatusOK, models.SuccessMsg{Message: "Word deleted"})
}

const defaultResultPerPage = 30

// @Summary Browse words in dictionary
// @Description 
// @Tags globalDictOp
// @Security JWTAuth
// @Param dictName path string true "Dictionary name"
// @Param page query int false "Page number"
// @Param RPP query int false "Results per page"
// @Produce json
// @Success 200 {object} models.SearchResult[models.JapaneseWord]
// @Failure 400 {object} models.ErrorMsg "Invalid dictionary name"
// @Failure 500 {object} models.ErrorMsg "Database error"
// @Router /api/words/{dictName}/get [get]
func (h *WordHandler) GetDict(c *gin.Context) {
	dictName := c.Param("dictName")
    if _, ok := availableDicts[dictName]; !ok {
        c.AbortWithStatusJSON(400, models.ErrorMsg{Error: "Invalid dictionary name"})
        return
    }
	page := c.Query("page")
	resultPerPageStr := c.Query("RPP")
	var pageInt int
	var err error
	pageInt, err = strconv.Atoi(page)
	if err != nil || pageInt < 1 {
		pageInt = 1
	}
	var resultPerPage int
	resultPerPage, err = strconv.Atoi(resultPerPageStr)
	if err != nil || resultPerPage < 1 || resultPerPage > 100 {
		resultPerPage = defaultResultPerPage
	}
	query := h.db.Preload("Examples").Model(&models.JapaneseWord{})
	if dictName != "all" {
		query = query.Where("dict_name = ?", dictName)
	}

	var words []models.JapaneseWord
    var count int64
    if err := query.Count(&count).Error; err != nil {
        c.AbortWithStatusJSON(500, models.ErrorMsg{Error: "Database error"})
        return
    }
	if err := query.
		Limit(resultPerPage).
		Offset((pageInt - 1) * resultPerPage).
		Find(&words).Error; 
		err != nil {
		c.AbortWithStatusJSON(500, models.ErrorMsg{Error: "Database error"})
		return
	}

	c.JSON(200, models.SearchResult[models.JapaneseWord]{
        Count: count,
        Page: pageInt,
        PageSize: resultPerPage,
        Results: words,
    })
}