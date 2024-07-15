package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/raihan1405/go-restapi/models"
	"github.com/raihan1405/go-restapi/db"
	"golang.org/x/crypto/bcrypt"
)

func Register(c *fiber.Ctx) error {
	var data map[string] string

	err := c.BodyParser(&data)

	if err != nil {
		return err
	}

	password, _ := bcrypt.GenerateFromPassword([]byte(data["password"]),14)
	user:= models.User{
		Username : data["username"],
		Email : data["email"],
		Password: password,
	}
	
	db.DB.Create(&user)
	

	return c.JSON(user)
}