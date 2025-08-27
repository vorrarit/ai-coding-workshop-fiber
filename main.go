package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// JWT Secret key - in production, use environment variable
var jwtSecret = []byte("your-secret-key-change-this-in-production")

// Database instance
var db *gorm.DB

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

// Request structures
type RegisterRequest struct {
	Email       string `json:"email" validate:"required,email"`
	Password    string `json:"password" validate:"required,min=6"`
	FirstName   string `json:"first_name" validate:"required"`
	LastName    string `json:"last_name" validate:"required"`
	PhoneNumber string `json:"phone_number"`
	DOB         string `json:"dob"` // Format: "2006-01-02"
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type TransferRequest struct {
	ToLBKCode string `json:"to_lbk_code" validate:"required"`
	Amount    uint   `json:"amount" validate:"required,min=1"`
	Message   string `json:"message"`
}

// Response structures
type HelloResponse struct {
	Message string `json:"message"`
}

type LoginResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type PointBalanceResponse struct {
	LBKCode      string `json:"lbk_code"`
	PointBalance uint   `json:"point_balance"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
}

type TransferResponse struct {
	TransferID uint   `json:"transfer_id"`
	Message    string `json:"message"`
	FromUser   struct {
		LBKCode   string `json:"lbk_code"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
	} `json:"from_user"`
	ToUser struct {
		LBKCode   string `json:"lbk_code"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
	} `json:"to_user"`
	Amount uint   `json:"amount"`
	Status string `json:"status"`
}

type UserSearchResponse struct {
	LBKCode   string `json:"lbk_code"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

// JWT Claims
type Claims struct {
	UserID uint   `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

// Initialize database
func initDatabase() {
	var err error
	db, err = gorm.Open(sqlite.Open("users.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Auto migrate the schema
	err = db.AutoMigrate(&User{}, &Transfer{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}
}

// Generate LBK code
func generateLBKCode() string {
	// Generate a random 6-digit number for LBK code
	// In production, you might want a more sophisticated generation logic
	timestamp := time.Now().Unix()
	return fmt.Sprintf("LBK%06d", timestamp%1000000)
}

// Hash password
func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// Check password
func checkPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// Generate JWT token
func generateToken(userID uint, email string) (string, error) {
	claims := Claims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// JWT Middleware
func jwtMiddleware(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(401).JSON(ErrorResponse{Error: "Missing authorization header"})
	}

	// Check if the header starts with "Bearer "
	if !strings.HasPrefix(authHeader, "Bearer ") {
		return c.Status(401).JSON(ErrorResponse{Error: "Invalid authorization header format"})
	}

	// Extract the token
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	// Parse and validate the token
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil {
		return c.Status(401).JSON(ErrorResponse{Error: "Invalid token"})
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		// Store user info in context
		c.Locals("userID", claims.UserID)
		c.Locals("email", claims.Email)
		return c.Next()
	}

	return c.Status(401).JSON(ErrorResponse{Error: "Invalid token"})
}

// Register endpoint
func register(c *fiber.Ctx) error {
	var req RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(ErrorResponse{Error: "Invalid request body"})
	}

	// Basic validation
	if req.Email == "" || req.Password == "" || req.FirstName == "" || req.LastName == "" {
		return c.Status(400).JSON(ErrorResponse{Error: "Email, password, first_name, and last_name are required"})
	}

	// Check if user already exists
	var existingUser User
	if err := db.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		return c.Status(400).JSON(ErrorResponse{Error: "User already exists"})
	}

	// Hash password
	hashedPassword, err := hashPassword(req.Password)
	if err != nil {
		return c.Status(500).JSON(ErrorResponse{Error: "Failed to hash password"})
	}

	// Parse DOB
	var dob time.Time
	if req.DOB != "" {
		dob, err = time.Parse("2006-01-02", req.DOB)
		if err != nil {
			return c.Status(400).JSON(ErrorResponse{Error: "Invalid date format. Use YYYY-MM-DD"})
		}
	}

	// Create user
	user := User{
		Email:        req.Email,
		Password:     hashedPassword,
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		PhoneNumber:  req.PhoneNumber,
		DOB:          dob,
		LBKCode:      generateLBKCode(),
		PointBalance: 1000, // Give new users 1000 points to start
	}

	if err := db.Create(&user).Error; err != nil {
		return c.Status(500).JSON(ErrorResponse{Error: "Failed to create user"})
	}

	// Generate token
	token, err := generateToken(user.ID, user.Email)
	if err != nil {
		return c.Status(500).JSON(ErrorResponse{Error: "Failed to generate token"})
	}

	return c.JSON(LoginResponse{Token: token, User: user})
}

// Login endpoint
func login(c *fiber.Ctx) error {
	var req LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(ErrorResponse{Error: "Invalid request body"})
	}

	// Find user
	var user User
	if err := db.Where("email = ?", req.Email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(401).JSON(ErrorResponse{Error: "Invalid credentials"})
		}
		return c.Status(500).JSON(ErrorResponse{Error: "Database error"})
	}

	// Check password
	if !checkPassword(req.Password, user.Password) {
		return c.Status(401).JSON(ErrorResponse{Error: "Invalid credentials"})
	}

	// Generate token
	token, err := generateToken(user.ID, user.Email)
	if err != nil {
		return c.Status(500).JSON(ErrorResponse{Error: "Failed to generate token"})
	}

	return c.JSON(LoginResponse{Token: token, User: user})
}

// Get user profile endpoint
func getMe(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	var user User
	if err := db.First(&user, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(404).JSON(ErrorResponse{Error: "User not found"})
		}
		return c.Status(500).JSON(ErrorResponse{Error: "Database error"})
	}

	return c.JSON(user)
}

// Get point balance endpoint
func getPointBalance(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	var user User
	if err := db.Select("lbk_code, point_balance, first_name, last_name").First(&user, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(404).JSON(ErrorResponse{Error: "User not found"})
		}
		return c.Status(500).JSON(ErrorResponse{Error: "Database error"})
	}

	response := PointBalanceResponse{
		LBKCode:      user.LBKCode,
		PointBalance: user.PointBalance,
		FirstName:    user.FirstName,
		LastName:     user.LastName,
	}

	return c.JSON(response)
}

// Search user by LBK code endpoint
func searchUserByLBK(c *fiber.Ctx) error {
	lbkCode := c.Query("lbk_code")
	if lbkCode == "" {
		return c.Status(400).JSON(ErrorResponse{Error: "lbk_code query parameter is required"})
	}

	var user User
	if err := db.Select("lbk_code, first_name, last_name").Where("lbk_code = ?", lbkCode).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(404).JSON(ErrorResponse{Error: "User not found"})
		}
		return c.Status(500).JSON(ErrorResponse{Error: "Database error"})
	}

	response := UserSearchResponse{
		LBKCode:   user.LBKCode,
		FirstName: user.FirstName,
		LastName:  user.LastName,
	}

	return c.JSON(response)
}

// Transfer points endpoint
func transferPoints(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	var req TransferRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(ErrorResponse{Error: "Invalid request body"})
	}

	// Basic validation
	if req.ToLBKCode == "" || req.Amount == 0 {
		return c.Status(400).JSON(ErrorResponse{Error: "to_lbk_code and amount are required"})
	}

	// Start transaction
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Get sender user
	var fromUser User
	if err := tx.First(&fromUser, userID).Error; err != nil {
		tx.Rollback()
		return c.Status(500).JSON(ErrorResponse{Error: "Failed to get sender information"})
	}

	// Check if sender has enough points
	if fromUser.PointBalance < req.Amount {
		tx.Rollback()
		return c.Status(400).JSON(ErrorResponse{Error: "Insufficient points"})
	}

	// Find recipient user
	var toUser User
	if err := tx.Where("lbk_code = ?", req.ToLBKCode).First(&toUser).Error; err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(404).JSON(ErrorResponse{Error: "Recipient user not found"})
		}
		return c.Status(500).JSON(ErrorResponse{Error: "Database error"})
	}

	// Check if not transferring to self
	if fromUser.ID == toUser.ID {
		tx.Rollback()
		return c.Status(400).JSON(ErrorResponse{Error: "Cannot transfer points to yourself"})
	}

	// Update balances
	if err := tx.Model(&fromUser).Update("point_balance", fromUser.PointBalance-req.Amount).Error; err != nil {
		tx.Rollback()
		return c.Status(500).JSON(ErrorResponse{Error: "Failed to update sender balance"})
	}

	if err := tx.Model(&toUser).Update("point_balance", toUser.PointBalance+req.Amount).Error; err != nil {
		tx.Rollback()
		return c.Status(500).JSON(ErrorResponse{Error: "Failed to update recipient balance"})
	}

	// Create transfer record
	transfer := Transfer{
		FromUserID: fromUser.ID,
		ToUserID:   toUser.ID,
		Amount:     req.Amount,
		Message:    req.Message,
		Status:     "completed",
	}

	if err := tx.Create(&transfer).Error; err != nil {
		tx.Rollback()
		return c.Status(500).JSON(ErrorResponse{Error: "Failed to create transfer record"})
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return c.Status(500).JSON(ErrorResponse{Error: "Failed to complete transfer"})
	}

	// Prepare response
	response := TransferResponse{
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

	return c.JSON(response)
}

// Get transfer history endpoint
func getTransferHistory(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	var transfers []Transfer
	if err := db.Preload("FromUser").Preload("ToUser").
		Where("from_user_id = ? OR to_user_id = ?", userID, userID).
		Order("created_at DESC").
		Limit(50). // Limit to last 50 transfers
		Find(&transfers).Error; err != nil {
		return c.Status(500).JSON(ErrorResponse{Error: "Failed to get transfer history"})
	}

	return c.JSON(fiber.Map{
		"transfers": transfers,
		"count":     len(transfers),
	})
}

func main() {
	// Initialize database
	initDatabase()

	// Get JWT secret from environment if available
	if secret := os.Getenv("JWT_SECRET"); secret != "" {
		jwtSecret = []byte(secret)
	}

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

	// Authentication endpoints
	app.Post("/register", register)
	app.Post("/login", login)

	// Protected endpoints
	app.Get("/me", jwtMiddleware, getMe)
	app.Get("/points/balance", jwtMiddleware, getPointBalance)
	app.Get("/users/search", jwtMiddleware, searchUserByLBK)
	app.Post("/points/transfer", jwtMiddleware, transferPoints)
	app.Get("/points/history", jwtMiddleware, getTransferHistory)
	app.Get("/point-balance", jwtMiddleware, getPointBalance)
	app.Get("/search-user", jwtMiddleware, searchUserByLBK)
	app.Post("/transfer-points", jwtMiddleware, transferPoints)
	app.Get("/transfer-history", jwtMiddleware, getTransferHistory)

	// Start server on port 3000
	log.Printf("Server starting on port 3000...")
	log.Fatal(app.Listen(":3000"))
}
