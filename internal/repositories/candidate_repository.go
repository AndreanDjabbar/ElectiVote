package repositories

import (
	"github.com/AndreanDjabbar/ElectiVote/internal/db"
	"github.com/AndreanDjabbar/ElectiVote/internal/models"
)

func AddCandidate(newCandidate models.Candidate) (models.Candidate, error) {
	err := db.DB.Create(&newCandidate).Error
	return newCandidate, err
}