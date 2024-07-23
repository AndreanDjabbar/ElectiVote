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
	FirstName   string `gorm:"type:varchar(255);not null"`
	LastName    string `gorm:"type:varchar(255);not null"`
	Age         uint   `gorm:"type:int;not null"`
	PhoneNumber string `gorm:"type:varchar(255);not null"`
	Birthday 	time.Time `gorm:"type:date;not null"`
	UserID      uint
	User      User   `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;"`
}