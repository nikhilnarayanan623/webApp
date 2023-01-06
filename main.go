package main

import (
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
	router.NoRoute(func(ctx *gin.Context) {
		ctx.Writer.Write([]byte(`
		<h1>Invalid Url</h1>
		`))
	})

	router.Run(":8000")
}
