package reviewer

// We do not do a further abstraction because we want
// different weight algorithm for word and grammar
import (
	"backend/pkg/auth"
	"backend/pkg/models"
	"errors"
	"sort"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"math/rand"
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

const time_threshold = 120

func getReviewWordsSeq(db *gorm.DB, review_cnt int64, userID uint, batch_size int) ([]models.UserWord, error) {
	recent_threshold := review_cnt - time_threshold
	var userWords []models.UserWord
	err := db.Where("user_id = ? AND last_seen >= ? AND Familiarity > 0", userID, recent_threshold).
		Find(&userWords).Error
	if err != nil {
		return nil, err
	}
	if len(userWords) == 0 {
		return []models.UserWord{}, nil
	}
	if len(userWords) <= batch_size {
		return userWords, nil
	}
	// sort with calcWeight(familiarity, lastSeenTillNow)
	sort.Slice(userWords, func(i, j int) bool {
		return calcWeight(userWords[i].Familiarity, int(review_cnt - userWords[i].LastSeen)) < calcWeight(userWords[j].Familiarity, int(review_cnt - userWords[i].LastSeen))
	})

	return userWords[:batch_size], nil
}

// fast log2 by shifting
// log2(1) = 0
// log2(2) = 1
// log2(3) = 2
// log2(4) = 2
// log2(5) = 3 
func log2_shift(n int) int {
	n--
	if n <= 0 {
		return 0
	}
	count := 1
	for n > 1 {
		n >>= 1
		count++
	}
	return count
}

type segmentTree struct {
	tree []int
	original_start int // the start index of the original array
	// say:  22
	//    6      16
	//  3   3   4   12
	// 3 0 2 1 4 0 5 7
	// original_start is 7
}

func (st *segmentTree) build(arr []int) {
	// `tree` size (as an array): 2^h - 1
	// h = ceil(log2(n)) + 1
	h := log2_shift(len(arr)) + 1
	st.tree = make([]int, (1 << h) - 1)
	st.original_start = (1 << (h - 1)) - 1
	for i := 0; i < len(arr); i++ {
		st.tree[st.original_start + i] = arr[i]
	}
	for i := st.original_start - 1; i >= 0; i-- {
		st.tree[i] = st.tree[2 * i + 1] + st.tree[2 * i + 2]
	}
}

func (st *segmentTree) search(target_sum int) int {
	// find the index of the first element that is greater than or equal to target_sum
	n := len(st.tree)
	node := 0
	for node < n {
		left := 2 * node + 1
		right := 2 * node + 2
		if left >= n {
			break
		}
		if st.tree[left] >= target_sum {
			node = left
		} else {
			target_sum -= st.tree[left]
			node = right
		}
	}
	return node
}

func (st *segmentTree) setZero(index int) {
	// set the index to 0
	node := st.original_start + index
	st.tree[node] = 0
	for node > 0 {
		node = (node - 1) / 2
		st.tree[node] = st.tree[2 * node + 1] + st.tree[2 * node + 2]
	}
}


// we will use a segment tree based approach to do weighted random sampling.
func getReviewWordsRand(db *gorm.DB, review_cnt int64, userID uint, batch_size int) ([]models.UserWord, error) {
	recent_threshold := review_cnt - time_threshold

	var userWords []models.UserWord
	err := db.Where("user_id = ? AND last_seen >= ? AND Familiarity > 0", userID, recent_threshold).
		Find(&userWords).Error
	if err != nil {
		return nil, err
	}
	if len(userWords) == 0 {
		return []models.UserWord{}, nil
	}
	if len(userWords) <= batch_size {
		return userWords, nil
	}
	weights := make([]int, len(userWords))
	total_weight := 0
	for i := 0; i < len(userWords); i++ {
		weight := calcWeight(userWords[i].Familiarity, int(review_cnt - userWords[i].LastSeen))
		weights[i] = weight
		total_weight += weight
	}
	// build the segment tree
	st := &segmentTree{}
	st.build(weights)
	choices := make([]models.UserWord, batch_size)
	for i := 0; i < batch_size; i++ {
		weight := rand.Intn(total_weight)
		index := st.search(weight)
		choices[i] = userWords[index]
		st.setZero(index)
		total_weight -= weights[index]
		weights[index] = 0
	}
	return choices, nil
}