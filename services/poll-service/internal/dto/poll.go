package dto

// CreatePollRequest represents the data required to create a new poll.
// Validation tags enforce structural constraints at the request boundary.
type CreatePollRequest struct {
	Title       string   `json:"title" validate:"required,min=3,max=200"`
	Description string   `json:"description,omitempty" validate:"max=1000"`
	Options     []string `json:"options" validate:"required,min=2,max=10,dive,required,min=1,max=200"`
}

// UpdatePollRequest represents the data for updating an existing poll.
type UpdatePollRequest struct {
	Title       *string  `json:"title,omitempty" validate:"omitempty,min=3,max=200"`
	Description *string  `json:"description,omitempty" validate:"omitempty,max=1000"`
	Options     []string `json:"options,omitempty" validate:"omitempty,min=2,max=10,dive,required,min=1,max=200"`
}

// PollResponse is the public-facing representation of a poll.
type PollResponse struct {
	ID          string           `json:"id"`
	Title       string           `json:"title"`
	Description string           `json:"description,omitempty"`
	Options     []OptionResponse `json:"options"`
	CreatorID   string           `json:"creatorId"`
	Status      string           `json:"status"`
	ExpiresAt   *string          `json:"expiresAt,omitempty"`
	CreatedAt   string           `json:"createdAt"`
}

// OptionResponse is the public-facing representation of a poll option.
type OptionResponse struct {
	ID   string `json:"id"`
	Text string `json:"text"`
}
