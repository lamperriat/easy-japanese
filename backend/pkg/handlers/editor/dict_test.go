package editor_test

import (
	"backend/pkg/test"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)


func TestGetDict(t *testing.T) {
	router := gin.Default()
	wordHandler := test.GetTestWordHandler()
	router.GET("/api/dict/:dictName/get", wordHandler.GetDict)

	req, err := http.NewRequest("GET", "/api/dict/all/get", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	println(rr.Body.String())
	if rr.Code != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			rr.Code, http.StatusOK)
	}
}