package realtime

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/redis/go-redis/v9"
	"github.com/votify/backend/internal/polls"
)

type RealtimeOption struct {
	ID        string `json:"id"`
	VoteCount int    `json:"voteCount"`
}

type RealtimeUpdateMessage struct {
	Type       string           `json:"type"`
	PollID     string           `json:"pollId"`
	TotalVotes int              `json:"totalVotes"`
	Options    []RealtimeOption `json:"options"`
}

type RealtimeService struct {
	mu          sync.RWMutex
	listeners   map[string][]func(msg RealtimeUpdateMessage)
	redisClient *redis.Client
	pollService polls.PollService
}

func NewRealtimeService(redisAddr, redisPassword string, redisDB int, pollService polls.PollService) *RealtimeService {
	var client *redis.Client
	if redisAddr != "" {
		client = redis.NewClient(&redis.Options{
			Addr:     redisAddr,
			Password: redisPassword,
			DB:       redisDB,
		})
	}

	return &RealtimeService{
		listeners:   make(map[string][]func(msg RealtimeUpdateMessage)),
		redisClient: client,
		pollService: pollService,
	}
}

func (s *RealtimeService) Subscribe(pollID string, callback func(msg RealtimeUpdateMessage)) func() {
	s.mu.Lock()
	s.listeners[pollID] = append(s.listeners[pollID], callback)
	firstListener := len(s.listeners[pollID]) == 1
	s.mu.Unlock()

	// If using Redis, listen on Redis channel
	var ctxCancel context.CancelFunc
	if s.redisClient != nil && firstListener {
		ctx, cancel := context.WithCancel(context.Background())
		ctxCancel = cancel
		pubsub := s.redisClient.Subscribe(ctx, "poll:results:"+pollID)

		go func() {
			defer pubsub.Close()
			ch := pubsub.Channel()
			for {
				select {
				case <-ctx.Done():
					return
				case msg, ok := <-ch:
					if !ok {
						return
					}
					var update RealtimeUpdateMessage
					if err := json.Unmarshal([]byte(msg.Payload), &update); err == nil {
						s.notifyLocalListeners(pollID, update)
					}
				}
			}
		}()
	}

	return func() {
		if ctxCancel != nil {
			ctxCancel()
		}
		s.mu.Lock()
		defer s.mu.Unlock()
		list := s.listeners[pollID]
		for i, fn := range list {
			// Compare pointer or remove
			_ = fn
			s.listeners[pollID] = append(list[:i], list[i+1:]...)
			break
		}
	}
}

func (s *RealtimeService) PublishUpdate(ctx context.Context, pollID string) {
	msg, err := s.BuildSnapshot(ctx, pollID)
	if err != nil {
		return
	}

	data, err := json.Marshal(msg)
	if err != nil {
		return
	}

	// Publish to Redis Pub/Sub if active
	if s.redisClient != nil {
		_ = s.redisClient.Publish(ctx, "poll:results:"+pollID, string(data)).Err()
	} else {
		// Local broadcast fallback when Redis is offline
		s.notifyLocalListeners(pollID, *msg)
	}
}

func (s *RealtimeService) notifyLocalListeners(pollID string, msg RealtimeUpdateMessage) {
	s.mu.RLock()
	listeners := append([]func(msg RealtimeUpdateMessage){}, s.listeners[pollID]...)
	s.mu.RUnlock()

	for _, fn := range listeners {
		go fn(msg)
	}
}

func (s *RealtimeService) BuildSnapshot(ctx context.Context, pollID string) (*RealtimeUpdateMessage, error) {
	results, err := s.pollService.GetPollResults(ctx, pollID)
	if err != nil {
		return nil, fmt.Errorf("fetch poll results: %w", err)
	}

	var opts []RealtimeOption
	for _, opt := range results.Options {
		opts = append(opts, RealtimeOption{
			ID:        opt.ID,
			VoteCount: opt.VoteCount,
		})
	}

	return &RealtimeUpdateMessage{
		Type:       "POLL_RESULTS_UPDATED",
		PollID:     results.PollID,
		TotalVotes: results.TotalVotes,
		Options:    opts,
	}, nil
}
