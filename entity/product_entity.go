package entity

type Product struct {
	Product_id string `json:"id"`
	Product_name string `json:"name"`
	Price int `json:"price"`
	Unit string `json:"unit"`
}