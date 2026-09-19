package dto

import "time"

// CastVoteRequest represents the request payload to submit a vote.
type CastVoteRequest struct {
	PollID   string `json:"pollId" binding:"required"`
	OptionID string `json:"optionId" binding:"required"`
}

// VoteResponse is returned after a vote is successfully recorded.
type VoteResponse struct {
	ID        string    `json:"id"`
	PollID    string    `json:"pollId"`
	OptionID  string    `json:"optionId"`
	CreatedAt time.Time `json:"createdAt"`
}
