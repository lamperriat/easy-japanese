package editor_test

import (
	"backend/pkg/auth"
	"backend/pkg/models"
	"backend/pkg/test"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"runtime"
	"testing"

	"github.com/gin-gonic/gin"
)

// TODO: Add unit test for user_dict related functions
func init() {
	_, filename, _, _ := runtime.Caller(0)
	dir := path.Join(path.Dir(filename), "../../..")
	err := os.Chdir(dir)
	if err != nil {
		panic(err)
	}
	os.Setenv("API_KEYS", "TEST_USE_API_KEY")
}


func TestUserDictOps(t *testing.T) {
	router := gin.Default()
	handler := test.GetTestWordHandler()
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
	// Test starts here
	word1 := models.UserWord{
		Kanji: "",
		Chinese: "测试11",
		Katakana: "テスト1",
		Hiragana: "",
		Type: "n",
		Examples: []models.UserWordExample{
			{
				Example: "This is a test sentence.",
				Chinese: "这是一个测试句子。",
			}, 
			{
				Example: "This is another test sentence.",
				Chinese: "这是另一个测试句子。",
			}, 
		}, 
	}

	word2 := models.UserWord{
		Kanji: "test2",
		Chinese: "测试2",
		Katakana: "テスト2",
		Hiragana: "",
		Type: "n",
		Examples: []models.UserWordExample{
			{
				Example: "This is a test sentence.",
				Chinese: "这是一个测试句子。",
			},
		}, 
	}

	t.Run("Add word 1", test.CreateTest(
		router, apikey, word1, "/api/user/words/add", "POST", http.StatusCreated,
	))
	t.Run("Add word 2", test.CreateTest(
		router, apikey, word2, "/api/user/words/add", "POST", http.StatusCreated,
	))
	t.Run("Accurate Search 1", test.CreateTest(
		router, apikey, models.UserWord{
			Kanji: "test2", 
		}, "/api/user/words/accurate-search", "POST", http.StatusOK,
	))
	t.Run("Accurate Search 2", test.CreateTest(
		router, apikey, models.UserWord{
			Katakana: "does not exist",
		}, "/api/user/words/accurate-search", "POST", http.StatusOK,
	))
	t.Run("Accurate Search 3", test.CreateTest(
		router, apikey, models.UserWord{
			Katakana: "テスト2",
		}, "/api/user/words/accurate-search", "POST", http.StatusOK,
	))
	t.Run("Fuzzy Search 1", test.CreateTest(
		router, apikey, gin.H{}, "/api/user/words/fuzzy-search?query=テス", "GET", http.StatusOK,
	))
	t.Run("Fuzzy Search 2", test.CreateTest(
		router, apikey, gin.H{}, "/api/user/words/fuzzy-search?query=测试", "GET", http.StatusOK,
	))
	t.Run("Fuzzy Search 3", test.CreateTest(
		router, apikey, gin.H{}, "/api/user/words/fuzzy-search?query=notExist", "GET", http.StatusOK,
	))
	t.Run("Get dict 1", test.CreateTest(
		router, apikey, gin.H{}, "/api/user/words/get", "GET", http.StatusOK,
	))
	t.Run("Edit word 1", test.CreateTest(
		router, apikey, models.UserWord{
			ID: 1, // TODO: This may fail unless we reload the db each test. 
			// Consider having some way to get the ID
			Kanji: "test1",
			Chinese: "测试1",
			Katakana: "テスト1",
			Hiragana: "",
			Type: "n",
			Examples: []models.UserWordExample{
				{
					Example: "This is a test sentence.",
					Chinese: "这是一个测试句子。",
				},
			},
		}, "/api/user/words/edit", "POST", http.StatusOK,
	))
	t.Run("Get dict 2", test.CreateTest(
		router, apikey, gin.H{}, "/api/user/words/get", "GET", http.StatusOK,
	))
	t.Run("Delete word 1", test.CreateTest(
		router, apikey, models.UserWord{
			ID: 1,
		}, "/api/user/words/delete", "POST", http.StatusOK,
	))
	t.Run("Get dict 3", test.CreateTest(
		router, apikey, gin.H{}, "/api/user/words/get", "GET", http.StatusOK,
	))
	t.Run("Delete word 2", test.CreateTest(
		router, apikey, models.UserWord{
			ID: 2,
		}, "/api/user/words/delete", "POST", http.StatusOK,
	))

	// Test stops here
	req = test.NewRequest(t, "GET", "/api/user/delete", nil, apikey)
	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	println(rr.Body.String())
	if rr.Code != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			rr.Code, http.StatusOK)
	}
}