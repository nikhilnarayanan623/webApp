package middleware

import (
	"fmt"
	"net/http"
	"time"
	"github.com/nikhilnarayanan623/webApp/pkg/controllers"
	"github.com/nikhilnarayanan623/webApp/pkg/db"
	"github.com/nikhilnarayanan623/webApp/pkg/helper"
	"github.com/nikhilnarayanan623/webApp/pkg/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func AdminAuth(ctx *gin.Context) {
	fmt.Println("Admin Auth")

	//get token using function that check the cookie is got and cookie alredy in black list

	token, ok := helper.GetToken(ctx, "admin")

	if !ok {
		ctx.Abort()
		controllers.LoginAdmin(ctx)
		return
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid { //valid token
		//chek the token time is over
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			fmt.Println("the token is timeouted")
			ctx.Abort()
			ctx.Redirect(http.StatusSeeOther, "/admin")
			return
		}
		//chekc the claims token userId is in database

		var admin models.Admin

		adminId := uint(claims["userId"].(float64))

		db.DB.Find(&admin, "id = ?", adminId)

		if admin.ID == 0 { //admin id not matching
			ctx.Abort()
			controllers.LoginAdmin(ctx)
			return
		}

		ctx.Set("adminId", adminId) //set the admin details in ctx

		if ctx.Request.URL.Path == "/admin" {
			ctx.Abort()
			ctx.Redirect(http.StatusSeeOther, "/admin/home")
			return
		}
		ctx.Next()

	} else {
		ctx.Abort()
		ctx.Redirect(http.StatusSeeOther, "/admin")
	}

}
