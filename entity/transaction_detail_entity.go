package entity

type Transaction_detail struct {
	Transaction_detail_id int `json:"id"`
	Transaction_id int `json:"billId"`
	Product_id int `json:"productId"`
	Product_price int `json:"productPrice"`
	Qty int `json:"qty"`
}