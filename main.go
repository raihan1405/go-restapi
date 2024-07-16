package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
	"github.com/raihan1405/go-restapi/db"
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
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "https://go-restapi-production.up.railway.app:8080", 
		AllowCredentials: true,
		AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
		AllowHeaders:     "Origin,Content-Type,Accept,Authorization",
	}))
	

	db.Init()
	models.Setup(db.DB)
	routes.Setup(app)

	log.Fatal(app.Listen(getPort()))
}
