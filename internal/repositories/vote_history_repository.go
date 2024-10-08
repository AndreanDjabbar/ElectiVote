package repositories

import (
	"github.com/AndreanDjabbar/ElectiVote/internal/db"
	"github.com/AndreanDjabbar/ElectiVote/internal/models"
)

func GetVoteHistoriesByUserID(userID uint) ([]models.VoteHistory, error) {
	voteHistories := []models.VoteHistory{}
	err := db.DB.Where("moderator_id = ?", userID).Find(&voteHistories).Error
	if err != nil {
		return voteHistories, err
	}
	return voteHistories, nil
}

func CreateVoteHistory(voteHistory *models.VoteHistory) error {
	err := db.DB.Create(voteHistory).Error
	if err != nil {
		return err
	}
	return nil
}

func GetVoteHistoryByVoteHistoryID(voteHistoryID uint) (models.VoteHistory, error) {
	voteHistory := models.VoteHistory{}
	err := db.DB.Where("vote_history_id = ?", voteHistoryID).First(&voteHistory).Error
	if err != nil {
		return voteHistory, err
	}
	return voteHistory, nil
}