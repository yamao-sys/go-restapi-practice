package main

import (
	"app/db"
	"app/models"

	"gorm.io/gorm"
)

func migrate(db *gorm.DB) {
	db.AutoMigrate(&models.User{}, &models.Todo{})
}

func main() {
	dbCon := db.Init()

	defer db.Close(dbCon)

	migrate(dbCon)
}
