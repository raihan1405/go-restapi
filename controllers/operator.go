package controllers

import (
	"strconv"
	"time"
	"log"
	"fmt"
	

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/raihan1405/go-restapi/db"
	"github.com/raihan1405/go-restapi/models"
	"github.com/raihan1405/go-restapi/validators"
	// "golang.org/x/crypto/bcrypt"
	"os"
)

var operatorSecretKey = os.Getenv("JWT_SECRET_OPERATOR") // Pastikan Anda memiliki JWT_SECRET_OPERATOR di .env

// LoginOperatorResponse defines the structure of a successful login response for operator
type LoginOperatorResponse struct {
	Message string        `json:"message"`
	Token   string        `json:"token"`
	Operator OperatorInfo `json:"operator"`
}

type OperatorInfo struct {
	ID         uint   `json:"id"`
	OperatorID string `json:"operator_id"`
	Name       string `json:"name"`
	Email      string `json:"email"`
	Phone      string `json:"phone"`
}


// LoginOperator godoc
// @Summary Log in an operator
// @Description Log in an operator with the provided Operator ID and return operator data
// @Tags auth
// @Accept json
// @Produce json
// @Param login body validators.OperatorLoginInput true "Operator login details"
// @Success 200 {object} LoginOperatorResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/operator/login [post]
func LoginOperator(c *fiber.Ctx) error {
    // Define multiple secret keys
    mapKey := map[string]string{
        "user":    os.Getenv("JWT_SECRET"),
        "operator": os.Getenv("JWT_SECRET_OPERATOR"),
    }

    // Check if the operator secret key is set
    operatorSecretKey, exists := mapKey["operator"]
    if !exists || operatorSecretKey == "" {
        return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
            Message: "Server configuration error",
            Error:   "JWT_SECRET_OPERATOR is not set",
        })
    }

    var data validators.OperatorLoginInput

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

    // Find operator by OperatorID
    var operator models.Operator
    if err := db.DB.Where("operator_id = ?", data.OperatorID).First(&operator).Error; err != nil {
        return c.Status(fiber.StatusNotFound).JSON(models.ErrorResponse{
            Message: "Operator not found",
            Error:   "No operator with the given Operator ID",
        })
    }

    // If you want to use password for operator, add verification here
    /*
        if err := bcrypt.CompareHashAndPassword([]byte(operator.Password), []byte(data.Password)); err != nil {
            return c.Status(fiber.StatusUnauthorized).JSON(models.ErrorResponse{
                Message: "Incorrect Operator ID or Password",
                Error:   "Authentication failed",
            })
        }
    */

    // Create JWT claims
    claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
        Subject:   strconv.Itoa(int(operator.ID)),
        ExpiresAt: time.Now().Add(time.Hour * 24).Unix(), // Token valid for 24 hours
    })

    // Set the key ID for the token
    claims.Header["kid"] = "operator"

    // Sign the token with the selected secret key
    token, err := claims.SignedString([]byte(operatorSecretKey))
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
            Message: "Could not generate token",
            Error:   err.Error(),
        })
    }

    // Set cookie with the JWT token
    cookie := fiber.Cookie{
        Name:     "jwt_operator",
        Value:    token,
        Expires:  time.Now().Add(time.Hour * 24),
        HTTPOnly: true,
        Secure:   false, // Set to true if using HTTPS
        SameSite: "Lax",
        Path:     "/",
    }
    c.Cookie(&cookie)
    log.Printf("jwt_operator cookie set: %s\n", token)

    // Return success response
    return c.JSON(LoginOperatorResponse{
        Message: "Login successful",
        Token:   token,
        Operator: OperatorInfo{
            ID:         operator.ID,
            OperatorID: operator.OperatorID,
            Name:       operator.Name,
            Email:      operator.Email,
            Phone:      operator.Phone,
        },
    })
}

func LogoutOperator(c *fiber.Ctx) error {
    // Hapus cookie jwt_operator dengan mengatur expired date ke masa lalu
	fmt.Println("LogoutOperator called")

    cookie := fiber.Cookie{
        Name:     "jwt_operator",
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

func GetAllProductsOperator(c *fiber.Ctx) error {
	var products []models.Product

	// Retrieve all products from the database
	if err := db.DB.Find(&products).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(map[string]interface{}{"error": "Cannot retrieve products"})
	}

	return c.JSON(products)
}