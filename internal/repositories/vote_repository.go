package repositories

import (
	"github.com/AndreanDjabbar/ElectiVote/internal/db"
	"github.com/AndreanDjabbar/ElectiVote/internal/models"
)

func CreateVote(vote models.Vote) (models.Vote, error) {
	logger.Info("Vote Repository - Create Vote")
	err := db.DB.Create(&vote).Error
	if err != nil {
		logger.Error(
			"Vote Repository - Error Creating Vote",
			"error", err,
		)
		return vote, err
	}
	logger.Info(
		"Vote Repository - Vote Created",
	)
	return vote, nil
}

func IsUniqueCode(voteCode string) (bool) {
	logger.Info("Vote Repository - Is Unique Code")
	votes := []models.Vote{}
	db.DB.Where("vote_code = ?", voteCode).Find(&votes)
	if len(votes) > 0 {
		logger.Warn(
			"Vote Repository - Vote Code Already Exists",
		)
		return false
	}
	logger.Info(
		"Vote Repository - Vote Code is Unique",
	)
	return true
}

func GetVotesDataByUsername(username string) ([]models.Vote, error) {
	logger.Info("Vote Repository - Get Votes Data By Username")
	userID, err := GetUserIdByUsername(username)
	if err != nil {
		logger.Error(
			"Vote Repository - Error Get User ID By Username",
			"error", err,
		)
		return nil, err
	}
	votes := []models.Vote{}
	err = db.DB.Where("moderator_id = ?", uint(userID)).Find(&votes).Error
	if err != nil {
		logger.Error(
			"Vote Repository - Error Get Votes Data By Username",
			"error", err,
		)
		return votes, err
	}
	logger.Info(
		"Vote Repository - Votes Data Found",
	)
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
	logger.Info("Vote Repository - Get Vote Data By Vote ID")
	vote := models.Vote{}
	err := db.DB.Where("vote_id = ?", voteID).Find(&vote).Error
	if err != nil {
		logger.Error(
			"Vote Repository - Error Get Vote Data By Vote ID",
			"error", err,
		)
		return vote, err
	}
	logger.Info(
		"Vote Repository - Vote Data Found",
	)
	return vote, nil
}

func GetModeratorIDByVoteID(voteID uint) (uint, error) {
	logger.Info("Vote Repository - Get Moderator ID By Vote ID")
	vote := models.Vote{}
	err := db.DB.Where("vote_id = ?", voteID).Find(&vote).Error
	if err != nil {
		logger.Error(
			"Vote Repository - Error Get Moderator ID By Vote ID",
			"error", err,
		)
		return 0, err
	}
	logger.Info(
		"Vote Repository - Moderator ID Found",
	)
	return vote.ModeratorID, nil
}

func IsValidVoteModerator(username string, voteID uint) bool {
	logger.Info("Vote Repository - Is Valid Vote Moderator")
	userID, err := GetUserIdByUsername(username)
	if err != nil {
		logger.Warn(
			"Vote Repository - Error Get User ID By Username",
		)
		return false
	}
	moderatorID, err := GetModeratorIDByVoteID(voteID)
	if err != nil {
		logger.Warn(
			"Vote Repository - Error Get Moderator ID By Vote ID",
		)
		return false
	}
	if uint(userID) == moderatorID {
		logger.Info(
			"Vote Repository - User is Moderator",
		)
		return true
	}	
	logger.Warn(
		"Vote Repository - User is Not Moderator",
	)
	return false
}

func UpdateVote(voteID uint, vote models.Vote) (models.Vote, error) {
	logger.Info("Vote Repository - Update Vote")
	err := db.DB.Model(&vote).Where("vote_id = ?", voteID).Updates(vote).Error
	if err != nil {
		logger.Error(
			"Vote Repository - Error Updating Vote",
			"error", err,
		)
		return vote, err
	}
	logger.Info(
		"Vote Repository - Vote Updated",
	)
	return vote, nil	
}

func DeleteVote(voteID uint) error {
	logger.Info("Vote Repository - Delete Vote")
	vote := models.Vote{}
	err := db.DB.Where("vote_id = ?", voteID).Delete(&vote).Error
	if err != nil {
		logger.Error(
			"Vote Repository - Error Deleting Vote",
			"error", err,
		)
		return err
	}
	logger.Info(
		"Vote Repository - Vote Deleted",
	)
	return nil
}

func GetVoteByVoteCode(voteCode string) (models.Vote, error) {
	logger.Info("Vote Repository - Get Vote By Vote Code")
	vote := models.Vote{}
	err := db.DB.Where("BINARY vote_code = ?", voteCode).First(&vote).Error
	if err != nil {
		logger.Error(
			"Vote Repository - Error Get Vote By Vote Code",
			"error", err,
		)
		return vote, err
	}
	logger.Info(
		"Vote Repository - Vote Found",
	)
	return vote, nil
}

func GetVoteIDByVoteCode(voteCode string) (uint, error) {
	logger.Info("Vote Repository - Get Vote ID By Vote Code")
	vote := models.Vote{}
	err := db.DB.Where("BINARY vote_code = ?", voteCode).First(&vote).Error
	if err != nil {
		logger.Error(
			"Vote Repository - Error Get Vote ID By Vote Code",
			"error", err,
		)
		return 0, err
	}
	logger.Info(
		"Vote Repository - Vote ID Found",
	)
	return vote.VoteID, nil
}