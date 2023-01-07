package controllers

import (
	"fmt"
	"net/http"
	"time"
	"webApp/pkg/db"
	"webApp/pkg/helper"
	"webApp/pkg/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"gorm.io/gorm/clause"
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
	ctx.Header("Cache-Control", "no-cache, no-store, must-revalidate")

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

	if !helper.JwtSetUp(ctx, "user", userVal) {
		ctx.Redirect(http.StatusSeeOther, "/")
		return
	}

	ctx.Redirect(http.StatusSeeOther, "/home")
}

// user home page
func HomeUser(ctx *gin.Context) {
	fmt.Println("Home user")

	ctx.Header("Cache-Control", "no-cache, no-store, must-revalidate")

	userId, _ := ctx.Get("userId")

	ctx.HTML(http.StatusOK, "userHome.html", userId)

}
func LogoutUser(ctx *gin.Context) {
	fmt.Println("logout user")

	cookieVal, ok := helper.GetCookieVal(ctx, "user")

	if !ok {
		ctx.Redirect(http.StatusSeeOther, "/")
		return
	}

	//get the token and check the token is expired
	if token, ok := helper.GetToken(ctx, "user"); ok {

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			//check the time is over if its not then add it black list
			if float64(time.Now().Unix()) < claims["exp"].(float64) {

				//add the cookieVal to black list
				db.DB.Clauses(clause.OnConflict{DoNothing: true}).Create(&models.JwtBlackList{
					TokenString: cookieVal,
					EndTime:     claims["exp"].(float64),
				})
			}
		}
	}
	//atlast redirect to login page
	ctx.Redirect(http.StatusSeeOther, "/")
}
