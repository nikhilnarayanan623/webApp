package helper

import (
	"fmt"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func JwtSetUp(ctx *gin.Context, name string, user interface{}) bool {
	fmt.Println("jwt setup")

	cookieTime := time.Now().Add(1 * time.Minute).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		name:  user,
		"exp": cookieTime,
	})

	//create signed token string using env vaiable
	if tokenString, err := token.SignedString([]byte(os.Getenv("JWTCODE"))); err == nil {
		//set cookie

		ctx.SetCookie("jwt-auth", tokenString, 1*60, "", "", false, true)
		fmt.Println("successfully setup jwt cookie")
		return true
	} else {
		fmt.Println(err)
	}

	fmt.Println("faild to setup jwt")
	return false
}

func GetCookieVal(ctx *gin.Context) (string, bool) {

	if cookieVal, err := ctx.Cookie("jwt-auth"); err == nil {
		return cookieVal, true
	}

	fmt.Println("faild to get cookie")
	return "", false
}
