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

type UpdateUserInput struct {
    Username    string `json:"username" validate:"required"`
    Email       string `json:"email" validate:"required,email"`
    PhoneNumber string `json:"phone_number" validate:"required"`
}

type UpdatePasswordInput struct {
    OldPassword string `json:"old_password" validate:"required"`
    NewPassword string `json:"new_password" validate:"required,min=8"`
}

type AddProductInput struct {
    ProductName string `json:"productName" validate:"required"`
    BrandName   string `json:"brandName" validate:"required"`
    Price       int    `json:"price" validate:"required"`
    Quantity    int    `json:"quantity" validate:"required"`
}

// EditProductInput represents the input data for editing an existing product
type EditProductInput struct {
    ProductName string  `json:"product_name" validate:"required,min=2,max=100"`
    BrandName   string  `json:"brand_name" validate:"required,min=2,max=100"`
    Price       float64 `json:"price" validate:"required,gt=0"`
    Quantity    int     `json:"quantity" validate:"required,min=0"`
}

type AddToCartInput struct {
	ProductID int `json:"productId" validate:"required"`
	Quantity  int `json:"quantity" validate:"required,min=1"`
}

type UpdateCartItemInput struct {
	Quantity int `json:"quantity" validate:"required,min=1"`
}
