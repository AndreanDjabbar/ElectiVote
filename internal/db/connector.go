package db

import (
	"fmt"
	"os"

	"github.com/AndreanDjabbar/ElectiVote/config"
	"github.com/AndreanDjabbar/ElectiVote/internal/models"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectToDatabase() {
	logger := config.SetUpLogger()
	err := godotenv.Load()
	if err != nil {
		logger.Error("Error loading .env file")
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
		logger.Error("Error connecting to database", "error", err)
		panic(err.Error())
	}
	err = database.AutoMigrate(
		&models.User{},
		&models.Profile{},
		&models.Vote{},
		&models.Candidate{},
		&models.VoteRecord{},
		&models.VoteHistory{},
	)
	if err != nil {
		logger.Error("Error migrating database", "error", err)
		panic(err.Error())
	}
	
	DB = database
	logger.Info("Database connected")
}