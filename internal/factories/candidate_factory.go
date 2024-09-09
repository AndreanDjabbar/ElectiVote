package factories

import (
	"github.com/AndreanDjabbar/ElectiVote/internal/models"
)

func CandidateFactory(candidateName, candidateDescription, candidatePicture string, voteID uint) models.Candidate {
	return models.Candidate{
		CandidateName: candidateName,
		VoteId:        voteID,
		CandidateDescription: candidateDescription,
		CandidatePicture: candidatePicture,
	}
}