package main

import (
	"log"
	"os"

	
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

	// err := godotenv.Load()
	// if err != nil {
	// 	log.Fatal("Error loading .env file")
	// }
	
	// db := db.Init()

	app := fiber.New()
	// app.Use(logger.New())
	// models.Setup(db)

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hei")
	})

	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	log.Fatal(app.Listen("0.0.0.0:" + port))
}
