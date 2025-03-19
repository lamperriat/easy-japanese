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

func TestUserReadingOps(t *testing.T) {
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
	router.POST("/api/user/reading-material/add", auth.APIKeyAuth(), handler.AddReadingMaterialUser)
	router.POST("/api/user/reading-material/edit", auth.APIKeyAuth(), handler.EditReadingMaterialUser)
	router.POST("/api/user/reading-material/delete", auth.APIKeyAuth(), handler.DeleteReadingMaterialUser)
	router.GET("/api/user/reading-material/get", auth.APIKeyAuth(), handler.GetReadingMaterialUser)
	router.GET("/api/user/reading-material/search", auth.APIKeyAuth(), handler.FuzzySearchReadingMaterialUser)
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
	reading1 := models.UserReadingMaterial{
		Title:     "title1",
		Content:   "test_content1",
		Chinese:   "test_chinese1",
	}
	reading2 := models.UserReadingMaterial{
		Title:     "title2",
		Content:   "test_content2",
		Chinese:   "test_chinese2",
	}
	reading3 := models.UserReadingMaterial{
		Title:     "title3",
		Content:   "test_content3",
		Chinese:   "test_chinese3",
	}
	// Test starts here
	t.Run("Add reading 1", test.CreateTest(
		router, apikey, reading1, "/api/user/reading-material/add", "POST", http.StatusCreated,
	))
	t.Run("Add reading 2", test.CreateTest(
		router, apikey, reading2, "/api/user/reading-material/add", "POST", http.StatusCreated,
	))
	t.Run("Add reading 3", test.CreateTest(
		router, apikey, reading3, "/api/user/reading-material/add", "POST", http.StatusCreated,
	))
	t.Run("Get reading", test.CreateTest(
		router, apikey, nil, "/api/user/reading-material/get", "GET", http.StatusOK,
	))
	t.Run("Search reading", test.CreateTest(
		router, apikey, nil, "/api/user/reading-material/search?query=content2", "GET", http.StatusOK,
	))
	t.Run("Edit reading", test.CreateTest(
		router, apikey, models.UserReadingMaterial{
			ID:       1,
			Title:    "title1_updated",
			Content:  "test_content1_updated",
			Chinese:  "test_chinese1_updated",
		}, "/api/user/reading-material/edit", "POST", http.StatusOK,
	))
	t.Run("Search reading", test.CreateTest(
		router, apikey, nil, "/api/user/reading-material/search?query=test", "GET", http.StatusOK,
	))
	t.Run("Delete reading", test.CreateTest(
		router, apikey, models.UserReadingMaterial{
			ID: 1,
		}, "/api/user/reading-material/delete", "POST", http.StatusOK,
	))
	t.Run("Get reading", test.CreateTest(
		router, apikey, nil, "/api/user/reading-material/get", "GET", http.StatusOK,
	))
	t.Run("Search reading", test.CreateTest(
		router, apikey, nil, "/api/user/reading-material/search?query=test", "GET", http.StatusOK,
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