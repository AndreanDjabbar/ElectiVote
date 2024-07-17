package db

import (
	"path/filepath"

	"github.com/AndreanDjabbar/CaysAPIHub/internal/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectToDatabase() {
	absPath, err := filepath.Abs("../../internal/db/electivote.db")
	database, err := gorm.Open(
		sqlite.Open("file:" + absPath + "?cache=shared&_loc=auto"),
		&gorm.Config{},
	)
	if err != nil {
		panic(err.Error())
	}
	database.AutoMigrate(&models.AuthToken{})
	DB = database
}