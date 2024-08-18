package factories

import "github.com/AndreanDjabbar/ElectiVote/internal/models"

func VoteRecordFactory(voteID, userID, candidateID uint, start models.CustomTime) models.VoteRecord {
	return models.VoteRecord{
		VoteId:      voteID,
		UserId:      userID,
		CandidateId: candidateID,
		VotedTime:  start,
	}
}