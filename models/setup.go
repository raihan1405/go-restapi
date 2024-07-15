package models

import(
	"gorm.io/gorm"
	"gorm.io/driver/mysql"
)

var DB *gorm.DB

func ConnectDatabase(){
	database, err := gorm.Open(mysql.Open())
}