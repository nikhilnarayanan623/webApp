package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nikhilnarayanan623/webApp/pkg/db"
	"github.com/nikhilnarayanan623/webApp/pkg/helper"
	"github.com/nikhilnarayanan623/webApp/pkg/initilizer"
	"github.com/nikhilnarayanan623/webApp/pkg/routes"
)

func init() {
	initilizer.LoadEnv()
	db.ConnnectToDb()
	db.MigrateToDB()
	helper.CreateAdmin()
}

func main() {

	router := gin.Default()
	//parse all templates
	router.LoadHTMLGlob("templates/*.html")

	//setup router for user and admin in different function
	routes.Admin(router)
	routes.User(router)

	//if invalid url found then show user login page
	router.NoRoute(func(ctx *gin.Context) {

		fmt.Println(ctx.Request.Method, "method")

		fmt.Println("invalid url so showing login page")
		ctx.Redirect(http.StatusSeeOther, "/")
	})

	router.Run(":8000")

}
