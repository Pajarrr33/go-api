package controller

import (
	"net/http"
	"strconv"
	"strings"
	"submission-project-enigma-laundry/entity"
	"submission-project-enigma-laundry/repository"
	"github.com/gin-gonic/gin"
)

type CustomerController interface {
	CreateCustomer(ctx *gin.Context)
	GetAllCustomer(ctx *gin.Context)
	GetDetailCustomer(ctx *gin.Context)
	UpdateCustomer(ctx *gin.Context)
	DeleteCustomer(ctx *gin.Context)
}

type CustomerResponse struct {
	Message string `json:"message"`
	Data  entity.Customer
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
		Data: *createdCustomer,
	}

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
		ctx.JSON(http.StatusBadRequest, gin.H{"message" : "Failed convert id. Make sure id is number", "details" : err.Error()})
		return
	}

	customer := entity.Customer{}

	detailCustomer, err := cc.CustomerRepository.GetDetailCustomer(convertedId,&customer)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message" : "Failed to get detail customer data", "details" : err.Error()})
		return
	}

	// Create The response struct 
	response := CustomerResponse{
		Message: "Successfuly Get Customer Detail",
		Data: *detailCustomer,
	}
	
	ctx.JSON(http.StatusOK, response)
}

func (cc *customerController) UpdateCustomer(ctx *gin.Context) {
	id := ctx.Param("id")
  
	convertedId, err := strconv.Atoi(id)
	if err != nil {
	  ctx.JSON(http.StatusBadRequest, gin.H{"message": "Failed convert id. Make sure id is number", "details": err.Error()})
	  return
	}
  
	customer := entity.Customer{}
  
	detailCustomer, err := cc.CustomerRepository.GetDetailCustomer(convertedId, &customer)
	if err != nil {
	  ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to get customer data", "details": err.Error()})
	  return
	}

	updateCustomer := entity.Customer{}
  
	err = ctx.ShouldBind(&updateCustomer)
	if err != nil {
	  ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid Input", "details": err.Error()})
	  return
	}
  
	// Update existing customer data
	if strings.TrimSpace(updateCustomer.Name) != "" {
	  detailCustomer.Name = updateCustomer.Name
	}
	if strings.TrimSpace(updateCustomer.Phone_number) != "" {
	  detailCustomer.Phone_number = updateCustomer.Phone_number
	}
	if strings.TrimSpace(updateCustomer.Address) != "" {
	  detailCustomer.Address = updateCustomer.Address
	}
  
	updatedCustomer, err := cc.CustomerRepository.UpdateCustomer(convertedId,detailCustomer) // Assuming UpdateCustomer function exists
	if err != nil {
	  ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to update customer data", "details": err.Error()})
	  return
	}
  
	// Create The response struct 
	response := CustomerResponse{
	  Message: "Successfully Updated Customer Data",
	  Data: *updatedCustomer,
	}
  
	ctx.JSON(http.StatusOK, response) // Updated status code to 200 (OK)
  }

func (cc *customerController) DeleteCustomer(ctx *gin.Context) {
	id := ctx.Param("id")
  
	convertedId, err := strconv.Atoi(id)
	if err != nil {
	  ctx.JSON(http.StatusBadRequest, gin.H{"message": "Failed convert id. Make sure id is number", "details": err.Error()})
	  return
	}

	customer := entity.Customer{}

	isCustomerExist,err := cc.CustomerRepository.IsCustomerExist(convertedId,&customer)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message" : "Error While Checking Customer", "details" : err.Error()})
		return
	}
	if !isCustomerExist {
		ctx.JSON(http.StatusNotFound, gin.H{"message" : "customer not found"})
		return
	}

	transaction := entity.Transaction{}

	isCustomerInTransaction,err := cc.CustomerRepository.CustomerInTransaction(convertedId,&transaction)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError,gin.H{"message" : "Error While Checking Customer in Transaction", "details" : err.Error()})
		return
	} 
	if isCustomerInTransaction {
		ctx.JSON(http.StatusConflict,gin.H{"message" : "Customer is being used in transaction. Please delete the transaction first"})
		return
	}

	_,err = cc.CustomerRepository.DeleteCustomer(convertedId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message" : "Error Deleting Customer", "details" :err.Error()})
		return
	}

	response := struct {
		Message string `json:"message"`
		Data string `json:"data"`
	}{
		Message: "Successfully deleted data",
		Data: "OK",
	}

	ctx.JSON(http.StatusOK,response)
}