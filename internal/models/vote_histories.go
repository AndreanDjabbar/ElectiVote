package models

type VoteHistory struct {
	VoteHistoryID 	uint `gorm:"primary_key"`
	ModeratorID 	uint `gorm:"not null"` 
	ModeratorName 	string `gorm:"type:varchar(255);not null"`
	VoteTitle    	string `gorm:"type:varchar(255);not null"`
	VoteDescription string `gorm:"type:text;default:NULL"`
	CandidateWinnerName string `gorm:"type:varchar(255);default:None"`
	CandidateWinnerPicture string `gorm:"type:varchar(255);default:None"`
	TotalVotes uint `gorm:"type:int;default:0"`
	Start         	CustomTime `gorm:"type:datetime;default:NULL"`
	End           	CustomTime `gorm:"type:datetime;default:NULL"`
}