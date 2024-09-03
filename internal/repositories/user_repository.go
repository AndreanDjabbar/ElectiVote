package repositories

import (
	"github.com/AndreanDjabbar/ElectiVote/internal/db"
	"github.com/AndreanDjabbar/ElectiVote/internal/models"
	"golang.org/x/crypto/bcrypt"
)

func RegisterUser(newUser models.User) (models.User, error) {
	logger.Info("User Repository - Register User")
	err := db.DB.Create(&newUser).Error
	if err != nil {
		logger.Error(
			"User Repository - Error Registering User",
			"error", err,
		)
		return newUser, err
	}
	logger.Info(
		"User Repository - User Registered",
	)
	return newUser, nil
}

func GetUserByUsername(username string) (models.User, error) {
	logger.Info("User Repository - Get User By Username")
	var user models.User
	err := db.DB.Where("username = ?", username).First(&user).Error
	if err != nil {
		logger.Error(
			"User Repository - Error Get User By Username",
			"error", err,
		)
		return user, err
	}
	logger.Info(
		"User Repository - User Found",
	)
	return user, nil
}

func CheckPasswordByUSername(username, password string) (bool, error) {
	logger.Info("User Repository - Check Password By Username")
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
			"User Repository - Error Compare Hash Password",
			"error", err,
		)
		return false, err
	}
	logger.Info(
		"User Repository - Password Match",
	)
	return true, nil
}

func GetUserIdByUsername(username string) (int, error) {
	logger.Info("User Repository - Get User ID By Username")
	user := models.User{}
	err := db.DB.Where("username = ?", username).First(&user).Error
	if err != nil {
		logger.Error(
			"User Repository - Error Get User ID By Username",
			"error", err,
		)
		return 0, err
	}
	logger.Info(
		"User Repository - User ID Found",
	)
	return int(user.ID), nil
}

func GetUserEmailByUsername(username string) (string, error) {
	logger.Info("User Repository - Get User Email By Username")
	user := models.User{}
	err := db.DB.Where("username = ?", username).First(&user).Error
	if err != nil {
		logger.Error(
			"User Repository - Error Get User Email By Username",
			"error", err,
		)
		return "", err
	}
	logger.Info(
		"User Repository - User Email Found",
	)
	return user.Email, nil
}

func GetUserByEmail(email string) (models.User, error) {
	logger.Info("User Repository - Get User By Email")
	var user models.User
	err := db.DB.Where("email = ?", email).First(&user).Error
	if err != nil {
		logger.Error(
			"User Repository - Error Get User By Email",
			"error", err,
		)
		return user, err
	}
	logger.Info(
		"User Repository - User Found",
	)
	return user, nil
}

func UpdatePasswordByEmail(email, passwordHashed string) (models.User, error) {
	logger.Info("User Repository - Update Password By Email")
	var user models.User
	err := db.DB.Model(&user).Where("email = ?", email).Update("password", passwordHashed).Error
	if err != nil {
		logger.Error(
			"User Repository - Error Update Password By Email",
			"error", err,
		)
		return user, err
	}
	logger.Info(
		"User Repository - Password Updated",
	)
	return user, nil
}