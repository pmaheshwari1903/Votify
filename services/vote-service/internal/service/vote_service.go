package service

// VoteService defines the interface for vote business logic.
type VoteService interface {
	// CastVote records a vote and publishes a VoteCreated Kafka event.
	// Business rules:
	// - Poll must exist and be in "published" status (cross-service check via HTTP or cache)
	// - Option must belong to the poll
	// - Duplicate vote prevention (based on voter ID or IP hash)
	// CastVote(ctx context.Context, req dto.CastVoteRequest, voterID string, ipHash string) (*dto.VoteResponse, error)

	// GetTally returns aggregated vote counts for a poll.
	// GetTally(ctx context.Context, pollID string) (*dto.TallyResponse, error)
}
