package controllers

import (
	"time"

	"github.com/gofiber/fiber"
	"github.com/golang-jwt/jwt/v4"

	"github.com/sandipbera35/jwt_authservice/database"
	"github.com/sandipbera35/jwt_authservice/models"
	"golang.org/x/crypto/bcrypt"
)

func Login(c *fiber.Ctx) {
	data := new(struct {
		EmailID  string `json:"email_id"`
		Password string `json:"password"`
	})

	if err := c.BodyParser(data); err != nil {
		c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
		return
	}

	var user models.User
	if err := database.Connect.Where("email_id = ?", data.EmailID).First(&user).Error; err != nil {
		c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid credentials",
		})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.UserPassword), []byte(data.Password)); err != nil {
		c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid credentials",
		})
		return
	}
	ExpiresAt := jwt.NewNumericDate(time.Now().Add(time.Hour * 1))
	IssuedAt := jwt.NewNumericDate(time.Now())
	token, err := GenerateJWT(user, ExpiresAt, IssuedAt)
	if err != nil {
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Could not login",
		})
		return
	}

	c.JSON(fiber.Map{
		"user_id":    user.ID,
		"expires_at": ExpiresAt,
		"issued_at":  IssuedAt,
		"token":      token,
	})
}

func GetProfile(c *fiber.Ctx) {
	token := c.Get("Authorization")
	claims, err := VerifyJWT(token)
	if err != nil {
		c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
		return
	}

	var file models.File
	if err := database.Connect.Where("user_id = ?", claims.UserId).Where("type = ?", "profile").First(&file).Error; err != nil {
		c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Profile not found",
		})
		return
	}
	claims.ProfileImage = &file

	c.Status(fiber.StatusOK).JSON(claims)
}
