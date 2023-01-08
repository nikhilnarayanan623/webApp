package main

import (
	"net/http"
	"webApp/pkg/db"
	"webApp/pkg/helper"
	"webApp/pkg/initilizer"
	"webApp/pkg/routes"

	"github.com/gin-gonic/gin"
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
		ctx.Redirect(http.StatusSeeOther, "/")
	})

	router.Run(":8000")
}
