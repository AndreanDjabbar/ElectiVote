package repositories

import (
	"github.com/AndreanDjabbar/ElectiVote/internal/db"
	"github.com/AndreanDjabbar/ElectiVote/internal/models"
)

func CreateProfile(newProfile models.Profile) (models.Profile, error) {
	logger.Info("Profile Repository - Create Profile")
	err := db.DB.Create(&newProfile).Error
	if err != nil {
		logger.Error(
			"Profile Repository - Error Creating Profile",
		)
		return newProfile, err
	}
	logger.Info(
		"Profile Repository - Profile Created",
	)
	return newProfile, nil
}

func UpdateProfileByUsername(username string, newProfile models.Profile) (models.Profile, error) {
	logger.Info("Profile Repository - Update Profile By Username")
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
			"Profile Repository - Error Get Profile By User ID",
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
	logger.Info(
		"Profile Repository - Profile Updated",
	)
	return profile, nil
}

func GetProfilesByUsername(username string) (models.Profile, error) {
	logger.Info("Profile Repository - Get Profile By Username")
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
			"Profile Repository - Error Get Profile By User ID",
		)
		return models.Profile{}, err
	}
	logger.Info(
		"Profile Repository - Profile Retrieved",
	)
	return profile, nil
}