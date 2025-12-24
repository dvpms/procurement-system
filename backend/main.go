package main

import (
	"procurement-system/database"
	"procurement-system/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	database.Connect()

	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*", // Boleh semua origin
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))

	routes.Setup(app)

	app.Listen(":3000")
}
