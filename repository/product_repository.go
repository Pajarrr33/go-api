package repository

import (
	"database/sql"
	"errors"
	"submission-project-enigma-laundry/entity"
	_ "github.com/lib/pq"
)

type ProductRepository interface {
	GetProduct() (*sql.Rows, error)
	GetProductByName(name string) (*sql.Rows, error)
	GetDetailProduct(id int, product *entity.Product) (*entity.Product, error)
	IsProductExist(id int, product *entity.Product) (bool, error)
	ProductInTransactionDetail(id int, transactionDetail *entity.Transaction_detail) (bool, error)
	CreateProduct(product *entity.Product) (*entity.Product, error)
	UpdateProduct(id int, product *entity.Product) (*entity.Product, error)
	DeleteProduct(id int) (bool, error)
}

type productRepository struct {
	DB *sql.DB
}

func NewProductRepo(db *sql.DB) ProductRepository {
	return &productRepository{DB: db}
}

func (pr *productRepository) IsProductExist(id int, product *entity.Product) (bool, error) {
	query := "SELECT product_id FROM product WHERE product_id = $1"

	// Execute the query and scan the result
	err := pr.DB.QueryRow(query, id).Scan(&product.Product_id)
	if err != nil {
		if err == sql.ErrNoRows {
			// No product found
			return false, nil // No error, just return false
		}
		// Return any other errors encountered
		return false, err
	}

	// product exists
	return true, nil
}

func (pr *productRepository) ProductInTransactionDetail(id int, transactionDetail *entity.Transaction_detail) (bool, error) {
	query := "SELECT product_id FROM transaction_detail WHERE product_id = $1"

	err := pr.DB.QueryRow(query, id).Scan(&transactionDetail.Product_id)
	if err != nil {
		if err == sql.ErrNoRows {
			// No product found
			return false, nil // No error, just return false
		}
		// Return any other errors encountered
		return false, err
	}

	// product exists
	return true, nil
}


func (pr *productRepository) CreateProduct(product *entity.Product) (*entity.Product, error) {
	// insert product data into db
	insert_query := "INSERT INTO product (product_name,unit,price) VALUES ($1, $2, $3) RETURNING product_id;"

	err := pr.DB.QueryRow(insert_query, product.Product_name, product.Unit, product.Price).Scan(&product.Product_id)
	if err != nil {
		return product, err // Handle error if the query fails
	}
	return product, nil
}

func (pr *productRepository) GetProduct() (*sql.Rows, error) {
	// Get all data from product table
	select_all := "SELECT product_id,product_name,unit,price FROM product;"

	rows, err := pr.DB.Query(select_all)
	if err != nil {
		return rows, err
	}
	return rows, nil
}

func (pr *productRepository) GetProductByName(name string) (*sql.Rows, error) {
	// Get all data from product table base on name
	query := "SELECT product_id,product_name,price,unit FROM product WHERE product_name LIKE $1;"

	rows, err := pr.DB.Query(query, name)
	if err != nil {
		return rows, err
	}
	return rows, nil
}

func (pr *productRepository) GetDetailProduct(id int, product *entity.Product) (*entity.Product, error) {
	select_by_id := "SELECT product_id,product_name,price,unit FROM product WHERE product_id = $1"

	err := pr.DB.QueryRow(select_by_id, id).Scan(&product.Product_id, &product.Product_name, &product.Price, &product.Unit)
	if err != nil {
		if err == sql.ErrNoRows {
			err = errors.New("product not found")
			return product, err
		}

		return product, err
	}

	return product, nil
}

func (pr *productRepository) UpdateProduct(id int, product *entity.Product) (*entity.Product, error) {
	update := "UPDATE product SET product_name = $2,unit = $3,price = $4 WHERE product_id = $1"

	_, err := pr.DB.Exec(update, id, product.Product_name, product.Unit, product.Price)
	if err != nil {
		return product, err
	}
	return product, nil
}

func (pr *productRepository) DeleteProduct(id int) (bool, error) {
	query := "DELETE FROM product WHERE product_id = $1"
	// Execute the query and scan the result
	_, err := pr.DB.Exec(query, id)
	if err != nil {
		// Return any other errors encountered
		return false, err
	}

	return true, nil
}
