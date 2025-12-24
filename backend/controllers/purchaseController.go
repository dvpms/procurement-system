package controllers

import (
	"procurement-system/database"
	"procurement-system/models"
	"time"

	"bytes"
	"encoding/json"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type PurchaseItemRequest struct {
	ItemID uint `json:"item_id" validate:"required,gt=0"` // ID harus ada
	Qty    int  `json:"qty" validate:"required,min=1"`    // Qty minimal 1
}

type CreatePurchaseRequest struct {
	SupplierID uint                  `json:"supplier_id" validate:"required,gt=0"`
	Items      []PurchaseItemRequest `json:"items" validate:"required,min=1,dive"` // Minimal beli 1 barang
}

// --- CONTROLLER ---

func CreatePurchase(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint) // Ambil dari Middleware

	// 1. Parsing Input
	req := new(CreatePurchaseRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid JSON format"})
	}

	// 2. Validation (Clean Code & Security)
	if err := Validate.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Validation failed",
			"errors":  err.Error(),
		})
	}

	// 3. Begin Transaction (ACID)
	tx := database.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var grandTotal int64 = 0
	var details []models.PurchasingDetail

	// 4. Logic Processing
	for _, itemReq := range req.Items {
		var item models.Item

		// Lock Row (Concurrency Control)
		if err := tx.Set("gorm:query_option", "FOR UPDATE").First(&item, itemReq.ItemID).Error; err != nil {
			tx.Rollback()
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "Item not found or unavailable"})
		}

		// Hitung Subtotal (Server Side Calculation - WAJIB)
		subTotal := int64(itemReq.Qty) * item.Price
		grandTotal += subTotal

		details = append(details, models.PurchasingDetail{
			ItemID:   item.ID,
			Qty:      itemReq.Qty,
			SubTotal: subTotal,
		})

		// Update Stock
		item.Stock += itemReq.Qty
		if err := tx.Save(&item).Error; err != nil {
			tx.Rollback()
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to update stock"})
		}
	}

	// 5. Create Header
	purchase := models.Purchasing{
		Date:       time.Now(),
		SupplierID: req.SupplierID,
		UserID:     userID,
		GrandTotal: grandTotal,
		Details:    details,
	}

	if err := tx.Create(&purchase).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to create transaction"})
	}

	// Commit Transaction
	tx.Commit()

	//Kirim Webhook (Non-blocking / Async via Goroutine)
	go sendWebhook(purchase)

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Transaction created successfully",
		"data":    purchase,
	})
}

// --- WEBHOOK HELPER ---
func sendWebhook(purchase models.Purchasing) {
	// 1. URL Dummy (Bisa pakai webhook.site untuk tes nyata, atau placeholder saja)
	url := "https://webhook.site/uuid-anda-disini" // Ganti jika ingin tes live, biarkan jika demo code

	// 2. Siapkan Payload JSON
	payload, _ := json.Marshal(map[string]interface{}{
		"event":       "PURCHASE_CREATED",
		"purchase_id": purchase.ID,
		"total":       purchase.GrandTotal,
		"date":        purchase.Date,
		"items_count": len(purchase.Details),
	})

	// 3. Kirim Request (Background)
	// Kita abaikan error karena ini hanya notifikasi fire-and-forget
	http.Post(url, "application/json", bytes.NewBuffer(payload))
}
