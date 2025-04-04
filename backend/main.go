package main

import (
	_ "backend/docs"
	"backend/pkg/auth"
	"backend/pkg/db"
	"backend/pkg/handlers/authmid"
	"backend/pkg/handlers/editor"
	"backend/pkg/handlers/reviewer"
	"backend/pkg/handlers/user"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

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
	userHandler := user.NewUserHandler(db)
	reviewHandler := reviewer.NewReviewHandler(db)
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST"},
		AllowHeaders:     []string{"Origin", "X-API-KEY", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// r.GET("/api/random", auth.APIKeyAuth(), editor.GetRandomNumber)
	r.POST("/api/auth/token", authmid.GetToken)
	// User routes
	userGroup := r.Group("/api/user", auth.JWTAuth())
	{
		userGroup.POST("/register", userHandler.RegisterUser)
		userGroup.POST("/update", userHandler.UpdateUserName)
		userGroup.GET("/delete", userHandler.DeleteUser)
		grammarGroup := userGroup.Group("/grammar")
		{
			grammarGroup.POST("/add", wordHandler.AddGrammarUser)
			grammarGroup.POST("/edit", wordHandler.EditGrammarUser)
			grammarGroup.POST("/delete", wordHandler.DeleteGrammarUser)
			grammarGroup.GET("/get", wordHandler.GetGrammarUser)
			grammarGroup.GET("/search", wordHandler.SearchGrammarUser)
		}
		readingGroup := userGroup.Group("/reading-material")
		{
			readingGroup.POST("/add", wordHandler.AddReadingMaterialUser)
			readingGroup.POST("/edit", wordHandler.EditReadingMaterialUser)
			readingGroup.POST("/delete", wordHandler.DeleteReadingMaterialUser)
			readingGroup.GET("/get", wordHandler.GetReadingMaterialUser)
			readingGroup.GET("/search", wordHandler.FuzzySearchReadingMaterialUser) // full text search
		}
		wordsGroup := userGroup.Group("/words")
		{
			wordsGroup.POST("/accurate-search", wordHandler.AccurateSearchWordUser)
			wordsGroup.GET("/fuzzy-search", wordHandler.FuzzySearchWordUser)
			wordsGroup.POST("/add", wordHandler.AddWordUser)
			wordsGroup.POST("/edit", wordHandler.EditWordUser)
			wordsGroup.POST("/delete", wordHandler.DeleteWordUser)
			wordsGroup.GET("/get", wordHandler.GetWordsUser)
		}
		reviewGroup := userGroup.Group("/review")
		{
			reviewWordGroup := reviewGroup.Group("/word")
			{
				reviewWordGroup.POST("/correct", reviewHandler.CorrectWord)
				reviewWordGroup.POST("/incorrect", reviewHandler.IncorrectWord)
				reviewWordGroup.GET("/get", reviewHandler.GetWords)
			}
			reviewGrammarGroup := reviewGroup.Group("/grammar")
			{
				reviewGrammarGroup.POST("/correct", reviewHandler.CorrectGrammar)
				reviewGrammarGroup.POST("/incorrect", reviewHandler.IncorrectGrammar)
				reviewGrammarGroup.GET("/get", reviewHandler.GetGrammar)
			}
		}
	}

	// Answer routes
	answerGroup := r.Group("/api/answer", auth.JWTAuth())
	{
		answerGroup.POST("/correct/:wordId", editor.UpdateWordWeightCorrect)
		answerGroup.POST("/wrong/:wordId", editor.UpdateWordWeightIncorrect)
	}

	// Dictionary/Words routes
	dictGroup := r.Group("/api/words/:dictName", auth.JWTAuth())
	{
		dictGroup.POST("/accurate-search", wordHandler.AccurateSearchWord)
		dictGroup.GET("/fuzzy-search", wordHandler.FuzzySearchWord)
		dictGroup.POST("/add", wordHandler.AddWord)
		dictGroup.POST("/edit", wordHandler.EditWord)
		dictGroup.POST("/delete", wordHandler.DeleteWord)
		dictGroup.GET("/get", wordHandler.GetDict)
	}
	
	// Reading Material routes
	readingGroup := r.Group("/api/reading-material", auth.JWTAuth())
	{
		readingGroup.POST("/add", wordHandler.AddReadingMaterial)
		readingGroup.POST("/edit", wordHandler.EditReadingMaterial)
		readingGroup.POST("/delete", wordHandler.DeleteReadingMaterial)
		readingGroup.GET("/get", wordHandler.GetReadingMaterial)
		readingGroup.GET("/search", wordHandler.FuzzySearchReadingMaterial)
	}
	
	// Grammar routes
	grammarGroup := r.Group("/api/grammar", auth.JWTAuth())
	{
		grammarGroup.POST("/add", wordHandler.AddGrammar)
		grammarGroup.POST("/edit", wordHandler.EditGrammar)
		grammarGroup.POST("/delete", wordHandler.DeleteGrammar)
		grammarGroup.GET("/get", wordHandler.GetGrammar)
		grammarGroup.GET("/search", wordHandler.FuzzySearchGrammar)
	}

	
	r.Run(":8080")
}