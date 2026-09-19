package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/votify/pkg/errors"
	"github.com/votify/pkg/events"
	pollDto "github.com/votify/poll-service/internal/dto"
	pollRoutes "github.com/votify/poll-service/internal/routes"

	"github.com/votify/vote-service/internal/broker"
	voteDto "github.com/votify/vote-service/internal/dto"
	"github.com/votify/vote-service/internal/model"
	"github.com/votify/vote-service/internal/repository"
)

type VoteService interface {
	CastVote(ctx context.Context, userID, clientIP string, req voteDto.CastVoteRequest) (*voteDto.VoteResponse, error)
}

type DefaultVoteService struct {
	repo           repository.VoteRepository
	producer       *broker.KafkaProducer
	pollServiceURL string
}

func NewVoteService(repo repository.VoteRepository, producer *broker.KafkaProducer, pollServiceURL string) *DefaultVoteService {
	return &DefaultVoteService{
		repo:           repo,
		producer:       producer,
		pollServiceURL: pollServiceURL,
	}
}

func (s *DefaultVoteService) CastVote(ctx context.Context, userID, clientIP string, req voteDto.CastVoteRequest) (*voteDto.VoteResponse, error) {
	pollID := strings.TrimSpace(req.PollID)
	optionID := strings.TrimSpace(req.OptionID)

	if pollID == "" || optionID == "" {
		return nil, errors.NewBadRequestError("pollId and optionId are required", "MISSING_REQUIRED_FIELDS")
	}

	// 1. Fetch Poll state authoritatively
	pubPoll, err := s.fetchPollState(ctx, pollID)
	if err != nil {
		return nil, err
	}

	// 2. Verify Poll status is OPEN
	if pubPoll.Status != "open" {
		return nil, errors.NewConflictError("This poll is closed for voting", "POLL_CLOSED")
	}

	// 3. Verify option belongs to this poll
	optionValid := false
	for _, opt := range pubPoll.Options {
		if opt.ID == optionID {
			optionValid = true
			break
		}
	}
	if !optionValid {
		return nil, errors.NewBadRequestError(fmt.Sprintf("Option '%s' does not exist in poll '%s'", optionID, pollID), "INVALID_OPTION_ID")
	}

	// 4. Determine voter identity for duplicate protection
	voterIdentity := userID
	if voterIdentity == "" {
		voterIdentity = "ip_" + clientIP
	}

	// 5. Check duplicate vote
	alreadyVoted, err := s.repo.HasVoted(ctx, pollID, voterIdentity)
	if err != nil {
		return nil, errors.NewInternalError(err)
	}
	if alreadyVoted {
		return nil, errors.NewConflictError("You have already voted on this poll", "DUPLICATE_VOTE")
	}

	// 6. Record vote in MongoDB / Vote Repository
	vote := &model.Vote{
		PollID:        pollID,
		OptionID:      optionID,
		VoterIdentity: voterIdentity,
		VoterIP:       clientIP,
		UserID:        userID,
	}

	if err := s.repo.Create(ctx, vote); err != nil {
		if strings.Contains(err.Error(), "duplicate vote") {
			return nil, errors.NewConflictError("You have already voted on this poll", "DUPLICATE_VOTE")
		}
		return nil, errors.NewInternalError(err)
	}

	// 7. Increment vote count in Poll Service
	s.incrementPollVote(ctx, pollID, optionID)

	// 8. Publish VoteCreated event to Kafka AFTER successful database persistence
	if s.producer != nil {
		_ = s.producer.PublishVoteCreated(ctx, events.VoteCreatedPayload{
			PollID:    pollID,
			OptionID:  optionID,
			VoterIP:   clientIP,
			UserID:    userID,
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		})
	}

	return &voteDto.VoteResponse{
		ID:        vote.ID,
		PollID:    vote.PollID,
		OptionID:  vote.OptionID,
		CreatedAt: vote.CreatedAt,
	}, nil
}

func (s *DefaultVoteService) fetchPollState(ctx context.Context, pollID string) (*dto.PublicPollResponse, error) {
	sharedPollRepo := pollRoutes.GetSharedPollRepository()
	p, _ := sharedPollRepo.FindByID(ctx, pollID)
	if p != nil {
		return &dto.PublicPollResponse{
			ID:          p.ID,
			Question:    p.Question,
			Description: p.Description,
			Options:     p.Options,
			Status:      p.Status,
			CreatedAt:   p.CreatedAt,
		}, nil
	}

	url := fmt.Sprintf("%s/polls/public/%s", s.pollServiceURL, pollID)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, errors.NewInternalError(err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		return nil, errors.NewNotFoundError("Poll", pollID)
	}
	defer resp.Body.Close()

	var apiResp struct {
		Success bool                   `json:"success"`
		Data    dto.PublicPollResponse `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil || !apiResp.Success {
		return nil, errors.NewNotFoundError("Poll", pollID)
	}

	return &apiResp.Data, nil
}

func (s *DefaultVoteService) incrementPollVote(ctx context.Context, pollID, optionID string) {
	sharedPollRepo := pollRoutes.GetSharedPollRepository()
	_ = sharedPollRepo.IncrementVote(ctx, pollID, optionID)

	url := fmt.Sprintf("%s/internal/polls/%s/options/%s/vote", s.pollServiceURL, pollID, optionID)
	req, err := http.NewRequestWithContext(ctx, "POST", url, nil)
	if err == nil {
		resp, err := http.DefaultClient.Do(req)
		if err == nil && resp != nil {
			resp.Body.Close()
		}
	}
}
