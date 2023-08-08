package models

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	gorm.Model
	Name     string    `gorm:"column:name"`
	Password string    `gorm:"column:password"`
	Date     time.Time `gorm:"column:date"`
	Role     int       `gorm:"column:role"`
}

func (*User) TableName() string {
	return "users"
}
