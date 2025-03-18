package reviewer_test

import "fmt"

import (
	"backend/pkg/auth"
	"backend/pkg/handlers/reviewer"
	"backend/pkg/models"
	"backend/pkg/test"
	"encoding/json"
	"math"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"runtime"
	"testing"

	"math/rand"

	"github.com/gin-gonic/gin"
)


func log_builtin(n int) int {
	return int(math.Ceil(math.Log2(float64(n))))
}

func log_shift(n int) int {
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
func BenchmarkLogBuiltin(b *testing.B) {
	input := generateInput(b.N)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		log_builtin(input[i])
	}
}

func BenchmarkLogShift(b *testing.B) {
	input := generateInput(b.N)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		log_shift(input[i])
	}
}

func TestLogShift2(t *testing.T) {
	t.Run("log 1", func(t *testing.T) {
		if log_shift(1) != 0 {
			t.Errorf("Expected 0, got %d", log_shift(1))
		}
	})
	t.Run("log 2", func(t *testing.T) {
		if log_shift(2) != 1 {
			t.Errorf("Expected 1, got %d", log_shift(2))
		}
	})
	t.Run("log 3", func(t *testing.T) {
		if log_shift(3) != 2 {
			t.Errorf("Expected 2, got %d", log_shift(3))
		}
	})
	t.Run("log 4", func(t *testing.T) {
		if log_shift(4) != 2 {
			t.Errorf("Expected 2, got %d", log_shift(4))
		}
	})
}

func generateInput(n int) []int {
	r := rand.New(rand.NewSource(123)) // Seed the random number generator
	input := make([]int, n)
	for i := 0; i < n; i++ {
		input[i] = r.Intn(1000000) // Generate random integers
	}
	return input
}


// below is the test for our algorithm
func init() {
	_, filename, _, _ := runtime.Caller(0)
	dir := path.Join(path.Dir(filename), "../../..")
	err := os.Chdir(dir)
	if err != nil {
		panic(err)
	}
	os.Setenv("API_KEYS", "TEST_USE_API_KEY")
}

func TestReviewWord(t *testing.T) {
	router := gin.Default()
	handler := test.GetTestWordHandler()
	reviewer := reviewer.NewReviewHandler(test.GetTestDB())
	defer func() {
		err := os.Remove("data/japanese_test.db")
		if err != nil {
			panic(err)
		}
	} ()
	handler_user := test.GetTestUserHandler()
	router.POST("/api/user/register", auth.APIKeyAuth(), handler_user.RegisterUser)
	router.POST("/api/user/update", auth.APIKeyAuth(), handler_user.UpdateUserName)
	router.GET("/api/user/delete", auth.APIKeyAuth(), handler_user.DeleteUser)
	router.POST("/api/user/words/accurate-search", auth.APIKeyAuth(), handler.AccurateSearchWordUser)
	router.GET("/api/user/words/fuzzy-search", auth.APIKeyAuth(), handler.FuzzySearchWordUser)
	router.POST("/api/user/words/add", auth.APIKeyAuth(), handler.AddWordUser)
	router.POST("/api/user/words/edit", auth.APIKeyAuth(), handler.EditWordUser)
	router.POST("/api/user/words/delete", auth.APIKeyAuth(), handler.DeleteWordUser)
	router.GET("/api/user/words/get", auth.APIKeyAuth(), handler.GetDictUser)
	router.GET("/api/user/review/get", auth.APIKeyAuth(), reviewer.GetWords)
	router.POST("/api/user/review/correct", auth.APIKeyAuth(), reviewer.CorrectWord)
	router.POST("/api/user/review/incorrect", auth.APIKeyAuth(), reviewer.IncorrectWord)
	

	apikey := "TEST_USE_API_KEY"
	test_user := models.User{
		Username: "test_user",
	}

	body, err := json.Marshal(test_user)
	if err != nil {
		t.Fatalf("Failed to marshal user: %v", err)
	}
	req := test.NewRequest(t, "POST", "/api/user/register", body, apikey)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	println(rr.Body.String())
	for i := 0; i < 100; i++ {
		idx := fmt.Sprint(i)
		word := models.UserWord{
			Kanji: "",
			Chinese: "测试" + idx,
			Katakana: "テスト" + idx,
			Hiragana: "",
			Type: "n",
			Examples: []models.UserWordExample{
				{
					Example: "This is a test sentence." + idx,
					Chinese: "这是一个测试句子。" + idx,
				}, 
				{
					Example: "This is another test sentence." + idx,
					Chinese: "这是另一个测试句子。" + idx,
				}, 
			}, 
		}
		t.Run("Add word", test.CreateTest(
			router, apikey, word, "/api/user/words/add", "POST", http.StatusCreated, 
		))
	}
	t.Run("Get all", test.CreateTest(
		router, apikey, nil, "/api/user/words/get?RPP=100", "GET", http.StatusOK, 
	))
	t.Run("Get words", test.CreateTest(
		router, apikey, nil, "/api/user/review/get", "GET", http.StatusOK,
	))
	t.Run("Get words seq", test.CreateTest(
		router, apikey, nil, "/api/user/review/get?seq=true", "GET", http.StatusOK,
	))

	// Test ends here
	req = test.NewRequest(t, "GET", "/api/user/delete", nil, apikey)
	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	println(rr.Body.String())
	if rr.Code != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			rr.Code, http.StatusOK)
	}
}