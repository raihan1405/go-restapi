package routes

import(
	"github.com/raihan1405/go-restapi/controllers"
	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App){
	app.Post("/api/register",controllers.Register)
}