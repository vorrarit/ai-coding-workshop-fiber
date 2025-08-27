package models

import (
	"time"
)

// User model
type User struct {
	ID           uint      `json:"id" gorm:"primarykey"`
	Email        string    `json:"email" gorm:"unique;not null"`
	Password     string    `json:"-" gorm:"not null"` // "-" excludes from JSON
	FirstName    string    `json:"first_name" gorm:"not null"`
	LastName     string    `json:"last_name" gorm:"not null"`
	PhoneNumber  string    `json:"phone_number"`
	DOB          time.Time `json:"dob"`
	LBKCode      string    `json:"lbk_code" gorm:"unique;not null"` // LBK identification code
	PointBalance uint      `json:"point_balance" gorm:"default:0"`  // Point balance
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// Transfer model for point transfers
type Transfer struct {
	ID         uint      `json:"id" gorm:"primarykey"`
	FromUserID uint      `json:"from_user_id" gorm:"not null"`
	ToUserID   uint      `json:"to_user_id" gorm:"not null"`
	FromUser   User      `json:"from_user" gorm:"foreignKey:FromUserID"`
	ToUser     User      `json:"to_user" gorm:"foreignKey:ToUserID"`
	Amount     uint      `json:"amount" gorm:"not null"`
	Message    string    `json:"message"`
	Status     string    `json:"status" gorm:"default:'completed'"` // completed, failed, pending
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
