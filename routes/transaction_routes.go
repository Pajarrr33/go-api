package routes

import (
	"submission-project-enigma-laundry/controller"

	"github.com/gin-gonic/gin"
)


func Transaction(router *gin.Engine, tc controller.TransactionController) {
	transactionRoutes := router.Group("/transactions")
	{
		transactionRoutes.POST("/",tc.CreateTransaction)
	}
}