package repository

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/votify/vote-service/internal/model"
)

// VoteRepository defines data access methods for votes.
type VoteRepository interface {
	Create(ctx context.Context, vote *model.Vote) error
	HasVoted(ctx context.Context, pollID, voterIdentity string) (bool, error)
}

// MemoryVoteRepository is an in-memory thread-safe implementation of VoteRepository
// using compound unique key `(pollID, voterIdentity)` to prevent duplicate votes.
type MemoryVoteRepository struct {
	mu            sync.RWMutex
	votes         map[string]*model.Vote
	votedKeyMap   map[string]bool // key: "pollID:voterIdentity"
	counter       int64
}

func NewMemoryVoteRepository() *MemoryVoteRepository {
	return &MemoryVoteRepository{
		votes:       make(map[string]*model.Vote),
		votedKeyMap: make(map[string]bool),
	}
}

func (r *MemoryVoteRepository) HasVoted(ctx context.Context, pollID, voterIdentity string) (bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	key := fmt.Sprintf("%s:%s", pollID, voterIdentity)
	return r.votedKeyMap[key], nil
}

func (r *MemoryVoteRepository) Create(ctx context.Context, vote *model.Vote) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	key := fmt.Sprintf("%s:%s", vote.PollID, vote.VoterIdentity)
	if r.votedKeyMap[key] {
		return fmt.Errorf("duplicate vote detected for poll '%s' and voter '%s'", vote.PollID, vote.VoterIdentity)
	}

	r.counter++
	vote.ID = fmt.Sprintf("vt_%d_%d", time.Now().UnixNano(), r.counter)
	vote.CreatedAt = time.Now().UTC()

	r.votedKeyMap[key] = true
	copy := *vote
	r.votes[vote.ID] = &copy

	return nil
}
