package repository

import (
	"database/sql"
	_ "github.com/lib/pq"
	"submission-project-enigma-laundry/entity"
)

type CustomerRepository interface {
	CreateCustomer(customer *entity.Customer) (*entity.Customer, error)
}

type customerRepository struct {
	DB *sql.DB
}

func NewCustomerRepo(db *sql.DB) CustomerRepository {
	return &customerRepository{DB: db}
}

func (cr *customerRepository) CreateCustomer(customer *entity.Customer) (*entity.Customer,error) {
	// insert customer data into db
	insert_query := "INSERT INTO customer (name,phone_number,address) VALUES ($1, $2, $3) RETURNING customer_id;"

	err := cr.DB.QueryRow(insert_query, customer.Name, customer.Phone_number,customer.Address).Scan(&customer.Customer_id)
	if err != nil {
		return customer , err // Handle error if the query fails
	} 
	return customer , nil
}