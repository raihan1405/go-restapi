package db

import(
	"fmt"
	"log"
	"os"
	"gorm.io/gorm"
	"gorm.io/driver/mysql"
)

var DB *gorm.DB

func Init() *gorm.DB {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		os.Getenv("MYSQLUSER"),
		os.Getenv("MYSQLPASSWORD"),
		os.Getenv("MYSQLHOST"),
		os.Getenv("MYSQLPORT"),
		os.Getenv("MYSQLDATABASE"),
	)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal(err)
		panic("Failed to connect to database")
	}

	return db

}

func Setyup(db *gorm.DB) {
	db.AutoMigrate()
}