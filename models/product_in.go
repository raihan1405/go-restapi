package models

import (
	"time"
	"gorm.io/gorm"
)

// ProductIn menyimpan transaksi produk yang masuk
type ProductIn struct {
	ID        uint      `gorm:"primaryKey"`
	ProductID int      `json:"productId"` // ID produk yang masuk
	Quantity  int       `json:"quantity"`  // Jumlah yang masuk
	OperatorID  string `json:"operator_id"`
	CreatedAt time.Time `json:"createdAt"` // Waktu transaksi
}

// Setup untuk otomatis migrasi tabel ProductIn
func (ProductIn) Setup(db *gorm.DB) {
	db.AutoMigrate(&ProductIn{})
}
