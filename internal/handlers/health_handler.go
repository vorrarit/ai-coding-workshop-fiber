package handlers

import (
	"fiber-api/internal/models"

	"github.com/gofiber/fiber/v2"
)

type HealthHandler struct{}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

// Hello endpoint
func (h *HealthHandler) Hello(c *fiber.Ctx) error {
	response := models.HelloResponse{
		Message: "hello world",
	}
	return c.JSON(response)
}

// Health check endpoint
func (h *HealthHandler) Health(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"status": "ok",
	})
}
