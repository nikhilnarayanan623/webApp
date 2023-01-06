package middleware

import (
	"fmt"
	"net/http"
	"time"
	"webApp/pkg/controllers"
	"webApp/pkg/helper"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func UserAuthentication(ctx *gin.Context) {

	token, ok := helper.GetToken(ctx, "user")

	if !ok {
		fmt.Println("aborted 1")
		ctx.Abort()
		controllers.LoginUser(ctx)
		return
	}

	//if the token is valid
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		//check its time is not over

		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			ctx.Abort()
			controllers.LoginUser(ctx)
			return
		}

		ctx.Set("userId", claims["userId"])

		if ctx.Request.URL.Path != "/home" {
			ctx.Redirect(http.StatusSeeOther, "/home")
			return
		}
		ctx.Next()
	} else {
		ctx.Abort()
		ctx.Redirect(http.StatusSeeOther, "/")
	}
}
