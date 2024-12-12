package models

type User struct {
	ID          int    `json:"id"`
	Email       string `json:"email" validate:"required,email"`
	PhoneNumber string `json:"phoneNumber" validate:"required"`
	Username    string `json:"username" validate:"required"`
	Password    []byte `json:"password" validate:"required"`
	Invoices  []Invoice `json:"invoices" gorm:"foreignkey:UserID"`
}
