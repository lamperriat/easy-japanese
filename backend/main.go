package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"backend/pkg/handlers"
	"backend/pkg/auth"
)

func main() {
	r := gin.Default()

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
	r.POST("/api/words/:dictName/check", auth.APIKeyAuth(), handlers.CheckSimilarWords)
	r.POST("/api/words/:dictName/submit", auth.APIKeyAuth(), handlers.AddWord)	
	r.Run(":8080")
}