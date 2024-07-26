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

func GetUserIdByUsername(username string) (int, error) {
	user := models.User{}
	err := db.DB.Where("username = ?", username).First(&user).Error
	if err != nil {
		return 0, err
	}
	return int(user.ID), nil
}

func GetUserEmailByUsername(username string) (string, error) {
	user := models.User{}
	err := db.DB.Where("username = ?", username).First(&user).Error
	if err != nil {
		return "", err
	}
	return user.Email, nil
}

func CreateProfile(newProfile models.Profile) (models.Profile, error) {
	err := db.DB.Create(&newProfile).Error
	if err != nil {
		return newProfile, err
	}
	return newProfile, nil
}

func UpdateProfileByUsername(username string, newProfile models.Profile) (models.Profile, error) {
	userId, err := GetUserIdByUsername(username)
	if err != nil {
		return models.Profile{}, err
	}
	var profile models.Profile
	err = db.DB.Where("user_id = ?", userId).First(&profile).Error
	if err != nil {
		return models.Profile{}, err
	}
	err = db.DB.Model(&profile).Updates(newProfile).Error
	if err != nil {
		return models.Profile{}, err
	}
	return profile, nil
}

func GetProfilesByUsername(username string) (models.Profile, error) {
    userId, err := GetUserIdByUsername(username)
    if err != nil {
        return models.Profile{}, err
    }

    profile := models.Profile{}
    err = db.DB.Where("user_id = ?", userId).First(&profile).Error
    if err != nil {
        return models.Profile{}, err
    }
    return profile, nil
}