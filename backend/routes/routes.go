package routes

import (
	"procurement-system/controllers"
	"procurement-system/middleware"

	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App) {
	api := app.Group("/api")

	// --- Public Routes ---
	api.Post("/register", controllers.Register)
	api.Post("/login", controllers.Login)

	// --- Protected Routes (Harus Login) ---
	// Semua route di bawah ini akan melewati middleware Auth
	protected := api.Group("/", middleware.IsAuthenticated)

	protected.Post("/purchase", controllers.CreatePurchase)

	// 2. Master Data: Items
	protected.Get("/items", controllers.GetItems)          // Staff & Admin boleh lihat
	protected.Post("/items", controllers.CreateItem)       // Admin only
	protected.Put("/items/:id", controllers.UpdateItem)    // Admin only
	protected.Delete("/items/:id", controllers.DeleteItem) // Admin only

	// 3. Master Data: Suppliers
	protected.Get("/suppliers", controllers.GetSuppliers)
	protected.Post("/suppliers", controllers.CreateSupplier)
}
