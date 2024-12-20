package models

import (
	"gorm.io/gorm"
)


// Struct untuk laporan produk
type ProductReport struct {
    ProductName      string `json:"productName"`
    ProductID        string `json:"productId"`
    Category         string `json:"category"`
    BrandName        string `json:"brandName"`
    InitialStock     int    `json:"initialStock"`
    FirstInStock     int    `json:"firstInStock"`
    FirstOutStock    int    `json:"firstOutStock"`
    StockAvailability int   `json:"stockAvailability"`
}


func (ProductReport) Setup(db *gorm.DB) {
	db.AutoMigrate(&ProductReport{})
}