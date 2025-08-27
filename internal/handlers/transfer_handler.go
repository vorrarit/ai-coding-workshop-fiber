package handlers

import (
	"fiber-api/internal/models"
	"fiber-api/internal/services"

	"github.com/gofiber/fiber/v2"
)

type TransferHandler struct {
	transferService *services.TransferService
}

func NewTransferHandler(transferService *services.TransferService) *TransferHandler {
	return &TransferHandler{
		transferService: transferService,
	}
}

// Transfer points endpoint
func (h *TransferHandler) TransferPoints(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	var req models.TransferRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(models.ErrorResponse{Error: "Invalid request body"})
	}

	// Basic validation
	if req.ToLBKCode == "" || req.Amount == 0 {
		return c.Status(400).JSON(models.ErrorResponse{Error: "to_lbk_code and amount are required"})
	}

	response, err := h.transferService.TransferPoints(userID, req)
	if err != nil {
		switch err.Error() {
		case "insufficient points":
			return c.Status(400).JSON(models.ErrorResponse{Error: err.Error()})
		case "recipient user not found":
			return c.Status(404).JSON(models.ErrorResponse{Error: err.Error()})
		case "cannot transfer points to yourself":
			return c.Status(400).JSON(models.ErrorResponse{Error: err.Error()})
		default:
			return c.Status(500).JSON(models.ErrorResponse{Error: err.Error()})
		}
	}

	return c.JSON(response)
}

// Get transfer history endpoint
func (h *TransferHandler) GetTransferHistory(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	response, err := h.transferService.GetTransferHistory(userID)
	if err != nil {
		return c.Status(500).JSON(models.ErrorResponse{Error: err.Error()})
	}

	return c.JSON(response)
}
