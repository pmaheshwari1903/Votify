package repository

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/votify/poll-service/internal/model"
)

// PollRepository defines data access methods for polls.
type PollRepository interface {
	Create(ctx context.Context, poll *model.Poll) error
	FindByID(ctx context.Context, id string) (*model.Poll, error)
	FindByOwnerID(ctx context.Context, ownerID string) ([]*model.Poll, error)
	Update(ctx context.Context, poll *model.Poll) error
	Delete(ctx context.Context, id string) error
	IncrementVote(ctx context.Context, pollID string, optionID string) error
}

// MemoryPollRepository is an in-memory thread-safe implementation of PollRepository.
type MemoryPollRepository struct {
	mu      sync.RWMutex
	polls   map[string]*model.Poll
	counter int64
}

func NewMemoryPollRepository() *MemoryPollRepository {
	return &MemoryPollRepository{
		polls: make(map[string]*model.Poll),
	}
}

func (r *MemoryPollRepository) Create(ctx context.Context, poll *model.Poll) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.counter++
	poll.ID = fmt.Sprintf("poll_%d_%d", time.Now().UnixNano(), r.counter)
	poll.CreatedAt = time.Now().UTC()
	poll.UpdatedAt = time.Now().UTC()

	copy := *poll
	r.polls[poll.ID] = &copy
	return nil
}

func (r *MemoryPollRepository) FindByID(ctx context.Context, id string) (*model.Poll, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	poll, exists := r.polls[id]
	if !exists {
		return nil, nil
	}
	copy := *poll
	return &copy, nil
}

func (r *MemoryPollRepository) FindByOwnerID(ctx context.Context, ownerID string) ([]*model.Poll, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []*model.Poll
	for _, poll := range r.polls {
		if poll.OwnerID == ownerID {
			copy := *poll
			result = append(result, &copy)
		}
	}
	return result, nil
}

func (r *MemoryPollRepository) Update(ctx context.Context, poll *model.Poll) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.polls[poll.ID]; !exists {
		return fmt.Errorf("poll '%s' not found", poll.ID)
	}

	poll.UpdatedAt = time.Now().UTC()
	copy := *poll
	r.polls[poll.ID] = &copy
	return nil
}

func (r *MemoryPollRepository) Delete(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.polls[id]; !exists {
		return fmt.Errorf("poll '%s' not found", id)
	}

	delete(r.polls, id)
	return nil
}

func (r *MemoryPollRepository) IncrementVote(ctx context.Context, pollID string, optionID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	poll, exists := r.polls[pollID]
	if !exists {
		return fmt.Errorf("poll '%s' not found", pollID)
	}

	optionFound := false
	for i := range poll.Options {
		if poll.Options[i].ID == optionID {
			poll.Options[i].VoteCount++
			poll.TotalVotes++
			optionFound = true
			break
		}
	}

	if !optionFound {
		return fmt.Errorf("option '%s' not found in poll '%s'", optionID, pollID)
	}

	poll.UpdatedAt = time.Now().UTC()
	return nil
}
