package models

import (
	"time"

	"gorm.io/gorm"
)

type CartItem struct {
	ID        int    `json:"id"`
	ProductID int    `json:"productId"`
	UserID    string `json:"userId"`
	Quantity  int    `json:"quantity"`
	Product   Product   `gorm:"foreignKey:ProductID"`
	TotalPrice float64   `json:"total_price" gorm:"-"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
}

func (CartItem) Setup(db *gorm.DB) {
	db.AutoMigrate(&CartItem{})
}