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

func GetVotesDataByUsername(username string) ([]models.Vote, error) {
	userID, err := GetUserIdByUsername(username)
	if err != nil {
		return nil, err
	}
	votes := []models.Vote{}
	err = db.DB.Where("moderator_id = ?", uint(userID)).Find(&votes).Error
	if err != nil {
		return votes, err
	}
	return votes, nil
}

func GetVoteIDByUserID(userID uint) (uint, error) {
	vote := models.Vote{}
	err := db.DB.Where("moderator_id = ?", userID).Find(&vote).Error
	if err != nil {
		return 0, err
	}
	return vote.VoteID, nil
}

func GetVoteDataByVoteID(voteID uint) (models.Vote, error) {
	vote := models.Vote{}
	err := db.DB.Where("vote_id = ?", voteID).Find(&vote).Error
	if err != nil {
		return vote, err
	}
	return vote, nil
}

func GetModeratorIDByVoteID(voteID uint) (uint, error) {
	vote := models.Vote{}
	err := db.DB.Where("vote_id = ?", voteID).Find(&vote).Error
	if err != nil {
		return 0, err
	}
	return vote.ModeratorID, nil
}

func IsValidVoteModerator(username string, voteID uint) bool {
	userID, err := GetUserIdByUsername(username)
	if err != nil {
		return false
	}
	moderatorID, err := GetModeratorIDByVoteID(voteID)
	if err != nil {
		return false
	}
	if uint(userID) == moderatorID {
		return true
	}
	return false
}

func UpdateVote(voteID uint, vote models.Vote) (models.Vote, error) {
	err := db.DB.Model(&vote).Where("vote_id = ?", voteID).Updates(vote).Error
	if err != nil {
		return vote, err
	}
	return vote, nil	
}

func DeleteVote(voteID uint) error {
	vote := models.Vote{}
	err := db.DB.Where("vote_id = ?", voteID).Delete(&vote).Error
	if err != nil {
		return err
	}
	return nil
}

func GetVoteByVoteCode(voteCode string) (models.Vote, error) {
	vote := models.Vote{}
	err := db.DB.Where("BINARY vote_code = ?", voteCode).First(&vote).Error
	if err != nil {
		return vote, err
	}
	return vote, nil
}

func GetVoteIDByVoteCode(voteCode string) (uint, error) {
	vote := models.Vote{}
	err := db.DB.Where("BINARY vote_code = ?", voteCode).First(&vote).Error
	if err != nil {
		return 0, err
	}
	return vote.VoteID, nil
}