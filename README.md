
# Golang API Challenge
This project is a simple Go-based API for managing a laundry shop. It provides CRUD (Create, Read, Update, Delete) operations for customers,employee,product, and transaction and includes validations for data consistency.


## Tech Stack

- Golang https://github.com/golang/go
- Postgree sql https://www.postgresql.org/
- Gin (HTTP Framework) https://github.com/gin-gonic/gin
- Godotenv (enviroment) https://github.com/johogodotenv

## Installation
1. Clone this repository

```bash
	https://git.enigmacamp.com/enigma-20/ahmad-fajar-shidik/challenge-goapi.git
```

2. Create a database 

```bash
	CREATE DATABASE example_name
```

3. Run this DDL query or copy it from DDL.sql File

```bash
	CREATE TABLE customer (
        customer_id SERIAL PRIMARY KEY,
        name VARCHAR(255) NOT NULL,
        phone_number VARCHAR(255) NOT NULL,
        address VARCHAR(255) DEFAULT '',
        created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
    );


    CREATE TABLE employee (
        employee_id SERIAL PRIMARY KEY,
        name VARCHAR(255) NOT NULL,
        phone_number VARCHAR(255) NOT NULL,
        address VARCHAR(255) DEFAULT '',
        created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
    );

    CREATE TABLE product (
        product_id SERIAL PRIMARY KEY,
        product_name VARCHAR(255) NOT NULL,
        unit VARCHAR(255) NOT NULL,
        price INT NOT NULL,
        created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
    );

    CREATE TABLE transaction (
        transaction_id SERIAL PRIMARY KEY,
        customer_id INT NOT NULL,
        employee_id INT NOT NULL,
        bill_date VARCHAR(255),
        entry_date VARCHAR(255),
        finish_date VARCHAR(255),
        created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
        FOREIGN KEY (customer_id) REFERENCES customer(customer_id),
        FOREIGN KEY (employee_id) REFERENCES employee(employee_id)
    );

    CREATE TABLE transaction_detail (
        transaction_detail_id SERIAL PRIMARY KEY,
        transaction_id INT NOT NULL,
        product_id INT NOT NULL,
        product_price INT NOT NULL,
        qty INT NOT NULL,
        FOREIGN KEY (transaction_id) REFERENCES transaction(transaction_id),
        FOREIGN KEY (product_id) REFERENCES product(product_id)
    );
```

4. Run this DML query or copy it from DML.sql File
```bash
	INSERT INTO customer (name, phone_number, address)
    VALUES
    ('John Doe', '555-1234', '123 Elm St'),
    ('Jane Smith', '555-5678', '456 Oak St'),
    ('Michael Johnson', '555-8765', '789 Pine St'),
    ('Emily Davis', '555-4321', '321 Maple Ave'),
    ('Robert Brown', '555-1111', '654 Cedar St');

    INSERT INTO employee (name, phone_number, address)
    VALUES
    ('Alice Williams', '555-2222', '987 Willow St'),
    ('David Harris', '555-3333', '123 Birch St'),
    ('Sophia Martinez', '555-4444', '456 Redwood St'),
    ('James Wilson', '555-5555', '789 Palm St'),
    ('Olivia Garcia', '555-6666', '321 Cypress Ave');

    INSERT INTO product (product_name, unit, price)
    VALUES
    ('Shampoo', 'bottle', 10000),
    ('Soap', 'bar', 5000),
    ('Toothpaste', 'tube', 15000),
    ('Conditioner', 'bottle', 12000),
    ('Body Lotion', 'bottle', 25000);

    INSERT INTO transaction (customer_id, employee_id, bill_date, entry_date, finish_date) 
    VALUES 
    (1, 1, '01-10-2024', '01-10-2024', '05-10-2024'),
    (2, 2, '02-10-2024', '02-10-2024', '06-10-2024'),
    (3, 3, '03-10-2024', '03-10-2024', '07-10-2024'),
    (4, 4, '04-10-2024', '04-10-2024', '08-10-2024'),
    (5, 5, '05-10-2024', '05-10-2024', '09-10-2024');

    INSERT INTO transaction_detail (transaction_id, product_id, product_price, qty)
    VALUES
    (1, 1, 10000, 2),
    (2, 2, 5000, 5),
    (3, 3, 15000, 3),
    (4, 4, 12000, 4),
    (5, 5, 25000, 1);

```

5. Configure Your database in env file and change the env file name to .env
```bash
DB_HOST=localhost
DB_PORT=5432
DB_USER=username
DB_PASSWORD=password
DB_NAME=example_name
```

6. Navigate to the project directory
```bash
cd challenge-goapi
```

7. Install necessary dependencies
```bash
go mod tidy
```

7. Run the application
```bash
go run main.go
```
    
## Features

- Customer Menu
    - Create Customer
    - View List Of Customer
    - View Customer By Id
    - Update Customer
    - Delete Customer

- Employee Menu
    - Create Employee
    - View List Of Employee
    - View Employee By Id
    - Update Employee
    - Delete Employee

- Product Menu
    - Create Product
    - View List Of Product
    - View Product by Id
    - Update Product
    - Delete Product

- Transaction Menu
    - Create Transaction
    - View List Of Transaction
    - View Transaction By Id

## API Spec
### Customer API

#### Create Customer

Request :

- Method : `POST`
- Endpoint : `/customers`
- Header :
  - Content-Type : application/json
  - Accept : application/json
- Body :

```json
{
  "name": "string",
  "phoneNumber": "string",
  "address": "string"
}
```

Response :

- Status : 201 Created
- Body :

```json
{
  "message": "string",
  "data": {
    "id": "string",
    "name": "string",
    "phoneNumber": "string",
    "address": "string"
  }
}
```

#### Get Customer

Request :

- Method : GET
- Endpoint : `/customers/:id`
- Header :
  - Accept : application/json

Response :

- Status : 200 OK
- Body :

```json
{
  "message": "string",
  "data": {
    "id": "string",
    "name": "string",
    "phoneNumber": "string",
    "address": "string"
  }
}
```

#### Update Customer

Request :

- Method : PUT
- Endpoint : `/customers/:id`
- Header :
  - Content-Type : application/json
  - Accept : application/json
- Body :

```json
{
  "name": "string",
  "phoneNumber": "string",
  "address": "string"
}
```

Response :

- Status : 200 OK
- Body :

```json
{
  "message": "string",
  "data": {
    "id": "string",
    "name": "string",
    "phoneNumber": "string",
    "address": "string"
  }
}
```

#### Delete Customer

Request :

- Method : DELETE
- Endpoint : `/customers/:id`
- Header :
  - Accept : application/json
- Body :

Response :

- Status : 200 OK
- Body :

```json
{
  "message": "string",
  "data": "OK"
}
```

### Employee API

#### Create Employee

Request :

- Method : `POST`
- Endpoint : `/employees`
- Header :
  - Content-Type : application/json
  - Accept : application/json
- Body :

```json
{
  "name": "string",
  "phoneNumber": "string",
  "address": "string"
}
```

Response :

- Status : 201 Created
- Body :

```json
{
  "message": "string",
  "data": {
    "id": "string",
    "name": "string",
    "phoneNumber": "string",
    "address": "string"
  }
}
```

#### Get Employee

Request :

- Method : GET
- Endpoint : `/employees/:id`
- Header :
  - Accept : application/json

Response :

- Status : 200 OK
- Body :

```json
{
  "message": "string",
  "data": {
    "id": "string",
    "name": "string",
    "phoneNumber": "string",
    "address": "string"
  }
}
```

#### Update Employee

Request :

- Method : PUT
- Endpoint : `/employees/:id`
- Header :
  - Content-Type : application/json
  - Accept : application/json
- Body :

```json
{
  "name": "string",
  "phoneNumber": "string",
  "address": "string"
}
```

Response :

- Status : 200 OK
- Body :

```json
{
  "message": "string",
  "data": {
    "id": "string",
    "name": "string",
    "phoneNumber": "string",
    "address": "string"
  }
}
```

#### Delete Employee

Request :

- Method : DELETE
- Endpoint : `/employees/:id`
- Header :
  - Accept : application/json
- Body :

Response :

- Status : 200 OK
- Body :

```json
{
  "message": "string",
  "data": "OK"
}
```

### Product API

#### Create Product

Request :

- Method : POST
- Endpoint : `/products`
- Header :
  - Content-Type : application/json
  - Accept : application/json
- Body :

```json
{
	"name": "string",
  "price": int,
  "unit": "string" (satuan product,cth: Buah atau Kg)
}
```

Response :

- Status Code: 201 Created
- Body:

```json
{
	"message": "string",
	"data": {
		"id": "string",
		"name": "string",
		"price": int,
		"unit": "string" (satuan product,cth: Buah atau Kg)
	}
}
```

#### List Product

Request :

- Method : GET
- Endpoint : `/products`
  - Header :
  - Accept : application/json
- Query Param :
  - productName : string `optional`,

Response :

- Status Code : 200 OK
- Body:

```json
{
	"message": "string",
	"data": [
		{
			"id": "string",
			"name": "string",
			"price": int,
			"unit": "string" (satuan product,cth: Buah atau Kg)
		},
		{
			"id": "string",
			"name": "string",
			"price": int,
			"unit": "string" (satuan product,cth: Buah atau Kg)
		}
	]
}
```

#### Product By Id

Request :

- Method : GET
- Endpoint : `/products/:id`
- Header :
  - Accept : application/json

Response :

- Status Code: 200 OK
- Body :

```json
{
	"message": "string",
	"data": {
		"id": "string",
		"name": "string",
		"price": int,
		"unit": "string" (satuan product,cth: Buah atau Kg)
	}
}
```

#### Update Product

Request :

- Method : PUT
- Endpoint : `/products/:id`
- Header :
  - Content-Type : application/json
  - Accept : application/json
- Body :

```json
{
	"name": "string",
	"price": int,
	"unit": "string" (satuan product,cth: Buah atau Kg)
}
```

Response :

- Status Code: 200 OK
- Body :

```json
{
	"message": "string",
	"data": {
		"id": "string",
		"name": "string",
		"price": int,
		"unit": "string" (satuan product,cth: Buah atau Kg)
	}
}
```

#### Delete Product

Request :

- Method : DELETE
- Endpoint : `/products/:id`
- Header :
  - Accept : application/json
- Body :

Response :

- Status : 200 OK
- Body :

```json
{
  "message": "string",
  "data": "OK"
}
```

### Transaction API

#### Create Transaction

Request :

- Method : POST
- Endpoint : `/transactions`
- Header :
  - Content-Type : application/json
  - Accept : application/json
- Body :

```json
{
	"billDate": "string",
	"entryDate": "string",
	"finishDate": "string",
	"employeeId": "string",
	"customerId": "string",
	"billDetails": [
		{
			"productId": "string",
			"qty": int
		}
	]
}
```

Request :

- Status Code: 201 Created
- Body :

```json
{
	"message": "string",
	"data":  {
		"id":  "string",
		"billDate":  "string",
		"entryDate":  "string",
		"finishDate":  "string",
		"employeeId":  "string",
		"customerId":  "string",
		"billDetails":  [
			{
				"id":	"string",
				"billId":  "string",
				"productId":  "string",
				"productPrice": int,
				"qty": int
			}
		]
	}
}
```

#### Get Transaction

Request :

- Method : GET
- Endpoint : `/transactions/:id_bill`
- Header :
  - Accept : application/json
- Body :

Response :

- Status Code: 200 OK
- Body :

```json
{
	"message": "string",
  "data": {
    "id": "string",
    "billDate": "string",
    "entryDate": "string",
    "finishDate": "string",
    "employee": {
      "id": "string",
      "name": "string",
      "phoneNumber": "string",
      "address": "string"
    },
    "customer": {
      "id": "string",
      "name": "string",
      "phoneNumber": "string",
      "address": "string"
    },
    "billDetails": [
      {
        "id": "string",
        "billId": "string",
        "product": {
          "id": "string",
          "name": "string",
          "price": int,
          "unit": "string" (satuan product,cth: Buah atau Kg)
        },
        "productPrice": int,
        "qty": int
      }
    ],
    "totalBill": int
  }
}
```

#### List Transaction

Pattern string date : `dd-MM-yyyy`

Request :

- Method : GET
- Endpoint : `/transactions`
- Header :
  - Accept : application/json
- Query Param :
  - startDate : string `optional`
  - endDate : string `optional`
  - productName : string `optional`
- Body :

Response :

- Status Code: 200 OK
- Body :

```json
{
	"message": "string",
  "data": [
    {
      "id": "string",
      "billDate": "string",
      "entryDate": "string",
      "finishDate": "string",
      "employee": {
        "id": "string",
        "name": "string",
        "phoneNumber": "string",
        "address": "string"
      },
      "customer": {
        "id": "string",
        "name": "string",
        "phoneNumber": "string",
        "address": "string"
      },
      "billDetails": [
        {
          "id": "string",
          "billId": "string",
          "product": {
            "id": "string",
            "name": "string",
            "price": int,
            "unit": "string" (satuan product,cth: Buah atau Kg)
          },
          "productPrice": int,
          "qty": int
        }
      ],
      "totalBill": int
    }
  ]
}
```