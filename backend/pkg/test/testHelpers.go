package test

import (
	"backend/pkg/db"
	"backend/pkg/handlers/editor"
	"backend/pkg/handlers/reviewer"
	"backend/pkg/handlers/user"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)


func GetTestDB() *gorm.DB {
	db, err := db.InitDBTest()
	if err != nil {
		panic(err)
	}
	return db
}

func GetTestWordHandler() *editor.WordHandler {
	return editor.NewWordHandler(GetTestDB())
}

func GetTestUserHandler() *user.UserHandler {
	return user.NewUserHandler(GetTestDB())
}

func GetTestReviewHandler() *reviewer.ReviewHandler {
	return reviewer.NewReviewHandler(GetTestDB())
}

func NewRequest(t *testing.T, method, url string, body []byte, token string) *http.Request {
	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", token)
	return req
}

// if expected is 0, then the status code will not be checked
func CreateTest(router *gin.Engine, token string, accept interface{}, 
		apipath string, method string, expected int) func(*testing.T) {
	return func(t *testing.T) {
		body, err := json.Marshal(accept)
		if err != nil {
			t.Fatalf("Failed to marshal word: %v", err)
		}
		if method == "GET" {
			body = nil
		}
		req := NewRequest(t, method, apipath, body, token)
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)
		println(rr.Body.String())
		if expected > 0 && rr.Code != expected {
			t.Errorf("handler returned wrong status code: got %v want %v",
				rr.Code, http.StatusCreated)
		}
	}
}
