package main

import (
	"log"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/swagger"
	"github.com/joho/godotenv"
	"github.com/raihan1405/go-restapi/db"
	"github.com/raihan1405/go-restapi/models"
	"github.com/raihan1405/go-restapi/routes"
	_ "github.com/raihan1405/go-restapi/docs"
)

func getPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	return "0.0.0.0:" + port

}
//	@title			Swagger Example API
//	@version		1.0
//	@description	This is a sample server celler server.
//	@termsOfService	http://swagger.io/terms/
func main() {

	
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	
	
	app := fiber.New()
	allowedOrigins := "http://localhost:5173,https://go-restapi-production.up.railway.app:8080","https://sjr-app-dev-6rs5.vercel.app/"
	allowedOriginsList := strings.Split(allowedOrigins, ",")

	// Set up CORS middleware with dynamic origin checking
	app.Use(cors.New(cors.Config{
		AllowOriginsFunc: func(origin string) bool {
			for _, allowedOrigin := range allowedOriginsList {
				if origin == allowedOrigin {
					return true
				}
			}
			return false
		},
		AllowCredentials: true,
		AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
		AllowHeaders:     "Origin,Content-Type,Accept,Authorization",
	}))
	

	db.Init()
	models.Setup(db.DB)
	routes.Setup(app)

	app.Get("/swagger/*", swagger.HandlerDefault) // default

	log.Fatal(app.Listen(getPort()))
}
