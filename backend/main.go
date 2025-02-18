package main

import (
	"backend/pkg/auth"
	"backend/pkg/db"
	"backend/pkg/handlers"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	db, err := db.InitDB()
	if err != nil {
		panic(err)
	}
	wordHandler := handlers.NewWordHandler(db)
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST"},
		AllowHeaders:     []string{"Origin", "X-API-KEY", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	r.GET("/api/random", auth.APIKeyAuth(), handlers.GetRandomNumber)
	r.POST("/api/answer/correct/:wordId", auth.APIKeyAuth(), handlers.UpdateWordWeightCorrect)
	r.POST("/api/answer/wrong/:wordId", auth.APIKeyAuth(), handlers.UpdateWordWeightIncorrect)
	r.POST("/api/words/:dictName/check", auth.APIKeyAuth(), wordHandler.CheckSimilarWords)
	r.POST("/api/words/:dictName/submit", auth.APIKeyAuth(), wordHandler.AddWord)
	r.GET("/api/dict/:dictName/get", auth.APIKeyAuth(), handlers.GetDict)	
	r.Run(":8080")
}