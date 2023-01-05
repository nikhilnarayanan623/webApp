package helper

import (
	"fmt"
	"os"
	"webApp/pkg/db"
	"webApp/pkg/models"

	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm/clause"
)

// to validate a admin
func ValidateAdmin(form struct {
	Email    string `validate:"required,email"`
	Password string `validate:"required"`
}) (interface{}, bool) {

	//validate given values

	//crate obj of vlaidatior
	validate := validator.New()

	if err := validate.Struct(form); err != nil { //error in

		var errorMessge = map[string]bool{}

		for _, er := range err.(validator.ValidationErrors) {

			errorMessge[er.Field()] = true
		}

		//return the error map and returrn false
		fmt.Println("error on validation of login form")
		return errorMessge, false
	}

	//check the admin in database
	var admin models.Admin

	db.DB.First(&admin, "email = ?", form.Email)

	if admin.ID == 0 { //user not found
		fmt.Println("admin not found in db")
		return map[string]bool{"Email": true}, false
	}

	//hash the password and check the password

	if err := bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(form.Password)); err != nil {
		fmt.Println("password no match")
		return map[string]bool{"Password": true}, false
	}

	//valid admin
	fmt.Println("valid admin login")
	return admin, true
}
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
