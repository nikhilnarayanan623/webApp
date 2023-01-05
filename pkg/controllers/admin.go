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

func Login(ctx *gin.Context) {
	fmt.Println("At Admin Login")

	ctx.Header("Cache-Control", "no-cache, no-store, must-revalidate")
	ctx.HTML(200, "adminLogin.html", LoginMessage)

}

func Submit(ctx *gin.Context) {
	fmt.Println("At Admin Submit")

	validRes, ok := helper.ValidateAdmin(struct { //call validation helper function
		Email    string `validate:"required,email"`
		Password string `validate:"required"`
	}{
		Email:    ctx.Request.PostFormValue("email"),
		Password: ctx.Request.PostFormValue("password"),
	})

	if !ok {
		LoginMessage = validRes
		Login(ctx)
		return
	}

	adminDetails = validRes

	//set the jwt
	if !helper.JwtSetUp(ctx, "admin", validRes) { //func to setup the jwt
		//error to setup the token
		Login(ctx)
		return
	}

	//valid admin
	//ctx.Header("Cache-Control", "no-cache, no-store, must-revalidate")

	ctx.Redirect(http.StatusSeeOther, "admin//home")
}

var adminDetails interface{}

func Home(ctx *gin.Context) {
	fmt.Println("At Admin Home")
	ctx.Header("Cache-Control", "no-cache, no-store, must-revalidate")
	ctx.HTML(200, "adminHome.html", adminDetails)
}

func Logout(ctx *gin.Context) {
	fmt.Println("At Admin Logout")
	//ctx.Header("Cache-Control", "no-cache, no-store, must-revalidate")
	cookieVal, ok := helper.GetCookieVal(ctx)

	if !ok {
		ctx.Redirect(http.StatusTemporaryRedirect, "/admin")
		return
	}

	//if token is there then add it blacklist if its time not out
	token, ok := helper.GetToken(ctx)

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
