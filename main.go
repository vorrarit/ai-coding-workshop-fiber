package main

import (
	"errors"
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
	ID          uint      `json:"id" gorm:"primarykey"`
	Email       string    `json:"email" gorm:"unique;not null"`
	Password    string    `json:"-" gorm:"not null"` // "-" excludes from JSON
	FirstName   string    `json:"first_name" gorm:"not null"`
	LastName    string    `json:"last_name" gorm:"not null"`
	PhoneNumber string    `json:"phone_number"`
	DOB         time.Time `json:"dob"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
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
	err = db.AutoMigrate(&User{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}
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
		Email:       req.Email,
		Password:    hashedPassword,
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		PhoneNumber: req.PhoneNumber,
		DOB:         dob,
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

	// Start server on port 3000
	log.Printf("Server starting on port 3000...")
	log.Fatal(app.Listen(":3000"))
}
