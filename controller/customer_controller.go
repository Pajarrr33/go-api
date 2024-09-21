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

type CustomerResponse struct {
	Message string `json:"message"`
	Data  struct {
		Customer_id int `json:"id"`
		Name string `json:"name"`
		Phone_number string `json:"phoneNumber"`
		Address string `json:"address"`
	} `json:"data"`
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
	// Create The response struct 
	response := CustomerResponse{
		Message: "Successfuly Create Customer",
	}
	response.Data.Customer_id = createdCustomer.Customer_id
	response.Data.Name = createdCustomer.Name
	response.Data.Phone_number = createdCustomer.Phone_number
	response.Data.Address = createdCustomer.Address

	ctx.JSON(http.StatusCreated, response)
}