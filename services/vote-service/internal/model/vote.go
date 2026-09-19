package model

import "time"

// Vote represents a recorded vote in the `votes` MongoDB collection.
type Vote struct {
	ID            string    `json:"id" bson:"_id,omitempty"`
	PollID        string    `json:"pollId" bson:"poll_id"`
	OptionID      string    `json:"optionId" bson:"option_id"`
	VoterIdentity string    `json:"voterIdentity" bson:"voter_identity"` // User ID or IP hash
	VoterIP       string    `json:"voterIp,omitempty" bson:"voter_ip,omitempty"`
	UserID        string    `json:"userId,omitempty" bson:"user_id,omitempty"`
	CreatedAt     time.Time `json:"createdAt" bson:"created_at"`
}
