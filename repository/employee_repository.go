package repository

import (
	"database/sql"
	"errors"
	"submission-project-enigma-laundry/entity"
	_ "github.com/lib/pq"
)

type EmployeeRepository interface {
	GetEmployee() (*sql.Rows, error)
	GetDetailEmployee(id int,Employee *entity.Employee) (*entity.Employee,error)
	IsEmployeeExist(id int,Employee *entity.Employee) (bool,error)
	EmployeeInTransaction(id int,transaction *entity.Transaction)  (bool,error)
	CreateEmployee(Employee *entity.Employee) (*entity.Employee, error)	
	UpdateEmployee(id int,Employee *entity.Employee) (*entity.Employee,error)
	DeleteEmployee(id int) (bool,error)
}

type employeeRepository struct {
	DB *sql.DB
}

func NewEmployeeRepo(db *sql.DB) EmployeeRepository {
	return &employeeRepository{DB: db}
}

func (er *employeeRepository) IsEmployeeExist(id int, employee *entity.Employee) (bool, error) {
	query := "SELECT employee_id FROM employee WHERE employee_id = $1"
	
	// Execute the query and scan the result
	err := er.DB.QueryRow(query, id).Scan(&employee.Employee_id)
	if err != nil {
		if err == sql.ErrNoRows {
			// No employee found
			return false, nil // No error, just return false
		}
		// Return any other errors encountered
		return false, err
	}

	// Employee exists
	return true, nil
}

func (er *employeeRepository) EmployeeInTransaction(id int,transaction *entity.Transaction)  (bool,error) {
	query := "SELECT employee_id FROM transaction WHERE employee_id = $1"

	err := er.DB.QueryRow(query,id).Scan(&transaction.Employee_id)
	if err != nil {
		if err == sql.ErrNoRows {
			// No employee found
			return false, nil // No error, just return false
		}
		// Return any other errors encountered
		return false, err
	}

	// employee exists
	return true, nil
}


func (er *employeeRepository) CreateEmployee(employee *entity.Employee) (*entity.Employee, error) {
	// insert employee data into db
	insert_query := "INSERT INTO employee (name,phone_number,address) VALUES ($1, $2, $3) RETURNING employee_id;"

	err := er.DB.QueryRow(insert_query, employee.Name,employee.Phone_number,employee.Address).Scan(&employee.Employee_id)
	if err != nil {
		return employee, err // Handle error if the query fails
	}
	return employee, nil
}

func (er *employeeRepository) GetEmployee() (*sql.Rows, error) {
	// Get all data from customer table
	select_all := "SELECT employee_id,name,phone_number,address FROM employee;"

	rows,err := er.DB.Query(select_all)
	if err != nil {
		return rows,err
	}
	return rows,nil
}

func (er *employeeRepository) GetDetailEmployee(id int,employee *entity.Employee) (*entity.Employee,error) {
	select_by_id := "SELECT employee_id,name,phone_number,address FROM employee WHERE employee_id = $1"
	
	err := er.DB.QueryRow(select_by_id,id).Scan(&employee.Employee_id,&employee.Name,&employee.Phone_number,&employee.Address)
	if err != nil {
		if err == sql.ErrNoRows {
			err = errors.New("employee not found")
			return employee , err
		}

		return employee , err
	}

	return employee , nil
}

func (er *employeeRepository) UpdateEmployee(id int,employee *entity.Employee) (*entity.Employee,error) {
	update := "UPDATE employee SET name = $2,phone_number = $3,address = $4 WHERE employee_id = $1"

	_, err := er.DB.Exec(update,id,employee.Name,employee.Phone_number,employee.Address)
	if err != nil {
		return employee,err
	}
	return employee,nil
}

func (er *employeeRepository) DeleteEmployee(id int) (bool,error) {
	query := "DELETE FROM employee WHERE employee_id = $1"
	// Execute the query and scan the result
	_,err := er.DB.Exec(query,id)
	if err != nil {
		// Return any other errors encountered
		return false, err
	}

	return true, nil
}