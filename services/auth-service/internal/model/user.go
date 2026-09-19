package model

import "time"

// User represents the domain model for an authenticated user account.
// Owned exclusively by Auth Service and stored in the `users` MongoDB collection.
type User struct {
	ID           string    `json:"id" bson:"_id,omitempty"`
	Name         string    `json:"name" bson:"name"`
	Email        string    `json:"email" bson:"email"`
	PasswordHash string    `json:"-" bson:"password_hash"`
	Role         string    `json:"role" bson:"role"`
	CreatedAt    time.Time `json:"createdAt" bson:"created_at"`
	UpdatedAt    time.Time `json:"updatedAt" bson:"updated_at"`
}
