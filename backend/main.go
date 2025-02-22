package main

import (
    _ "backend/docs"
    swaggerFiles "github.com/swaggo/files"
    ginSwagger "github.com/swaggo/gin-swagger"
	"backend/pkg/auth"
	"backend/pkg/db"
	"backend/pkg/handlers/editor"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)


// @title Easy Japanese API
// @version 0.1
// @description 
// @license.name MIT
// @license.url http://opensource.org/licenses/MIT
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name X-API-KEY
// @host localhost:8080
// @BasePath /
func main() {
	r := gin.Default()
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	db, err := db.InitDB()
	if err != nil {
		panic(err)
	}
	wordHandler := editor.NewWordHandler(db)
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST"},
		AllowHeaders:     []string{"Origin", "X-API-KEY", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	r.GET("/api/random", auth.APIKeyAuth(), editor.GetRandomNumber)
	r.POST("/api/answer/correct/:wordId", auth.APIKeyAuth(), editor.UpdateWordWeightCorrect)
	r.POST("/api/answer/wrong/:wordId", auth.APIKeyAuth(), editor.UpdateWordWeightIncorrect)
	r.POST("/api/words/:dictName/accurate-search", auth.APIKeyAuth(), wordHandler.AccurateSearchWord)
	r.GET("/api/words/:dictName/fuzzy-search", auth.APIKeyAuth(), wordHandler.FuzzySearchWord)
	r.POST("/api/words/:dictName/add", auth.APIKeyAuth(), wordHandler.AddWord)
	r.POST("/api/words/:dictName/edit", auth.APIKeyAuth(), wordHandler.EditWord)
	r.POST("/api/words/:dictName/delete", auth.APIKeyAuth(), wordHandler.DeleteWord)
	r.GET("/api/words/:dictName/get", auth.APIKeyAuth(), wordHandler.GetDict)
	
	r.POST("/api/reading-material/add", auth.APIKeyAuth(), wordHandler.AddReadingMaterial)
	r.POST("/api/reading-material/edit", auth.APIKeyAuth(), wordHandler.EditReadingMaterial)
	r.POST("/api/reading-material/delete", auth.APIKeyAuth(), wordHandler.DeleteReadingMaterial)
	r.GET("/api/reading-material/get", auth.APIKeyAuth(), wordHandler.GetReadingMaterial)
	r.GET("/api/reading-material/search", auth.APIKeyAuth(), wordHandler.FuzzySearchReadingMaterial)
	
	r.POST("/api/grammar/add", auth.APIKeyAuth(), wordHandler.AddGrammar)
	r.POST("/api/grammar/edit", auth.APIKeyAuth(), wordHandler.EditGrammar)
	r.POST("/api/grammar/delete", auth.APIKeyAuth(), wordHandler.DeleteGrammar)
	r.GET("/api/grammar/get", auth.APIKeyAuth(), wordHandler.GetGrammar)
	r.GET("/api/grammar/search", auth.APIKeyAuth(), wordHandler.FuzzySearchGrammar)
	r.Run(":8080")
}