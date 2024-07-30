package factories

import (
	"time"

	"github.com/AndreanDjabbar/ElectiVote/internal/models"
)

func StartVoteFactory(voteTitle, voteDescription, voteCode string, moderatorID uint, start time.Time) (models.Vote) {
	newVote := models.Vote{
		VoteTitle:       voteTitle,
		VoteDescription: voteDescription,
		VoteCode:        voteCode,
		ModeratorID:     moderatorID,
		Start:           start,
	}
	return newVote
}