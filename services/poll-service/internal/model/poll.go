package model

import "time"

// Poll represents a poll in the Votify system.
// This model is owned exclusively by the Poll Service.
type Poll struct {
	ID          string     `json:"id" bson:"_id,omitempty"`
	Title       string     `json:"title" bson:"title"`
	Description string     `json:"description,omitempty" bson:"description,omitempty"`
	Options     []Option   `json:"options" bson:"options"`
	CreatorID   string     `json:"creatorId" bson:"creator_id"`
	Status      string     `json:"status" bson:"status"`
	ExpiresAt   *time.Time `json:"expiresAt,omitempty" bson:"expires_at,omitempty"`
	CreatedAt   time.Time  `json:"createdAt" bson:"created_at"`
	UpdatedAt   time.Time  `json:"updatedAt" bson:"updated_at"`
}

// Option represents a single choice within a poll.
type Option struct {
	ID   string `json:"id" bson:"id"`
	Text string `json:"text" bson:"text"`
}

// Poll statuses
const (
	StatusDraft     = "draft"
	StatusPublished = "published"
	StatusClosed    = "closed"
)
