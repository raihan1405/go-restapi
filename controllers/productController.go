package controllers

import (
	"strconv"
	"log"
	"os"
	"fmt"


	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	//"github.com/golang-jwt/jwt/v4" // Menggunakan jwt dari golang-jwt/jwt/v4
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
    userInterface := c.Locals("user")
    if userInterface == nil {
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
    }

    user, ok := userInterface.(*jwt.Token)
    if !ok {
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid user token"})
    }

    claims, ok := user.Claims.(jwt.MapClaims)
    if !ok {
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token claims"})
    }

    operatorID, ok := claims["sub"].(string)
    if !ok || operatorID == "" {
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid operator ID in token"})
    }

    var data validators.AddProductInput

    // Parse data into the structure
    if err := c.BodyParser(&data); err != nil {
        log.Printf("Error parsing JSON: %v\n", err)
        return c.Status(fiber.StatusBadRequest).JSON(map[string]interface{}{"error": "Cannot parse JSON"})
    }

    // Validate input data
    if err := validators.Validate.Struct(data); err != nil {
        log.Printf("Validation error: %v\n", err)
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
        Category:    data.Category, // Menyimpan Category
        OperatorID:  operatorID,     // Menyimpan OperatorID
    }

    // Save product to database
    if err := db.DB.Create(&product).Error; err != nil {
        log.Printf("Database error: %v\n", err)
        return c.Status(fiber.StatusInternalServerError).JSON(map[string]interface{}{"error": "Cannot save product"})
    }

    return c.Status(fiber.StatusCreated).JSON(product)
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
    // Get token from cookie
    cookie := c.Cookies("jwt_operator")

    // Parse the token to validate it and extract claims
    token, err := jwt.ParseWithClaims(cookie, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
        // Use the appropriate signing key based on the kid
        keyID, ok := token.Header["kid"].(string)
        if !ok {
            return nil, jwt.NewValidationError("missing kid header", jwt.ValidationErrorClaimsInvalid)
        }
		fmt.Println(keyID)

        switch keyID {
        case "operator":
            return []byte(os.Getenv("JWT_SECRET_OPERATOR")), nil
        case "user":
            return []byte(os.Getenv("JWT_SECRET")), nil
        default:
            return nil, jwt.NewValidationError("invalid kid", jwt.ValidationErrorClaimsInvalid)
        }
    })

    if err != nil || !token.Valid {
        return c.Status(fiber.StatusUnauthorized).JSON(map[string]interface{}{"error": "unauthenticated"})
    }

    // Get product ID from the URL parameter
    id, err := strconv.Atoi(c.Params("id"))
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(map[string]interface{}{"error": "Invalid product ID"})
    }

    // Parse the incoming request body
    var data validators.EditProductInput
    if err := c.BodyParser(&data); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(map[string]interface{}{"error": "Cannot parse JSON"})
    }

    // Validate the input
    if err := validators.Validate.Struct(data); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(map[string]interface{}{"error": err.Error()})
    }

    // Validate quantity is not negative
    if data.Quantity < 0 {
        return c.Status(fiber.StatusBadRequest).JSON(map[string]interface{}{"error": "Quantity cannot be negative"})
    }

    // Find product by ID
    var product models.Product
    if err := db.DB.First(&product, id).Error; err != nil {
        return c.Status(fiber.StatusNotFound).JSON(map[string]interface{}{"error": "Product not found"})
    }

    // Update product details
    product.ProductName = data.ProductName
    product.BrandName = data.BrandName
    product.Category = data.Category
    product.Price = int(data.Price)
    product.Quantity = data.Quantity
    product.Status = data.Quantity > 0

    // Save the updated product
    if err := db.DB.Save(&product).Error; err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(map[string]interface{}{"error": "Cannot update product"})
    }

    // Return the updated product
    return c.JSON(product)
}