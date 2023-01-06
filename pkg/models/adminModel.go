package models

import "gorm.io/gorm"

// model for admin
type Admin struct {
	gorm.Model

	Email    string `gorm:"unique; not null"`
	Password string `grom:"not null"`
}

// model to store balck listed jwt in db
type JwtBlackList struct {
	ID          uint    `gorm:"primarykey"`
	TokenString string  `gorm:"not null"`
	EndTime     float64 `gorm:"not null"`
}
