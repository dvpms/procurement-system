package controllers

import (
	"procurement-system/database"
	"procurement-system/models"

	"github.com/gofiber/fiber/v2"
)

type ItemInput struct {
	Name  string `json:"name" validate:"required,min=3"`
	Stock int    `json:"stock" validate:"required,min=0"`
	Price int64  `json:"price" validate:"required,min=1"`
}

type SupplierInput struct {
	Name    string `json:"name" validate:"required,min=3"`
	Email   string `json:"email" validate:"required,email"`
	Address string `json:"address" validate:"required"`
}

// --- HELPER FUNCTION ---

func isAdmin(c *fiber.Ctx) bool {
	role := c.Locals("user_role")
	return role == "admin"
}

// ==========================================
// 📦 ITEMS CONTROLLER
// ==========================================

// GET /api/items (Read) - Bisa diakses Staff/Admin
func GetItems(c *fiber.Ctx) error {
	var items []models.Item
	// Urutkan dari yang stoknya paling sedikit (biar kelihatan mana yang butuh restock)
	database.DB.Order("stock asc").Find(&items)
	return c.JSON(fiber.Map{
		"message": "Success",
		"data":    items,
	})
}

// POST /api/items (Create) - ADMIN ONLY
func CreateItem(c *fiber.Ctx) error {
	if !isAdmin(c) {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"message": "Forbidden: Only Admin can add items"})
	}

	input := new(ItemInput)
	if err := c.BodyParser(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid JSON"})
	}

	// Pakai Shared Validator
	if err := Validate.Struct(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"errors": err.Error()})
	}

	newItem := models.Item{
		Name:  input.Name,
		Stock: input.Stock,
		Price: input.Price,
	}

	database.DB.Create(&newItem)
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "Item added", "data": newItem})
}

// PUT /api/items/:id (Update) - ADMIN ONLY
func UpdateItem(c *fiber.Ctx) error {
	if !isAdmin(c) {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"message": "Forbidden"})
	}

	id := c.Params("id")
	var item models.Item

	if err := database.DB.First(&item, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "Item not found"})
	}

	input := new(ItemInput)
	if err := c.BodyParser(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid JSON"})
	}

	if err := Validate.Struct(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"errors": err.Error()})
	}

	// Update data
	item.Name = input.Name
	item.Stock = input.Stock
	item.Price = input.Price

	database.DB.Save(&item)
	return c.JSON(fiber.Map{"message": "Item updated", "data": item})
}

// DELETE /api/items/:id (Delete) - ADMIN ONLY
func DeleteItem(c *fiber.Ctx) error {
	if !isAdmin(c) {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"message": "Forbidden"})
	}

	id := c.Params("id")
	var item models.Item
	if err := database.DB.First(&item, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "Item not found"})
	}

	database.DB.Delete(&item)
	return c.JSON(fiber.Map{"message": "Item deleted"})
}

// ==========================================
// 🚚 SUPPLIERS CONTROLLER
// ==========================================

func GetSuppliers(c *fiber.Ctx) error {
	var suppliers []models.Supplier
	database.DB.Order("name asc").Find(&suppliers)
	return c.JSON(fiber.Map{"data": suppliers})
}

func CreateSupplier(c *fiber.Ctx) error {
	if !isAdmin(c) {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"message": "Forbidden"})
	}

	input := new(SupplierInput)
	if err := c.BodyParser(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid JSON"})
	}

	if err := Validate.Struct(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"errors": err.Error()})
	}

	supplier := models.Supplier{
		Name:    input.Name,
		Email:   input.Email,
		Address: input.Address,
	}

	database.DB.Create(&supplier)
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "Supplier added", "data": supplier})
}
