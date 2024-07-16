package controllers

import (
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/raihan1405/go-restapi/db"
	"github.com/raihan1405/go-restapi/models"
	"golang.org/x/crypto/bcrypt"
)

const secretKey = "secret"

func Register(c *fiber.Ctx) error {
	var data map[string] string

	err := c.BodyParser(&data)

	if err != nil {
		return err
	}

	password, _ := bcrypt.GenerateFromPassword([]byte(data["password"]),14)
	user:= models.User{
		Username : data["username"],
		Email : data["email"],
		Password: password,
	}

	db.DB.Create(&user)
	return c.JSON(user)
}

func Login(c *fiber.Ctx) error{
	var data map[string] string

	err := c.BodyParser(&data)

	if err != nil {
		return err
	}

	var user models.User
	db.DB.Where("email = ?",data["email"]).First(&user)

	if user.ID == 0 {
		c.Status(fiber.StatusNotFound)
		return c.JSON(fiber.Map{
			"message" : "user not found",
		})
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password),[]byte(data["password"]))
	
	if err != nil{
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message" : "incorrect password",
		})
	}

	claims := jwt.NewWithClaims(jwt.SigningMethodES256,jwt.StandardClaims{
		Issuer: strconv.Itoa(int(user.ID)),
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
	})

	token, err := claims.SignedString([]byte(secretKey))

	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message" : "could not login",
		})
	}

	cookie := fiber.Cookie{
		Name : "jwt",
		Value : token,
		Expires : time.Now().Add(time.Hour * 24),
		HTTPOnly: true,
	}
	c.Cookie(&cookie)

	return c.JSON(fiber.Map{
		"message" : "success",
	})
}

func User(c *fiber.Ctx) error{
	cookie := c.Cookies("jwt")
	token,err:=jwt.ParseWithClaims(cookie,&jwt.StandardClaims{},func(token *jwt.Token)(interface{},error){
		return []byte(secretKey),nil
	})

	if err!= nil{
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message" : "unauthenticated",
		})
	}
	claims := token.Claims.(jwt.StandardClaims)

	var user models.User
	db.DB.Where("id = ?",claims.Issuer).First(&user)

	return c.JSON(user)
}

func Logout(c *fiber.Ctx) error{
	cookie := fiber.Cookie{
		Name : "jwt",
		Value : "",
		Expires : time.Now().Add(-time.Hour),
		HTTPOnly: true,
	}

	c.Cookie(&cookie)

	return c.JSON(fiber.Map{
		"message" : "logout success",
	})
}