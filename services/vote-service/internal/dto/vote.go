package dto

// CastVoteRequest represents the data required to cast a vote.
type CastVoteRequest struct {
	PollID   string `json:"pollId" validate:"required"`
	OptionID string `json:"optionId" validate:"required"`
}

// VoteResponse is the public confirmation of a cast vote.
type VoteResponse struct {
	ID        string `json:"id"`
	PollID    string `json:"pollId"`
	OptionID  string `json:"optionId"`
	CreatedAt string `json:"createdAt"`
}

// TallyResponse contains the aggregated vote counts for a poll.
type TallyResponse struct {
	PollID  string              `json:"pollId"`
	Results []OptionTallyResult `json:"results"`
	Total   int64               `json:"total"`
}

// OptionTallyResult is the vote count for a single option.
type OptionTallyResult struct {
	OptionID string  `json:"optionId"`
	Count    int64   `json:"count"`
	Percent  float64 `json:"percent"`
}
