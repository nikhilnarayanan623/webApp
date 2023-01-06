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

var LoginMessage interface{}

func LoginAdmin(ctx *gin.Context) {
	fmt.Println("At Admin Login")

	ctx.Header("Cache-Control", "no-cache, no-store, must-revalidate")
	ctx.HTML(200, "adminLogin.html", LoginMessage)
	userMessage = nil

}

func SubmitAdmin(ctx *gin.Context) {
	fmt.Println("At Admin Submit")

	adminVal, ok := helper.ValidateAdmin(struct { //call validation helper function
		Email    string `validate:"required,email"`
		Password string `validate:"required"`
	}{
		Email:    ctx.Request.PostFormValue("email"),
		Password: ctx.Request.PostFormValue("password"),
	})

	if !ok {
		LoginMessage = adminVal
		LoginAdmin(ctx)
		return
	}

	//set the jwt // admiVal is actully admin id
	if !helper.JwtSetUp(ctx, "admin", adminVal) { //func to setup the jwt
		//error to setup the token
		LoginAdmin(ctx)
		return
	}

	//valid admin

	ctx.Redirect(http.StatusSeeOther, "admin//home")
}

func HomeAdmin(ctx *gin.Context) {
	fmt.Println("At Admin Home")

	ctx.Header("Cache-Control", "no-cache, no-store, must-revalidate")

	value, _ := ctx.Get("adminId")

	ctx.HTML(200, "adminHome.html", value)
}

// logout
func LogoutAdmin(ctx *gin.Context) {
	fmt.Println("At Admin Logout")

	cookieVal, ok := helper.GetCookieVal(ctx, "admin")

	if !ok {
		ctx.Redirect(http.StatusTemporaryRedirect, "/admin")
		return
	}

	//if token is there then add it blacklist if its time not out
	token, ok := helper.GetToken(ctx, "admin")

	if ok {
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			//chekc the token time not over
			tokenTime := claims["exp"].(float64)

			if float64(time.Now().Unix()) < tokenTime { //if time is not over then add this to black listdb

				db.DB.Clauses(clause.OnConflict{DoNothing: true}).Create(&models.JwtBlackList{
					TokenString: cookieVal,
					EndTime:     tokenTime,
				})
			}
		}
	}

	ctx.Redirect(http.StatusTemporaryRedirect, "/admin")
}
