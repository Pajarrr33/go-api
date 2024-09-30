package controller

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"submission-project-enigma-laundry/entity"
	"submission-project-enigma-laundry/repository"
	"time"

	"github.com/gin-gonic/gin"
)

type TransactionController interface {
	CreateTransaction(ctx *gin.Context)
	GetTransaction(ctx *gin.Context)
	ListTransaction(ctx *gin.Context)
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
		BillDetails []struct {
			Id             string `json:"id"`
			Transaction_id string    `json:"billId"`
			Product 	   entity.Product `json:"product"`
			Product_price  int    `json:"productPrice"`
			Qty            int    `json:"qty"`
		} `json:"billDetails"`
		Total_bill int `json:"totalBill"`
	} `json:"data"`
}

type TransactionResponseSlice struct {
	Message string `json:"message"`
	Data    []struct {
		Id         string `json:"id"`
		BillDate   string `json:"billDate"`
		EntryDate  string `json:"entryDate"`
		FinishDate string `json:"finishDate"`
		Employee   entity.Employee `json:"employee"`
		Customer   entity.Customer `json:"customer"`
		BillDetails []struct {
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
	response.Data.BillDetails = append(response.Data.BillDetails, struct{
		Id string "json:\"id\""; 
		Transaction_id string "json:\"billId\""; 
		Product entity.Product "json:\"product\""; 
		Product_price int "json:\"productPrice\""; 
		Qty int "json:\"qty\""
	}{
		Id: detailTransaction.Bill_detail[0].Transaction_detail_id,
		Transaction_id: detailTransaction.Bill_detail[0].Transaction_id,
		// Product
		Product: entity.Product{
			Product_id: detailTransaction.Bill_detail[0].Product.Product_id,
			Product_name: detailTransaction.Bill_detail[0].Product.Product_name,
			Price: detailTransaction.Bill_detail[0].Product.Price,
			Unit: detailTransaction.Bill_detail[0].Product.Unit,
		},
		Product_price: detailTransaction.Bill_detail[0].Product_price,
		Qty: detailTransaction.Bill_detail[0].Qty,
	})

	ctx.JSON(http.StatusOK, response)
}

func (tc *transactionController) ListTransaction(ctx *gin.Context) {
	transactions := []entity.Transaction{}

	var condition []string
	var conditionalQuery string
	paramsCount := 0

	startDate := ctx.Query("startDate")
	if startDate != "" {
		parsedDate,err := time.Parse("02-01-2006",startDate)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"message" : "Failed convert start date. Make sure start date format is dd-MM-yyyy", "details" : err.Error()})
			return
		}
		formatedDate := parsedDate.Format("2006-01-02")
		condition = append(condition,fmt.Sprintf(`t.entry_date >= '%s'`,formatedDate))
		paramsCount++
	}
	endDate := ctx.Query("endDate")
	if endDate != "" {
		parsedDate,err := time.Parse("02-01-2006",endDate)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"message" : "Failed convert end date. Make sure end date format is dd-MM-yyyy", "details" : err.Error()})
			return
		}
		formatedDate := parsedDate.Format("2006-01-02")
		condition = append(condition,fmt.Sprintf(`t.finish_date >= '%s'`,formatedDate))
		paramsCount++
	}

	productName := ctx.Query("productName")
	if productName != "" {
		condition = append(condition,fmt.Sprintf(`p.product_name LIKE '%s'`,productName))
		paramsCount++
	}

	if len(condition) > 0 {
		conditionalQuery = " WHERE " + strings.Join(condition, " AND ")
	}

	rows,err := tc.transactionRepository.ListTransaction(conditionalQuery)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message" : "Failed to get List Transaction", "details" : err.Error()})
		return
	}

	defer rows.Close()

	for rows.Next() {
		transaction := entity.Transaction{}
		err = rows.Scan(&transaction.Transaction_id,&transaction.Bill_date,&transaction.Entry_date,&transaction.Finish_date,&transaction.Employee.Employee_id,&transaction.Employee.Name,&transaction.Employee.Phone_number,&transaction.Employee.Address,&transaction.Customer.Customer_id,&transaction.Customer.Name,&transaction.Customer.Phone_number,&transaction.Customer.Address,&transaction.Bill_detail[0].Transaction_detail_id,&transaction.Bill_detail[0].Transaction_id,&transaction.Bill_detail[0].Product_price,&transaction.Bill_detail[0].Qty,&transaction.Bill_detail[0].Product.Product_id,&transaction.Bill_detail[0].Product.Product_name,&transaction.Bill_detail[0].Product.Price,&transaction.Bill_detail[0].Product.Unit,&transaction.Total_bill)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message" : "Failed scanning transaction", "details" : err.Error()})
			return
		}
		transactions = append(transactions, transaction)
	}

	err = rows.Err()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message" : "Error encountred during iteration", "details" : err.Error()})
		return
	}

	var response TransactionResponseSlice
	response.Message = "Successfuly Get Transaction"
	if response.Data == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message" : "transaction not found"})
		return
	}
	
	for _, transaction := range transactions {
		response.Data = append(response.Data, struct{
			Id string "json:\"id\""; 
			BillDate string "json:\"billDate\""; 
			EntryDate string "json:\"entryDate\""; 
			FinishDate string "json:\"finishDate\""; 
			Employee entity.Employee "json:\"employee\""; 
			Customer entity.Customer "json:\"customer\""; 
			BillDetails []struct{
				Id string "json:\"id\""; 
				Transaction_id string "json:\"billId\""; 
				Product entity.Product "json:\"product\""; 
				Product_price int "json:\"productPrice\""; 
				Qty int "json:\"qty\""
			} "json:\"billDetails\""; 
			Total_bill int "json:\"totalBill\""
		}{
			Id: transaction.Transaction_id,
			BillDate: transaction.Bill_date,
			EntryDate: transaction.Entry_date,
			FinishDate: transaction.Finish_date,
			Employee: entity.Employee{
				Employee_id: transaction.Employee.Employee_id,
				Name: transaction.Employee.Name,
				Phone_number: transaction.Employee.Phone_number,
				Address: transaction.Employee.Address,
			},
			Customer: entity.Customer{
				Customer_id: transaction.Customer.Customer_id,
				Name: transaction.Customer.Name,
				Phone_number: transaction.Customer.Phone_number,
				Address: transaction.Customer.Address,
			},
			BillDetails: []struct{
				Id string "json:\"id\""; 
				Transaction_id string "json:\"billId\""; 
				Product entity.Product "json:\"product\""; 
				Product_price int "json:\"productPrice\""; 
				Qty int "json:\"qty\""
			}{
				{
					Id: transaction.Bill_detail[0].Transaction_detail_id,
					Transaction_id: transaction.Bill_detail[0].Transaction_id,
					// Product
					Product: entity.Product{
						Product_id: transaction.Bill_detail[0].Product.Product_id,
						Product_name: transaction.Bill_detail[0].Product.Product_name,
						Price: transaction.Bill_detail[0].Product.Price,
						Unit: transaction.Bill_detail[0].Product.Unit,
					},
					Product_price: transaction.Bill_detail[0].Product_price,
					Qty: transaction.Bill_detail[0].Qty,
				},
			},
			Total_bill: transaction.Total_bill,
		})
	}

	ctx.JSON(http.StatusOK, response)
}