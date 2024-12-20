package routes

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
	"github.com/raihan1405/go-restapi/controllers"
)

func Setup(app *fiber.App) {

	app.Post("/api/register", controllers.Register)
	app.Post("/api/login", controllers.Login)
	app.Post("/admin/login", controllers.LoginAdmin)
	app.Post("/operator/loginOperator", controllers.LoginOperator)
	
	api := app.Group("/api", jwtware.New(jwtware.Config{
		SigningKey:  []byte(os.Getenv("JWT_SECRET")), 
		TokenLookup: "cookie:jwt",
	}))

	api.Get("/userProducts", controllers.GetAllProducts)
	api.Get("/user", controllers.GetUser)
	api.Post("/logoutUser", controllers.LogoutUser)
	api.Put("/user", controllers.UpdateProfile)
	api.Put("/user/password", controllers.UpdatePassword)
	api.Get("/Products", controllers.GetAllProducts)
	api.Post("/addToCart", controllers.AddToCart)
	api.Get("/itemCart", controllers.GetCart)
	api.Put("/itemCart/edit/:id", controllers.UpdateCartItem)
	api.Delete("deleteCart/:id", controllers.RemoveFromCart)
	api.Post("/createInvoice", controllers.CreateInvoice)
	api.Get("/getInvoice", controllers.GetAllInvoices)
	

	apiOperator := app.Group("/operator", jwtware.New(jwtware.Config{
		SigningKey:   []byte(os.Getenv("JWT_SECRET_OPERATOR")),
		TokenLookup:  "cookie:jwt_operator",                    
		ErrorHandler: jwtError,
	}))

	apiOperator.Get("/Products", controllers.GetAllProducts)
	apiOperator.Get("/dashboard", controllers.OperatorDashboard)
	apiOperator.Post("/products", controllers.AddProduct)
	apiOperator.Post("/logoutOperator", controllers.LogoutOperator)
	apiOperator.Put("/products/edit/:id", controllers.EditProduct)
	apiOperator.Get("/getAllInvoice", controllers.GetAllInvoicesForOperator)
	apiOperator.Put("/invoices/approve", controllers.ApproveInvoices)
	apiOperator.Put("/invoices/reject", controllers.RejectInvoices)
	apiOperator.Get("/invoices/accepted", controllers.GetAcceptInvoice)
	apiOperator.Put("/invoices/updateShipment", controllers.UpdateStatusInvoice)


	apiAdmin := app.Group("/admin", jwtware.New(jwtware.Config{
		SigningKey:   []byte(os.Getenv("JWT_SECRET_ADMIN")),
		TokenLookup:  "cookie:jwt_admin",                    
		ErrorHandler: jwtError,
	}))

	apiAdmin.Get("/adminProducts", controllers.GetAllProducts)
	apiAdmin.Post("/logoutAdmin", controllers.LogoutAdmin)
	apiAdmin.Get("/productAdmin", controllers.GetAllProducts)
	apiAdmin.Get("/getAllInvoiceAdmin", controllers.GetAllInvoicesForAdmin)
	apiAdmin.Get("/getProductReport/:id", controllers.GenerateProductReport)

}

func jwtError(c *fiber.Ctx, err error) error {
	log.Printf("JWT Error: %v\n", err)
	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
		"message": "unauthenticated",
		"error":   err.Error(),
	})
}
