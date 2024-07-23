package main

import (
	"github.com/gofiber/fiber"
	"github.com/sandipbera35/jwt_authservice/controllers"
	"github.com/sandipbera35/jwt_authservice/database"
)

func init() {
	database.DbConnect()
}
func main() {
	println("Server strated .....!")

	app := fiber.New()

	route := app.Group("/api/v1")

	route.Post("/register", func(c *fiber.Ctx) {
		err := controllers.Register(c)

		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return
		}
		c.Status(fiber.StatusOK)
	})
	route.Patch("/upload/profile/image", func(c *fiber.Ctx) {
		err := controllers.AddUploadProfilePic(c)

		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return
		}
		c.Status(fiber.StatusOK)
	})
	route.Get("/get/profile/image", func(c *fiber.Ctx) {
		err := controllers.GetProfilePic(c)

		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return
		}
		c.Status(fiber.StatusOK)
	})
	route.Get("/get/profile/image/by/id", func(c *fiber.Ctx) {
		err := controllers.GetPublicProfilePicById(c)

		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return
		}
		c.Status(fiber.StatusOK)
	})
	route.Patch("/update/profile", func(c *fiber.Ctx) {
		err := controllers.UpdateProfileDetails(c)

		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return
		}
		c.Status(fiber.StatusOK)
	})
	route.Post("/login", func(c *fiber.Ctx) {
		err := controllers.Login(c)

		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			c.JSON("Something went wrong")
			return
		}
		c.Status(fiber.StatusOK)
	})
	route.Get("/profile", func(c *fiber.Ctx) {
		err := controllers.GetProfile(c)

		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			c.JSON("Something went wrong")
			return
		}
		c.Status(fiber.StatusOK)
	})
	app.Listen(":8091")

}
