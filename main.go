package main

import (
	"log"
	"os"

	"github.com/raihan1405/go-restapi/models"
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()
	models.ConnectDatabase()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	port := os.Getenv("PORT")

	if port == "" {
		port = "3000"
	}

	log.Fatal(app.Listen("0.0.0.0:" + port))
}
