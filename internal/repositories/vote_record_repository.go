package repositories

import (
	"github.com/AndreanDjabbar/ElectiVote/internal/db"
	"github.com/AndreanDjabbar/ElectiVote/internal/models"
)

func CreateVoteRecord(voteRecord models.VoteRecord) (models.VoteRecord, error) {
	logger.Info("Vote Record Repository - Create Vote Record")
	err := db.DB.Create(&voteRecord).Error
	if err != nil {
		logger.Error(
			"Vote Record Repository - Error Creating Vote Record",
			"error", err,
		)
		return voteRecord, err
	}
	logger.Info(
		"Vote Record Repository - Vote Record Created",
	)
	return voteRecord, nil
}

func IsVoted(userID uint, voteCode string) (bool) {
	logger.Info("Vote Record Repository - Is Voted")
	voteID, err := GetVoteIDByVoteCode(voteCode)
	if err != nil {
		logger.Error(
			"Vote Record Repository - Error Get Vote ID By Vote Code",
			"error", err,
		)
		return false
	}
	voteRecord := models.VoteRecord{}
	err = db.DB.Where("user_id = ? AND vote_id = ?", userID, voteID).First(&voteRecord).Error
	if err != nil {
		logger.Error(
			"Vote Record Repository - Error Get Vote Record By User ID And Vote ID",
			"error", err,
		)
		return false
	}
	logger.Info(
		"Vote Record Repository - Vote Record Found",
	)
	return true
}