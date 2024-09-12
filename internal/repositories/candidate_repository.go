package repositories

import (
	"github.com/AndreanDjabbar/ElectiVote/internal/db"
	"github.com/AndreanDjabbar/ElectiVote/internal/models"
)

func AddCandidate(newCandidate models.Candidate) (models.Candidate, error) {
	err := db.DB.Create(&newCandidate).Error
	return newCandidate, err
}

func GetCandidatesByVoteID(voteID uint) ([]models.Candidate, error) {
	var candidates []models.Candidate
	err := db.DB.Where("vote_id = ?", voteID).Find(&candidates).Error
	return candidates, err
}

func GetCandidateByCandidateID(candidateID uint) (models.Candidate, error) {
	candidate := models.Candidate{}
	err := db.DB.Where("candidate_id = ?", candidateID).Find(&candidate).Error
	return candidate, err
}

func GetVoteIDByCandidateID(candidateID uint) (uint, error) {
	candidate := models.Candidate{}
	err := db.DB.Where("candidate_id = ?", candidateID).Find(&candidate).Error
	return candidate.VoteId, err
}

func IsValidCandidateModerator(username string, candidateID uint) bool {
	userID, err := GetUserIdByUsername(username)
	if err != nil {
		return false
	}
	userVoteID, err := GetVoteIDByUserID(uint(userID))
	if err != nil {
		return false
	}
	candidateVoteID, err := GetVoteIDByCandidateID(candidateID)
	if err != nil {
		return false
	}
	if userVoteID == candidateVoteID {
		return true
	}
	return false
}

func UpdateCandidate(candidateID uint, candidate models.Candidate) (models.Candidate, error) {
	err := db.DB.Model(&models.Candidate{}).Where("candidate_id = ?", candidateID).Updates(&candidate).Error
	return candidate, err
}

func DeleteCandidate(candidateID uint) error {
	err := db.DB.Where("candidate_id = ?", candidateID).Delete(&models.Candidate{}).Error
	return err
}

func IncrementCandidateVote(candidateID uint) error {
	candidate := models.Candidate{}
	err := db.DB.Where("candidate_id = ?", candidateID).Find(&candidate).Error
	if err != nil {
		return err
	}
	candidate.TotalVotes += 1
	err = db.DB.Model(&models.Candidate{}).Where("candidate_id = ?", candidateID).Updates(&candidate).Error
	return err
}