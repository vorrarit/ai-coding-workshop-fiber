package utils

import (
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// Hash password
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// Check password
func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// Generate LBK code
func GenerateLBKCode() string {
	// Generate a random 6-digit number for LBK code
	// In production, you might want a more sophisticated generation logic
	timestamp := time.Now().Unix()
	return fmt.Sprintf("LBK%06d", timestamp%1000000)
}
