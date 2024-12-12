package controllers

import (
	"strconv"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/raihan1405/go-restapi/db"
	"github.com/raihan1405/go-restapi/models"
	"os"
)

// OperatorDashboard godoc
// @Summary Get operator dashboard
// @Description Get dashboard data for the authenticated operator
// @Tags operator
// @Produce json
// @Success 200 {object} models.Operator
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /api/operator/dashboard [get]
func OperatorDashboard(c *fiber.Ctx) error {
    // Ambil token dari cookie
    cookie := c.Cookies("jwt_operator")

    // Parse token dengan claims
    token, err := jwt.ParseWithClaims(cookie, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
        // Use multiple keys if necessary
        keyID, ok := token.Header["kid"].(string)
        if !ok {
            return nil, jwt.NewValidationError("missing kid header", jwt.ValidationErrorClaimsInvalid)
        }

        // Return the appropriate signing key based on the kid
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
        return c.Status(fiber.StatusUnauthorized).JSON(ErrorResponse{
            Message: "unauthenticated",
            Error:   "Invalid or expired token",
        })
    }

    // Ambil claims
    claims, ok := token.Claims.(*jwt.StandardClaims)
    if !ok {
        return c.Status(fiber.StatusUnauthorized).JSON(ErrorResponse{
            Message: "unauthenticated",
            Error:   "Invalid token claims",
        })
    }

    // Convert Subject to integer (Operator ID)
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

    // Kembalikan data operator
    return c.JSON(operator)
}
