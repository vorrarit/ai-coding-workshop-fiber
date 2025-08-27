package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

// Response structure for the hello endpoint
type HelloResponse struct {
	Message string `json:"message"`
}

func main() {
	// Create a new Fiber instance
	app := fiber.New(fiber.Config{
		AppName: "Fiber API Server v1.0.0",
	})

	// Add middleware
	app.Use(cors.New())

	// Define the /api/hello endpoint
	app.Get("/api/hello", func(c *fiber.Ctx) error {
		response := HelloResponse{
			Message: "hello world",
		}
		return c.JSON(response)
	})

	// Health check endpoint
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "ok",
		})
	})

	// Start server on port 3000
	log.Fatal(app.Listen(":3000"))
}
