package middleware

import (
	"fmt"
	"net/http"
	"time"
	"webApp/pkg/controllers"
	"webApp/pkg/db"
	"webApp/pkg/helper"
	"webApp/pkg/models"

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
			ctx.Redirect(http.StatusSeeOther, "/")
			return
		}

		//chekc the claims token userId is in database
		var user models.User
		db.DB.First(&user, "id = ?", claims["userId"])

		if user.ID == 0 || !user.Status { //admin id not matching
			ctx.Abort()
			fmt.Println("user not found but jwt is there admin deleted user")
			controllers.LoginUser(ctx)
			// ctx.Redirect(http.StatusSeeOther, "/")
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
