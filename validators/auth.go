package validators

import "github.com/go-playground/validator/v10"

var Validate = validator.New()

type RegisterInput struct {
	Email       string `json:"email" validate:"required,email"`
	PhoneNumber string `json:"phoneNumber" validate:"required"`
	Username    string `json:"username" validate:"required"`
	Password    string `json:"password" validate:"required,min=6"`
}

type LoginInput struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}