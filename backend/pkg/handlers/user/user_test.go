package user_test

import (
	"backend/pkg/auth"
	// "backend/pkg/handlers/authmid"
	"backend/pkg/models"
	"backend/pkg/test"
	"bytes"
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
}

func TestUserOps(t *testing.T) {
	router := gin.Default()
	handler := test.GetTestUserHandler()

	router.POST("/api/user/register", auth.JWTAuth(), handler.RegisterUser)
	router.POST("/api/user/update", auth.JWTAuth(), handler.UpdateUserName)
	router.GET("/api/user/delete", auth.JWTAuth(), handler.DeleteUser)
	// router.POST("/api/auth/token", authmid.GetToken(test.GetTestDB()))
	apikey := "TEST_USE_API_KEY"
	var token string
	t.Run("GetToken", func(t *testing.T) {
		req, err := http.NewRequest("POST", "/api/auth/token", nil)
		req.Header.Set("X-API-KEY", apikey)
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)
		println(rr.Body.String())
		if rr.Code != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v",
				rr.Code, http.StatusOK)
		}
		var tokenResponse models.TokenResponse
		err = json.Unmarshal(rr.Body.Bytes(), &tokenResponse)
		if err != nil {
			t.Fatalf("Failed to unmarshal token response: %v", err)
		}
		if tokenResponse.Token == "" {
			t.Errorf("handler returned empty token")
		}
		if tokenResponse.ExpiresIn == 0 {
			t.Errorf("handler returned empty expiresIn")
		}
		token = tokenResponse.Token
		println("Token: ", token)
	})
	t.Run("RegisterUser", func(t *testing.T) {
		toRegister := models.User{
			Username: "test_user",
		}
	
		body, err := json.Marshal(toRegister)
		if err != nil {
			t.Fatalf("Failed to marshal user: %v", err)
		}
	
		req, err := http.NewRequest("POST", "/api/user/register", bytes.NewBuffer(body))
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", token)
	
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)
		println(rr.Body.String())
		if rr.Code != http.StatusCreated {
			t.Errorf("handler returned wrong status code: got %v want %v",
				rr.Code, http.StatusCreated)
		}	
	})
	
	t.Run("UpdateUserName", func(t *testing.T) {
		toUpdate := models.User{
			Username: "test_user_updated",
		}

		body, err := json.Marshal(toUpdate)
		if err != nil {
			t.Fatalf("Failed to marshal user: %v", err)
		}

		req, err := http.NewRequest("POST", "/api/user/update", bytes.NewBuffer(body))
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", token)

		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)
		println(rr.Body.String())
		if rr.Code != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v",
				rr.Code, http.StatusOK)
		}
	})

	t.Run("RegisterUserFail", func(t *testing.T) {
		toRegister := models.User{
			Username: "test_user",
		}
	
		body, err := json.Marshal(toRegister)
		if err != nil {
			t.Fatalf("Failed to marshal user: %v", err)
		}
	
		req, err := http.NewRequest("POST", "/api/user/register", bytes.NewBuffer(body))
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}
		req.Header.Set("Content-Type", "application/json")	
		req.Header.Set("Authorization", token)
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)
		println(rr.Body.String())
		if rr.Code != http.StatusConflict {
			t.Errorf("handler returned wrong status code: got %v want %v",
				rr.Code, http.StatusCreated)
		}	
	})

	// t.Run("RegisterFailDueToDuplicateUsername", func(t *testing.T) {
	// 	toRegister := models.User{
	// 		Username: "test_user_updated",
	// 	}

	// 	new_apikey := "TEST_USE_API_KEY2"
	// 	body, err := json.Marshal(toRegister)
	// 	if err != nil {
	// 		t.Fatalf("Failed to marshal user: %v", err)
	// 	}

	// 	req, err := http.NewRequest("POST", "/api/user/register", bytes.NewBuffer(body))
	// 	if err != nil {
	// 		t.Fatalf("Failed to create request: %v", err)
	// 	}
	// 	req.Header.Set("Content-Type", "application/json")
	// 	req.Header.Set("X-API-KEY", new_apikey)

	// 	rr := httptest.NewRecorder()
	// 	router.ServeHTTP(rr, req)
	// 	println(rr.Body.String())
	// 	if rr.Code != http.StatusConflict {
	// 		t.Errorf("handler returned wrong status code: got %v want %v",
	// 			rr.Code, http.StatusConflict)
	// 	}
	// })

	t.Run("DeleteUser", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/api/user/delete", nil)
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}
		req.Header.Set("Authorization", token)

		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)
		println(rr.Body.String())
		if rr.Code != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v",
				rr.Code, http.StatusOK)
		}
	})
}


