package repositories

import (
	"github.com/AndreanDjabbar/CaysAPIHub/internal/db"
	"github.com/AndreanDjabbar/CaysAPIHub/internal/models"
)

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