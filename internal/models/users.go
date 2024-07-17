package models

import "time"

type AuthToken struct {
	ID       uint   "gorm:primary_key"
	Username string `gorm:unique;type:varchar(255);not null`
	Token    string `gorm:unique;type:varchar(255);not null`
	Created  time.Time
	Expired  time.Time
}