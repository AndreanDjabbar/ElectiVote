package repositories

import (
	"github.com/AndreanDjabbar/ElectiVote/internal/db"
	"github.com/AndreanDjabbar/ElectiVote/internal/models"
)

func SaveSupport(support models.Support) error {
	err := db.DB.Create(&support).Error
	if err != nil {
		return err
	}
	return nil
}