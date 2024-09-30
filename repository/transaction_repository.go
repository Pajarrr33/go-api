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
	ListTransaction(addtionalQuery string) (*sql.Rows, error)
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
	
	getPrice := "SELECT price FROM product WHERE product_id = $1;"
	
	err = tx.QueryRow(getPrice,transaction.Bill_detail[0].Product_id).Scan(&transaction.Bill_detail[0].Product_price)
	if err != nil {
		err = fmt.Errorf("failed get price from product , %s",err)
		tx.Rollback()
		return transaction, err // Handle error if the query fails
	}

	createTransactionDetail := "INSERT INTO transaction_detail (transaction_id,product_id,product_price,qty) VALUES ($1,$2,$3,$4) RETURNING transaction_detail_id"

	err = tx.QueryRow(createTransactionDetail,transaction.Transaction_id,transaction.Bill_detail[0].Product_id,transaction.Bill_detail[0].Product_price,transaction.Bill_detail[0].Qty).Scan(&transaction.Bill_detail[0].Transaction_detail_id)
	if err != nil {
		err = fmt.Errorf("failed insert Into transaction detail , %s",err)
		tx.Rollback()
		return transaction, err // Handle error if the query fails
	}

	transaction.Bill_detail[0].Transaction_id = transaction.Transaction_id
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
	query := "SELECT t.transaction_id,t.bill_date,t.entry_date,t.finish_date,e.employee_id,e.name,e.phone_number,e.address,c.customer_id,c.name,c.phone_number,c.address,td.transaction_detail_id,td.transaction_id,td.product_price,td.qty,p.product_id,p.product_name,p.price,p.unit,p.price * td.qty AS total_bill FROM transaction AS t INNER JOIN employee AS e ON t.employee_id = e.employee_id INNER JOIN customer AS c ON t.customer_id = c.customer_id INNER JOIN transaction_detail AS td ON t.transaction_id = td.transaction_id INNER JOIN product AS p ON td.product_id = p.product_id WHERE t.transaction_id = $1;"

	err := tr.DB.QueryRow(query,id).Scan(&transaction.Transaction_id,&transaction.Bill_date,&transaction.Entry_date,&transaction.Finish_date,&transaction.Employee.Employee_id,&transaction.Employee.Name,&transaction.Employee.Phone_number,&transaction.Employee.Address,&transaction.Customer.Customer_id,&transaction.Customer.Name,&transaction.Customer.Phone_number,&transaction.Customer.Address,&transaction.Bill_detail[0].Transaction_detail_id,&transaction.Bill_detail[0].Transaction_id,&transaction.Bill_detail[0].Product_price,&transaction.Bill_detail[0].Qty,&transaction.Bill_detail[0].Product.Product_id,&transaction.Bill_detail[0].Product.Product_name,&transaction.Bill_detail[0].Product.Price,&transaction.Bill_detail[0].Product.Unit,&transaction.Total_bill)
	if err != nil {
		if err == sql.ErrNoRows {
			err = errors.New("transaction not found")
			return transaction, err
		}

		return transaction, err
	}

	return transaction, nil
}

func (tr *transactionRepository)ListTransaction(addtionalQuery string) (*sql.Rows, error) {
	query := `
		SELECT t.transaction_id, t.bill_date, t.entry_date, t.finish_date, e.employee_id, e.name, e.phone_number, e.address,
		       c.customer_id, c.name, c.phone_number, c.address, td.transaction_detail_id, td.transaction_id, 
		       td.product_price, td.qty, p.product_id, p.product_name, p.price, p.unit, p.price * td.qty AS total_bill
		FROM transaction AS t
		INNER JOIN employee AS e ON t.employee_id = e.employee_id
		INNER JOIN customer AS c ON t.customer_id = c.customer_id
		INNER JOIN transaction_detail AS td ON t.transaction_id = td.transaction_id
		INNER JOIN product AS p ON td.product_id = p.product_id
	`
	if addtionalQuery != ""{
		query += addtionalQuery
		rows, err := tr.DB.Query(query)
			if err != nil {
				return rows, err
			}
		return rows, nil
	}
	rows, err := tr.DB.Query(query)
		if err != nil {
			return rows, err
		}
	return rows, nil
}

