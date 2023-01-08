package routes

import (
	"github.com/nikhilnarayanan623/webApp/pkg/controllers"
	"github.com/nikhilnarayanan623/webApp/pkg/middleware"

	"github.com/gin-gonic/gin"
)

func User(router *gin.Engine) {

	router.GET("/signup", middleware.UserAuthentication, controllers.SignupUser)
	router.POST("/signup", controllers.SigupSubmitUser)

	router.GET("/", middleware.UserAuthentication, controllers.LoginUser)
	router.POST("/", controllers.LoginSubmitUser)

	router.GET("/home", middleware.UserAuthentication, controllers.HomeUser)
	router.GET("/logout", controllers.LogoutUser)

	//add to cart

	router.GET("/home/addToCart/:pid", middleware.UserAuthentication, controllers.AddToCartUser)

	router.GET("/cart", middleware.UserAuthentication, controllers.ShowCartUser)
	router.GET("/cart/removeFromCart/:pid", middleware.UserAuthentication, controllers.RemoveFromCartUser)
}
