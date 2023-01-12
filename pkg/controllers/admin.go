package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/nikhilnarayanan623/webApp/pkg/db"
	"github.com/nikhilnarayanan623/webApp/pkg/helper"
	"github.com/nikhilnarayanan623/webApp/pkg/models"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v4"
	"gorm.io/gorm/clause"
)

var adminMessage interface{}

// func to render the login page for admin
func LoginAdmin(ctx *gin.Context) {
	fmt.Println("At Admin Login")

	ctx.Header("Cache-Control", "no-cache, no-store, must-revalidate")

	ctx.HTML(200, "adminLogin.html", adminMessage)
	adminMessage = map[string]string{}

}

// func to chekc the amdin form details when admin login
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
		adminMessage = adminVal
		ctx.Redirect(http.StatusSeeOther, "/admin")
		return
	}

	//set the jwt // admiVal is actully admin id
	if !helper.JwtSetUp(ctx, "admin", adminVal) { //func to setup the jwt
		//error to setup the token
		LoginAdmin(ctx)
		return
	}

	//valid admin

	ctx.Redirect(http.StatusSeeOther, "admin/home")
}

// to render the home page for admin and show the user details
func HomeAdmin(ctx *gin.Context) {
	fmt.Println("At Admin Home")

	ctx.Header("Cache-Control", "no-cache, no-store, must-revalidate")

	//adminId, _ := ctx.Get("adminId")

	var record []models.User

	db.DB.Find(&record)

	//struct to store user details and this fiels is in template
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

	ctx.HTML(200, "adminHome.html", arrayOfField)
}

// func to logout admin and store the token as token if its expire time not over
func LogoutAdmin(ctx *gin.Context) {
	fmt.Println("At Admin Logout")

	cookieVal, ok := helper.GetCookieVal(ctx, "admin") //to store on black list if time not  over

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

// func to delete specific user using params as user id
func DeleteUserAdmin(ctx *gin.Context) {

	userId := ctx.Param("id")
	//permenant delte changed to soft delete
	db.DB.Clauses(clause.OnConflict{DoNothing: true}).Delete(&models.User{}, "id = ?", userId)

	ctx.Redirect(http.StatusSeeOther, "/admin/home")
}

// block or unblock specific user using params as user id
func BlockUserAdmin(ctx *gin.Context) {
	fmt.Println("at admin block")

	userId := ctx.Params.ByName("id")

	if ctx.Params.ByName("status") == "block" { //chekc block or unblock
		db.DB.Model(&models.User{}).Where("id = ?", userId).Update("status", false)
	} else {
		db.DB.Model(&models.User{}).Where("id = ?", userId).Update("status", true)
	}

	ctx.Redirect(http.StatusSeeOther, "/admin/home")
}

// func to render a form page form admin to add product
func AddProductGet(ctx *gin.Context) {
	fmt.Println("at add product")
	ctx.Header("Cache-Control", "no-cache, no-store, must-revalidate")

	ctx.HTML(http.StatusOK, "addProduct.html", adminMessage)
	adminMessage = nil
}

// func to check value that get from value and validate and store on database
func AddProductPost(ctx *gin.Context) {

	price, err := strconv.Atoi(ctx.Request.PostFormValue("price"))

	//invalid value on price like string // -1 is not accepted in custon function
	if err != nil && ctx.Request.PostFormValue("price") != "" {
		price = -1
	}

	var form = struct {
		ProductName string  `validate:"CustomValidForAddProduct"`
		Description string  `validate:"CustomValidForAddProduct"`
		Price       float64 `validate:"CustomValidAddProductPrice,numeric"`
	}{
		ProductName: ctx.Request.PostFormValue("pname"),
		Price:       float64(price),
		Description: ctx.Request.PostFormValue("descritption"),
	}

	//ceate an instance of validator and register custom function

	validate := validator.New()
	validate.RegisterValidation("CustomValidForAddProduct", helper.CustomValidForAddProduct)
	validate.RegisterValidation("CustomValidAddProductPrice", helper.CustomValidAddProductPrice)

	//chekc if there is an error in the form according to custom function
	if err := validate.Struct(form); err != nil { //form have error then find the error and show it on page

		var errorToTempl = map[string]string{}

		for _, er := range err.(validator.ValidationErrors) {

			errorToTempl[er.Field()] = "Enter " + er.Namespace() + "Properly"
		}

		adminMessage = errorToTempl

		//check the param and go to appropripriate page because two page have same handler
		from := ctx.Param("from")
		if from == "add" {
			ctx.Redirect(http.StatusSeeOther, "/admin/addProduct")
			return
		}

		ctx.Redirect(http.StatusSeeOther, "/admin/products")
		return
	}

	//if any conflict like same product is already exist do nothin using the cluase
	db.DB.Clauses(clause.OnConflict{DoNothing: true}).Create(&models.Product{ //product name field is unique
		ProductName: form.ProductName,
		Price:       form.Price,
		Description: form.Description,
		StockIn:     true,
	})
	//after all ho to product page
	ctx.Redirect(http.StatusSeeOther, "/admin/products")
}

// func to render a  page that can show all product available in database
func ShowProductsAdmin(ctx *gin.Context) {

	fmt.Println("at admin show products")

	ctx.Header("Cache-Control", "no-cache, no-store, must-revalidate")

	var records []models.Product // to store all products in an array

	db.DB.Find(&records) //find all products and store the array

	//struct to store all value we want to show for a sing product including serial number as array indext
	type Product struct {
		ID          int
		PID         uint
		ProductName string
		Price       float64
		Description string
		StockIn     bool
	}
	//arra for all products
	var arrayProducts []Product

	for i, res := range records {

		arrayProducts = append(arrayProducts, Product{
			ID:          i + 1,
			PID:         res.PID,
			ProductName: res.ProductName,
			Price:       res.Price,
			Description: res.Description,
			StockIn:     res.StockIn,
		})
	}

	// render the html page
	ctx.HTML(http.StatusOK, "adminProducts.html", arrayProducts)

}

// to delet or block a product using the product id as params
func BlockOrDeleteProductAdmin(ctx *gin.Context) {
	fmt.Println("at block product")

	//get the params of product id and the status as block / unblock / delete
	pid := ctx.Param("pid")
	status := ctx.Param("status")

	//convet pid to integer form

	if pidInt, err := strconv.Atoi(pid); err == nil {

		if status == "block" { //block product
			db.DB.Model(&models.Product{}).Where("p_id = ?", pidInt).Update("stock_in", false)
		} else if status == "unblock" { //ublock product
			db.DB.Model(&models.Product{}).Where("p_id = ?", pidInt).Update("stock_in", true)
		} else if status == "delete" { //delete product
			db.DB.Clauses(clause.OnConflict{DoNothing: true}).Unscoped().Delete(&models.Product{}, "p_id = ?", pidInt)
		}
	}

	//redirect to the product page
	ctx.Redirect(http.StatusSeeOther, "/admin/products")
}

// variable of stuct to store product details and error when the form submiting
var forProductPage = struct {
	Product interface{} //to store key as string and value as any
	Error   interface{}
}{Product: map[string]interface{}{}, Error: map[string]string{}} //assign empty string/string map

// edit product
func EditProductGet(ctx *gin.Context) {
	fmt.Println("edit product get")

	ctx.Header("Cache-Control", "no-cache, no-store, must-revalidate")

	pid := ctx.Param("pid")

	if pidInt, err := strconv.Atoi(pid); err == nil { //if no error to convert string of pid to int
		//find the product
		var prodcut models.Product

		db.DB.Find(&prodcut, "p_id = ?", pidInt)

		//create strcut that can hold all value of product and store it on forProductPage
		forProductPage.Product = map[string]interface{}{
			"PID":         prodcut.PID,
			"ProductName": prodcut.ProductName,
			"Price":       prodcut.Price,
			"Description": prodcut.Description,
		}
	}

	ctx.HTML(http.StatusOK, "adminEditProduct.html", forProductPage)

	//clear old value that sored in maps
	forProductPage.Product = map[string]string{}
	forProductPage.Error = map[string]string{}

}

// func to get values from a the from and check there is an error or not if no error the update it
func EditProductPost(ctx *gin.Context) {
	fmt.Println("edit product post")

	pid := ctx.Param("pid")

	//first get the product price and convert into string
	price, err := strconv.Atoi(ctx.Request.PostFormValue("price"))
	//if error on converting to int then assing to -1 so the custom validation have condtion to check
	if err != nil && ctx.Request.PostFormValue("price") != "" {
		price = -1
	}

	var form = struct {
		ProductName string  `validate:"CustomValidForUpdate"`
		Description string  `validate:"CustomValidForUpdate"`
		Price       float64 `validate:"CustomValidProductPrice,numeric"`
	}{
		ProductName: ctx.Request.PostFormValue("pname"),
		Price:       float64(price),
		Description: ctx.Request.PostFormValue("descritption"),
	}

	//check if all fleild is empty then no need to update
	if form.ProductName == "" && form.Description == "" && price == 0 {

		forProductPage.Error = map[string]string{"Alert": "Enter one of the field to update", "Color": "text-danger"}
		EditProductGet(ctx)
		return
	}
	//validate form using validate function
	validate := validator.New()

	validate.RegisterValidation("CustomValidForUpdate", helper.CustomValidForUpdate)
	validate.RegisterValidation("CustomValidProductPrice", helper.CustomValidProductPrice)

	//validate the form
	if err := validate.Struct(form); err != nil {

		var formErrors = map[string]string{}

		for _, er := range err.(validator.ValidationErrors) {

			formErrors[er.Namespace()] = "Enter " + er.Field() + " Properly"
		}

		forProductPage.Error = formErrors //assign to datasToEditPage error
		EditProductGet(ctx)
		return
	}

	// if form is valid then convert pid to int and update the form value in database
	if pidInt, err := strconv.Atoi(pid); err == nil {

		db.DB.Model(&models.Product{}).Where("p_id = ?", pidInt).Updates(&models.Product{
			ProductName: form.ProductName,
			Price:       float64(price),
			Description: form.Description,
		})

		// sent an successfull message to page
		forProductPage.Error = map[string]string{"Alert": "Successfully Updated Details", "Color": "text-success"}
		EditProductGet(ctx)
		return
	}

	//if an error to convert pid to in then show cant update
	forProductPage.Error = map[string]string{"Alert": "Can't Update Details", "Color": "text-danger"}
	EditProductGet(ctx)
}
