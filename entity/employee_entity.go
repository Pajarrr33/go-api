package entity

type Employee struct {
	Employee_id int `json:"id"`
	Name string `json:"name"`
	Phone_number string `json:"phoneNumber"`
	Address string `json:"address"`
}