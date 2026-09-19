package model

import "time"

// User represents an authenticated user in the Votify system.
// This model is owned exclusively by the Auth Service.
type User struct {
	ID        string    `json:"id" bson:"_id,omitempty"`
	Email     string    `json:"email" bson:"email"`
	Name      string    `json:"name" bson:"name"`
	Password  string    `json:"-" bson:"password"` // Never serialized to JSON
	Role      string    `json:"role" bson:"role"`
	CreatedAt time.Time `json:"createdAt" bson:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updated_at"`
}

// Roles
const (
	RoleUser  = "user"
	RoleAdmin = "admin"
)
