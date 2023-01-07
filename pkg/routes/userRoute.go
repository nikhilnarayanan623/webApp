package routes

import (
	"webApp/pkg/controllers"
	"webApp/pkg/middleware"

	"github.com/gin-gonic/gin"
)

func User(router *gin.Engine) {

	router.GET("/signup", middleware.UserAuthentication, controllers.SignupUser)
	router.POST("/signup", controllers.SigupSubmitUser)

	router.GET("/", middleware.UserAuthentication, controllers.LoginUser)
	router.POST("/", controllers.LoginSubmitUser)

	router.GET("/home", middleware.UserAuthentication, controllers.HomeUser)
	router.GET("/logout", controllers.LogoutUser)
}
