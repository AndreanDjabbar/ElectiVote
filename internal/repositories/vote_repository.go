package repositories

import (
	"github.com/AndreanDjabbar/ElectiVote/internal/db"
	"github.com/AndreanDjabbar/ElectiVote/internal/models"
)

func CreateVote(vote models.Vote) (models.Vote, error) {
	err := db.DB.Create(&vote).Error
	if err != nil {
		return vote, err
	}
	return vote, nil
}

func IsUniqueCode(voteCode string) (bool) {
	votes := []models.Vote{}
	db.DB.Where("vote_code = ?", voteCode).Find(&votes)
	if len(votes) > 0 {
		return false
	}
	return true
}