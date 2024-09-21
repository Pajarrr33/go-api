package routes

import (
	"submission-project-enigma-laundry/controller"

	"github.com/gin-gonic/gin"
)


func Customer(router *gin.Engine, cc controller.CustomerController) {
	customerRoutes := router.Group("/customers")
	{
		customerRoutes.POST("/", cc.CreateCustomer)
	}
}