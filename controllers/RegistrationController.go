package controllers

import (
	"strings"

	"github.com/gofiber/fiber"
	"github.com/google/uuid"
	"github.com/sandipbera35/jwt_authservice/database"
	"github.com/sandipbera35/jwt_authservice/models"
	"golang.org/x/crypto/bcrypt"
)

// var jwtSecret = []byte("supersecretkey")

// Register endpoint
func Register(c *fiber.Ctx) {
	userUiModel := new(models.UserUiModel)

	if err := c.BodyParser(userUiModel); err != nil {
		c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
		return
	}

	user := new(models.User)
	user.FirstName = userUiModel.FirstName
	user.LastName = userUiModel.LastName
	user.Gender = userUiModel.Gender
	user.BirthDate = userUiModel.BirthDate
	user.UserName = userUiModel.UserName
	user.EmailID = userUiModel.EmailID
	user.MobileNo = userUiModel.MobileNo

	user.ID = uuid.New()

	chkQ := database.Connect.Model(models.User{}).Where("user_name = ?", user.UserName).Or("email_id = ?", user.EmailID).Or("mobile_no = ?", user.MobileNo).Find(&user)
	if chkQ.RowsAffected > 0 {
		c.Status(fiber.StatusBadRequest).JSON(fiber.Map{

			"status":  fiber.StatusBadRequest,
			"message": "User already exists with this email or mobile number",
		})
		return
	}
	//check password is vlid or not
	if len(userUiModel.UserPassword) < 6 || strings.TrimSpace(userUiModel.UserPassword) == "" {
		c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"message": "Password must be at least 6 characters",
		})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userUiModel.UserPassword), bcrypt.DefaultCost)
	if err != nil {
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to hash password",
		})
		return
	}
	user.UserPassword = string(hashedPassword)

	if err := database.Connect.Create(&user).Error; err != nil {
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Could not register user",
		})
		return
	}

	// c.Status(fiber.StatusOK)
	c.JSON(fiber.Map{
		"status":  fiber.StatusOK,
		"message": "User registered successfully",
		"data":    user,
	})
}
