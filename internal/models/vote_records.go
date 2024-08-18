package models

type VoteRecord struct {
	VoteRecordID uint `gorm:"primary_key"`
	VoteId       uint
	Vote         Vote `gorm:"foreignKey:VoteId;constraint:OnDelete:CASCADE;"`
	UserId 	 uint
	User         User `gorm:"foreignKey:UserId;constraint:OnDelete:CASCADE;"`
	CandidateId  uint
	Candidate    Candidate `gorm:"foreignKey:CandidateId;constraint:OnDelete:CASCADE;"`
	VotedTime   CustomTime `gorm:"type:datetime;default:NULL"`
}