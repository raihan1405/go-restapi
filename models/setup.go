package models

import "gorm.io/gorm"

var db *gorm.DB

func Setup() {
	db.AutoMigrate(&User{})
}
