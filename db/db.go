package db

import (
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Init() {
	// Hardcoded values for local MySQL
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		"root",      // MySQL Username, biasanya "root" untuk lokal
		"password",  // MySQL Password, masukkan password MySQL Anda
		"127.0.0.1", // Host, gunakan "127.0.0.1" atau "localhost" untuk lokal
		"3306",      // MySQL port, default adalah "3306"
		"kp",        // Nama database yang ingin Anda hubungkan
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal(err)
		panic("Failed to connect to database")
	}

	DB = db
}
