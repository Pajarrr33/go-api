package entity


type Transaction struct {
	Transaction_id     	string				`json:"id"`
	Customer_id        	string 				`json:"customerId"`
	Employee_id        	string 				`json:"employeeId"`
	Bill_date          	string 				`json:"billDate"`
	Entry_date         	string 				`json:"entryDate"`
	Finish_date        	string 				`json:"finishDate"`
	Employee 			Employee 			`json:"employee"`
	Customer 			Customer			`json:"customer"`
	Bill_detail 		[]Transaction_detail  `json:"billDetails"`
	Total_bill			int					`json:"totalBill"`
}