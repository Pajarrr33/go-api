package controller

import (
	"net/http"
	"strconv"
	"submission-project-enigma-laundry/entity"
	"submission-project-enigma-laundry/repository"
	"github.com/gin-gonic/gin"
)

type CustomerController interface {
	CreateCustomer(ctx *gin.Context)
	GetAllCustomer(ctx *gin.Context)
	GetDetailCustomer(ctx *gin.Context)
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

type CustomerResponseSlice struct {
	Message string `json:"message"`
	Data []entity.Customer
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

func (cc *customerController) GetAllCustomer(ctx *gin.Context) {
	customers := []entity.Customer{}
	rows, err := cc.CustomerRepository.GetCustomer()
	
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message" : "Failed to get all customer data", "details" : err.Error()})
		return
	}

	defer rows.Close()

	for rows.Next() {
		customer := entity.Customer{}
		err = rows.Scan(&customer.Customer_id,&customer.Name,&customer.Phone_number,&customer.Address)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message" : "Failed scanning customer data", "details" : err.Error()})
			return
		}
		customers = append(customers, customer)
	}

	err = rows.Err()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message" : "Error encountred during iteration", "details" : err.Error()})
		return
	}

	response := CustomerResponseSlice{
		Message: "Successfully get all data from customer",
		Data: customers,
	}

	ctx.JSON(http.StatusOK, response)
}

func (cc *customerController) GetDetailCustomer(ctx *gin.Context) {
	id := ctx.Param("id")
	
	convertedId,err := strconv.Atoi(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message" : "Failed convert id. Make sure id is number", "details" : err})
		return
	}

	customer := entity.Customer{}

	detailCustomer, err := cc.CustomerRepository.GetDetailCustomer(convertedId,&customer)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message" : "Failed to get detail customer data", "details" : err})
		return
	}

	// Create The response struct 
	response := CustomerResponse{
		Message: "Successfuly Get Customer Detail",
	}
	response.Data.Customer_id = detailCustomer.Customer_id
	response.Data.Name = detailCustomer.Name
	response.Data.Phone_number = detailCustomer.Phone_number
	response.Data.Address = detailCustomer.Address

	ctx.JSON(http.StatusCreated, response)
}