package entity

type Transaction_detail struct {
	Transaction_detail_id 	string `json:"id"`
	Transaction_id 			string `json:"billId"`
	Product 				Product `json:"product"`
	Product_id 				string `json:"productId"`
	Product_price 			int `json:"productPrice"`
	Qty 					int `json:"qty"`
}