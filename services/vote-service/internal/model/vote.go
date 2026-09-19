package model

import "time"

// Vote represents a single vote cast on a poll option.
// This model is owned exclusively by the Vote Service.
type Vote struct {
	ID        string    `json:"id" bson:"_id,omitempty"`
	PollID    string    `json:"pollId" bson:"poll_id"`
	OptionID  string    `json:"optionId" bson:"option_id"`
	VoterID   string    `json:"voterId,omitempty" bson:"voter_id,omitempty"` // Empty for anonymous votes
	IPHash    string    `json:"-" bson:"ip_hash"`                            // Hashed IP for duplicate detection
	CreatedAt time.Time `json:"createdAt" bson:"created_at"`
}

// VoteTally represents the aggregated vote count for a poll option.
type VoteTally struct {
	OptionID string `json:"optionId" bson:"option_id"`
	Count    int64  `json:"count" bson:"count"`
}
