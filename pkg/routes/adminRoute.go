package routes

import (
	"webApp/pkg/controllers"

	"github.com/gin-gonic/gin"
)

func Admin(router *gin.Engine) {

	router.GET("/admin", controllers.Login)
	router.POST("/admin", controllers.Submit)

	router.GET("/admin/home", controllers.Home)
	router.GET("/admin/logout", controllers.Logout)
}
