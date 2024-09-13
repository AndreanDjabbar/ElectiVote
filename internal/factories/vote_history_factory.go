package factories

import "github.com/AndreanDjabbar/ElectiVote/internal/models"

func VoteHistoryFactory(moderatorID, totalVotes uint, moderatorName, voteTitle, voteDescription, candidateWinnerName, candidateWinnerPicture string,  start models.CustomTime, end models.CustomTime) *models.VoteHistory {
	return &models.VoteHistory{
		ModeratorID:     moderatorID,
		TotalVotes:      totalVotes,
		ModeratorName:   moderatorName,
		VoteTitle:       voteTitle,
		VoteDescription: voteDescription,
		CandidateWinnerName: candidateWinnerName,
		CandidateWinnerPicture: candidateWinnerPicture,
		Start:           start,
		End:             end,
	}

}