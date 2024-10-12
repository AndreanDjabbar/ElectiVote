package models

type Support struct {
	SupportID uint `gorm:"primary_key"`
	DonatorName string `gorm:"type:varchar(255);default:NULL"`
	DonatorEmail string `gorm:"type:varchar(255);default:NULL"`
	TransactionID string `gorm:"type:varchar(255);default:NULL"`
	Amount      float64 `gorm:"type:float;default:NULL"`
	SupportedTime   CustomTime `gorm:"type:datetime;default:NULL"`
	Message    string `gorm:"type:text;default:NULL"`
}