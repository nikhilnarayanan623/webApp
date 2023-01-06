package middleware

import (
	"fmt"
	"time"
	"webApp/pkg/controllers"
	"webApp/pkg/helper"

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
			controllers.LoginAdmin(ctx)
			return
		}
		//chekc the claims token useId is valid

		ctx.Set("adminId", claims["userId"]) //set the admin details in ctx

		if ctx.Request.URL.Path != "/admin/home" {
			ctx.Abort()
			ctx.Redirect(300, "/admin/home")
			return
		}
		ctx.Next()

	} else {
		ctx.Abort()
		ctx.Redirect(300, "/admin")
	}

}
