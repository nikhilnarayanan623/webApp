package helper

import (
	"os"
	"webApp/pkg/db"
	"webApp/pkg/models"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm/clause"
)

func CreateAdmin() {

	aEmail := os.Getenv("ADMINEMAIL")
	aPass := os.Getenv("ADMINPASS")

	//hash the password and if no error create admin
	if hashPass, err := bcrypt.GenerateFromPassword([]byte(aPass), 10); err == nil {

		db.DB.Clauses(clause.OnConflict{DoNothing: true}).Create(&models.Admin{
			Email:    aEmail,
			Password: string(hashPass),
		})
	}

}
