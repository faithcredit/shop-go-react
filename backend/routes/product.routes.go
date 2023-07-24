package routes

import (
	"gmrg/controllers"
	"gmrg/middleware"
	"gmrg/services"

	"github.com/gin-gonic/gin"
)

type ProductRouteController struct {
	productController controllers.ProductController
}

func NewProductControllerRoute(productController controllers.ProductController) ProductRouteController {
	return ProductRouteController{productController}
}

func (r *ProductRouteController) ProductRoute(rg *gin.RouterGroup, userService services.UserService) {
	protected := rg.Group("/auth/products")
	protected.Use(middleware.DeserializeUser(userService))
	router := rg.Group("/products")

	router.GET("/search", r.productController.FindProducts)
	router.GET("/:productId", r.productController.FindProductById)
	protected.POST("/", r.productController.CreateProduct)
	protected.PATCH("/:productId", r.productController.UpdateProduct)
	protected.DELETE("/:productId", r.productController.DeleteProduct)
}
