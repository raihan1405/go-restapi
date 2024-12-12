package controllers

import (
	"os"
	"strconv"
	"time"
    "fmt"

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
    // Retrieve the JWT token from the cookies
    cookie := c.Cookies("jwt")

    // Parse the JWT token and extract the claims
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

    // Handle error if token is invalid or parsing fails
    if err != nil || !token.Valid {
        return c.Status(fiber.StatusUnauthorized).JSON(ErrorResponse{
            Message: "unauthenticated",
            Error:   "Invalid or expired token",
        })
    }

    // Extract the claims and cast to the correct type
    claims, ok := token.Claims.(*jwt.StandardClaims)
    if !ok {
        return c.Status(fiber.StatusUnauthorized).JSON(ErrorResponse{
            Message: "unauthenticated",
            Error:   "Invalid token claims",
        })
    }

    // Convert the Subject (which is the user ID in string format) to an integer
    userID, err := strconv.Atoi(claims.Subject)
    if err != nil {
        return c.Status(fiber.StatusUnauthorized).JSON(ErrorResponse{
            Message: "Invalid token",
            Error:   "Token contains invalid user ID",
        })
    }

    // Retrieve the user from the database using the converted user ID
    var user models.User
    db.DB.Where("id = ?", userID).First(&user)

    // If user is not found, return a 404 error
    if user.ID == 0 {
        return c.Status(fiber.StatusNotFound).JSON(ErrorResponse{
            Message: "user not found",
            Error:   "No user with the given ID",
        })
    }

    // Return the user details as the response
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

    if err != nil {
        return c.Status(fiber.StatusUnauthorized).JSON(ErrorResponse{"unauthenticated", err.Error()})
    }
    claims := token.Claims.(*jwt.StandardClaims)

    // Convert claims.Subject to integer
    userID, err := strconv.Atoi(claims.Subject)
    if err != nil {
        return c.Status(fiber.StatusUnauthorized).JSON(ErrorResponse{"Invalid token", "Token contains invalid user ID"})
    }

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
    db.DB.Where("id = ?", userID).First(&user)
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

    if err != nil {
        return c.Status(fiber.StatusUnauthorized).JSON(ErrorResponse{"unauthenticated", err.Error()})
    }
    claims := token.Claims.(*jwt.StandardClaims)

    // Convert claims.Subject to integer
    userID, err := strconv.Atoi(claims.Subject)
    if err != nil {
        return c.Status(fiber.StatusUnauthorized).JSON(ErrorResponse{"Invalid token", "Token contains invalid user ID"})
    }

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
    db.DB.Where("id = ?", userID).First(&user)
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

    keyMap := map[string]string{
        "user":    os.Getenv("JWT_SECRET"),
        "operator": os.Getenv("JWT_SECRET_OPERATOR"),
    }

    secretKey, exists := keyMap["user"]
    if !exists || secretKey == "" {
        return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
            Message: "Server configuration error",
            Error:   "JWT_SECRET is not set",
        })
    }

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

    // Set the key ID for the token
    claims.Header["kid"] = "user"

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
        Secure:   true,  // This ensures the cookie is only sent over HTTPS
        SameSite: "Lax",
		Path:     "/", // Allows the cookie to be sent cross-domain
	}

	c.Cookie(&cookie)
    fmt.Printf("jwt_operator cookie set: %s\n", token)

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
			ID:          uint(user.ID),
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
func LogoutUser(c *fiber.Ctx) error {
    fmt.Println("logoutUser called")
    
	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
        SameSite: "Lax",
		Path:     "/",
	}

	c.Cookie(&cookie)

	return c.JSON(map[string]interface{}{"message": "logout success"})

}
