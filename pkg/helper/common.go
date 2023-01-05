package helper

import (
	"fmt"
	"os"
	"time"
	"webApp/pkg/db"
	"webApp/pkg/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

// func to setup the jwt token parametes are:
// gin context for setup cookie
// name to store key and userdetails as value in jwt map
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
	}

	fmt.Println("faild to setup jwt")
	return false
}

// check the tokenstrng is in black list db
// func LogoutedUser(ctx *gin.Context) bool {

// 	if cookieval, ok := GetCookieVal(ctx); ok {
// 		//chekc the cookieval on black db
// 		var jwtBlack models.JwtBlackList

// 		db.DB.First(&jwtBlack, "token = ?", cookieval)

// 		if jwtBlack.ID == 0 {
// 			return true // this cookie not in db
// 		}
// 	}

// 	//if cookie didnt get or user in black list return flase
// 	return false
// }

func GetToken(ctx *gin.Context) (*jwt.Token, bool) {

	//get the cookie
	cookieval, ok := GetCookieVal(ctx)

	if !ok { //problem to get cookie so return flase
		return nil, false
	}
	//check the user in black list or not
	var jwtBlack models.JwtBlackList
	db.DB.First(&jwtBlack, "token = ?", cookieval)

	if jwtBlack.ID != 0 {
		return nil, false //this user is in black list
	}

	//parse the cookieval to jwt token
	token, err := jwt.Parse(cookieval, func(token *jwt.Token) (interface{}, error) {

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(os.Getenv("JWTCODE")), nil
	})

	if err != nil {
		fmt.Println("faild to parse the cookie to token")
		return nil, false
	}

	return token, true
}

// to get cookie from client side
func GetCookieVal(ctx *gin.Context) (string, bool) {

	if cookieVal, err := ctx.Cookie("jwt-auth"); err == nil {
		return cookieVal, true
	}

	fmt.Println("faild to get cookie")
	return "", false
}
