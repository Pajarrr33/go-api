package main

import (
	"submission-project-enigma-laundry/config"
	"submission-project-enigma-laundry/controller"
	"submission-project-enigma-laundry/repository"
	"submission-project-enigma-laundry/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	// Connect into the database
	db := config.ConnectDb()

	defer db.Close()

	var (
		// Implement Dependency Injection
		// Repository
		customerRepository repository.CustomerRepository = repository.NewCustomerRepo(db)
		employeeRepository repository.EmployeeRepository = repository.NewEmployeeRepo(db)

		// Controller
		customerController controller.CustomerController = controller.NewCustomerController(customerRepository)
		employeeController controller.EmployeeController = controller.NewEmployeeController(employeeRepository)
	)

	server := gin.Default()

	// Routes
	routes.Customer(server,customerController)
	routes.Employee(server,employeeController)

	server.Run(":8080")
}