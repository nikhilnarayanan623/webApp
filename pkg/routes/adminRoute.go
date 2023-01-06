package routes

import (
	"webApp/pkg/controllers"
	"webApp/pkg/middleware"

	"github.com/gin-gonic/gin"
)

func Admin(router *gin.Engine) {

	router.GET("/admin", middleware.AdminAuth, controllers.LoginAdmin)
	router.POST("/admin", controllers.SubmitAdmin)

	router.GET("/admin/home", middleware.AdminAuth, controllers.HomeAdmin)
	router.GET("/admin/logout", controllers.LogoutAdmin)
}
