package events

import (
	"time"
)

// Standard topic names for Kafka domain event streams.
const (
	TopicVoteCreated         = "votify.votes.created"
	TopicPollCreated         = "votify.polls.created"
	TopicPollPublished       = "votify.polls.published"
	TopicPollClosed          = "votify.polls.closed"
	TopicUserRegistered      = "votify.auth.user-registered"
	TopicSubscriptionCreated = "votify.payments.subscription-created"
	TopicPaymentCompleted    = "votify.payments.completed"
)

// Event type constants.
const (
	EventTypeVoteCreated         = "VoteCreated"
	EventTypePollCreated         = "PollCreated"
	EventTypePollPublished       = "PollPublished"
	EventTypePollClosed          = "PollClosed"
	EventTypeUserRegistered       = "UserRegistered"
	EventTypeSubscriptionCreated  = "SubscriptionCreated"
	EventTypePaymentCompleted    = "PaymentCompleted"
)

// Current event specification version.
const DefaultEventVersion = "1.0"

// Event is the Cloud-ready standard envelope for all Kafka messages.
// It includes versioning, correlation identifiers, and producer metadata.
type Event struct {
	// EventID is a unique identifier (UUIDv4) for deduplication.
	EventID string `json:"event_id"`

	// EventType identifies the domain event (e.g., "VoteCreated").
	EventType string `json:"event_type"`

	// EventVersion specifies the schema version (e.g., "1.0").
	EventVersion string `json:"event_version"`

	// ProducerService identifies the producing microservice (e.g., "vote-service").
	ProducerService string `json:"producer_service"`

	// OccurredAt is the ISO 8601 UTC timestamp when the event occurred.
	OccurredAt string `json:"occurred_at"`

	// PartitionKey is the key used for Kafka partitioning (e.g., poll_id or user_id).
	PartitionKey string `json:"partition_key"`

	// Payload contains the event-specific domain data.
	Payload interface{} `json:"payload"`
}

// NewEvent constructs a standard Event envelope with default metadata.
func NewEvent(eventID, eventType, producerService, partitionKey string, payload interface{}) Event {
	return Event{
		EventID:         eventID,
		EventType:       eventType,
		EventVersion:    DefaultEventVersion,
		ProducerService: producerService,
		OccurredAt:      time.Now().UTC().Format(time.RFC3339),
		PartitionKey:    partitionKey,
		Payload:         payload,
	}
}

// VoteCreatedPayload represents the domain event data when a vote is cast.
type VoteCreatedPayload struct {
	PollID    string `json:"poll_id"`
	OptionID  string `json:"option_id"`
	VoterIP   string `json:"voter_ip,omitempty"`
	UserID    string `json:"user_id,omitempty"`
	Timestamp string `json:"timestamp"`
}

// PollCreatedPayload represents the domain event data when a poll is created.
type PollCreatedPayload struct {
	PollID    string `json:"poll_id"`
	CreatorID string `json:"creator_id"`
	Title     string `json:"title"`
	Timestamp string `json:"timestamp"`
}

// UserRegisteredPayload represents the domain event data when a new user registers.
type UserRegisteredPayload struct {
	UserID    string `json:"user_id"`
	Email     string `json:"email"`
	Timestamp string `json:"timestamp"`
}
