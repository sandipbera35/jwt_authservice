package controllers

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"gorm.io/gorm/clause"

	"github.com/sandipbera35/jwt_authservice/database"
	"github.com/sandipbera35/jwt_authservice/models"
	"golang.org/x/crypto/bcrypt"
)

func Login(c *fiber.Ctx) error {
	data := new(struct {
		EmailID  string `json:"email_id"`
		Password string `json:"password"`
	})

	if err := c.BodyParser(data); err != nil {
		c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
		return nil
	}

	var user models.User
	if err := database.Connect.Where("email_id = ?", data.EmailID).Find(&user).Error; err != nil {
		c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid credentials",
		})
		return nil
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.UserPassword), []byte(data.Password)); err != nil {
		c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid credentials",
		})
		return nil
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
		return nil
	}
	refresh_token, err := GenerateJWT(user, rExpiresAt, rIssuedAt)
	if err != nil {
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Could not login",
		})
		return nil
	}

	c.JSON(fiber.Map{
		"user_id":           user.ID,
		"expires_at":        ExpiresAt,
		"issued_at":         IssuedAt,
		"access_token":      token,
		"refresh_token_exp": rExpiresAt,
		"refresh_token_iss": rIssuedAt,
		"refresh_token":     refresh_token,
	})
	return nil
}

func GetProfile(c *fiber.Ctx) error {
	token := c.Get("Authorization")
	claims, err := VerifyJWT(token)
	if err != nil {
		c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
		return nil
	}

	var user models.User
	dbQ := database.Connect.Where("id = ?", claims.UserId).Preload(clause.Associations).Find(&user)

	if dbQ.Error != nil {
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Could not get profile",
		})
		return nil
	}
	// claims.ProfileImage = user.ProfileImage
	// claims.CoverImage = user.CoverImage

	c.Status(fiber.StatusOK).JSON(user)
	return nil
}

func UpdateProfileDetails(c *fiber.Ctx) error {

	token := c.Get("Authorization")
	// fmt.Printf("token: %v\n", token)
	claims, err := VerifyJWT(token)
	if err != nil {
		c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
		return nil
	}
	var profile models.User
	// fmt.Printf("profile: %v\n", claims)
	findQ := database.Connect.Where("id = ?", claims.UserId).Find(&profile)

	if findQ.Error != nil {
		c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  fiber.StatusNotFound,
			"message": "Profile not found",
			"data":    nil,
		})
		return nil
	}

	var usermodel models.UserUiModel
	if err := c.BodyParser(&usermodel); err != nil {
		c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
		return nil
	}

	// fmt.Printf("usermodel: %v\n", usermodel)

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
		return nil
	}

	var profile_image models.ProfileImage
	ProfileImageQ := database.Connect.Where("user_id = ?", claims.UserId).Find(&profile_image)
	if ProfileImageQ.RowsAffected == 0 {
		c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  fiber.StatusNotFound,
			"message": "Profile picture not found",
			"data":    nil,
		})
		return nil
	}
	var cover_image models.CoverImage
	CoverImageQ := database.Connect.Where("user_id = ?", claims.UserId).Find(&cover_image)
	if CoverImageQ.RowsAffected == 0 {
		c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  fiber.StatusNotFound,
			"message": "Profile picture not found",
			"data":    nil,
		})
		return nil
	}

	// file.IsPublic = ispub
	fileSaveQ := database.Connect.Model(models.ProfileImage{}).Where("id = ?", profile_image.ID).UpdateColumn("is_public", usermodel.ProfilePicStatus)

	if fileSaveQ.Error != nil {
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  fiber.StatusInternalServerError,
			"message": "Failed to update profile details",
			"data":    nil,
		})
		return nil
	}
	coverSaveQ := database.Connect.Model(models.CoverImage{}).Where("id = ?", cover_image.ID).UpdateColumn("is_public", usermodel.CoverPicStatus)

	if coverSaveQ.Error != nil {
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  fiber.StatusInternalServerError,
			"message": "Failed to update profile details",
			"data":    nil,
		})
		return nil
	}

	GetProfileQ := database.Connect.Where("id = ?", claims.UserId).Preload(clause.Associations).Find(&profile)
	if GetProfileQ.RowsAffected == 0 {
		c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  fiber.StatusInternalServerError,
			"message": "Internal server errors",
			"data":    nil,
		})
		return nil
	}

	c.Status(fiber.StatusOK).JSON(profile)

	return nil
}
