package controllers

import (
	"log"
	"strconv"
	"time"
	//"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/raihan1405/go-restapi/db"
	"github.com/raihan1405/go-restapi/models"
	"github.com/raihan1405/go-restapi/validators"
	"os"
)

var adminSecretKey = os.Getenv("JWT_SECRET_ADMIN") // Pastikan Anda memiliki JWT_SECRET_ADMIN di .env

// LoginAdminResponse defines the structure of a successful login response for admin
type LoginAdminResponse struct {
	Message string    `json:"message"`
	Token   string    `json:"token"`
	Admin   AdminInfo `json:"admin"`
}

type AdminInfo struct {
	ID      uint   `json:"id"`
	AdminID string `json:"admin_id"`
	Name    string `json:"name"`
	Email   string `json:"email"`
	Phone   string `json:"phone"`
}

// LoginAdmin godoc
// @Summary Log in an admin
// @Description Log in an admin with the provided Admin ID and return admin data
// @Tags auth
// @Accept json
// @Produce json
// @Param login body validators.AdminLoginInput true "Admin login details"
// @Success 200 {object} LoginAdminResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/admin/login [post]
func LoginAdmin(c *fiber.Ctx) error {
	// Define multiple secret keys
	mapKey := map[string]string{
		"admin": os.Getenv("JWT_SECRET_ADMIN"),
	}

	// Check if the admin secret key is set
	adminSecretKey, exists := mapKey["admin"]
	if !exists || adminSecretKey == "" {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Message: "Server configuration error",
			Error:   "JWT_SECRET_ADMIN is not set",
		})
	}

	var data validators.AdminLoginInput

	// Parse JSON body
	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Message: "Cannot parse JSON",
			Error:   err.Error(),
		})
	}

	// Validate input
	if err := validators.Validate.Struct(data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Message: "Validation error",
			Error:   err.Error(),
		})
	}

	// Find admin by AdminID
	var admin models.Admin
	if err := db.DB.Where("admin_id = ?", data.AdminID).First(&admin).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(models.ErrorResponse{
			Message: "Admin not found",
			Error:   "No admin with the given Admin ID",
		})
	}

	// If you want to use password for admin, add verification here
	/*
	   if err := bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(data.Password)); err != nil {
	       return c.Status(fiber.StatusUnauthorized).JSON(models.ErrorResponse{
	           Message: "Incorrect Admin ID or Password",
	           Error:   "Authentication failed",
	       })
	   }
	*/

	// Create JWT claims
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Subject:   strconv.Itoa(int(admin.ID)),
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(), // Token valid for 24 hours
	})

	// Set the key ID for the token
	claims.Header["kid"] = "admin"

	// Sign the token with the selected secret key
	token, err := claims.SignedString([]byte(adminSecretKey))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Message: "Could not generate token",
			Error:   err.Error(),
		})
	}

	// Set cookie with the JWT token
	cookie := fiber.Cookie{
		Name:     "jwt_admin",
		Value:    token,
		Expires:  time.Now().Add(time.Hour * 24),
		HTTPOnly: true,
		Secure:   false, // Set to true if using HTTPS
		SameSite: "Lax",
		Path:     "/",
	}
	c.Cookie(&cookie)
	log.Printf("jwt_admin cookie set: %s\n", token)

	// Return success response
	return c.JSON(LoginAdminResponse{
		Message: "Login successful",
		Token:   token,
		Admin: AdminInfo{
			ID:      admin.ID,
			AdminID: admin.AdminID,
			Name:    admin.Name,
			Email:   admin.Email,
			Phone:   admin.Phone,
		},
	})
}

func GetAllInvoicesForAdmin(c *fiber.Ctx) error {
	// Get the JWT token and verify the user
	cookie := c.Cookies("jwt_admin")

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
		case "admin":
			return []byte(os.Getenv("JWT_SECRET_ADMIN")), nil
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
	adminID, err := strconv.Atoi(claims.Subject)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(ErrorResponse{
			Message: "Invalid token",
			Error:   "Token contains invalid operator ID",
		})
	}

	// Cari operator di database
	var operator models.Operator
	if err := db.DB.Where("id = ?", adminID).First(&operator).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(ErrorResponse{
			Message: "operator not found",
			Error:   "No operator with the given ID",
		})
	}

	// Retrieve all invoices for the logged-in user, preload related InvoiceItems, Products, and User
	var invoices []models.Invoice
	if err := db.DB.Preload("InvoiceItems.Product").Preload("User").Find(&invoices).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Message: "failed to retrieve invoices",
			Error:   "Failed to retrieve invoices from the database",
		})
	}

	// Return the list of invoices
	return c.JSON(invoices)
}

func LogoutAdmin(c *fiber.Ctx) error {
    // Hapus cookie jwt_operator dengan mengatur expired date ke masa lalu

    cookie := fiber.Cookie{
        Name:     "jwt_admin",
        Value:    "",
        Expires:  time.Now().Add(-time.Hour), // Set expired
        HTTPOnly: true,
        Secure:   false, // Set ke true jika menggunakan HTTPS di produksi
        SameSite: "Lax",
		Path:     "/",
    }
    c.Cookie(&cookie)

    return c.JSON(map[string]interface{}{"message": "Logout successful"})
}
