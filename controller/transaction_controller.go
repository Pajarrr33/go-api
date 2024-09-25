package controller

import (
	"submission-project-enigma-laundry/entity"
	"submission-project-enigma-laundry/repository"
	"net/http"
	"strconv"
	"github.com/gin-gonic/gin"
)

type TransactionController interface {
	CreateTransaction(ctx *gin.Context)
}

type TransactionResponse struct {
	Message string `json:"message"`
	Data    struct {
		Id         string `json:"id"`
		BillDate   string `json:"billDate"`
		EntryDate  string `json:"entryDate"`
		FinishDate string `json:"finishDate"`
		EmployeeId string `json:"employeeId"`
		CustomerId string `json:"customerId"`
		BillDetails struct {
			Id             string `json:"id"`
			Transaction_id string    `json:"billId"`
			Product_id     string `json:"productId"`
			Product_price  int    `json:"productPrice"`
			Qty            int    `json:"qty"`
		} `json:"billDetails"`
	} `json:"data"`
}

type transactionController struct {
	customerRepository 		repository.CustomerRepository
	employeeRepository 		repository.EmployeeRepository
	productRepository 		repository.ProductRepository
	transactionRepository 	repository.TransactionRepository
}

func NewTransactionController(cr repository.CustomerRepository,er repository.EmployeeRepository,pr repository.ProductRepository,tr repository.TransactionRepository) TransactionController {
	return &transactionController{customerRepository: cr,employeeRepository: er,productRepository: pr,transactionRepository: tr}
}

func (tc *transactionController) CreateTransaction(ctx *gin.Context) {
	var newTransaction entity.Transaction
	err := ctx.ShouldBind(&newTransaction)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message" : "Invalid Input","details" : err.Error()})
		return
	}

	converIdCustomer,err := strconv.Atoi(newTransaction.Customer_id) 
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Failed convert customer id. Make sure customer id is number", "details": err.Error()})
		return
	}

	isCustomerExist,err := tc.customerRepository.IsCustomerExist(converIdCustomer,&newTransaction.Customer)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message" : "Error While Checking Customer", "details" : err.Error()})
		return
	}
	if !isCustomerExist {
		ctx.JSON(http.StatusNotFound, gin.H{"message" : "customer not found"})
		return
	}

	converIdEmployee,err := strconv.Atoi(newTransaction.Employee_id) 
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Failed convert employee id. Make sure employee id is number", "details": err.Error()})
		return
	}

	isEmployeeExist,err := tc.employeeRepository.IsEmployeeExist(converIdEmployee,&newTransaction.Employee)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message" : "Error While Checking Employee", "details" : err.Error()})
		return
	}
	if !isEmployeeExist {
		ctx.JSON(http.StatusNotFound, gin.H{"message" : "employee not found"})
		return
	}

	converIdProduct,err := strconv.Atoi(newTransaction.Bill_detail[0].Product_id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Failed convert product id. Make sure product id is number", "details": err.Error()})
		return
	}

	isProductExist,err := tc.productRepository.IsProductExist(converIdProduct,&newTransaction.Bill_detail[0].Product)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message" : "Error While Checking Product", "details" : err.Error()})
		return
	}
	if !isProductExist {
		ctx.JSON(http.StatusNotFound, gin.H{"message" : "product not found"})
		return
	}

	createdTransaction,err := tc.transactionRepository.CreateTransaction(&newTransaction) 
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message" : "Failed to create transaction", "details" : err.Error()})
		return
	}

	// Create The response struct 
	var response TransactionResponse
	response.Message = "Successfuly Create Transaction"
	response.Data.Id = createdTransaction.Transaction_id
	response.Data.BillDate = createdTransaction.Bill_date
	response.Data.EntryDate = createdTransaction.Entry_date
	response.Data.FinishDate = createdTransaction.Finish_date
	response.Data.EmployeeId = createdTransaction.Employee_id
	response.Data.CustomerId = createdTransaction.Customer_id

	// Insert data into the nested BillDetails struct
	response.Data.BillDetails.Id = createdTransaction.Bill_detail[0].Transaction_detail_id
	response.Data.BillDetails.Transaction_id = createdTransaction.Bill_detail[0].Transaction_id
	response.Data.BillDetails.Product_id = createdTransaction.Bill_detail[0].Product_id
	response.Data.BillDetails.Product_price = createdTransaction.Bill_detail[0].Product_price
	response.Data.BillDetails.Qty = createdTransaction.Bill_detail[0].Qty
	ctx.JSON(http.StatusCreated, response)
}