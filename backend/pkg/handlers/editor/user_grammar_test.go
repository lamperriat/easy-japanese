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

func init() {
	_, filename, _, _ := runtime.Caller(0)
	dir := path.Join(path.Dir(filename), "../../..")
	err := os.Chdir(dir)
	if err != nil {
		panic(err)
	}
	os.Setenv("API_KEYS", "TEST_USE_API_KEY")
}

func TestUserGrammarOps(t *testing.T) {
	router := gin.Default()
	handler := test.GetTestWordHandler()
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
	router.POST("/api/user/grammar/add", auth.APIKeyAuth(), handler.AddGrammarUser)
	router.POST("/api/user/grammar/edit", auth.APIKeyAuth(), handler.EditGrammarUser)
	router.POST("/api/user/grammar/delete", auth.APIKeyAuth(), handler.DeleteGrammarUser)
	router.GET("/api/user/grammar/get", auth.APIKeyAuth(), handler.GetGrammarUser)
	router.GET("/api/user/grammar/search", auth.APIKeyAuth(), handler.SearchGrammarUser)
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
	grammar1 := models.UserGrammar{
		Description: "测试语法1",
		Examples: []models.UserGrammarExample{
			{
				Example: "Test example 1",
				Chinese: "测试例句中文1",
			}, 
			{
				Example: "Test example 2",
				Chinese: "测试例句中文2",
			}, 
			{
				Example: "Test example 3",
				Chinese: "测试例句中文3",
			},
		},
	}
	grammar2 := models.UserGrammar{
		Description: "测试语法2",
		Examples: []models.UserGrammarExample{}, 
	};
	grammar3 := models.UserGrammar{
		Description: "搜索关键词",
		Examples: []models.UserGrammarExample{
			{
				Example: "Test example 1",
				Chinese: "测试例句中文1",
			}, 
		}, 
	}
	t.Run("AddGrammar 1", test.CreateTest(
		router, apikey, grammar1, "/api/user/grammar/add", "POST", http.StatusCreated, 
	))
	t.Run("AddGrammar 2", test.CreateTest(
		router, apikey, grammar2, "/api/user/grammar/add", "POST", http.StatusCreated,
	))
	t.Run("AddGrammar 3", test.CreateTest(
		router, apikey, grammar3, "/api/user/grammar/add", "POST", http.StatusCreated,
	))
	t.Run("GetGrammar", test.CreateTest(
		router, apikey, nil, "/api/user/grammar/get", "GET", http.StatusOK,
	))
	t.Run("SearchGrammar", test.CreateTest(
		router, apikey, nil, "/api/user/grammar/search?query=关键", "GET", http.StatusOK,
	))
	t.Run("SearchGrammar 2", test.CreateTest(
		router, apikey, nil, "/api/user/grammar/search?query=nonexist", "GET", http.StatusOK,
	))
	t.Run("EditGrammar", test.CreateTest(
		router, apikey, models.UserGrammar{
			ID: 1,
			Description: "edited grammar",
			Examples: []models.UserGrammarExample{
				{
					Example: "edited example 1",
					Chinese: "编辑后的例句中文1",
				}, 
				{
					Example: "edited example 2",
					Chinese: "编辑后的例句中文2",
				},
			},
		}, "/api/user/grammar/edit", "POST", http.StatusOK,
	))
	t.Run("EditGrammar 2", test.CreateTest(
		router, apikey, models.UserGrammar{
			ID: 3,
			Description: "edited grammar 2",
			Examples: []models.UserGrammarExample{}, 
		}, "/api/user/grammar/edit", "POST", http.StatusOK,
	))
	t.Run("Current Grammar", test.CreateTest(
		router, apikey, nil, "/api/user/grammar/get", "GET", http.StatusOK,
	))
	t.Run("DeleteGrammar", test.CreateTest(
		router, apikey, models.UserGrammar{
			ID: 1,
		}, "/api/user/grammar/delete", "POST", http.StatusOK,
	))
	t.Run("DeleteGrammar 2", test.CreateTest(
		router, apikey, models.UserGrammar{
			ID: 2,
		}, "/api/user/grammar/delete", "POST", http.StatusOK,
	))
	t.Run("Current Grammar", test.CreateTest(
		router, apikey, nil, "/api/user/grammar/get", "GET", http.StatusOK,
	))
	t.Run("DeleteGrammar 3", test.CreateTest(
		router, apikey, models.UserGrammar{
			ID: 3,
		}, "/api/user/grammar/delete", "POST", http.StatusOK,
	))
	t.Run("Current Grammar", test.CreateTest(
		router, apikey, nil, "/api/user/grammar/get", "GET", http.StatusOK,
	))
	// Test ends here
	req = test.NewRequest(t, "GET", "/api/user/delete", nil, apikey)
	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	println(rr.Body.String())
	// Clean up
	if rr.Code != http.StatusOK {
		t.Fatalf("Failed to delete user: %v", rr.Body.String())
	}
	// remove test database file
}