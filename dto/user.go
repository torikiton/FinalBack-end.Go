package dto

import "time"

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}
type LoginResponse struct {
	CustomerID  int
	FirstName   string
	LastName    string
	Email       string
	PhoneNumber string
	Address     string
	Password    string `json:"-"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type UpdateAddressRequest struct {
	CustomerID int    `json:"customer_id" binding:"required"`
	Address    string `json:"address" binding:"required"`
}
type RegisterRequest struct {
	Email       string `json:"email" binding:"required"`
	Password    string `json:"password" binding:"required"`
	FirstName   string `json:"first_name" binding:"required"`
	LastName    string `json:"last_name" binding:"required"`
	Address     string `json:"address" binding:"required"`
	PhoneNumber string `json:"phone_number" binding:"required"`
}
