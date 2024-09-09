package repositories

import (
	"github.com/AndreanDjabbar/ElectiVote/internal/db"
	"github.com/AndreanDjabbar/ElectiVote/internal/models"
)

func CreateProfile(newProfile models.Profile) (models.Profile, error) {
	err := db.DB.Create(&newProfile).Error
	if err != nil {
		logger.Error(
			"Profile Repository - Error Creating Profile",
			"error", err,
		)
		return newProfile, err
	}
	return newProfile, nil
}

func UpdateProfileByUsername(username string, newProfile models.Profile) (models.Profile, error) {
	userId, err := GetUserIdByUsername(username)
	if err != nil {
		logger.Error(
			"Profile Repository - Error Get User ID By Username",
			"error", err,
		)
		return models.Profile{}, err
	}
	var profile models.Profile
	err = db.DB.Where("user_id = ?", userId).First(&profile).Error
	if err != nil {
		logger.Error(
			"Profile Repository - Error Get Profile",
			"error", err,
		)
		return models.Profile{}, err
	}
	err = db.DB.Model(&profile).Updates(newProfile).Error
	if err != nil {
		logger.Error(
			"Profile Repository - Error Updating Profile",
			"error", err,
		)
		return models.Profile{}, err
	}
	return profile, nil
}

func GetProfilesByUsername(username string) (models.Profile, error) {
	userId, err := GetUserIdByUsername(username)
	if err != nil {
		logger.Error(
			"Profile Repository - Error Get User ID By Username",
			"error", err,
		)
		return models.Profile{}, err
	}

	profile := models.Profile{}
	err = db.DB.Where("user_id = ?", userId).First(&profile).Error
	if err != nil {
		logger.Error(
			"Profile Repository - Error Get Profile",
			"error", err,
		)
		return models.Profile{}, err
	}
	return profile, nil
}