// models/models.go
package models

import (
	"time"

	_ "gorm.io/gorm"
)

type SocialMedia struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `json:"name"`
	URL       string    `json:"social_media_url"`
	UserID    uint      `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	User      User      `gorm:"foreignKey:UserID" json:"-"`
}

type User struct {
	ID           uint          `gorm:"primaryKey" json:"id"`
	Username     string        `json:"username"`
	Email        string        `json:"email"`
	Password     string        `json:"password"`
	Age          int           `json:"age"`
	ImageURL     string        `json:"profile_image_url"`
	CreatedAt    time.Time     `json:"created_at"`
	UpdatedAt    time.Time     `json:"updated_at"`
	SocialMedias []SocialMedia `json:"social_medias,omitempty"`
	Photos       []Photo       `json:"photos,omitempty"`
}

type Photo struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Title     string    `json:"title"`
	Caption   string    `json:"caption"`
	URL       string    `json:"photo_url"`
	UserID    uint      `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	User      User      `gorm:"foreignKey:UserID" json:"-"`
	Comments  []Comment `json:"comments,omitempty"`
}

type Comment struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `json:"user_id"`
	PhotoID   uint      `json:"photo_id"`
	Message   string    `json:"message"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	User      User      `gorm:"foreignKey:UserID" json:"-"`
	Photo     Photo     `gorm:"foreignKey:PhotoID" json:"-"`
}
