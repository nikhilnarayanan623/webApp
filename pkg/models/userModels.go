package models

import (
	"github.com/lib/pq"
	"gorm.io/gorm"
)

// model for user
type User struct {
	gorm.Model

	FirstName string `gorm:"not null"`
	LastName  string `gorm:"not null"`

	Email    string `gorm:"type:varchar(100);unique"`
	Password string `gorm:"not null"`
	Status   bool   `gorm:"not null"`

	Products pq.Int64Array `gorm:"type:integer[]"`
}

// ProductId []int `gorm:"type:integer[]"`
