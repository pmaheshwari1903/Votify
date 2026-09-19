package auth

import "time"

// RegisterRequest represents the data required to register a new user.
type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email,max=255"`
	Name     string `json:"name" binding:"required,min=2,max=100"`
	Password string `json:"password" binding:"required,min=8,max=128"`
}

// LoginRequest represents the data required to log in.
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// UserResponse is the public-facing representation of a user.
type UserResponse struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"createdAt"`
}

// AuthResponse is returned upon successful authentication.
type AuthResponse struct {
	Token string       `json:"token"`
	User  UserResponse `json:"user"`
}
