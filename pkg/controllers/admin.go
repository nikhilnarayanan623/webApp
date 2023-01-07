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

	//adminId, _ := ctx.Get("adminId")

	var record []models.User
	db.DB.Find(&record)

	//type to store uer details in order way
	type field struct {
		ID        int
		UserId    uint
		FirstName string
		LastName  string
		Email     string
		Status    bool
	}
	//slice to store all user
	var arrayOfField []field

	//range through record and add it on slice

	for i, res := range record {
		arrayOfField = append(arrayOfField, field{
			ID:        i + 1,
			UserId:    res.ID,
			FirstName: res.FirstName,
			LastName:  res.LastName,
			Email:     res.Email,
			Status:    res.Status,
		})
	}

	fmt.Println("test")
	ctx.HTML(200, "adminHome.html", arrayOfField)
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

//delte user

func DeleteUserAdmin(ctx *gin.Context) {

	userId := ctx.Param("id")

	db.DB.Clauses(clause.OnConflict{DoNothing: true}).Unscoped().Delete(&models.User{}, "id = ?", userId)

	ctx.Redirect(http.StatusSeeOther, "/admin/home")
}
func BlockUserAdmin(ctx *gin.Context) {
	fmt.Println("at admin block")

	fmt.Println(ctx.Params)

	userId := ctx.Params.ByName("id")

	if ctx.Params.ByName("status") == "block" {
		db.DB.Model(&models.User{}).Where("id = ?", userId).Update("status", false)
	} else {
		db.DB.Model(&models.User{}).Where("id = ?", userId).Update("status", true)
	}

	ctx.Redirect(http.StatusSeeOther, "/admin/home")
}
