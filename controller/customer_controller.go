package controller

import (
	"net/http"
	"submission-project-enigma-laundry/entity"
	"submission-project-enigma-laundry/repository"
	"github.com/gin-gonic/gin"
)

type CustomerController interface {
	CreateCustomer(ctx *gin.Context)
}

type customerController struct {
	CustomerRepository repository.CustomerRepository
}

func NewCustomerController(repo repository.CustomerRepository) CustomerController {
	return &customerController{CustomerRepository: repo}
}

func (cc *customerController) CreateCustomer(ctx *gin.Context) {
	var newCustomer entity.Customer
	err := ctx.ShouldBind(&newCustomer)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message" : "Invalid Input","details" : err.Error()})
		return
	}

	createdCustomer,err := cc.CustomerRepository.CreateCustomer(&newCustomer)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message" : "Failed to create customer", "details" : err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"data" : createdCustomer})
}