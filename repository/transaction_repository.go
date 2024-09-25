package repository

import (
	"database/sql"
	"submission-project-enigma-laundry/entity"
	"fmt"
	_ "github.com/lib/pq"
)

type TransactionRepository interface {
	CreateTransaction(transaction *entity.Transaction) (*entity.Transaction,error)
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
