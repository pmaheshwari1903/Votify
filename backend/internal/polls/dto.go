package polls

import "time"

// CreatePollRequest represents request data to create a new poll.
type CreatePollRequest struct {
	Question    string   `json:"question" binding:"required,min=5,max=300"`
	Description string   `json:"description,omitempty" binding:"max=1000"`
	Options     []string `json:"options" binding:"required,min=2,max=10"`
}

// UpdatePollRequest represents request data to edit an existing poll.
type UpdatePollRequest struct {
	Question    string `json:"question,omitempty" binding:"omitempty,min=5,max=300"`
	Description string `json:"description,omitempty" binding:"omitempty,max=1000"`
}

// PublicPollResponse represents sanitized data for public voting.
type PublicPollResponse struct {
	ID          string       `json:"id"`
	Question    string       `json:"question"`
	Description string       `json:"description,omitempty"`
	Options     []PollOption `json:"options"`
	Status      string       `json:"status"`
	CreatedAt   time.Time    `json:"createdAt"`
}

// PollResultsResponse represents poll vote counts for display.
type PollResultsResponse struct {
	PollID     string       `json:"pollId"`
	Question   string       `json:"question"`
	Status     string       `json:"status"`
	TotalVotes int          `json:"totalVotes"`
	Options    []PollOption `json:"options"`
}
