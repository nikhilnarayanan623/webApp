package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/nikhilnarayanan623/webApp/pkg/db"
	"github.com/nikhilnarayanan623/webApp/pkg/helper"
	"github.com/nikhilnarayanan623/webApp/pkg/models"
	"golang.org/x/crypto/bcrypt"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var userMessage interface{} // to store all message that want to show in login and signup page

// singup hanler
func SignupUser(ctx *gin.Context) {
	fmt.Println("signup user")

	ctx.Header("Cache-Control", "no-cache, no-store, must-revalidate")

	ctx.HTML(http.StatusOK, "userSignup.html", userMessage)
	userMessage = nil
}

func SigupSubmitUser(ctx *gin.Context) {
	fmt.Println("signup submit user")

	//validte the form value using a function that use validator package
	message, ok := helper.ValidateUserSubmit(struct {
		FirstName string `validate:"required"`
		LastName  string `validate:"required"`
		Email     string `validate:"required,email"`
		Password  string `validate:"required"`
	}{
		FirstName: ctx.Request.PostFormValue("fname"),
		LastName:  ctx.Request.PostFormValue("lname"),
		Email:     ctx.Request.PostFormValue("email"),
		Password:  ctx.Request.PostFormValue("password"),
	})

	if !ok {
		fmt.Println("not ok on form submit")
		userMessage = message
		SignupUser(ctx)
		return
	}

	//if is a valid form then the function will sore datasToEditPage on database

	userMessage = message
	//there is no error then see the login page
	ctx.Redirect(http.StatusSeeOther, "/")

}

// login user
func LoginUser(ctx *gin.Context) {
	fmt.Println("login user")
	ctx.Header("Cache-Control", "no-cache, no-store, must-revalidate")

	ctx.HTML(http.StatusOK, "userLogin.html", userMessage)

	userMessage = nil //after render html then dlete message
}

// login submit
func LoginSubmitUser(ctx *gin.Context) {
	fmt.Println("login submit user")

	//validate user
	userVal, ok := helper.ValidateUserLogin(struct {
		Email    string `validate:"required,email"`
		Password string `validate:"required"`
	}{
		Email:    ctx.Request.PostFormValue("email"),
		Password: ctx.Request.PostFormValue("password"),
	})

	//if any probleme when user validation then show it on login page
	if !ok {
		userMessage = userVal
		ctx.Redirect(http.StatusSeeOther, "/")
		return
	}

	// if a valid user setyp jwt

	if !helper.JwtSetUp(ctx, "user", userVal) {
		ctx.Redirect(http.StatusSeeOther, "/")
		return
	}

	ctx.Redirect(http.StatusSeeOther, "/home")
}

// user home page
func HomeUser(ctx *gin.Context) {
	fmt.Println("Home user")

	ctx.Header("Cache-Control", "no-cache, no-store, must-revalidate")

	//get all products that stock in db and show it

	var products []models.Product
	db.DB.Find(&products, "stock_in = ?", true) //find all product

	type fields struct {
		UserName    string
		PID         uint
		ProductName string
		Price       float64
		Description string
	}

	var results []fields

	for _, res := range products {

		results = append(results, fields{ //append all value to that slice
			PID:         res.PID,
			ProductName: res.Name,
			Price:       res.Price,
			Description: res.Description,
		})
	}

	//find the user and add the username to results
	userId, _ := ctx.Get("userId") // user id from context

	var user models.User
	db.DB.Find(&user, "id = ?", userId)

	//create a struct that can hold products and userNam

	var PassValue = struct { //to pass value to template
		UserName string
		Products []fields
	}{
		UserName: user.FirstName,
		Products: results,
	}

	ctx.HTML(http.StatusOK, "userHome.html", PassValue)

}

// Logout
func LogoutUser(ctx *gin.Context) {
	fmt.Println("logout user")

	cookieVal, ok := helper.GetCookieVal(ctx, "user")

	if !ok {
		ctx.Redirect(http.StatusSeeOther, "/")
		return
	}

	//get the token and check the token is expired
	if token, ok := helper.GetToken(ctx, "user"); ok {

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			//check the time is over if its not then add it black list
			if float64(time.Now().Unix()) < claims["exp"].(float64) {

				//add the cookieVal to black list
				db.DB.Clauses(clause.OnConflict{DoNothing: true}).Create(&models.JwtBlackList{
					TokenString: cookieVal,
					EndTime:     claims["exp"].(float64),
				})
			}
		}
	}
	//atlast redirect to login page
	ctx.Redirect(http.StatusSeeOther, "/")
}

//add to cart for user

func AddToCartUser(ctx *gin.Context) {
	fmt.Println("add to cart")

	userId, _ := ctx.Get("userId") // user id from context

	pid := ctx.Params.ByName("pid")

	//convver pid to integer
	if pidInt, err := strconv.Atoi(pid); err == nil {

		// u := uint(userId.(float64))
		//append the pid to users product column array
		db.DB.Model(&models.User{}).Where("id = ?", userId).Update("Products", gorm.Expr("array_append(Products, ?)", pidInt))

	}
	fmt.Println("user id at add to cart ", userId)

	ctx.Redirect(http.StatusSeeOther, "/home")
}

// show the cart
func ShowCartUser(ctx *gin.Context) {

	fmt.Println("at Show cart")
	ctx.Header("Cache-Control", "no-cache, no-store, must-revalidate")

	userId, _ := ctx.Get("userId") // user id from context

	var user = models.User{}

	db.DB.Find(&user, "id = ?", userId) //find the user from database

	type Cart struct { // to store each details we want from database
		PID         uint
		ProductName string
		Price       float64
		StockIn     bool
	}

	var arrayOfCart []Cart //to store all detail of product

	for _, res := range user.Products { //range through the products id array we got from user

		var product models.Product

		db.DB.Find(&product, "p_id = ?", res)

		//check the product is deleted or not

		if product.PID == 0 {
			continue
		}

		arrayOfCart = append(arrayOfCart, Cart{
			PID:         product.PID,
			ProductName: product.Name,
			Price:       product.Price,
			StockIn:     product.StockIn,
		})
	}
	fmt.Println("last of show cart")
	ctx.HTML(http.StatusOK, "userCart.html", arrayOfCart)

}

// remove from cart
func RemoveFromCartUser(ctx *gin.Context) {

	fmt.Println("at remove from cart")

	//get the product id from params

	pid := ctx.Param("pid")

	//get the user id form context
	if userId, ok := ctx.Get("userId"); ok {
		//conver pid string to integer

		if pidInt, err := strconv.Atoi(pid); err == nil { //there is no error to parse the pid to int

			db.DB.Model(&models.User{}).Where("id = ?", userId).Update("Products", gorm.Expr("array_remove(Products, ?)", pidInt))
		}
	}

	ctx.Redirect(http.StatusSeeOther, "/cart")
}

// edit user get
type data struct {
	UserFristName string
	UserLastName  string
	UserEmail     string

	//for error
	Error interface{}
}

var datasToEditPage = data{Error: map[string]string{}} //emtpy map for if no error on form

// ediuser get func
func EditUserGet(ctx *gin.Context) {
	fmt.Println("at edit user")
	ctx.Header("Cache-Control", "no-cache, no-store, must-revalidate")

	//type to store user details and error to show

	// get user id
	if userId, ok := ctx.Get("userId"); ok {
		//get the user from database

		var user models.User
		db.DB.Find(&user, "id = ?", userId)

		datasToEditPage.UserFristName = user.FirstName
		datasToEditPage.UserLastName = user.LastName
		datasToEditPage.UserEmail = user.Email
	}

	ctx.HTML(http.StatusOK, "userEditProfile.html", datasToEditPage)
	datasToEditPage.Error = map[string]string{} //clear errors
}

// edit user post
func EditUserPost(ctx *gin.Context) {
	fmt.Println("edit post user")

	userId, ok := ctx.Get("userId")

	if !ok { //didnt get userid the sent an alert to page
		datasToEditPage.Error = map[string]string{"Alert": "Can't Updated Details", "Color": "text-success"}
		ctx.Redirect(http.StatusSeeOther, "/edituser")
		return
	}

	//validte the form value using a function that use validator package
	var form = struct {
		FirstName string `validate:"CustomValidForUpdate"`
		LastName  string `validate:"CustomValidForUpdate"`
		Email     string `validate:"CustomValidForUpdate,email"`
		Password  string
	}{
		FirstName: ctx.Request.PostFormValue("fname"),
		LastName:  ctx.Request.PostFormValue("lname"),
		Email:     ctx.Request.PostFormValue("email"),
		Password:  ctx.Request.PostFormValue("password"),
	}
	// check all field is empty if empty no need to validate or update
	if form.Email == "" && form.FirstName == "" && form.LastName == "" && form.Password == "" {
		datasToEditPage.Error = map[string]string{"Alert": "Enter one of the field to update", "Color": "text-danger"}
		ctx.Redirect(http.StatusSeeOther, "/edituser")
		return
	}

	//if some field have valued then validate it

	validate := validator.New()
	validate.RegisterValidation("CustomValidForUpdate", helper.CustomValidForUpdate)

	if err := validate.Struct(form); err != nil {

		var formErrors = map[string]string{}

		for _, er := range err.(validator.ValidationErrors) {

			formErrors[er.Namespace()] = "Enter " + er.Field() + " Properly"
		}

		datasToEditPage.Error = formErrors //assign to datasToEditPage error
		ctx.Redirect(http.StatusSeeOther, "/edituser")
		return
	}

	//check if password is empty then no need to hash the pass and update the pass

	var result *gorm.DB

	if form.Password == "" {
		result = db.DB.Model(&models.User{}).Where("id = ?", userId).Updates(&models.User{
			FirstName: form.FirstName,
			LastName:  form.LastName,
			Email:     form.Email,
		})

	} else { // hash the password and update

		if hashPass, err := bcrypt.GenerateFromPassword([]byte(form.Password), 10); err == nil { //no error to hash the password

			result = db.DB.Model(&models.User{}).Where("id = ?", userId).Updates(&models.User{
				FirstName: form.FirstName,
				LastName:  form.LastName,
				Email:     form.Email,

				Password: string(hashPass),
			})
		}

	}

	if result.Error != nil { //error user is already exist

		datasToEditPage.Error = map[string]string{"Alert": "User Alredy Exist", "Color": "text-danger"}
		ctx.Redirect(http.StatusSeeOther, "/edituser")
		return
	}

	//sent an successfull message to page
	datasToEditPage.Error = map[string]string{"Alert": "Successfully Updated Details", "Color": "text-success"}
	ctx.Redirect(http.StatusSeeOther, "/edituser")
}
