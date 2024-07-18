package repositories

import (
	"github.com/AndreanDjabbar/CaysAPIHub/internal/db"
	"github.com/AndreanDjabbar/CaysAPIHub/internal/models"
	"golang.org/x/crypto/bcrypt"
)

func RegisterUser(newUser models.User) (models.User, error) {
	err := db.DB.Create(&newUser).Error
	if err != nil {
		return newUser, err
	}
	return newUser, nil
}

func GetUserByUsername(username string) (models.User, error) {
	var user models.User
	err := db.DB.Where("username = ?", username).First(&user).Error
	if err != nil {
		return user, err
	}
	return user, nil
}

func CheckPasswordByUSername(username, password string) (bool, error) {
	user, err := GetUserByUsername(username)
	if err != nil {
		return false, err
	}
	err = bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(password),
	)
	if err != nil {
		return false, err
	}
	return true, nil
}