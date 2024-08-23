package models

import "gorm.io/gorm"

func Setup(db *gorm.DB) {
	db.Migrator().DropTable(
		&User{},
	)
	db.AutoMigrate(
		&User{},
		&Product{},
	)
}
