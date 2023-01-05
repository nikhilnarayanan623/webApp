package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Login(ctx *gin.Context) {
	fmt.Println("At Admin Login")
	ctx.HTML(200, "adminLogin.html", nil)
}

func Submit(ctx *gin.Context) {
	fmt.Println("At Admin Submit")

	//valid
	Home(ctx)
}

func Home(ctx *gin.Context) {
	fmt.Println("At Admin Home")

	ctx.HTML(200, "adminHome.html", nil)
}

func Logout(ctx *gin.Context) {
	fmt.Println("At Admin Logout")
	ctx.Redirect(http.StatusTemporaryRedirect, "/admin")
}
