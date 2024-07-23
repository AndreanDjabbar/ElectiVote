package db

import (
	"github.com/AndreanDjabbar/CaysAPIHub/internal/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectToDatabase() {

	database, err := gorm.Open(
		mysql.Open(
			"ElectiVote:ElectiVote12@tcp(localhost:3306)/ElectiVote",
		),
	)
	if err != nil {
		panic(err.Error())
	}
	database.AutoMigrate(
		&models.User{},
		&models.Profile{},
	)
	DB = database
}