package models

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
