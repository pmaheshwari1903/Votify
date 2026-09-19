package pubsub

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
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

type RedisRealtimeStore struct {
	mu             sync.RWMutex
	listeners      map[string][]func(msg RealtimeUpdateMessage)
	redisAddr      string
	redisPassword  string
	pollServiceURL string
}

func NewRedisRealtimeStore(
	redisAddr string,
	password string,
	pollServiceURL string,
) *RedisRealtimeStore {
	return &RedisRealtimeStore{
		listeners:      make(map[string][]func(msg RealtimeUpdateMessage)),
		redisAddr:      redisAddr,
		redisPassword:  password,
		pollServiceURL: pollServiceURL,
	}
}

func (s *RedisRealtimeStore) Subscribe(
	pollID string,
	fn func(msg RealtimeUpdateMessage),
) func() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.listeners[pollID] = append(s.listeners[pollID], fn)

	return func() {
		s.mu.Lock()
		defer s.mu.Unlock()

		list := s.listeners[pollID]

		for i := range list {
			_ = list[i]
			s.listeners[pollID] = append(list[:i], list[i+1:]...)
			break
		}
	}
}

func (s *RedisRealtimeStore) PublishUpdate(
	ctx context.Context,
	pollID string,
) {
	msg, err := s.BuildSnapshot(ctx, pollID)
	if err != nil {
		return
	}

	s.mu.RLock()
	listeners := append(
		[]func(msg RealtimeUpdateMessage){},
		s.listeners[pollID]...,
	)
	s.mu.RUnlock()

	for _, fn := range listeners {
		go fn(*msg)
	}
}

func (s *RedisRealtimeStore) BuildSnapshot(
	ctx context.Context,
	pollID string,
) (*RealtimeUpdateMessage, error) {

	url := fmt.Sprintf(
		"%s/polls/public/%s",
		strings.TrimRight(s.pollServiceURL, "/"),
		pollID,
	)

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		url,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("create poll request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetch poll from poll service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(
			"poll service returned status %d",
			resp.StatusCode,
		)
	}

	var apiResp struct {
		Success bool `json:"success"`
		Data    struct {
			ID      string `json:"id"`
			Options []struct {
				ID        string `json:"id"`
				VoteCount int    `json:"voteCount"`
			} `json:"options"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, fmt.Errorf("decode poll response: %w", err)
	}

	if !apiResp.Success {
		return nil, fmt.Errorf("poll service returned unsuccessful response")
	}

	var options []RealtimeOption
	totalVotes := 0

	for _, opt := range apiResp.Data.Options {
		options = append(options, RealtimeOption{
			ID:        opt.ID,
			VoteCount: opt.VoteCount,
		})

		totalVotes += opt.VoteCount
	}

	return &RealtimeUpdateMessage{
		Type:       "POLL_RESULTS_UPDATED",
		PollID:     apiResp.Data.ID,
		TotalVotes: totalVotes,
		Options:    options,
	}, nil
}