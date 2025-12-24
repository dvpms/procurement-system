package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"procurement-system/database"
	"procurement-system/routes"
)

func main() {
	database.Connect()

	app := fiber.New()

	// CORS agar Frontend jQuery bisa akses
	app.Use(cors.New(cors.Config{
		AllowCredentials: true,
		AllowOrigins:     "http://localhost:5500", // Sesuaikan port frontend nanti
	}))

	routes.Setup(app)

	app.Listen(":3000")
}