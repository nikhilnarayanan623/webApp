package routes

import (
	"github.com/nikhilnarayanan623/webApp/pkg/controllers"
	"github.com/nikhilnarayanan623/webApp/pkg/middleware"

	"github.com/gin-gonic/gin"
)

func User(router *gin.Engine) {
	// signup and signup submit
	router.GET("/signup", middleware.UserAuthentication, controllers.SignupUser)
	router.POST("/signup", controllers.SigupSubmitUser)

	//login and login submit
	router.GET("/", middleware.UserAuthentication, controllers.LoginUser)
	router.POST("/", controllers.LoginSubmitUser)

	//home page and logut
	router.GET("/home", middleware.UserAuthentication, controllers.HomeUser)
	router.GET("/logout", controllers.LogoutUser)

	//show cart
	router.GET("/cart", middleware.UserAuthentication, controllers.ShowCartUser)

	//add to cart and remove from cart
	router.GET("/home/addToCart/:pid", middleware.UserAuthentication, controllers.AddToCartUser)
	router.GET("/cart/removeFromCart/:pid", middleware.UserAuthentication, controllers.RemoveFromCartUser)

	//edit user details get and post for get page and submit details
	router.GET("/edituser/", middleware.UserAuthentication, controllers.EditUserGet)
	router.POST("/edituser", middleware.UserAuthentication, controllers.EditUserPost)

}
