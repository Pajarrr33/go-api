package config

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func ConnectDb() *sql.DB {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	// Database connection constants
	host     := os.Getenv("DB_HOST")  		    	
	port,err := strconv.Atoi(os.Getenv("DB_PORT"))  
	if err != nil {
		panic(err)
	}
	user     := os.Getenv("DB_USERNAME")  			
	password := os.Getenv("DB_PASSWORD")   			
	dbname   := os.Getenv("DB_DATABASE")        	
	dbconnection := os.Getenv("DB_CONNECTION") 

	// Connection string for PostgreSQL
	var psqlInfo = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
    				host, port, user, password, dbname)

	db, err := sql.Open(dbconnection, psqlInfo)  // Open a connection using the connection string

	if err != nil {
		panic(err)  // Handle error if connection fails
	}

	// Test the database connection
	err = db.Ping()
	if err != nil {
		panic(err)  // Handle error if ping fails
	}

	return db
}