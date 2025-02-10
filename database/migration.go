package database

import (
	"article-golang/models"
	"fmt"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) {
	err := db.AutoMigrate(&models.Post{})
	if err != nil {
		fmt.Println("Migration failed:", err)
	} else {
		fmt.Println("Migration successful")
	}
}