package controllers

import (
	"os"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/raihan1405/go-restapi/db"
	"github.com/raihan1405/go-restapi/models"
	"github.com/raihan1405/go-restapi/validators"
	"golang.org/x/crypto/bcrypt"
)

var secretKey = os.Getenv("JWT_SECRET")

// Definisikan struktur respons sukses
type LoginResponse struct {
	Message string `json:"message"`
	Token   string `json:"token"`
	User    struct {
		ID          uint   `json:"id"`
		Username    string `json:"username"`
		Email       string `json:"email"`
		PhoneNumber string `json:"phoneNumber"`
	} `json:"user"`
}

// Definisikan struktur respons kesalahan
type ErrorResponse struct {
	Message string `json:"message"`
	Error   string `json:"error"`
}

// GetUser godoc
// @Summary Get authenticated user details
// @Description Get details of the authenticated user based on the JWT token
// @Tags user
// @Produce json
// @Success 200 {object} models.User
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /api/user [get]
func GetUser(c *fiber.Ctx) error {
	cookie := c.Cookies("jwt")
	token, err := jwt.ParseWithClaims(cookie, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(ErrorResponse{"unauthenticated", err.Error()})
	}
	claims := token.Claims.(*jwt.StandardClaims)

	var user models.User
	db.DB.Where("id = ?", claims.Issuer).First(&user)
	if user.ID == 0 {
		return c.Status(fiber.StatusNotFound).JSON(ErrorResponse{"user not found", "No user with the given ID"})
	}

	return c.JSON(user)
}

// UpdatePassword godoc
// @Summary Update user password
// @Description Update user password with the provided old and new passwords
// @Tags user
// @Accept json
// @Produce json
// @Param update body validators.UpdatePasswordInput true "User update password details"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /api/user/password [put]
func UpdatePassword(c *fiber.Ctx) error {
	cookie := c.Cookies("jwt")
	token, err := jwt.ParseWithClaims(cookie, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(ErrorResponse{"unauthenticated", err.Error()})
	}
	claims := token.Claims.(*jwt.StandardClaims)

	var data validators.UpdatePasswordInput
	err = c.BodyParser(&data)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{"Cannot parse JSON", err.Error()})
	}

	err = validators.Validate.Struct(data)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{"Validation error", err.Error()})
	}

	var user models.User
	db.DB.Where("id = ?", claims.Issuer).First(&user)
	if user.ID == 0 {
		return c.Status(fiber.StatusNotFound).JSON(ErrorResponse{"user not found", "No user with the given ID"})
	}

	// Verify old password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(data.OldPassword))
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(ErrorResponse{"incorrect old password", err.Error()})
	}

	// Generate new hashed password
	newPassword, err := bcrypt.GenerateFromPassword([]byte(data.NewPassword), 14)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{"Cannot hash new password", err.Error()})
	}

	user.Password = newPassword
	db.DB.Save(&user)

	return c.JSON(map[string]interface{}{"message": "password updated successfully"})
}

// UpdateProfile godoc
// @Summary Update user details
// @Description Update user details with the provided information
// @Tags user
// @Accept json
// @Produce json
// @Param update body validators.UpdateUserInput true "User update details"
// @Success 200 {object} models.User
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /api/user [put]
func UpdateProfile(c *fiber.Ctx) error {
	cookie := c.Cookies("jwt")
	token, err := jwt.ParseWithClaims(cookie, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(ErrorResponse{"unauthenticated", err.Error()})
	}
	claims := token.Claims.(*jwt.StandardClaims)

	var data validators.UpdateUserInput
	err = c.BodyParser(&data)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{"Cannot parse JSON", err.Error()})
	}

	err = validators.Validate.Struct(data)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{"Validation error", err.Error()})
	}

	var user models.User
	db.DB.Where("id = ?", claims.Issuer).First(&user)
	if user.ID == 0 {
		return c.Status(fiber.StatusNotFound).JSON(ErrorResponse{"user not found", "No user with the given ID"})
	}

	user.Username = data.Username
	user.Email = data.Email
	user.PhoneNumber = data.PhoneNumber

	db.DB.Save(&user)

	return c.JSON(user)
}

// Register godoc
// @Summary Register a new user
// @Description Register a new user with the provided details
// @Tags auth
// @Accept json
// @Produce json
// @Param register body validators.RegisterInput true "User registration details"
// @Success 200 {object} models.User
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/register [post]
func Register(c *fiber.Ctx) error {
	var data validators.RegisterInput

	// Parse data into the structure
	err := c.BodyParser(&data)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{"Cannot parse JSON", err.Error()})
	}

	// Validate input data
	err = validators.Validate.Struct(data)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{"Validation error", err.Error()})
	}

	// Generate hashed password
	password, err := bcrypt.GenerateFromPassword([]byte(data.Password), 14)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{"Cannot hash password", err.Error()})
	}

	// Create user
	user := models.User{
		Username:    data.Username,
		Email:       data.Email,
		PhoneNumber: data.PhoneNumber,
		Password:    password,
	}

	// Save user to database
	db.DB.Create(&user)
	return c.JSON(user)
}

// Login godoc
// @Summary Log in a user
// @Description Log in a user with the provided credentials and return user data
// @Tags auth
// @Accept json
// @Produce json
// @Param login body validators.LoginInput true "User login details"
// @Success 200 {object} LoginResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/login [post]
func Login(c *fiber.Ctx) error {
	var data validators.LoginInput

	// Parse JSON data
	err := c.BodyParser(&data)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{"Cannot parse JSON", err.Error()})
	}

	// Validate input
	err = validators.Validate.Struct(data)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{"Validation error", err.Error()})
	}

	// Find user by email
	var user models.User
	db.DB.Where("email = ?", data.Email).First(&user)

	// Check if user exists
	if user.ID == 0 {
		return c.Status(fiber.StatusNotFound).JSON(ErrorResponse{"User not found", "No user with the given email"})
	}

	// Compare password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(data.Password))
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(ErrorResponse{"Incorrect password", err.Error()})
	}

	// Generate JWT token
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Subject:   strconv.Itoa(int(user.ID)),
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(), // Token expires in 24 hours
	})

	token, err := claims.SignedString([]byte(secretKey))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{"Could not login", err.Error()})
	}

	// Set cookie with JWT token
	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    token,
		Expires:  time.Now().Add(time.Hour * 24),
		HTTPOnly: true,
	}
	c.Cookie(&cookie)

	// Return user data along with the token
	return c.JSON(LoginResponse{
		Message: "Login successful",
		Token:   token,
		User: struct {
			ID          uint   `json:"id"`
			Username    string `json:"username"`
			Email       string `json:"email"`
			PhoneNumber string `json:"phoneNumber"`
		}{
			ID:           uint(user.ID),
			Username:    user.Username,
			Email:       user.Email,
			PhoneNumber: user.PhoneNumber,
		},
	})
}

// Logout godoc
// @Summary Log out the authenticated user
// @Description Log out the authenticated user by clearing the JWT cookie
// @Tags auth
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /api/logout [post]
func Logout(c *fiber.Ctx) error {
	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
	}

	c.Cookie(&cookie)

	return c.JSON(map[string]interface{}{"message": "logout success"})
}
