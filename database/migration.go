package database

import (
	"github.com/faazabilamri7/mygram/models"
)

func AutoMigrate() {
	db := GetDB()
	db.AutoMigrate(&models.User{})
	db.AutoMigrate(&models.SocialMedia{})
	db.AutoMigrate(&models.Photo{})
	db.AutoMigrate(&models.Comment{})
}
