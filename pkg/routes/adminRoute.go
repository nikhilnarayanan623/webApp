package routes

import (
	"github.com/nikhilnarayanan623/webApp/pkg/controllers"
	"github.com/nikhilnarayanan623/webApp/pkg/middleware"

	"github.com/gin-gonic/gin"
)

func Admin(router *gin.Engine) {

	router.GET("/admin", middleware.AdminAuth, controllers.LoginAdmin)
	router.POST("/admin", controllers.SubmitAdmin)

	router.GET("/admin/home", middleware.AdminAuth, controllers.HomeAdmin)
	router.GET("/admin/logout", controllers.LogoutAdmin)

	router.GET("/admin/:deleteuser/:id", controllers.DeleteUserAdmin)
	router.GET("/admin/blockuser/:status/:id", controllers.BlockUserAdmin)

	//admin products
	router.GET("/admin/products", middleware.AdminAuth, controllers.ShowProductsAdmin)

	router.GET("/admin/products/:status/:pid", middleware.AdminAuth, controllers.BlockOrDeleteProductAdmin)

	//add produtct

	router.GET("/admin/addProduct", middleware.AdminAuth, controllers.AddProductGet)
	router.POST("/admin/addProduct/:from", controllers.AddProductPost)

	//edit product
	router.GET("/admin/editdProduct/:pid", middleware.AdminAuth, controllers.EditProductGet)
}
