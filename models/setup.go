package models

import "gorm.io/gorm"

func Setup(db *gorm.DB) {
	db.AutoMigrate(
		&User{},
		&Product{},
	)
}
