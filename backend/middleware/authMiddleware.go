package middleware

import (
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func IsAuthenticated(c *fiber.Ctx) error {
	cookie := c.Cookies("jwt")
	if cookie == "" {
		return c.Status(401).JSON(fiber.Map{"message": "Unauthenticated"})
	}

	SecretKey := os.Getenv("JWT_SECRET")
	if SecretKey == "" {
		return c.Status(500).JSON(fiber.Map{"message": "JWT secret not configured"})
	}

	token, err := jwt.ParseWithClaims(cookie, &jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(SecretKey), nil
	})

	if err != nil || !token.Valid {
		return c.Status(401).JSON(fiber.Map{"message": "Unauthenticated"})
	}

	claims := token.Claims.(*jwt.MapClaims)
	id := uint((*claims)["iss"].(float64))

	c.Locals("user_id", id)

	return c.Next()
}
