package models

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

type TransferHistoryResponse struct {
	Transfers []Transfer `json:"transfers"`
	Count     int        `json:"count"`
}
