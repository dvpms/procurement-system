package routes

import (
	"github.com/gofiber/fiber/v2"
	"procurement-system/controllers"
	"procurement-system/middleware"
)

func Setup(app *fiber.App) {
	api := app.Group("/api")

	api.Post("/register", controllers.Register)
	api.Post("/login", controllers.Login)
	
	// Protected Routes
	app.Use(middleware.IsAuthenticated)
	
	api.Post("/purchase", controllers.CreatePurchase)
}