package test

import (
	"backend/pkg/db"
	"backend/pkg/handlers/editor"
	"backend/pkg/handlers/user"

	"gorm.io/gorm"
)


func GetTestDB() *gorm.DB {
	db, err := db.InitDB()
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