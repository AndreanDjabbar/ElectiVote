package models

import (
	"time"
)

type Vote struct {
	VoteID          uint   `gorm:"primary_key"`
	VoteTitle       string `gorm:"type:varchar(255);not null"`
	VoteDescription string `gorm:"type:text;default:NULL"`
	VoteCode        string `gorm:"type:varchar(255);not null"`
	ModeratorID     uint
	User            User     `gorm:"foreignKey:ModeratorID;constraint:OnDelete:CASCADE;"`
	Start           time.Time `gorm:"type:datetime;default:NULL"`
	End             time.Time `gorm:"type:datetime;default:NULL"`
}