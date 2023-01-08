package helper

import (
	"fmt"

	"github.com/nikhilnarayanan623/webApp/pkg/db"
	"github.com/nikhilnarayanan623/webApp/pkg/models"

	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm/clause"
)

func ValidateUserLogin(form struct {
	Email    string `validate:"required,email"`
	Password string `validate:"required"`
}) (interface{}, bool) {

	//chekc the form is valid using validator package
	validate := validator.New()

	if err := validate.Struct(form); err != nil {
		var templateMessage = map[string]string{}

		for _, er := range err.(validator.ValidationErrors) {

			templateMessage[er.Field()] = "Enter " + er.Field() + " Properly"
		}

		return templateMessage, false
	}

	//chekc the user is in database
	var user models.User

	db.DB.First(&user, "email = ?", form.Email)

	if user.ID == 0 { //user not found
		return map[string]string{
			"Alert": "You are not a registered user you can signup",
			"Color": "text-danger",
		}, false
	}
	//hash the password and check it on db pass
	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(form.Password)) != nil {
		return map[string]string{
			"Password": "Wrong Password",
		}, false
	}

	//check the user is blocked or not
	if !user.Status {
		return map[string]string{
			"Color": "text-danger",
			"Alert": "You are blocked by admin",
		}, false
	}

	//valid user so return user id
	return user.ID, true
}

//validate user signup

func ValidateUserSubmit(form struct {
	FirstName string `validate:"required"`
	LastName  string `validate:"required"`
	Email     string `validate:"required,email"`
	Password  string `validate:"required"`
}) (interface{}, bool) {

	//chekc the from is valid or not
	validate := validator.New()

	if err := validate.Struct(form); err != nil {

		TempMessagese := map[string]string{}

		for _, er := range err.(validator.ValidationErrors) {

			fmt.Println("***", er.Error(), "***")

			TempMessagese[er.Namespace()] = "Enter " + er.Namespace() + "Properly"
		}

		return TempMessagese, false
	}

	//chek the user already exist

	var user models.User

	db.DB.First(&user, "email = ?", form.Email)

	fmt.Println("test1")

	if user.ID != 0 { //the user alredy exist
		fmt.Println("usre alredy exist")
		return map[string]string{"Alert": "user alredy exist"}, false
	}

	//user not exist then hash the pass and store it database

	hasPass, err := bcrypt.GenerateFromPassword([]byte(form.Password), 10)

	if err != nil {
		fmt.Println("hash error")
		return map[string]string{"Password": "Error"}, false
	}

	//there is no error to hash the pass

	db.DB.Clauses(clause.OnConflict{DoNothing: true}).Create(&models.User{
		FirstName: form.FirstName,
		LastName:  form.LastName,
		Email:     form.Email,
		Password:  string(hasPass),
		Status:    true,
	})
	return map[string]string{"Color": "text-success",
		"Alert": "Sucessfully Account Created You Can Login",
	}, true //everyting ok
}
