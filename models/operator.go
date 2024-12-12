package models

import (
	"time"

	"gorm.io/gorm"
)

type Operator struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	OperatorID string         `gorm:"unique;not null" json:"operator_id"`
	Name       string         `json:"name"`
	Email      string         `gorm:"unique" json:"email"`
	Phone      string         `json:"phone"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`

}

func (Operator) Setup(db *gorm.DB) {
	db.AutoMigrate(&Operator{})
}