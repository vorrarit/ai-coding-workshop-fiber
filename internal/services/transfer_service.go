package services

import (
	"errors"
	"fiber-api/internal/models"

	"gorm.io/gorm"
)

type TransferService struct {
	db *gorm.DB
}

func NewTransferService(db *gorm.DB) *TransferService {
	return &TransferService{db: db}
}

func (s *TransferService) TransferPoints(fromUserID uint, req models.TransferRequest) (*models.TransferResponse, error) {
	// Start transaction
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Get sender user
	var fromUser models.User
	if err := tx.First(&fromUser, fromUserID).Error; err != nil {
		tx.Rollback()
		return nil, errors.New("failed to get sender information")
	}

	// Check if sender has enough points
	if fromUser.PointBalance < req.Amount {
		tx.Rollback()
		return nil, errors.New("insufficient points")
	}

	// Find recipient user
	var toUser models.User
	if err := tx.Where("lbk_code = ?", req.ToLBKCode).First(&toUser).Error; err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("recipient user not found")
		}
		return nil, errors.New("database error")
	}

	// Check if not transferring to self
	if fromUser.ID == toUser.ID {
		tx.Rollback()
		return nil, errors.New("cannot transfer points to yourself")
	}

	// Update balances
	if err := tx.Model(&fromUser).Update("point_balance", fromUser.PointBalance-req.Amount).Error; err != nil {
		tx.Rollback()
		return nil, errors.New("failed to update sender balance")
	}

	if err := tx.Model(&toUser).Update("point_balance", toUser.PointBalance+req.Amount).Error; err != nil {
		tx.Rollback()
		return nil, errors.New("failed to update recipient balance")
	}

	// Create transfer record
	transfer := models.Transfer{
		FromUserID: fromUser.ID,
		ToUserID:   toUser.ID,
		Amount:     req.Amount,
		Message:    req.Message,
		Status:     "completed",
	}

	if err := tx.Create(&transfer).Error; err != nil {
		tx.Rollback()
		return nil, errors.New("failed to create transfer record")
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return nil, errors.New("failed to complete transfer")
	}

	// Prepare response
	response := &models.TransferResponse{
		TransferID: transfer.ID,
		Message:    "Transfer completed successfully",
		FromUser: struct {
			LBKCode   string `json:"lbk_code"`
			FirstName string `json:"first_name"`
			LastName  string `json:"last_name"`
		}{
			LBKCode:   fromUser.LBKCode,
			FirstName: fromUser.FirstName,
			LastName:  fromUser.LastName,
		},
		ToUser: struct {
			LBKCode   string `json:"lbk_code"`
			FirstName string `json:"first_name"`
			LastName  string `json:"last_name"`
		}{
			LBKCode:   toUser.LBKCode,
			FirstName: toUser.FirstName,
			LastName:  toUser.LastName,
		},
		Amount: req.Amount,
		Status: "completed",
	}

	return response, nil
}

func (s *TransferService) GetTransferHistory(userID uint) (*models.TransferHistoryResponse, error) {
	var transfers []models.Transfer
	if err := s.db.Preload("FromUser").Preload("ToUser").
		Where("from_user_id = ? OR to_user_id = ?", userID, userID).
		Order("created_at DESC").
		Limit(50). // Limit to last 50 transfers
		Find(&transfers).Error; err != nil {
		return nil, errors.New("failed to get transfer history")
	}

	return &models.TransferHistoryResponse{
		Transfers: transfers,
		Count:     len(transfers),
	}, nil
}
