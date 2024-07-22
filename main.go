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

	app.Post("/register", func(c *fiber.Ctx) {
		err := controllers.Register(c)

		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return
		}
		c.Status(fiber.StatusOK)
	})
	app.Post("/login", func(c *fiber.Ctx) {
		err := controllers.Login(c)

		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			c.JSON("Something went wrong")
			return
		}
		c.Status(fiber.StatusOK)
	})
	app.Get("/profile", func(c *fiber.Ctx) {
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
