package controller

import (
	"fmt"
	"net/http"
	"strconv"
	"submission-project-enigma-laundry/entity"
	"submission-project-enigma-laundry/repository"
	"regexp"
	"errors"
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
		BillDetails []struct {
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

	for _, billDetail := range newTransaction.Bill_detail {
		converIdProduct,err := strconv.Atoi(billDetail.Product_id)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "Failed convert product id. Make sure product id is number", "details": err.Error()})
			return
		}

		isProductExist,err := tc.productRepository.IsProductExist(converIdProduct,&billDetail.Product)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message" : "Error While Checking Product", "details" : err.Error()})
			return
		}
		if !isProductExist {
			ctx.JSON(http.StatusNotFound, gin.H{"message" : "product not found"})
			return
		}
	}

	isDateValid,err := isValidDate(newTransaction.Entry_date,newTransaction.Finish_date,newTransaction.Bill_date)
	if !isDateValid && err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "entryDate, finishDate, billDate format is wrong", "details": err.Error()})
		return
	}

	if newTransaction.Bill_date != newTransaction.Entry_date {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Bill date in future"})
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
	for _, billDetail := range createdTransaction.Bill_detail {
		response.Data.BillDetails = append(response.Data.BillDetails, struct{
			Id string "json:\"id\""; 
			Transaction_id string "json:\"billId\""; 
			Product_id string "json:\"productId\""; 
			Product_price int "json:\"productPrice\""; 
			Qty int "json:\"qty\""
		}{
			Id: billDetail.Transaction_detail_id,
			Transaction_id: billDetail.Transaction_id,
			Product_id: billDetail.Product_id,
			Product_price: billDetail.Product_price,
			Qty: billDetail.Qty,
		})
	}
	
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
	// Insert data into the nested BillDetails struct
	for _, billDetail := range detailTransaction.Bill_detail {
		response.Data.BillDetails = append(response.Data.BillDetails, struct{
			Id string "json:\"id\""; 
			Transaction_id string "json:\"billId\""; 
			Product entity.Product "json:\"product\""; 
			Product_price int "json:\"productPrice\""; 
			Qty int "json:\"qty\""
		}{
			Id: billDetail.Transaction_detail_id,
			Transaction_id: billDetail.Transaction_id,
			// Product
			Product: entity.Product{
				Product_id: billDetail.Product.Product_id,
				Product_name: billDetail.Product.Product_name,
				Price: billDetail.Product.Price,
				Unit: billDetail.Product.Unit,
			},
			Product_price: billDetail.Product_price,
			Qty: billDetail.Qty,
		})
	}
	
	ctx.JSON(http.StatusOK, response)
}

func (tc *transactionController) ListTransaction(ctx *gin.Context) {
	transactions := []entity.Transaction{}

	var transactionQueryParam, transactionDetailQueryParam string

    // Build the product query if productName is provided
	startDate := ctx.Query("startDate")
	endDate := ctx.Query("endDate")
    productName := ctx.Query("productName")
	// Build the transaction date query
    if startDate != "" || endDate != "" || productName != ""{
        transactionQueryParam = " WHERE " + buildDateQuery(startDate, endDate, productName)
    }
    if productName != "" {
        transactionDetailQueryParam = " WHERE " + buildProductQuery(productName)
    }
	
	rows,err := tc.transactionRepository.ListTransaction(transactionQueryParam)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message" : "Failed to get List Transaction", "details" : err.Error()})
		return
	}

	defer rows.Close()

	for rows.Next() {
		transaction := entity.Transaction{}
		err = rows.Scan(&transaction.Transaction_id,&transaction.Bill_date,&transaction.Entry_date,&transaction.Finish_date,&transaction.Employee.Employee_id,&transaction.Employee.Name,&transaction.Employee.Phone_number,&transaction.Employee.Address,&transaction.Customer.Customer_id,&transaction.Customer.Name,&transaction.Customer.Phone_number,&transaction.Customer.Address)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message" : "Failed scanning transaction", "details" : err.Error()})
			return
		}
		transactions = append(transactions, transaction)
	}

	err = rows.Err()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message" : "Error encountred during iteration transaction", "details" : err.Error()})
		return
	}

	rows,err = tc.transactionRepository.TransactionDetails(transactionDetailQueryParam)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message" : "Failed to get List Transaction detail", "details" : err.Error()})
		return
	}

	transaction_details := []entity.Transaction_detail{}
	for rows.Next() {
		transaction_detail := entity.Transaction_detail{}
		err = rows.Scan(&transaction_detail.Transaction_detail_id,&transaction_detail.Transaction_id,&transaction_detail.Product_price,&transaction_detail.Qty,&transaction_detail.Product.Product_id,&transaction_detail.Product.Product_name,&transaction_detail.Product.Price,&transaction_detail.Product.Unit)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message" : "Failed scanning transaction detail", "details" : err.Error()})
			return
		}
		transaction_details = append(transaction_details, transaction_detail)
	}
	err = rows.Err()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message" : "Error encountred during iteration transaction", "details" : err.Error()})
		return
	}


	for i := 0; i < len(transactions); i++ {
		for j := 0; j < len(transaction_details); j++ {
			// Check if the transaction_id matches
			if transactions[i].Transaction_id == transaction_details[j].Transaction_id {
				// Update the total bill
				transactions[i].Total_bill += transaction_details[j].Product_price * transaction_details[j].Qty
				
				// Append the matching transaction detail to the Bill_detail slice
				transactions[i].Bill_detail = append(transactions[i].Bill_detail, transaction_details[j])
			}
		}
	}

	var response TransactionResponseSlice
	response.Message = "Successfully Get Transaction"

	for _, transaction := range transactions {
		// Reinitialize billDetails for each transaction to ensure it's empty
		var billDetails []struct {
			Id             string        `json:"id"`
			Transaction_id string        `json:"billId"`
			Product        entity.Product `json:"product"`
			Product_price  int           `json:"productPrice"`
			Qty            int           `json:"qty"`
		}

		// Loop through the bill details of the current transaction
		for _, detail := range transaction.Bill_detail {
			billDetails = append(billDetails, struct {
				Id             string        `json:"id"`
				Transaction_id string        `json:"billId"`
				Product        entity.Product `json:"product"`
				Product_price  int           `json:"productPrice"`
				Qty            int           `json:"qty"`
			}{
				Id:             detail.Transaction_detail_id,
				Transaction_id: detail.Transaction_id,
				Product: entity.Product{
					Product_id:   detail.Product.Product_id,
					Product_name: detail.Product.Product_name,
					Price:        detail.Product.Price,
					Unit:         detail.Product.Unit,
				},
				Product_price: detail.Product_price,
				Qty:           detail.Qty,
			})
		}

		// Append the transaction and its bill details to the response
		response.Data = append(response.Data, struct {
			Id             string        `json:"id"`
			BillDate       string        `json:"billDate"`
			EntryDate      string        `json:"entryDate"`
			FinishDate     string        `json:"finishDate"`
			Employee       entity.Employee `json:"employee"`
			Customer       entity.Customer `json:"customer"`
			BillDetails    []struct {
				Id             string        `json:"id"`
				Transaction_id string        `json:"billId"`
				Product        entity.Product `json:"product"`
				Product_price  int           `json:"productPrice"`
				Qty            int           `json:"qty"`
			} `json:"billDetails"`
			Total_bill int `json:"totalBill"`
		}{
			Id: transaction.Transaction_id,
			BillDate: transaction.Bill_date,
			EntryDate: transaction.Entry_date,
			FinishDate: transaction.Finish_date,
			Employee: entity.Employee{
				Employee_id:   transaction.Employee.Employee_id,
				Name:          transaction.Employee.Name,
				Phone_number:  transaction.Employee.Phone_number,
				Address:       transaction.Employee.Address,
			},
			Customer: entity.Customer{
				Customer_id:   transaction.Customer.Customer_id,
				Name:          transaction.Customer.Name,
				Phone_number:  transaction.Customer.Phone_number,
				Address:       transaction.Customer.Address,
			},
			BillDetails: billDetails, // Properly populated for each transaction
			Total_bill:  transaction.Total_bill,
		})
	}

	// If no data is found, return a "not found" response
	if len(response.Data) == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "transaction not found"})
		return
	}

	// Return the populated response
	ctx.JSON(http.StatusOK, response)
}

func isValidDate(dates ...string) (bool,error) {
	// Regular expression for dd/MM/yyyy format
	var datePattern = `^(0[1-9]|[12][0-9]|3[01])-(0[1-9]|1[0-2])-\d{4}$`

	// Compile the regular expression
	re := regexp.MustCompile(datePattern)

	// Iterate over each date and check if it matches the pattern
	for _, date := range dates {
		if !re.MatchString(date) {
			err := errors.New("the entryDate / finishDate / billDate format is invalid make sure the date format is dd-mm-yyyy")
			return false,err
		}
	}

	// Return true if the date matches the pattern, false otherwise
	return true,nil
}


func buildDateQuery(startDate, endDate, productName string) string {
    query := ""
    if startDate != "" {
		if query != "" {
            query += " AND"
        }
        query += fmt.Sprintf(" TO_DATE(t.entry_date, 'DD-MM-YY') >= TO_DATE('%s', 'DD-MM-YY')", startDate)
    }
    if endDate != "" {
        if query != "" {
            query += " AND"
        }
        query += fmt.Sprintf(" TO_DATE(t.finish_date, 'DD-MM-YY') <= TO_DATE('%s', 'DD-MM-YY')", endDate)
    }
	if productName != "" {
		if query != "" {
            query += " AND"
        }
		query += fmt.Sprintf(" p.product_name LIKE '%%%s%%'", productName)
	}
    return query
}

func buildProductQuery(productName string) string {
    if productName != "" {
        return fmt.Sprintf(" p.product_name LIKE '%%%s%%'", productName)
    }
    return ""
}