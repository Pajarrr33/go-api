package routes

import (
	"submission-project-enigma-laundry/controller"

	"github.com/gin-gonic/gin"
)


func Employee(router *gin.Engine, ec controller.EmployeeController) {
	employeeRoutes := router.Group("/employees")
	{
		employeeRoutes.GET("/",ec.GetAllEmployee)
		employeeRoutes.GET("/:id",ec.GetDetailEmployee)
		employeeRoutes.POST("/", ec.CreateEmployee)
		employeeRoutes.PUT("/:id",ec.UpdateEmployee)
		employeeRoutes.DELETE("/:id",ec.DeleteEmployee)
	}
}