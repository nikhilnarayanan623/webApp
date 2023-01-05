package models

import "gorm.io/gorm"

// model for admin
type Admin struct {
	gorm.Model

	Email    string `gorm:"unique; not null"`
	Password string `grom:"not null"`
}

// model for user
type User struct {
	gorm.Model

	FirstName string `gorm:"not null"`
	LastName  string `gorm:"not null"`

	Email    string `grom:"unique; not null"`
	Password string `gorm:"unique; not null"`
}

// model to store balck listed jwt in db
type JwtBlackList struct {
	Token   string  `gorm:"not null"`
	EndTime float64 `gorm:"not null"`
}
