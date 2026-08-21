package middleware

import (
	"os"
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
)

func JWTProtected(c fiber.Ctx) error {

	authHeader := c.Get("Authorization")

	if authHeader == "" {
		return c.Status(401).JSON(fiber.Map{
			"error": "Token gerekli",
		})
	}

	parts := strings.Split(authHeader, " ")

	if len(parts) != 2 || parts[0] != "Bearer" {
		return c.Status(401).JSON(fiber.Map{
			"error": "Gecersiz token",
		})
	}

	tokenString := parts[1]

	token, err := jwt.Parse(
		tokenString,
		func(token *jwt.Token) (interface{}, error) {

			if token.Method != jwt.SigningMethodHS256 {
				return nil, jwt.ErrSignatureInvalid
			}

			return []byte(os.Getenv("JWT_SECRET")), nil
		},
	)

	if err != nil || !token.Valid {
		return c.Status(401).JSON(fiber.Map{
			"error": "Gecersiz veya suresi dolmus token",
		})
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok {
		return c.Status(401).JSON(fiber.Map{
			"error": "Token bilgileri okunamadi",
		})
	}

	userID, ok := claims["user_id"].(float64)

	if !ok {
		return c.Status(401).JSON(fiber.Map{
			"error": "Token icinde kullanici bilgisi yok",
		})
	}

	c.Locals("userID", uint(userID))

	return c.Next()
}