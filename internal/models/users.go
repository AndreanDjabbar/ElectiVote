package models

import "time"

type User struct {
	ID       uint   `gorm:"primary_key"`
	Username string `gorm:"unique;type:varchar(255);not null"`
	Password string `gorm:"type:varchar(255);not null"`
	Email    string `gorm:"unique;type:varchar(255);not null"`
	Role string `gorm:type:enum('admin', 'user');not null`
}

type AuthToken struct {
	ID       uint   "gorm:primary_key"
	Username string `gorm:unique;type:varchar(255);not null`
	Token    string `gorm:unique;type:varchar(255);not null`
	Created  time.Time
	Expired  time.Time
}