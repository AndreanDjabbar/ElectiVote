package models

import (
	"database/sql/driver"
	"fmt"
	"time"
)

type User struct {
	ID       uint   `gorm:"primary_key"`
	Username string `gorm:"unique;type:varchar(255);not null"`
	Password string `gorm:"type:varchar(255);not null"`
	Email    string `gorm:"unique;type:varchar(255);not null"`
	Role     string `gorm:type:enum('admin', 'user');not null`
}

type Profile struct {
    ProfileID   uint       `gorm:"primary_key"`
    Picture     string     `gorm:"type:varchar(255)"`
    FirstName   string     `gorm:"type:varchar(255)"`
    LastName    string     `gorm:"type:varchar(255)"`
    Age         uint       `gorm:"type:int;default:NULL"`
    PhoneNumber string     `gorm:"type:varchar(255)"`
    Birthday    NullTime   `gorm:"type:date;default:NULL"`
    UserID      uint
    User        User       `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;"`
}

type NullTime struct {
    Time  time.Time
    Valid bool
}

func (nt *NullTime) Scan(value interface{}) error {
    if value == nil {
        nt.Time, nt.Valid = time.Time{}, false
        return nil
    }
    nt.Valid = true
    switch v := value.(type) {
    case time.Time:
        nt.Time = v
    case []uint8:
        t, err := time.Parse("2006-01-02", string(v))
        if err != nil {
            return err
        }
        nt.Time = t
    default:
        return fmt.Errorf("unsupported Scan, storing driver.Value type %T into type *NullTime", value)
    }
    return nil
}

func (nt NullTime) Value() (driver.Value, error) {
    if !nt.Valid {
        return nil, nil
    }
    return nt.Time, nil
}