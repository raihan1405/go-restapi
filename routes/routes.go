package routes

import(
	"github.com/raihan1405/go-restapi/controllers"
	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App){

	app.Post("/api/register",controllers.Register)
	app.Post("/api/login",controllers.Login)
	app.Get("/api/user",controllers.User)
	app.Post("/api/logout",controllers.Logout)
	app.Put("/api/user", controllers.UpdateProfile)
}