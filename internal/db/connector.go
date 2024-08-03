package db

import (
	"fmt"
	"log"
	"os"

	"github.com/AndreanDjabbar/ElectiVote/internal/models"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectToDatabase() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	database, err := gorm.Open(		
		mysql.Open(dsn),
		&gorm.Config{},
	)
	if err != nil {
		panic(err.Error())
	}
	database.AutoMigrate(
		&models.Vote{},
		&models.Candidate{},
		&models.User{},
		&models.Profile{},
	)
	DB = database
}