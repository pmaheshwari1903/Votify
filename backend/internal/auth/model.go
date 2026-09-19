package auth

import "time"

// User represents the domain model for an authenticated user account.
type User struct {
	ID           string    `json:"id" bson:"_id,omitempty"`
	Name         string    `json:"name" bson:"name"`
	Email        string    `json:"email" bson:"email"`
	PasswordHash string    `json:"-" bson:"password_hash"`
	Role         string    `json:"role" bson:"role"`
	OIDCIssuer   string    `json:"oidcIssuer,omitempty" bson:"oidc_issuer,omitempty"`
	OIDCSubject  string    `json:"oidcSubject,omitempty" bson:"oidc_subject,omitempty"`
	CreatedAt    time.Time `json:"createdAt" bson:"created_at"`
	UpdatedAt    time.Time `json:"updatedAt" bson:"updated_at"`
}
