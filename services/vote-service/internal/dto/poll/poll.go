package poll

import "time"

type PublicPollResponse struct {
	ID          string       `json:"id"`
	Question    string       `json:"question"`
	Description string       `json:"description,omitempty"`
	Options     []PollOption `json:"options"`
	Status      string       `json:"status"`
	CreatedAt   time.Time    `json:"createdAt"`
}

type PollOption struct {
	ID        string `json:"id"`
	Text      string `json:"text"`
	VoteCount int    `json:"voteCount"`
}