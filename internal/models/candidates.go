package models

type Candidate struct {
	CandidateID          uint   `gorm:"primary_key"`
	CandidateName        string `gorm:"type:varchar(255);not null"`
	CandidateDescription string `gorm:"type:text;default:NULL"`
	TotalVotes           uint   `gorm:"type:int;default:0"`
	CandidatePicture     string `gorm:"type:varchar(255);default:NULL"`
	VoteId               uint
	Vote                 Vote   `gorm:"foreignKey:VoteId;constraint:OnDelete:CASCADE;"`
}