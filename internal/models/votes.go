package models

import (
	"database/sql/driver"
	"fmt"
	"time"
)

type CustomTime struct {
	time.Time
}

// Scan implements the sql.Scanner interface
func (ct *CustomTime) Scan(value interface{}) error {
	if value == nil {
		ct.Time = time.Time{}
		return nil
	}
	val, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("cannot convert %v to CustomTime", value)
	}
	t, err := time.Parse("2006-01-02 15:04:05", string(val))
	if err != nil {
		return err
	}
	ct.Time = t
	return nil
}

// Value implements the driver.Valuer interface
func (ct CustomTime) Value() (driver.Value, error) {
	return ct.Time.Format("2006-01-02 15:04:05"), nil
}

type Vote struct {
	VoteID          uint                  `gorm:"primary_key"`
	VoteTitle       string                `gorm:"type:varchar(255);not null"`
	VoteDescription string                `gorm:"type:text;default:NULL"`
	VoteCode        string                `gorm:"type:varchar(255);not null"`
	ModeratorID     uint
	User            User                  `gorm:"foreignKey:ModeratorID;constraint:OnDelete:CASCADE;"`
	Start           CustomTime `gorm:"type:datetime;default:NULL"`
}