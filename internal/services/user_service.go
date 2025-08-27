package services

import (
	"errors"
	"fiber-api/internal/models"
	"fiber-api/internal/utils"
	"time"

	"gorm.io/gorm"
)

type UserService struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{db: db}
}

func (s *UserService) CreateUser(req models.RegisterRequest) (*models.User, error) {
	// Check if user already exists
	var existingUser models.User
	if err := s.db.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		return nil, errors.New("user already exists")
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	// Parse DOB
	var dob time.Time
	if req.DOB != "" {
		dob, err = time.Parse("2006-01-02", req.DOB)
		if err != nil {
			return nil, errors.New("invalid date format. Use YYYY-MM-DD")
		}
	}

	// Create user
	user := models.User{
		Email:        req.Email,
		Password:     hashedPassword,
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		PhoneNumber:  req.PhoneNumber,
		DOB:          dob,
		LBKCode:      utils.GenerateLBKCode(),
		PointBalance: 1000, // Give new users 1000 points to start
	}

	if err := s.db.Create(&user).Error; err != nil {
		return nil, errors.New("failed to create user")
	}

	return &user, nil
}

func (s *UserService) AuthenticateUser(req models.LoginRequest) (*models.User, error) {
	var user models.User
	if err := s.db.Where("email = ?", req.Email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("invalid credentials")
		}
		return nil, errors.New("database error")
	}

	// Check password
	if !utils.CheckPassword(req.Password, user.Password) {
		return nil, errors.New("invalid credentials")
	}

	return &user, nil
}

func (s *UserService) GetUserByID(userID uint) (*models.User, error) {
	var user models.User
	if err := s.db.First(&user, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, errors.New("database error")
	}
	return &user, nil
}

func (s *UserService) SearchUserByLBK(lbkCode string) (*models.User, error) {
	var user models.User
	if err := s.db.Select("lbk_code, first_name, last_name").Where("lbk_code = ?", lbkCode).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, errors.New("database error")
	}
	return &user, nil
}
