package routes

import (
	"submission-project-enigma-laundry/controller"

	"github.com/gin-gonic/gin"
)


func Product(router *gin.Engine, pc controller.ProductController) {
	productRoutes := router.Group("/products")
	{
		productRoutes.GET("/",pc.ListProduct)
		productRoutes.GET("/:id",pc.GetDetailProduct)
		productRoutes.POST("/", pc.CreateProduct)
		productRoutes.PUT("/:id",pc.UpdateProduct)
		productRoutes.DELETE("/:id",pc.DeleteProduct)
	}
}