package repositories

import (
	"github.com/AndreanDjabbar/ElectiVote/internal/db"
	"github.com/AndreanDjabbar/ElectiVote/internal/models"
	"golang.org/x/crypto/bcrypt"
)

func RegisterUser(newUser models.User) (models.User, error) {
	err := db.DB.Create(&newUser).Error
	if err != nil {
		logger.Error(
			"User Repository - Error Register User",
			"error", err,
		)
		return newUser, err
	}
	return newUser, nil
}

func GetUserByUsername(username string) (models.User, error) {
	var user models.User
	err := db.DB.Where("username = ?", username).First(&user).Error
	if err != nil {
		logger.Error(
			"User Repository - Error Get User By Username",
			"error", err,
		)
		return user, err
	}
	return user, nil
}

func CheckPasswordByUSername(username, password string) (bool, error) {
	user, err := GetUserByUsername(username)
	if err != nil {
		logger.Error(
			"User Repository - Error Get User By Username",
			"error", err,
		)
		return false, err
	}
	err = bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(password),
	)
	if err != nil {
		logger.Error(
			"User Repository - Error Compare Password",
			"error", err,
		)
		return false, err
	}
	return true, nil
}

func GetUserIdByUsername(username string) (int, error) {
	user := models.User{}
	err := db.DB.Where("username = ?", username).First(&user).Error
	if err != nil {
		logger.Error(
			"User Repository - Error Get User ID By Username",
			"error", err,
		)
		return 0, err
	}
	return int(user.ID), nil
}

func GetUserEmailByUsername(username string) (string, error) {
	user := models.User{}
	err := db.DB.Where("username = ?", username).First(&user).Error
	if err != nil {
		logger.Error(
			"User Repository - Error Get User Email By Username",
			"error", err,
		)
		return "", err
	}
	return user.Email, nil
}

func GetUserByEmail(email string) (models.User, error) {
	var user models.User
	err := db.DB.Where("email = ?", email).First(&user).Error
	if err != nil {
		logger.Error(
			"User Repository - Error Get User By Email",
			"error", err,
		)
		return user, err
	}
	return user, nil
}

func UpdatePasswordByEmail(email, passwordHashed string) (models.User, error) {
	var user models.User
	err := db.DB.Model(&user).Where("email = ?", email).Update("password", passwordHashed).Error
	if err != nil {
		logger.Error(
			"User Repository - Error Update Password By Email",
			"error", err,
		)
		return user, err
	}
	return user, nil
}