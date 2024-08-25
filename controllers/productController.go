package controllers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4" // Menggunakan jwt dari golang-jwt/jwt/v4
	"github.com/raihan1405/go-restapi/db"
	"github.com/raihan1405/go-restapi/models"
	"github.com/raihan1405/go-restapi/validators"
)

// AddProduct godoc
// @Summary Add a new product
// @Description Add a new product with the provided details
// @Tags product
// @Accept json
// @Produce json
// @Param product body validators.AddProductInput true "Product details"
// @Success 200 {object} models.Product
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/products [post]
func AddProduct(c *fiber.Ctx) error {
	// Mendapatkan user dari context (yang di-set oleh middleware JWT)
	user, ok := c.Locals("user").(*jwt.Token)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	claims, ok := user.Claims.(jwt.MapClaims)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token claims"})
	}

	userID, ok := claims["sub"].(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid user ID in token"})
	}

	var data validators.AddProductInput

	// Parse data into the structure
	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(map[string]interface{}{"error": "Cannot parse JSON"})
	}

	// Validate input data
	if err := validators.Validate.Struct(data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(map[string]interface{}{"error": err.Error()})
	}

	// Set status based on quantity
	status := data.Quantity > 0

	// Create product
	product := models.Product{
		ProductName: data.ProductName,
		BrandName:   data.BrandName,
		Price:       int(data.Price),
		Status:      status,
		Quantity:    data.Quantity,
		UserID:      userID,
	}

	// Save product to database
	if err := db.DB.Create(&product).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(map[string]interface{}{"error": "Cannot save product"})
	}

	return c.JSON(product)
}


// GetAllProducts godoc
// @Summary Get all products
// @Description Get a list of all products
// @Tags product
// @Produce json
// @Success 200 {array} models.Product
// @Failure 500 {object} map[string]interface{}
// @Router /api/products [get]
func GetAllProducts(c *fiber.Ctx) error {
	var products []models.Product

	// Retrieve all products from the database
	if err := db.DB.Find(&products).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(map[string]interface{}{"error": "Cannot retrieve products"})
	}

	return c.JSON(products)
}

// EditProduct godoc
// @Summary Edit an existing product
// @Description Edit an existing product with the provided details
// @Tags product
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Param product body validators.EditProductInput true "Product details"
// @Success 200 {object} models.Product
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/products/{id} [put]
func EditProduct(c *fiber.Ctx) error {
	var data validators.EditProductInput
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(map[string]interface{}{"error": "Invalid product ID"})
	}

	// Parse data into the structure
	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(map[string]interface{}{"error": "Cannot parse JSON"})
	}

	// Validate input data
	if err := validators.Validate.Struct(data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(map[string]interface{}{"error": err.Error()})
	}

	// Find the product by ID
	var product models.Product
	if err := db.DB.First(&product, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(map[string]interface{}{"error": "Product not found"})
	}

	// Update product details
	product.ProductName = data.ProductName
	product.BrandName = data.BrandName
	product.Price = int(data.Price)
	product.Quantity = data.Quantity
	product.Status = data.Quantity > 0

	// Save changes to database
	if err := db.DB.Save(&product).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(map[string]interface{}{"error": "Cannot update product"})
	}

	return c.JSON(product)
}
