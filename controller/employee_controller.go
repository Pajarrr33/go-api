package controller

import (
	"net/http"
	"strconv"
	"strings"
	"submission-project-enigma-laundry/entity"
	"submission-project-enigma-laundry/repository"
	"github.com/gin-gonic/gin"
)

type EmployeeController interface {
	CreateEmployee(ctx *gin.Context)
	GetAllEmployee(ctx *gin.Context)
	GetDetailEmployee(ctx *gin.Context)
	UpdateEmployee(ctx *gin.Context)
	DeleteEmployee(ctx *gin.Context)
}

type employeeController struct {
	employeeRepository repository.EmployeeRepository
}

type EmployeeResponse struct {
	Message string `json:"message"`
	Data entity.Employee `json:"data"`
}

type EmployeeResponseSlice struct {
	Message string `json:"message"`
	Data []entity.Employee `json:"data"`
}

func NewEmployeeController(repo repository.EmployeeRepository) EmployeeController {
	return &employeeController{employeeRepository: repo}
}

func (ec *employeeController) CreateEmployee(ctx *gin.Context) {
	var newEmployee entity.Employee
	err := ctx.ShouldBind(&newEmployee)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message" : "Invalid Input","details" : err.Error()})
		return
	}

	createdEmployee,err := ec.employeeRepository.CreateEmployee(&newEmployee)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message" : "Failed to create employee", "details" : err.Error()})
		return
	}
	// Create The response struct 
	response := EmployeeResponse{
		Message: "Successfuly Create Employee",
		Data: *createdEmployee,
	}

	ctx.JSON(http.StatusCreated, response)
}

func (ec *employeeController) GetAllEmployee(ctx *gin.Context) {
	employees := []entity.Employee{}
	rows, err := ec.employeeRepository.GetEmployee()
	
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message" : "Failed to get all employee data", "details" : err.Error()})
		return
	}

	defer rows.Close()

	for rows.Next() {
		employee := entity.Employee{}
		err = rows.Scan(&employee.Employee_id,&employee.Name,&employee.Phone_number,&employee.Address)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message" : "Failed scanning employee data", "details" : err.Error()})
			return
		}
		employees = append(employees, employee)
	}

	err = rows.Err()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message" : "Error encountred during iteration", "details" : err.Error()})
		return
	}

	response :=EmployeeResponseSlice{
		Message: "Successfully get all data from employee",
		Data: employees,
	}

	ctx.JSON(http.StatusOK, response)
}

func (ec *employeeController) GetDetailEmployee(ctx *gin.Context) {
	id,err := strconv.Atoi(ctx.Param("id"))

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message" : "Failed convert id. Make sure id is number", "details" : err.Error()})
		return
	}

	employee := entity.Employee{}

	detailEmployee, err := ec.employeeRepository.GetDetailEmployee(id,&employee)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message" : "Failed to get detail employee data", "details" : err.Error()})
		return
	}

	// Create The response struct 
	response := EmployeeResponse{
		Message: "Successfuly Get Employee Detail",
		Data: *detailEmployee,
	}
	
	ctx.JSON(http.StatusOK, response)
}

func (ec *employeeController) UpdateEmployee(ctx *gin.Context) {
	convertedId, err := strconv.Atoi(ctx.Param("id"))
	
	if err != nil {
	  ctx.JSON(http.StatusBadRequest, gin.H{"message": "Failed convert id. Make sure id is number", "details": err.Error()})
	  return
	}
  
	employee := entity.Employee{}
  
	detailEmployee, err := ec.employeeRepository.GetDetailEmployee(convertedId, &employee)
	if err != nil {
	  ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to get employee data", "details": err.Error()})
	  return
	}

	updateEmployee := entity.Employee{}
  
	err = ctx.ShouldBind(&updateEmployee)
	if err != nil {
	  ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid Input", "details": err.Error()})
	  return
	}
  
	// Update existing employee data
	if strings.TrimSpace(updateEmployee.Name) != "" {
	  detailEmployee.Name = updateEmployee.Name
	}
	if strings.TrimSpace(updateEmployee.Phone_number) != "" {
	  detailEmployee.Phone_number = updateEmployee.Phone_number
	}
	if strings.TrimSpace(updateEmployee.Address) != "" {
	  detailEmployee.Address = updateEmployee.Address
	}
  
	updatedEmployee, err := ec.employeeRepository.UpdateEmployee(convertedId,detailEmployee) // Assuming updateEmployee function exists
	if err != nil {
	  ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to update employee data", "details": err.Error()})
	  return
	}
  
	// Create The response struct 
	response := EmployeeResponse{
	  Message: "Successfully Updated Employee Data",
	  Data: *updatedEmployee,
	}
  
	ctx.JSON(http.StatusOK, response) // Updated status code to 200 (OK)
  }

  func (ec *employeeController) DeleteEmployee(ctx *gin.Context) {
	convertedId, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
	  ctx.JSON(http.StatusBadRequest, gin.H{"message": "Failed convert id. Make sure id is number", "details": err.Error()})
	  return
	}

	employee := entity.Employee{}

	isEmployeeExist,err := ec.employeeRepository.IsEmployeeExist(convertedId,&employee)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message" : "Error While Checking Employee", "details" : err.Error()})
		return
	}
	if !isEmployeeExist {
		ctx.JSON(http.StatusNotFound, gin.H{"message" : "employee not found"})
		return
	}

	transaction := entity.Transaction{}

	isEmployeeInTransaction,err := ec.employeeRepository.EmployeeInTransaction(convertedId,&transaction)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError,gin.H{"message" : "Error While Checking Employee in Transaction", "details" : err.Error()})
		return
	} 
	if isEmployeeInTransaction {
		ctx.JSON(http.StatusConflict,gin.H{"message" : "Employee is being used in transaction. Please delete the transaction first"})
		return
	}

	_,err = ec.employeeRepository.DeleteEmployee(convertedId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message" : "Error Deleting Employee", "details" :err.Error()})
		return
	}

	response := struct {
		Message string `json:"message"`
		Data string `json:"data"`
	}{
		Message: "Successfully deleted data",
		Data: "OK",
	}

	ctx.JSON(http.StatusOK,response)
}