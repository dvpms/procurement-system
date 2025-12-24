package controllers

import (
	"procurement-system/database"
	"procurement-system/models"
	"time"

	"github.com/gofiber/fiber/v2"
)

// Struct untuk menampung JSON request dari frontend
type CreatePurchaseRequest struct {
	SupplierID uint `json:"supplier_id"`
	Items      []struct {
		ItemID uint `json:"item_id"`
		Qty    int  `json:"qty"`
	} `json:"items"`
}

func CreatePurchase(c *fiber.Ctx) error {
	// 1. Ambil User ID dari Token (disimpan di Locals oleh Middleware)
	userID := c.Locals("user_id").(uint)

	// 2. Parse Request Body
	req := new(CreatePurchaseRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "Invalid Data"})
	}

	// === MEMULAI DATABASE TRANSACTION ===
	tx := database.DB.Begin()

	// Error handling panic recovery (jika ada crash di tengah transaksi)
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var grandTotal int64 = 0
	var details []models.PurchasingDetail

	// 3. Loop items untuk hitung harga SERVER-SIDE
	for _, itemReq := range req.Items {
		var item models.Item
		// Kunci row database agar tidak ada race condition
		if err := tx.Set("gorm:query_option", "FOR UPDATE").First(&item, itemReq.ItemID).Error; err != nil {
			tx.Rollback()
			return c.Status(404).JSON(fiber.Map{"message": "Item tidak ditemukan"})
		}

		// Hitung Subtotal
		subTotal := int64(itemReq.Qty) * item.Price
		grandTotal += subTotal

		// Tambah Detail ke memory slice
		details = append(details, models.PurchasingDetail{
			ItemID:   item.ID,
			Qty:      itemReq.Qty,
			SubTotal: subTotal,
		})

		// 4. Update Stock Item
		item.Stock += itemReq.Qty
		if err := tx.Save(&item).Error; err != nil {
			tx.Rollback()
			return c.Status(500).JSON(fiber.Map{"message": "Gagal update stock"})
		}
	}

	// 5. Buat Header Transaksi
	purchase := models.Purchasing{
		Date:       time.Now(),
		SupplierID: req.SupplierID,
		UserID:     userID,
		GrandTotal: grandTotal,
		Details:    details, // GORM akan otomatis insert details juga karena relasi HasMany
	}

	if err := tx.Create(&purchase).Error; err != nil {
		tx.Rollback()
		return c.Status(500).JSON(fiber.Map{"message": "Gagal simpan transaksi"})
	}

	// === COMMIT TRANSAKSI (Jika semua lancar, simpan permanen) ===
	tx.Commit()

	return c.JSON(fiber.Map{
		"message": "Transaksi Sukses",
		"data":    purchase,
	})
}
