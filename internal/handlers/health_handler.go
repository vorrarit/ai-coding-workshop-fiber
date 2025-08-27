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
// @Summary Hello World
// @Description Get hello world message
// @Tags Health
// @Accept json
// @Produce json
// @Success 200 {object} models.HelloResponse
// @Router /api/hello [get]
func (h *HealthHandler) Hello(c *fiber.Ctx) error {
	response := models.HelloResponse{
		Message: "hello world",
	}
	return c.JSON(response)
}

// Health check endpoint
// @Summary Health Check
// @Description Check if the server is healthy
// @Tags Health
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string
// @Router /health [get]
func (h *HealthHandler) Health(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"status": "ok",
	})
}
