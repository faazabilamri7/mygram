// database/database.go
package database

import (
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
)

var db *gorm.DB

func InitDB() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	host := os.Getenv("DATABASE_HOST")
	user := os.Getenv("DATABASE_USER")
	password := os.Getenv("DATABASE_PASSWORD")
	dbName := os.Getenv("DATABASE_NAME")

	// Periksa apakah semua variabel lingkungan telah diatur
	if host == "" || user == "" || password == "" || dbName == "" {
		panic("One or more environment variables are not set")
	}

	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local", user, password, host, dbName)

	conn, err := gorm.Open("mysql", dataSourceName)
	if err != nil {
		panic(err.Error())
	}

	// Set max open connections
	conn.DB().SetMaxOpenConns(20)

	// Simpan koneksi database
	db = conn

	// Ping database to check connectivity
	err = conn.DB().Ping()
	if err != nil {
		panic(err.Error())
	}

	// Auto Migrate tabel
}

func GetDB() *gorm.DB {
	return db
}
