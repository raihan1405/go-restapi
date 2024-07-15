package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
	"github.com/raihan1405/go-restapi/models"
	"github.com/raihan1405/go-restapi/routes"
)

func getPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	return "0.0.0.0:" + port

}

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	
	
	app := fiber.New()
	app.Use(logger.New())
	

	models.Setup()
	routes.Setup(app)

	log.Fatal(app.Listen(getPort()))
}
