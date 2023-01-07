package models

import "gorm.io/gorm"

// model for user
type User struct {
	gorm.Model

	FirstName string `gorm:"not null"`
	LastName  string `gorm:"not null"`

	Email    string `grom:"unique; not null"`
	Password string `gorm:"unique; not null"`
	Status   bool   `gorm:"not null"`
}
