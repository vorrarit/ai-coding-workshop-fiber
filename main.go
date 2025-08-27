package main

import (
	"fiber-api/internal/config"
	"fiber-api/internal/database"
	"fiber-api/internal/handlers"
	"fiber-api/internal/middleware"
	"fiber-api/internal/services"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Initialize database
	db := database.NewDatabase(cfg.DatabasePath)

	// Initialize services
	userService := services.NewUserService(db.GetDB())
	transferService := services.NewTransferService(db.GetDB())

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(userService, cfg.JWTSecret)
	userHandler := handlers.NewUserHandler(userService)
	transferHandler := handlers.NewTransferHandler(transferService)
	healthHandler := handlers.NewHealthHandler()

	// Create Fiber app
	app := fiber.New(fiber.Config{
		AppName: cfg.AppName,
	})

	// Add middleware
	app.Use(cors.New())

	// Initialize JWT middleware
	jwtMiddleware := middleware.JWTMiddleware(cfg.JWTSecret)

	// Public routes
	app.Get("/api/hello", healthHandler.Hello)
	app.Get("/health", healthHandler.Health)

	// Authentication routes
	app.Post("/register", authHandler.Register)
	app.Post("/login", authHandler.Login)

	// Protected routes
	app.Get("/me", jwtMiddleware, userHandler.GetMe)
	app.Get("/points/balance", jwtMiddleware, userHandler.GetPointBalance)
	app.Get("/users/search", jwtMiddleware, userHandler.SearchUserByLBK)
	app.Post("/points/transfer", jwtMiddleware, transferHandler.TransferPoints)
	app.Get("/points/history", jwtMiddleware, transferHandler.GetTransferHistory)

	// Start server
	log.Printf("Server starting on port %s...", cfg.ServerPort)
	log.Fatal(app.Listen(cfg.ServerPort))
}
