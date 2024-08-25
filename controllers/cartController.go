package controllers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4" // Menggunakan jwt dari golang-jwt/jwt/v4
	"github.com/raihan1405/go-restapi/db"
	"github.com/raihan1405/go-restapi/models"
	"github.com/raihan1405/go-restapi/validators"
)

// ErrorResponse digunakan untuk mengembalikan pesan error yang lebih jelas
type ErrorResponse struct {
	Error string `json:"error"`
}

// SuccessResponse digunakan untuk mengembalikan pesan sukses
type SuccessResponse struct {
	Message string `json:"message"`
}

// AddToCart godoc
// @Summary Add a product to cart
// @Description Add a product to the user's cart
// @Tags cart
// @Accept json
// @Produce json
// @Param cart body validators.AddToCartInput true "Cart item details"
// @Success 200 {object} models.CartItem
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/cart [post]
func AddToCart(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(*jwt.Token)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(ErrorResponse{Error: "Unauthorized"})
	}

	claims, ok := user.Claims.(jwt.MapClaims)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(ErrorResponse{Error: "Invalid token claims"})
	}

	userID, ok := claims["sub"].(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(ErrorResponse{Error: "Invalid user ID in token"})
	}

	var data validators.AddToCartInput

	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{Error: "Cannot parse JSON"})
	}

	if err := validators.Validate.Struct(data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{Error: err.Error()})
	}

	// Find the product
	var product models.Product
	if err := db.DB.First(&product, data.ProductID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(ErrorResponse{Error: "Product not found"})
	}

	// Create cart item
	cartItem := models.CartItem{
		ProductID: data.ProductID,
		UserID:    userID,
		Quantity:  data.Quantity,
	}

	// Save cart item to database
	if err := db.DB.Create(&cartItem).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Error: "Cannot add product to cart"})
	}

	return c.JSON(cartItem)
}

// GetCart godoc
// @Summary Get all items in the cart
// @Description Get a list of all items in the user's cart
// @Tags cart
// @Produce json
// @Success 200 {array} models.CartItem
// @Failure 500 {object} ErrorResponse
// @Router /api/cart [get]
func GetCart(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(*jwt.Token)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(ErrorResponse{Error: "Unauthorized"})
	}

	claims, ok := user.Claims.(jwt.MapClaims)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(ErrorResponse{Error: "Invalid token claims"})
	}

	userID, ok := claims["sub"].(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(ErrorResponse{Error: "Invalid user ID in token"})
	}

	var cartItems []models.CartItem

	// Retrieve all cart items for the user from the database
	if err := db.DB.Where("user_id = ?", userID).Find(&cartItems).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Error: "Cannot retrieve cart items"})
	}

	return c.JSON(cartItems)
}

// RemoveFromCart godoc
// @Summary Remove an item from the cart
// @Description Remove an item from the user's cart by ID
// @Tags cart
// @Param id path int true "Cart Item ID"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/cart/{id} [delete]
func RemoveFromCart(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(*jwt.Token)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(ErrorResponse{Error: "Unauthorized"})
	}

	claims, ok := user.Claims.(jwt.MapClaims)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(ErrorResponse{Error: "Invalid token claims"})
	}

	userID, ok := claims["sub"].(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(ErrorResponse{Error: "Invalid user ID in token"})
	}

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{Error: "Invalid cart item ID"})
	}

	var cartItem models.CartItem
	if err := db.DB.Where("id = ? AND user_id = ?", id, userID).First(&cartItem).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(ErrorResponse{Error: "Cart item not found"})
	}

	// Delete the cart item
	if err := db.DB.Delete(&cartItem).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Error: "Cannot remove cart item"})
	}

	return c.JSON(SuccessResponse{Message: "Item removed from cart"})
}

// UpdateCartItem godoc
// @Summary Update an item in the cart
// @Description Update the quantity of an item in the user's cart
// @Tags cart
// @Accept json
// @Produce json
// @Param id path int true "Cart Item ID"
// @Param cart body validators.UpdateCartItemInput true "Updated cart item details"
// @Success 200 {object} models.CartItem
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/cart/{id} [put]
func UpdateCartItem(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(*jwt.Token)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(ErrorResponse{Error: "Unauthorized"})
	}

	claims, ok := user.Claims.(jwt.MapClaims)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(ErrorResponse{Error: "Invalid token claims"})
	}

	userID, ok := claims["sub"].(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(ErrorResponse{Error: "Invalid user ID in token"})
	}

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{Error: "Invalid cart item ID"})
	}

	var data validators.UpdateCartItemInput
	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{Error: "Cannot parse JSON"})
	}

	if err := validators.Validate.Struct(data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{Error: err.Error()})
	}

	var cartItem models.CartItem
	if err := db.DB.Where("id = ? AND user_id = ?", id, userID).First(&cartItem).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(ErrorResponse{Error: "Cart item not found"})
	}

	// Update the cart item quantity
	cartItem.Quantity = data.Quantity

	if err := db.DB.Save(&cartItem).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Error: "Cannot update cart item"})
	}

	return c.JSON(cartItem)
}
