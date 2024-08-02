package factories

import (
	"github.com/AndreanDjabbar/ElectiVote/internal/models"
)

func StartVoteFactory(voteTitle, voteDescription, voteCode string, moderatorID uint, start models.CustomTime) (models.Vote) {
	newVote := models.Vote{
		VoteTitle:       voteTitle,
		VoteDescription: voteDescription,
		VoteCode:        voteCode,
		ModeratorID:     moderatorID,
		Start:           start,
	}
	return newVote
}

func UpdateVoteFactory(voteTitle, voteDescription string) (models.Vote) {
	updatedVote := models.Vote{
		VoteTitle:       voteTitle,
		VoteDescription: voteDescription,
	}
	return updatedVote
}