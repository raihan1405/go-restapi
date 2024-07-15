package models

import(
	"fmt"
	"log"
	"os"
	"gorm.io/gorm"
	"gorm.io/driver/mysql"
)

var DB *gorm.DB

func ConnectDatabase(){

	dbUser := os.Getenv("root")
	dbPass := os.Getenv("dygXZHHXbSWhwoTmYCieafeSJIJknQqn")
	dbHost := os.Getenv("monorail.proxy.rlwy.net")
	dbName := os.Getenv("railway")

	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", dbUser, dbPass, dbHost, dbName)
	database, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	DB = database
	log.Println("Database connection successfully established")
}