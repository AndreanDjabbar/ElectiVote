package repositories

import (
	"log/slog"

	"github.com/AndreanDjabbar/ElectiVote/config"
	"github.com/AndreanDjabbar/ElectiVote/internal/db"
	"github.com/AndreanDjabbar/ElectiVote/internal/models"
)

var logger *slog.Logger = config.SetUpLogger()

func AddCandidate(newCandidate models.Candidate) (models.Candidate, error) {
	err := db.DB.Create(&newCandidate).Error
	logger.Info("Candidate Repository - Add Candidate")
	return newCandidate, err
}

func GetCandidatesByVoteID(voteID uint) ([]models.Candidate, error) {
	var candidates []models.Candidate
	err := db.DB.Where("vote_id = ?", voteID).Find(&candidates).Error
	logger.Info("Candidate Repository - Get Candidates By Vote ID")
	return candidates, err
}

func GetCandidateByCandidateID(candidateID uint) (models.Candidate, error) {
	candidate := models.Candidate{}
	err := db.DB.Where("candidate_id = ?", candidateID).Find(&candidate).Error
	logger.Info("Candidate Repository - Get Candidate By Candidate ID")
	return candidate, err
}

func GetVoteIDByCandidateID(candidateID uint) (uint, error) {
	candidate := models.Candidate{}
	err := db.DB.Where("candidate_id = ?", candidateID).Find(&candidate).Error
	logger.Info("Candidate Repository - Get Vote ID By Candidate ID")
	return candidate.VoteId, err
}

func IsValidCandidateModerator(username string, candidateID uint) bool {
	userID, err := GetUserIdByUsername(username)
	if err != nil {
		logger.Warn(
			"Candidate Repository - Error Get User ID By Username",
			"error", err,
		)
		return false
	}
	userVoteID, err := GetVoteIDByUserID(uint(userID))
	if err != nil {
		logger.Warn(
			"Candidate Repository - Error Get Vote ID By User ID",
		)
		return false
	}
	candidateVoteID, err := GetVoteIDByCandidateID(candidateID)
	if err != nil {
		logger.Warn(
			"Candidate Repository - Error Get Vote ID By Candidate ID",
		)
		return false
	}
	if userVoteID == candidateVoteID {
		return true
	}
	return false
}

func UpdateCandidate(candidateID uint, candidate models.Candidate) (models.Candidate, error) {
	err := db.DB.Model(&models.Candidate{}).Where("candidate_id = ?", candidateID).Updates(&candidate).Error
	logger.Info("Candidate Repository - Update Candidate")
	return candidate, err
}

func DeleteCandidate(candidateID uint) error {
	err := db.DB.Where("candidate_id = ?", candidateID).Delete(&models.Candidate{}).Error
	logger.Info("Candidate Repository - Delete Candidate")
	return err
}

func IncrementCandidateVote(candidateID uint) error {
	candidate := models.Candidate{}
	err := db.DB.Where("candidate_id = ?", candidateID).Find(&candidate).Error
	if err != nil {
		logger.Warn(
			"Candidate Repository - Error Increment Candidate Vote",
		)
		return err
	}
	candidate.TotalVotes += 1
	err = db.DB.Model(&models.Candidate{}).Where("candidate_id = ?", candidateID).Updates(&candidate).Error
	logger.Info("Candidate Repository - Increment Candidate Vote")
	return err
}