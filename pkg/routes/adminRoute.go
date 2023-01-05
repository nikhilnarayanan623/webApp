package routes

import (
	"webApp/pkg/controllers"
	"webApp/pkg/middleware"

	"github.com/gin-gonic/gin"
)

func Admin(router *gin.Engine) {

	router.GET("/admin", middleware.AdminAuth, controllers.Login)
	router.POST("/admin", controllers.Submit)

	router.GET("/admin/home", middleware.AdminAuth, controllers.Home)
	router.GET("/admin/logout", middleware.AdminAuth, controllers.Logout)
}
