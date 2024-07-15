package main

import (
	"log"
	"os"

	"github.com/raihan1405/go-restapi/models"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/raihan1405/go-restapi/db"
	"github.com/gofiber/fiber/v2"
)

func getPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	return "0.0.0.0:" + port

}

func main() {
	db := db.Init()
	app := fiber.New()
	app.Use(logger.New())
	models.Setup(db)

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hei")
	})

	log.Fatal(app.Listen(getPort()))
}
