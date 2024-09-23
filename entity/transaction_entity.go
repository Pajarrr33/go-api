package entity

type Transaction struct {
	Transaction_id int `json:"id"`
	Customer_id string `json:"customerId"`
	Employee_id string `json:"employeeId"`
	Bill_date string `json:"billDate"`
	Entry_date string `json:"entryDate"`
	Finish_date string `json:"finishDate"`
}