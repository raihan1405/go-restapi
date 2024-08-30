package routes

import (
	"os"

	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
	"github.com/raihan1405/go-restapi/controllers"
)

func Setup(app *fiber.App) {

	// Rute publik yang tidak membutuhkan autentikasi
	app.Post("/api/register", controllers.Register)
	app.Post("/api/login", controllers.Login)
	app.Post("/api/products", controllers.AddProduct)
	app.Get("/api/products", controllers.GetAllProducts)
	app.Put("/api/products/:id", controllers.EditProduct)

	// Middleware JWT untuk melindungi rute di bawah ini
	api := app.Group("/api", jwtware.New(jwtware.Config{
		SigningKey: []byte(os.Getenv("JWT_SECRET")), // Gantilah dengan kunci rahasia yang sebenarnya
		TokenLookup: "cookie:jwt",
	}))

	// Rute yang dilindungi oleh JWT middleware
	api.Get("/user", controllers.GetUser)
	api.Post("/logout", controllers.Logout)
	api.Put("/user", controllers.UpdateProfile)
	api.Put("/user/password", controllers.UpdatePassword)

	// Rute untuk produk
	
}
