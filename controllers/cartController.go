package controllers

import (
	"strconv"
	"time"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4" // Menggunakan jwt dari golang-jwt/jwt/v4
	"github.com/raihan1405/go-restapi/db"
	"github.com/raihan1405/go-restapi/models"
	"github.com/raihan1405/go-restapi/validators"
)

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
	// Verify JWT token and get user info from token
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

	// Parse request body to extract AddToCartInput data
	var data validators.AddToCartInput
	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{Error: "Cannot parse JSON"})
	}

	// Validate input data
	if err := validators.Validate.Struct(data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{Error: err.Error()})
	}

	// Check if the product exists in the database
	var product models.Product
	if err := db.DB.First(&product, data.ProductID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(ErrorResponse{Error: "Product not found"})
	}

	// Check if the product is in stock
	if product.Quantity <= 0 {
		return c.Status(fiber.StatusConflict).JSON(ErrorResponse{Error: "Product is out of stock"})
	}

	// Check if the cart item already exists
	var existingCartItem models.CartItem
	if err := db.DB.Where("user_id = ? AND product_id = ?", userID, data.ProductID).First(&existingCartItem).Error; err == nil {
		// If the product is already in the cart, update the quantity
		existingCartItem.Quantity += data.Quantity
		if err := db.DB.Save(&existingCartItem).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Error: "Failed to update cart"})
		}
		return c.JSON(existingCartItem)
	}

	// Create a new cart item if not already in the cart
	cartItem := models.CartItem{
		ProductID: data.ProductID,
		UserID:    userID,
		Quantity:  data.Quantity,
	}

	// Save cart item to the database
	if err := db.DB.Create(&cartItem).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Error: "Cannot add product to cart"})
	}

	// Return the cart item response
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

	// Retrieve all cart items for the user from the database along with related product details
	if err := db.DB.Preload("Product").Where("user_id = ?", userID).Find(&cartItems).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Error: "Cannot retrieve cart items"})
	}

	// Calculate TotalPrice for each cart item dynamically
	for i := range cartItems {
		// Calculate the total price (Quantity * Product Price)
		cartItems[i].TotalPrice = float64(cartItems[i].Quantity * cartItems[i].Product.Price)
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
	// Get the JWT token and verify the user
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

	// Parse the cart item ID from the URL parameter
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{Error: "Invalid cart item ID"})
	}

	// Parse the new data for updating the cart item
	var data validators.UpdateCartItemInput
	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{Error: "Cannot parse JSON"})
	}

	// Validate the input data
	if err := validators.Validate.Struct(data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{Error: err.Error()})
	}

	// Retrieve the cart item from the database
	var cartItem models.CartItem
	if err := db.DB.Where("id = ? AND user_id = ?", id, userID).First(&cartItem).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(ErrorResponse{Error: "Cart item not found"})
	}

	// Check if the quantity is valid (greater than zero)
	if data.Quantity <= 0 {
		// Optionally, you could delete the cart item if the quantity is zero
		if err := db.DB.Delete(&cartItem).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Error: "Cannot remove cart item"})
		}
		return c.JSON(SuccessResponse{Message: "Item removed from cart"})
	}

	// Retrieve product to check stock availability
	var product models.Product
	if err := db.DB.First(&product, cartItem.ProductID).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Error: "Product not found"})
	}

	// Check if the new quantity exceeds available stock
	if data.Quantity > product.Quantity {
		return c.Status(fiber.StatusConflict).JSON(ErrorResponse{Error: "Insufficient stock"})
	}

	// Update the cart item quantity
	cartItem.Quantity = data.Quantity

	// Save the updated cart item to the database
	if err := db.DB.Save(&cartItem).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Error: "Cannot update cart item"})
	}

	// Return the updated cart item
	return c.JSON(cartItem)
}


// CreateInvoice godoc
// @Summary Create an invoice from the user's cart
// @Description Create an invoice from all selected items in the user's cart
// @Tags invoice
// @Accept json
// @Produce json
// @Success 201 {object} models.Invoice
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/invoice [post]
func CreateInvoice(c *fiber.Ctx) error {
	// Verify JWT token and get user info from token
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

	// Retrieve all items in the user's cart
	var cartItems []models.CartItem
	if err := db.DB.Preload("Product").Where("user_id = ?", userID).Find(&cartItems).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Error: "Failed to retrieve cart items"})
	}

	// Calculate total price of the items in the cart
	var totalPrice float64
	for _, item := range cartItems {
		totalPrice += float64(item.Quantity * item.Product.Price)
	}

	// Create a new invoice
	invoice := models.Invoice{
		UserID:     userID,
		TotalPrice: totalPrice,
		CreatedAt:  time.Now(),
		Status:     "Pending", // You can set status to "Pending" or as appropriate
	}

	// Save the invoice to the database
	if err := db.DB.Create(&invoice).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Error: "Failed to create invoice"})
	}

	// Create invoice items for each cart item
	for _, item := range cartItems {
		invoiceItem := models.InvoiceItem{
			InvoiceID: invoice.ID,
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     float64(item.Product.Price),
			Total:     float64(item.Quantity * item.Product.Price),
		}

		// Save the invoice item to the database
		if err := db.DB.Create(&invoiceItem).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Error: "Failed to create invoice items"})
		}
	}

	// Optionally, clear the user's cart after creating the invoice
	if err := db.DB.Where("user_id = ?", userID).Delete(&models.CartItem{}).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Error: "Failed to clear cart after creating invoice"})
	}

	// Return the created invoice
	return c.Status(fiber.StatusCreated).JSON(invoice)
}

// GetAllInvoices godoc
// @Summary Get all invoices for the logged-in user
// @Description Get a list of all invoices associated with the logged-in user
// @Tags invoice
// @Produce json
// @Success 200 {array} models.Invoice
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/invoices [get]
// GetAllInvoices godoc
// @Summary Get all invoices for the logged-in user
// @Description Get a list of all invoices associated with the logged-in user
// @Tags invoice
// @Produce json
// @Success 200 {array} models.Invoice
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/invoices [get]
// GetAllInvoices godoc
// @Summary Get all invoices for the logged-in user
// @Description Get a list of all invoices associated with the logged-in user
// @Tags invoice
// @Produce json
// @Success 200 {array} models.Invoice
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/invoices [get]
func GetAllInvoices(c *fiber.Ctx) error {
	// Get the JWT token and verify the user
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

	// Retrieve all invoices for the logged-in user, preload related InvoiceItems, Products, and User
	var invoices []models.Invoice
	if err := db.DB.Preload("InvoiceItems.Product").Preload("User").Where("user_id = ?", userID).Find(&invoices).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Error: "Failed to retrieve invoices"})
	}

	// Return the list of invoices
	return c.JSON(invoices)
}

// GetAllInvoicesForOperator godoc
// @Summary Get all invoices from all users (for operator)
// @Description Get a list of all invoices for all users
// @Tags invoice
// @Produce json
// @Success 200 {array} models.Invoice
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/operator/invoices [get]
// GetAllInvoices godoc
// @Summary Get all invoices for all users (operator access)
// @Description Get a list of all invoices associated with all users (accessible by operator)
// @Tags invoice
// @Produce json
// @Success 200 {array} models.Invoice
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/invoices [get]
func GetAllInvoicesForOperator(c *fiber.Ctx) error {
    // Ambil token JWT dari cookie
    cookie := c.Cookies("jwt_operator")

    // Jika tidak ada token di cookie
    if cookie == "" {
        return c.Status(fiber.StatusUnauthorized).JSON(ErrorResponse{
            Message: "unauthenticated",
            Error:   "No JWT token found in cookie",
        })
    }

    // Parse token dengan claims
    token, err := jwt.ParseWithClaims(cookie, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
        // Mengambil key ID (kid) dari header JWT
        keyID, ok := token.Header["kid"].(string)
        if !ok {
            return nil, jwt.NewValidationError("missing kid header", jwt.ValidationErrorClaimsInvalid)
        }

        // Return signing key berdasarkan kid
        switch keyID {
        case "operator":
            return []byte(os.Getenv("JWT_SECRET_OPERATOR")), nil
        case "user":
            return []byte(os.Getenv("JWT_SECRET")), nil
        default:
            return nil, jwt.NewValidationError("invalid kid", jwt.ValidationErrorClaimsInvalid)
        }
    })

    // Cek error token atau invalid token
    if err != nil || !token.Valid {
        return c.Status(fiber.StatusUnauthorized).JSON(ErrorResponse{
            Message: "unauthenticated",
            Error:   "Invalid or expired token",
        })
    }

    // Ambil claims dari token
    claims, ok := token.Claims.(*jwt.StandardClaims)
    if !ok {
        return c.Status(fiber.StatusUnauthorized).JSON(ErrorResponse{
            Message: "unauthenticated",
            Error:   "Invalid token claims",
        })
    }

    // Konversi Subject menjadi integer (Operator ID)
    operatorID, err := strconv.Atoi(claims.Subject)
    if err != nil {
        return c.Status(fiber.StatusUnauthorized).JSON(ErrorResponse{
            Message: "Invalid token",
            Error:   "Token contains invalid operator ID",
        })
    }

    // Cari operator di database
    var operator models.Operator
    if err := db.DB.Where("id = ?", operatorID).First(&operator).Error; err != nil {
        return c.Status(fiber.StatusNotFound).JSON(ErrorResponse{
            Message: "operator not found",
            Error:   "No operator with the given ID",
        })
    }

    // Ambil semua invoice dari database, preload data terkait InvoiceItems dan Products
    var invoices []models.Invoice
    if err := db.DB.Preload("InvoiceItems.Product").Preload("User").Find(&invoices).Error; err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
            Message: "failed to retrieve invoices",
            Error:   "Failed to retrieve invoices from the database",
        })
    }

    // Return all invoices as JSON response
    return c.JSON(invoices)
}

// ApproveMultipleInvoices memungkinkan operator untuk menyetujui beberapa pesanan sekaligus
func ApproveInvoices(c *fiber.Ctx) error {
    // Ambil token dari cookie
    cookie := c.Cookies("jwt_operator")

    // Parse token dengan claims
    token, err := jwt.ParseWithClaims(cookie, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
        keyID, ok := token.Header["kid"].(string)
        if !ok {
            return nil, jwt.NewValidationError("missing kid header", jwt.ValidationErrorClaimsInvalid)
        }

        switch keyID {
        case "operator":
            return []byte(os.Getenv("JWT_SECRET_OPERATOR")), nil
        default:
            return nil, jwt.NewValidationError("invalid kid", jwt.ValidationErrorClaimsInvalid)
        }
    })

    if err != nil || !token.Valid {
        return c.Status(fiber.StatusUnauthorized).JSON(ErrorResponse{
            Message: "unauthenticated",
            Error:   "Invalid or expired token",
        })
    }

    // Ambil operatorID dari claims
    claims, ok := token.Claims.(*jwt.StandardClaims)
    if !ok {
        return c.Status(fiber.StatusUnauthorized).JSON(ErrorResponse{
            Message: "unauthenticated",
            Error:   "Invalid token claims",
        })
    }

	operatorID := claims.Subject
	var operator models.Operator
    if err := db.DB.Where("id = ?", operatorID).First(&operator).Error; err != nil {
        return c.Status(fiber.StatusNotFound).JSON(ErrorResponse{
            Message: "operator not found",
            Error:   "No operator with the given ID",
        })
    }

    // Parse daftar order IDs yang ingin di-approve
    var orderIDs []int
    if err := c.BodyParser(&orderIDs); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
            Message: "invalid input",
            Error:   "Failed to parse order IDs",
        })
    }

    // Update status setiap order
    for _, orderID := range orderIDs {
        if err := db.DB.Model(&models.Invoice{}).Where("id = ?", orderID).Update("status", "Approved").Error; err != nil {
            return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
                Message: "failed to update status",
                Error:   "Failed to update status for order ID " + strconv.Itoa(orderID),
            })
        }
    }

    return c.JSON(fiber.Map{
        "message": "Orders successfully approved",
    })
}

func RejectInvoices(c *fiber.Ctx) error {
    // Ambil token dari cookie
    cookie := c.Cookies("jwt_operator")

    // Parse token dengan claims
    token, err := jwt.ParseWithClaims(cookie, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
        keyID, ok := token.Header["kid"].(string)
        if !ok {
            return nil, jwt.NewValidationError("missing kid header", jwt.ValidationErrorClaimsInvalid)
        }

        switch keyID {
        case "operator":
            return []byte(os.Getenv("JWT_SECRET_OPERATOR")), nil
        default:
            return nil, jwt.NewValidationError("invalid kid", jwt.ValidationErrorClaimsInvalid)
        }
    })

    if err != nil || !token.Valid {
        return c.Status(fiber.StatusUnauthorized).JSON(ErrorResponse{
            Message: "unauthenticated",
            Error:   "Invalid or expired token",
        })
    }

    // Ambil operatorID dari claims
    claims, ok := token.Claims.(*jwt.StandardClaims)
    if !ok {
        return c.Status(fiber.StatusUnauthorized).JSON(ErrorResponse{
            Message: "unauthenticated",
            Error:   "Invalid token claims",
        })
    }

    operatorID := claims.Subject
    var operator models.Operator
    if err := db.DB.Where("id = ?", operatorID).First(&operator).Error; err != nil {
        return c.Status(fiber.StatusNotFound).JSON(ErrorResponse{
            Message: "operator not found",
            Error:   "No operator with the given ID",
        })
    }

    // Parse daftar order IDs yang ingin di-reject
    var orderIDs []int
    if err := c.BodyParser(&orderIDs); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
            Message: "invalid input",
            Error:   "Failed to parse order IDs",
        })
    }

    // Update status setiap order menjadi "Rejected"
    for _, orderID := range orderIDs {
        if err := db.DB.Model(&models.Invoice{}).Where("id = ?", orderID).Update("status", "Rejected").Error; err != nil {
            return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
                Message: "failed to update status",
                Error:   "Failed to update status for order ID " + strconv.Itoa(orderID),
            })
        }
    }

    return c.JSON(fiber.Map{
        "message": "Orders successfully rejected",
    })
}






