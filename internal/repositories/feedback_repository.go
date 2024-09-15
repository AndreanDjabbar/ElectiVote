package repositories

import (
	"github.com/AndreanDjabbar/ElectiVote/internal/db"
	"github.com/AndreanDjabbar/ElectiVote/internal/models"
)

func CreateFeedback(feedback models.Feedback) (models.Feedback, error) {
	err := db.DB.Create(&feedback).Error
	if err != nil {
		return feedback, err
	}
	return feedback, nil
}