package models

import (
	"time"
	"gorm.io/gorm"
)

// InvoiceDetailResponse mewakili response untuk detail invoice
type InvoiceDetailResponse struct {
	Invoice      Invoice        `json:"invoice"`
	InvoiceItems []InvoiceItem `json:"invoice_items"`
}

// Invoice mewakili invoice yang dibuat dari keranjang
type Invoice struct {
	ID          int           `json:"id"`
	UserID      string        `json:"user_id"`
	User         User          `json:"user" gorm:"foreignkey:UserID"`
	TotalPrice  float64       `json:"total_price"`
	CreatedAt   time.Time     `json:"created_at"`
	Status      string        `json:"status"`
	StatusShipment string       `json:"status_shipment"`
	InvoiceItems []InvoiceItem `json:"invoice_items" gorm:"foreignkey:InvoiceID"`
}


// InvoiceItem mewakili item yang ada dalam invoice
type InvoiceItem struct {
	ID        int     `json:"id"`
	InvoiceID int     `json:"invoice_id"`
	ProductID int     `json:"product_id"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
	Total     float64 `json:"total"` // total = quantity * price
	Product   Product `json:"product" gorm:"foreignkey:ProductID"` // Preload the Product details
}

func (Invoice) Setup(db *gorm.DB) {
	db.AutoMigrate(&Invoice{})
}

func (InvoiceItem) Setup(db *gorm.DB) {
	db.AutoMigrate(&InvoiceItem{})
}