package models

import "gorm.io/gorm"

var db *gorm.DB

func Setup(db *gorm.DB) {
	db.AutoMigrate(&User{})
}
