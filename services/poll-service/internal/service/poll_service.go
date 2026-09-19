package service

// PollService defines the interface for poll business logic.
type PollService interface {
	// Create creates a new poll.
	// Business rules: validate options uniqueness, set default status to "draft",
	// assign creator from authenticated context.
	// Create(ctx context.Context, creatorID string, req dto.CreatePollRequest) (*dto.PollResponse, error)

	// GetByID retrieves a poll by its ID.
	// GetByID(ctx context.Context, id string) (*dto.PollResponse, error)

	// Update updates an existing poll.
	// Business rules: only the creator can update, poll must be in "draft" status.
	// Update(ctx context.Context, id string, creatorID string, req dto.UpdatePollRequest) (*dto.PollResponse, error)

	// Delete deletes a poll.
	// Business rules: only the creator can delete.
	// Delete(ctx context.Context, id string, creatorID string) error

	// List retrieves polls with pagination.
	// List(ctx context.Context, page, perPage int) ([]dto.PollResponse, int64, error)

	// Publish transitions a poll from "draft" to "published".
	// Publishes a PollPublished Kafka event.
	// Publish(ctx context.Context, id string, creatorID string) (*dto.PollResponse, error)

	// Close transitions a poll from "published" to "closed".
	// Publishes a PollClosed Kafka event.
	// Close(ctx context.Context, id string, creatorID string) (*dto.PollResponse, error)
}
