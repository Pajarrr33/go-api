package repository

import (
	"database/sql"
	"errors"
	"submission-project-enigma-laundry/entity"
	_ "github.com/lib/pq"
)

type CustomerRepository interface {
	GetCustomer() (*sql.Rows, error)
	GetDetailCustomer(id int,customer *entity.Customer) (*entity.Customer,error)
	IsCustomerExist(id int,customer *entity.Customer) (bool,error)
	CustomerInTransaction(customerId int,transaction *entity.Transaction)  (bool,error)
	CreateCustomer(customer *entity.Customer) (*entity.Customer, error)	
	UpdateCustomer(id int,customer *entity.Customer) (*entity.Customer,error)
	DeleteCustomer(id int) (bool,error)
}

type customerRepository struct {
	DB *sql.DB
}

func NewCustomerRepo(db *sql.DB) CustomerRepository {
	return &customerRepository{DB: db}
}

func (cr *customerRepository) IsCustomerExist(id int, customer *entity.Customer) (bool, error) {
	query := "SELECT customer_id FROM customer WHERE customer_id = $1"
	
	// Execute the query and scan the result
	err := cr.DB.QueryRow(query, id).Scan(&customer.Customer_id)
	if err != nil {
		if err == sql.ErrNoRows {
			// No customer found
			return false, nil // No error, just return false
		}
		// Return any other errors encountered
		return false, err
	}

	// Customer exists
	return true, nil
}

func (cr *customerRepository) CustomerInTransaction(customerId int,transaction *entity.Transaction)  (bool,error) {
	query := "SELECT customer_id FROM transaction WHERE customer_id = $1"

	err := cr.DB.QueryRow(query,customerId).Scan(&transaction.Customer_id)
	if err != nil {
		if err == sql.ErrNoRows {
			// No customer found
			return false, nil // No error, just return false
		}
		// Return any other errors encountered
		return false, err
	}

	// Customer exists
	return true, nil
}


func (cr *customerRepository) CreateCustomer(customer *entity.Customer) (*entity.Customer, error) {
	// insert customer data into db
	insert_query := "INSERT INTO customer (name,phone_number,address) VALUES ($1, $2, $3) RETURNING customer_id;"

	err := cr.DB.QueryRow(insert_query, customer.Name, customer.Phone_number, customer.Address).Scan(&customer.Customer_id)
	if err != nil {
		return customer, err // Handle error if the query fails
	}
	return customer, nil
}

func (cr *customerRepository) GetCustomer() (*sql.Rows, error) {
	// Get all data from customer table
	select_all := "SELECT customer_id,name,phone_number,address FROM customer;"

	rows,err := cr.DB.Query(select_all)
	if err != nil {
		return rows,err
	}
	return rows,nil
}

func (cr *customerRepository) GetDetailCustomer(id int,customer *entity.Customer) (*entity.Customer,error) {
	select_by_id := "SELECT customer_id,name,phone_number,address FROM customer WHERE customer_id = $1"
	
	err := cr.DB.QueryRow(select_by_id,id).Scan(&customer.Customer_id,&customer.Name,&customer.Phone_number,&customer.Address)
	if err != nil {
		if err == sql.ErrNoRows {
			err = errors.New("customer not found")
			return customer , err
		}

		return customer , err
	}

	return customer , nil
}

func (cr *customerRepository) UpdateCustomer(id int,customer *entity.Customer) (*entity.Customer,error) {
	update := "UPDATE customer SET name = $2,phone_number = $3,address = $4 WHERE customer_id = $1"
	_, err := cr.DB.Exec(update,id,customer.Name,customer.Phone_number,customer.Address)
	if err != nil {
		return customer,err
	}
	return customer,nil
}

func (cr *customerRepository) DeleteCustomer(id int) (bool,error) {
	query := "DELETE FROM customer WHERE customer_id = $1"
	// Execute the query and scan the result
	_,err := cr.DB.Exec(query,id)
	if err != nil {
		// Return any other errors encountered
		return false, err
	}

	// Customer exists
	return true, nil
}