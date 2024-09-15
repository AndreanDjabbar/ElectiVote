package factories

import "github.com/AndreanDjabbar/ElectiVote/internal/models"

func FeedbackFactory(userID uint, feedbackMessage string, feedbackRate uint, feedbackDate models.CustomTime) models.Feedback {
	return models.Feedback{
		FeedbackMessage: feedbackMessage,
		FeedbackRate:    feedbackRate,
		FeedbackDate:    feedbackDate,
		UserId:          userID,
	}
}