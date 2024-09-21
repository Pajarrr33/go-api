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

		// Controller
		customerController controller.CustomerController = controller.NewCustomerController(customerRepository)
	)

	server := gin.Default()

	// Routes
	routes.Customer(server,customerController)

	server.Run(":8080")
}