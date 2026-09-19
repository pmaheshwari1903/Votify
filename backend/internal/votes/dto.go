package votes

import "time"

// CastVoteRequest represents the payload to submit a vote for a poll option.
type CastVoteRequest struct {
	PollID   string `json:"pollId" binding:"required"`
	OptionID string `json:"optionId" binding:"required"`
}

// VoteResponse represents the confirmation envelope after a vote is recorded.
type VoteResponse struct {
	ID        string    `json:"id"`
	PollID    string    `json:"pollId"`
	OptionID  string    `json:"optionId"`
	CreatedAt time.Time `json:"createdAt"`
}
