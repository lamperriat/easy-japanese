package reviewer
// We do not do a further abstraction because we want
// different weight algorithm for word and grammar
import (
	"backend/pkg/auth"
	"backend/pkg/models"
	"errors"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ReviewHandler struct {
	db *gorm.DB
}
func NewReviewHandler(db *gorm.DB) *ReviewHandler {
	return &ReviewHandler{db: db}
}

func updateFamiliarityCorrect(familiarity int) int {
	if familiarity <= 0 {
		return 0
	} else if familiarity < 15 {
		return familiarity - 1
	} else if familiarity < 80 {
		return familiarity - 2
	} else if familiarity < 120 {
		return familiarity - 3
	} else {
		return familiarity - 5
	}
}

func updateFamiliarityIncorrect(familiarity int) int {
	if familiarity <= 0 {
		return 5
	} else if familiarity < 10 {
		return familiarity + 5
	} else if familiarity < 80 {
		return familiarity + 3
	} else {
		return familiarity + 1
	}
}

func calcWeight(familiarity int, lastSeenTillNow int) int {
	clamped_time := min(lastSeenTillNow, 1500)
	clamped_familiarity := min(familiarity, 80)
	return clamped_time * clamped_time / 10000 + clamped_familiarity * clamped_familiarity
}

// @Summary User correctly answer the word
// @Description 
// @Tags globalDictOp
// @Security APIKeyAuth
// @Accept json
// @Produce json
// @Success 200 {object} models.SuccessMsg "Success"
// @Failure 404 {object} models.ErrorMsg "User word not found"
// @Failure 400 {object} models.ErrorMsg "Invalid JSON format"
// @Failure 500 {object} models.ErrorMsg "Database error"
// @Router /api/user/review/correct [post]
func (h* ReviewHandler) CorrectWord(c *gin.Context) {
	h.updateWord(c, true)
}

// @Summary User incorrectly answer the word
// @Description 
// @Tags globalDictOp
// @Security APIKeyAuth
// @Accept json
// @Produce json
// @Success 200 {object} models.SuccessMsg "Success"
// @Failure 404 {object} models.ErrorMsg "User word not found"
// @Failure 400 {object} models.ErrorMsg "Invalid JSON format"
// @Failure 500 {object} models.ErrorMsg "Database error"
// @Router /api/user/review/incorrect [post]
func (h* ReviewHandler) IncorrectWord(c *gin.Context) {
	h.updateWord(c, false)
}

func (h* ReviewHandler) updateWord(c *gin.Context, correct bool) {
	var updateFunc func(int) int
	if correct {
		updateFunc = updateFamiliarityCorrect
	} else {
		updateFunc = updateFamiliarityIncorrect
	}
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

	var userWord models.UserWord
	if err := c.ShouldBindJSON(&userWord); err != nil {
		c.AbortWithStatusJSON(400, models.ErrorMsg{Error: "Invalid JSON format"})
		return
	}
	userWord.UserID = user.ID
	err := h.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("user_id = ? AND id = ?", user.ID, userWord.ID).
			First(&userWord).Error; err != nil {
				return err
		}
		// 3 steps: 
		// 1. Get and update User.ReviewCount
		// 2. Get and update UserWord.Familiarity
		// 3. Update UserWord.LastSeen
		if err := tx.Model(&user).Update("review_count", user.ReviewCount + 1).Error; err != nil {
			return err
		}
		if err := tx.
			Model(&userWord).
			Update("familiarity", updateFunc(userWord.Familiarity)).
			Error; err != nil {
			return err
		}
		if err := tx.
			Model(&userWord).
			Update("last_seen", user.ReviewCount).
			Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(404, models.ErrorMsg{Error: "User word not found"})
		} else {
			c.JSON(500, models.ErrorMsg{Error: "Database error"})
		}
		return
	}
	c.JSON(200, models.SuccessMsg{Message: "User word updated"})
}