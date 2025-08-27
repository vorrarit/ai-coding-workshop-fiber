package handlers

import (
	"fiber-api/internal/models"
	"fiber-api/internal/services"
	"fiber-api/internal/utils"

	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	userService *services.UserService
	jwtSecret   []byte
}

func NewAuthHandler(userService *services.UserService, jwtSecret []byte) *AuthHandler {
	return &AuthHandler{
		userService: userService,
		jwtSecret:   jwtSecret,
	}
}

// Register endpoint
// @Summary User Registration
// @Description Register a new user account
// @Tags Authentication
// @Accept json
// @Produce json
// @Param user body models.RegisterRequest true "User registration details"
// @Success 200 {object} models.LoginResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /register [post]
func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req models.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(models.ErrorResponse{Error: "Invalid request body"})
	}

	// Basic validation
	if req.Email == "" || req.Password == "" || req.FirstName == "" || req.LastName == "" {
		return c.Status(400).JSON(models.ErrorResponse{Error: "Email, password, first_name, and last_name are required"})
	}

	user, err := h.userService.CreateUser(req)
	if err != nil {
		return c.Status(400).JSON(models.ErrorResponse{Error: err.Error()})
	}

	// Generate token
	token, err := utils.GenerateToken(user.ID, user.Email, h.jwtSecret)
	if err != nil {
		return c.Status(500).JSON(models.ErrorResponse{Error: "Failed to generate token"})
	}

	return c.JSON(models.LoginResponse{Token: token, User: *user})
}

// Login endpoint
// @Summary User Login
// @Description Authenticate user and return JWT token
// @Tags Authentication
// @Accept json
// @Produce json
// @Param credentials body models.LoginRequest true "User login credentials"
// @Success 200 {object} models.LoginResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /login [post]
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req models.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(models.ErrorResponse{Error: "Invalid request body"})
	}

	user, err := h.userService.AuthenticateUser(req)
	if err != nil {
		return c.Status(401).JSON(models.ErrorResponse{Error: err.Error()})
	}

	// Generate token
	token, err := utils.GenerateToken(user.ID, user.Email, h.jwtSecret)
	if err != nil {
		return c.Status(500).JSON(models.ErrorResponse{Error: "Failed to generate token"})
	}

	return c.JSON(models.LoginResponse{Token: token, User: *user})
}
