package factories

import (
	"github.com/AndreanDjabbar/ElectiVote/config"
	"github.com/AndreanDjabbar/ElectiVote/internal/models"
)

func CandidateFactory(candidateName, candidateDescription, candidatePicture string, voteID uint) models.Candidate {
	logger := config.SetUpLogger()
	logger.Info(
		"Candidate Factory - Creating Candidate",
	)
	return models.Candidate{
		CandidateName: candidateName,
		VoteId:        voteID,
		CandidateDescription: candidateDescription,
		CandidatePicture: candidatePicture,
	}
}