package models

type User struct {
	ID       int
	Email string `json:"email"`
	Username string `json:"username"`
	Password []byte `json:"password"`
}

