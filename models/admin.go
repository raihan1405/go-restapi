package models

import (
	"time"

	"gorm.io/gorm"
)

type Admin struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	AdminID string         `gorm:"unique;not null" json:"admin_id"`
	Name       string         `json:"name"`
	Email      string         `gorm:"unique" json:"email"`
	Phone      string         `json:"phone"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`

}

func (Admin) Setup(db *gorm.DB) {
	db.AutoMigrate(&Admin{})
}