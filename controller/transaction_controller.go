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
	GetTransaction(ctx *gin.Context)
}

type CreatedTransactionResponse struct {
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

type TransactionResponse struct {
	Message string `json:"message"`
	Data    struct {
		Id         string `json:"id"`
		BillDate   string `json:"billDate"`
		EntryDate  string `json:"entryDate"`
		FinishDate string `json:"finishDate"`
		Employee   entity.Employee `json:"employee"`
		Customer   entity.Customer `json:"customer"`
		BillDetails [1]struct {
			Id             string `json:"id"`
			Transaction_id string    `json:"billId"`
			Product 	   entity.Product `json:"product"`
			Product_price  int    `json:"productPrice"`
			Qty            int    `json:"qty"`
		} `json:"billDetails"`
		Total_bill int `json:"totalBill"`
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
	var response CreatedTransactionResponse
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

func (tc *transactionController) GetTransaction(ctx *gin.Context) {
	id,err := strconv.Atoi(ctx.Param("id_bill"))

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message" : "Failed convert id bill. Make sure id bill is number", "details" : err.Error()})
		return
	}

	isTransactionExist,err := tc.transactionRepository.IsTransactionExist(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message" : "Error While Checking Transaction", "details" : err.Error()})
		return
	}
	if !isTransactionExist {
		ctx.JSON(http.StatusNotFound, gin.H{"message" : "transaction not found"})
		return
	}

	isTransactionDetailExist,err := tc.transactionRepository.IsTransactionDetailExist(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message" : "Error While Checking Transaction Detail", "details" : err.Error()})
		return
	}
	if !isTransactionDetailExist {
		ctx.JSON(http.StatusNotFound, gin.H{"message" : "transaction detail not found"})
		return
	}

	transaction := entity.Transaction{}
	detailTransaction,err := tc.transactionRepository.GetTransaction(&transaction,id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message" : "Failed to get transaction data", "details" : err.Error()})
		return
	}

	// response struct 
	var response TransactionResponse
	response.Message = "Successfuly Get Transaction"

	// Data Response
	// Transaction
	response.Data.Id = detailTransaction.Transaction_id
	response.Data.BillDate = detailTransaction.Bill_date
	response.Data.EntryDate = detailTransaction.Entry_date
	response.Data.FinishDate = detailTransaction.Finish_date
	response.Data.Total_bill = detailTransaction.Total_bill

	// Employee
	response.Data.Employee.Employee_id = detailTransaction.Employee.Employee_id
	response.Data.Employee.Name = detailTransaction.Employee.Name
	response.Data.Employee.Phone_number = detailTransaction.Employee.Phone_number
	response.Data.Employee.Address = detailTransaction.Employee.Address

	// Customer
	response.Data.Customer.Customer_id = detailTransaction.Customer.Customer_id
	response.Data.Customer.Name = detailTransaction.Customer.Name
	response.Data.Customer.Phone_number = detailTransaction.Customer.Phone_number
	response.Data.Customer.Address = detailTransaction.Customer.Address

	// Bill Details
	response.Data.BillDetails[0].Id = detailTransaction.Bill_detail[0].Transaction_detail_id
	response.Data.BillDetails[0].Transaction_id = detailTransaction.Bill_detail[0].Transaction_id
	response.Data.BillDetails[0].Product_price = detailTransaction.Bill_detail[0].Product_price
	response.Data.BillDetails[0].Qty = detailTransaction.Bill_detail[0].Qty

	// Product
	response.Data.BillDetails[0].Product.Product_id = detailTransaction.Bill_detail[0].Product.Product_id
	response.Data.BillDetails[0].Product.Product_name = detailTransaction.Bill_detail[0].Product.Product_name
	response.Data.BillDetails[0].Product.Price = detailTransaction.Bill_detail[0].Product.Price
	response.Data.BillDetails[0].Product.Unit = detailTransaction.Bill_detail[0].Product.Unit

	ctx.JSON(http.StatusOK, response)
}