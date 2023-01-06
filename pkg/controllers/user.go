package controllers

import (
	"fmt"
	"net/http"
	"webApp/pkg/helper"

	"github.com/gin-gonic/gin"
)

var userMessage interface{} // to store all message that want to show in login and signup page

// singup hanler
func SignupUser(ctx *gin.Context) {
	fmt.Println("signup user")

	ctx.HTML(http.StatusOK, "userSignup.html", userMessage)
	userMessage = nil
}

func SigupSubmitUser(ctx *gin.Context) {
	fmt.Println("signup submit user")

	//validte the form value using a function that use validator package
	message, ok := helper.ValidateUserSubmit(struct {
		FirstName string `validate:"required"`
		LastName  string `validate:"required"`
		Email     string `validate:"required,email"`
		Password  string `validate:"required"`
	}{
		FirstName: ctx.Request.PostFormValue("fname"),
		LastName:  ctx.Request.PostFormValue("lname"),
		Email:     ctx.Request.PostFormValue("email"),
		Password:  ctx.Request.PostFormValue("password"),
	})

	if !ok {
		fmt.Println("not ok on form submit")
		userMessage = message
		SignupUser(ctx)
		return
	}

	//if is a valid form then the function will sore datas on database

	userMessage = message
	//there is no error then see the login page
	ctx.Redirect(http.StatusSeeOther, "/")

}

// login user
func LoginUser(ctx *gin.Context) {
	fmt.Println("login user")
	ctx.HTML(http.StatusOK, "userLogin.html", userMessage)

	userMessage = nil //after render html then dlete message
}

func LoginSubmitUser(ctx *gin.Context) {
	fmt.Println("login submit user")

	//validate user
	userVal, ok := helper.ValidateUserLogin(struct {
		Email    string `validate:"required,email"`
		Password string `validate:"required"`
	}{
		Email:    ctx.Request.PostFormValue("email"),
		Password: ctx.Request.PostFormValue("password"),
	})

	//if any probleme when user validation then show it
	if !ok {
		userMessage = userVal
		LoginUser(ctx)
		return
	}

	// if a valid user setyp jwt

	if !helper.JwtSetUp(ctx, userVal) {
		ctx.Redirect(http.StatusSeeOther, "/")
		return
	}

	ctx.Redirect(http.StatusSeeOther, "/home")
}

// user home page
func HomeUser(ctx *gin.Context) {
	fmt.Println("Home user")
	ctx.HTML(http.StatusOK, "userHome.html", nil)

}
func LogoutUser(ctx *gin.Context) {
	fmt.Println("logout user")

	ctx.Redirect(http.StatusSeeOther, "/")
}
