package middleware

import (
	"fiber-api/internal/models"
	"fiber-api/internal/utils"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func JWTMiddleware(jwtSecret []byte) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(401).JSON(models.ErrorResponse{Error: "Missing authorization header"})
		}

		// Check if the header starts with "Bearer "
		if !strings.HasPrefix(authHeader, "Bearer ") {
			return c.Status(401).JSON(models.ErrorResponse{Error: "Invalid authorization header format"})
		}

		// Extract the token
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// Parse and validate the token
		claims, err := utils.ParseToken(tokenString, jwtSecret)
		if err != nil {
			return c.Status(401).JSON(models.ErrorResponse{Error: "Invalid token"})
		}

		// Store user info in context
		c.Locals("userID", claims.UserID)
		c.Locals("email", claims.Email)
		return c.Next()
	}
}
