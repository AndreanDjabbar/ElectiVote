package repositories

import (
	"github.com/AndreanDjabbar/CaysAPIHub/internal/db"
	"github.com/AndreanDjabbar/CaysAPIHub/internal/models"
)

func RegisterUser(newUser models.User) (models.User, error) {
	err := db.DB.Create(&newUser).Error
	if err != nil {
		return newUser, err
	}
	return newUser, nil
}