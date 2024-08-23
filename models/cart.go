package models

type CartItem struct {
	ID        int    `json:"id"`
	ProductID int    `json:"productId"`
	UserID    string `json:"userId"`
	Quantity  int    `json:"quantity"`
}
