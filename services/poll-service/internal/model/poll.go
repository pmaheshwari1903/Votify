package model

import "time"

// PollOption represents a single voting option in a poll.
type PollOption struct {
	ID        string `json:"id" bson:"id"`
	Text      string `json:"text" bson:"text"`
	VoteCount int    `json:"voteCount" bson:"vote_count"`
}

// Poll represents the core domain model for a poll.
// Owned exclusively by Poll Service and stored in the `polls` MongoDB collection.
type Poll struct {
	ID          string       `json:"id" bson:"_id,omitempty"`
	OwnerID     string       `json:"ownerId" bson:"owner_id"`
	Question    string       `json:"question" bson:"question"`
	Description string       `json:"description,omitempty" bson:"description,omitempty"`
	Options     []PollOption `json:"options" bson:"options"`
	Status      string       `json:"status" bson:"status"` // "draft", "open", "closed"
	TotalVotes  int          `json:"totalVotes" bson:"total_votes"`
	CreatedAt   time.Time    `json:"createdAt" bson:"created_at"`
	UpdatedAt   time.Time    `json:"updatedAt" bson:"updated_at"`
}
