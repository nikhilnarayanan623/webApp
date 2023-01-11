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

	//setup router
	routes.Admin(router)
	routes.User(router)

	//no  rout found
	//invalid url then redirect to login page of user that middleware check if user is logged in then show home page
	//otherwise show login page
	router.NoRoute(func(ctx *gin.Context) {

		fmt.Println("invalid url so showing login page")
		ctx.Redirect(http.StatusSeeOther, "/")
	})

	router.Run(":8000")

}
