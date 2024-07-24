package models

import "time"

type User struct {
	ID       uint   `gorm:"primary_key"`
	Username string `gorm:"unique;type:varchar(255);not null"`
	Password string `gorm:"type:varchar(255);not null"`
	Email    string `gorm:"unique;type:varchar(255);not null"`
	Role     string `gorm:type:enum('admin', 'user');not null`
}

type Profile struct {
	ProfileID   uint   `gorm:"primary_key"`
	Picture	 string `gorm:"type:varchar(255)"`
	FirstName   string `gorm:"type:varchar(255)"`
	LastName    string `gorm:"type:varchar(255)"`
	Age         uint   `gorm:"type:int;default:NULL"`
	PhoneNumber string `gorm:"type:varchar(255)"`
	Birthday 	*time.Time `gorm:"type:date;default:NULL"`
	UserID      uint
	User      User   `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;"`
}