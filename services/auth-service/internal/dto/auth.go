package dto

// RegisterRequest represents the data required to register a new user.
// Validation tags define the structural constraints that the backend
// enforces independently of any frontend validation.
type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email,max=255"`
	Name     string `json:"name" validate:"required,min=2,max=100"`
	Password string `json:"password" validate:"required,min=8,max=128"`
}

// LoginRequest represents the data required to log in.
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// AuthResponse is returned after successful authentication.
// It contains the JWT token and basic user info.
type AuthResponse struct {
	Token string       `json:"token"`
	User  UserResponse `json:"user"`
}

// UserResponse is the public-facing representation of a user.
// It excludes sensitive fields like password.
type UserResponse struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
	Role  string `json:"role"`
}
