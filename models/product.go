package models


type Product struct{
	ID int    `json:"id"`
	ProductName string `json:"productName"`
	BrandName string `json:"brandName"`
	Price int `json:"price"`
	Status bool `json:"status"`
	Quantity int `json:"quantity"`
	Category    string `json:"Category"`
	UserID      string `json:"userId"`
}
