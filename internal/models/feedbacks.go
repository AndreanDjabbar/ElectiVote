package models

type Feedback struct {
	FeedbackID uint `gorm:"primary_key"`
	FeedbackMessage   string `gorm:"type:text;not null"`
	FeedbackRate uint `gorm:"type:int;not null"`
	FeedbackDate CustomTime `gorm:"type:datetime;default:NULL"`
	UserId uint
	User User `gorm:"foreignKey:UserId;constraint:OnDelete:CASCADE;"`
}