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

	app.Get("/", func(c *fiber.Ctx) {
		c.JSON(map[string]interface{}{
			"message": "Jwt AuthServer API by Sandip Bera",
		})
	})

	route := app.Group("/api/v1")
	route.Post("/register", controllers.Register)
	route.Patch("/upload/profile/image", controllers.AddUploadProfilePic)
	route.Get("/get/profile/image", controllers.GetProfilePic)
	route.Get("/get/profile/image/by/id", controllers.GetPublicProfilePicById)
	route.Patch("/update/profile", controllers.UpdateProfileDetails)
	route.Post("/login", controllers.Login)
	route.Get("/profile", controllers.GetProfile)

	app.Listen(":8091")

}
