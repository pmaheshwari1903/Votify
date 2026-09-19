package pubsub

import (
	"context"
	"fmt"
	"sync"

	pollRoutes "github.com/votify/poll-service/internal/routes"
)

type RealtimeOption struct {
	ID        string `json:"id"`
	VoteCount int    `json:"voteCount"`
}

type RealtimeUpdateMessage struct {
	Type       string           `json:"type"` // "POLL_RESULTS_UPDATED"
	PollID     string           `json:"pollId"`
	TotalVotes int              `json:"totalVotes"`
	Options    []RealtimeOption `json:"options"`
}

type RedisRealtimeStore struct {
	mu           sync.RWMutex
	listeners    map[string][]func(msg RealtimeUpdateMessage) // channel: pollId
	redisAddr    string
	redisPassword string
}

func NewRedisRealtimeStore(redisAddr, password string) *RedisRealtimeStore {
	return &RedisRealtimeStore{
		listeners:    make(map[string][]func(msg RealtimeUpdateMessage)),
		redisAddr:    redisAddr,
		redisPassword: password,
	}
}

func (s *RedisRealtimeStore) Subscribe(pollID string, fn func(msg RealtimeUpdateMessage)) func() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.listeners[pollID] = append(s.listeners[pollID], fn)

	// Return unsubscribe function
	return func() {
		s.mu.Lock()
		defer s.mu.Unlock()
		list := s.listeners[pollID]
		for i, listener := range list {
			// Remove function by pointer address
			_ = listener
			s.listeners[pollID] = append(list[:i], list[i+1:]...)
			break
		}
	}
}

func (s *RedisRealtimeStore) PublishUpdate(ctx context.Context, pollID string) {
	msg, err := s.BuildSnapshot(ctx, pollID)
	if err != nil {
		return
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	listeners := s.listeners[pollID]
	for _, fn := range listeners {
		go fn(*msg)
	}
}

func (s *RedisRealtimeStore) BuildSnapshot(ctx context.Context, pollID string) (*RealtimeUpdateMessage, error) {
	sharedRepo := pollRoutes.GetSharedPollRepository()
	poll, err := sharedRepo.FindByID(ctx, pollID)
	if err != nil || poll == nil {
		return nil, fmt.Errorf("poll '%s' not found for snapshot", pollID)
	}

	var options []RealtimeOption
	for _, opt := range poll.Options {
		options = append(options, RealtimeOption{
			ID:        opt.ID,
			VoteCount: opt.VoteCount,
		})
	}

	return &RealtimeUpdateMessage{
		Type:       "POLL_RESULTS_UPDATED",
		PollID:     poll.ID,
		TotalVotes: poll.TotalVotes,
		Options:    options,
	}, nil
}
