package models

type Candidate struct {
	CandidateID uint   `gorm:"primary_key"`
	CandidateName string `gorm:"type:varchar(255)"`
	CandidateDescription string `gorm:"type:text;default:NULL"`
	VoteID uint
	Vote Vote `gorm:"foreignKey:VoteID;constraint:OnDelete:CASCADE;"`
	TotalVotes uint `gorm:"type:int;default:0"`
	CandidatePicture string `gorm:"type:varchar(255)"`
}