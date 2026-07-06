package auth

import "stock-trading/internal/profile"

type RegisterRequest struct {
	Name         string                   `json:"name" binding:"required,min=3,max=50"`
	Mobile       string                   `json:"mobile" binding:"required,len=10,numeric"`
	Email        string                   `json:"email" binding:"required,email"`
	Password     string                   `json:"password" binding:"required"`
	PAN          string                   `json:"pan" binding:"required"`
	Aadhaar      string                   `json:"aadhaar" binding:"required"`
	Address      string                   `json:"address" binding:"required"`
	IncomeRange  string                   `json:"income_range" binding:"required"`
	Occupation   string                   `json:"occupation" binding:"required"`
	IPVPhoto     string                   `json:"ipv_photo" binding:"required"`
	IPVLatitude  string                   `json:"ipv_latitude" binding:"required"`
	IPVLongitude string                   `json:"ipv_longitude" binding:"required"`
	DOB          string                   `json:"dob" binding:"required"` // Format: YYYY-MM-DD
	BankAccounts []profile.BankAccountDTO `json:"bank_accounts" binding:"required,min=1"`
	Nominees     []profile.NomineeDTO     `json:"nominees" binding:"omitempty,max=3"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type UserResponse struct {
	ID           string                   `json:"id"`
	Name         string                   `json:"name"`
	Email        string                   `json:"email"`
	Mobile       string                   `json:"mobile"`
	Role         string                   `json:"role"`
	Status       string                   `json:"status"`
	DOB          string                   `json:"dob,omitempty"`
	Address      string                   `json:"address,omitempty"`
	PAN          string                   `json:"pan,omitempty"`
	Aadhaar      string                   `json:"aadhaar,omitempty"`
	IncomeRange  string                   `json:"income_range,omitempty"`
	Occupation   string                   `json:"occupation,omitempty"`
	IPVPhoto     string                   `json:"ipv_photo,omitempty"`
	IPVLatitude  string                   `json:"ipv_latitude,omitempty"`
	IPVLongitude string                   `json:"ipv_longitude,omitempty"`
	FatherName   string                   `json:"father_name,omitempty"`
	MotherName   string                   `json:"mother_name,omitempty"`
	Country      string                   `json:"country,omitempty"`
	State        string                   `json:"state,omitempty"`
	City         string                   `json:"city,omitempty"`
	Pincode      string                   `json:"pincode,omitempty"`
	BankAccounts []profile.BankAccountDTO `json:"bank_accounts,omitempty"`
	Nominees     []profile.NomineeDTO     `json:"nominees,omitempty"`
	CreatedAt    string                   `json:"created_at"`
}

type LoginResponse struct {
	AccessToken  string        `json:"access_token"`
	RefreshToken string        `json:"refresh_token"`
	User         *UserResponse `json:"data"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type VerifyResetTokenRequest struct {
	Token string `json:"token" binding:"required"`
}

type ResetPasswordRequest struct {
	Token    string `json:"token" binding:"required"`
	Password string `json:"password" binding:"required"`
}
