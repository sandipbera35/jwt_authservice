package controllers

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber"
	"github.com/golang-jwt/jwt/v4"
	"gorm.io/gorm/clause"

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
	rExpiresAt := jwt.NewNumericDate(time.Now().Add(time.Hour * 24))
	rIssuedAt := jwt.NewNumericDate(time.Now())
	token, err := GenerateJWT(user, ExpiresAt, IssuedAt)
	if err != nil {
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Could not login",
		})
		return
	}
	refresh_token, err := GenerateJWT(user, rExpiresAt, rIssuedAt)
	if err != nil {
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Could not login",
		})
		return
	}

	c.JSON(fiber.Map{
		"user_id":           user.ID,
		"expires_at":        ExpiresAt,
		"issued_at":         IssuedAt,
		"access_token":      token,
		"refresh_token_ex":  rExpiresAt,
		"refresh_token_iss": rIssuedAt,
		"refresh_token":     refresh_token,
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

	var user models.User
	dbQ := database.Connect.Where("id = ?", claims.UserId).Preload(clause.Associations).First(&user)

	if dbQ.Error != nil {
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Could not get profile",
		})
		return
	}
	// claims.ProfileImage = user.ProfileImage
	// claims.CoverImage = user.CoverImage

	c.Status(fiber.StatusOK).JSON(user)
}

func UpdateProfileDetails(c *fiber.Ctx) {

	token := c.Get("Authorization")
	// fmt.Printf("token: %v\n", token)
	claims, err := VerifyJWT(token)
	if err != nil {
		c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
		return
	}
	var profile models.User
	// fmt.Printf("profile: %v\n", claims)
	findQ := database.Connect.Where("id = ?", claims.UserId).First(&profile)

	if findQ.Error != nil {
		c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  fiber.StatusNotFound,
			"message": "Profile not found",
			"data":    nil,
		})
		return
	}

	var usermodel models.UserUiModel
	if err := c.BodyParser(&usermodel); err != nil {
		c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
		return
	}

	fmt.Printf("usermodel: %v\n", usermodel)

	profile.FirstName = usermodel.FirstName
	profile.LastName = usermodel.LastName
	profile.Gender = usermodel.Gender
	profile.BirthDate = usermodel.BirthDate

	saveQ := database.Connect.Where("id = ?", claims.UserId).Save(&profile)

	if saveQ.Error != nil {
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  fiber.StatusInternalServerError,
			"message": "Failed to update profile details",
			"data":    nil,
		})
		return
	}

	var profile_image models.ProfileImage
	ProfileImageQ := database.Connect.Where("user_id = ?", claims.UserId).First(&profile_image)
	if ProfileImageQ.RowsAffected == 0 {
		c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  fiber.StatusNotFound,
			"message": "Profile picture not found",
			"data":    nil,
		})
		return
	}
	var cover_image models.CoverImage
	CoverImageQ := database.Connect.Where("user_id = ?", claims.UserId).First(&cover_image)
	if CoverImageQ.RowsAffected == 0 {
		c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  fiber.StatusNotFound,
			"message": "Profile picture not found",
			"data":    nil,
		})
		return
	}

	// file.IsPublic = ispub
	fileSaveQ := database.Connect.Model(models.ProfileImage{}).Where("id = ?", profile_image.ID).UpdateColumn("is_public", usermodel.ProfilePicStatus)

	if fileSaveQ.Error != nil {
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  fiber.StatusInternalServerError,
			"message": "Failed to update profile details",
			"data":    nil,
		})
		return
	}
	coverSaveQ := database.Connect.Model(models.CoverImage{}).Where("id = ?", cover_image.ID).UpdateColumn("is_public", usermodel.CoverPicStatus)

	if coverSaveQ.Error != nil {
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  fiber.StatusInternalServerError,
			"message": "Failed to update profile details",
			"data":    nil,
		})
		return
	}

	GetProfileQ := database.Connect.Where("id = ?", claims.UserId).Preload(clause.Associations).Find(&profile)
	if GetProfileQ.RowsAffected == 0 {
		c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  fiber.StatusInternalServerError,
			"message": "Internal server errors",
			"data":    nil,
		})
		return
	}

	c.Status(fiber.StatusOK).JSON(profile)
}
