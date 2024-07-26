package main

import (
	"runtime"

	"github.com/gofiber/fiber"
	"github.com/sandipbera35/jwt_authservice/controllers"
	"github.com/sandipbera35/jwt_authservice/database"
)

func init() {
	database.ConnectDatabase()
}
func main() {
	println("Server strated .....!")

	app := fiber.New(&fiber.Settings{
		Concurrency: runtime.NumCPU(),
		BodyLimit:   1024 * 1024 * 1024,
	})

	app.Get("/", func(c *fiber.Ctx) {
		c.Status(fiber.StatusOK)
		c.JSON(map[string]interface{}{
			"message": "Jwt AuthServer API by Sandip Bera",
		})
	})

	adminroute := app.Group("/admin")
	adminroute.Post("/register", controllers.AddAdmin)
	adminroute.Get("/getadmins", controllers.GetAdmins)

	route := app.Group("/api/v1")
	route.Post("/register", controllers.Register)
	route.Patch("/upload/profile/image", controllers.AddUploadProfilePic)
	route.Patch("/upload/cover/image", controllers.AddUploadCoverPic)
	route.Get("/get/profile/image", controllers.GetProfilePic)
	route.Get("/get/cover/image", controllers.GetCoverPic)
	route.Get("/get/profile/image/by/id", controllers.GetPublicProfilePicById)
	route.Get("/get/cover/image/by/id", controllers.GetPublicCoverPicById)
	route.Delete("/delete/profile/image", controllers.DeleteProfilePic)
	route.Delete("/delete/cover/image", controllers.DeleteCoverPic)
	route.Patch("/update/profile", controllers.UpdateProfileDetails)
	route.Post("/login", controllers.Login)
	route.Get("/profile", controllers.GetProfile)

	app.Listen(":8091")

}
