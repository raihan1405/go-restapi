package db

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Init() {
	// Load the .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Retrieve the MySQL credentials from the .env file
	mysqlUser := os.Getenv("MYSQLUSER")
	mysqlPassword := os.Getenv("MYSQLPASSWORD")
	mysqlHost := os.Getenv("MYSQLHOST")
	mysqlPort := os.Getenv("MYSQLPORT")
	mysqlDatabase := os.Getenv("MYSQLDATABASE")

	// Construct the DSN (Data Source Name) with the environment variables
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		mysqlUser,     // MySQL Username
		mysqlPassword, // MySQL Password
		mysqlHost,     // MySQL Host
		mysqlPort,     // MySQL Port
		mysqlDatabase, // MySQL Database name
	)

	// Open a connection to the MySQL database
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	// Handle any errors that occur during connection
	if err != nil {
		log.Fatal(err)
		panic("Failed to connect to database")
	}else{
		fmt.Println("success to connect")
	}

	// Assign the database connection to the global DB variable
	DB = db
}
