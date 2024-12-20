package models

import (
	"time"
	"gorm.io/gorm"
)

// ProductOut menyimpan transaksi produk yang keluar
type ProductOut struct {
	ID        uint      `gorm:"primaryKey"`
	ProductID uint      `json:"productId"` // ID produk yang keluar
	Quantity  int       `json:"quantity"`  // Jumlah yang keluar
	OperatorID  string `json:"operator_id"`
	CreatedAt time.Time `json:"createdAt"` // Waktu transaksi
}

// Setup untuk otomatis migrasi tabel ProductOut
func (ProductOut) Setup(db *gorm.DB) {
	db.AutoMigrate(&ProductOut{})
}
