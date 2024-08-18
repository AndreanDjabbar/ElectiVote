package repositories

import (
	"github.com/AndreanDjabbar/ElectiVote/internal/db"
	"github.com/AndreanDjabbar/ElectiVote/internal/models"
)

func CreateVoteRecord(voteRecord models.VoteRecord) (models.VoteRecord, error) {
	err := db.DB.Create(&voteRecord).Error
	if err != nil {
		return voteRecord, err
	}
	return voteRecord, nil
}

func IsVoted(userID uint, voteCode string) (bool) {
	voteID, err := GetVoteIDByVoteCode(voteCode)
	if err != nil {
		return false
	}
	voteRecord := models.VoteRecord{}
	err = db.DB.Where("user_id = ? AND vote_id = ?", userID, voteID).First(&voteRecord).Error
	if err != nil {
		return false
	}
	return true
}