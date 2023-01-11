package helper

import (
	"fmt"
	"os"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/nikhilnarayanan623/webApp/pkg/db"
	"github.com/nikhilnarayanan623/webApp/pkg/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

// func to setup the jwt token parametes are:
// gin context for setup cookie
// name to store key and userdetails as value in jwt map
func JwtSetUp(ctx *gin.Context, name string, userId interface{}) bool {
	fmt.Println("jwt setup")

	cookieTime := time.Now().Add(10 * time.Minute).Unix()
	fmt.Println("jwt setup ", userId)

	// v := reflect.ValueOf(user)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": userId, //store the user id on token
		"exp":    cookieTime,
	})

	//create signed token string using env vaiable
	if tokenString, err := token.SignedString([]byte(os.Getenv("JWTCODE"))); err == nil {
		//set cookie

		ctx.SetCookie(name, tokenString, 10*60, "", "", false, true)
		fmt.Println("successfully setup jwt cookie")
		return true
	}

	fmt.Println("faild to setup jwt")
	return false
}

// get token if token is not in black list of dtabase
func GetToken(ctx *gin.Context, name string) (*jwt.Token, bool) {
	//delete expired token from black list database
	db.DeleteBlackListToken()

	//get the cookie
	cookieval, ok := GetCookieVal(ctx, name)

	if !ok { //problem to get cookie so return flase
		return nil, false
	}
	//check the user in black list or not
	var jwtBlack models.JwtBlackList

	db.DB.Find(&jwtBlack, "token_string = ?", cookieval)

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

	///return token with true
	return token, true
}

// to get cookie from client side
func GetCookieVal(ctx *gin.Context, name string) (string, bool) {

	if cookieVal, err := ctx.Cookie(name); err == nil {
		return cookieVal, true
	}

	fmt.Println("faild to get cookie")
	return "", false
}

// func to help validation
func CustomValidForUpdate(fl validator.FieldLevel) bool {

	value := fl.Field().String()

	return value == "" || len(value) > 3 || len(value) < 40 //if any of this condition true then dont take it as error
}
