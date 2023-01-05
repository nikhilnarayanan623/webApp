package controllers

import (
	"fmt"
	"net/http"
	"webApp/pkg/helper"

	"github.com/gin-gonic/gin"
)

var LoginMessage interface{}

func Login(ctx *gin.Context) {
	fmt.Println("At Admin Login")
	ctx.HTML(200, "adminLogin.html", LoginMessage)
}

func Submit(ctx *gin.Context) {
	fmt.Println("At Admin Submit")

	validRes, ok := helper.ValidateAdmin(struct {
		Email    string `validate:"required,email"`
		Password string `validate:"required"`
	}{
		Email:    ctx.Request.PostFormValue("email"),
		Password: ctx.Request.PostFormValue("password"),
	})

	if !ok {
		LoginMessage = validRes
		Login(ctx)
	}

	adminDetails = validRes
	//valid admin
	Home(ctx)
}

var adminDetails interface{}

func Home(ctx *gin.Context) {
	fmt.Println("At Admin Home")

	ctx.HTML(200, "adminHome.html", adminDetails)
}

func Logout(ctx *gin.Context) {
	fmt.Println("At Admin Logout")
	ctx.Redirect(http.StatusTemporaryRedirect, "/admin")
}
