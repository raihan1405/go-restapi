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
		os.Getenv("root"),
		os.Getenv("dygXZHHXbSWhwoTmYCieafeSJIJknQqn"),
		os.Getenv("monorail.proxy.rlwy.net"),
		os.Getenv("44256"),
		os.Getenv("railway"),
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