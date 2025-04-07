package reviewer

// We do not do a further abstraction because we want
// different weight algorithm for word and grammar
import (
	"backend/pkg/models"
	"errors"
	"fmt"
	"sort"
	"strconv"

	"math/rand"

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
// @Tags reviewer 
// @Security JWTAuth
// @Accept json
// @Produce json
// @Success 200 {object} models.SuccessMsg "Success"
// @Failure 404 {object} models.ErrorMsg "User word not found"
// @Failure 400 {object} models.ErrorMsg "Invalid JSON format"
// @Failure 500 {object} models.ErrorMsg "Database error"
// @Router /api/user/review/word/correct [post]
func (h* ReviewHandler) CorrectWord(c *gin.Context) {
	var model models.UserWord
	updateLearnable(h.db, c, updateFamiliarityCorrect, &model)
}

// @Summary User incorrectly answer the word
// @Description 
// @Tags reviewer
// @Security JWTAuth
// @Accept json
// @Produce json
// @Success 200 {object} models.SuccessMsg "Success"
// @Failure 404 {object} models.ErrorMsg "User word not found"
// @Failure 400 {object} models.ErrorMsg "Invalid JSON format"
// @Failure 500 {object} models.ErrorMsg "Database error"
// @Router /api/user/review/word/incorrect [post]
func (h* ReviewHandler) IncorrectWord(c *gin.Context) {
	var model models.UserWord
	updateLearnable(h.db, c, updateFamiliarityIncorrect, &model)
}

// @Summary User correctly answer the grammar
// @Description 
// @Tags reviewer
// @Security JWTAuth
// @Accept json
// @Produce json
// @Success 200 {object} models.SuccessMsg "Success"
// @Failure 404 {object} models.ErrorMsg "User grammar not found"
// @Failure 400 {object} models.ErrorMsg "Invalid JSON format"
// @Failure 500 {object} models.ErrorMsg "Database error"
// @Router /api/user/review/grammar/correct [post]
func (h* ReviewHandler) CorrectGrammar(c *gin.Context) {
	var model models.UserGrammar
	updateLearnable(h.db, c, updateFamiliarityCorrect, &model)
}

// @Summary User incorrectly answer the grammar
// @Description 
// @Tags reviewer
// @Security JWTAuth
// @Accept json
// @Produce json
// @Success 200 {object} models.SuccessMsg "Success"
// @Failure 404 {object} models.ErrorMsg "User grammar not found"
// @Failure 400 {object} models.ErrorMsg "Invalid JSON format"
// @Failure 500 {object} models.ErrorMsg "Database error"
// @Router /api/user/review/grammar/incorrect [post]
func (h* ReviewHandler) IncorrectGrammar(c *gin.Context) {
	var model models.UserGrammar
	updateLearnable(h.db, c, updateFamiliarityIncorrect, &model)
}

// abstraction
func updateLearnable[T models.Learnable](db *gorm.DB, c *gin.Context, updateFunc func(int) int, model T) {
    keyhash_,  _ := c.Get("keyhash")
    keyhash, _ := keyhash_.(string)
    var user models.User
    if err := db.Where("keyhash = ?", keyhash).
        First(&user).Error; err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            c.JSON(404, models.ErrorMsg{Error: "User not found"})
        } else {
            c.JSON(500, models.ErrorMsg{Error: "Database error"})
        } 
        return
    }

	if err := c.ShouldBindJSON(model); err != nil {
		c.AbortWithStatusJSON(400, models.ErrorMsg{Error: "Invalid JSON format"})
		return
	}
	model.SetUserID(user.ID)
	err := db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("user_id = ? AND id = ?", user.ID, model.GetID()).
			First(model).Error; err != nil {
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
			Model(model).
			Update("familiarity", updateFunc(model.GetFamiliarity())).
			Error; err != nil {
			return err
		}
		if err := tx.
			Model(model).
			Update("last_seen", user.ReviewCount).
			Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(404, models.ErrorMsg{Error: fmt.Sprintf("User or %s found", model.GetName())})
		} else {
			c.JSON(500, models.ErrorMsg{Error: "Database error"})
		}
		return
	}
	c.JSON(200, models.SuccessMsg{Message: fmt.Sprintf("User %s updated", model.GetName())})
}

const time_threshold = 90

func getReviewLearnableSeq[T models.Learnable](db *gorm.DB, review_cnt int64, userID uint, batch_size int) ([]T, error) {
	recent_threshold := review_cnt - time_threshold
	var learnables []T
	// try to preload examples
	preload := true
	err := db.Preload("Examples").
		Where("user_id = ? AND last_seen <= ? AND Familiarity > 0", userID, recent_threshold).
		Find(&learnables).Error
	if err != nil {
		// if fail, try to find without preloading
		preload = false
		err = db.Where("user_id = ? AND Familiarity > 0", userID).
			Find(&learnables).Error
		if err != nil {
			return nil, err
		}
	}
	if len(learnables) < batch_size {
		if preload {
			err := db.Preload("Examples").
				Where("user_id = ? AND Familiarity > 0", userID).
				Find(&learnables).Error
			if err != nil {
				return nil, err
			}
		} else {
			err := db.Where("user_id = ? AND Familiarity > 0", userID).
				Find(&learnables).Error
			if err != nil {
				return nil, err
			}
		}
		if len(learnables) == 0 {
			return []T{}, nil
		}
	}
	if len(learnables) <= batch_size {
		return learnables, nil
	}
	sort.Slice(learnables, func(i, j int) bool {
		// descending order
		return calcWeight(learnables[i].GetFamiliarity(), int(review_cnt - learnables[i].GetLastSeen())) > calcWeight(learnables[j].GetFamiliarity(), int(review_cnt - learnables[j].GetLastSeen()))
	})

	return learnables[:batch_size], nil
}

func getReviewLearnableRand[T models.Learnable](db *gorm.DB, review_cnt int64, userID uint, batch_size int) ([]T, error) {
	recent_threshold := review_cnt - time_threshold
	var learnables []T
	preload := true
	err := db.Preload("Examples").
		Where("user_id = ? AND last_seen <= ? AND Familiarity > 0", userID, recent_threshold).
		Find(&learnables).Error
	if err != nil {
		preload = false
		err = db.Where("user_id = ? AND Familiarity > 0", userID).
			Find(&learnables).Error
		if err != nil {
			return nil, err
		}
	}
	if len(learnables) < batch_size {
		if preload {
			err := db.Preload("Examples").
				Where("user_id = ? AND Familiarity > 0", userID).
				Find(&learnables).Error
			if err != nil {
				return nil, err
			}
		} else {
			err := db.Where("user_id = ? AND Familiarity > 0", userID).
				Find(&learnables).Error
			if err != nil {
				return nil, err
			}
		}
		if len(learnables) == 0 {
			return []T{}, nil
		}
	}
	if len(learnables) <= batch_size {
		return learnables, nil
	}
	weights := make([]int, len(learnables))
	total_weight := 0
	for i := 0; i < len(learnables); i++ {
		weight := calcWeight(learnables[i].GetFamiliarity(), int(review_cnt - learnables[i].GetLastSeen()))
		weights[i] = weight
		total_weight += weight
	}
	// build the segment tree
	st := &segmentTree{}
	st.build(weights)
	choices := make([]T, batch_size)
	for i := 0; i < batch_size; i++ {
		weight := rand.Intn(total_weight) + 1
		index := st.search(weight)
		choices[i] = learnables[index]
		st.setZero(index)
		total_weight -= weights[index]
		weights[index] = 0
	}
	return choices, nil
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

// return: index in the original array
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
	return node - st.original_start
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



// @Summary Get batched words for review 
// @Description 
// @Tags reviewer
// @Security JWTAuth
// @Accept json
// @Produce json
// @Param batch query int false "Batch size (default 20)"
// @Param seq query bool false "Use sequential sampling (default false)"
// @Success 200 {object} []models.UserWord "Success"
// @Failure 404 {object} models.ErrorMsg "User word not found"
// @Failure 400 {object} models.ErrorMsg "Invalid JSON format"
// @Failure 500 {object} models.ErrorMsg "Database error"
// @Router /api/user/review/word/get [get]
func (h *ReviewHandler) GetWords(c *gin.Context) {
    keyhash_,  _ := c.Get("keyhash")
    keyhash, _ := keyhash_.(string)
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

	batch_size_str := c.Query("batch")
	batch_size, err := strconv.Atoi(batch_size_str)
	if err != nil || batch_size < 1 {
		batch_size = 20
	}
	seq_str := c.Query("seq")
	seq, err := strconv.ParseBool(seq_str)
	if err != nil {
		seq = false
	}
	var userWords []*models.UserWord
	if seq {
		userWords, err = getReviewLearnableSeq[*models.UserWord](h.db, user.ReviewCount, user.ID, batch_size)
	} else {
		userWords, err = getReviewLearnableRand[*models.UserWord](h.db, user.ReviewCount, user.ID, batch_size)
	}
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(404, models.ErrorMsg{Error: "User word not found"})
		} else {
			c.JSON(500, models.ErrorMsg{Error: "Database error"})
		}
		return
	}
	if len(userWords) == 0 {
		c.JSON(200, []models.UserWord{})
		return
	}
	// We do NOT need to care about updating the LastSeen of words here.
	// When user answers, `updateWord` will do the job

	c.JSON(200, userWords)
}

// @Summary Get batched words for review 
// @Description 
// @Tags reviewer
// @Security JWTAuth
// @Accept json
// @Produce json
// @Param batch query int false "Batch size (default 20)"
// @Param seq query bool false "Use sequential sampling (default false)"
// @Success 200 {object} []models.UserGrammar "Success"
// @Failure 404 {object} models.ErrorMsg "User grammar not found"
// @Failure 400 {object} models.ErrorMsg "Invalid JSON format"
// @Failure 500 {object} models.ErrorMsg "Database error"
// @Router /api/user/review/grammar/get [get]
func (h *ReviewHandler) GetGrammar(c *gin.Context) {
	// yes the function is almost exactly the same as GetWords
	// we do not do a further abstraction, sticking to the principle of
	// "low in coupling, high in cohesion"
    keyhash_,  _ := c.Get("keyhash")
    keyhash, _ := keyhash_.(string)
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

	batch_size_str := c.Query("batch")
	batch_size, err := strconv.Atoi(batch_size_str)
	if err != nil || batch_size < 1 {
		batch_size = 20
	}
	seq_str := c.Query("seq")
	seq, err := strconv.ParseBool(seq_str)
	if err != nil {
		seq = false
	}
	var userGrammars []*models.UserGrammar
	if seq {
		userGrammars, err = getReviewLearnableSeq[*models.UserGrammar](h.db, user.ReviewCount, user.ID, batch_size)
	} else {
		userGrammars, err = getReviewLearnableRand[*models.UserGrammar](h.db, user.ReviewCount, user.ID, batch_size)
	}
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(404, models.ErrorMsg{Error: "User word not found"})
		} else {
			c.JSON(500, models.ErrorMsg{Error: "Database error"})
		}
		return
	}
	if len(userGrammars) == 0 {
		c.JSON(200, []models.UserGrammar{})
		return
	}
	// We do NOT need to care about updating the LastSeen of words here.
	// When user answers, `updateWord` will do the job

	c.JSON(200, userGrammars)
}	

// since our segment tree is not exposed, we have to write unit test in the same package
func TestSegTree() {
	arr := []int{1,2,1,1,2,1,1,1,10,20,30}
	st := &segmentTree{}
	st.build(arr)
	n_choices := 5
	total := 0
	for _, v := range arr {
		total += v
	}
	fmt.Printf("tree arr len: %d\n", len(st.tree))
	for i := 0; i < n_choices; i++ {
		weight := rand.Intn(total) + 1
		index := st.search(weight)
		println("index:", index)
		st.setZero(index)
		total -= arr[index]
		arr[index] = 0
		println("tree: ", st.tree[0])
		fmt.Printf("%v\n", st.tree[1:3])
		fmt.Printf("%v\n", st.tree[3:7])
		fmt.Printf("%v\n", st.tree[7:15])
		fmt.Printf("%v\n", st.tree[15:31])
	}
}