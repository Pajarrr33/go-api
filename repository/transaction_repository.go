package repository

import (
	"database/sql"
	"submission-project-enigma-laundry/entity"
	"fmt"
	"errors"
	_ "github.com/lib/pq"
)

type TransactionRepository interface {
	CreateTransaction(transaction *entity.Transaction) (*entity.Transaction,error)
	GetTransaction(transaction *entity.Transaction,id int) (*entity.Transaction,error)
	ListTransaction(transactionQueryParam string) (*sql.Rows, error)
	TransactionDetails(transactionDetailQueryParam string) (*sql.Rows, error)
	IsTransactionExist(id int)(bool, error) 
	IsTransactionDetailExist(id int) (bool, error)
}

type transactionRepository struct {
	DB *sql.DB
}

func NewTransactionRepo(db *sql.DB) TransactionRepository {
	return &transactionRepository{DB: db}
}

func (tr *transactionRepository) CreateTransaction(transaction *entity.Transaction) (*entity.Transaction,error) {
	tx,err := tr.DB.Begin()

	if err != nil {
		err = fmt.Errorf("failed starting transaction , %s",err)
		return transaction, err // Handle error if the query fails
	}
	
	createTransaction := "INSERT INTO transaction (customer_id,employee_id,bill_date,entry_date,finish_date) VALUES ($1,$2,$3,$4,$5) RETURNING transaction_id"

	err = tx.QueryRow(createTransaction, transaction.Customer_id, transaction.Employee_id, transaction.Bill_date, transaction.Entry_date, transaction.Finish_date).Scan(&transaction.Transaction_id)
	if err != nil {
		err = fmt.Errorf("failed insert into transaction , %s",err)
		tx.Rollback()
		return transaction, err // Handle error if the query fails
	}

	for i := range transaction.Bill_detail {
		billDetail := &transaction.Bill_detail[i] // Get pointer to the original element
	
		getPrice := "SELECT price FROM product WHERE product_id = $1;"
		err = tx.QueryRow(getPrice, billDetail.Product_id).Scan(&billDetail.Product_price)
		if err != nil {
			err = fmt.Errorf("failed to get price from product, %s", err)
			tx.Rollback()
			return transaction, err
		}
	
		createTransactionDetail := "INSERT INTO transaction_detail (transaction_id, product_id, product_price, qty) VALUES ($1, $2, $3, $4) RETURNING transaction_detail_id"
		err = tx.QueryRow(createTransactionDetail, transaction.Transaction_id, billDetail.Product_id, billDetail.Product_price, billDetail.Qty).Scan(&billDetail.Transaction_detail_id)
		if err != nil {
			err = fmt.Errorf("failed to insert into transaction detail, %s", err)
			tx.Rollback()
			return transaction, err
		}
	
		billDetail.Transaction_id = transaction.Transaction_id
	}	
	
	err = tx.Commit()
	if err != nil {
		err = fmt.Errorf("failed commit transaction , %s",err)
		return transaction, err // Handle error if the query fails
	}

	return transaction,nil
}

func (tr *transactionRepository) IsTransactionExist(id int) (bool, error) {
	query := "SELECT transaction_id FROM transaction WHERE transaction_id = $1"

	// Execute the query and scan the result
	err := tr.DB.QueryRow(query, id).Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			// No transaction found
			return false, nil // No error, just return false
		}
		// Return any other errors encountered
		return false, err
	}

	// transaction exists
	return true, nil
}

func (tr *transactionRepository) IsTransactionDetailExist(id int) (bool, error) {
	query := "SELECT transaction_id FROM transaction_detail WHERE transaction_id = $1"

	// Execute the query and scan the result
	err := tr.DB.QueryRow(query, id).Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			// No transaction Detail found
			return false, nil // No error, just return false
		}
		// Return any other errors encountered
		return false, err
	}

	// transaction detail exists
	return true, nil
}

func (tr *transactionRepository) GetTransaction(transaction *entity.Transaction,id int) (*entity.Transaction,error) {
	select_transaction_by_id := `SELECT 
	t.transaction_id,t.bill_date,t.entry_date,t.finish_date,
	e.employee_id,e.name,e.phone_number,e.address,
	c.customer_id,c.name,c.phone_number,c.address
	FROM transaction AS t 
	INNER JOIN employee AS e ON t.employee_id = e.employee_id 
	INNER JOIN customer AS c ON t.customer_id = c.customer_id 
	WHERE t.transaction_id = $1;`

	err := tr.DB.QueryRow(select_transaction_by_id,id).Scan(&transaction.Transaction_id,&transaction.Bill_date,&transaction.Entry_date,&transaction.Finish_date,&transaction.Employee.Employee_id,&transaction.Employee.Name,&transaction.Employee.Phone_number,&transaction.Employee.Address,&transaction.Customer.Customer_id,&transaction.Customer.Name,&transaction.Customer.Phone_number,&transaction.Customer.Address)
	if err != nil {
		if err == sql.ErrNoRows {
			err = errors.New("transaction not found")
			return transaction, err
		}

		return transaction, err
	}

	select_transaction_detail_by_transaction_id := `SELECT 
	td.transaction_detail_id,td.transaction_id,td.product_price,td.qty,
	p.product_id,p.product_name,p.price,p.unit
	FROM transaction_detail AS td
	INNER JOIN product AS p ON td.product_id = p.product_id WHERE transaction_id = $1;`

	rows, err := tr.DB.Query(select_transaction_detail_by_transaction_id,id)
	if err != nil {
		return transaction, err
	}
	
	defer rows.Close()

	for rows.Next() {
		transaction_detail := entity.Transaction_detail{}
		err = rows.Scan(&transaction_detail.Transaction_detail_id,&transaction_detail.Transaction_id,&transaction_detail.Product_price,&transaction_detail.Qty,&transaction_detail.Product.Product_id,&transaction_detail.Product.Product_name,&transaction_detail.Product.Price,&transaction_detail.Product.Unit)
		if err != nil {
			return transaction, err
		}
		transaction.Total_bill += transaction_detail.Product_price * transaction_detail.Qty
		transaction.Bill_detail = append(transaction.Bill_detail, transaction_detail)
	}

	return transaction, nil
}

func (tr *transactionRepository)ListTransaction(transactionQueryParam string) (*sql.Rows, error) {
	query := `
		SELECT t.transaction_id, t.bill_date, t.entry_date, t.finish_date, e.employee_id, e.name, e.phone_number, e.address,
		       c.customer_id, c.name, c.phone_number, c.address
		FROM transaction AS t
		INNER JOIN employee AS e ON t.employee_id = e.employee_id
		INNER JOIN customer AS c ON t.customer_id = c.customer_id
		INNER JOIN transaction_detail AS td ON t.transaction_id = td.transaction_id
		INNER JOIN product AS p ON td.product_id = p.product_id`

	if transactionQueryParam != ""{
		query += transactionQueryParam
	}

	rows, err := tr.DB.Query(query)
	if err != nil {
		return nil, err
	}

	return rows, nil
}

func (tr *transactionRepository) TransactionDetails(transactionDetailQueryParam string) (*sql.Rows, error) {
	query := `SELECT 
	td.transaction_detail_id,td.transaction_id,td.product_price,td.qty,
	p.product_id,p.product_name,p.price,p.unit
	FROM transaction_detail AS td
	INNER JOIN product AS p ON td.product_id = p.product_id`

	if transactionDetailQueryParam != "" {
		query += transactionDetailQueryParam
	}

	rows,err := tr.DB.Query(query)
	if err != nil {
		return nil,err
	}

	return rows,nil
}

