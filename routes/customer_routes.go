package routes

import (
	"submission-project-enigma-laundry/controller"

	"github.com/gin-gonic/gin"
)


func Customer(router *gin.Engine, cc controller.CustomerController) {
	customerRoutes := router.Group("/customers")
	{
		customerRoutes.GET("/",cc.GetAllCustomer)
		customerRoutes.GET("/:id",cc.GetDetailCustomer)
		customerRoutes.POST("/", cc.CreateCustomer)
		customerRoutes.PUT("/:id",cc.UpdateCustomer)
		customerRoutes.DELETE("/:id",cc.DeleteCustomer)
	}
}