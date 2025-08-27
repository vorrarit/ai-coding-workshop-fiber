package handlers

import (
	"fiber-api/internal/models"
	"fiber-api/internal/services"

	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	userService *services.UserService
}

func NewUserHandler(userService *services.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// Get user profile endpoint
// @Summary Get User Profile
// @Description Get authenticated user's profile information
// @Tags User
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.User
// @Failure 401 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /me [get]
func (h *UserHandler) GetMe(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	user, err := h.userService.GetUserByID(userID)
	if err != nil {
		if err.Error() == "user not found" {
			return c.Status(404).JSON(models.ErrorResponse{Error: err.Error()})
		}
		return c.Status(500).JSON(models.ErrorResponse{Error: err.Error()})
	}

	return c.JSON(user)
}

// Get point balance endpoint
// @Summary Get Point Balance
// @Description Get authenticated user's point balance
// @Tags User
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.PointBalanceResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /points/balance [get]
func (h *UserHandler) GetPointBalance(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	user, err := h.userService.GetUserByID(userID)
	if err != nil {
		if err.Error() == "user not found" {
			return c.Status(404).JSON(models.ErrorResponse{Error: err.Error()})
		}
		return c.Status(500).JSON(models.ErrorResponse{Error: err.Error()})
	}

	response := models.PointBalanceResponse{
		LBKCode:      user.LBKCode,
		PointBalance: user.PointBalance,
		FirstName:    user.FirstName,
		LastName:     user.LastName,
	}

	return c.JSON(response)
}

// Search user by LBK code endpoint
// @Summary Search User by LBK Code
// @Description Find a user by their LBK identification code
// @Tags User
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param lbk_code query string true "LBK identification code"
// @Success 200 {object} models.UserSearchResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /users/search [get]
func (h *UserHandler) SearchUserByLBK(c *fiber.Ctx) error {
	lbkCode := c.Query("lbk_code")
	if lbkCode == "" {
		return c.Status(400).JSON(models.ErrorResponse{Error: "lbk_code query parameter is required"})
	}

	user, err := h.userService.SearchUserByLBK(lbkCode)
	if err != nil {
		if err.Error() == "user not found" {
			return c.Status(404).JSON(models.ErrorResponse{Error: err.Error()})
		}
		return c.Status(500).JSON(models.ErrorResponse{Error: err.Error()})
	}

	response := models.UserSearchResponse{
		LBKCode:   user.LBKCode,
		FirstName: user.FirstName,
		LastName:  user.LastName,
	}

	return c.JSON(response)
}
