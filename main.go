package main

import (
	"webApp/pkg/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.New()

	//parse all templae
	router.LoadHTMLGlob("templates/*.html")

	routes.Admin(router)

	//no  rout found
	router.NoRoute(func(ctx *gin.Context) {
		ctx.Writer.Write([]byte(`
		<h1>Invalid Url</h1>
		`))
	})

	router.Run(":8000")
}
