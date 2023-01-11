package helper

import (
	"fmt"
	"os"

	"github.com/nikhilnarayanan623/webApp/pkg/db"
	"github.com/nikhilnarayanan623/webApp/pkg/models"

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

		var errorMessge = map[string]string{}

		for _, er := range err.(validator.ValidationErrors) {

			errorMessge[er.Field()] = "Enter " + er.Field() + " Properly"
		}

		//return the error map and returrn false
		fmt.Println("error on validation of login form")
		return errorMessge, false
	}

	//check the admin in database
	var admin models.Admin

	db.DB.Find(&admin, "email = ?", form.Email)

	if admin.ID == 0 { //user not found

		return map[string]string{
			"Alert": "You are not an admin",
			"Color": "text-danger",
		}, false
	}

	//hash the password and check the password

	if err := bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(form.Password)); err != nil {

		fmt.Println("password no match")
		return map[string]string{"Password": "Wrong Password"}, false
	}

	//valid admin
	fmt.Println("valid admin login")

	return admin.ID, true //return the id of user in database
}

// create admin
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
