package controllers

import (
	"fmt"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/sandipbera35/jwt_authservice/database"
	"github.com/sandipbera35/jwt_authservice/models"
	"gorm.io/gorm"
)

var Roles []string = []string{"SUPERUSER", "ADMIN", "EDITOR"}

func AddAdmin(c *fiber.Ctx) error {

	key := os.Getenv("ADMINKEY")

	if key != strings.TrimSpace(c.Get("Authorization")) {
		c.Status(fiber.StatusUnauthorized)
		c.JSON("Unauthorized")
		return nil
	}

	userID := c.FormValue("user_id")
	role := c.FormValue("role")

	if role == "" {
		c.Status(fiber.StatusBadRequest)
		c.JSON("Role is required")
		return nil
	}
	if !IsValidRole(role) {
		c.Status(fiber.StatusBadRequest)
		c.JSON("Invalid role")
		return nil
	}

	var findsuper models.Admin

	superQ := database.Connect.Model(&models.Admin{}).Where("role = ?", "SUPERUSER").Find(&findsuper)

	if superQ.RowsAffected > 0 {
		c.Status(fiber.StatusBadRequest)
		c.JSON("Superuser already exists more than one superuser is not allowed")
		return nil
	}
	user := models.User{}
	userQ := database.Connect.Where("id = ?", userID).Find(&user)

	if userQ.Error != nil {
		c.Status(fiber.StatusNotFound)
		c.JSON("User not found")
		return nil
	}
	var admin models.Admin
	adminQ := database.Connect.Where("user_id = ?", user.ID).Find(&admin)

	if adminQ.RowsAffected > 0 {
		c.Status(fiber.StatusBadRequest)
		msg := fmt.Sprintf("%v already exists", admin.Role)
		c.JSON(msg)
		return nil
	}

	admin.ID = uuid.New()
	admin.Role = role
	admin.UserID = user.ID

	createAdminQ := database.Connect.Model(&models.Admin{}).Create(&admin)
	if createAdminQ.Error != nil {
		c.Status(fiber.StatusInternalServerError)
		c.JSON("Internal Server error")
		return nil
	}
	if createAdminQ.RowsAffected == 0 {
		c.Status(fiber.StatusInternalServerError)
		c.JSON("Something went wrong")
		return nil
	}

	c.Status(fiber.StatusOK)
	msg := fmt.Sprintf("%v added successfully", role)
	c.JSON(msg)
	return nil
}

func GetAdmins(c *fiber.Ctx) error {
	key := os.Getenv("ADMINKEY")

	if key != strings.TrimSpace(c.Get("Authorization")) {
		c.Status(fiber.StatusUnauthorized)
		c.JSON("Unauthorized")
		return nil
	}
	admins := []models.Admin{}
	database.Connect.Preload("User", func(db *gorm.DB) *gorm.DB {
		return db.Preload("ProfileImage").Preload("CoverImage")

	}).Find(&admins)
	c.Status(fiber.StatusOK)
	c.JSON(admins)
	return nil
}

func IsValidRole(role string) bool {
	for _, v := range Roles {
		if strings.ToUpper(role) == v {
			return true
		}
	}
	return false
}
