package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"
	"webApp/pkg/db"
	"webApp/pkg/helper"
	"webApp/pkg/models"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v4"
	"gorm.io/gorm/clause"
)

var adminMessage interface{}

func LoginAdmin(ctx *gin.Context) {
	fmt.Println("At Admin Login")

	ctx.Header("Cache-Control", "no-cache, no-store, must-revalidate")
	ctx.HTML(200, "adminLogin.html", adminMessage)
	adminMessage = nil

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

// delte user
func DeleteUserAdmin(ctx *gin.Context) {

	userId := ctx.Param("id")

	db.DB.Clauses(clause.OnConflict{DoNothing: true}).Unscoped().Delete(&models.User{}, "id = ?", userId)

	ctx.Redirect(http.StatusSeeOther, "/admin/home")
}

// block usr
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

// add product
func AddProductGet(ctx *gin.Context) {
	fmt.Println("at add product")
	ctx.Header("Cache-Control", "no-cache, no-store, must-revalidate")

	ctx.HTML(http.StatusOK, "addProduct.html", adminMessage)
	adminMessage = nil
}

// add product post !! two page have same handler so with the path i just redirect to page
func AddProductPost(ctx *gin.Context) {

	price, _ := strconv.Atoi(ctx.Request.PostFormValue("price"))

	// if err != nil {
	// 	fmt.Println("error on parsing price value from form")
	// 	ctx.Redirect(http.StatusSeeOther, "/admin/addProduct")

	// 	adminMessage = map[string]string{
	// 		"Price": "Enter Price Properly",
	// 	}
	// 	return
	// }
	fmt.Println("price int", price)

	var form = struct {
		ProductName string  `validate:"required"`
		Description string  `validate:"required"`
		Price       float64 `validate:"required,numeric"`
	}{
		ProductName: ctx.Request.PostFormValue("pname"),
		Price:       float64(price),
		Description: ctx.Request.PostFormValue("descritption"),
	}

	//add validation later

	validate := validator.New()

	if err := validate.Struct(form); err != nil { //form have error then find the error and show it on page

		var errorToTempl = map[string]string{}

		for _, er := range err.(validator.ValidationErrors) {

			errorToTempl[er.Field()] = "Enter " + er.Namespace() + "Properly"
		}

		adminMessage = errorToTempl

		//check the param and go to appropripriate page param have a value its where from
		from := ctx.Param("from")
		if from == "add" {
			ctx.Redirect(http.StatusSeeOther, "/admin/addProduct")
			return
		}

		ctx.Redirect(http.StatusSeeOther, "/admin/products")
		fmt.Println("herre")
		return

	}

	//chekc the product alredy in database if not then add it to database

	db.DB.Clauses(clause.OnConflict{DoNothing: true}).Create(&models.Product{
		Name:        form.ProductName,
		Price:       form.Price,
		Description: form.Description,
		StockIn:     true,
	})

	ctx.Redirect(http.StatusSeeOther, "/admin/products")
}

//products

func ShowProductsAdmin(ctx *gin.Context) {

	fmt.Println("at admin show products")

	ctx.Header("Cache-Control", "no-cache, no-store, must-revalidate")

	type Product struct { //this in the template name
		ID          int
		PID         uint
		ProductName string
		Price       float64
		Description string
		StockIn     bool
	}

	var records []models.Product //to store all products in an array form

	db.DB.Find(&records) //find all products

	//store all products in a new array product that we created

	var arrayProducts []Product

	for i, res := range records {

		arrayProducts = append(arrayProducts, Product{
			ID:          i + 1,
			PID:         res.PID,
			ProductName: res.Name,
			Price:       res.Price,
			Description: res.Description,
			StockIn:     res.StockIn,
		})
	}

	// show the page

	ctx.HTML(http.StatusOK, "adminProducts.html", arrayProducts)

}

func BlockOrDeleteProductAdmin(ctx *gin.Context) {
	fmt.Println("at block product")

	//get the params
	pid := ctx.Param("pid")
	status := ctx.Param("status")

	//convet pid to int

	if pidInt, err := strconv.Atoi(pid); err == nil { //there is no error

		if status == "block" { //block product
			db.DB.Model(&models.Product{}).Where("p_id = ?", pidInt).Update("stock_in", false)
		} else if status == "unblock" { //ublock product
			db.DB.Model(&models.Product{}).Where("p_id = ?", pidInt).Update("stock_in", true)
		} else if status == "delete" { //delete product
			db.DB.Clauses(clause.OnConflict{DoNothing: true}).Unscoped().Delete(&models.Product{}, "p_id = ?", pidInt)
		}
	}

	//after redirect to product page

	ctx.Redirect(http.StatusSeeOther, "/admin/products")
}
