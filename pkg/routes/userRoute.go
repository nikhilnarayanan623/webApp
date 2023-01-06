package routes

import (
	"webApp/pkg/controllers"

	"github.com/gin-gonic/gin"
)

func User(router *gin.Engine) {

	router.GET("/signup", controllers.SignupUser)
	router.POST("/signup", controllers.SigupSubmitUser)

	router.GET("/", controllers.LoginUser)
	router.POST("/", controllers.LoginSubmitUser)
	router.GET("/home", controllers.HomeUser)
	router.GET("/logout", controllers.LogoutUser)
}
