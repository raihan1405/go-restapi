package models

type User struct {
	
	ID       int
	Email string `json:"email"`
	PhoneNumber string `json:"phoneNumber"`
	Username string `json:"username"`
	Password []byte `json:"password"`
}




