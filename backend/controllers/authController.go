package controllers

import (
	"os"
	"procurement-system/database"
	"procurement-system/models"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func Register(c *fiber.Ctx) error {
	var data map[string]string
	if err := c.BodyParser(&data); err != nil {
		return err
	}

	password, _ := bcrypt.GenerateFromPassword([]byte(data["password"]), 14)
	user := models.User{
		Username: data["username"],
		Password: string(password),
		Role:     "staff", // Default
	}

	if err := database.DB.Create(&user).Error; err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "Gagal register, username mungkin duplikat"})
	}
	return c.JSON(user)
}

func Login(c *fiber.Ctx) error {
	var data map[string]string
	if err := c.BodyParser(&data); err != nil {
		return err
	}

	var user models.User
	// Cari user berdasarkan username
	database.DB.Where("username = ?", data["username"]).First(&user)

	if user.ID == 0 {
		return c.Status(404).JSON(fiber.Map{"message": "User not found"})
	}

	// Cek Password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(data["password"])); err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "Incorrect password"})
	}

	// 1. Ambil Secret Key
	secretKey := os.Getenv("JWT_SECRET")
	if secretKey == "" {
		// Safety check: Beritahu jika .env belum terbaca
		return c.Status(500).JSON(fiber.Map{
			"message": "Server Error: JWT_SECRET is missing in .env",
		})
	}

	// 2. Buat Claims
	claims := jwt.MapClaims{
		"iss": user.ID,
		"exp": time.Now().Add(time.Hour * 24).Unix(), // Expire 1 hari
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 3. Generate Token String (JANGAN pakai tanda _ )
	t, err := token.SignedString([]byte(secretKey))
	if err != nil {
		// Jika gagal, return error agar kita tahu kenapa
		return c.Status(500).JSON(fiber.Map{
			"message": "Could not generate token",
			"error":   err.Error(),
		})
	}

	// Jika sukses, kirim token
	return c.JSON(fiber.Map{
		"token": t,
		"user":  user,
	})
}
